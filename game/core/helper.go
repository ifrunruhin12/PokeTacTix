package core

import (
	"fmt"
	"pokemon-cli/game/models"
	"pokemon-cli/internal/pokemon"
)

// Function to show battle result of 1v1 battle
func showBattleResult(state *models.GameState, playerCard, aiCard *pokemon.Card) {
	fmt.Println("\n--- Battle Over ---")
	if playerCard.HP <= 0 && aiCard.HP <= 0 {
		fmt.Println("It's a draw! Both Pokémon were knocked out.")
	} else if playerCard.HP <= 0 {
		fmt.Printf("You lost the battle. %s was knocked out.\n", playerCard.Name)
	} else if aiCard.HP <= 0 {
		fmt.Printf("You won the battle! AI's %s was knocked out.\n", aiCard.Name)
	} else if state.PlayerSurrendered {
		fmt.Println("You surrendered the battle. AI wins.")
	} else {
		fmt.Println("Battle ended unexpectedly.")
	}
}

// Show round summary
func showRoundSummary(state *models.GameState, playerCard, aiCard *pokemon.Card) {
	fmt.Println("\n--- Round Summary ---")
	if playerCard.HP <= 0 && aiCard.HP <= 0 {
		fmt.Println("Both Pokémon were knocked out! It's a draw for this round.")
	} else if playerCard.HP <= 0 {
		fmt.Printf("You lost the round. %s is knocked out.\n", playerCard.Name)
	} else if aiCard.HP <= 0 {
		fmt.Printf("You won the round! AI's %s is knocked out.\n", aiCard.Name)
	} else if state.RoundOver {
		if playerCard.HP <= 0 {
			fmt.Printf("You lost the round. %s is knocked out.\n", playerCard.Name)
		} else {
			fmt.Printf("You won the round! AI's %s is knocked out.\n", aiCard.Name)
		}
	} else {
		fmt.Println("Round ended unexpectedly.")
	}
}

// GetDefendCost calculates the stamina cost for defending based on the Pokemon's max HP
// Formula: (HPMax + 1) / 2
func GetDefendCost(hpMax int) int {
	return (hpMax + 1) / 2
}

// Helper function to reset battle state
func resetBattleState(state *models.GameState) {
	state.BattleStarted = false
	state.InBattle = false
	state.HaveCard = false
	state.Round = 0
	state.PlayerActiveIdx = 0
	state.AIActiveIdx = 0
	state.CardMovePlayer = 0
	state.CardMoveAI = 0
	state.CurrentMovetype = ""
	state.RoundStarted = false
	state.SwitchedThisRound = false
	state.BattleOver = false
	state.RoundOver = false
	state.SacrificeCount = nil
	state.PlayerSurrendered = false
	state.JustSwitched = false
	state.HasPlayedRound = false
	state.TurnNumber = 0
	state.BattleMode = ""
	state.PlayerName = ""
}
