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
