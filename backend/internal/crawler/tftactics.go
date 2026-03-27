package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

const (
	tacticsToolsSource = "tacticstools"
	tacticsToolsURL    = "https://tactics.tools/team-compositions"
)

// TacticsToolsScraper scrapes composition data from tactics.tools via __NEXT_DATA__.
type TacticsToolsScraper struct {
	client *HTTPClient
}

// NewTacticsToolsScraper creates a tactics.tools scraper.
func NewTacticsToolsScraper(client *HTTPClient) *TacticsToolsScraper {
	return &TacticsToolsScraper{client: client}
}

func (s *TacticsToolsScraper) Name() string { return tacticsToolsSource }

func (s *TacticsToolsScraper) Scrape(ctx context.Context) (*ScrapeResult, error) {
	body, err := s.client.Get(ctx, tacticsToolsURL)
	if err != nil {
		return nil, fmt.Errorf("fetch tactics.tools: %w", err)
	}

	comps, patch, err := parseTacticsToolsHTML(body)
	if err != nil {
		return nil, fmt.Errorf("parse tactics.tools: %w", err)
	}

	if len(comps) == 0 {
		return nil, fmt.Errorf("tactics.tools returned 0 compositions")
	}

	return &ScrapeResult{
		Source:       tacticsToolsSource,
		Patch:        patch,
		Compositions: comps,
		ScrapedAt:    time.Now(),
	}, nil
}

var nextDataRe = regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)

func parseTacticsToolsHTML(html []byte) ([]RawComposition, string, error) {
	matches := nextDataRe.FindSubmatch(html)
	if len(matches) < 2 {
		return nil, "", fmt.Errorf("__NEXT_DATA__ not found in HTML")
	}

	var nextData tacticsToolsNextData
	if err := json.Unmarshal(matches[1], &nextData); err != nil {
		return nil, "", fmt.Errorf("unmarshal __NEXT_DATA__: %w", err)
	}

	data := nextData.Props.PageProps.InitialData
	if data == nil {
		return nil, "", fmt.Errorf("no initialData in page props")
	}

	// Determine patch from aperture
	patch := "unknown"

	// Process compositions from nonSpatComps (main comp list)
	rawComps := data.NonSpatComps
	if len(rawComps) == 0 && len(data.Groups) > 0 {
		// Fall back to groups
		for _, g := range data.Groups {
			rawComps = append(rawComps, g.Full.Comps...)
		}
	}

	comps := make([]RawComposition, 0, len(rawComps))
	for _, rc := range rawComps {
		if len(rc.Units) == 0 {
			continue
		}

		// Validate data
		if rc.Place < 1 || rc.Place > 8 {
			continue
		}

		// Calculate win rate and play rate
		winRate := 0.0
		if rc.Count > 0 {
			winRate = float64(rc.Win) / float64(rc.Count) * 100
		}

		top4Rate := 0.0
		if rc.Count > 0 {
			top4Rate = float64(rc.Top4) / float64(rc.Count) * 100
		}

		// Assign tier based on average placement
		tier := placementToTier(rc.Place)

		// Name from champion list (no names in data)
		name := buildCompName(rc.Units)

		comps = append(comps, RawComposition{
			Name:         name,
			Tier:         tier,
			WinRate:      winRate,
			PlayRate:     top4Rate, // Using top4 rate as "play rate" metric
			AvgPlacement: rc.Place,
			ChampionIDs:  rc.Units,
			CoreItems:    make(map[string][]string),
		})
	}

	return comps, patch, nil
}

// placementToTier maps average placement to a tier letter.
func placementToTier(avgPlace float64) string {
	switch {
	case avgPlace <= 3.2:
		return "S"
	case avgPlace <= 3.8:
		return "A"
	case avgPlace <= 4.3:
		return "B"
	default:
		return "C"
	}
}

// buildCompName creates a human-readable name from champion IDs.
func buildCompName(units []string) string {
	if len(units) == 0 {
		return "Unknown"
	}
	// Take first 2-3 champs, strip the TFT16_ prefix
	names := make([]string, 0, 3)
	for i, u := range units {
		if i >= 3 {
			break
		}
		name := u
		if idx := len("TFT16_"); len(u) > idx && u[:idx] == "TFT16_" {
			name = u[idx:]
		}
		names = append(names, name)
	}
	result := ""
	for i, n := range names {
		if i > 0 {
			result += " "
		}
		result += n
	}
	return result
}

// --- __NEXT_DATA__ response types ---

type tacticsToolsNextData struct {
	Props struct {
		PageProps struct {
			InitialData *tacticsToolsData `json:"initialData"`
		} `json:"pageProps"`
	} `json:"props"`
}

type tacticsToolsData struct {
	Count        int                    `json:"count"`
	Place        float64                `json:"place"`
	Top4         int                    `json:"top4"`
	Win          int                    `json:"win"`
	NonSpatComps []tacticsToolsComp     `json:"nonSpatComps"`
	Groups       []tacticsToolsGroup    `json:"groups"`
}

type tacticsToolsGroup struct {
	Full struct {
		Comps []tacticsToolsComp `json:"comps"`
	} `json:"full"`
}

type tacticsToolsComp struct {
	Units       []string `json:"units"`
	SpatItems   []string `json:"spatItems"`
	ExtraTraits []string `json:"extraTraits"`
	Count       int      `json:"count"`
	Place       float64  `json:"place"`
	Top4        int      `json:"top4"`
	Win         int      `json:"win"`
}
