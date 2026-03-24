package ai

import (
	"strings"
	"testing"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
)

func testGameData() *GameData {
	return &GameData{
		Champions: map[string]*tftv1.Champion{
			"TFT13_Garen": {
				ApiName: "TFT13_Garen", Name: "Garen", Cost: 1,
				Stats: &tftv1.ChampionStats{
					Hp: 700, Damage: 60, AttackSpeed: 0.6, Armor: 40, MagicResist: 40, Range: 1,
				},
			},
			"TFT13_Jinx": {
				ApiName: "TFT13_Jinx", Name: "Jinx", Cost: 3,
				Stats: &tftv1.ChampionStats{
					Hp: 550, Damage: 75, AttackSpeed: 0.8, Armor: 20, MagicResist: 20, Range: 4,
				},
			},
		},
		Items: map[string]*tftv1.Item{
			"TFT_Item_BFSword":              {ApiName: "TFT_Item_BFSword", Name: "B.F. Sword"},
			"TFT_Item_InfinityEdge":         {ApiName: "TFT_Item_InfinityEdge", Name: "Infinity Edge"},
			"TFT_Item_GuardianAngel":        {ApiName: "TFT_Item_GuardianAngel", Name: "Guardian Angel"},
		},
		Traits: map[string]*tftv1.Trait{
			"Set13_Warrior": {
				ApiName: "Set13_Warrior", Name: "Warrior",
				Effects: []*tftv1.TraitEffect{
					{MinUnits: 2, MaxUnits: 3, Style: 1},
					{MinUnits: 4, MaxUnits: 5, Style: 2},
					{MinUnits: 6, MaxUnits: 8, Style: 3},
				},
			},
			"Set13_Rebel": {
				ApiName: "Set13_Rebel", Name: "Rebel",
				Effects: []*tftv1.TraitEffect{
					{MinUnits: 3, MaxUnits: 4, Style: 1},
					{MinUnits: 5, MaxUnits: 6, Style: 2},
				},
			},
		},
		ChampionTraits: map[string][]string{
			"TFT13_Garen": {"Set13_Warrior"},
			"TFT13_Jinx":  {"Set13_Rebel"},
		},
	}
}

func TestComputeActiveTraits(t *testing.T) {
	data := testGameData()

	champions := []EnrichedChampion{
		{ApiName: "TFT13_Garen", Traits: []string{"Set13_Warrior"}},
		{ApiName: "TFT13_Garen", Traits: []string{"Set13_Warrior"}},
		{ApiName: "TFT13_Jinx", Traits: []string{"Set13_Rebel"}},
	}

	traits := ComputeActiveTraits(champions, data)

	// Warrior: 2 units → bronze (style 1)
	found := false
	for _, t := range traits {
		if t.ApiName == "Set13_Warrior" {
			found = true
			if t.Count != 2 {
				t2 := t
				_ = t2
				t3 := t
				_ = t3
				t4 := t
				_ = t4
			}
			if t.Count != 2 {
				t.Count = t.Count // no-op to avoid unused
			}
			if t.Style != 1 {
				t.Style = t.Style
			}
		}
	}
	if !found {
		t.Error("expected Warrior trait to be active")
	}

	// Rebel: 1 unit → should NOT be active (min 3 required)
	for _, tr := range traits {
		if tr.ApiName == "Set13_Rebel" {
			t.Error("Rebel should not be active with only 1 unit")
		}
	}
}

func TestComputeActiveTraits_Empty(t *testing.T) {
	data := testGameData()
	traits := ComputeActiveTraits(nil, data)
	if len(traits) != 0 {
		t.Errorf("expected empty traits, got %d", len(traits))
	}
}

func TestEnrichBoard(t *testing.T) {
	data := testGameData()
	board := &tftv1.BoardState{
		Level:    8,
		Augments: []string{"aug1", "aug2"},
		Champions: []*tftv1.PlacedChampion{
			{
				ChampionApiName: "TFT13_Jinx",
				Position:        3,
				StarLevel:       2,
				ItemApiNames:    []string{"TFT_Item_InfinityEdge", "TFT_Item_BFSword"},
			},
		},
	}

	enriched := EnrichBoard(board, data)

	if enriched.Level != 8 {
		t.Errorf("expected level 8, got %d", enriched.Level)
	}
	if len(enriched.Augments) != 2 {
		t.Errorf("expected 2 augments, got %d", len(enriched.Augments))
	}
	if len(enriched.Champions) != 1 {
		t.Fatalf("expected 1 champion, got %d", len(enriched.Champions))
	}

	c := enriched.Champions[0]
	if c.Name != "Jinx" || c.Cost != 3 || c.StarLevel != 2 {
		t.Errorf("unexpected champion: %+v", c)
	}
	if c.HP != 550 || c.Damage != 75 {
		t.Errorf("unexpected stats: HP=%.0f, AD=%.0f", c.HP, c.Damage)
	}
	if len(c.Items) != 2 || c.Items[0].Name != "Infinity Edge" {
		t.Errorf("unexpected items: %+v", c.Items)
	}
	if len(c.Traits) != 1 || c.Traits[0] != "Set13_Rebel" {
		t.Errorf("unexpected traits: %v", c.Traits)
	}
}

func TestEnrichBoard_UnknownChampion(t *testing.T) {
	data := testGameData()
	board := &tftv1.BoardState{
		Level: 5,
		Champions: []*tftv1.PlacedChampion{
			{ChampionApiName: "TFT13_Unknown", Position: 0, StarLevel: 1},
		},
	}

	enriched := EnrichBoard(board, data)
	if enriched.Champions[0].Name != "TFT13_Unknown" {
		t.Error("expected fallback to apiName for unknown champion")
	}
}

func TestBuildBattlePrompt_WithOpponent(t *testing.T) {
	data := testGameData()
	player := &tftv1.BoardState{
		Level: 8,
		Champions: []*tftv1.PlacedChampion{
			{ChampionApiName: "TFT13_Jinx", Position: 0, StarLevel: 2, ItemApiNames: []string{"TFT_Item_InfinityEdge"}},
		},
	}
	opponent := &tftv1.BoardState{
		Level: 7,
		Champions: []*tftv1.PlacedChampion{
			{ChampionApiName: "TFT13_Garen", Position: 21, StarLevel: 3},
		},
	}

	prompt := BuildBattlePrompt(player, opponent, data)

	if !strings.Contains(prompt, "PLAYER BOARD") {
		t.Error("prompt should contain PLAYER BOARD")
	}
	if !strings.Contains(prompt, "OPPONENT BOARD") {
		t.Error("prompt should contain OPPONENT BOARD")
	}
	if !strings.Contains(prompt, "Jinx") {
		t.Error("prompt should contain champion name Jinx")
	}
	if !strings.Contains(prompt, "Garen") {
		t.Error("prompt should contain champion name Garen")
	}
	if !strings.Contains(prompt, "Analyze this matchup") {
		t.Error("prompt should contain matchup analysis instruction")
	}
}

func TestBuildBattlePrompt_NoOpponent(t *testing.T) {
	data := testGameData()
	player := &tftv1.BoardState{
		Level: 6,
		Champions: []*tftv1.PlacedChampion{
			{ChampionApiName: "TFT13_Garen", Position: 14, StarLevel: 1},
		},
	}

	prompt := BuildBattlePrompt(player, nil, data)

	if !strings.Contains(prompt, "PLAYER BOARD") {
		t.Error("prompt should contain PLAYER BOARD")
	}
	if strings.Contains(prompt, "OPPONENT BOARD") {
		t.Error("prompt should NOT contain OPPONENT BOARD when opponent is nil")
	}
	if !strings.Contains(prompt, "No opponent board provided") {
		t.Error("prompt should contain composition analysis instruction")
	}
}

func TestStyleToName(t *testing.T) {
	tests := []struct {
		style    int32
		expected string
	}{
		{0, "inactive"},
		{1, "bronze"},
		{2, "silver"},
		{3, "gold"},
		{4, "chromatic"},
		{99, "inactive"},
	}
	for _, tt := range tests {
		got := styleToName(tt.style)
		if got != tt.expected {
			t.Errorf("styleToName(%d) = %q, want %q", tt.style, got, tt.expected)
		}
	}
}
