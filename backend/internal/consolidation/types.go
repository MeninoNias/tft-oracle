package consolidation

// NormalizedComp represents a composition from one source after normalization.
type NormalizedComp struct {
	Source       string
	Name         string
	Tier         string
	Score        float64
	WinRate      float64
	PlayRate     float64
	AvgPlacement float64
	ChampionIDs  []string
	CoreItems    map[string][]string
}

// MatchGroup represents compositions from multiple sources that are the same comp.
type MatchGroup struct {
	Name         string            // canonical name (from first source found)
	Sources      []NormalizedComp  // one per source (max 3)
	ChampionIDs  []string          // union of all champion lists
	CoreItems    map[string][]string
}

// ConsolidatedResult is the output of the scoring step for one composition.
type ConsolidatedResult struct {
	Name              string
	ConsolidatedTier  string
	ConsolidatedScore float64
	Confidence        string
	Consensus         string
	MetaTFTTier       string
	TFTacticsTier     string
	MobalyticsTier    string
	AvgWinRate        float64
	AvgPlayRate       float64
	AvgPlacement      float64
	CoreChampions     []string
	RecommendedItems  map[string][]string
}
