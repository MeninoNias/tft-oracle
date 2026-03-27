package crawler

import (
	"context"
	"fmt"
)

const metatftSource = "metatft"

// MetaTFTScraper is a placeholder for MetaTFT scraping.
// MetaTFT is a client-side SPA with no public API — scraping requires
// a headless browser (Puppeteer/Chromedp), which is planned for a future iteration.
type MetaTFTScraper struct {
	client *HTTPClient
}

// NewMetaTFTScraper creates a MetaTFT scraper.
func NewMetaTFTScraper(client *HTTPClient) *MetaTFTScraper {
	return &MetaTFTScraper{client: client}
}

func (s *MetaTFTScraper) Name() string { return metatftSource }

func (s *MetaTFTScraper) Scrape(ctx context.Context) (*ScrapeResult, error) {
	return nil, fmt.Errorf("metatft scraper not yet implemented — requires headless browser")
}
