package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
)

// Store handles persisting raw tier data to the database.
type Store struct {
	db      *pgxpool.Pool
	queries *generated.Queries
}

// NewStore creates a new crawler store.
func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		queries: generated.New(db),
	}
}

// SaveResult persists a scrape result to the database.
// It deletes old data for the source+patch and inserts new rows in a single transaction.
func (s *Store) SaveResult(ctx context.Context, result *ScrapeResult) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// Delete old data for this source+patch
	if err := qtx.DeleteRawTierDataBySourceAndPatch(ctx, generated.DeleteRawTierDataBySourceAndPatchParams{
		Source: result.Source,
		Patch:  result.Patch,
	}); err != nil {
		return fmt.Errorf("delete old data: %w", err)
	}

	// Insert new compositions
	for _, comp := range result.Compositions {
		coreItemsJSON, err := json.Marshal(comp.CoreItems)
		if err != nil {
			return fmt.Errorf("marshal core items: %w", err)
		}

		championIDs := comp.ChampionIDs
		if championIDs == nil {
			championIDs = []string{}
		}

		_, err = qtx.InsertRawTierData(ctx, generated.InsertRawTierDataParams{
			Source:          result.Source,
			Patch:           result.Patch,
			CompositionName: comp.Name,
			Tier:            textFromString(comp.Tier),
			WinRate:         numericFromFloat(comp.WinRate),
			PlayRate:        numericFromFloat(comp.PlayRate),
			AvgPlacement:    numericFromFloat(comp.AvgPlacement),
			ChampionIds:     championIDs,
			CoreItems:       coreItemsJSON,
		})
		if err != nil {
			return fmt.Errorf("insert comp %q: %w", comp.Name, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	log.Printf("crawler: stored %d compositions from %s (patch %s)",
		len(result.Compositions), result.Source, result.Patch)
	return nil
}

func textFromString(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func numericFromFloat(f float64) pgtype.Numeric {
	if f == 0 {
		return pgtype.Numeric{Valid: false}
	}
	// Convert float to *big.Int with 2 decimal places
	scaled := int64(f * 100)
	return pgtype.Numeric{
		Int:   big.NewInt(scaled),
		Exp:   -2,
		Valid: true,
	}
}
