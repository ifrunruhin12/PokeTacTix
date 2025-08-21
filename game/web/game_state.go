package web

import "pokemon-cli/game/models"

// TurnState contains web-only turn-based fields
// This keeps web transport/session concerns out of core game state.
type TurnState struct {
	PendingPlayerMove    string
	PendingPlayerMoveIdx int
	PendingAIMove        string
	PendingAIMoveIdx     int
	WhoseTurn            string // "player" or "ai"
}

// buildWebState returns a JSON-serializable map of the current state, turn info and log entries
func buildWebState(state *models.GameState, turn *TurnState, logEntries []string) map[string]interface{} {
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
	return map[string]interface{}{
		"state":      state,
		"turn":       turn,
		"log":        logEntries,
		"battleOver": state.BattleOver,
		"whoseTurn":  turn.WhoseTurn,
		"winner":     winner,
	}
}
