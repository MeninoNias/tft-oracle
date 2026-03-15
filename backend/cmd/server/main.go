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
	"github.com/MeninoNias/tft-oracle/backend/internal/cdragon"
	"github.com/MeninoNias/tft-oracle/backend/internal/config"
	"github.com/MeninoNias/tft-oracle/backend/internal/database"
	"github.com/MeninoNias/tft-oracle/backend/internal/patch"
	"github.com/MeninoNias/tft-oracle/backend/internal/player"
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

	// Set up Connect RPC handlers
	mux := http.NewServeMux()

	interceptors := connect.WithInterceptors(newLoggingInterceptor())

	patchPath, patchHandler := tftv1connect.NewPatchServiceHandler(
		patch.NewService(pool),
		interceptors,
	)
	mux.Handle(patchPath, patchHandler)

	playerPath, playerHandler := tftv1connect.NewPlayerServiceHandler(
		player.NewService(),
		interceptors,
	)
	mux.Handle(playerPath, playerHandler)

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

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
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
