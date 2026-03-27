package consolidation

import "math"

// TierToScore converts a letter tier to a numeric score (0-100).
var TierToScore = map[string]float64{
	"S": 95,
	"A": 80,
	"B": 60,
	"C": 40,
	"D": 20,
}

// ScoreToTier converts a numeric score back to a letter tier.
func ScoreToTier(score float64) string {
	switch {
	case score >= 88:
		return "S"
	case score >= 70:
		return "A"
	case score >= 50:
		return "B"
	case score >= 30:
		return "C"
	default:
		return "D"
	}
}

// DefaultWeights are the source weights for cross-ranking.
var DefaultWeights = map[string]float64{
	"mobalytics":   0.50,
	"tacticstools": 0.35,
	"metatft":      0.15,
}

// ScoreGroup computes the consolidated result for a match group.
func ScoreGroup(group MatchGroup, weights map[string]float64) ConsolidatedResult {
	if weights == nil {
		weights = DefaultWeights
	}

	// Collect per-source data
	tiers := make(map[string]string)       // source -> tier
	scores := make(map[string]float64)     // source -> numeric score
	winRates := make([]float64, 0)
	playRates := make([]float64, 0)
	placements := make([]float64, 0)

	for _, src := range group.Sources {
		tiers[src.Source] = src.Tier
		scores[src.Source] = TierToScore[src.Tier]

		if src.WinRate > 0 {
			winRates = append(winRates, src.WinRate)
		}
		if src.PlayRate > 0 {
			playRates = append(playRates, src.PlayRate)
		}
		if src.AvgPlacement > 0 {
			placements = append(placements, src.AvgPlacement)
		}
	}

	// Weighted score
	totalWeight := 0.0
	weightedScore := 0.0
	for source, score := range scores {
		w := weights[source]
		if w == 0 {
			w = 1.0 / float64(len(scores))
		}
		weightedScore += score * w
		totalWeight += w
	}
	if totalWeight > 0 {
		weightedScore /= totalWeight
	}

	// Confidence and consensus
	confidence := computeConfidence(tiers)
	consensus := computeConsensus(tiers)

	return ConsolidatedResult{
		Name:              group.Name,
		ConsolidatedTier:  ScoreToTier(weightedScore),
		ConsolidatedScore: math.Round(weightedScore*100) / 100,
		Confidence:        confidence,
		Consensus:         consensus,
		MetaTFTTier:       tiers["metatft"],
		TFTacticsTier:     tiers["tftactics"],
		MobalyticsTier:    tiers["mobalytics"],
		AvgWinRate:        avg(winRates),
		AvgPlayRate:       avg(playRates),
		AvgPlacement:      avg(placements),
		CoreChampions:     group.ChampionIDs,
		RecommendedItems:  group.CoreItems,
	}
}

// computeConfidence determines confidence based on source agreement.
func computeConfidence(tiers map[string]string) string {
	if len(tiers) <= 1 {
		return "low"
	}

	// Count how many sources agree within 1 tier
	tierValues := make([]float64, 0, len(tiers))
	for _, t := range tiers {
		tierValues = append(tierValues, TierToScore[t])
	}

	// Check if all are within 20 points (1 tier difference)
	allClose := true
	for i := 0; i < len(tierValues); i++ {
		for j := i + 1; j < len(tierValues); j++ {
			if math.Abs(tierValues[i]-tierValues[j]) > 20 {
				allClose = false
			}
		}
	}

	if allClose && len(tiers) >= 3 {
		return "high"
	}

	// Check if majority (2/3) agree
	tierCounts := make(map[string]int)
	for _, t := range tiers {
		tierCounts[t]++
	}
	for _, count := range tierCounts {
		if count >= 2 {
			return "medium"
		}
	}

	return "low"
}

// computeConsensus determines the agreement level across sources.
func computeConsensus(tiers map[string]string) string {
	if len(tiers) <= 1 {
		return "single_source"
	}

	tierCounts := make(map[string]int)
	for _, t := range tiers {
		tierCounts[t]++
	}

	if len(tierCounts) == 1 {
		return "unanimous"
	}

	for _, count := range tierCounts {
		if count >= 2 {
			return "majority"
		}
	}

	return "split"
}

func avg(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return math.Round(sum/float64(len(values))*100) / 100
}
