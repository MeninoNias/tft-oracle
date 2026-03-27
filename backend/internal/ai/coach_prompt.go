package ai

import (
	"fmt"
	"sort"
	"strings"
)

// --- Helper types for prompt building ---

// MatchMeta holds match-level metadata for prompt context.
type MatchMeta struct {
	GameLength float32
	LastRound  int32
	GameType   string
}

// TierListEntry holds a meta composition for prompt context.
type TierListEntry struct {
	CompositionName  string
	Tier             string
	Score            float64
	CoreChampions    []string
	RecommendedItems map[string][]string // champion apiName → item display names
	WinRate          float64
	AvgPlacement     float64
}

// MatchParticipantData holds participant data for prompt building.
type MatchParticipantData struct {
	Puuid                 string
	Placement             int32
	Level                 int32
	GoldLeft              int32
	LastRound             int32
	TimeEliminated        float32
	TotalDamageToPlayers  int32
	PlayersEliminated     int32
	Augments              []string
	Traits                []ParticipantTrait
	Units                 []ParticipantUnit
}

// ParticipantTrait is a trait from a match participant's JSONB data.
type ParticipantTrait struct {
	ApiName     string
	NumUnits    int32
	Style       int32
	TierCurrent int32
	TierTotal   int32
}

// ParticipantUnit is a unit from a match participant's JSONB data.
type ParticipantUnit struct {
	CharacterID string
	Name        string
	Tier        int32 // star level
	Rarity      int32 // cost tier (0-4 = 1-5 gold)
	Items       []string
}

// MatchSummary is a condensed match record for history prompt building.
type MatchSummary struct {
	MatchID   string
	Placement int32
	Level     int32
	GoldLeft  int32
	LastRound int32
	Augments  []string
	Units     []ParticipantUnit
	Traits    []ParticipantTrait
}

// HistoryAggregates holds computed stats across multiple matches.
type HistoryAggregates struct {
	AvgPlacement  float64
	Top4Rate      float64
	WinRate       float64
	PlacementDist map[int32]int
	TopComps      []CompFrequency
	TopChampions  []ChampFrequency
	TopItems      []ItemFrequency
	TopAugments   []AugmentFrequency
}

// CompFrequency tracks how often a composition was played.
type CompFrequency struct {
	TraitCombo   string
	Games        int
	AvgPlacement float64
}

// ChampFrequency tracks how often a champion was played.
type ChampFrequency struct {
	Name         string
	Games        int
	AvgPlacement float64
}

// ItemFrequency tracks how often an item was built.
type ItemFrequency struct {
	Name  string
	Count int
}

// AugmentFrequency tracks how often an augment was picked.
type AugmentFrequency struct {
	Name         string
	Games        int
	AvgPlacement float64
}

// BuildMatchCoachPrompt constructs the user prompt for per-match coaching analysis.
func BuildMatchCoachPrompt(
	player MatchParticipantData,
	allParticipants []MatchParticipantData,
	meta MatchMeta,
	data *GameData,
	tierList []TierListEntry,
	playerRank string,
) string {
	var sb strings.Builder

	// Player section
	sb.WriteString("=== YOUR MATCH RESULT ===\n")
	sb.WriteString(fmt.Sprintf("Placement: %d | Level: %d | Gold Left: %d | Last Round: %s\n",
		player.Placement, player.Level, player.GoldLeft, formatRound(player.LastRound)))
	sb.WriteString(fmt.Sprintf("Rank: %s\n", playerRank))
	sb.WriteString(fmt.Sprintf("Damage Dealt: %d | Players Eliminated: %d\n",
		player.TotalDamageToPlayers, player.PlayersEliminated))
	gameLenMin := meta.GameLength / 60
	sb.WriteString(fmt.Sprintf("Game Length: %.0fm | Game Type: %s\n", gameLenMin, meta.GameType))

	if len(player.Augments) > 0 {
		augNames := resolveAugmentNames(player.Augments, data)
		sb.WriteString(fmt.Sprintf("\nAugments: %s\n", strings.Join(augNames, ", ")))
	}

	// Player's units
	sb.WriteString("\nChampions:\n")
	for i, u := range player.Units {
		name := resolveChampionName(u.CharacterID, data)
		cost := resolveChampionCost(u.CharacterID, data)
		items := resolveItemNames(u.Items, data)
		sb.WriteString(fmt.Sprintf("  %d. %s (%dg, %d*)", i+1, name, cost, u.Tier))
		if len(items) > 0 {
			sb.WriteString(fmt.Sprintf(" — Items: %s", strings.Join(items, ", ")))
		}
		sb.WriteString("\n")
	}

	// Player's active traits
	activeTraits := filterActiveTraits(player.Traits)
	if len(activeTraits) > 0 {
		sb.WriteString("\nActive Traits:\n")
		for _, t := range activeTraits {
			name := resolveTraitName(t.ApiName, data)
			sb.WriteString(fmt.Sprintf("  - %s: %d units (%s)\n", name, t.NumUnits, styleToName(t.Style)))
		}
	}

	// Lobby section
	sb.WriteString("\n=== LOBBY (all 8 players) ===\n")
	for _, p := range allParticipants {
		if p.Puuid == player.Puuid {
			sb.WriteString(fmt.Sprintf("%s %d. (YOU) ", placementEmoji(p.Placement), p.Placement))
		} else {
			sb.WriteString(fmt.Sprintf("%s %d. ", placementEmoji(p.Placement), p.Placement))
		}
		// Condensed unit list (top 5 by cost)
		sortedUnits := sortUnitsByCost(p.Units)
		unitStrs := make([]string, 0, 5)
		for i, u := range sortedUnits {
			if i >= 5 {
				break
			}
			name := resolveChampionName(u.CharacterID, data)
			unitStrs = append(unitStrs, fmt.Sprintf("%s %d*", name, u.Tier))
		}
		sb.WriteString(strings.Join(unitStrs, ", "))

		// Condensed traits
		at := filterActiveTraits(p.Traits)
		if len(at) > 0 {
			traitStrs := make([]string, 0, 3)
			for i, t := range at {
				if i >= 3 {
					break
				}
				name := resolveTraitName(t.ApiName, data)
				traitStrs = append(traitStrs, fmt.Sprintf("%s(%d)", name, t.NumUnits))
			}
			sb.WriteString(" | " + strings.Join(traitStrs, " "))
		}
		sb.WriteString("\n")
	}

	// Meta context
	if len(tierList) > 0 {
		sb.WriteString("\n=== META CONTEXT (current patch) ===\n")
		limit := 15
		if len(tierList) < limit {
			limit = len(tierList)
		}
		for _, tl := range tierList[:limit] {
			sb.WriteString(fmt.Sprintf("%s-tier: %s", tl.Tier, tl.CompositionName))
			if len(tl.CoreChampions) > 0 {
				sb.WriteString(fmt.Sprintf(" (core: %s)", strings.Join(tl.CoreChampions, ", ")))
			}
			if tl.WinRate > 0 {
				sb.WriteString(fmt.Sprintf(" | WR: %.0f%% | Avg: %.1f", tl.WinRate*100, tl.AvgPlacement))
			}
			// Show BIS for core champions
			if len(tl.RecommendedItems) > 0 {
				bisStrs := make([]string, 0)
				for champ, items := range tl.RecommendedItems {
					if len(items) > 0 {
						bisStrs = append(bisStrs, fmt.Sprintf("%s→%s", champ, strings.Join(items, "/")))
					}
				}
				if len(bisStrs) > 0 {
					sb.WriteString(fmt.Sprintf(" | BIS: %s", strings.Join(bisStrs, ", ")))
				}
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n=== INSTRUCTIONS ===\n")
	sb.WriteString(fmt.Sprintf("Analyze this completed match. The player placed %d%s.\n",
		player.Placement, ordinalSuffix(player.Placement)))
	sb.WriteString("Provide insights for each category (composition, itemization, economy, augments, adaptability).\n")
	sb.WriteString("Compare against the meta tier list. Analyze lobby context for contested comps.\n")
	sb.WriteString("Give specific, actionable suggestions prioritized by impact.\n")

	return sb.String()
}

// BuildHistoryCoachPrompt constructs the user prompt for multi-match coaching analysis.
func BuildHistoryCoachPrompt(
	matches []MatchSummary,
	aggregates HistoryAggregates,
	data *GameData,
	tierList []TierListEntry,
	playerRank string,
) string {
	var sb strings.Builder

	sb.WriteString("=== PLAYER PROFILE ===\n")
	sb.WriteString(fmt.Sprintf("Rank: %s | Matches analyzed: %d\n", playerRank, len(matches)))

	// Match list
	sb.WriteString("\n=== MATCH HISTORY ===\n")
	for i, m := range matches {
		sb.WriteString(fmt.Sprintf("Match %d: %s | Lv%d | %dg left | Round %s",
			i+1, formatPlacement(m.Placement), m.Level, m.GoldLeft, formatRound(m.LastRound)))
		// Top 3 units
		sorted := sortUnitsByCost(m.Units)
		unitStrs := make([]string, 0, 3)
		for j, u := range sorted {
			if j >= 3 {
				break
			}
			name := resolveChampionName(u.CharacterID, data)
			unitStrs = append(unitStrs, fmt.Sprintf("%s %d*", name, u.Tier))
		}
		if len(unitStrs) > 0 {
			sb.WriteString(" | " + strings.Join(unitStrs, ", "))
		}
		// Augments
		if len(m.Augments) > 0 {
			augNames := resolveAugmentNames(m.Augments, data)
			sb.WriteString(" | Aug: " + strings.Join(augNames, ", "))
		}
		sb.WriteString("\n")
	}

	// Aggregate stats
	sb.WriteString("\n=== AGGREGATE STATS ===\n")
	sb.WriteString(fmt.Sprintf("Avg Placement: %.1f | Top 4: %.0f%% | Win Rate: %.0f%%\n",
		aggregates.AvgPlacement, aggregates.Top4Rate*100, aggregates.WinRate*100))
	sb.WriteString("Placement Distribution: ")
	for p := int32(1); p <= 8; p++ {
		count := aggregates.PlacementDist[p]
		sb.WriteString(fmt.Sprintf("%d%s=%d ", p, ordinalSuffix(p), count))
	}
	sb.WriteString("\n")

	// Comp frequency
	if len(aggregates.TopComps) > 0 {
		sb.WriteString("\n=== MOST PLAYED COMPS ===\n")
		for _, c := range aggregates.TopComps {
			sb.WriteString(fmt.Sprintf("  %s: %d games (avg %.1f)\n", c.TraitCombo, c.Games, c.AvgPlacement))
		}
	}

	// Champion frequency
	if len(aggregates.TopChampions) > 0 {
		sb.WriteString("\n=== MOST PLAYED CHAMPIONS ===\n")
		for _, c := range aggregates.TopChampions {
			sb.WriteString(fmt.Sprintf("  %s: %d games (avg %.1f)\n", c.Name, c.Games, c.AvgPlacement))
		}
	}

	// Item frequency
	if len(aggregates.TopItems) > 0 {
		sb.WriteString("\n=== MOST BUILT ITEMS ===\n")
		for _, it := range aggregates.TopItems {
			sb.WriteString(fmt.Sprintf("  %s: %d times\n", it.Name, it.Count))
		}
	}

	// Augment frequency
	if len(aggregates.TopAugments) > 0 {
		sb.WriteString("\n=== AUGMENT PATTERNS ===\n")
		for _, a := range aggregates.TopAugments {
			sb.WriteString(fmt.Sprintf("  %s: %d games (avg %.1f)\n", a.Name, a.Games, a.AvgPlacement))
		}
	}

	// Meta context
	if len(tierList) > 0 {
		sb.WriteString("\n=== META CONTEXT (current patch) ===\n")
		limit := 10
		if len(tierList) < limit {
			limit = len(tierList)
		}
		for _, tl := range tierList[:limit] {
			sb.WriteString(fmt.Sprintf("  %s-tier: %s", tl.Tier, tl.CompositionName))
			if tl.WinRate > 0 {
				sb.WriteString(fmt.Sprintf(" | WR: %.0f%% | Avg: %.1f", tl.WinRate*100, tl.AvgPlacement))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n=== INSTRUCTIONS ===\n")
	sb.WriteString("Analyze this player's match history for recurring patterns.\n")
	sb.WriteString("Rate each skill radar dimension (0-100) based on evidence.\n")
	sb.WriteString("Identify trends by comparing first half vs second half of matches.\n")
	sb.WriteString("Provide a concrete 3-5 step improvement plan prioritized by impact.\n")

	return sb.String()
}

// ComputeHistoryAggregates computes aggregate statistics from match summaries.
func ComputeHistoryAggregates(matches []MatchSummary, data *GameData) HistoryAggregates {
	agg := HistoryAggregates{
		PlacementDist: make(map[int32]int),
	}

	if len(matches) == 0 {
		return agg
	}

	totalPlacement := 0
	top4Count := 0
	winCount := 0

	// Track frequencies
	champGames := make(map[string][]int32)     // name → placements
	itemCounts := make(map[string]int)          // name → count
	augmentGames := make(map[string][]int32)    // name → placements
	compGames := make(map[string][]int32)       // trait combo → placements

	for _, m := range matches {
		totalPlacement += int(m.Placement)
		agg.PlacementDist[m.Placement]++
		if m.Placement <= 4 {
			top4Count++
		}
		if m.Placement == 1 {
			winCount++
		}

		// Champion frequency
		for _, u := range m.Units {
			name := resolveChampionName(u.CharacterID, data)
			champGames[name] = append(champGames[name], m.Placement)

			// Item frequency
			for _, itemAPI := range u.Items {
				itemName := resolveItemName(itemAPI, data)
				itemCounts[itemName]++
			}
		}

		// Augment frequency
		for _, aug := range m.Augments {
			augName := resolveItemName(aug, data) // augments are in items table
			augmentGames[augName] = append(augmentGames[augName], m.Placement)
		}

		// Comp frequency (dominant trait combo — top 2 active traits by style)
		combo := dominantTraitCombo(m.Traits, data)
		if combo != "" {
			compGames[combo] = append(compGames[combo], m.Placement)
		}
	}

	n := float64(len(matches))
	agg.AvgPlacement = float64(totalPlacement) / n
	agg.Top4Rate = float64(top4Count) / n
	agg.WinRate = float64(winCount) / n

	// Build sorted frequency lists
	agg.TopComps = buildCompFrequency(compGames, 10)
	agg.TopChampions = buildChampFrequency(champGames, 15)
	agg.TopItems = buildItemFrequency(itemCounts, 15)
	agg.TopAugments = buildAugmentFrequency(augmentGames, 10)

	return agg
}

// --- Helper functions ---

func resolveChampionName(apiName string, data *GameData) string {
	if data != nil {
		if c, ok := data.Champions[apiName]; ok {
			return c.Name
		}
	}
	return apiName
}

func resolveChampionCost(apiName string, data *GameData) int32 {
	if data != nil {
		if c, ok := data.Champions[apiName]; ok {
			return c.Cost
		}
	}
	return 0
}

func resolveItemName(apiName string, data *GameData) string {
	if data != nil {
		if it, ok := data.Items[apiName]; ok {
			return it.Name
		}
	}
	return apiName
}

func resolveItemNames(apiNames []string, data *GameData) []string {
	names := make([]string, 0, len(apiNames))
	for _, api := range apiNames {
		names = append(names, resolveItemName(api, data))
	}
	return names
}

func resolveAugmentNames(apiNames []string, data *GameData) []string {
	// Augments are stored in the items table with "augment" tag
	return resolveItemNames(apiNames, data)
}

func resolveTraitName(apiName string, data *GameData) string {
	if data != nil {
		if t, ok := data.Traits[apiName]; ok {
			return t.Name
		}
	}
	return apiName
}

func filterActiveTraits(traits []ParticipantTrait) []ParticipantTrait {
	active := make([]ParticipantTrait, 0)
	for _, t := range traits {
		if t.Style > 0 {
			active = append(active, t)
		}
	}
	sort.Slice(active, func(i, j int) bool {
		if active[i].Style != active[j].Style {
			return active[i].Style > active[j].Style
		}
		return active[i].NumUnits > active[j].NumUnits
	})
	return active
}

func sortUnitsByCost(units []ParticipantUnit) []ParticipantUnit {
	sorted := make([]ParticipantUnit, len(units))
	copy(sorted, units)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Rarity != sorted[j].Rarity {
			return sorted[i].Rarity > sorted[j].Rarity
		}
		return sorted[i].Tier > sorted[j].Tier
	})
	return sorted
}

func dominantTraitCombo(traits []ParticipantTrait, data *GameData) string {
	active := filterActiveTraits(traits)
	if len(active) == 0 {
		return ""
	}
	limit := 2
	if len(active) < limit {
		limit = len(active)
	}
	names := make([]string, 0, limit)
	for _, t := range active[:limit] {
		names = append(names, resolveTraitName(t.ApiName, data))
	}
	return strings.Join(names, " + ")
}

func buildCompFrequency(compGames map[string][]int32, limit int) []CompFrequency {
	result := make([]CompFrequency, 0, len(compGames))
	for combo, placements := range compGames {
		avg := avgPlacement(placements)
		result = append(result, CompFrequency{TraitCombo: combo, Games: len(placements), AvgPlacement: avg})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Games > result[j].Games })
	if len(result) > limit {
		result = result[:limit]
	}
	return result
}

func buildChampFrequency(champGames map[string][]int32, limit int) []ChampFrequency {
	result := make([]ChampFrequency, 0, len(champGames))
	for name, placements := range champGames {
		avg := avgPlacement(placements)
		result = append(result, ChampFrequency{Name: name, Games: len(placements), AvgPlacement: avg})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Games > result[j].Games })
	if len(result) > limit {
		result = result[:limit]
	}
	return result
}

func buildItemFrequency(itemCounts map[string]int, limit int) []ItemFrequency {
	result := make([]ItemFrequency, 0, len(itemCounts))
	for name, count := range itemCounts {
		result = append(result, ItemFrequency{Name: name, Count: count})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Count > result[j].Count })
	if len(result) > limit {
		result = result[:limit]
	}
	return result
}

func buildAugmentFrequency(augGames map[string][]int32, limit int) []AugmentFrequency {
	result := make([]AugmentFrequency, 0, len(augGames))
	for name, placements := range augGames {
		avg := avgPlacement(placements)
		result = append(result, AugmentFrequency{Name: name, Games: len(placements), AvgPlacement: avg})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Games > result[j].Games })
	if len(result) > limit {
		result = result[:limit]
	}
	return result
}

func avgPlacement(placements []int32) float64 {
	if len(placements) == 0 {
		return 0
	}
	sum := 0
	for _, p := range placements {
		sum += int(p)
	}
	return float64(sum) / float64(len(placements))
}

func formatRound(lastRound int32) string {
	if lastRound <= 3 {
		return fmt.Sprintf("1-%d", lastRound)
	}
	adjusted := lastRound - 3
	stage := adjusted/7 + 2
	round := adjusted%7 + 1
	return fmt.Sprintf("%d-%d", stage, round)
}

func formatPlacement(p int32) string {
	return fmt.Sprintf("%d%s", p, ordinalSuffix(p))
}

func ordinalSuffix(n int32) string {
	if n >= 11 && n <= 13 {
		return "th"
	}
	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func placementEmoji(p int32) string {
	if p <= 4 {
		return "+"
	}
	return "-"
}
