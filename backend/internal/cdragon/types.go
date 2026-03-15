package cdragon

// CDragonData is the top-level response from CommunityDragon TFT JSON.
type CDragonData struct {
	SetData []CDragonSetData `json:"setData"`
	Items   []CDragonItem    `json:"items"`
}

// CDragonSetData represents a TFT set with its champions and traits.
type CDragonSetData struct {
	Number    int               `json:"number"`
	Name      string            `json:"name"`
	Mutator   string            `json:"mutator"`
	Champions []CDragonChampion `json:"champions"`
	Traits    []CDragonTrait    `json:"traits"`
}

// CDragonChampion represents a champion from CommunityDragon.
type CDragonChampion struct {
	APIName       string              `json:"apiName"`
	Name          string              `json:"name"`
	Cost          int                 `json:"cost"`
	Traits        []string            `json:"traits"` // Display names, NOT apiNames
	Stats         CDragonChampionStat `json:"stats"`
	Ability       CDragonAbility      `json:"ability"`
	Icon          string              `json:"icon"`
	SquareIcon    string              `json:"squareIcon"`
	TileIcon      string              `json:"tileIcon"`
	CharacterName string              `json:"characterName"`
}

// CDragonChampionStat holds champion numeric stats.
type CDragonChampionStat struct {
	HP             float64 `json:"hp"`
	Armor          float64 `json:"armor"`
	MagicResist    float64 `json:"magicResist"`
	Damage         float64 `json:"damage"`
	AttackSpeed    float64 `json:"attackSpeed"`
	Range          float64 `json:"range"`
	Mana           float64 `json:"mana"`
	InitialMana    float64 `json:"initialMana"`
	CritChance     float64 `json:"critChance"`
	CritMultiplier float64 `json:"critMultiplier"`
}

// CDragonAbility represents a champion's ability.
type CDragonAbility struct {
	Name      string                  `json:"name"`
	Desc      string                  `json:"desc"`
	Icon      string                  `json:"icon"`
	Variables []CDragonAbilityVariable `json:"variables"`
}

// CDragonAbilityVariable holds an ability scaling variable.
type CDragonAbilityVariable struct {
	Name   string     `json:"name"`
	Values []*float64 `json:"value"` // Note: CDragon uses "value" not "values", and values can be null
}

// CDragonTrait represents a trait from CommunityDragon.
type CDragonTrait struct {
	APIName string              `json:"apiName"`
	Name    string              `json:"name"`
	Desc    string              `json:"desc"`
	Icon    string              `json:"icon"`
	Effects []CDragonTraitEffect `json:"effects"`
}

// CDragonTraitEffect represents a breakpoint in a trait.
type CDragonTraitEffect struct {
	MinUnits  int                `json:"minUnits"`
	MaxUnits  int                `json:"maxUnits"`
	Style     int                `json:"style"`
	Variables map[string]float64 `json:"variables"`
}

// CDragonItem represents an item from CommunityDragon.
type CDragonItem struct {
	APIName            string             `json:"apiName"`
	Name               string             `json:"name"`
	Desc               string             `json:"desc"`
	Composition        []string           `json:"composition"`
	Effects            map[string]float64 `json:"effects"`
	Icon               string             `json:"icon"`
	AssociatedTraits   []string           `json:"associatedTraits"`
	IncompatibleTraits []string           `json:"incompatibleTraits"`
	Tags               []string           `json:"tags"`
	Unique             bool               `json:"unique"`
}
