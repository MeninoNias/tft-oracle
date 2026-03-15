package patch

import (
	"context"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
	"github.com/MeninoNias/tft-oracle/backend/gen/tft/v1/tftv1connect"
	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
)

var _ tftv1connect.PatchServiceHandler = (*Service)(nil)

type Service struct {
	db      *pgxpool.Pool
	queries *generated.Queries
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		db:      db,
		queries: generated.New(db),
	}
}

func (s *Service) GetPatchData(
	ctx context.Context,
	req *connect.Request[tftv1.GetPatchDataRequest],
) (*connect.Response[tftv1.GetPatchDataResponse], error) {
	setNumber := req.Msg.SetNumber

	// If set_number is 0, get latest
	var set generated.Set
	var err error
	if setNumber == 0 {
		set, err = s.queries.GetLatestSet(ctx)
	} else {
		set, err = s.queries.GetSetByNumber(ctx, setNumber)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	// Fetch all data in parallel could be done, but for simplicity sequential
	champions, err := s.queries.GetChampionsBySet(ctx, set.Number)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	traits, err := s.queries.GetTraitsBySet(ctx, set.Number)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	items, err := s.queries.GetItemsBySet(ctx, set.Number)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Build champion trait map
	championTraits := make(map[string][]string)
	for _, c := range champions {
		ct, err := s.queries.GetTraitsByChampion(ctx, generated.GetTraitsByChampionParams{
			ChampionApiName: c.ApiName,
			SetNumber:       set.Number,
		})
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		traitNames := make([]string, 0, len(ct))
		for _, t := range ct {
			traitNames = append(traitNames, t.TraitApiName)
		}
		championTraits[c.ApiName] = traitNames
	}

	// Map to protobuf
	pbChampions := make([]*tftv1.Champion, 0, len(champions))
	for _, c := range champions {
		pbChampions = append(pbChampions, &tftv1.Champion{
			ApiName:       c.ApiName,
			Name:          c.Name,
			Cost:          c.Cost,
			Role:          c.Role,
			TraitApiNames: championTraits[c.ApiName],
			Stats:         mapStatsToProto(c.Stats),
			Ability:       mapAbilityToProto(c.Ability),
			IconUrl:       c.IconUrl,
			SquareIconUrl: c.SquareIconUrl,
			TileIconUrl:   c.TileIconUrl,
		})
	}

	pbTraits := make([]*tftv1.Trait, 0, len(traits))
	for _, t := range traits {
		pbTraits = append(pbTraits, &tftv1.Trait{
			ApiName: t.ApiName,
			Name:    t.Name,
			Desc:    t.Description,
			IconUrl: t.IconUrl,
			Effects: mapTraitEffectsToProto(t.Effects),
		})
	}

	pbItems := make([]*tftv1.Item, 0, len(items))
	for _, i := range items {
		pbItems = append(pbItems, &tftv1.Item{
			ApiName:            i.ApiName,
			Name:               i.Name,
			Desc:               i.Description,
			Composition:        i.Composition,
			Effects:            mapItemEffectsToProto(i.Effects),
			IconUrl:            i.IconUrl,
			AssociatedTraits:   i.AssociatedTraits,
			IncompatibleTraits: i.IncompatibleTraits,
			Tags:               i.Tags,
			Unique:             i.IsUnique,
		})
	}

	return connect.NewResponse(&tftv1.GetPatchDataResponse{
		Set: &tftv1.SetMetadata{
			Number:  set.Number,
			Name:    set.Name,
			Mutator: set.Mutator,
		},
		Champions: pbChampions,
		Items:     pbItems,
		Traits:    pbTraits,
	}), nil
}
