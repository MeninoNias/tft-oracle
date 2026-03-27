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

// --- Coach analysis types ---

// MatchCoachAnalysis is the structured response for per-match coaching.
type MatchCoachAnalysis struct {
	OverallGrade   string            `json:"overall_grade"`
	Summary        string            `json:"summary"`
	Insights       []CoachingInsight `json:"insights"`
	MetaComparison MetaComparison    `json:"meta_comparison"`
	LobbyContext   LobbyContext      `json:"lobby_context"`
	Suggestions    []CoachSuggestion `json:"suggestions"`
}

// CoachingInsight is a per-category coaching evaluation.
type CoachingInsight struct {
	Category string `json:"category"` // "composition", "itemization", "economy", "augments", "adaptability"
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Grade    string `json:"grade"` // "good", "okay", "poor"
}

// MetaComparison compares the player's board against the current meta.
type MetaComparison struct {
	ClosestMetaComp string   `json:"closest_meta_comp"`
	MetaTier        string   `json:"meta_tier"`
	MissingUnits    []string `json:"missing_units"`
	SuboptimalItems []string `json:"suboptimal_items"`
	Assessment      string   `json:"assessment"`
}

// LobbyContext describes how contested the player's comp was in the lobby.
type LobbyContext struct {
	ContestedCount          int32    `json:"contested_count"`
	LobbyStrengthAssessment string   `json:"lobby_strength_assessment"`
	ContestedDetails        []string `json:"contested_details"`
}

// CoachSuggestion is a prioritized coaching action item.
type CoachSuggestion struct {
	Description string `json:"description"`
	Priority    string `json:"priority"` // "high", "medium", "low"
	Category    string `json:"category"` // "composition", "itemization", "economy", "augments", "adaptability"
}

// HistoryCoachAnalysis is the structured response for multi-match coaching.
type HistoryCoachAnalysis struct {
	OverallSummary  string            `json:"overall_summary"`
	SkillRadar      SkillRadar        `json:"skill_radar"`
	Patterns        []PatternInsight  `json:"patterns"`
	Trends          []TrendItem       `json:"trends"`
	ImprovementPlan []CoachSuggestion `json:"improvement_plan"`
}

// SkillRadar rates 5 skill dimensions from 0 to 100.
type SkillRadar struct {
	Economy      float64 `json:"economy"`
	Itemization  float64 `json:"itemization"`
	Composition  float64 `json:"composition"`
	Adaptability float64 `json:"adaptability"`
	Consistency  float64 `json:"consistency"`
}

// PatternInsight describes a recurring pattern across matches.
type PatternInsight struct {
	Category  string `json:"category"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	Sentiment string `json:"sentiment"` // "positive", "neutral", "negative"
}

// TrendItem describes a performance trend direction.
type TrendItem struct {
	Metric    string `json:"metric"`
	Direction string `json:"direction"` // "improving", "stable", "declining"
	Detail    string `json:"detail"`
}
