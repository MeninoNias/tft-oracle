package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/MeninoNias/tft-oracle/backend/gen/tft/v1/tftv1connect"
	"github.com/MeninoNias/tft-oracle/backend/internal/ai"
	"github.com/MeninoNias/tft-oracle/backend/internal/auth"
	"github.com/MeninoNias/tft-oracle/backend/internal/coach"
	"github.com/MeninoNias/tft-oracle/backend/internal/cache"
	"github.com/MeninoNias/tft-oracle/backend/internal/cdragon"
	"github.com/MeninoNias/tft-oracle/backend/internal/config"
	"github.com/MeninoNias/tft-oracle/backend/internal/consolidation"
	"github.com/MeninoNias/tft-oracle/backend/internal/crawler"
	"github.com/MeninoNias/tft-oracle/backend/internal/database"
	"github.com/MeninoNias/tft-oracle/backend/internal/patch"
	"github.com/MeninoNias/tft-oracle/backend/internal/player"
	"github.com/MeninoNias/tft-oracle/backend/internal/riot"
	"github.com/MeninoNias/tft-oracle/backend/internal/simulation"
	"github.com/MeninoNias/tft-oracle/backend/internal/tierlist"
)

func main() {
	syncFlag := flag.Bool("sync", false, "Force sync CommunityDragon data on startup")
	flag.Parse()

	// Load .env file (optional, for development)
	_ = godotenv.Load()
	_ = godotenv.Load(filepath.Join("..", ".env"))

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run database migrations
	migrationsPath := filepath.Join("..", "migrations")
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		// Try from project root (when running from backend/)
		migrationsPath = "migrations"
	}
	log.Println("running migrations...")
	if err := database.RunMigrations(cfg.DatabaseURL, migrationsPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Connect to database
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// CommunityDragon sync
	syncer := cdragon.NewSyncer(pool)
	if *syncFlag {
		if err := syncer.Sync(ctx, "en_us"); err != nil {
			log.Fatalf("cdragon sync failed: %v", err)
		}
	} else {
		if err := syncer.SyncIfEmpty(ctx, "en_us"); err != nil {
			log.Printf("warning: cdragon sync failed: %v", err)
		}
	}

	// Redis (optional)
	var cacheClient *cache.Client
	if cfg.RedisURL != "" {
		var err error
		cacheClient, err = cache.NewClient(cfg.RedisURL)
		if err != nil {
			log.Printf("warning: redis connection failed: %v (continuing without cache)", err)
		} else if cacheClient != nil {
			log.Println("redis: connected")
			defer cacheClient.Close()
		}
	}

	// Riot API client
	riotClient := riot.NewClient(cfg.RiotAPIKey)
	if riotClient.Available() {
		if err := riotClient.HealthCheck(ctx); err != nil {
			log.Printf("riot api: WARNING — key validation failed: %v", err)
		} else {
			log.Println("riot api: key valid")
		}
	} else {
		log.Println("riot api: not configured (player features disabled)")
	}

	// Auth (optional — works without JWT_SECRET, but auth endpoints return errors)
	var jwtMgr *auth.JWTManager
	if cfg.JWTSecret != "" {
		var err error
		jwtMgr, err = auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)
		if err != nil {
			log.Fatalf("failed to init JWT: %v", err)
		}
		log.Println("auth: JWT configured")
	} else {
		log.Println("auth: JWT_SECRET not set — auth features disabled")
	}

	// OpenAI client (optional — works without key)
	aiClient := ai.NewClient(cfg.OpenAIAPIKey)
	if aiClient.Available() {
		log.Println("openai: configured")
	} else {
		log.Println("openai: OPENAI_API_KEY not set — simulation features disabled")
	}

	// Crawler & consolidation (Phase 4)
	consolidationEngine := consolidation.NewEngine(pool, nil)
	httpClient := crawler.NewHTTPClient()
	scrapers := []crawler.Scraper{
		crawler.NewMobalyticsScraper(httpClient),
		crawler.NewTacticsToolsScraper(httpClient),
		crawler.NewMetaTFTScraper(httpClient),
	}

	crawlerInterval := 24 * time.Hour
	if d, err := time.ParseDuration(cfg.CrawlerInterval); err == nil {
		crawlerInterval = d
	}

	crawl := crawler.New(pool, scrapers, crawlerInterval, consolidationEngine.Consolidate)
	if cfg.CrawlerEnabled {
		go crawl.StartScheduler(ctx)
		log.Println("crawler: enabled")
	} else {
		log.Println("crawler: disabled")
	}

	// Set up Connect RPC handlers
	mux := http.NewServeMux()

	interceptors := connect.WithInterceptors(
		auth.NewAuthInterceptor(jwtMgr),
		newLoggingInterceptor(),
	)

	patchPath, patchHandler := tftv1connect.NewPatchServiceHandler(
		patch.NewService(pool),
		interceptors,
	)
	mux.Handle(patchPath, patchHandler)

	playerPath, playerHandler := tftv1connect.NewPlayerServiceHandler(
		player.NewService(pool, riotClient, cacheClient),
		interceptors,
	)
	mux.Handle(playerPath, playerHandler)

	authPath, authHandler := tftv1connect.NewAuthServiceHandler(
		auth.NewService(pool, riotClient, jwtMgr),
		interceptors,
	)
	mux.Handle(authPath, authHandler)

	simPath, simHandler := tftv1connect.NewSimulationServiceHandler(
		simulation.NewService(pool, aiClient),
		interceptors,
	)
	mux.Handle(simPath, simHandler)

	tierListPath, tierListHandler := tftv1connect.NewTierListServiceHandler(
		tierlist.NewService(pool, crawl),
		interceptors,
	)
	mux.Handle(tierListPath, tierListHandler)

	coachPath, coachHandler := tftv1connect.NewCoachServiceHandler(
		coach.NewService(pool, aiClient),
		interceptors,
	)
	mux.Handle(coachPath, coachHandler)

	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://localhost:1420",
			"https://tauri.localhost",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
			"Connect-Protocol-Version",
			"Connect-Timeout-Ms",
			"Grpc-Timeout",
			"X-Grpc-Web",
			"X-User-Agent",
		},
		ExposedHeaders: []string{
			"Grpc-Status",
			"Grpc-Message",
			"Grpc-Status-Details-Bin",
		},
		AllowCredentials: false,
	})

	server := &http.Server{
		Addr:         cfg.ListenAddr(),
		Handler:      corsHandler.Handler(h2c.NewHandler(mux, &http2.Server{})),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Graceful shutdown (os.Interrupt works reliably on Windows)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("shutting down...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	log.Printf("server listening on %s", cfg.ListenAddr())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

// loggingInterceptor logs Connect RPC calls.
type loggingInterceptor struct{}

func newLoggingInterceptor() *loggingInterceptor {
	return &loggingInterceptor{}
}

func (i *loggingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		start := time.Now()
		resp, err := next(ctx, req)
		log.Printf("%s %s (%v)", req.Spec().Procedure, statusFromError(err), time.Since(start))
		return resp, err
	}
}

func (i *loggingInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *loggingInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

func statusFromError(err error) string {
	if err == nil {
		return "ok"
	}
	if connectErr, ok := err.(*connect.Error); ok {
		return connectErr.Code().String()
	}
	return "unknown"
}
