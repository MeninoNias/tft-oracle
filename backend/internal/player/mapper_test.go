package player

import (
	"testing"

	"github.com/MeninoNias/tft-oracle/backend/internal/riot"
)

func TestMapAccountToProto(t *testing.T) {
	dto := &riot.AccountDTO{PUUID: "abc-123", GameName: "Player", TagLine: "BR1"}
	p := mapAccountToProto(dto)
	if p.Puuid != "abc-123" || p.GameName != "Player" || p.TagLine != "BR1" {
		t.Errorf("unexpected proto: %+v", p)
	}
}

func TestMapSummonerToProto(t *testing.T) {
	dto := &riot.SummonerDTO{
		ID: "s1", PUUID: "p1", ProfileIconID: 42, SummonerLevel: 200, RevisionDate: 1700000000,
	}
	p := mapSummonerToProto(dto)
	if p.Id != "s1" || p.Puuid != "p1" || p.ProfileIconId != 42 || p.SummonerLevel != 200 {
		t.Errorf("unexpected proto: %+v", p)
	}
}

func TestMapLeagueEntriesToProto(t *testing.T) {
	entries := []riot.LeagueEntryDTO{
		{
			QueueType: "RANKED_TFT", Tier: "DIAMOND", Rank: "II",
			LeaguePoints: 45, Wins: 100, Losses: 80,
			HotStreak: true, Veteran: false,
		},
		{
			QueueType: "RANKED_TFT_TURBO", Tier: "GOLD", Rank: "I",
			LeaguePoints: 0, Wins: 10, Losses: 5,
			MiniSeries: &riot.MiniSeriesDTO{
				Progress: "WLN", Wins: 1, Losses: 1, Target: 2,
			},
		},
	}

	result := mapLeagueEntriesToProto(entries)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}

	// First entry — no mini series
	if result[0].Tier != "DIAMOND" || result[0].Rank != "II" || !result[0].HotStreak {
		t.Errorf("unexpected first entry: %+v", result[0])
	}
	if result[0].MiniSeries != nil {
		t.Error("expected nil mini series for first entry")
	}

	// Second entry — with mini series
	if result[1].MiniSeries == nil {
		t.Fatal("expected mini series for second entry")
	}
	if result[1].MiniSeries.Progress != "WLN" || result[1].MiniSeries.Target != 2 {
		t.Errorf("unexpected mini series: %+v", result[1].MiniSeries)
	}
}

func TestMapLeagueEntriesToProto_Empty(t *testing.T) {
	result := mapLeagueEntriesToProto(nil)
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d", len(result))
	}
}

func TestMapMatchToProto(t *testing.T) {
	dto := &riot.MatchDTO{
		Metadata: riot.MatchMetadataDTO{
			DataVersion:       "5",
			MatchID:           "BR1_123",
			ParticipantPUUIDs: []string{"p1", "p2"},
		},
		Info: riot.MatchInfoDTO{
			GameDatetime: 1700000000000,
			GameLength:   1800.5,
			GameVersion:  "14.1",
			TFTSetNumber: 13,
			TFTGameType:  "standard",
			Participants: []riot.ParticipantDTO{
				{
					PUUID: "p1", Placement: 1, Level: 9,
					Augments: []string{"aug1", "aug2"},
					Traits:   []riot.TraitDTO{{Name: "Set13_Trait1", NumUnits: 3, Style: 2}},
					Units:    []riot.UnitDTO{{CharacterID: "TFT13_Garen", Tier: 2, ItemNames: []string{"TFT_Item_BFSword"}}},
				},
				{
					PUUID: "p2", Placement: 8, Level: 5,
				},
			},
		},
	}

	p := mapMatchToProto(dto)

	if p.Metadata.MatchId != "BR1_123" || p.Metadata.DataVersion != "5" {
		t.Errorf("unexpected metadata: %+v", p.Metadata)
	}
	if len(p.Metadata.ParticipantPuuids) != 2 {
		t.Errorf("expected 2 puuids, got %d", len(p.Metadata.ParticipantPuuids))
	}
	if p.Info.TftSetNumber != 13 || p.Info.GameLength != 1800.5 {
		t.Errorf("unexpected info: %+v", p.Info)
	}
	if len(p.Info.Participants) != 2 {
		t.Fatalf("expected 2 participants, got %d", len(p.Info.Participants))
	}

	// First participant
	p1 := p.Info.Participants[0]
	if p1.Placement != 1 || p1.Level != 9 {
		t.Errorf("unexpected participant 1: %+v", p1)
	}
	if len(p1.Augments) != 2 || p1.Augments[0] != "aug1" {
		t.Errorf("unexpected augments: %v", p1.Augments)
	}
	if len(p1.Traits) != 1 || p1.Traits[0].ApiName != "Set13_Trait1" {
		t.Errorf("unexpected traits: %v", p1.Traits)
	}
	if len(p1.Units) != 1 || p1.Units[0].CharacterId != "TFT13_Garen" {
		t.Errorf("unexpected units: %v", p1.Units)
	}

	// Second participant — nil augments
	p2 := p.Info.Participants[1]
	if p2.Placement != 8 {
		t.Errorf("unexpected participant 2 placement: %d", p2.Placement)
	}
}

func TestMapParticipantToProto_NilAugments(t *testing.T) {
	p := &riot.ParticipantDTO{
		PUUID: "p1", Placement: 4,
		Augments: nil,
		Traits:   nil,
		Units:    nil,
	}
	result := mapParticipantToProto(p)
	if result.Augments == nil {
		t.Error("expected non-nil augments slice")
	}
	if len(result.Augments) != 0 {
		t.Errorf("expected empty augments, got %d", len(result.Augments))
	}
}
