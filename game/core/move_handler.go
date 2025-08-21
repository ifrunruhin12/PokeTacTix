package core

import (
	"bufio"
	"fmt"
	"pokemon-cli/game/models"
	"pokemon-cli/pokemon"
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
			defendCost := playerCard.Defense - int(float64(playerCard.Defense)*0.75)
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
			if count == 0 {
				hpCost = 10
			} else if count == 1 {
				hpCost = 15
			} else if count == 2 {
				hpCost = 20
			} else {
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

// Handle the sacrifice mechanic
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
