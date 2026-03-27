package coach

import (
	"encoding/json"
	"math/big"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
	"github.com/MeninoNias/tft-oracle/backend/internal/ai"
	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

// --- AI response → Proto ---

func mapMatchAnalysisToProto(matchID string, placement int32, a *ai.MatchCoachAnalysis) *tftv1.AnalyzeMatchResponse {
	insights := make([]*tftv1.CoachingInsight, 0, len(a.Insights))
	for _, ins := range a.Insights {
		insights = append(insights, &tftv1.CoachingInsight{
			Category: ins.Category,
			Title:    ins.Title,
			Detail:   ins.Detail,
			Grade:    ins.Grade,
		})
	}

	suggestions := make([]*tftv1.CoachSuggestion, 0, len(a.Suggestions))
	for _, s := range a.Suggestions {
		suggestions = append(suggestions, &tftv1.CoachSuggestion{
			Description: s.Description,
			Priority:    s.Priority,
			Category:    s.Category,
		})
	}

	return &tftv1.AnalyzeMatchResponse{
		MatchId:      matchID,
		Placement:    placement,
		OverallGrade: a.OverallGrade,
		Summary:      a.Summary,
		Insights:     insights,
		MetaComparison: &tftv1.MetaComparison{
			ClosestMetaComp: a.MetaComparison.ClosestMetaComp,
			MetaTier:        a.MetaComparison.MetaTier,
			MissingUnits:    a.MetaComparison.MissingUnits,
			SuboptimalItems: a.MetaComparison.SuboptimalItems,
			Assessment:      a.MetaComparison.Assessment,
		},
		LobbyContext: &tftv1.LobbyContext{
			ContestedCount:          a.LobbyContext.ContestedCount,
			LobbyStrengthAssessment: a.LobbyContext.LobbyStrengthAssessment,
			ContestedDetails:        a.LobbyContext.ContestedDetails,
		},
		Suggestions: suggestions,
	}
}

func mapHistoryAnalysisToProto(matchCount int32, a *ai.HistoryCoachAnalysis) *tftv1.AnalyzeHistoryResponse {
	patterns := make([]*tftv1.PatternInsight, 0, len(a.Patterns))
	for _, p := range a.Patterns {
		patterns = append(patterns, &tftv1.PatternInsight{
			Category:  p.Category,
			Title:     p.Title,
			Detail:    p.Detail,
			Sentiment: p.Sentiment,
		})
	}

	trends := make([]*tftv1.TrendItem, 0, len(a.Trends))
	for _, t := range a.Trends {
		trends = append(trends, &tftv1.TrendItem{
			Metric:    t.Metric,
			Direction: t.Direction,
			Detail:    t.Detail,
		})
	}

	plan := make([]*tftv1.CoachSuggestion, 0, len(a.ImprovementPlan))
	for _, s := range a.ImprovementPlan {
		plan = append(plan, &tftv1.CoachSuggestion{
			Description: s.Description,
			Priority:    s.Priority,
			Category:    s.Category,
		})
	}

	return &tftv1.AnalyzeHistoryResponse{
		MatchesAnalyzed: matchCount,
		OverallSummary:  a.OverallSummary,
		SkillRadar: &tftv1.SkillRadar{
			Economy:      float32(a.SkillRadar.Economy),
			Itemization:  float32(a.SkillRadar.Itemization),
			Composition:  float32(a.SkillRadar.Composition),
			Adaptability: float32(a.SkillRadar.Adaptability),
			Consistency:  float32(a.SkillRadar.Consistency),
		},
		Patterns:        patterns,
		Trends:          trends,
		ImprovementPlan: plan,
	}
}

// --- DB → AI types ---

// dbTrait mirrors the JSONB structure stored in match_participants.traits.
type dbTrait struct {
	Name        string `json:"name"`
	NumUnits    int32  `json:"num_units"`
	Style       int32  `json:"style"`
	TierCurrent int32  `json:"tier_current"`
	TierTotal   int32  `json:"tier_total"`
}

// dbUnit mirrors the JSONB structure stored in match_participants.units.
type dbUnit struct {
	CharacterID string   `json:"character_id"`
	Name        string   `json:"name"`
	Tier        int32    `json:"tier"`
	Rarity      int32    `json:"rarity"`
	Items       []string `json:"items"`
}

func mapParticipantToCoachData(p generated.MatchParticipant) ai.MatchParticipantData {
	var traits []dbTrait
	_ = json.Unmarshal(p.Traits, &traits)

	var units []dbUnit
	_ = json.Unmarshal(p.Units, &units)

	aiTraits := make([]ai.ParticipantTrait, 0, len(traits))
	for _, t := range traits {
		aiTraits = append(aiTraits, ai.ParticipantTrait{
			ApiName:     t.Name, // stored as trait name/apiName in JSONB
			NumUnits:    t.NumUnits,
			Style:       t.Style,
			TierCurrent: t.TierCurrent,
			TierTotal:   t.TierTotal,
		})
	}

	aiUnits := make([]ai.ParticipantUnit, 0, len(units))
	for _, u := range units {
		aiUnits = append(aiUnits, ai.ParticipantUnit{
			CharacterID: u.CharacterID,
			Name:        u.Name,
			Tier:        u.Tier,
			Rarity:      u.Rarity,
			Items:       u.Items,
		})
	}

	return ai.MatchParticipantData{
		Puuid:                p.Puuid,
		Placement:            p.Placement,
		Level:                p.Level,
		GoldLeft:             p.GoldLeft,
		LastRound:            p.LastRound,
		TimeEliminated:       p.TimeEliminated,
		TotalDamageToPlayers: p.TotalDamageToPlayers,
		PlayersEliminated:    p.PlayersEliminated,
		Augments:             p.Augments,
		Traits:               aiTraits,
		Units:                aiUnits,
	}
}

func mapParticipantToSummary(p generated.MatchParticipant) ai.MatchSummary {
	data := mapParticipantToCoachData(p)
	return ai.MatchSummary{
		MatchID:   p.MatchID,
		Placement: data.Placement,
		Level:     data.Level,
		GoldLeft:  data.GoldLeft,
		LastRound: data.LastRound,
		Augments:  data.Augments,
		Units:     data.Units,
		Traits:    data.Traits,
	}
}

// --- Tier list DB → AI types ---

func mapTierListForPrompt(entries []generated.ConsolidatedTierList) []ai.TierListEntry {
	result := make([]ai.TierListEntry, 0, len(entries))
	for _, e := range entries {
		entry := ai.TierListEntry{
			CompositionName: e.CompositionName,
			Tier:            e.ConsolidatedTier,
			Score:           numericToFloat(e.ConsolidatedScore),
			CoreChampions:   e.CoreChampions,
			WinRate:         numericToFloat(e.AvgWinRate),
			AvgPlacement:    numericToFloat(e.AvgPlacement),
		}

		// Decode recommended items JSONB
		if len(e.RecommendedItems) > 0 {
			var items map[string][]string
			if err := json.Unmarshal(e.RecommendedItems, &items); err == nil {
				entry.RecommendedItems = items
			}
		}

		result = append(result, entry)
	}
	return result
}

func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Int.Float64()
	if n.Exp != 0 {
		// Adjust for exponent
		exp := new(big.Float).SetFloat64(1)
		for i := int32(0); i < -n.Exp; i++ {
			exp.Mul(exp, new(big.Float).SetFloat64(10))
		}
		result := new(big.Float).SetFloat64(f)
		result.Quo(result, exp)
		r, _ := result.Float64()
		return r
	}
	return f
}
