package storage

import (
	"time"

	"pokemon-cli/internal/pokemon"
)

// GameState represents the complete state of the CLI game
type GameState struct {
	PlayerName    string         `json:"player_name"`
	Coins         int            `json:"coins"`
	Collection    []PlayerCard   `json:"collection"`
	Deck          []int          `json:"deck"` // Card IDs (indices in Collection)
	Stats         PlayerStats    `json:"stats"`
	ShopState     ShopState      `json:"shop_state"`
	BattleHistory []BattleRecord `json:"battle_history,omitempty"`
	Settings      GameSettings   `json:"settings"`
	LastSaved     time.Time      `json:"last_saved"`
	Version       string         `json:"version"`
}

// GameSettings stores user preferences
type GameSettings struct {
	QuickBattle bool   `json:"quick_battle"` // Skip animations and delays
	BattleSpeed string `json:"battle_speed"` // "slow", "normal", "fast"
}

// PlayerCard represents a Pokemon card owned by the player in CLI mode
type PlayerCard struct {
	ID           int            `json:"id"`
	PokemonID    int            `json:"pokemon_id"` // ID from pokemon database
	Name         string         `json:"name"`
	Level        int            `json:"level"`
	XP           int            `json:"xp"`
	BaseHP       int            `json:"base_hp"`
	BaseAttack   int            `json:"base_attack"`
	BaseDefense  int            `json:"base_defense"`
	BaseSpeed    int            `json:"base_speed"`
	Types        []string       `json:"types"`
	Moves        []pokemon.Move `json:"moves"`
	Sprite       string         `json:"sprite"`
	IsLegendary  bool           `json:"is_legendary"`
	IsMythical   bool           `json:"is_mythical"`
	AcquiredAt   time.Time      `json:"acquired_at"`
}

// PlayerStats tracks battle statistics for the player
type PlayerStats struct {
	TotalBattles1v1  int `json:"total_battles_1v1"`
	Wins1v1          int `json:"wins_1v1"`
	Losses1v1        int `json:"losses_1v1"`
	Draws1v1         int `json:"draws_1v1"`
	TotalBattles5v5  int `json:"total_battles_5v5"`
	Wins5v5          int `json:"wins_5v5"`
	Losses5v5        int `json:"losses_5v5"`
	Draws5v5         int `json:"draws_5v5"`
	HighestLevel     int `json:"highest_level"`
	TotalCoinsEarned int `json:"total_coins_earned"`
	TotalPokemon     int `json:"total_pokemon_collected"`
}

// ShopState tracks the shop inventory and refresh status
type ShopState struct {
	Inventory           []ShopItem `json:"inventory"`
	LastRefresh         time.Time  `json:"last_refresh"`
	BattlesSinceRefresh int        `json:"battles_since_refresh"`
}

// ShopItem represents a Pokemon available for purchase in the shop
type ShopItem struct {
	PokemonID   int            `json:"pokemon_id"`
	Name        string         `json:"name"`
	Types       []string       `json:"types"`
	BaseHP      int            `json:"base_hp"`
	BaseAttack  int            `json:"base_attack"`
	BaseDefense int            `json:"base_defense"`
	BaseSpeed   int            `json:"base_speed"`
	Moves       []pokemon.Move `json:"moves"`
	Sprite      string         `json:"sprite"`
	Price       int            `json:"price"`
	Rarity      string         `json:"rarity"`
	IsLegendary bool           `json:"is_legendary"`
	IsMythical  bool           `json:"is_mythical"`
}

// BattleRecord represents a single battle in the history
type BattleRecord struct {
	Mode        string    `json:"mode"`        // "1v1" or "5v5"
	Result      string    `json:"result"`      // "win", "loss", "draw"
	CoinsEarned int       `json:"coins_earned"`
	Duration    int       `json:"duration,omitempty"` // in seconds
	Timestamp   time.Time `json:"timestamp"`
}

// CardStats represents computed stats based on level
type CardStats struct {
	HP      int `json:"hp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"`
	Stamina int `json:"stamina"`
}

// GetCurrentStats calculates current stats based on level for PlayerCard
func (c *PlayerCard) GetCurrentStats() CardStats {
	levelMultiplier := float64(c.Level - 1)

	hp := int(float64(c.BaseHP) * (1.0 + levelMultiplier*0.03))
	attack := int(float64(c.BaseAttack) * (1.0 + levelMultiplier*0.02))
	defense := int(float64(c.BaseDefense) * (1.0 + levelMultiplier*0.02))
	speed := int(float64(c.BaseSpeed) * (1.0 + levelMultiplier*0.01))
	stamina := speed * 2

	return CardStats{
		HP:      hp,
		Attack:  attack,
		Defense: defense,
		Speed:   speed,
		Stamina: stamina,
	}
}

// ToCard converts a PlayerCard to a pokemon.Card for battle use
func (c *PlayerCard) ToCard() pokemon.Card {
	stats := c.GetCurrentStats()
	return pokemon.Card{
		Name:        c.Name,
		HP:          stats.HP,
		HPMax:       stats.HP,
		Stamina:     stats.Stamina,
		Defense:     stats.Defense,
		Attack:      stats.Attack,
		Speed:       stats.Speed,
		Moves:       c.Moves,
		Types:       c.Types,
		Sprite:      c.Sprite,
		Level:       c.Level,
		XP:          c.XP,
		IsLegendary: c.IsLegendary,
		IsMythical:  c.IsMythical,
	}
}
