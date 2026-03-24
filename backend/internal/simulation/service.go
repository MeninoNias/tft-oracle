package simulation

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
	"github.com/MeninoNias/tft-oracle/backend/gen/tft/v1/tftv1connect"
	"github.com/MeninoNias/tft-oracle/backend/internal/ai"
	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
)

var _ tftv1connect.SimulationServiceHandler = (*Service)(nil)

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

func (s *Service) SimulateBattle(
	ctx context.Context,
	req *connect.Request[tftv1.SimulateBattleRequest],
) (*connect.Response[tftv1.SimulateBattleResponse], error) {
	// Validate request
	if req.Msg.PlayerBoard == nil || len(req.Msg.PlayerBoard.Champions) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("player board must have at least one champion"))
	}

	if !s.ai.Available() {
		return nil, connect.NewError(connect.CodeUnavailable,
			fmt.Errorf("AI simulation is not available — OPENAI_API_KEY not configured"))
	}

	// Fetch game data for enrichment
	gameData, err := s.loadGameData(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("load game data: %w", err))
	}

	// Build prompt
	prompt := ai.BuildBattlePrompt(req.Msg.PlayerBoard, req.Msg.OpponentBoard, gameData)

	// Call AI
	analysis, err := s.ai.AnalyzeBattle(ctx, prompt)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("ai analysis: %w", err))
	}

	// Map to proto response
	return connect.NewResponse(mapAnalysisToProto(analysis)), nil
}

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

	// Build lookup maps
	data := &ai.GameData{
		Champions:      make(map[string]*tftv1.Champion, len(champions)),
		Items:          make(map[string]*tftv1.Item, len(items)),
		Traits:         make(map[string]*tftv1.Trait, len(traits)),
		ChampionTraits: make(map[string][]string),
	}

	for _, c := range champions {
		data.Champions[c.ApiName] = &tftv1.Champion{
			ApiName: c.ApiName,
			Name:    c.Name,
			Cost:    c.Cost,
			Stats:   mapStatsFromDB(c.Stats),
		}

		// Fetch traits for each champion
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
		data.Items[i.ApiName] = &tftv1.Item{
			ApiName: i.ApiName,
			Name:    i.Name,
		}
	}

	for _, t := range traits {
		data.Traits[t.ApiName] = &tftv1.Trait{
			ApiName: t.ApiName,
			Name:    t.Name,
			Effects: mapTraitEffectsFromDB(t.Effects),
		}
	}

	return data, nil
}
