package patch

import (
	"encoding/json"
	"fmt"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
)

type dbStats struct {
	HP             float64 `json:"hp"`
	Armor          float64 `json:"armor"`
	MagicResist    float64 `json:"magic_resist"`
	Damage         float64 `json:"damage"`
	AttackSpeed    float64 `json:"attack_speed"`
	Range          float64 `json:"range"`
	Mana           int32   `json:"mana"`
	InitialMana    float64 `json:"initial_mana"`
	CritChance     float64 `json:"crit_chance"`
	CritMultiplier float64 `json:"crit_multiplier"`
}

type dbAbility struct {
	Name      string              `json:"name"`
	Desc      string              `json:"desc"`
	IconURL   string              `json:"icon_url"`
	Variables []dbAbilityVariable `json:"variables"`
}

type dbAbilityVariable struct {
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
}

type dbTraitEffect struct {
	MinUnits  int32              `json:"min_units"`
	MaxUnits  int32              `json:"max_units"`
	Style     int32              `json:"style"`
	Variables map[string]float64 `json:"variables"`
}

func mapStatsToProto(raw []byte) *tftv1.ChampionStats {
	var s dbStats
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil
	}
	return &tftv1.ChampionStats{
		Hp:             s.HP,
		Armor:          s.Armor,
		MagicResist:    s.MagicResist,
		Damage:         s.Damage,
		AttackSpeed:    s.AttackSpeed,
		Range:          s.Range,
		Mana:           s.Mana,
		InitialMana:    s.InitialMana,
		CritChance:     s.CritChance,
		CritMultiplier: s.CritMultiplier,
	}
}

func mapAbilityToProto(raw []byte) *tftv1.ChampionAbility {
	var a dbAbility
	if err := json.Unmarshal(raw, &a); err != nil {
		return nil
	}

	vars := make([]*tftv1.AbilityVariable, 0, len(a.Variables))
	for _, v := range a.Variables {
		vars = append(vars, &tftv1.AbilityVariable{
			Name:   v.Name,
			Values: v.Values,
		})
	}

	return &tftv1.ChampionAbility{
		Name:      a.Name,
		Desc:      a.Desc,
		IconUrl:   a.IconURL,
		Variables: vars,
	}
}

func mapTraitEffectsToProto(raw []byte) []*tftv1.TraitEffect {
	var effects []dbTraitEffect
	if err := json.Unmarshal(raw, &effects); err != nil {
		return nil
	}

	result := make([]*tftv1.TraitEffect, 0, len(effects))
	for _, e := range effects {
		result = append(result, &tftv1.TraitEffect{
			MinUnits:  e.MinUnits,
			MaxUnits:  e.MaxUnits,
			Style:     e.Style,
			Variables: e.Variables,
		})
	}
	return result
}

func mapItemEffectsToProto(raw []byte) map[string]float64 {
	var effects map[string]float64
	if err := json.Unmarshal(raw, &effects); err != nil {
		return nil
	}
	return effects
}

func statsToJSON(stats *tftv1.ChampionStats) ([]byte, error) {
	if stats == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(dbStats{
		HP:             stats.Hp,
		Armor:          stats.Armor,
		MagicResist:    stats.MagicResist,
		Damage:         stats.Damage,
		AttackSpeed:    stats.AttackSpeed,
		Range:          stats.Range,
		Mana:           stats.Mana,
		InitialMana:    stats.InitialMana,
		CritChance:     stats.CritChance,
		CritMultiplier: stats.CritMultiplier,
	})
}

func abilityToJSON(ability *tftv1.ChampionAbility) ([]byte, error) {
	if ability == nil {
		return []byte("{}"), nil
	}

	vars := make([]dbAbilityVariable, 0, len(ability.Variables))
	for _, v := range ability.Variables {
		vars = append(vars, dbAbilityVariable{
			Name:   v.Name,
			Values: v.Values,
		})
	}

	return json.Marshal(dbAbility{
		Name:      ability.Name,
		Desc:      ability.Desc,
		IconURL:   ability.IconUrl,
		Variables: vars,
	})
}

func traitEffectsToJSON(effects []*tftv1.TraitEffect) ([]byte, error) {
	if effects == nil {
		return []byte("[]"), nil
	}

	dbEffects := make([]dbTraitEffect, 0, len(effects))
	for _, e := range effects {
		dbEffects = append(dbEffects, dbTraitEffect{
			MinUnits:  e.MinUnits,
			MaxUnits:  e.MaxUnits,
			Style:     e.Style,
			Variables: e.Variables,
		})
	}
	return json.Marshal(dbEffects)
}

func itemEffectsToJSON(effects map[string]float64) ([]byte, error) {
	if effects == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(effects)
}

func mustJSON(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	}
	return string(b)
}
