// Package models contains the game state, card and player models for the Pokémon CLI application.
package models

type GameState struct {
	BattleStarted     bool
	Player            *Player
	AI                *Player
	PlayerName        string
	InBattle          bool
	HaveCard          bool
	Round             int
	PlayerActiveIdx   int
	AIActiveIdx       int
	CardMovePlayer    int
	CardMoveAI        int
	CurrentMovetype   string
	RoundStarted      bool
	SwitchedThisRound bool
	BattleOver        bool
	RoundOver         bool
	SacrificeCount    map[int]int // key: PlayerActiveIdx, value: number of sacrifices for that Pokémon
	LastHpLost        int
	LastStaminaLost   int
	LastDamageDealt   int
	PlayerSurrendered bool // Track if player surrendered the whole battle
	JustSwitched      bool // true if player just switched and hasn't played a round yet
	HasPlayedRound    bool
	TurnNumber        int
	BattleMode        string // "1v1" or "5v5"
	// --- Web turn-based fields ---
	PendingPlayerMove    string // Player's chosen move for this turn (web)
	PendingPlayerMoveIdx int    // Player's chosen move index (web)
	PendingAIMove        string // AI's chosen move for this turn (web)
	PendingAIMoveIdx     int    // AI's chosen move index (web)
	WhoseTurn            string // "player" or "ai" (web)
}
