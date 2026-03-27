package ai

import (
	"strings"
	"testing"
)

func testParticipantData(puuid string, placement int32) MatchParticipantData {
	return MatchParticipantData{
		Puuid:                puuid,
		Placement:            placement,
		Level:                8,
		GoldLeft:             12,
		LastRound:            35, // stage 6-4
		TotalDamageToPlayers: 45,
		PlayersEliminated:    1,
		Augments:             []string{"TFT_Item_BFSword"},
		Traits: []ParticipantTrait{
			{ApiName: "Set13_Warrior", NumUnits: 4, Style: 2, TierCurrent: 2, TierTotal: 3},
			{ApiName: "Set13_Rebel", NumUnits: 3, Style: 1, TierCurrent: 1, TierTotal: 2},
		},
		Units: []ParticipantUnit{
			{CharacterID: "TFT13_Garen", Name: "Garen", Tier: 2, Rarity: 0, Items: []string{"TFT_Item_GuardianAngel"}},
			{CharacterID: "TFT13_Jinx", Name: "Jinx", Tier: 2, Rarity: 2, Items: []string{"TFT_Item_InfinityEdge", "TFT_Item_BFSword"}},
		},
	}
}

func testTierList() []TierListEntry {
	return []TierListEntry{
		{
			CompositionName:  "Warrior Garen",
			Tier:             "S",
			Score:            95,
			CoreChampions:    []string{"Garen", "Vi", "Darius"},
			RecommendedItems: map[string][]string{"Garen": {"Guardian Angel", "Warmog's"}},
			WinRate:          0.52,
			AvgPlacement:     3.2,
		},
		{
			CompositionName: "Rebel Jinx",
			Tier:            "A",
			Score:           80,
			CoreChampions:   []string{"Jinx", "Ekko"},
			WinRate:         0.48,
			AvgPlacement:    3.8,
		},
	}
}

func TestBuildMatchCoachPrompt(t *testing.T) {
	data := testGameData()
	player := testParticipantData("player-1", 3)
	opponent := testParticipantData("opponent-1", 1)
	allParticipants := []MatchParticipantData{opponent, player}
	meta := MatchMeta{GameLength: 1800, LastRound: 35, GameType: "standard"}
	tierList := testTierList()

	prompt := BuildMatchCoachPrompt(player, allParticipants, meta, data, tierList, "Gold II")

	// Should contain player section
	if !strings.Contains(prompt, "YOUR MATCH RESULT") {
		t.Error("prompt should contain YOUR MATCH RESULT section")
	}
	if !strings.Contains(prompt, "Gold II") {
		t.Error("prompt should contain player rank")
	}
	if !strings.Contains(prompt, "Placement: 3") {
		t.Error("prompt should contain placement")
	}

	// Should contain champion names (resolved from GameData)
	if !strings.Contains(prompt, "Jinx") {
		t.Error("prompt should contain champion name Jinx")
	}
	if !strings.Contains(prompt, "Garen") {
		t.Error("prompt should contain champion name Garen")
	}

	// Should contain lobby section
	if !strings.Contains(prompt, "LOBBY") {
		t.Error("prompt should contain LOBBY section")
	}

	// Should contain meta context
	if !strings.Contains(prompt, "META CONTEXT") {
		t.Error("prompt should contain META CONTEXT section")
	}
	if !strings.Contains(prompt, "Warrior Garen") {
		t.Error("prompt should contain meta comp name")
	}

	// Should contain instructions
	if !strings.Contains(prompt, "INSTRUCTIONS") {
		t.Error("prompt should contain INSTRUCTIONS section")
	}
}

func TestBuildMatchCoachPrompt_NoTierList(t *testing.T) {
	data := testGameData()
	player := testParticipantData("player-1", 5)
	allParticipants := []MatchParticipantData{player}
	meta := MatchMeta{GameLength: 1500, LastRound: 30, GameType: "standard"}

	prompt := BuildMatchCoachPrompt(player, allParticipants, meta, data, nil, "Silver I")

	if strings.Contains(prompt, "META CONTEXT") {
		t.Error("prompt should NOT contain META CONTEXT when tier list is empty")
	}
	if !strings.Contains(prompt, "YOUR MATCH RESULT") {
		t.Error("prompt should still contain player section")
	}
}

func TestBuildHistoryCoachPrompt(t *testing.T) {
	data := testGameData()
	matches := []MatchSummary{
		{
			MatchID: "match-1", Placement: 2, Level: 8, GoldLeft: 10, LastRound: 35,
			Augments: []string{"TFT_Item_BFSword"},
			Units: []ParticipantUnit{
				{CharacterID: "TFT13_Jinx", Tier: 2, Rarity: 2, Items: []string{"TFT_Item_InfinityEdge"}},
				{CharacterID: "TFT13_Garen", Tier: 3, Rarity: 0},
			},
			Traits: []ParticipantTrait{
				{ApiName: "Set13_Warrior", NumUnits: 4, Style: 2},
			},
		},
		{
			MatchID: "match-2", Placement: 6, Level: 7, GoldLeft: 3, LastRound: 28,
			Units: []ParticipantUnit{
				{CharacterID: "TFT13_Garen", Tier: 2, Rarity: 0},
			},
			Traits: []ParticipantTrait{
				{ApiName: "Set13_Warrior", NumUnits: 2, Style: 1},
			},
		},
	}

	aggregates := ComputeHistoryAggregates(matches, data)
	prompt := BuildHistoryCoachPrompt(matches, aggregates, data, testTierList(), "Gold II")

	if !strings.Contains(prompt, "PLAYER PROFILE") {
		t.Error("prompt should contain PLAYER PROFILE section")
	}
	if !strings.Contains(prompt, "Matches analyzed: 2") {
		t.Error("prompt should contain match count")
	}
	if !strings.Contains(prompt, "MATCH HISTORY") {
		t.Error("prompt should contain MATCH HISTORY section")
	}
	if !strings.Contains(prompt, "AGGREGATE STATS") {
		t.Error("prompt should contain AGGREGATE STATS section")
	}
	if !strings.Contains(prompt, "META CONTEXT") {
		t.Error("prompt should contain META CONTEXT section")
	}
}

func TestComputeHistoryAggregates(t *testing.T) {
	data := testGameData()
	matches := []MatchSummary{
		{Placement: 1, Units: []ParticipantUnit{{CharacterID: "TFT13_Jinx", Tier: 2, Rarity: 2, Items: []string{"TFT_Item_InfinityEdge"}}}, Traits: []ParticipantTrait{{ApiName: "Set13_Rebel", NumUnits: 5, Style: 2}}},
		{Placement: 3, Units: []ParticipantUnit{{CharacterID: "TFT13_Jinx", Tier: 2, Rarity: 2}}, Traits: []ParticipantTrait{{ApiName: "Set13_Rebel", NumUnits: 5, Style: 2}}},
		{Placement: 5, Units: []ParticipantUnit{{CharacterID: "TFT13_Garen", Tier: 3, Rarity: 0}}, Traits: []ParticipantTrait{{ApiName: "Set13_Warrior", NumUnits: 4, Style: 2}}},
		{Placement: 7, Units: []ParticipantUnit{{CharacterID: "TFT13_Garen", Tier: 2, Rarity: 0}}, Traits: []ParticipantTrait{{ApiName: "Set13_Warrior", NumUnits: 2, Style: 1}}},
	}

	agg := ComputeHistoryAggregates(matches, data)

	// Avg placement: (1+3+5+7)/4 = 4.0
	if agg.AvgPlacement != 4.0 {
		t.Errorf("expected avg placement 4.0, got %.1f", agg.AvgPlacement)
	}
	// Top 4: 2/4 = 0.5
	if agg.Top4Rate != 0.5 {
		t.Errorf("expected top4 rate 0.5, got %.2f", agg.Top4Rate)
	}
	// Win rate: 1/4 = 0.25
	if agg.WinRate != 0.25 {
		t.Errorf("expected win rate 0.25, got %.2f", agg.WinRate)
	}
	// Placement distribution
	if agg.PlacementDist[1] != 1 || agg.PlacementDist[3] != 1 || agg.PlacementDist[5] != 1 || agg.PlacementDist[7] != 1 {
		t.Errorf("unexpected placement distribution: %v", agg.PlacementDist)
	}
	// Champion frequency: Jinx=2, Garen=2
	if len(agg.TopChampions) < 2 {
		t.Fatalf("expected at least 2 top champions, got %d", len(agg.TopChampions))
	}
	// Item frequency: Infinity Edge=1
	if len(agg.TopItems) < 1 {
		t.Fatalf("expected at least 1 top item, got %d", len(agg.TopItems))
	}
}

func TestComputeHistoryAggregates_Empty(t *testing.T) {
	data := testGameData()
	agg := ComputeHistoryAggregates(nil, data)
	if agg.AvgPlacement != 0 {
		t.Errorf("expected 0 avg placement for empty matches, got %.1f", agg.AvgPlacement)
	}
}

func TestFormatRound(t *testing.T) {
	tests := []struct {
		round    int32
		expected string
	}{
		{1, "1-1"},
		{3, "1-3"},
		{4, "2-2"},
		{10, "3-1"},
		{35, "6-5"},
	}
	for _, tt := range tests {
		got := formatRound(tt.round)
		if got != tt.expected {
			t.Errorf("formatRound(%d) = %q, want %q", tt.round, got, tt.expected)
		}
	}
}

func TestOrdinalSuffix(t *testing.T) {
	tests := []struct {
		n        int32
		expected string
	}{
		{1, "st"}, {2, "nd"}, {3, "rd"}, {4, "th"},
		{11, "th"}, {12, "th"}, {13, "th"}, {21, "st"},
	}
	for _, tt := range tests {
		got := ordinalSuffix(tt.n)
		if got != tt.expected {
			t.Errorf("ordinalSuffix(%d) = %q, want %q", tt.n, got, tt.expected)
		}
	}
}
