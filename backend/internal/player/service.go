package player

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
	"github.com/MeninoNias/tft-oracle/backend/gen/tft/v1/tftv1connect"
	"github.com/MeninoNias/tft-oracle/backend/internal/cache"
	"github.com/MeninoNias/tft-oracle/backend/internal/riot"
	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
)

var _ tftv1connect.PlayerServiceHandler = (*Service)(nil)

type Service struct {
	db      *pgxpool.Pool
	queries *generated.Queries
	riot    riot.RiotAPI
	cache   *cache.Client // can be nil
}

func NewService(db *pgxpool.Pool, riotClient riot.RiotAPI, cacheClient *cache.Client) *Service {
	return &Service{
		db:      db,
		queries: generated.New(db),
		riot:    riotClient,
		cache:   cacheClient,
	}
}

func (s *Service) GetPlayerProfile(
	ctx context.Context,
	req *connect.Request[tftv1.GetPlayerProfileRequest],
) (*connect.Response[tftv1.GetPlayerProfileResponse], error) {
	gameName := req.Msg.GetGameName()
	tagLine := req.Msg.GetTagLine()
	region := req.Msg.GetRegion()

	if gameName == "" || tagLine == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game_name and tag_line are required"))
	}
	if region == "" {
		region = "br"
	}

	// Resolve region → routing region + platform
	server := riot.ResolveServer(region)

	if !s.riot.Available() {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("player features unavailable: RIOT_API_KEY not configured"))
	}

	// 1. Resolve PUUID (cache → Riot API)
	puuid, err := s.cache.GetPUUID(ctx, gameName, tagLine)
	if err != nil {
		account, err := s.riot.GetAccountByRiotID(ctx, server.Region, gameName, tagLine)
		if err != nil {
			return nil, fmt.Errorf("get account: %w", err)
		}
		puuid = account.PUUID
		s.cache.SetPUUID(ctx, gameName, tagLine, puuid)
	}

	// 3. Check profile cache
	if cached, err := s.cache.GetPlayerProfile(ctx, puuid); err == nil {
		var resp tftv1.GetPlayerProfileResponse
		if json.Unmarshal(cached, &resp) == nil {
			return connect.NewResponse(&resp), nil
		}
	}

	// 4. Fetch summoner + league in parallel (both use platform routing)
	var summoner *riot.SummonerDTO
	var entries []riot.LeagueEntryDTO

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		summoner, err = s.riot.GetSummonerByPUUID(gCtx, server.Platform, puuid)
		return err
	})
	g.Go(func() error {
		var err error
		entries, err = s.riot.GetLeagueByPUUID(gCtx, server.Platform, puuid)
		return err
	})
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("fetch profile: %w", err)
	}

	// 6. Build response
	resp := &tftv1.GetPlayerProfileResponse{
		Account: &tftv1.Account{
			Puuid:    puuid,
			GameName: gameName,
			TagLine:  tagLine,
		},
		Summoner:      mapSummonerToProto(summoner),
		RankedEntries: mapLeagueEntriesToProto(entries),
	}

	// 6. Cache profile
	if data, err := json.Marshal(resp); err == nil {
		s.cache.SetPlayerProfile(ctx, puuid, data)
	}

	// 7. Persist player in DB (best effort)
	tier, rank, lp, wins, losses := extractRankedData(entries)
	if _, err := s.queries.UpsertPlayer(ctx, generated.UpsertPlayerParams{
		Puuid:         puuid,
		GameName:      gameName,
		TagLine:       tagLine,
		Region:        server.Region,
		Platform:      server.Platform,
		SummonerID:    summoner.ID,
		ProfileIconID: summoner.ProfileIconID,
		SummonerLevel: summoner.SummonerLevel,
		Tier:          tier,
		Rank:          rank,
		LeaguePoints:  lp,
		Wins:          wins,
		Losses:        losses,
	}); err != nil {
		log.Printf("warning: upsert player failed: %v", err)
	}

	return connect.NewResponse(resp), nil
}

func (s *Service) GetMatchHistory(
	ctx context.Context,
	req *connect.Request[tftv1.GetMatchHistoryRequest],
) (*connect.Response[tftv1.GetMatchHistoryResponse], error) {
	puuid := req.Msg.GetPuuid()
	region := req.Msg.GetRegion()
	count := req.Msg.GetCount()
	start := req.Msg.GetStart()

	if puuid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("puuid is required"))
	}
	if region == "" {
		region = "br"
	}
	// Resolve to routing region for match API
	server := riot.ResolveServer(region)
	region = server.Region

	if count == 0 {
		count = 20
	}
	if count > 200 {
		count = 200
	}

	if !s.riot.Available() {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("player features unavailable: RIOT_API_KEY not configured"))
	}

	// 1. Get match IDs (cache → Riot API)
	// Only use cache if it has enough IDs for the requested count
	matchIDs, err := s.cache.GetMatchIDs(ctx, puuid)
	if err != nil || int32(len(matchIDs)) < count {
		matchIDs, err = s.riot.GetMatchIDsByPUUID(ctx, region, puuid, count, start)
		if err != nil {
			return nil, fmt.Errorf("get match ids: %w", err)
		}
		s.cache.SetMatchIDs(ctx, puuid, matchIDs)
	}

	// 2. Fetch each match (cache → DB → Riot API)
	matches := make([]*tftv1.Match, 0, len(matchIDs))
	for _, matchID := range matchIDs {
		match, err := s.fetchMatch(ctx, region, matchID)
		if err != nil {
			log.Printf("warning: failed to fetch match %s: %v", matchID, err)
			continue
		}
		matches = append(matches, match)
	}

	return connect.NewResponse(&tftv1.GetMatchHistoryResponse{
		Matches: matches,
	}), nil
}

// fetchMatch tries cache → DB → Riot API, persisting new matches.
func (s *Service) fetchMatch(ctx context.Context, region, matchID string) (*tftv1.Match, error) {
	// Try Redis cache
	if cached, err := s.cache.GetMatchDetail(ctx, matchID); err == nil {
		var dto riot.MatchDTO
		if json.Unmarshal(cached, &dto) == nil {
			return mapMatchToProto(&dto), nil
		}
	}

	// Try DB
	exists, err := s.queries.MatchExists(ctx, matchID)
	if err == nil && exists {
		return s.loadMatchFromDB(ctx, matchID)
	}

	// Fetch from Riot API
	dto, err := s.riot.GetMatch(ctx, region, matchID)
	if err != nil {
		return nil, fmt.Errorf("fetch match: %w", err)
	}

	// Cache in Redis
	if data, err := json.Marshal(dto); err == nil {
		s.cache.SetMatchDetail(ctx, matchID, data)
	}

	// Persist in DB
	if err := s.persistMatch(ctx, dto); err != nil {
		log.Printf("warning: persist match %s failed: %v", matchID, err)
	}

	return mapMatchToProto(dto), nil
}

func (s *Service) loadMatchFromDB(ctx context.Context, matchID string) (*tftv1.Match, error) {
	match, err := s.queries.GetMatchByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("get match from db: %w", err)
	}

	participants, err := s.queries.GetParticipantsByMatch(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("get participants from db: %w", err)
	}

	return mapMatchFromDB(match, participants), nil
}

func (s *Service) persistMatch(ctx context.Context, dto *riot.MatchDTO) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	puuids := dto.Metadata.ParticipantPUUIDs
	if puuids == nil {
		puuids = []string{}
	}

	_, err = qtx.UpsertMatch(ctx, generated.UpsertMatchParams{
		MatchID:           dto.Metadata.MatchID,
		DataVersion:       dto.Metadata.DataVersion,
		GameDatetime:      dto.Info.GameDatetime,
		GameLength:        dto.Info.GameLength,
		GameVersion:       dto.Info.GameVersion,
		QueueID:           dto.Info.QueueID,
		TftSetNumber:      dto.Info.TFTSetNumber,
		TftGameType:       dto.Info.TFTGameType,
		ParticipantPuuids: puuids,
	})
	// UpsertMatch uses DO NOTHING — no row returned on conflict is expected (pgx.ErrNoRows)
	if err != nil && err.Error() != "no rows in result set" {
		return fmt.Errorf("upsert match: %w", err)
	}

	for _, p := range dto.Info.Participants {
		traitsJSON, err := traitsToJSON(p.Traits)
		if err != nil {
			return fmt.Errorf("marshal traits: %w", err)
		}

		unitsJSON, err := unitsToJSON(p.Units)
		if err != nil {
			return fmt.Errorf("marshal units: %w", err)
		}

		companionJSON, err := companionToJSON(p.Companion)
		if err != nil {
			return fmt.Errorf("marshal companion: %w", err)
		}

		augments := p.Augments
		if augments == nil {
			augments = []string{}
		}

		err = qtx.UpsertMatchParticipant(ctx, generated.UpsertMatchParticipantParams{
			MatchID:              dto.Metadata.MatchID,
			Puuid:                p.PUUID,
			Placement:            p.Placement,
			Level:                p.Level,
			GoldLeft:             p.GoldLeft,
			LastRound:            p.LastRound,
			TimeEliminated:       p.TimeEliminated,
			TotalDamageToPlayers: p.TotalDamageToPlayers,
			PlayersEliminated:    p.PlayersEliminated,
			Augments:             augments,
			Companion:            companionJSON,
			Traits:               traitsJSON,
			Units:                unitsJSON,
		})
		if err != nil {
			return fmt.Errorf("upsert participant: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// extractRankedData finds the RANKED_TFT entry and returns tier/rank/lp/wins/losses.
func extractRankedData(entries []riot.LeagueEntryDTO) (string, string, int32, int32, int32) {
	for _, e := range entries {
		if e.QueueType == "RANKED_TFT" {
			return e.Tier, e.Rank, e.LeaguePoints, e.Wins, e.Losses
		}
	}
	return "", "", 0, 0, 0
}
