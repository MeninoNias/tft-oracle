package ai

import (
	"fmt"
	"sort"
	"strings"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
)

// GameData holds lookup maps for enriching board states with full game data.
type GameData struct {
	Champions map[string]*tftv1.Champion // apiName → Champion
	Items     map[string]*tftv1.Item     // apiName → Item
	Traits    map[string]*tftv1.Trait    // apiName → Trait
	// ChampionTraits maps champion apiName → list of trait apiNames
	ChampionTraits map[string][]string
}

// BuildBattlePrompt constructs the user prompt for the AI model from two board states.
func BuildBattlePrompt(playerBoard, opponentBoard *tftv1.BoardState, data *GameData) string {
	var sb strings.Builder

	playerEnriched := EnrichBoard(playerBoard, data)
	sb.WriteString(formatBoard("PLAYER", playerEnriched))

	if opponentBoard != nil && len(opponentBoard.Champions) > 0 {
		opponentEnriched := EnrichBoard(opponentBoard, data)
		sb.WriteString("\n")
		sb.WriteString(formatBoard("OPPONENT", opponentEnriched))
		sb.WriteString("\nAnalyze this matchup. Who wins and why? Provide positioning tips and improvement suggestions.")
	} else {
		sb.WriteString("\nNo opponent board provided. Analyze this composition's strengths, weaknesses, and potential counters. Suggest improvements.")
	}

	return sb.String()
}

// EnrichBoard joins a proto BoardState with full game data from the DB.
func EnrichBoard(board *tftv1.BoardState, data *GameData) *EnrichedBoard {
	enriched := &EnrichedBoard{
		Level:    board.Level,
		Augments: board.Augments,
	}

	for _, pc := range board.Champions {
		ec := EnrichedChampion{
			ApiName:   pc.ChampionApiName,
			StarLevel: pc.StarLevel,
			Position:  pc.Position,
		}

		if champ, ok := data.Champions[pc.ChampionApiName]; ok {
			ec.Name = champ.Name
			ec.Cost = champ.Cost
			if champ.Stats != nil {
				ec.HP = champ.Stats.Hp
				ec.Damage = champ.Stats.Damage
				ec.AttackSpd = champ.Stats.AttackSpeed
				ec.Armor = champ.Stats.Armor
				ec.MagicRes = champ.Stats.MagicResist
				ec.Range = champ.Stats.Range
			}
		} else {
			ec.Name = pc.ChampionApiName // fallback to apiName
		}

		if traits, ok := data.ChampionTraits[pc.ChampionApiName]; ok {
			ec.Traits = traits
		}

		for _, itemAPI := range pc.ItemApiNames {
			info := ItemInfo{ApiName: itemAPI}
			if item, ok := data.Items[itemAPI]; ok {
				info.Name = item.Name
			} else {
				info.Name = itemAPI
			}
			ec.Items = append(ec.Items, info)
		}

		enriched.Champions = append(enriched.Champions, ec)
	}

	enriched.Traits = ComputeActiveTraits(enriched.Champions, data)
	return enriched
}

// ComputeActiveTraits calculates which traits are active and at what tier.
func ComputeActiveTraits(champions []EnrichedChampion, data *GameData) []ActiveTrait {
	// Count units per trait
	traitCounts := make(map[string]int32)
	for _, c := range champions {
		for _, t := range c.Traits {
			traitCounts[t]++
		}
	}

	var active []ActiveTrait
	for apiName, count := range traitCounts {
		at := ActiveTrait{
			ApiName: apiName,
			Count:   count,
		}

		if trait, ok := data.Traits[apiName]; ok {
			at.Name = trait.Name
			// Find the highest threshold met
			for _, effect := range trait.Effects {
				if count >= effect.MinUnits {
					at.Threshold = effect.MinUnits
					at.Style = effect.Style
				}
			}
		} else {
			at.Name = apiName
		}

		if at.Style > 0 { // only include active traits (not 0-style/inactive)
			active = append(active, at)
		}
	}

	// Sort: higher style first, then by count desc
	sort.Slice(active, func(i, j int) bool {
		if active[i].Style != active[j].Style {
			return active[i].Style > active[j].Style
		}
		return active[i].Count > active[j].Count
	})

	return active
}

func formatBoard(label string, board *EnrichedBoard) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("=== %s BOARD (Level %d) ===\n", label, board.Level))

	if len(board.Augments) > 0 {
		sb.WriteString(fmt.Sprintf("Augments: %s\n", strings.Join(board.Augments, ", ")))
	}

	sb.WriteString("\nChampions:\n")
	for i, c := range board.Champions {
		row := c.Position / 7
		col := c.Position % 7

		items := make([]string, 0, len(c.Items))
		for _, item := range c.Items {
			items = append(items, item.Name)
		}

		sb.WriteString(fmt.Sprintf("  %d. %s (%dg, %d*) @ row %d col %d\n",
			i+1, c.Name, c.Cost, c.StarLevel, row, col))

		if len(items) > 0 {
			sb.WriteString(fmt.Sprintf("     Items: %s\n", strings.Join(items, ", ")))
		}

		sb.WriteString(fmt.Sprintf("     Stats: HP=%.0f AD=%.0f AS=%.2f Armor=%.0f MR=%.0f Range=%.0f\n",
			c.HP, c.Damage, c.AttackSpd, c.Armor, c.MagicRes, c.Range))

		if len(c.Traits) > 0 {
			traitNames := make([]string, 0, len(c.Traits))
			for _, t := range c.Traits {
				traitNames = append(traitNames, t)
			}
			sb.WriteString(fmt.Sprintf("     Traits: %s\n", strings.Join(traitNames, ", ")))
		}
	}

	if len(board.Traits) > 0 {
		sb.WriteString("\nActive Traits:\n")
		for _, t := range board.Traits {
			styleName := styleToName(t.Style)
			sb.WriteString(fmt.Sprintf("  - %s: %d units (%s)\n", t.Name, t.Count, styleName))
		}
	}

	return sb.String()
}

func styleToName(style int32) string {
	switch style {
	case 1:
		return "bronze"
	case 2:
		return "silver"
	case 3:
		return "gold"
	case 4:
		return "chromatic"
	default:
		return "inactive"
	}
}
