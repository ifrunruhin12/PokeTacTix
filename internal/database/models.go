package database

import (
	"encoding/json"
	"time"
)

// User represents a user account
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Coins        int       `json:"coins"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PlayerCard represents a Pokemon card owned by a player
type PlayerCard struct {
	ID           int             `json:"id"`
	UserID       int             `json:"user_id"`
	PokemonName  string          `json:"pokemon_name"`
	Level        int             `json:"level"`
	XP           int             `json:"xp"`
	BaseHP       int             `json:"base_hp"`
	BaseAttack   int             `json:"base_attack"`
	BaseDefense  int             `json:"base_defense"`
	BaseSpeed    int             `json:"base_speed"`
	Types        json.RawMessage `json:"types"`
	Moves        json.RawMessage `json:"moves"`
	Sprite       string          `json:"sprite"`
	IsLegendary  bool            `json:"is_legendary"`
	IsMythical   bool            `json:"is_mythical"`
	InDeck       bool            `json:"in_deck"`
	DeckPosition *int            `json:"deck_position,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// CardStats represents computed stats for a card at its current level
type CardStats struct {
	HP      int `json:"hp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"`
	Stamina int `json:"stamina"`
}

// GetCurrentStats calculates current stats based on level
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

// BattleHistory represents a battle record
type BattleHistory struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Mode        string    `json:"mode"`
	Result      string    `json:"result"`
	CoinsEarned int       `json:"coins_earned"`
	Duration    *int      `json:"duration,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// PlayerStats represents player statistics
type PlayerStats struct {
	UserID           int       `json:"user_id"`
	TotalBattles1v1  int       `json:"total_battles_1v1"`
	Wins1v1          int       `json:"wins_1v1"`
	Losses1v1        int       `json:"losses_1v1"`
	TotalBattles5v5  int       `json:"total_battles_5v5"`
	Wins5v5          int       `json:"wins_5v5"`
	Losses5v5        int       `json:"losses_5v5"`
	TotalCoinsEarned int       `json:"total_coins_earned"`
	HighestLevel     int       `json:"highest_level"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Achievement represents an achievement definition
type Achievement struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Icon             string `json:"icon"`
	RequirementType  string `json:"requirement_type"`
	RequirementValue int    `json:"requirement_value"`
}

// UserAchievement represents an unlocked achievement
type UserAchievement struct {
	UserID        int       `json:"user_id"`
	AchievementID int       `json:"achievement_id"`
	UnlockedAt    time.Time `json:"unlocked_at"`
}

// AchievementWithStatus includes unlock status
type AchievementWithStatus struct {
	Achievement
	Unlocked   bool       `json:"unlocked"`
	UnlockedAt *time.Time `json:"unlocked_at,omitempty"`
}
