package simulation

import (
	"encoding/json"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
	"github.com/MeninoNias/tft-oracle/backend/internal/ai"
)

func mapAnalysisToProto(a *ai.BattleAnalysis) *tftv1.SimulateBattleResponse {
	changes := make([]*tftv1.SuggestedChange, 0, len(a.SuggestedChanges))
	for _, c := range a.SuggestedChanges {
		changes = append(changes, &tftv1.SuggestedChange{
			Description: c.Description,
			Priority:    c.Priority,
			Category:    c.Category,
		})
	}

	return &tftv1.SimulateBattleResponse{
		WinProbability:   a.WinProbability,
		Confidence:       a.Confidence,
		Analysis:         a.Analysis,
		PositioningTip:   a.PositioningTip,
		KeyFactors:       a.KeyFactors,
		SuggestedChanges: changes,
	}
}

// DB JSONB → Proto helpers (reused from patch service pattern)

type dbStats struct {
	HP          float64 `json:"hp"`
	Armor       float64 `json:"armor"`
	MagicResist float64 `json:"magic_resist"`
	Damage      float64 `json:"damage"`
	AttackSpeed float64 `json:"attack_speed"`
	Range       float64 `json:"range"`
}

func mapStatsFromDB(raw []byte) *tftv1.ChampionStats {
	var s dbStats
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil
	}
	return &tftv1.ChampionStats{
		Hp:          s.HP,
		Armor:       s.Armor,
		MagicResist: s.MagicResist,
		Damage:      s.Damage,
		AttackSpeed: s.AttackSpeed,
		Range:       s.Range,
	}
}

type dbTraitEffect struct {
	MinUnits int32 `json:"min_units"`
	MaxUnits int32 `json:"max_units"`
	Style    int32 `json:"style"`
}

func mapTraitEffectsFromDB(raw []byte) []*tftv1.TraitEffect {
	var effects []dbTraitEffect
	if err := json.Unmarshal(raw, &effects); err != nil {
		return nil
	}
	result := make([]*tftv1.TraitEffect, 0, len(effects))
	for _, e := range effects {
		result = append(result, &tftv1.TraitEffect{
			MinUnits: e.MinUnits,
			MaxUnits: e.MaxUnits,
			Style:    e.Style,
		})
	}
	return result
}
