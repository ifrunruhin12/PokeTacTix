// Package battle provides battle state management for the web API.
package battle

import (
	"pokemon-cli/game/models"
	"pokemon-cli/internal/pokemon"
	"time"
)

// BattleState represents the complete state of a battle session
// This is the enhanced model that supports both 1v1 and 5v5 modes
type BattleState struct {
	ID                   string       `json:"id"`
	UserID               int          `json:"user_id"`
	Mode                 string       `json:"mode"` // "1v1" or "5v5"
	PlayerDeck           []BattleCard `json:"player_deck"`
	AIDeck               []BattleCard `json:"ai_deck"`
	PlayerActiveIdx      int          `json:"player_active_idx"`
	AIActiveIdx          int          `json:"ai_active_idx"`
	TurnNumber           int          `json:"turn_number"`
	RoundNumber          int          `json:"round_number"` // For 5v5 battles
	WhoseTurn            string       `json:"whose_turn"`   // "player" or "ai"
	BattleOver           bool         `json:"battle_over"`
	Winner               string       `json:"winner"`             // "player", "ai", "draw"
	RewardClaimed        bool         `json:"reward_claimed"`     // Track if 5v5 reward has been claimed
	ConsecutivePasses    int          `json:"consecutive_passes"` // Track consecutive passes by both players
	PendingPlayerMove    string       `json:"pending_player_move"`
	PendingPlayerMoveIdx int          `json:"pending_player_move_idx"`
	PendingAIMove        string       `json:"pending_ai_move"`
	PendingAIMoveIdx     int          `json:"pending_ai_move_idx"`
	SacrificeCount       map[int]int  `json:"sacrifice_count"` // Track sacrifices per Pokemon
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
}

// BattleCard represents a Pokemon card in battle with current state
type BattleCard struct {
	CardID       int            `json:"card_id"`
	Name         string         `json:"name"`
	HP           int            `json:"hp"`
	HPMax        int            `json:"hp_max"`
	Stamina      int            `json:"stamina"`
	StaminaMax   int            `json:"stamina_max"`
	Attack       int            `json:"attack"`
	Defense      int            `json:"defense"`
	Speed        int            `json:"speed"`
	Types        []string       `json:"types"`
	Moves        []pokemon.Move `json:"moves"`
	Sprite       string         `json:"sprite"`
	IsKnockedOut bool           `json:"is_knocked_out"`
	Level        int            `json:"level"`
}

// TurnState contains web-only turn-based fields (legacy support)
type TurnState struct {
	PendingPlayerMove    string
	PendingPlayerMoveIdx int
	PendingAIMove        string
	PendingAIMoveIdx     int
	WhoseTurn            string // "player" or "ai"
}

// ConvertToBattleCard converts a pokemon.Card to a BattleCard
func ConvertToBattleCard(card pokemon.Card, cardID int) BattleCard {
	return BattleCard{
		CardID:       cardID,
		Name:         card.Name,
		HP:           card.HP,
		HPMax:        card.HPMax,
		Stamina:      card.Stamina,
		StaminaMax:   card.Speed * 2,
		Attack:       card.Attack,
		Defense:      card.Defense,
		Speed:        card.Speed,
		Types:        card.Types,
		Moves:        card.Moves,
		Sprite:       card.Sprite,
		IsKnockedOut: card.HP <= 0,
		Level:        1, // Default level, will be updated from database
	}
}

// ConvertFromBattleCard converts a BattleCard back to pokemon.Card
func ConvertFromBattleCard(bc BattleCard) pokemon.Card {
	return pokemon.Card{
		Name:    bc.Name,
		HP:      bc.HP,
		HPMax:   bc.HPMax,
		Stamina: bc.Stamina,
		Attack:  bc.Attack,
		Defense: bc.Defense,
		Speed:   bc.Speed,
		Types:   bc.Types,
		Moves:   bc.Moves,
		Sprite:  bc.Sprite,
	}
}

// GetActivePlayerCard returns the active player's BattleCard
func (bs *BattleState) GetActivePlayerCard() *BattleCard {
	if bs.PlayerActiveIdx >= 0 && bs.PlayerActiveIdx < len(bs.PlayerDeck) {
		return &bs.PlayerDeck[bs.PlayerActiveIdx]
	}
	return nil
}

// GetActiveAICard returns the active AI's BattleCard
func (bs *BattleState) GetActiveAICard() *BattleCard {
	if bs.AIActiveIdx >= 0 && bs.AIActiveIdx < len(bs.AIDeck) {
		return &bs.AIDeck[bs.AIActiveIdx]
	}
	return nil
}

// HasPlayerPokemonAlive checks if player has any Pokemon with HP > 0
func (bs *BattleState) HasPlayerPokemonAlive() bool {
	for _, card := range bs.PlayerDeck {
		if card.HP > 0 {
			return true
		}
	}
	return false
}

// HasAIPokemonAlive checks if AI has any Pokemon with HP > 0
func (bs *BattleState) HasAIPokemonAlive() bool {
	for _, card := range bs.AIDeck {
		if card.HP > 0 {
			return true
		}
	}
	return false
}

// CheckBattleEnd checks if the battle should end and updates state accordingly
func (bs *BattleState) CheckBattleEnd() {
	switch bs.Mode {
	case "1v1":
		playerCard := bs.GetActivePlayerCard()
		aiCard := bs.GetActiveAICard()
		if playerCard != nil && aiCard != nil {
			if playerCard.HP <= 0 && aiCard.HP <= 0 {
				bs.BattleOver = true
				bs.Winner = "draw"
			} else if playerCard.HP <= 0 {
				bs.BattleOver = true
				bs.Winner = "ai"
			} else if aiCard.HP <= 0 {
				bs.BattleOver = true
				bs.Winner = "player"
			}
		}
	case "5v5":
		if !bs.HasPlayerPokemonAlive() && !bs.HasAIPokemonAlive() {
			bs.BattleOver = true
			bs.Winner = "draw"
		} else if !bs.HasPlayerPokemonAlive() {
			bs.BattleOver = true
			bs.Winner = "ai"
		} else if !bs.HasAIPokemonAlive() {
			bs.BattleOver = true
			bs.Winner = "player"
		}
	}
}

// buildWebState returns a JSON-serializable map of the current state, turn info and log entries
func buildWebState(state *models.GameState, turn *TurnState, logEntries []string) map[string]any {
	winner := ""
	if state.BattleOver {
		playerCard := &state.Player.Deck[state.PlayerActiveIdx]
		aiCard := &state.AI.Deck[state.AIActiveIdx]
		if playerCard.HP <= 0 && aiCard.HP <= 0 {
			winner = "draw"
		} else if playerCard.HP <= 0 || state.PlayerSurrendered {
			winner = "ai"
		} else if aiCard.HP <= 0 {
			winner = "player"
		}
	}
	return map[string]any{
		"state":      state,
		"turn":       turn,
		"log":        logEntries,
		"battleOver": state.BattleOver,
		"whoseTurn":  turn.WhoseTurn,
		"winner":     winner,
	}
}

// BuildBattleResponse creates a response for the enhanced battle state
func BuildBattleResponse(bs *BattleState, logEntries []string, hideAICards bool) map[string]any {
	// Create a copy of the battle state for response
	response := map[string]any{
		"id":                bs.ID,
		"user_id":           bs.UserID,
		"mode":              bs.Mode,
		"player_active_idx": bs.PlayerActiveIdx,
		"ai_active_idx":     bs.AIActiveIdx,
		"turn_number":       bs.TurnNumber,
		"round_number":      bs.RoundNumber,
		"whose_turn":        bs.WhoseTurn,
		"battle_over":       bs.BattleOver,
		"winner":            bs.Winner,
		"reward_claimed":    bs.RewardClaimed,
		"log":               logEntries,
		"created_at":        bs.CreatedAt,
		"updated_at":        bs.UpdatedAt,
	}

	// Always show full player deck
	response["player_deck"] = bs.PlayerDeck

	// For 5v5 battles, hide inactive AI Pokemon details if requested
	if bs.Mode == "5v5" && hideAICards {
		aiDeck := make([]map[string]any, len(bs.AIDeck))
		for i, card := range bs.AIDeck {
			if i == bs.AIActiveIdx {
				// Show full details for active AI Pokemon
				aiDeck[i] = map[string]any{
					"card_id":        card.CardID,
					"name":           card.Name,
					"hp":             card.HP,
					"hp_max":         card.HPMax,
					"stamina":        card.Stamina,
					"stamina_max":    card.StaminaMax,
					"attack":         card.Attack,
					"defense":        card.Defense,
					"speed":          card.Speed,
					"types":          card.Types,
					"moves":          card.Moves,
					"sprite":         card.Sprite,
					"is_knocked_out": card.IsKnockedOut,
					"level":          card.Level,
					"is_active":      true,
					"is_face_down":   false,
				}
			} else {
				// Hide details for inactive AI Pokemon (face-down)
				aiDeck[i] = map[string]any{
					"card_id":        card.CardID,
					"is_knocked_out": card.IsKnockedOut,
					"is_active":      false,
					"is_face_down":   !card.IsKnockedOut, // Show face-down only if not KO'd
				}
			}
		}
		response["ai_deck"] = aiDeck
	} else {
		// For 1v1 or when not hiding, show all AI cards
		response["ai_deck"] = bs.AIDeck
	}

	return response
}
