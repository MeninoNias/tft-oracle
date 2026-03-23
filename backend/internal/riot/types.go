package riot

// AccountDTO represents a Riot Account (Account V1).
type AccountDTO struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

// SummonerDTO represents a TFT Summoner (TFT Summoner V1).
type SummonerDTO struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	ProfileIconID int32  `json:"profileIconId"`
	SummonerLevel int64  `json:"summonerLevel"`
	RevisionDate  int64  `json:"revisionDate"`
}

// LeagueEntryDTO represents a ranked league entry (TFT League V1).
type LeagueEntryDTO struct {
	QueueType    string        `json:"queueType"`
	Tier         string        `json:"tier"`
	Rank         string        `json:"rank"`
	LeaguePoints int32         `json:"leaguePoints"`
	Wins         int32         `json:"wins"`
	Losses       int32         `json:"losses"`
	HotStreak    bool          `json:"hotStreak"`
	Veteran      bool          `json:"veteran"`
	FreshBlood   bool          `json:"freshBlood"`
	Inactive     bool          `json:"inactive"`
	SummonerID   string        `json:"summonerId"`
	MiniSeries   *MiniSeriesDTO `json:"miniSeries,omitempty"`
}

// MiniSeriesDTO represents a promotion series.
type MiniSeriesDTO struct {
	Progress string `json:"progress"`
	Wins     int32  `json:"wins"`
	Losses   int32  `json:"losses"`
	Target   int32  `json:"target"`
}

// LeagueListDTO represents a league list (TFT League V1 Challenger/Grandmaster/Master).
type LeagueListDTO struct {
	Tier     string          `json:"tier"`
	LeagueID string          `json:"leagueId"`
	Queue    string          `json:"queue"`
	Name     string          `json:"name"`
	Entries  []LeagueItemDTO `json:"entries"`
}

// LeagueItemDTO represents a single entry in a league list.
type LeagueItemDTO struct {
	SummonerID   string         `json:"summonerId"`
	LeaguePoints int32          `json:"leaguePoints"`
	Rank         string         `json:"rank"`
	Wins         int32          `json:"wins"`
	Losses       int32          `json:"losses"`
	HotStreak    bool           `json:"hotStreak"`
	Veteran      bool           `json:"veteran"`
	FreshBlood   bool           `json:"freshBlood"`
	Inactive     bool           `json:"inactive"`
	MiniSeries   *MiniSeriesDTO `json:"miniSeries,omitempty"`
}

// MatchDTO represents a full TFT match (TFT Match V1).
type MatchDTO struct {
	Metadata MatchMetadataDTO `json:"metadata"`
	Info     MatchInfoDTO     `json:"info"`
}

// MatchMetadataDTO contains match metadata.
type MatchMetadataDTO struct {
	DataVersion      string   `json:"data_version"`
	MatchID          string   `json:"match_id"`
	ParticipantPUUIDs []string `json:"participants"`
}

// MatchInfoDTO contains match gameplay data.
type MatchInfoDTO struct {
	GameDatetime   int64            `json:"game_datetime"`
	GameLength     float32          `json:"game_length"`
	GameVersion    string           `json:"game_version"`
	GameVariation  string           `json:"game_variation"`
	QueueID        int32            `json:"queue_id"`
	TFTSetNumber   int32            `json:"tft_set_number"`
	TFTSetCoreName string           `json:"tft_set_core_name"`
	TFTGameType    string           `json:"tft_game_type"`
	EndOfGameResult string          `json:"endOfGameResult"`
	Participants   []ParticipantDTO `json:"participants"`
}

// ParticipantDTO represents a participant in a match.
type ParticipantDTO struct {
	PUUID                string        `json:"puuid"`
	Placement            int32         `json:"placement"`
	Level                int32         `json:"level"`
	GoldLeft             int32         `json:"gold_left"`
	LastRound            int32         `json:"last_round"`
	TimeEliminated       float32       `json:"time_eliminated"`
	TotalDamageToPlayers int32         `json:"total_damage_to_players"`
	PlayersEliminated    int32         `json:"players_eliminated"`
	PartnerGroupID       int32         `json:"partner_group_id"`
	Augments             []string      `json:"augments"`
	Companion            CompanionDTO  `json:"companion"`
	Traits               []TraitDTO    `json:"traits"`
	Units                []UnitDTO     `json:"units"`
}

// CompanionDTO represents a Little Legend.
type CompanionDTO struct {
	ContentID string `json:"content_ID"`
	Species   string `json:"species"`
	ItemID    int32  `json:"item_ID"`
	SkinID    int32  `json:"skin_ID"`
}

// TraitDTO represents a trait in a match.
type TraitDTO struct {
	Name        string `json:"name"`
	NumUnits    int32  `json:"num_units"`
	Style       int32  `json:"style"`
	TierCurrent int32  `json:"tier_current"`
	TierTotal   int32  `json:"tier_total"`
}

// UnitDTO represents a unit on a player's board.
type UnitDTO struct {
	CharacterID string   `json:"character_id"`
	Name        string   `json:"name"`
	Tier        int32    `json:"tier"`
	Rarity      int32    `json:"rarity"`
	ItemNames   []string `json:"itemNames"` // Item apiName strings (e.g. "TFT_Item_BFSword")
}
