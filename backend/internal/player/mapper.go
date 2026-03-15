package player

import (
	"encoding/json"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
	"github.com/MeninoNias/tft-oracle/backend/internal/riot"
	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
)

// --- Riot DTO → Proto ---

func mapAccountToProto(dto *riot.AccountDTO) *tftv1.Account {
	return &tftv1.Account{
		Puuid:    dto.PUUID,
		GameName: dto.GameName,
		TagLine:  dto.TagLine,
	}
}

func mapSummonerToProto(dto *riot.SummonerDTO) *tftv1.Summoner {
	return &tftv1.Summoner{
		Id:            dto.ID,
		Puuid:         dto.PUUID,
		ProfileIconId: dto.ProfileIconID,
		SummonerLevel: dto.SummonerLevel,
		RevisionDate:  dto.RevisionDate,
	}
}

func mapLeagueEntriesToProto(entries []riot.LeagueEntryDTO) []*tftv1.RankedEntry {
	result := make([]*tftv1.RankedEntry, 0, len(entries))
	for _, e := range entries {
		entry := &tftv1.RankedEntry{
			QueueType:    e.QueueType,
			Tier:         e.Tier,
			Rank:         e.Rank,
			LeaguePoints: e.LeaguePoints,
			Wins:         e.Wins,
			Losses:       e.Losses,
			HotStreak:    e.HotStreak,
			Veteran:      e.Veteran,
			FreshBlood:   e.FreshBlood,
			Inactive:     e.Inactive,
		}
		if e.MiniSeries != nil {
			entry.MiniSeries = &tftv1.MiniSeries{
				Progress: e.MiniSeries.Progress,
				Wins:     e.MiniSeries.Wins,
				Losses:   e.MiniSeries.Losses,
				Target:   e.MiniSeries.Target,
			}
		}
		result = append(result, entry)
	}
	return result
}

func mapMatchToProto(dto *riot.MatchDTO) *tftv1.Match {
	participants := make([]*tftv1.Participant, 0, len(dto.Info.Participants))
	for _, p := range dto.Info.Participants {
		participants = append(participants, mapParticipantToProto(&p))
	}

	return &tftv1.Match{
		Metadata: &tftv1.MatchMetadata{
			DataVersion:      dto.Metadata.DataVersion,
			MatchId:          dto.Metadata.MatchID,
			ParticipantPuuids: dto.Metadata.ParticipantPUUIDs,
		},
		Info: &tftv1.MatchInfo{
			GameDatetime:   dto.Info.GameDatetime,
			GameLength:     dto.Info.GameLength,
			GameVersion:    dto.Info.GameVersion,
			GameVariation:  dto.Info.GameVariation,
			QueueId:        dto.Info.QueueID,
			TftSetNumber:   dto.Info.TFTSetNumber,
			TftSetCoreName: dto.Info.TFTSetCoreName,
			TftGameType:    dto.Info.TFTGameType,
			EndOfGameResult: dto.Info.EndOfGameResult,
			Participants:   participants,
		},
	}
}

func mapParticipantToProto(p *riot.ParticipantDTO) *tftv1.Participant {
	traits := make([]*tftv1.MatchTrait, 0, len(p.Traits))
	for _, t := range p.Traits {
		traits = append(traits, &tftv1.MatchTrait{
			ApiName:     t.Name,
			NumUnits:    t.NumUnits,
			Style:       t.Style,
			TierCurrent: t.TierCurrent,
			TierTotal:   t.TierTotal,
		})
	}

	units := make([]*tftv1.MatchUnit, 0, len(p.Units))
	for _, u := range p.Units {
		units = append(units, &tftv1.MatchUnit{
			CharacterId: u.CharacterID,
			Name:        u.Name,
			Tier:        u.Tier,
			Rarity:      u.Rarity,
			ItemNames:   u.ItemNames,
		})
	}

	augments := p.Augments
	if augments == nil {
		augments = []string{}
	}

	return &tftv1.Participant{
		Puuid:                p.PUUID,
		Placement:            p.Placement,
		Level:                p.Level,
		GoldLeft:             p.GoldLeft,
		LastRound:            p.LastRound,
		TimeEliminated:       p.TimeEliminated,
		TotalDamageToPlayers: p.TotalDamageToPlayers,
		PlayersEliminated:    p.PlayersEliminated,
		PartnerGroupId:       p.PartnerGroupID,
		Augments:             augments,
		Companion: &tftv1.Companion{
			ContentId: p.Companion.ContentID,
			Species:   p.Companion.Species,
			ItemId:    p.Companion.ItemID,
			SkinId:    p.Companion.SkinID,
		},
		Traits: traits,
		Units:  units,
	}
}

// --- Riot DTO → JSONB (for DB storage) ---

type dbTrait struct {
	Name        string `json:"name"`
	NumUnits    int32  `json:"num_units"`
	Style       int32  `json:"style"`
	TierCurrent int32  `json:"tier_current"`
	TierTotal   int32  `json:"tier_total"`
}

type dbUnit struct {
	CharacterID string   `json:"character_id"`
	Name        string   `json:"name"`
	Tier        int32    `json:"tier"`
	Rarity      int32    `json:"rarity"`
	ItemIDs     []int32  `json:"item_ids"`
	Items       []string `json:"items"`
}

type dbCompanion struct {
	ContentID string `json:"content_id"`
	Species   string `json:"species"`
	ItemID    int32  `json:"item_id"`
	SkinID    int32  `json:"skin_id"`
}

func traitsToJSON(traits []riot.TraitDTO) ([]byte, error) {
	dbTraits := make([]dbTrait, 0, len(traits))
	for _, t := range traits {
		dbTraits = append(dbTraits, dbTrait{
			Name:        t.Name,
			NumUnits:    t.NumUnits,
			Style:       t.Style,
			TierCurrent: t.TierCurrent,
			TierTotal:   t.TierTotal,
		})
	}
	return json.Marshal(dbTraits)
}

func unitsToJSON(units []riot.UnitDTO) ([]byte, error) {
	dbUnits := make([]dbUnit, 0, len(units))
	for _, u := range units {
		itemNames := u.ItemNames
		if itemNames == nil {
			itemNames = []string{}
		}
		dbUnits = append(dbUnits, dbUnit{
			CharacterID: u.CharacterID,
			Name:        u.Name,
			Tier:        u.Tier,
			Rarity:      u.Rarity,
			ItemIDs:     []int32{},
			Items:       itemNames,
		})
	}
	return json.Marshal(dbUnits)
}

func companionToJSON(c riot.CompanionDTO) ([]byte, error) {
	return json.Marshal(dbCompanion{
		ContentID: c.ContentID,
		Species:   c.Species,
		ItemID:    c.ItemID,
		SkinID:    c.SkinID,
	})
}

// --- DB → Proto (reconstructing match from stored data) ---

func mapMatchFromDB(match generated.Match, participants []generated.MatchParticipant) *tftv1.Match {
	pbParticipants := make([]*tftv1.Participant, 0, len(participants))
	for _, p := range participants {
		pbParticipants = append(pbParticipants, mapParticipantFromDB(p))
	}

	return &tftv1.Match{
		Metadata: &tftv1.MatchMetadata{
			DataVersion:      match.DataVersion,
			MatchId:          match.MatchID,
			ParticipantPuuids: match.ParticipantPuuids,
		},
		Info: &tftv1.MatchInfo{
			GameDatetime: match.GameDatetime,
			GameLength:   match.GameLength,
			GameVersion:  match.GameVersion,
			QueueId:      match.QueueID,
			TftSetNumber: match.TftSetNumber,
			TftGameType:  match.TftGameType,
			Participants: pbParticipants,
		},
	}
}

func mapParticipantFromDB(p generated.MatchParticipant) *tftv1.Participant {
	var traits []dbTrait
	_ = json.Unmarshal(p.Traits, &traits)

	var units []dbUnit
	_ = json.Unmarshal(p.Units, &units)

	var comp dbCompanion
	_ = json.Unmarshal(p.Companion, &comp)

	pbTraits := make([]*tftv1.MatchTrait, 0, len(traits))
	for _, t := range traits {
		pbTraits = append(pbTraits, &tftv1.MatchTrait{
			ApiName:     t.Name,
			NumUnits:    t.NumUnits,
			Style:       t.Style,
			TierCurrent: t.TierCurrent,
			TierTotal:   t.TierTotal,
		})
	}

	pbUnits := make([]*tftv1.MatchUnit, 0, len(units))
	for _, u := range units {
		pbUnits = append(pbUnits, &tftv1.MatchUnit{
			CharacterId: u.CharacterID,
			Name:        u.Name,
			Tier:        u.Tier,
			Rarity:      u.Rarity,
			ItemIds:     u.ItemIDs,
			ItemNames:   u.Items,
		})
	}

	augments := p.Augments
	if augments == nil {
		augments = []string{}
	}

	return &tftv1.Participant{
		Puuid:                p.Puuid,
		Placement:            p.Placement,
		Level:                p.Level,
		GoldLeft:             p.GoldLeft,
		LastRound:            p.LastRound,
		TimeEliminated:       p.TimeEliminated,
		TotalDamageToPlayers: p.TotalDamageToPlayers,
		PlayersEliminated:    p.PlayersEliminated,
		Augments:             augments,
		Companion: &tftv1.Companion{
			ContentId: comp.ContentID,
			Species:   comp.Species,
			ItemId:    comp.ItemID,
			SkinId:    comp.SkinID,
		},
		Traits: pbTraits,
		Units:  pbUnits,
	}
}
