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

// EvolutionChain represents an evolution chain entity stored in PostgreSQL & Redis
type EvolutionChain struct {
	ID               int       `json:"id"`
	MemberSpeciesIDs []int     `json:"member_species_ids"`
	FetchedAt        time.Time `json:"fetched_at"`
}

// ToCard converts a Pokemon domain entity into a battle Card
func (p *Pokemon) ToCard() Card {
	hp := p.BaseStats.HP
	if hp <= 0 {
		hp = 50
	}
	attack := p.BaseStats.Attack
	if attack <= 0 {
		attack = 40
	}
	defense := p.BaseStats.Defense
	if defense <= 0 {
		defense = 40
	}
	speed := p.BaseStats.Speed
	if speed <= 0 {
		speed = 50
	}

	stamina := speed * 2

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
		IsLegendary: (hp + attack + defense + speed) >= 450,
		IsMythical:  false,
	}
}
