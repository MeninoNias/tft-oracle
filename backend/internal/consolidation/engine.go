package consolidation

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
)

// Engine consolidates raw tier data from multiple sources into a unified tier list.
type Engine struct {
	db      *pgxpool.Pool
	queries *generated.Queries
	weights map[string]float64
}

// NewEngine creates a new consolidation engine.
func NewEngine(db *pgxpool.Pool, weights map[string]float64) *Engine {
	if weights == nil {
		weights = DefaultWeights
	}
	return &Engine{
		db:      db,
		queries: generated.New(db),
		weights: weights,
	}
}

// Consolidate fetches raw tier data for a patch, matches compositions,
// scores them, and stores the results.
func (e *Engine) Consolidate(ctx context.Context, patch string) error {
	log.Printf("consolidation: processing patch %s...", patch)

	// 1. Fetch all raw data for this patch
	rawData, err := e.queries.GetRawTierDataByPatch(ctx, patch)
	if err != nil {
		return fmt.Errorf("fetch raw data: %w", err)
	}

	if len(rawData) == 0 {
		log.Printf("consolidation: no raw data for patch %s", patch)
		return nil
	}

	// 2. Group by source and normalize
	bySource := make(map[string][]NormalizedComp)
	for _, row := range rawData {
		comp := NormalizedComp{
			Source:      row.Source,
			Name:        row.CompositionName,
			Tier:        textToString(row.Tier),
			Score:       TierToScore[textToString(row.Tier)],
			WinRate:     numericToFloat(row.WinRate),
			PlayRate:    numericToFloat(row.PlayRate),
			AvgPlacement: numericToFloat(row.AvgPlacement),
			ChampionIDs: row.ChampionIds,
		}

		// Parse core items from JSONB
		if len(row.CoreItems) > 0 {
			var items map[string][]string
			if err := json.Unmarshal(row.CoreItems, &items); err == nil {
				comp.CoreItems = items
			}
		}
		if comp.CoreItems == nil {
			comp.CoreItems = make(map[string][]string)
		}

		bySource[row.Source] = append(bySource[row.Source], comp)
	}

	// 3. Match compositions across sources
	groups := MatchCompositions(bySource)
	log.Printf("consolidation: matched %d raw entries into %d composition groups", len(rawData), len(groups))

	// 4. Score each group
	results := make([]ConsolidatedResult, 0, len(groups))
	for _, group := range groups {
		result := ScoreGroup(group, e.weights)
		results = append(results, result)
	}

	// 5. Store consolidated results
	for _, r := range results {
		itemsJSON, err := json.Marshal(r.RecommendedItems)
		if err != nil {
			return fmt.Errorf("marshal items for %q: %w", r.Name, err)
		}

		champions := r.CoreChampions
		if champions == nil {
			champions = []string{}
		}

		_, err = e.queries.UpsertConsolidatedTierEntry(ctx, generated.UpsertConsolidatedTierEntryParams{
			Patch:             patch,
			CompositionName:   r.Name,
			ConsolidatedTier:  r.ConsolidatedTier,
			ConsolidatedScore: floatToNumeric(r.ConsolidatedScore),
			Confidence:        r.Confidence,
			Consensus:         r.Consensus,
			MetatftTier:       stringToText(r.MetaTFTTier),
			TftacticsTier:     stringToText(r.TFTacticsTier),
			MobalyticsTier:    stringToText(r.MobalyticsTier),
			AvgWinRate:        floatToNumeric(r.AvgWinRate),
			AvgPlayRate:       floatToNumeric(r.AvgPlayRate),
			AvgPlacement:      floatToNumeric(r.AvgPlacement),
			CoreChampions:     champions,
			RecommendedItems:  itemsJSON,
		})
		if err != nil {
			return fmt.Errorf("upsert consolidated %q: %w", r.Name, err)
		}
	}

	log.Printf("consolidation: stored %d consolidated compositions for patch %s", len(results), patch)
	return nil
}

func textToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid || n.Int == nil {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}

func stringToText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func floatToNumeric(f float64) pgtype.Numeric {
	if f == 0 {
		return pgtype.Numeric{Valid: false}
	}
	scaled := int64(f * 100)
	return pgtype.Numeric{
		Int:   big.NewInt(scaled),
		Exp:   -2,
		Valid: true,
	}
}
