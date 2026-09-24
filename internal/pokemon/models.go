package pokemon

import (
	"encoding/json"
	"time"
)

// BaseStats represents normalized Pokémon stats
type BaseStats struct {
	HP             int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
}

// Pokemon represents a canonical Pokémon entity stored in PostgreSQL & Redis
type Pokemon struct {
	ID               int             `json:"id"`
	Name             string          `json:"name"`
	SpeciesID        int             `json:"species_id"`
	EvolutionChainID int             `json:"evolution_chain_id"`
	Generation       int             `json:"generation"`
	Types            []string        `json:"types"`
	BaseStats        BaseStats       `json:"base_stats"`
	Abilities        json.RawMessage `json:"abilities"`
	SpriteURL        string          `json:"sprite_url"`
	RawJSON          json.RawMessage `json:"raw_json"`
	FetchedAt        time.Time       `json:"fetched_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// Evolution trigger names as used by PokéAPI evolution details.
const (
	TriggerLevelUp = "level-up"
	TriggerUseItem = "use-item"
)

// Evolution methods exposed by EvolutionOption.
const (
	MethodLevel = "level"
	MethodItem  = "item"
)

// FriendshipEvolutionLevel is the level at which friendship-based evolutions
// (PokéAPI min_happiness, e.g. pichu -> pikachu) become level-up evolutions
// in this game. PokeTacTix has no friendship mechanic, so happiness is
// translated into a fixed level requirement during chain extraction; the
// constant lives here so the rule is adjustable in one place.
const FriendshipEvolutionLevel = 20

// CurrentEvolutionExtractionVersion tags freshly extracted evolution chains.
// Stored/cached chains with a lower version were extracted by older code and
// are treated as stale (e.g. they predate min_happiness synthesis) so they
// self-heal via a refresh from PokéAPI. Bump this whenever extraction output
// changes shape in a way older chains must pick up.
const CurrentEvolutionExtractionVersion = 3

// EvolutionLink represents a single directed edge in an evolution chain
// (e.g. charmander -> charmeleon at level 16 via level-up, or pikachu ->
// raichu via use-item with a thunder-stone).
type EvolutionLink struct {
	FromSpeciesID         int    `json:"from"`
	ToSpeciesID           int    `json:"to"`
	MinLevel              int    `json:"min_level,omitempty"`
	Trigger               string `json:"trigger,omitempty"` // e.g. "level-up", "use-item", "trade"
	Item                  string `json:"item,omitempty"`    // required item id (slug) for "use-item" triggers, e.g. "thunder-stone"
	TimeOfDay             string `json:"time_of_day,omitempty"`
	UnsupportedConditions bool   `json:"unsupported_conditions,omitempty"`
}

// PickEvolutionLink returns the level-up evolution edge from speciesID whose
// minimum level has been reached at the given level. Non-level-up triggers
// (items, trades, etc.) are ignored — this game only evolves through leveling.
// Branching chains resolve deterministically to the lowest qualifying level.
// Returns nil when no evolution applies.
func PickEvolutionLink(links []EvolutionLink, speciesID int, level int) *EvolutionLink {
	var best *EvolutionLink
	for i := range links {
		l := &links[i]
		if l.FromSpeciesID != speciesID {
			continue
		}
		if l.Trigger != "level-up" || l.MinLevel <= 0 || l.TimeOfDay != "" || l.UnsupportedConditions {
			continue
		}
		if level < l.MinLevel {
			continue
		}
		if best == nil || l.MinLevel < best.MinLevel {
			best = l
		}
	}
	return best
}

// PickEvolutionLinkForItem returns the use-item evolution edge from speciesID
// that requires the given item id (slug). Level-up edges are ignored — item
// evolution must never trigger from level alone. Branching chains resolve
// deterministically to the first matching edge. Returns nil when no edge
// applies.
func PickEvolutionLinkForItem(links []EvolutionLink, speciesID int, itemID string) *EvolutionLink {
	for i := range links {
		l := &links[i]
		if l.FromSpeciesID != speciesID {
			continue
		}
		if l.Trigger != TriggerUseItem || l.Item == "" || l.TimeOfDay != "" || l.UnsupportedConditions {
			continue
		}
		if l.Item != itemID {
			continue
		}
		return l
	}
	return nil
}

// EvolutionOption describes one way a Pokemon can evolve, normalized to a
// method tag so callers (deck UI, evolve endpoints) can treat level and item
// evolution uniformly.
type EvolutionOption struct {
	Method       string `json:"method"` // MethodLevel or MethodItem
	ToSpeciesID  int    `json:"to_species_id"`
	TargetName   string `json:"target_name"`
	TargetSprite string `json:"target_sprite"`
	MinLevel     int    `json:"min_level,omitempty"` // MethodLevel: required level
	Item         string `json:"item,omitempty"`      // MethodItem: required item id (slug)
}

// EvolutionOptionsFromLinks returns the evolution edges leaving speciesID as
// method-tagged options. Edges with an unusable trigger (trade, friendship,
// level-up without a minimum level, use-item without a known item) are left
// out — they cannot be satisfied by any mechanism this game implements.
func EvolutionOptionsFromLinks(links []EvolutionLink, speciesID int) []EvolutionOption {
	var options []EvolutionOption
	for _, l := range links {
		if l.FromSpeciesID != speciesID {
			continue
		}
		switch {
		case l.Trigger == TriggerLevelUp && l.MinLevel > 0 && l.TimeOfDay == "" && !l.UnsupportedConditions:
			options = append(options, EvolutionOption{
				Method:      MethodLevel,
				ToSpeciesID: l.ToSpeciesID,
				MinLevel:    l.MinLevel,
			})
		case l.Trigger == TriggerUseItem && l.Item != "" && l.TimeOfDay == "" && !l.UnsupportedConditions:
			options = append(options, EvolutionOption{
				Method:      MethodItem,
				ToSpeciesID: l.ToSpeciesID,
				Item:        l.Item,
			})
		}
	}
	return options
}

// EvolutionChain represents an evolution chain entity stored in PostgreSQL & Redis
type EvolutionChain struct {
	ID               int             `json:"id"`
	MemberSpeciesIDs []int           `json:"member_species_ids"`
	Links            []EvolutionLink `json:"links"`
	FetchedAt        time.Time       `json:"fetched_at"`
	// Version records which extraction code produced the links; chains below
	// CurrentEvolutionExtractionVersion are refreshed on first use.
	Version int `json:"version,omitempty"`
}

func cardHPFromBase(baseHP int) int {
	if baseHP <= 0 {
		baseHP = 50
	}
	return baseHP + baseHP/2
}

// CardBaseStats returns the four base stat values that get stored on player_cards,
// applying the same floor defaults and HP boost as ToCard(). Both ToCard() and
// evolvedBaseStats() in the battle package delegate here so the logic never
// diverges between the two code paths.
func (p *Pokemon) CardBaseStats() (hp, attack, defense, speed int) {
	hp = cardHPFromBase(p.BaseStats.HP)
	attack = p.BaseStats.Attack
	if attack <= 0 {
		attack = 40
	}
	defense = p.BaseStats.Defense
	if defense <= 0 {
		defense = 40
	}
	speed = p.BaseStats.Speed
	if speed <= 0 {
		speed = 50
	}
	return hp, attack, defense, speed
}

// ToCard converts a Pokemon domain entity into a battle Card
func (p *Pokemon) ToCard() Card {
	hp, attack, defense, speed := p.CardBaseStats()

	stamina := speed * 2
	isLegendary, isMythical := IsLegendaryOrMythical(p.Name)

	return Card{
		CardID:      p.ID,
		Name:        p.Name,
		HP:          hp,
		HPMax:       hp,
		Stamina:     stamina,
		Attack:      attack,
		Defense:     defense,
		Speed:       speed,
		Types:       p.Types,
		Sprite:      p.SpriteURL,
		Level:       1,
		XP:          0,
		IsLegendary: isLegendary,
		IsMythical:  isMythical,
	}
}
