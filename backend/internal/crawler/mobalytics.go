package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	mobalyticsSource     = "mobalytics"
	mobalyticsGraphQLURL = "https://mobalytics.gg/api/tft/v1/graphql/query"
)

// MobalyticsScraper scrapes composition tier data from Mobalytics via GraphQL.
type MobalyticsScraper struct {
	client *HTTPClient
}

// NewMobalyticsScraper creates a Mobalytics scraper.
func NewMobalyticsScraper(client *HTTPClient) *MobalyticsScraper {
	return &MobalyticsScraper{client: client}
}

func (s *MobalyticsScraper) Name() string { return mobalyticsSource }

func (s *MobalyticsScraper) Scrape(ctx context.Context) (*ScrapeResult, error) {
	var allComps []RawComposition
	var patch string

	for _, tier := range []string{"s", "a", "b", "c"} {
		comps, p, err := s.scrapeTier(ctx, tier)
		if err != nil {
			return nil, fmt.Errorf("scrape tier %s: %w", tier, err)
		}
		allComps = append(allComps, comps...)
		if p != "" {
			patch = p
		}
	}

	if len(allComps) == 0 {
		return nil, fmt.Errorf("mobalytics returned 0 compositions")
	}

	return &ScrapeResult{
		Source:       mobalyticsSource,
		Patch:        patch,
		Compositions: allComps,
		ScrapedAt:    time.Now(),
	}, nil
}

func (s *MobalyticsScraper) scrapeTier(ctx context.Context, tier string) ([]RawComposition, string, error) {
	query := `{
		tft {
			metaCompositions(filter: {set: "16", tier: "` + tier + `"}, first: 50) {
				items {
					name
					tier
					patch
					slug
					formation {
						positions {
							champion {
								champion { slug }
								items { slug }
								level
							}
						}
					}
				}
			}
		}
	}`

	body, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return nil, "", fmt.Errorf("marshal query: %w", err)
	}

	respBody, err := s.client.Post(ctx, mobalyticsGraphQLURL, body)
	if err != nil {
		return nil, "", fmt.Errorf("graphql request: %w", err)
	}

	var resp mobalyticsGQLResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, "", fmt.Errorf("unmarshal response: %w", err)
	}

	items := resp.Data.TFT.MetaCompositions.Items
	if len(items) == 0 {
		return nil, "", nil
	}

	patch := items[0].Patch
	comps := make([]RawComposition, 0, len(items))

	for _, item := range items {
		champIDs := make([]string, 0)
		coreItems := make(map[string][]string)

		for _, pos := range item.Formation.Positions {
			slug := pos.Champion.Champion.Slug
			apiName := slugToApiName(slug)
			champIDs = append(champIDs, apiName)

			if len(pos.Champion.Items) > 0 {
				itemSlugs := make([]string, 0, len(pos.Champion.Items))
				for _, it := range pos.Champion.Items {
					itemSlugs = append(itemSlugs, it.Slug)
				}
				coreItems[apiName] = itemSlugs
			}
		}

		comps = append(comps, RawComposition{
			Name:        item.Name,
			Tier:        normalizeTier(item.Tier),
			ChampionIDs: champIDs,
			CoreItems:   coreItems,
		})
	}

	return comps, patch, nil
}

// slugToApiName converts Mobalytics slug to Riot api_name format.
func slugToApiName(slug string) string {
	known := map[string]string{
		"kogmaw":      "TFT16_KogMaw",
		"missfortune": "TFT16_MissFortune",
		"luciansenna": "TFT16_LucianSenna",
		"kobukoyuumi": "TFT16_KobukoYuumi",
		"baronnashor": "TFT16_BaronNashor",
		"riftherald":  "TFT16_RiftHerald",
		"drmundo":     "TFT16_DrMundo",
		"reksai":      "TFT16_RekSai",
		"tahmkench":   "TFT16_TahmKench",
		"chogath":     "TFT16_ChoGath",
	}

	if name, ok := known[slug]; ok {
		return name
	}

	if len(slug) == 0 {
		return slug
	}
	return "TFT16_" + strings.ToUpper(slug[:1]) + slug[1:]
}

// --- GraphQL response types ---

type mobalyticsGQLResponse struct {
	Data struct {
		TFT struct {
			MetaCompositions struct {
				Items []mobalyticsCompItem `json:"items"`
			} `json:"metaCompositions"`
		} `json:"tft"`
	} `json:"data"`
}

type mobalyticsCompItem struct {
	Name      string `json:"name"`
	Tier      string `json:"tier"`
	Patch     string `json:"patch"`
	Slug      string `json:"slug"`
	Formation struct {
		Positions []mobalyticsPosition `json:"positions"`
	} `json:"formation"`
}

type mobalyticsPosition struct {
	Champion struct {
		Champion struct {
			Slug string `json:"slug"`
		} `json:"champion"`
		Items []struct {
			Slug string `json:"slug"`
		} `json:"items"`
		Level int `json:"level"`
	} `json:"champion"`
}
