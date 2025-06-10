package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"pokemon-cli/pokemon"
	"strings"
)

// Main turn/round loop
func StartTurnLoop(scanner *bufio.Scanner, state *GameState) {
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
		fmt.Printf("\nThe %s Round of the battle has started! It's Turn 1 for both player and AI.\n", ordinal(state.Round))
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
			playerMove, playerMoveIdx = getPlayerMove(scanner, state, playerCard)
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
			if playerMove == "sacrifice" {
				handleSacrifice(state, playerCard)
				continue
			}
			aiMove, aiMoveIdx = getAIMove(playerMove, aiCard, state, state.AIActiveIdx)
			// AI can use sacrifice if stamina < 50% of max
			if aiMove == "sacrifice" {
				handleSacrificeAI(aiCard, state)
				continue
			}
			if aiMove == "surrender" {
				if state.BattleMode == "1v1" {
					fmt.Println("AI surrendered the battle!")
					state.BattleOver = true
					break
				} else {
					fmt.Println("AI surrendered this round!")
					aiCard.HP = 0
					state.RoundOver = true
					break
				}
			}
			fmt.Printf("AI chose to %s.\n", aiMove)
			// Re-assign playerCard and aiCard after possible switch
			playerCard = &player.Deck[state.PlayerActiveIdx]
			aiCard = &ai.Deck[state.AIActiveIdx]
		} else {
			// Even turns: AI chooses first, then player
			aiMove, aiMoveIdx = getAIMove("", aiCard, state, state.AIActiveIdx)
			// AI can use sacrifice if stamina < 50% of max
			if aiMove == "sacrifice" {
				handleSacrificeAI(aiCard, state)
				continue
			}
			if aiMove == "surrender" {
				if state.BattleMode == "1v1" {
					fmt.Println("AI surrendered the battle!")
					state.BattleOver = true
					break
				} else {
					fmt.Println("AI surrendered this round!")
					aiCard.HP = 0
					state.RoundOver = true
					break
				}
			}
			fmt.Printf("AI chose to %s.\n", aiMove)
			playerMove, playerMoveIdx = getPlayerMove(scanner, state, playerCard)

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
			if playerMove == "sacrifice" {
				handleSacrifice(state, playerCard)
				continue
			}
			// Re-assign playerCard and aiCard after possible switch
			playerCard = &player.Deck[state.PlayerActiveIdx]
			aiCard = &ai.Deck[state.AIActiveIdx]
		}
		// Process moves
		processTurnResult(playerMove, aiMove, playerMoveIdx, aiMoveIdx, playerCard, aiCard, state)
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


//Function to show battle result of 1v1 battle
func showBattleResult(state *GameState, playerCard, aiCard *pokemon.Card) {
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

// Get the player's move (attack/defend/surrender/sacrifice/pass)
func getPlayerMove(scanner *bufio.Scanner, state *GameState, playerCard *pokemon.Card) (string, int) {
	for {
		if state.BattleMode == "5v5" {
			fmt.Print("Enter your move (attack/defend/surrender/sacrifice/pass). To see the card you are battling with use the 'card' command. To end/lose the game use 'surrender all' command. You can use the 'switch' command to switch to different pokemon if the requirements met: ")
		} else {
			fmt.Print("Enter your move (attack/defend/surrender/sacrifice/pass). To see the card you are battling with use the 'card' command. To end/lose the game use 'surrender all' command: ")
		}
		if !scanner.Scan() {
			return "surrender", 0 // treat EOF as surrender
		}
		move := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if move == "card" {
			PrintCard(*playerCard)
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
			defendCost := playerCard.Defense - int(float64(playerCard.Defense) * 0.75)
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
			if float64(playerCard.Stamina) >= 0.5 * float64(playerCard.Speed * 2) {
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
func handleSacrifice(state *GameState, playerCard *pokemon.Card) {
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

// AI move logic
func getAIMove(playerMove string, aiCard *pokemon.Card, state *GameState, aiIdx int) (string, int) {
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
	defendCost := aiCard.Defense - int(float64(aiCard.Defense) * 0.75)
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
		if count == 0 {
			hpCost = 10
		} else if count == 1 {
			hpCost = 15
		} else if count == 2 {
			hpCost = 20
		} else {
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

// Process the result of a turn
func processTurnResult(playerMove, aiMove string, playerMoveIdx, aiMoveIdx int, playerCard, aiCard *pokemon.Card, state *GameState) {
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
		if aiMove == "attack" {
			aiDmg := calculateDamage(aiCard, playerCard, false, aiMoveIdx)
			playerCard.HP -= aiDmg
			aiCard.Stamina -= aiCard.Moves[aiMoveIdx].StaminaCost
			state.LastHpLost = aiDmg
			state.LastStaminaLost = 0
			state.LastDamageDealt = 0
		} else if aiMove == "defend" {
			aiCard.Stamina -= aiDefendCost
			state.LastHpLost = 0
			state.LastStaminaLost = 0
			state.LastDamageDealt = 0
		}
		return
	}
	if aiMove == "pass" {
		// Player does their move, AI does nothing
		if playerMove == "attack" {
			playerDmg := calculateDamage(playerCard, aiCard, false, playerMoveIdx)
			aiCard.HP -= playerDmg
			playerCard.Stamina -= playerCard.Moves[playerMoveIdx].StaminaCost
			state.LastHpLost = 0
			state.LastStaminaLost = playerCard.Moves[playerMoveIdx].StaminaCost
			state.LastDamageDealt = playerDmg
		} else if playerMove == "defend" {
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

// Calculate damage (probabilities shift with attack stat)
func calculateDamage(attacker, defender *pokemon.Card, defenderDefending bool, moveIdx int) int {
	move := attacker.Moves[moveIdx]
	power := move.Power
	attackStat := attacker.Attack


	percent := rollDamagePercent(attackStat)
	baseDmg := int(float64(power) * percent)

	typeMultiplier := TypeMultiplier(move.Type, defender.Types, attacker.Name)
	baseDmg = int(float64(baseDmg) * typeMultiplier)


	if defenderDefending {
		baseDmg = int(float64(baseDmg) * 0.25)
	}
	return baseDmg
}

// Show round summary
func showRoundSummary(state *GameState, playerCard, aiCard *pokemon.Card) {
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

// Prepare for the next round or end the battle
func prepareNextRound(scanner *bufio.Scanner, state *GameState) {
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
	fmt.Printf("\nPrepare for the %s round.\n", ordinal(state.Round))
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

// AI sacrifice logic
func handleSacrificeAI(aiCard *pokemon.Card, state *GameState) {
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

// Helper to get ordinal suffix for round numbers
func ordinal(n int) string {
	if n%100 >= 11 && n%100 <= 13 {
		return fmt.Sprintf("%dth", n)
	}
	switch n % 10 {
	case 1:
		return fmt.Sprintf("%dst", n)
	case 2:
		return fmt.Sprintf("%dnd", n)
	case 3:
		return fmt.Sprintf("%drd", n)
	default:
		return fmt.Sprintf("%dth", n)
	}
}

// Helper function to reset battle state
func resetBattleState(state *GameState) {
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
	state.BattleMode = ""
}
