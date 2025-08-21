package web

import (
	"fmt"
	"pokemon-cli/game/core"
	"pokemon-cli/game/models"
)

// ProcessWebMove processes a move for the web API, updates the state, and returns a JSON-serializable result.
func ProcessWebMove(state *models.GameState, move string, moveIdx *int) (map[string]interface{}, error) {
	if state == nil || !state.BattleStarted || state.BattleOver {
		return nil, fmt.Errorf("no active battle or battle is over")
	}
	playerCard := &state.Player.Deck[state.PlayerActiveIdx]
	aiCard := &state.AI.Deck[state.AIActiveIdx]
	logEntries := []string{}

	// At the start of a new turn, add a turn log entry
	if state.PendingPlayerMove == "" && state.PendingAIMove == "" {
		whose := "Player's"
		if state.TurnNumber%2 == 0 {
			whose = "AI's"
		}
		logEntries = append(logEntries, fmt.Sprintf("Turn %d begins! %s move first.", state.TurnNumber, whose))
		if state.TurnNumber%2 == 1 {
			state.WhoseTurn = "player"
		} else {
			state.WhoseTurn = "ai"
		}
	}

	// Only accept moves from the correct actor
	if state.WhoseTurn == "player" {
		if move == "surrender" {
			state.BattleStarted = false
			state.InBattle = false
			state.BattleOver = true
			state.PlayerSurrendered = true
			logEntries = append(logEntries, "Player surrendered! AI wins the battle!")
			// Clear pending moves and turn
			state.PendingPlayerMove = ""
			state.PendingAIMove = ""
			state.PendingPlayerMoveIdx = 0
			state.PendingAIMoveIdx = 0
			state.WhoseTurn = ""
			return buildWebState(state, logEntries), nil
		}
		if move == "sacrifice" {
			oldHP := playerCard.HP
			oldStamina := playerCard.Stamina
			core.HandleSacrifice(state, playerCard)
			maxStamina := playerCard.Speed * 2
			if playerCard.Stamina > maxStamina {
				playerCard.Stamina = maxStamina
			}
			hpLost := oldHP - playerCard.HP
			staminaGained := playerCard.Stamina - oldStamina
			logEntries = append(logEntries, fmt.Sprintf("Player sacrificed %d HP and gained %d stamina.", hpLost, staminaGained))
			return buildWebState(state, logEntries), nil
		}
		// If this is a player turn after AI has acted (even turn), resolve the turn
		if state.PendingAIMove != "" {
			state.PendingPlayerMove = move
			if moveIdx != nil {
				state.PendingPlayerMoveIdx = *moveIdx
			} else {
				state.PendingPlayerMoveIdx = 0
			}
			// Log player move with move name if attack
			if move == "attack" {
				moveName := ""
				if state.PendingPlayerMoveIdx >= 0 && state.PendingPlayerMoveIdx < len(playerCard.Moves) {
					moveName = playerCard.Moves[state.PendingPlayerMoveIdx].Name
				}
				logEntries = append(logEntries, fmt.Sprintf("Player chose to attack with %s.", moveName))
			} else {
				logEntries = append(logEntries, fmt.Sprintf("Player chose %s.", move))
			}
			// Log AI move with move name if attack
			if state.PendingAIMove == "attack" {
				moveName := ""
				if state.PendingAIMoveIdx >= 0 && state.PendingAIMoveIdx < len(aiCard.Moves) {
					moveName = aiCard.Moves[state.PendingAIMoveIdx].Name
				}
				logEntries = append(logEntries, fmt.Sprintf("AI chose to attack with %s.", moveName))
			} else {
				logEntries = append(logEntries, fmt.Sprintf("AI chose %s.", state.PendingAIMove))
			}
			core.ProcessTurnResult(
				state.PendingPlayerMove, state.PendingAIMove,
				state.PendingPlayerMoveIdx, state.PendingAIMoveIdx,
				playerCard, aiCard, state,
			)
			if state.LastDamageDealt > 0 {
				logEntries = append(logEntries, fmt.Sprintf("Player dealt %d damage to AI.", state.LastDamageDealt))
			}
			if state.LastHpLost > 0 {
				logEntries = append(logEntries, fmt.Sprintf("AI dealt %d damage to Player.", state.LastHpLost))
			}
			if playerCard.HP <= 0 && aiCard.HP <= 0 {
				state.BattleOver = true
				logEntries = append(logEntries, "It's a draw! Both Pokémon were knocked out.")
			} else if playerCard.HP <= 0 {
				state.BattleOver = true
				logEntries = append(logEntries, "Player's Pokémon was knocked out! AI wins.")
			} else if aiCard.HP <= 0 {
				state.BattleOver = true
				logEntries = append(logEntries, "AI's Pokémon was knocked out! Player wins.")
			}
			state.TurnNumber++
			state.PendingPlayerMove = ""
			state.PendingAIMove = ""
			state.PendingPlayerMoveIdx = 0
			state.PendingAIMoveIdx = 0
			if !state.BattleOver {
				if state.TurnNumber%2 == 1 {
					state.WhoseTurn = "player"
				} else {
					state.WhoseTurn = "ai"
				}
			}
			return buildWebState(state, logEntries), nil
		}
		// Otherwise, this is a player-first turn: store move, set WhoseTurn to AI
		state.PendingPlayerMove = move
		if moveIdx != nil {
			state.PendingPlayerMoveIdx = *moveIdx
		} else {
			state.PendingPlayerMoveIdx = 0
		}
		state.WhoseTurn = "ai"
		return buildWebState(state, logEntries), nil
	} else if state.WhoseTurn == "ai" {
		for {
			aiMove, aiMoveIdx := core.GetAIMove(state.PendingPlayerMove, aiCard, state, state.AIActiveIdx)
			if aiMove == "surrender" {
				state.BattleStarted = false
				state.InBattle = false
				state.BattleOver = true
				logEntries = append(logEntries, "AI surrendered! Player wins the battle!")
				// Clear pending moves and turn
				state.PendingPlayerMove = ""
				state.PendingAIMove = ""
				state.PendingPlayerMoveIdx = 0
				state.PendingAIMoveIdx = 0
				state.WhoseTurn = ""
				return buildWebState(state, logEntries), nil
			}
			if aiMove == "sacrifice" {
				maxStamina := aiCard.Speed * 2
				if float64(aiCard.Stamina) >= 0.5*float64(maxStamina) {
					break
				}
				oldHP := aiCard.HP
				oldStamina := aiCard.Stamina
				core.HandleSacrificeAI(aiCard, state)
				if aiCard.Stamina > maxStamina {
					aiCard.Stamina = maxStamina
				}
				hpLost := oldHP - aiCard.HP
				staminaGained := aiCard.Stamina - oldStamina
				logEntries = append(logEntries, fmt.Sprintf("AI sacrificed %d HP and gained %d stamina.", hpLost, staminaGained))
				continue
			}
			// If this is an AI-first turn (even turn), store AI's move, set WhoseTurn to player, and return (do NOT resolve turn yet)
			if state.PendingPlayerMove == "" {
				state.PendingAIMove = aiMove
				state.PendingAIMoveIdx = aiMoveIdx
				// Log AI move with move name if attack
				if aiMove == "attack" {
					moveName := ""
					if aiMoveIdx >= 0 && aiMoveIdx < len(aiCard.Moves) {
						moveName = aiCard.Moves[aiMoveIdx].Name
					}
					logEntries = append(logEntries, fmt.Sprintf("AI chose to attack with %s.", moveName))
				} else {
					logEntries = append(logEntries, fmt.Sprintf("AI chose %s.", aiMove))
				}
				state.WhoseTurn = "player"
				return buildWebState(state, logEntries), nil
			}
			// Otherwise, this is AI acting second (player-first turn): resolve the turn
			state.PendingAIMove = aiMove
			state.PendingAIMoveIdx = aiMoveIdx
			// Log player move with move name if attack
			if state.PendingPlayerMove == "attack" {
				moveName := ""
				if state.PendingPlayerMoveIdx >= 0 && state.PendingPlayerMoveIdx < len(playerCard.Moves) {
					moveName = playerCard.Moves[state.PendingPlayerMoveIdx].Name
				}
				logEntries = append(logEntries, fmt.Sprintf("Player chose to attack with %s.", moveName))
			} else {
				logEntries = append(logEntries, fmt.Sprintf("Player chose %s.", state.PendingPlayerMove))
			}
			// Log AI move with move name if attack
			if aiMove == "attack" {
				moveName := ""
				if aiMoveIdx >= 0 && aiMoveIdx < len(aiCard.Moves) {
					moveName = aiCard.Moves[aiMoveIdx].Name
				}
				logEntries = append(logEntries, fmt.Sprintf("AI chose to attack with %s.", moveName))
			} else {
				logEntries = append(logEntries, fmt.Sprintf("AI chose %s.", aiMove))
			}
			core.ProcessTurnResult(
				state.PendingPlayerMove, state.PendingAIMove,
				state.PendingPlayerMoveIdx, state.PendingAIMoveIdx,
				playerCard, aiCard, state,
			)
			if state.LastDamageDealt > 0 {
				logEntries = append(logEntries, fmt.Sprintf("Player dealt %d damage to AI.", state.LastDamageDealt))
			}
			if state.LastHpLost > 0 {
				logEntries = append(logEntries, fmt.Sprintf("AI dealt %d damage to Player.", state.LastHpLost))
			}
			if playerCard.HP <= 0 && aiCard.HP <= 0 {
				state.BattleOver = true
				logEntries = append(logEntries, "It's a draw! Both Pokémon were knocked out.")
			} else if playerCard.HP <= 0 {
				state.BattleOver = true
				logEntries = append(logEntries, "Player's Pokémon was knocked out! AI wins.")
			} else if aiCard.HP <= 0 {
				state.BattleOver = true
				logEntries = append(logEntries, "AI's Pokémon was knocked out! Player wins.")
			}
			state.TurnNumber++
			state.PendingPlayerMove = ""
			state.PendingAIMove = ""
			state.PendingPlayerMoveIdx = 0
			state.PendingAIMoveIdx = 0
			if !state.BattleOver {
				if state.TurnNumber%2 == 1 {
					state.WhoseTurn = "player"
				} else {
					state.WhoseTurn = "ai"
				}
			}
			return buildWebState(state, logEntries), nil
		}
	} else {
		return buildWebState(state, logEntries), nil
	}
	return buildWebState(state, logEntries), nil
}
