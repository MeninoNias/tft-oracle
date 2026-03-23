package ai

// BattleAnalysis is the structured response from the AI model.
type BattleAnalysis struct {
	WinProbability   float64           `json:"win_probability"`
	Confidence       float64           `json:"confidence"`
	Analysis         string            `json:"analysis"`
	PositioningTip   string            `json:"positioning_tip"`
	KeyFactors       []string          `json:"key_factors"`
	SuggestedChanges []SuggestedChange `json:"suggested_changes"`
}

// SuggestedChange is a single actionable improvement suggestion.
type SuggestedChange struct {
	Description string `json:"description"`
	Priority    string `json:"priority"` // "high", "medium", "low"
	Category    string `json:"category"` // "items", "positioning", "composition", "economy"
}

// EnrichedBoard is a board state enriched with full game data from the DB.
type EnrichedBoard struct {
	Level     int32
	Augments  []string
	Champions []EnrichedChampion
	Traits    []ActiveTrait
}

// EnrichedChampion is a placed champion enriched with stats and trait info.
type EnrichedChampion struct {
	ApiName   string
	Name      string
	Cost      int32
	StarLevel int32
	Position  int32
	Items     []ItemInfo
	HP        float64
	Damage    float64
	AttackSpd float64
	Armor     float64
	MagicRes  float64
	Range     float64
	Traits    []string
}

// ItemInfo holds basic item data for prompt building.
type ItemInfo struct {
	ApiName string
	Name    string
}

// ActiveTrait represents a trait activation on a board.
type ActiveTrait struct {
	ApiName   string
	Name      string
	Count     int32 // how many units contribute
	Threshold int32 // breakpoint reached
	Style     int32 // 0=inactive, 1=bronze, 2=silver, 3=gold, 4=chromatic
}
