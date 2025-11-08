package core

import (
	"math/rand"
	"pokemon-cli/game/models"
	"pokemon-cli/pokemon"
)

// GetAIMove gets AI move logic
func GetAIMove(playerMove string, aiCard *pokemon.Card, state *models.GameState, aiIdx int) (string, int) {
	maxStamina := aiCard.Speed * 2
	// Check if AI can attack or defend
	canAttack := false
	minAttackCost := 9999
	for _, move := range aiCard.Moves {
		if aiCard.Stamina >= move.StaminaCost {
			canAttack = true
			break
		}
		if move.StaminaCost < minAttackCost {
			minAttackCost = move.StaminaCost
		}
	}
	defendCost := GetDefendCost(aiCard.HPMax)
	canDefend := aiCard.Stamina >= defendCost
	count := 0
	if state != nil && state.SacrificeCount != nil {
		count = state.SacrificeCount[aiIdx]
	}
	// If AI can't attack or defend
	if !canAttack && !canDefend {
		// Can AI sacrifice?
		canSacrifice := false
		var hpCost int
		switch count {
		case 0:
			hpCost = 10
		case 1:
			hpCost = 15
		case 2:
			hpCost = 20
		default:
			hpCost = 9999
		}
		if float64(aiCard.Stamina) < 0.5*float64(maxStamina) && aiCard.HP > hpCost && count < 3 {
			canSacrifice = true
		}
		if canSacrifice {
			if rand.Float64() < 0.99 {
				return "sacrifice", 0
			} else {
				return "surrender", 0
			}
		} else {
			if rand.Float64() < 0.95 {
				return "surrender", 0
			} else {
				return "pass", 0
			}
		}
	}
	// If player passed, AI always attacks if possible
	if playerMove == "pass" && canAttack {
		moveIdx := rand.Intn(len(aiCard.Moves))
		return "attack", moveIdx
	}
	// If player attacks, AI defends 66% of the time, attacks 34%
	if playerMove == "attack" {
		if rand.Float64() < 0.66 && canDefend {
			return "defend", 0
		}
		if canAttack {
			moveIdx := rand.Intn(len(aiCard.Moves))
			return "attack", moveIdx
		}
	}
	// If player defends, AI always attacks if possible
	if playerMove == "defend" && canAttack {
		moveIdx := rand.Intn(len(aiCard.Moves))
		return "attack", moveIdx
	}
	// Default: attack if possible, else pass
	if canAttack {
		moveIdx := rand.Intn(len(aiCard.Moves))
		return "attack", moveIdx
	}
	return "pass", 0
}

// HandleSacrificeAI handles AI sacrifice logic
func HandleSacrificeAI(aiCard *pokemon.Card, state *models.GameState) {
	aiIdx := state.AIActiveIdx
	if state.SacrificeCount == nil {
		state.SacrificeCount = make(map[int]int)
	}
	count := state.SacrificeCount[aiIdx]
	maxStamina := int(float64(aiCard.HPMax) * 2.5)
	if count >= 3 {
		return
	}
	var hpCost int
	var staminaGain float64
	switch count {
	case 0:
		hpCost = 10
		staminaGain = 0.5
	case 1:
		hpCost = 15
		staminaGain = 0.25
	case 2:
		hpCost = 20
		staminaGain = 0.15
	}
	if float64(aiCard.Stamina) >= 0.5*float64(maxStamina) {
		return
	}
	if aiCard.HP <= hpCost {
		return
	}
	aiCard.HP -= hpCost
	gain := int(float64(maxStamina) * staminaGain)
	aiCard.Stamina += gain
	state.SacrificeCount[aiIdx] = count + 1
}
