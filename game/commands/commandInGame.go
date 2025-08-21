package commands

import (
	"bufio"
	"fmt"
	"pokemon-cli/game/core"
	"pokemon-cli/game/models"
	"strings"
)

func CommandListInBattle(state *models.GameState) {
	if !state.BattleStarted {
		fmt.Println("You are not in a battle.")
		return
	}

	fmt.Println("Available commands in battle:")
	fmt.Println("1. attack    - Attack the opponent")
	fmt.Println("2. defend    - Defend against opponent's attack")
	fmt.Println("3. sacrifice - Sacrifice HP to gain stamina (when stamina < 50%)")
	fmt.Println("4. pass      - Do nothing this turn")

	if state.BattleMode == "1v1" {
		fmt.Println("5. surrender - Surrender the battle")
	} else {
		fmt.Println("5. surrender - Surrender current round")
		fmt.Println("6. surrender all - Surrender the entire battle")
	}

	fmt.Println("7. card      - View your current active card")

	if state.BattleMode == "5v5" {
		fmt.Println("8. card all  - View all your cards")
		fmt.Println("9. choose    - Choose which card to play")
		fmt.Println("10. switch   - Switch to a different Pokémon (only at start of round)")
	}
}

func CommandCardChooser(scanner *bufio.Scanner, state *models.GameState) {
	if !state.BattleStarted {
		fmt.Println("You need to start a battle first with the 'battle' command.")
		return
	}

	fmt.Print("Enter the number of the card you want to choose: ")
	if !scanner.Scan() {
		return
	}
	input := strings.TrimSpace(scanner.Text())
	if n, err := fmt.Sscanf(input, "%d", &state.PlayerActiveIdx); err == nil && n == 1 && state.PlayerActiveIdx >= 1 && state.PlayerActiveIdx <= 5 {
		state.PlayerActiveIdx-- // 1-based to 0-based
	} else {
		fmt.Println("Invalid card number. Please enter a number between 1 and 5.")
	}
	state.HaveCard = true
	core.StartTurnLoop(scanner, state)
}

func CommandMovesAttack(scanner *bufio.Scanner, state *models.GameState) {
	if !state.BattleStarted {
		fmt.Println("You need to start a battle first with the 'battle' command.")
		return
	}

	if !state.HaveCard {
		fmt.Println("You need to choose a card first with the 'choose' command and then attack.")
		return
	}

	fmt.Print("Choose a move to attack with. Enter the number: ")
	if !scanner.Scan() {
		return
	}

	input := strings.TrimSpace(scanner.Text())
	if n, err := fmt.Sscanf(input, "%d", &state.CardMovePlayer); err == nil && n == 1 && state.CardMovePlayer >= 1 && state.CardMovePlayer <= 4 {
		state.CardMovePlayer-- // 1-based to 0-based
	} else {
		fmt.Println("Invalid move number. Please enter a number between 1 and 4.")
	}

	state.CurrentMovetype = "attack"
}

func CommandDefendMove(state *models.GameState) {
	if !state.BattleStarted {
		fmt.Println("You need to start a battle first with the 'battle' command.")
		return
	}

	if !state.HaveCard {
		fmt.Println("You need to choose a card first with the 'choose' command and then attack.")
		return
	}
	state.CurrentMovetype = "defend"
}

func CommandCardAll(state *models.GameState) {
	if !state.BattleStarted {
		fmt.Println("You need to start a battle first with the 'battle' command.")
		return
	}
	fmt.Printf("%s, here are your cards:\n", state.PlayerName)
	for _, card := range state.Player.AllCards() {
		models.PrintCard(card)
	}
}

func CommandCurrentCard(state *models.GameState) {
	if !state.BattleStarted {
		fmt.Println("You are not in a battle.")
		return
	}

	if state.BattleMode == "1v1" {
		// In 1v1, there's only one card at index 0
		models.PrintCard(state.Player.Deck[0])
	} else {
		// In 5v5, show the active card
		if !state.HaveCard {
			fmt.Println("You haven't chosen a card yet. Use 'choose' to select one.")
			return
		}
		models.PrintCard(state.Player.Deck[state.PlayerActiveIdx])
	}
}

// Delegate to core implementation
func CommandSwitch(scanner *bufio.Scanner, state *models.GameState) {
	core.CommandSwitch(scanner, state)
}

// Delegate to core implementation
func CommandSurrender(scanner *bufio.Scanner, state *models.GameState, surrenderAll bool) {
	core.CommandSurrender(scanner, state, surrenderAll)
}
