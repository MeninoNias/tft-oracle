package crawler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

// Scraper defines the contract for a single-site scraper.
type Scraper interface {
	// Name returns the source identifier (e.g. "metatft").
	Name() string
	// Scrape fetches and parses tier list data from the source.
	Scrape(ctx context.Context) (*ScrapeResult, error)
}

// OnCompleteFunc is called after a crawl cycle completes with the patch version.
type OnCompleteFunc func(ctx context.Context, patch string) error

// Crawler orchestrates multiple scrapers and stores results.
type Crawler struct {
	scrapers   []Scraper
	store      *Store
	interval   time.Duration
	onComplete OnCompleteFunc

	mu       sync.Mutex
	statuses map[string]*SourceStatusInfo
}

// New creates a new Crawler.
func New(db *pgxpool.Pool, scrapers []Scraper, interval time.Duration, onComplete OnCompleteFunc) *Crawler {
	statuses := make(map[string]*SourceStatusInfo, len(scrapers))
	for _, s := range scrapers {
		statuses[s.Name()] = &SourceStatusInfo{
			Source: s.Name(),
			Status: "never",
		}
	}

	return &Crawler{
		scrapers:   scrapers,
		store:      NewStore(db),
		interval:   interval,
		onComplete: onComplete,
		statuses:   statuses,
	}
}

// Run executes a single crawl cycle: scrapes all sources in parallel, stores results,
// then triggers consolidation.
func (c *Crawler) Run(ctx context.Context) error {
	log.Println("crawler: starting crawl cycle...")
	start := time.Now()

	type result struct {
		scrapeResult *ScrapeResult
		err          error
		name         string
	}

	results := make([]result, len(c.scrapers))
	g, gCtx := errgroup.WithContext(ctx)

	for i, scraper := range c.scrapers {
		g.Go(func() error {
			log.Printf("crawler: scraping %s...", scraper.Name())
			sr, err := scraper.Scrape(gCtx)
			results[i] = result{scrapeResult: sr, err: err, name: scraper.Name()}
			// Don't propagate errors — one failure shouldn't block others
			return nil
		})
	}

	_ = g.Wait()

	// Process results
	var (
		successCount int
		lastPatch    string
	)

	for _, r := range results {
		c.mu.Lock()
		status := c.statuses[r.name]

		if r.err != nil {
			log.Printf("crawler: %s failed: %v", r.name, r.err)
			status.Status = "error"
			status.ErrorMessage = r.err.Error()
			c.mu.Unlock()
			continue
		}

		// Store the scraped data
		if err := c.store.SaveResult(ctx, r.scrapeResult); err != nil {
			log.Printf("crawler: failed to store %s results: %v", r.name, err)
			status.Status = "error"
			status.ErrorMessage = fmt.Sprintf("store failed: %v", err)
			c.mu.Unlock()
			continue
		}

		status.Status = "ok"
		status.LastScraped = r.scrapeResult.ScrapedAt
		status.ErrorMessage = ""
		c.mu.Unlock()

		successCount++
		if r.scrapeResult.Patch != "" && r.scrapeResult.Patch != "unknown" {
			lastPatch = r.scrapeResult.Patch
		} else if lastPatch == "" {
			lastPatch = r.scrapeResult.Patch
		}
	}

	// Normalize patches — use the best known patch for all "unknown" results
	if lastPatch != "" && lastPatch != "unknown" {
		for _, r := range results {
			if r.scrapeResult != nil && (r.scrapeResult.Patch == "" || r.scrapeResult.Patch == "unknown") {
				r.scrapeResult.Patch = lastPatch
				// Re-store with correct patch
				_ = c.store.SaveResult(ctx, r.scrapeResult)
			}
		}
	}

	log.Printf("crawler: cycle complete — %d/%d sources succeeded (%v)",
		successCount, len(c.scrapers), time.Since(start))

	// Trigger consolidation if we have any data
	if successCount > 0 && lastPatch != "" && c.onComplete != nil {
		log.Printf("crawler: triggering consolidation for patch %s...", lastPatch)
		if err := c.onComplete(ctx, lastPatch); err != nil {
			log.Printf("crawler: consolidation failed: %v", err)
			return fmt.Errorf("consolidation: %w", err)
		}
	}

	return nil
}

// StartScheduler runs the crawler on a periodic schedule.
// It blocks until the context is cancelled.
func (c *Crawler) StartScheduler(ctx context.Context) {
	log.Printf("crawler: scheduler started (interval: %v)", c.interval)

	// Run immediately on first start if data is stale
	if c.shouldRunNow(ctx) {
		if err := c.Run(ctx); err != nil {
			log.Printf("crawler: initial run failed: %v", err)
		}
	}

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("crawler: scheduler stopped")
			return
		case <-ticker.C:
			if err := c.Run(ctx); err != nil {
				log.Printf("crawler: scheduled run failed: %v", err)
			}
		}
	}
}

// shouldRunNow checks if the last scrape is older than the interval.
func (c *Crawler) shouldRunNow(ctx context.Context) bool {
	latest, err := c.store.queries.GetLatestScrapedAt(ctx)
	if err != nil || latest == nil {
		return true // No data yet
	}

	if t, ok := latest.(time.Time); ok {
		return time.Since(t) > c.interval
	}
	return true
}

// GetStatuses returns the current scraper statuses.
func (c *Crawler) GetStatuses() []SourceStatusInfo {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]SourceStatusInfo, 0, len(c.statuses))
	for _, s := range c.statuses {
		out = append(out, *s)
	}
	return out
}

// --- Shared helpers ---

// normalizeTier converts tier strings to a consistent uppercase format.
func normalizeTier(tier string) string {
	t := strings.ToUpper(strings.TrimSpace(tier))
	if t == "S+" {
		return "S"
	}
	return t
}
