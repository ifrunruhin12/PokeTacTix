package core

import (
	"bufio"
	"fmt"
	"pokemon-cli/game/models"
	"pokemon-cli/internal/pokemon"
	"strings"
)

// Get the player's move (attack/defend/surrender/sacrifice/pass)
func getPlayerMove(scanner *bufio.Scanner, state *models.GameState, playerCard *pokemon.Card) (string, int) {
	for {
		if state.BattleMode == "5v5" {
			fmt.Print("Enter your move (attack/defend/surrender/sacrifice/pass). To see the card you are battling with use the 'card' command. To end/lose the game use 'surrender all' command. You can use the 'switch' command to switch to different pokemon if the requirements met: ")
		} else {
			fmt.Print("Enter your move (attack/defend/surrender/sacrifice/pass). To see the card you are battling with use the 'card' command. To end/lose the game use 'surrender' command: ")
		}
		if !scanner.Scan() {
			return "surrender", 0 // treat EOF as surrender
		}
		move := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if move == "card" {
			models.PrintCard(*playerCard)
			continue // re-prompt without printing the prompt again
		}
		if move == "surrender all" {
			if state.BattleMode == "1v1" {
				fmt.Println("'surrender all' is not available in 1v1 battles. Use 'surrender' to end the battle.")
				continue
			}
			CommandSurrender(scanner, state, true)
			return "surrender all", 0
		}
		if move == "switch" {
			if state.BattleMode == "1v1" {
				fmt.Println("'switch' is not available in 1v1 battles. You only have one Pokémon.")
				continue
			}
			// Switch logic - only allow at turn 1, after playing a round, and if current pokemon won previous round
			if state.TurnNumber > 1 {
				fmt.Println("You can only switch at the beginning of a round (Turn 1). Wait for the next round.")
				continue
			}
			if state.JustSwitched {
				fmt.Println("You can't switch right after choosing a new Pokémon. Play at least one round first.")
				continue
			}
			if !state.HasPlayedRound {
				fmt.Println("You can't switch your Pokémon right now. Your current Pokémon needs to play at least one complete round first.")
				continue
			}
			if playerCard.HP <= 0 {
				fmt.Println("Your current Pokémon is knocked out. Use 'choose' to select a new one instead.")
				continue
			}
			CommandSwitch(scanner, state)
			// Update playerCard reference after potential switch
			playerCard = &state.Player.Deck[state.PlayerActiveIdx]
			continue
		}

		if move == "attack" {
			moveIdx := 0
			if len(playerCard.Moves) > 1 {
				fmt.Print("Choose a move number: ")
				if !scanner.Scan() {
					return "surrender", 0
				}
				input := strings.TrimSpace(scanner.Text())
				if n, err := fmt.Sscanf(input, "%d", &moveIdx); err != nil || n != 1 || moveIdx < 1 || moveIdx > len(playerCard.Moves) {
					fmt.Println("Invalid move number.")
					continue
				}
				moveIdx--
			}
			staminaCost := playerCard.Moves[moveIdx].StaminaCost
			if playerCard.Stamina < staminaCost {
				fmt.Println("Not enough stamina to attack. Use 'sacrifice', 'pass', or 'surrender'.")
				continue
			}
			state.CardMovePlayer = moveIdx
			state.JustSwitched = false
			return "attack", moveIdx
		}
		if move == "defend" {
			defendCost := GetDefendCost(playerCard.HPMax)
			if playerCard.Stamina < defendCost {
				fmt.Println("Not enough stamina to defend. Use 'sacrifice', 'pass', or 'surrender'.")
				continue
			}
			state.JustSwitched = false
			return "defend", 0
		}
		if move == "sacrifice" {
			idx := state.PlayerActiveIdx
			if state.SacrificeCount == nil {
				state.SacrificeCount = make(map[int]int)
			}
			count := state.SacrificeCount[idx]
			var hpCost int
			switch count {
			case 0:
				hpCost = 10
			case 1:
				hpCost = 15
			case 2:
				hpCost = 20
			default:
				fmt.Println("You can only sacrifice three times per round.")
				continue
			}
			if float64(playerCard.Stamina) >= 0.5*float64(playerCard.Speed*2) {
				fmt.Println("You can only use 'sacrifice' when your current stamina is less than 50% of max stamina.")
				continue
			}
			if playerCard.HP <= hpCost {
				fmt.Printf("Not enough HP to sacrifice. You need at least %d HP.\n", hpCost+1)
				continue
			}
			state.JustSwitched = false
			return "sacrifice", 0
		}
		if move == "pass" {
			state.JustSwitched = false
			return "pass", 0
		}
		if move == "surrender" {
			return "surrender", 0
		}
		if state.BattleMode == "5v5" {
			fmt.Println("Invalid move. Please enter 'attack', 'defend', 'sacrifice', 'pass', 'switch', or 'surrender'.")
		} else {
			fmt.Println("Invalid move. Please enter 'attack', 'defend', 'sacrifice', 'pass', or 'surrender'.")
		}
	}
}

// HandleSacrifice handles the sacrifice mechanic
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

// CommandSwitch allows the player to switch their active Pokémon before the round starts.
func CommandSwitch(scanner *bufio.Scanner, state *models.GameState) {
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
func CommandSurrender(scanner *bufio.Scanner, state *models.GameState, surrenderAll bool) {
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
