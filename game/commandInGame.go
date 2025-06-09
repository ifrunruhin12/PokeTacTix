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
	fmt.Println("5. switch   - Switch your active Pokémon before the round starts")
	fmt.Println("6. surrender - Surrender this round (you lose the round)")
	fmt.Println("7. surrender all - Surrender the whole battle (you lose the battle)")
	fmt.Println("8. defend   - Choose to defend against an attack")
	fmt.Println("9. command --in-battle - Show this command list")
	fmt.Println("10. exit     - Exit the game")
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
	StartTurnLoop(scanner, state)
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

func CommandDefendMove(state *GameState) {
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

// CommandSwitch allows the player to switch their active Pokémon before the round starts.
func CommandSwitch(scanner *bufio.Scanner, state *GameState) {
	if !state.BattleStarted {
		fmt.Println("You need to start a battle first with the 'battle' command.")
		return
	}
	if !state.HaveCard {
		fmt.Println("You need to choose a card first with the 'choose' command before you can switch.")
		return
	}
	
	// Check if we just switched or chose a new pokemon
	if state.JustSwitched {
		fmt.Println("You can't switch right after choosing a new Pokémon. Play at least one round first.")
		return
	}
	
	// Check if we're in the middle of a round (turn > 1)
	if state.TurnNumber > 1 {
		fmt.Println("You can only switch at the beginning of a round (Turn 1). Wait for the next round.")
		return
	}
	
	// Check if current pokemon hasn't played a complete round yet
	if !state.HasPlayedRound {
		fmt.Println("You can't switch your Pokémon right now. Your current Pokémon needs to play at least one complete round first.")
		return
	}
	
	// Check if current pokemon is knocked out (shouldn't be able to switch from KO'd pokemon)
	if state.Player.Deck[state.PlayerActiveIdx].HP <= 0 {
		fmt.Println("Your current Pokémon is knocked out. Use 'choose' to select a new one instead.")
		return
	}
	
	fmt.Print("Enter the number of the card you want to switch to: ")
	if !scanner.Scan() {
		return
	}
	input := strings.TrimSpace(scanner.Text())
	var idx int
	if n, err := fmt.Sscanf(input, "%d", &idx); err == nil && n == 1 && idx >= 1 && idx <= 5 {
		idx-- // 1-based to 0-based
		if idx == state.PlayerActiveIdx {
			fmt.Println("You are already using this Pokémon.")
			return
		}
		// Only allow switch if the chosen Pokémon is not knocked out
		if state.Player.Deck[idx].HP <= 0 {
			fmt.Println("You cannot switch to a knocked out Pokémon.")
			return
		}
		state.PlayerActiveIdx = idx
		fmt.Printf("You switched to %s. HP and Stamina remain as before.\n", state.Player.Deck[idx].Name)
		state.JustSwitched = true
		state.HasPlayedRound = false
	} else {
		fmt.Println("Invalid card number. Please enter a number between 1 and 5.")
	}
}

// CommandSurrender handles surrendering the round or the whole battle.
func CommandSurrender(scanner *bufio.Scanner, state *GameState, surrenderAll bool) {
	if !state.BattleStarted {
		fmt.Println("You need to start a battle first with the 'battle' command.")
		return
	}
	if surrenderAll {
		fmt.Println("You have surrendered the entire battle. You lose!")
		state.BattleStarted = false
		state.InBattle = false
		state.BattleOver = true
		state.PlayerSurrendered = true
		return
	}
	// Surrender just the round
	if state.RoundStarted {
		fmt.Println("You have surrendered this round. You lose the round!")
		state.Player.Deck[state.PlayerActiveIdx].HP = 0 // Mark current Pokémon as KO
		state.RoundOver = true
		// Next: prompt for choose for next round, handled in battle loop
	} else {
		fmt.Println("You can only surrender during a round.")
	}
}
