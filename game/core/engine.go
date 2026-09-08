// Package core contains the main game logic for the Pokémon battle system.
package core

import (
	"fmt"
	"pokemon-cli/game/models"
	"pokemon-cli/internal/pokemon"
)

// ProcessTurnResult processes the result of a turn — used by the legacy web battle path.
func ProcessTurnResult(playerMove, aiMove string, playerMoveIdx, aiMoveIdx int, playerCard, aiCard *pokemon.Card, state *models.GameState) {
	// Track for player feedback
	state.LastHpLost = 0
	state.LastStaminaLost = 0
	state.LastDamageDealt = 0

	playerDefendCost := GetDefendCost(playerCard.HPMax)
	aiDefendCost := GetDefendCost(aiCard.HPMax)

	if playerMove == "pass" && aiMove == "pass" {
		fmt.Println("Both passed. Nothing happened!")
		return
	}
	if playerMove == "pass" {
		switch aiMove {
		case "attack":
			aiDmg := CalculateDamage(aiCard, playerCard, false, aiMoveIdx)
			playerCard.HP -= aiDmg
			aiCard.Stamina -= aiCard.Moves[aiMoveIdx].StaminaCost
			state.LastHpLost = aiDmg
		case "defend":
			aiCard.Stamina -= aiDefendCost
		}
		return
	}
	if aiMove == "pass" {
		switch playerMove {
		case "attack":
			playerDmg := CalculateDamage(playerCard, aiCard, false, playerMoveIdx)
			aiCard.HP -= playerDmg
			playerCard.Stamina -= playerCard.Moves[playerMoveIdx].StaminaCost
			state.LastStaminaLost = playerCard.Moves[playerMoveIdx].StaminaCost
			state.LastDamageDealt = playerDmg
		case "defend":
			playerCard.Stamina -= playerDefendCost
			state.LastStaminaLost = playerDefendCost
		}
		return
	}
	if playerMove == "attack" && aiMove == "attack" {
		playerDmg := CalculateDamage(playerCard, aiCard, false, playerMoveIdx)
		aiDmg := CalculateDamage(aiCard, playerCard, false, aiMoveIdx)
		aiCard.HP -= playerDmg
		playerCard.HP -= aiDmg
		playerCard.Stamina -= playerCard.Moves[playerMoveIdx].StaminaCost
		aiCard.Stamina -= aiCard.Moves[aiMoveIdx].StaminaCost
		state.LastHpLost = aiDmg
		state.LastStaminaLost = playerCard.Moves[playerMoveIdx].StaminaCost
		state.LastDamageDealt = playerDmg
	} else if playerMove == "attack" && aiMove == "defend" {
		playerDmg := CalculateDamage(playerCard, aiCard, true, playerMoveIdx)
		aiCard.Stamina -= aiDefendCost
		playerCard.Stamina -= playerCard.Moves[playerMoveIdx].StaminaCost
		if playerDmg <= aiCard.Defense {
			state.LastStaminaLost = playerCard.Moves[playerMoveIdx].StaminaCost
		} else {
			aiCard.HP -= (playerDmg - aiCard.Defense)
			state.LastStaminaLost = playerCard.Moves[playerMoveIdx].StaminaCost
			state.LastDamageDealt = playerDmg - aiCard.Defense
		}
	} else if playerMove == "defend" && aiMove == "attack" {
		aiDmg := CalculateDamage(aiCard, playerCard, true, aiMoveIdx)
		playerCard.Stamina -= playerDefendCost
		aiCard.Stamina -= aiCard.Moves[aiMoveIdx].StaminaCost
		if aiDmg <= playerCard.Defense {
			state.LastStaminaLost = playerDefendCost
		} else {
			playerCard.HP -= (aiDmg - playerCard.Defense)
			state.LastHpLost = aiDmg - playerCard.Defense
			state.LastStaminaLost = playerDefendCost
		}
	} else if playerMove == "defend" && aiMove == "defend" {
		playerCard.Stamina -= playerDefendCost
		aiCard.Stamina -= aiDefendCost
		state.LastStaminaLost = playerDefendCost
	}
	// Clamp HP and stamina
	if playerCard.HP < 0 {
		playerCard.HP = 0
	}
	if aiCard.HP < 0 {
		aiCard.HP = 0
	}
	if playerCard.Stamina < 0 {
		playerCard.Stamina = 0
	}
	if aiCard.Stamina < 0 {
		aiCard.Stamina = 0
	}
}
