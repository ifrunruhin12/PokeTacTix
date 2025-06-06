package game

import (
	"bufio"
	"fmt"
	"strings"
)

func CommandListInBattle(state *GameState) {
	if !state.BattleStarted {
		fmt.Println("You are not in a battle yet. This commands only works when in a battle. Use the command 'battle' to start one.")
		return
	}

	fmt.Println("Available commands:")
	fmt.Println("1. card all - Show all your cards")
	fmt.Println("2. card     - Show the current card of the player")
	fmt.Println("3. attack   - Choose a move to attack with")
	fmt.Println("4. choose   - Choose a card to play")
	fmt.Println("5. command --in-battle - Show this command list")
	fmt.Println("6. exit     - Exit the game")
}

func CommandCardChooser(scanner *bufio.Scanner, state *GameState) {
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
	StartBattleLoop(scanner, state)
}

func CommandMovesAttack(scanner *bufio.Scanner, state *GameState) {
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

func CommandCardAll(state *GameState) {
	if !state.BattleStarted {
		fmt.Println("You need to start a battle first with the 'battle' command.")
		return
	}
	fmt.Printf("%s, here are your cards:\n", state.PlayerName)
	for _, card := range state.Player.AllCards() {
		PrintCard(card)
	}
}

func CommandCurrentCard(state *GameState) {
	if !state.BattleStarted {
		fmt.Println("You need to start a battle first with the 'battle' command.")
		return
	}

	if !state.HaveCard {
		fmt.Println("You need to choose a card first with the 'choose' command.")
		return
	}
	PrintCard(state.Player.AllCards()[state.PlayerActiveIdx])
}
