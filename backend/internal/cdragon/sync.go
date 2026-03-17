package cdragon

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
)

type Syncer struct {
	client  *Client
	db      *pgxpool.Pool
	queries *generated.Queries
}

func NewSyncer(db *pgxpool.Pool) *Syncer {
	return &Syncer{
		client:  NewClient(),
		db:      db,
		queries: generated.New(db),
	}
}

// Sync fetches CommunityDragon data and stores it in the database.
func (s *Syncer) Sync(ctx context.Context, locale string) error {
	log.Println("cdragon: fetching data...")
	data, err := s.client.Fetch(ctx, locale)
	if err != nil {
		return fmt.Errorf("fetch cdragon: %w", err)
	}

	log.Println("cdragon: parsing data...")
	parsed := Parse(data)
	if parsed == nil {
		return fmt.Errorf("no current set found in cdragon data")
	}

	log.Printf("cdragon: syncing set %d (%s) — %d champions, %d traits, %d items, %d augments",
		parsed.Number, parsed.Name, len(parsed.Champions), len(parsed.Traits), len(parsed.Items), len(parsed.Augments))

	if err := s.store(ctx, parsed); err != nil {
		return fmt.Errorf("store cdragon data: %w", err)
	}

	log.Println("cdragon: sync complete")
	return nil
}

// SyncIfEmpty runs sync only if the database has no sets.
func (s *Syncer) SyncIfEmpty(ctx context.Context, locale string) error {
	_, err := s.queries.GetLatestSet(ctx)
	if err == nil {
		log.Println("cdragon: data already exists, skipping sync")
		return nil
	}
	return s.Sync(ctx, locale)
}

func (s *Syncer) store(ctx context.Context, parsed *ParsedSet) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// 1. Upsert set
	_, err = qtx.UpsertSet(ctx, generated.UpsertSetParams{
		Number:  int32(parsed.Number),
		Name:    parsed.Name,
		Mutator: parsed.Mutator,
	})
	if err != nil {
		return fmt.Errorf("upsert set: %w", err)
	}

	setNumber := int32(parsed.Number)

	// 2. Upsert traits (must be before champions due to FK)
	for _, t := range parsed.Traits {
		effectsJSON, err := json.Marshal(t.Effects)
		if err != nil {
			return fmt.Errorf("marshal trait effects: %w", err)
		}

		_, err = qtx.UpsertTrait(ctx, generated.UpsertTraitParams{
			ApiName:     t.APIName,
			SetNumber:   setNumber,
			Name:        t.Name,
			Description: t.Desc,
			IconUrl:     t.IconURL,
			Effects:     effectsJSON,
		})
		if err != nil {
			return fmt.Errorf("upsert trait %s: %w", t.APIName, err)
		}
	}

	// 3. Delete old champions for this set, then re-insert.
	// This ensures stale entries (non-playable units from older syncs) are removed.
	if err := qtx.DeleteChampionsBySet(ctx, setNumber); err != nil {
		return fmt.Errorf("delete champions: %w", err)
	}

	for _, c := range parsed.Champions {
		statsJSON, err := json.Marshal(map[string]interface{}{
			"hp":              c.Stats.HP,
			"armor":           c.Stats.Armor,
			"magic_resist":    c.Stats.MagicResist,
			"damage":          c.Stats.Damage,
			"attack_speed":    c.Stats.AttackSpeed,
			"range":           c.Stats.Range,
			"mana":            int32(c.Stats.Mana),
			"initial_mana":    c.Stats.InitialMana,
			"crit_chance":     c.Stats.CritChance,
			"crit_multiplier": c.Stats.CritMultiplier,
		})
		if err != nil {
			return fmt.Errorf("marshal champion stats: %w", err)
		}

		abilityVars := make([]map[string]interface{}, 0, len(c.Ability.Variables))
		for _, v := range c.Ability.Variables {
			abilityVars = append(abilityVars, map[string]interface{}{
				"name":   v.Name,
				"values": v.Values,
			})
		}
		abilityJSON, err := json.Marshal(map[string]interface{}{
			"name":      c.Ability.Name,
			"desc":      c.Ability.Desc,
			"icon_url":  c.Ability.IconURL,
			"variables": abilityVars,
		})
		if err != nil {
			return fmt.Errorf("marshal champion ability: %w", err)
		}

		_, err = qtx.UpsertChampion(ctx, generated.UpsertChampionParams{
			ApiName:       c.APIName,
			SetNumber:     setNumber,
			Name:          c.Name,
			Cost:          int32(c.Cost),
			Role:          "", // CDragon doesn't provide role directly
			Stats:         statsJSON,
			Ability:       abilityJSON,
			IconUrl:       c.IconURL,
			SquareIconUrl: c.SquareIconURL,
			TileIconUrl:   c.TileIconURL,
		})
		if err != nil {
			return fmt.Errorf("upsert champion %s: %w", c.APIName, err)
		}
	}

	// 4. Delete old champion_traits for this set and re-insert
	err = qtx.DeleteChampionTraitsBySet(ctx, setNumber)
	if err != nil {
		return fmt.Errorf("delete champion traits: %w", err)
	}

	for _, c := range parsed.Champions {
		for _, traitAPI := range c.TraitAPINames {
			err = qtx.InsertChampionTrait(ctx, generated.InsertChampionTraitParams{
				ChampionApiName: c.APIName,
				TraitApiName:    traitAPI,
				SetNumber:       setNumber,
			})
			if err != nil {
				return fmt.Errorf("insert champion trait %s-%s: %w", c.APIName, traitAPI, err)
			}
		}
	}

	// 5. Delete old items for this set, then re-insert (same as champions).
	if err := qtx.DeleteItemsBySet(ctx, setNumber); err != nil {
		return fmt.Errorf("delete items: %w", err)
	}

	for _, item := range parsed.Items {
		effectsJSON, err := json.Marshal(item.Effects)
		if err != nil {
			return fmt.Errorf("marshal item effects: %w", err)
		}

		composition := item.Composition
		if composition == nil {
			composition = []string{}
		}
		assocTraits := item.AssociatedTraits
		if assocTraits == nil {
			assocTraits = []string{}
		}
		incompTraits := item.IncompatibleTraits
		if incompTraits == nil {
			incompTraits = []string{}
		}
		tags := item.Tags
		if tags == nil {
			tags = []string{}
		}

		_, err = qtx.UpsertItem(ctx, generated.UpsertItemParams{
			ApiName:            item.APIName,
			SetNumber:          setNumber,
			Name:               item.Name,
			Description:        item.Desc,
			Composition:        composition,
			Effects:            effectsJSON,
			IconUrl:            item.IconURL,
			AssociatedTraits:   assocTraits,
			IncompatibleTraits: incompTraits,
			Tags:               tags,
			IsUnique:           item.Unique,
		})
		if err != nil {
			return fmt.Errorf("upsert item %s: %w", item.APIName, err)
		}
	}

	// 6. Insert augments (stored in the items table, differentiated by "augment" tag)
	for _, aug := range parsed.Augments {
		effectsJSON, err := json.Marshal(aug.Effects)
		if err != nil {
			return fmt.Errorf("marshal augment effects: %w", err)
		}

		composition := aug.Composition
		if composition == nil {
			composition = []string{}
		}
		assocTraits := aug.AssociatedTraits
		if assocTraits == nil {
			assocTraits = []string{}
		}
		incompTraits := aug.IncompatibleTraits
		if incompTraits == nil {
			incompTraits = []string{}
		}
		tags := aug.Tags
		if tags == nil {
			tags = []string{}
		}

		_, err = qtx.UpsertItem(ctx, generated.UpsertItemParams{
			ApiName:            aug.APIName,
			SetNumber:          setNumber,
			Name:               aug.Name,
			Description:        aug.Desc,
			Composition:        composition,
			Effects:            effectsJSON,
			IconUrl:            aug.IconURL,
			AssociatedTraits:   assocTraits,
			IncompatibleTraits: incompTraits,
			Tags:               tags,
			IsUnique:           aug.Unique,
		})
		if err != nil {
			return fmt.Errorf("upsert augment %s: %w", aug.APIName, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
