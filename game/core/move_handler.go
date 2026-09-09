package core

import (
	"fmt"
	"pokemon-cli/game/models"
	"pokemon-cli/internal/pokemon"
)

// HandleSacrifice handles the sacrifice mechanic for the legacy web path.
func HandleSacrifice(state *models.GameState, playerCard *pokemon.Card) {
	idx := state.PlayerActiveIdx
	if state.SacrificeCount == nil {
		state.SacrificeCount = make(map[int]int)
	}
	count := state.SacrificeCount[idx]
	maxStamina := playerCard.Speed * 2
	if count >= 3 {
		fmt.Println("You can only sacrifice three times per round.")
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
	if float64(playerCard.Stamina) >= 0.5*float64(maxStamina) {
		fmt.Println("You can only use 'sacrifice' when your current stamina is less than 50% of max stamina.")
		return
	}
	if playerCard.HP <= hpCost {
		fmt.Printf("Not enough HP to sacrifice. You need at least %d HP.\n", hpCost+1)
		return
	}
	playerCard.HP -= hpCost
	gain := int(float64(maxStamina) * staminaGain)
	playerCard.Stamina += gain
	fmt.Printf("You sacrificed %d HP and gained %d stamina.\n", hpCost, gain)
	state.SacrificeCount[idx] = count + 1
}
