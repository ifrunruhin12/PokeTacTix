// Package core contains the main game logic for the Pokémon CLI battle system.
package core

import (
	"bufio"
	"fmt"
	"pokemon-cli/game/models"
	"pokemon-cli/game/utils"
	"pokemon-cli/pokemon"
	"strings"
)

func StartTurnLoop(scanner *bufio.Scanner, state *models.GameState) {
	turn := 1
	player := state.Player
	ai := state.AI
	// Reset sacrifice count for this round
	if state.SacrificeCount == nil {
		state.SacrificeCount = make(map[int]int)
	}
	state.SacrificeCount[state.PlayerActiveIdx] = 0

	if state.BattleMode == "1v1" {
		fmt.Printf("\nThe 1v1 battle has started! It's Turn 1 for both player and AI.\n")
		fmt.Println("This battle will continue until one Pokémon is knocked out or surrenders.")
	} else {
		fmt.Printf("\nThe %s Round of the battle has started! It's Turn 1 for both player and AI.\n", utils.Ordinal(state.Round))
	}

	fmt.Println("Choose your move with the correct command. To see all the in-battle commands type 'command --in-battle'.")

	for {
		playerCard := &player.Deck[state.PlayerActiveIdx]
		aiCard := &ai.Deck[state.AIActiveIdx]

		if playerCard.HP <= 0 || aiCard.HP <= 0 || state.RoundOver || state.BattleOver {
			break
		}

		fmt.Printf("\n--- Turn %d ---\n", turn)
		var playerMove, aiMove string
		var playerMoveIdx, aiMoveIdx int
		// Odd turns: player chooses first, then AI
		if turn%2 == 1 {
			// Player's move (with sacrifice as free action)
			for {
				playerMove, playerMoveIdx = getPlayerMove(scanner, state, playerCard)
				if playerMove == "sacrifice" {
					HandleSacrifice(state, playerCard)
					continue // re-prompt player for a real move
				}
				break
			}
			if playerMove == "surrender" {
				if state.BattleMode == "1v1" {
					fmt.Println("You surrendered the battle!")
					state.BattleStarted = false
					state.InBattle = false
					state.BattleOver = true
					state.PlayerSurrendered = true
					break
				} else {
					fmt.Println("You surrendered this round!")
					playerCard.HP = 0
					state.RoundOver = true
					break
				}
			}

			if playerMove == "surrender all" {
				break
			}
			// AI's move (with sacrifice as free action)
			for {
				aiMove, aiMoveIdx = GetAIMove(playerMove, aiCard, state, state.AIActiveIdx)
				if aiMove == "sacrifice" {
					HandleSacrificeAI(aiCard, state)
					continue // re-prompt AI for a real move
				}
				if aiMove == "surrender" {
					if state.BattleMode == "1v1" {
						fmt.Println("AI surrendered the battle!")
						state.BattleStarted = false
						state.InBattle = false
						state.BattleOver = true
						break
					} else {
						fmt.Println("AI surrendered this round!")
						aiCard.HP = 0
						state.RoundOver = true
						break
					}
				}
				break
			}
			fmt.Printf("AI chose to %s.\n", aiMove)
			// Re-assign playerCard and aiCard after possible switch
			playerCard = &player.Deck[state.PlayerActiveIdx]
			aiCard = &ai.Deck[state.AIActiveIdx]
		} else {
			// Even turns: AI chooses first, then player
			for {
				aiMove, aiMoveIdx = GetAIMove("", aiCard, state, state.AIActiveIdx)
				if aiMove == "sacrifice" {
					HandleSacrificeAI(aiCard, state)
					continue // re-prompt AI for a real move
				}
				if aiMove == "surrender" {
					if state.BattleMode == "1v1" {
						fmt.Println("AI surrendered the battle!")
						state.BattleStarted = false
						state.InBattle = false
						state.BattleOver = true
						break
					} else {
						fmt.Println("AI surrendered this round!")
						aiCard.HP = 0
						state.RoundOver = true
						break
					}
				}
				break
			}
			fmt.Printf("AI chose to %s.\n", aiMove)
			for {
				playerMove, playerMoveIdx = getPlayerMove(scanner, state, playerCard)
				if playerMove == "sacrifice" {
					HandleSacrifice(state, playerCard)
					continue // re-prompt player for a real move
				}
				break
			}
			if playerMove == "surrender" {
				if state.BattleMode == "1v1" {
					fmt.Println("You surrendered the battle!")
					state.BattleStarted = false
					state.InBattle = false
					state.BattleOver = true
					state.PlayerSurrendered = true
					break
				} else {
					fmt.Println("You surrendered this round!")
					playerCard.HP = 0
					state.RoundOver = true
					break
				}
			}

			if playerMove == "surrender all" {
				break
			}
			// Re-assign playerCard and aiCard after possible switch
			playerCard = &player.Deck[state.PlayerActiveIdx]
			aiCard = &ai.Deck[state.AIActiveIdx]
		}
		// Process moves
		ProcessTurnResult(playerMove, aiMove, playerMoveIdx, aiMoveIdx, playerCard, aiCard, state)
		// Show result (only player's HP/stamina lost and damage dealt)
		fmt.Printf("You lost %d HP and %d stamina this turn. Your current HP: %d, current stamina: %d\n", state.LastHpLost, state.LastStaminaLost, playerCard.HP, playerCard.Stamina)
		fmt.Printf("You dealt %d damage to the AI.\n", state.LastDamageDealt)
		// Check for KO
		if playerCard.HP <= 0 {
			fmt.Printf("Your %s was knocked out!\n", playerCard.Name)
			state.RoundOver = true
			break
		}
		if aiCard.HP <= 0 {
			fmt.Printf("AI's %s was knocked out!\n", aiCard.Name)
			state.RoundOver = true
			break
		}
		turn++
	}
	// End of round
	if state.BattleMode == "1v1" {
		showBattleResult(state, &player.Deck[0], &ai.Deck[0])
		resetBattleState(state)
	} else {
		showRoundSummary(state, &player.Deck[state.PlayerActiveIdx], &ai.Deck[state.AIActiveIdx])
		prepareNextRound(scanner, state)
	}
}

// ProcessTurnResult Process the result of a turn
func ProcessTurnResult(playerMove, aiMove string, playerMoveIdx, aiMoveIdx int, playerCard, aiCard *pokemon.Card, state *models.GameState) {
	// Track for player feedback
	state.LastHpLost = 0
	state.LastStaminaLost = 0
	state.LastDamageDealt = 0

	playerDefendCost := (playerCard.HPMax + 1) / 2
	aiDefendCost := (aiCard.HPMax + 1) / 2

	if playerMove == "pass" && aiMove == "pass" {
		fmt.Println("Both passed. Nothing happened!")
		return
	}
	if playerMove == "pass" {
		// AI does its move, player does nothing
		switch aiMove {
		case "attack":
			aiDmg := calculateDamage(aiCard, playerCard, false, aiMoveIdx)
			playerCard.HP -= aiDmg
			aiCard.Stamina -= aiCard.Moves[aiMoveIdx].StaminaCost
			state.LastHpLost = aiDmg
			state.LastStaminaLost = 0
			state.LastDamageDealt = 0
		case "defend":
			aiCard.Stamina -= aiDefendCost
			state.LastHpLost = 0
			state.LastStaminaLost = 0
			state.LastDamageDealt = 0
		}
		return
	}
	if aiMove == "pass" {
		// Player does their move, AI does nothing
		switch playerMove {
		case "attack":
			playerDmg := calculateDamage(playerCard, aiCard, false, playerMoveIdx)
			aiCard.HP -= playerDmg
			playerCard.Stamina -= playerCard.Moves[playerMoveIdx].StaminaCost
			state.LastHpLost = 0
			state.LastStaminaLost = playerCard.Moves[playerMoveIdx].StaminaCost
			state.LastDamageDealt = playerDmg
		case "defend":
			playerCard.Stamina -= playerDefendCost
			state.LastHpLost = 0
			state.LastStaminaLost = playerDefendCost
			state.LastDamageDealt = 0
		}
		return
	}
	if playerMove == "attack" && aiMove == "attack" {
		playerDmg := calculateDamage(playerCard, aiCard, false, playerMoveIdx)
		aiDmg := calculateDamage(aiCard, playerCard, false, aiMoveIdx)
		aiCard.HP -= playerDmg
		playerCard.HP -= aiDmg
		playerCard.Stamina -= playerCard.Moves[playerMoveIdx].StaminaCost
		aiCard.Stamina -= aiCard.Moves[aiMoveIdx].StaminaCost
		state.LastHpLost = aiDmg
		state.LastStaminaLost = playerCard.Moves[playerMoveIdx].StaminaCost
		state.LastDamageDealt = playerDmg
	} else if playerMove == "attack" && aiMove == "defend" {
		playerDmg := calculateDamage(playerCard, aiCard, true, playerMoveIdx)
		aiCard.Stamina -= aiDefendCost
		playerCard.Stamina -= playerCard.Moves[playerMoveIdx].StaminaCost
		if playerDmg <= aiCard.Defense {
			// AI blocks all damage
			state.LastHpLost = 0
			state.LastStaminaLost = playerCard.Moves[playerMoveIdx].StaminaCost
			state.LastDamageDealt = 0
		} else {
			aiCard.HP -= (playerDmg - aiCard.Defense)
			state.LastHpLost = 0
			state.LastStaminaLost = playerCard.Moves[playerMoveIdx].StaminaCost
			state.LastDamageDealt = playerDmg - aiCard.Defense
		}
	} else if playerMove == "defend" && aiMove == "attack" {
		aiDmg := calculateDamage(aiCard, playerCard, true, aiMoveIdx)
		playerCard.Stamina -= playerDefendCost
		aiCard.Stamina -= aiCard.Moves[aiMoveIdx].StaminaCost
		if aiDmg <= playerCard.Defense {
			// Player blocks all damage
			state.LastHpLost = 0
			state.LastStaminaLost = playerDefendCost
			state.LastDamageDealt = 0
		} else {
			playerCard.HP -= (aiDmg - playerCard.Defense)
			state.LastHpLost = aiDmg - playerCard.Defense
			state.LastStaminaLost = playerDefendCost
			state.LastDamageDealt = 0
		}
	} else if playerMove == "defend" && aiMove == "defend" {
		playerCard.Stamina -= playerDefendCost
		aiCard.Stamina -= aiDefendCost
		state.LastHpLost = 0
		state.LastStaminaLost = playerDefendCost
		state.LastDamageDealt = 0
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

// Prepare for the next round or end the battle
func prepareNextRound(scanner *bufio.Scanner, state *models.GameState) {
	if state.BattleMode == "1v1" {
		return
	}
	player := state.Player
	ai := state.AI
	// Check if player surrendered the whole battle
	if state.PlayerSurrendered {
		fmt.Println("\n--- Battle Over ---")
		fmt.Println("Player fully surrendered. AI won the battle.")
		// Reset all battle-related state
		resetBattleState(state)
		return
	}

	// Check if either side has usable Pokémon left
	playerAlive := false
	for _, c := range player.Deck {
		if c.HP > 0 {
			playerAlive = true
			break
		}
	}
	aiAlive := false
	for _, c := range ai.Deck {
		if c.HP > 0 {
			aiAlive = true
			break
		}
	}
	if !playerAlive || !aiAlive || state.BattleOver {
		fmt.Println("\n--- Battle Over ---")
		if !playerAlive && !aiAlive {
			fmt.Println("It's a draw! Both sides are out of Pokémon.")
		} else if !playerAlive {
			fmt.Println("You lost the battle. All your Pokémon are knocked out or surrendered.")
		} else {
			fmt.Println("You won the battle! The AI is out of Pokémon.")
		}
		// Reset all battle-related state
		resetBattleState(state)
		return
	}

	// Next round
	state.Round++
	state.RoundOver = false
	state.RoundStarted = false
	state.SwitchedThisRound = false
	fmt.Printf("\nPrepare for the %s round.\n", utils.Ordinal(state.Round))
	// Winner can switch, loser must choose
	if state.Player.Deck[state.PlayerActiveIdx].HP > 0 {
		fmt.Println("You can use 'switch' to change your Pokémon, or continue with the same one (HP and stamina are not restored).")
	} else {
		fmt.Println("Your Pokémon is knocked out. Use 'choose' to pick a new one (cannot pick knocked out Pokémon).")
		for {
			fmt.Print("Enter the number of the card you want to choose: ")
			if !scanner.Scan() {
				return
			}
			input := strings.TrimSpace(scanner.Text())
			var idx int
			if n, err := fmt.Sscanf(input, "%d", &idx); err == nil && n == 1 && idx >= 1 && idx <= 5 {
				idx--
				if state.Player.Deck[idx].HP > 0 {
					state.PlayerActiveIdx = idx
					state.JustSwitched = true
					break
				} else {
					fmt.Println("That Pokémon is knocked out. Choose another.")
				}
			} else {
				fmt.Println("Invalid card number. Please enter a number between 1 and 5.")
			}
		}
	}
	// AI chooses next available or switches if low HP/stamina
	aiCurrent := &state.AI.Deck[state.AIActiveIdx]
	aiShouldSwitch := false
	if aiCurrent.HP > 0 {
		lowHP := float64(aiCurrent.HP) < 0.3*float64(aiCurrent.HPMax)
		lowStamina := float64(aiCurrent.Stamina) < 0.3*float64(aiCurrent.Stamina)
		if lowHP || lowStamina {
			bestIdx := state.AIActiveIdx
			bestScore := aiCurrent.HP + aiCurrent.Stamina
			for i, c := range state.AI.Deck {
				if c.HP > 0 && (c.HP+c.Stamina) > bestScore {
					bestIdx = i
					bestScore = c.HP + c.Stamina
				}
			}
			if bestIdx != state.AIActiveIdx {
				state.AIActiveIdx = bestIdx
				aiShouldSwitch = true
			}
		}
	} else {
		// AI's current is KO, must pick next available
		for i, c := range state.AI.Deck {
			if c.HP > 0 {
				state.AIActiveIdx = i
				break
			}
		}
	}
	if aiShouldSwitch {
		fmt.Printf("AI switched to %s for this round.\n", state.AI.Deck[state.AIActiveIdx].Name)
	}
	// Start next round
	StartTurnLoop(scanner, state)
}
