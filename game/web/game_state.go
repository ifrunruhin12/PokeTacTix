package web

import "pokemon-cli/game/models"

// buildWebState returns a JSON-serializable map of the current state and log entries
func buildWebState(state *models.GameState, logEntries []string) map[string]interface{} {
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
		"log":        logEntries,
		"battleOver": state.BattleOver,
		"whoseTurn":  state.WhoseTurn,
		"winner":     winner,
	}
}
