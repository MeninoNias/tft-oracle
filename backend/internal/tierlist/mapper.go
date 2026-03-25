package tierlist

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
)

func mapToProto(row generated.ConsolidatedTierList) *tftv1.TierListEntry {
	entry := &tftv1.TierListEntry{
		CompositionName:   row.CompositionName,
		ConsolidatedTier:  row.ConsolidatedTier,
		ConsolidatedScore: numericToFloat64(row.ConsolidatedScore),
		Confidence:        row.Confidence,
		Consensus:         row.Consensus,
		MetatftTier:       textToString(row.MetatftTier),
		TftacticsTier:     textToString(row.TftacticsTier),
		MobalyticsTier:    textToString(row.MobalyticsTier),
		AvgWinRate:        numericToFloat64(row.AvgWinRate),
		AvgPlayRate:       numericToFloat64(row.AvgPlayRate),
		AvgPlacement:      numericToFloat64(row.AvgPlacement),
		CoreChampions:     row.CoreChampions,
	}

	// Parse recommended items from JSONB
	if len(row.RecommendedItems) > 0 {
		var items map[string][]string
		if err := json.Unmarshal(row.RecommendedItems, &items); err == nil {
			entry.RecommendedItems = make(map[string]*tftv1.ItemList, len(items))
			for champ, itemList := range items {
				entry.RecommendedItems[champ] = &tftv1.ItemList{
					ItemApiNames: itemList,
				}
			}
		}
	}

	return entry
}

func textToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid || n.Int == nil {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}
