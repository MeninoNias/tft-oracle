package crawler

import "time"

// RawComposition represents a single composition scraped from a tier list site.
type RawComposition struct {
	Name         string            `json:"name"`
	Tier         string            `json:"tier"`
	WinRate      float64           `json:"win_rate"`
	PlayRate     float64           `json:"play_rate"`
	AvgPlacement float64           `json:"avg_placement"`
	ChampionIDs  []string          `json:"champion_ids"`  // api_names
	CoreItems    map[string][]string `json:"core_items"`   // champion_api_name -> [item_api_names]
}

// ScrapeResult holds the output of a single scrape operation.
type ScrapeResult struct {
	Source       string           `json:"source"`
	Patch        string           `json:"patch"`
	Compositions []RawComposition `json:"compositions"`
	ScrapedAt    time.Time        `json:"scraped_at"`
}

// SourceStatusInfo tracks the health of each scraper.
type SourceStatusInfo struct {
	Source       string
	LastScraped  time.Time
	Status       string // "ok", "error", "never"
	ErrorMessage string
}
