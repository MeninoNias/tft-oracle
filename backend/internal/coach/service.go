package coach

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
	"github.com/MeninoNias/tft-oracle/backend/gen/tft/v1/tftv1connect"
	"github.com/MeninoNias/tft-oracle/backend/internal/ai"
	"github.com/MeninoNias/tft-oracle/backend/internal/auth"
	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
)

var _ tftv1connect.CoachServiceHandler = (*Service)(nil)

type Service struct {
	db      *pgxpool.Pool
	queries *generated.Queries
	ai      *ai.Client
}

func NewService(db *pgxpool.Pool, aiClient *ai.Client) *Service {
	return &Service{
		db:      db,
		queries: generated.New(db),
		ai:      aiClient,
	}
}

// GetMatchAnalysis returns a previously cached analysis without calling AI.
func (s *Service) GetMatchAnalysis(
	ctx context.Context,
	req *connect.Request[tftv1.GetMatchAnalysisRequest],
) (*connect.Response[tftv1.GetMatchAnalysisResponse], error) {
	if _, err := auth.RequireAuth(ctx); err != nil {
		return nil, err
	}

	matchID := req.Msg.MatchId
	puuid := req.Msg.Puuid
	if matchID == "" || puuid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("match_id and puuid are required"))
	}

	cached, err := s.queries.GetCoachMatchAnalysis(ctx, generated.GetCoachMatchAnalysisParams{
		MatchID: matchID,
		Puuid:   puuid,
	})
	if err != nil {
		return connect.NewResponse(&tftv1.GetMatchAnalysisResponse{Found: false}), nil
	}

	var analysis ai.MatchCoachAnalysis
	if err := json.Unmarshal(cached.Response, &analysis); err != nil {
		return connect.NewResponse(&tftv1.GetMatchAnalysisResponse{Found: false}), nil
	}

	placement := s.getPlayerPlacement(ctx, matchID, puuid)
	return connect.NewResponse(&tftv1.GetMatchAnalysisResponse{
		Found:    true,
		Analysis: mapMatchAnalysisToProto(matchID, placement, &analysis),
	}), nil
}

// AnalyzeMatch delivers personalized coaching for a single completed match.
func (s *Service) AnalyzeMatch(
	ctx context.Context,
	req *connect.Request[tftv1.AnalyzeMatchRequest],
) (*connect.Response[tftv1.AnalyzeMatchResponse], error) {
	// Require authentication
	if _, err := auth.RequireAuth(ctx); err != nil {
		return nil, err
	}

	matchID := req.Msg.MatchId
	puuid := req.Msg.Puuid

	if matchID == "" || puuid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("match_id and puuid are required"))
	}

	if !s.ai.Available() {
		return nil, connect.NewError(connect.CodeUnavailable,
			fmt.Errorf("AI coaching is not available — OPENAI_API_KEY not configured"))
	}

	// Check cache (skip if force re-generation requested)
	if !req.Msg.Force {
		cached, err := s.queries.GetCoachMatchAnalysis(ctx, generated.GetCoachMatchAnalysisParams{
			MatchID: matchID,
			Puuid:   puuid,
		})
		if err == nil {
			var analysis ai.MatchCoachAnalysis
			if err := json.Unmarshal(cached.Response, &analysis); err == nil {
				log.Printf("coach-match: cache hit for %s/%s", matchID, puuid[:8])
				placement := s.getPlayerPlacement(ctx, matchID, puuid)
				return connect.NewResponse(mapMatchAnalysisToProto(matchID, placement, &analysis)), nil
			}
		}
	}

	// Load match data
	participants, err := s.queries.GetParticipantsByMatch(ctx, matchID)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound,
			fmt.Errorf("match not found: %w", err))
	}

	// Find the player
	var playerData ai.MatchParticipantData
	var allData []ai.MatchParticipantData
	found := false
	for _, p := range participants {
		d := mapParticipantToCoachData(p)
		allData = append(allData, d)
		if p.Puuid == puuid {
			playerData = d
			found = true
		}
	}
	if !found {
		return nil, connect.NewError(connect.CodeNotFound,
			fmt.Errorf("player %s not found in match %s", puuid, matchID))
	}

	// Load match metadata
	match, err := s.queries.GetMatchByID(ctx, matchID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("load match metadata: %w", err))
	}
	meta := ai.MatchMeta{
		GameLength: match.GameLength,
		LastRound:  playerData.LastRound,
		GameType:   match.TftGameType,
	}

	// Load enrichment data
	gameData, err := s.loadGameData(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("load game data: %w", err))
	}

	tierList := s.loadTierList(ctx)
	playerRank := s.loadPlayerRank(ctx, puuid)

	// Build prompt and call AI
	prompt := ai.BuildMatchCoachPrompt(playerData, allData, meta, gameData, tierList, playerRank)
	analysis, tokensUsed, err := s.ai.AnalyzeMatchCoaching(ctx, prompt)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("ai analysis: %w", err))
	}

	// Cache result
	responseJSON, _ := json.Marshal(analysis)
	_, cacheErr := s.queries.UpsertCoachMatchAnalysis(ctx, generated.UpsertCoachMatchAnalysisParams{
		MatchID:    matchID,
		Puuid:      puuid,
		Response:   responseJSON,
		Model:      "gpt-4o-mini",
		TokensUsed: int32(tokensUsed),
	})
	if cacheErr != nil {
		log.Printf("warning: failed to cache coach match analysis: %v", cacheErr)
	}

	return connect.NewResponse(mapMatchAnalysisToProto(matchID, playerData.Placement, analysis)), nil
}

// AnalyzeHistory identifies patterns and trends across recent matches.
func (s *Service) AnalyzeHistory(
	ctx context.Context,
	req *connect.Request[tftv1.AnalyzeHistoryRequest],
) (*connect.Response[tftv1.AnalyzeHistoryResponse], error) {
	// Require authentication
	if _, err := auth.RequireAuth(ctx); err != nil {
		return nil, err
	}

	puuid := req.Msg.Puuid
	matchCount := req.Msg.MatchCount

	if puuid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("puuid is required"))
	}

	if !s.ai.Available() {
		return nil, connect.NewError(connect.CodeUnavailable,
			fmt.Errorf("AI coaching is not available — OPENAI_API_KEY not configured"))
	}

	// Default and cap match count
	if matchCount <= 0 {
		matchCount = 20
	}
	if matchCount > 50 {
		matchCount = 50
	}

	// Load participant rows
	participantRows, err := s.queries.GetParticipantsByPUUID(ctx, generated.GetParticipantsByPUUIDParams{
		Puuid:  puuid,
		Limit:  matchCount,
		Offset: 0,
	})
	if err != nil || len(participantRows) == 0 {
		return nil, connect.NewError(connect.CodeNotFound,
			fmt.Errorf("no matches found for player"))
	}

	actualCount := int32(len(participantRows))

	// Build match IDs for cache key
	matchIDs := make([]string, 0, len(participantRows))
	for _, p := range participantRows {
		matchIDs = append(matchIDs, p.MatchID)
	}

	// Check cache
	cached, err := s.queries.GetCoachHistoryAnalysis(ctx, generated.GetCoachHistoryAnalysisParams{
		Puuid:      puuid,
		MatchCount: actualCount,
	})
	if err == nil && matchIDsMatch(cached.MatchIds, matchIDs) {
		var analysis ai.HistoryCoachAnalysis
		if err := json.Unmarshal(cached.Response, &analysis); err == nil {
			log.Printf("coach-history: cache hit for %s (n=%d)", puuid[:8], actualCount)
			return connect.NewResponse(mapHistoryAnalysisToProto(actualCount, &analysis)), nil
		}
	}

	// Build match summaries
	summaries := make([]ai.MatchSummary, 0, len(participantRows))
	for _, p := range participantRows {
		summaries = append(summaries, mapParticipantToSummary(p))
	}

	// Load enrichment data
	gameData, err := s.loadGameData(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("load game data: %w", err))
	}

	aggregates := ai.ComputeHistoryAggregates(summaries, gameData)
	tierList := s.loadTierList(ctx)
	playerRank := s.loadPlayerRank(ctx, puuid)

	// Build prompt and call AI
	prompt := ai.BuildHistoryCoachPrompt(summaries, aggregates, gameData, tierList, playerRank)
	analysis, tokensUsed, err := s.ai.AnalyzeHistoryCoaching(ctx, prompt)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("ai analysis: %w", err))
	}

	// Cache result
	responseJSON, _ := json.Marshal(analysis)
	_, cacheErr := s.queries.UpsertCoachHistoryAnalysis(ctx, generated.UpsertCoachHistoryAnalysisParams{
		Puuid:      puuid,
		MatchCount: actualCount,
		MatchIds:   matchIDs,
		Response:   responseJSON,
		Model:      "gpt-4o-mini",
		TokensUsed: int32(tokensUsed),
	})
	if cacheErr != nil {
		log.Printf("warning: failed to cache coach history analysis: %v", cacheErr)
	}

	return connect.NewResponse(mapHistoryAnalysisToProto(actualCount, analysis)), nil
}

// --- Internal helpers ---

func (s *Service) loadGameData(ctx context.Context) (*ai.GameData, error) {
	set, err := s.queries.GetLatestSet(ctx)
	if err != nil {
		return nil, fmt.Errorf("get latest set: %w", err)
	}

	champions, err := s.queries.GetChampionsBySet(ctx, set.Number)
	if err != nil {
		return nil, fmt.Errorf("get champions: %w", err)
	}

	items, err := s.queries.GetItemsBySet(ctx, set.Number)
	if err != nil {
		return nil, fmt.Errorf("get items: %w", err)
	}

	traits, err := s.queries.GetTraitsBySet(ctx, set.Number)
	if err != nil {
		return nil, fmt.Errorf("get traits: %w", err)
	}

	// Also load augments (items with "augment" tag)
	augments, err := s.queries.GetAugmentsBySet(ctx, set.Number)
	if err != nil {
		// Not critical — augments are optional
		log.Printf("warning: failed to load augments: %v", err)
	}

	data := &ai.GameData{
		Champions:      make(map[string]*tftv1.Champion, len(champions)),
		Items:          make(map[string]*tftv1.Item, len(items)+len(augments)),
		Traits:         make(map[string]*tftv1.Trait, len(traits)),
		ChampionTraits: make(map[string][]string),
	}

	for _, c := range champions {
		data.Champions[c.ApiName] = &tftv1.Champion{
			ApiName: c.ApiName,
			Name:    c.Name,
			Cost:    c.Cost,
		}

		ct, err := s.queries.GetTraitsByChampion(ctx, generated.GetTraitsByChampionParams{
			ChampionApiName: c.ApiName,
			SetNumber:       set.Number,
		})
		if err == nil {
			traitNames := make([]string, 0, len(ct))
			for _, t := range ct {
				traitNames = append(traitNames, t.TraitApiName)
			}
			data.ChampionTraits[c.ApiName] = traitNames
		}
	}

	for _, i := range items {
		data.Items[i.ApiName] = &tftv1.Item{ApiName: i.ApiName, Name: i.Name}
	}
	for _, a := range augments {
		data.Items[a.ApiName] = &tftv1.Item{ApiName: a.ApiName, Name: a.Name}
	}

	for _, t := range traits {
		data.Traits[t.ApiName] = &tftv1.Trait{
			ApiName: t.ApiName,
			Name:    t.Name,
		}
	}

	return data, nil
}

func (s *Service) loadTierList(ctx context.Context) []ai.TierListEntry {
	patch, err := s.queries.GetLatestConsolidatedPatch(ctx)
	if err != nil {
		return nil
	}

	entries, err := s.queries.GetConsolidatedTierList(ctx, patch)
	if err != nil {
		return nil
	}

	return mapTierListForPrompt(entries)
}

func (s *Service) loadPlayerRank(ctx context.Context, puuid string) string {
	player, err := s.queries.GetPlayerByPUUID(ctx, puuid)
	if err != nil {
		return "Unknown"
	}
	if player.Tier == "" {
		return "Unranked"
	}
	return fmt.Sprintf("%s %s", player.Tier, player.Rank)
}

func (s *Service) getPlayerPlacement(ctx context.Context, matchID, puuid string) int32 {
	participants, err := s.queries.GetParticipantsByMatch(ctx, matchID)
	if err != nil {
		return 0
	}
	for _, p := range participants {
		if p.Puuid == puuid {
			return p.Placement
		}
	}
	return 0
}

func matchIDsMatch(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Ensure pgx.ErrNoRows is handled (used by cache checks).
var _ = pgx.ErrNoRows
