package cdragon

import (
	"strings"
)

const cdragonAssetBase = "https://raw.communitydragon.org/latest/game/"

// ParsedSet holds the parsed data for a single TFT set.
type ParsedSet struct {
	Number    int
	Name      string
	Mutator   string
	Champions []ParsedChampion
	Traits    []ParsedTrait
	Items     []ParsedItem
}

type ParsedChampion struct {
	APIName       string
	Name          string
	Cost          int
	TraitAPINames []string // Resolved to apiNames
	Stats         CDragonChampionStat
	Ability       ParsedAbility
	IconURL       string
	SquareIconURL string
	TileIconURL   string
}

type ParsedAbility struct {
	Name      string
	Desc      string
	IconURL   string
	Variables []ParsedAbilityVariable
}

type ParsedAbilityVariable struct {
	Name   string
	Values []float64
}

type ParsedTrait struct {
	APIName string
	Name    string
	Desc    string
	IconURL string
	Effects []CDragonTraitEffect
}

type ParsedItem struct {
	APIName            string
	Name               string
	Desc               string
	Composition        []string
	Effects            map[string]float64
	IconURL            string
	AssociatedTraits   []string
	IncompatibleTraits []string
	Tags               []string
	Unique             bool
}

// FindCurrentSet returns the highest set number that is a main set (no suffix mutator).
// Mutators like "PAIRS", "TURBO", "PVEMODE" indicate alternate game modes.
func FindCurrentSet(data *CDragonData) *CDragonSetData {
	var current *CDragonSetData
	for i := range data.SetData {
		sd := &data.SetData[i]
		if isAlternateMode(sd.Mutator) {
			continue
		}
		if current == nil || sd.Number > current.Number {
			current = sd
		}
	}
	return current
}

// isAlternateMode returns true if the mutator indicates a non-standard game mode.
func isAlternateMode(mutator string) bool {
	upper := strings.ToUpper(mutator)
	alternates := []string{"PAIRS", "TURBO", "PVEMODE"}
	for _, alt := range alternates {
		if strings.Contains(upper, alt) {
			return true
		}
	}
	return false
}

// Parse converts raw CDragon data into a ParsedSet for the current set.
func Parse(data *CDragonData) *ParsedSet {
	setData := FindCurrentSet(data)
	if setData == nil {
		return nil
	}

	// Build display name → apiName map for trait resolution.
	// CDragon champion.traits[] contains display names, not apiNames.
	traitDisplayToAPI := buildTraitDisplayMap(setData.Traits)

	// Build set of item apiNames from the set data champions' items
	// Actually, setData doesn't have items — items are in the global list.
	// We filter global items by checking if they belong to the current set via apiName prefix or tags.

	// Parse traits
	traits := make([]ParsedTrait, 0, len(setData.Traits))
	for _, t := range setData.Traits {
		traits = append(traits, ParsedTrait{
			APIName: t.APIName,
			Name:    t.Name,
			Desc:    t.Desc,
			IconURL: ConvertIconPath(t.Icon),
			Effects: t.Effects,
		})
	}

	// Parse champions
	champions := make([]ParsedChampion, 0, len(setData.Champions))
	for _, c := range setData.Champions {
		traitAPINames := resolveTraitNames(c.Traits, traitDisplayToAPI)

		ability := ParsedAbility{
			Name:    c.Ability.Name,
			Desc:    c.Ability.Desc,
			IconURL: ConvertIconPath(c.Ability.Icon),
		}
		for _, v := range c.Ability.Variables {
			values := make([]float64, 0, len(v.Values))
			for _, val := range v.Values {
				if val != nil {
					values = append(values, *val)
				} else {
					values = append(values, 0)
				}
			}
			ability.Variables = append(ability.Variables, ParsedAbilityVariable{
				Name:   v.Name,
				Values: values,
			})
		}

		champions = append(champions, ParsedChampion{
			APIName:       c.APIName,
			Name:          c.Name,
			Cost:          c.Cost,
			TraitAPINames: traitAPINames,
			Stats:         c.Stats,
			Ability:       ability,
			IconURL:       ConvertIconPath(c.Icon),
			SquareIconURL: ConvertIconPath(c.SquareIcon),
			TileIconURL:   ConvertIconPath(c.TileIcon),
		})
	}

	// Filter items for current set
	setPrefix := strings.ToLower(setData.Mutator)
	items := filterItems(data.Items, setPrefix, setData.Number)

	return &ParsedSet{
		Number:    setData.Number,
		Name:      setData.Name,
		Mutator:   setData.Mutator,
		Champions: champions,
		Traits:    traits,
		Items:     items,
	}
}

// buildTraitDisplayMap creates a mapping from trait display name → apiName.
func buildTraitDisplayMap(traits []CDragonTrait) map[string]string {
	m := make(map[string]string, len(traits))
	for _, t := range traits {
		m[t.Name] = t.APIName
	}
	return m
}

// resolveTraitNames converts display names to apiNames using the mapping.
func resolveTraitNames(displayNames []string, displayToAPI map[string]string) []string {
	result := make([]string, 0, len(displayNames))
	for _, name := range displayNames {
		if apiName, ok := displayToAPI[name]; ok {
			result = append(result, apiName)
		}
	}
	return result
}

// filterItems filters the global items list for items belonging to the current set.
func filterItems(allItems []CDragonItem, setMutator string, setNumber int) []ParsedItem {
	items := make([]ParsedItem, 0)
	setTag := strings.ToLower(setMutator)

	for _, item := range allItems {
		// Include base items (no set-specific tag needed) and set-specific items.
		// Items are included if they have no "Set" prefix in apiName OR if they match current set.
		apiLower := strings.ToLower(item.APIName)

		// Skip items from other sets
		if strings.HasPrefix(apiLower, "tft") {
			// Check if it belongs to a different set number
			belongsToCurrentSet := false

			// Base TFT items (e.g., "TFT_Item_BFSword") belong to all sets
			if !hasSetNumber(apiLower) {
				belongsToCurrentSet = true
			}

			// Set-specific items contain the set identifier
			if setTag != "" && strings.Contains(apiLower, strings.ToLower(setTag)) {
				belongsToCurrentSet = true
			}

			// Also check for set number in the apiName (e.g., "TFT13_...")
			setPrefix := strings.ToLower("TFT" + itoa(setNumber))
			if strings.HasPrefix(apiLower, setPrefix+"_") {
				belongsToCurrentSet = true
			}

			if !belongsToCurrentSet {
				continue
			}
		}

		items = append(items, ParsedItem{
			APIName:            item.APIName,
			Name:               item.Name,
			Desc:               item.Desc,
			Composition:        item.Composition,
			Effects:            item.Effects,
			IconURL:            ConvertIconPath(item.Icon),
			AssociatedTraits:   item.AssociatedTraits,
			IncompatibleTraits: item.IncompatibleTraits,
			Tags:               item.Tags,
			Unique:             item.Unique,
		})
	}

	return items
}

// hasSetNumber checks if an item apiName contains a set number suffix (e.g., TFT9_, TFT13_).
func hasSetNumber(apiName string) bool {
	if !strings.HasPrefix(apiName, "tft") {
		return false
	}
	// After "tft", check if next chars are digits followed by "_"
	rest := apiName[3:]
	i := 0
	for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
		i++
	}
	return i > 0 && i < len(rest) && rest[i] == '_'
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

// ConvertIconPath converts a CommunityDragon icon path to a full URL.
// CDragon paths look like: "ASSETS/Maps/Particles/TFT/Item_Icons/Standard/Deathblade.tex"
// Converted to: "https://raw.communitydragon.org/latest/game/assets/maps/particles/tft/item_icons/standard/deathblade.png"
func ConvertIconPath(path string) string {
	if path == "" {
		return ""
	}

	// Lowercase the entire path
	lower := strings.ToLower(path)

	// Replace .tex and .dds extensions with .png
	lower = strings.TrimSuffix(lower, ".tex")
	lower = strings.TrimSuffix(lower, ".dds")
	lower += ".png"

	return cdragonAssetBase + lower
}
