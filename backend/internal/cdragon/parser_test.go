package cdragon

import (
	"testing"
)

func TestConvertIconPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tex extension",
			input:    "ASSETS/Maps/Particles/TFT/Item_Icons/Standard/Deathblade.tex",
			expected: "https://raw.communitydragon.org/latest/game/assets/maps/particles/tft/item_icons/standard/deathblade.png",
		},
		{
			name:     "dds extension",
			input:    "ASSETS/Characters/TFT13_Garen/HUD/TFT13_Garen_Square.dds",
			expected: "https://raw.communitydragon.org/latest/game/assets/characters/tft13_garen/hud/tft13_garen_square.png",
		},
		{
			name:     "empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "already lowercase",
			input:    "assets/maps/tft/icons/item.tex",
			expected: "https://raw.communitydragon.org/latest/game/assets/maps/tft/icons/item.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertIconPath(tt.input)
			if got != tt.expected {
				t.Errorf("ConvertIconPath(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsAlternateMode(t *testing.T) {
	tests := []struct {
		mutator  string
		expected bool
	}{
		{"TFTSet13", false},
		{"TFTSet13_PAIRS", true},
		{"TFTSet12_TURBO", true},
		{"TFTSet11_PVEMODE", true},
		{"", false},
		{"TFTSet13_Stage2", false},
	}

	for _, tt := range tests {
		t.Run(tt.mutator, func(t *testing.T) {
			got := isAlternateMode(tt.mutator)
			if got != tt.expected {
				t.Errorf("isAlternateMode(%q) = %v, want %v", tt.mutator, got, tt.expected)
			}
		})
	}
}

func TestFindCurrentSet(t *testing.T) {
	data := &CDragonData{
		SetData: []CDragonSetData{
			{Number: 12, Name: "Set 12", Mutator: "TFTSet12"},
			{Number: 12, Name: "Set 12 Turbo", Mutator: "TFTSet12_TURBO"},
			{Number: 13, Name: "Set 13", Mutator: "TFTSet13"},
			{Number: 13, Name: "Set 13 Pairs", Mutator: "TFTSet13_PAIRS"},
		},
	}

	current := FindCurrentSet(data)
	if current == nil {
		t.Fatal("expected to find current set, got nil")
	}
	if current.Number != 13 {
		t.Errorf("expected set 13, got %d", current.Number)
	}
	if current.Mutator != "TFTSet13" {
		t.Errorf("expected mutator TFTSet13, got %s", current.Mutator)
	}
}

func TestResolveTraitNames(t *testing.T) {
	displayToAPI := map[string]string{
		"Warrior":  "Set13_Warrior",
		"Sorcerer": "Set13_Sorcerer",
		"Guardian": "Set13_Guardian",
	}

	tests := []struct {
		name         string
		displayNames []string
		expectedLen  int
	}{
		{
			name:         "all found",
			displayNames: []string{"Warrior", "Sorcerer"},
			expectedLen:  2,
		},
		{
			name:         "some not found",
			displayNames: []string{"Warrior", "Unknown"},
			expectedLen:  1,
		},
		{
			name:         "empty",
			displayNames: []string{},
			expectedLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveTraitNames(tt.displayNames, displayToAPI)
			if len(got) != tt.expectedLen {
				t.Errorf("resolveTraitNames() returned %d items, want %d", len(got), tt.expectedLen)
			}
		})
	}
}

func TestIsPlayableChampion(t *testing.T) {
	tests := []struct {
		name     string
		cost     int
		traits   []string
		expected bool
	}{
		{"normal 1-cost with traits", 1, []string{"Set16_Warrior"}, true},
		{"normal 5-cost with traits", 5, []string{"Set16_Sorcerer", "Set16_Guardian"}, true},
		{"anvil cost 8", 8, []string{}, false},
		{"summon cost 11", 11, []string{}, false},
		{"monster cost 1 no traits", 1, []string{}, false},
		{"dummy cost 1 no traits", 1, []string{}, false},
		{"cost 0 excluded", 0, []string{"Set16_Warrior"}, false},
		{"cost 6 excluded", 6, []string{"Set16_Warrior"}, false},
		{"negative cost excluded", -1, []string{"Set16_Warrior"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPlayableChampion(tt.cost, tt.traits)
			if got != tt.expected {
				t.Errorf("isPlayableChampion(%d, %v) = %v, want %v", tt.cost, tt.traits, got, tt.expected)
			}
		})
	}
}

func TestParseFiltersNonPlayableUnits(t *testing.T) {
	data := &CDragonData{
		SetData: []CDragonSetData{
			{
				Number:  16,
				Name:    "Set 16",
				Mutator: "TFTSet16",
				Traits: []CDragonTrait{
					{APIName: "Set16_Warrior", Name: "Warrior"},
					{APIName: "Set16_Sorcerer", Name: "Sorcerer"},
				},
				Champions: []CDragonChampion{
					// Playable champions
					{APIName: "TFT16_Garen", Name: "Garen", Cost: 1, Traits: []string{"Warrior"}},
					{APIName: "TFT16_Lux", Name: "Lux", Cost: 5, Traits: []string{"Sorcerer"}},
					// PvE monster (cost 1, no traits)
					{APIName: "TFT_ElderDragon", Name: "Elder Dragon", Cost: 1, Traits: []string{}},
					// PvE monster from old set (cost 1, traits won't resolve)
					{APIName: "TFT9_SLIME_Crab", Name: "Rift Scuttler", Cost: 1, Traits: []string{"OldSetTrait"}},
					// Training dummy (cost 1, no traits)
					{APIName: "TFT_TrainingDummy", Name: "Training Dummy", Cost: 1, Traits: []string{}},
					// Item anvil (cost 8)
					{APIName: "TFT_ArmoryKeyOrnn", Name: "Artifact Item Anvil", Cost: 8, Traits: []string{}},
					// Summon (cost 11)
					{APIName: "TFT16_AnnieTibbers", Name: "Tibbers", Cost: 11, Traits: []string{}},
					// Trait mechanic summon (cost 1, no traits)
					{APIName: "TFT16_PiltoverInvention", Name: "Piltover Invention", Cost: 1, Traits: []string{}},
				},
			},
		},
	}

	parsed := Parse(data)
	if parsed == nil {
		t.Fatal("expected parsed set, got nil")
	}

	if len(parsed.Champions) != 2 {
		names := make([]string, len(parsed.Champions))
		for i, c := range parsed.Champions {
			names[i] = c.APIName
		}
		t.Errorf("expected 2 playable champions, got %d: %v", len(parsed.Champions), names)
	}

	// Verify only real champions remain
	for _, c := range parsed.Champions {
		if c.APIName != "TFT16_Garen" && c.APIName != "TFT16_Lux" {
			t.Errorf("unexpected champion in results: %s", c.APIName)
		}
	}
}

func TestHasSetNumber(t *testing.T) {
	tests := []struct {
		apiName  string
		expected bool
	}{
		{"tft_item_bfsword", false},
		{"tft13_item_something", true},
		{"tft9_item_something", true},
		{"somethingelse", false},
		{"tft_", false},
	}

	for _, tt := range tests {
		t.Run(tt.apiName, func(t *testing.T) {
			got := hasSetNumber(tt.apiName)
			if got != tt.expected {
				t.Errorf("hasSetNumber(%q) = %v, want %v", tt.apiName, got, tt.expected)
			}
		})
	}
}
