package battle

import (
	"fmt"
	"pokemon-cli/game/core"
	"pokemon-cli/internal/pokemon"
	"time"

	"github.com/google/uuid"
)

// StartBattle initializes a new battle with the given mode and decks
func StartBattle(userID int, mode string, playerDeck []pokemon.Card, aiDeck []pokemon.Card) (*BattleState, error) {
	if mode != "1v1" && mode != "5v5" {
		return nil, fmt.Errorf("invalid battle mode: %s", mode)
	}

	if mode == "1v1" && (len(playerDeck) < 1 || len(aiDeck) < 1) {
		return nil, fmt.Errorf("1v1 mode requires at least 1 Pokemon per side")
	}

	if mode == "5v5" && (len(playerDeck) != 5 || len(aiDeck) != 5) {
		return nil, fmt.Errorf("5v5 mode requires exactly 5 Pokemon per side")
	}

	// Convert decks to BattleCards
	playerBattleDeck := make([]BattleCard, len(playerDeck))
	for i, card := range playerDeck {
		playerBattleDeck[i] = ConvertToBattleCard(card, i+1)
	}

	aiBattleDeck := make([]BattleCard, len(aiDeck))
	for i, card := range aiDeck {
		aiBattleDeck[i] = ConvertToBattleCard(card, i+1)
	}

	now := time.Now()
	battleState := &BattleState{
		ID:              uuid.New().String(),
		UserID:          userID,
		Mode:            mode,
		PlayerDeck:      playerBattleDeck,
		AIDeck:          aiBattleDeck,
		PlayerActiveIdx: 0,
		AIActiveIdx:     0,
		TurnNumber:      1,
		RoundNumber:     1,
		WhoseTurn:       "player", // Player always goes first in turn 1
		BattleOver:      false,
		Winner:          "",
		SacrificeCount:  make(map[int]int),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	return battleState, nil
}

// ProcessMove processes a player's move in the battle
func ProcessMove(bs *BattleState, move string, moveIdx *int) ([]string, error) {
	if bs.BattleOver {
		return nil, fmt.Errorf("battle is already over")
	}

	logEntries := []string{}

	// Get active cards first
	playerCard := bs.GetActivePlayerCard()
	aiCard := bs.GetActiveAICard()

	// Handle surrender
	if move == "surrender" {
		if bs.Mode == "1v1" {
			// In 1v1, surrender ends the entire battle
			bs.BattleOver = true
			bs.Winner = "ai"
			logEntries = append(logEntries, "Player surrendered! AI wins the battle!")
			return logEntries, nil
		} else if bs.Mode == "5v5" {
			// In 5v5, surrender only knocks out current Pokemon
			if playerCard != nil {
				playerCard.HP = 0
				playerCard.IsKnockedOut = true
				logEntries = append(logEntries, fmt.Sprintf("Player surrendered! %s was knocked out!", playerCard.Name))
			}

			// Check if player has any Pokemon left
			if !bs.HasPlayerPokemonAlive() {
				bs.BattleOver = true
				bs.Winner = "ai"
				logEntries = append(logEntries, "Player has no Pokemon left! AI wins the battle!")
			} else {
				// Player must switch to next available Pokemon
				for i, card := range bs.PlayerDeck {
					if card.HP > 0 && i != bs.PlayerActiveIdx {
						bs.PlayerActiveIdx = i
						bs.RoundNumber++
						logEntries = append(logEntries, fmt.Sprintf("Player must switch to another Pokemon. Round %d begins.", bs.RoundNumber))
						break
					}
				}
			}
			return logEntries, nil
		}
	}

	if playerCard == nil || aiCard == nil {
		return nil, fmt.Errorf("invalid active Pokemon")
	}

	// Handle sacrifice (free action)
	if move == "sacrifice" {
		oldHP := playerCard.HP
		oldStamina := playerCard.Stamina

		// Convert to pokemon.Card for core logic
		pCard := ConvertFromBattleCard(*playerCard)

		// Use existing sacrifice logic
		sacrificeCount := bs.SacrificeCount[bs.PlayerActiveIdx]
		if sacrificeCount >= 3 {
			return nil, fmt.Errorf("maximum sacrifices reached for this Pokemon")
		}

		var hpCost int
		var staminaGain float64
		switch sacrificeCount {
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

		if pCard.HP <= hpCost {
			return nil, fmt.Errorf("insufficient HP to sacrifice")
		}

		maxStamina := pCard.Speed * 2
		if float64(pCard.Stamina) >= 0.5*float64(maxStamina) {
			return nil, fmt.Errorf("stamina is already above 50%%")
		}

		pCard.HP -= hpCost
		gain := int(float64(maxStamina) * staminaGain)
		pCard.Stamina += gain
		if pCard.Stamina > maxStamina {
			pCard.Stamina = maxStamina
		}

		// Update battle card
		playerCard.HP = pCard.HP
		playerCard.Stamina = pCard.Stamina
		bs.SacrificeCount[bs.PlayerActiveIdx] = sacrificeCount + 1

		hpLost := oldHP - playerCard.HP
		staminaGained := playerCard.Stamina - oldStamina
		logEntries = append(logEntries, fmt.Sprintf("Player sacrificed %d HP and gained %d stamina.", hpLost, staminaGained))

		return logEntries, nil
	}

	// Check whose turn it is
	if bs.WhoseTurn != "player" {
		return nil, fmt.Errorf("it's not player's turn")
	}

	// Validate move
	switch move {
	case "attack":
		if moveIdx == nil {
			return nil, fmt.Errorf("move index required for attack")
		}
		if *moveIdx < 0 || *moveIdx >= len(playerCard.Moves) {
			return nil, fmt.Errorf("invalid move index")
		}
		if playerCard.Stamina < playerCard.Moves[*moveIdx].StaminaCost {
			return nil, fmt.Errorf("insufficient stamina for this move")
		}
	case "defend":
		defendCost := core.GetDefendCost(playerCard.HPMax)
		if playerCard.Stamina < defendCost {
			return nil, fmt.Errorf("insufficient stamina to defend")
		}
	}

	// Store player's move
	bs.PendingPlayerMove = move
	if moveIdx != nil {
		bs.PendingPlayerMoveIdx = *moveIdx
	}

	// Log player's move
	if move == "attack" && moveIdx != nil {
		moveName := playerCard.Moves[*moveIdx].Name
		logEntries = append(logEntries, fmt.Sprintf("Player chose to attack with %s.", moveName))
	} else {
		logEntries = append(logEntries, fmt.Sprintf("Player chose %s.", move))
	}

	// AI makes its move
	aiLogs := processAIMove(bs)
	logEntries = append(logEntries, aiLogs...)

	// Resolve the turn
	resolveLogs := resolveTurn(bs)
	logEntries = append(logEntries, resolveLogs...)

	// Check for knockouts and handle Pokemon switching
	switchLogs := handleKnockouts(bs)
	logEntries = append(logEntries, switchLogs...)

	// Check if battle is over
	bs.CheckBattleEnd()

	// Prepare for next turn
	if !bs.BattleOver {
		bs.TurnNumber++
		bs.WhoseTurn = "player"
		bs.PendingPlayerMove = ""
		bs.PendingPlayerMoveIdx = 0
		bs.PendingAIMove = ""
		bs.PendingAIMoveIdx = 0
	}

	return logEntries, nil
}

// processAIMove handles AI decision making with enhanced logic
func processAIMove(bs *BattleState) []string {
	logEntries := []string{}
	aiCard := bs.GetActiveAICard()

	if aiCard == nil {
		return logEntries
	}

	// Convert to pokemon.Card for core AI logic
	aCard := ConvertFromBattleCard(*aiCard)

	// Handle AI sacrifices
	for {
		// Use enhanced AI decision making
		aiMove, aiMoveIdx := GetEnhancedAIMove(bs, bs.PendingPlayerMove)

		if aiMove == "surrender" {
			if bs.Mode == "1v1" {
				// In 1v1, AI surrender ends the entire battle
				bs.BattleOver = true
				bs.Winner = "player"
				logEntries = append(logEntries, "AI surrendered! Player wins the battle!")
				return logEntries
			} else if bs.Mode == "5v5" {
				// In 5v5, AI surrender only knocks out current Pokemon
				aiCard.HP = 0
				aiCard.IsKnockedOut = true
				logEntries = append(logEntries, fmt.Sprintf("AI surrendered! %s was knocked out!", aiCard.Name))

				// Check if AI has any Pokemon left
				if !bs.HasAIPokemonAlive() {
					bs.BattleOver = true
					bs.Winner = "player"
					logEntries = append(logEntries, "AI has no Pokemon left! Player wins the battle!")
				} else {
					// AI switches to next available Pokemon
					for i, card := range bs.AIDeck {
						if card.HP > 0 && i != bs.AIActiveIdx {
							bs.AIActiveIdx = i
							logEntries = append(logEntries, fmt.Sprintf("AI switched to %s.", bs.AIDeck[i].Name))
							break
						}
					}
				}
				return logEntries
			}
		}

		if aiMove == "sacrifice" {
			maxStamina := aCard.Speed * 2
			if float64(aCard.Stamina) >= 0.5*float64(maxStamina) {
				break
			}

			sacrificeCount := bs.SacrificeCount[bs.AIActiveIdx]
			if sacrificeCount >= 3 {
				break
			}

			oldHP := aCard.HP
			oldStamina := aCard.Stamina

			var hpCost int
			var staminaGain float64
			switch sacrificeCount {
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

			if aCard.HP <= hpCost {
				break
			}

			aCard.HP -= hpCost
			gain := int(float64(maxStamina) * staminaGain)
			aCard.Stamina += gain
			if aCard.Stamina > maxStamina {
				aCard.Stamina = maxStamina
			}

			// Update battle card
			aiCard.HP = aCard.HP
			aiCard.Stamina = aCard.Stamina
			bs.SacrificeCount[bs.AIActiveIdx] = sacrificeCount + 1

			hpLost := oldHP - aCard.HP
			staminaGained := aCard.Stamina - oldStamina
			logEntries = append(logEntries, fmt.Sprintf("AI sacrificed %d HP and gained %d stamina.", hpLost, staminaGained))
			continue
		}

		// Store AI's move
		bs.PendingAIMove = aiMove
		bs.PendingAIMoveIdx = aiMoveIdx

		// Log AI's move
		if aiMove == "attack" {
			moveName := aCard.Moves[aiMoveIdx].Name
			logEntries = append(logEntries, fmt.Sprintf("AI chose to attack with %s.", moveName))
		} else {
			logEntries = append(logEntries, fmt.Sprintf("AI chose %s.", aiMove))
		}

		break
	}

	return logEntries
}

// resolveTurn resolves both player and AI moves
func resolveTurn(bs *BattleState) []string {
	logEntries := []string{}

	playerCard := bs.GetActivePlayerCard()
	aiCard := bs.GetActiveAICard()

	if playerCard == nil || aiCard == nil {
		return logEntries
	}

	// Convert to pokemon.Card for core logic
	pCard := ConvertFromBattleCard(*playerCard)
	aCard := ConvertFromBattleCard(*aiCard)

	playerMove := bs.PendingPlayerMove
	aiMove := bs.PendingAIMove
	playerMoveIdx := bs.PendingPlayerMoveIdx
	aiMoveIdx := bs.PendingAIMoveIdx

	playerDefendCost := core.GetDefendCost(pCard.HPMax)
	aiDefendCost := core.GetDefendCost(aCard.HPMax)

	// Track damage for logging
	playerDamage := 0
	aiDamage := 0

	// Process moves based on combination
	if playerMove == "pass" && aiMove == "pass" {
		bs.ConsecutivePasses++
		logEntries = append(logEntries, fmt.Sprintf("Both passed. Nothing happened! (Pass count: %d/3)", bs.ConsecutivePasses))

		// Check for stalemate - if both players pass 3 times in a row, end in draw
		if bs.ConsecutivePasses >= 3 {
			bs.BattleOver = true
			bs.Winner = "draw"
			logEntries = append(logEntries, "Stalemate! Both players passed 3 times in a row. Battle ends in a draw!")
		}
	} else if playerMove == "pass" {
		// Reset consecutive passes if only one player passed
		bs.ConsecutivePasses = 0
		switch aiMove {
		case "attack":
			aiDamage = core.CalculateDamage(&aCard, &pCard, false, aiMoveIdx)
			pCard.HP -= aiDamage
			aCard.Stamina -= aCard.Moves[aiMoveIdx].StaminaCost
			logEntries = append(logEntries, fmt.Sprintf("AI dealt %d damage to Player.", aiDamage))
		case "defend":
			aCard.Stamina -= aiDefendCost
		}
	} else if aiMove == "pass" {
		// Reset consecutive passes if only one player passed
		bs.ConsecutivePasses = 0
		switch playerMove {
		case "attack":
			playerDamage = core.CalculateDamage(&pCard, &aCard, false, playerMoveIdx)
			aCard.HP -= playerDamage
			pCard.Stamina -= pCard.Moves[playerMoveIdx].StaminaCost
			logEntries = append(logEntries, fmt.Sprintf("Player dealt %d damage to AI.", playerDamage))
		case "defend":
			pCard.Stamina -= playerDefendCost
		}
	} else if playerMove == "attack" && aiMove == "attack" {
		// Reset consecutive passes when both attack
		bs.ConsecutivePasses = 0
		playerDamage = core.CalculateDamage(&pCard, &aCard, false, playerMoveIdx)
		aiDamage = core.CalculateDamage(&aCard, &pCard, false, aiMoveIdx)
		aCard.HP -= playerDamage
		pCard.HP -= aiDamage
		pCard.Stamina -= pCard.Moves[playerMoveIdx].StaminaCost
		aCard.Stamina -= aCard.Moves[aiMoveIdx].StaminaCost
		logEntries = append(logEntries, fmt.Sprintf("Player dealt %d damage to AI.", playerDamage))
		logEntries = append(logEntries, fmt.Sprintf("AI dealt %d damage to Player.", aiDamage))
	} else if playerMove == "attack" && aiMove == "defend" {
		// Reset consecutive passes
		bs.ConsecutivePasses = 0
		playerDamage = core.CalculateDamage(&pCard, &aCard, true, playerMoveIdx)
		aCard.Stamina -= aiDefendCost
		pCard.Stamina -= pCard.Moves[playerMoveIdx].StaminaCost
		if playerDamage <= aCard.Defense {
			logEntries = append(logEntries, "AI blocked all damage!")
		} else {
			actualDamage := playerDamage - aCard.Defense
			aCard.HP -= actualDamage
			logEntries = append(logEntries, fmt.Sprintf("Player dealt %d damage to AI (after defense).", actualDamage))
		}
	} else if playerMove == "defend" && aiMove == "attack" {
		// Reset consecutive passes
		bs.ConsecutivePasses = 0
		aiDamage = core.CalculateDamage(&aCard, &pCard, true, aiMoveIdx)
		pCard.Stamina -= playerDefendCost
		aCard.Stamina -= aCard.Moves[aiMoveIdx].StaminaCost
		if aiDamage <= pCard.Defense {
			logEntries = append(logEntries, "Player blocked all damage!")
		} else {
			actualDamage := aiDamage - pCard.Defense
			pCard.HP -= actualDamage
			logEntries = append(logEntries, fmt.Sprintf("AI dealt %d damage to Player (after defense).", actualDamage))
		}
	} else if playerMove == "defend" && aiMove == "defend" {
		// Reset consecutive passes
		bs.ConsecutivePasses = 0
		pCard.Stamina -= playerDefendCost
		aCard.Stamina -= aiDefendCost
		logEntries = append(logEntries, "Both defended. No damage dealt.")
	}

	// Clamp HP and stamina
	if pCard.HP < 0 {
		pCard.HP = 0
	}
	if aCard.HP < 0 {
		aCard.HP = 0
	}
	if pCard.Stamina < 0 {
		pCard.Stamina = 0
	}
	if aCard.Stamina < 0 {
		aCard.Stamina = 0
	}

	// Update battle cards
	playerCard.HP = pCard.HP
	playerCard.Stamina = pCard.Stamina
	playerCard.IsKnockedOut = pCard.HP <= 0
	aiCard.HP = aCard.HP
	aiCard.Stamina = aCard.Stamina
	aiCard.IsKnockedOut = aCard.HP <= 0

	return logEntries
}

// handleKnockouts handles Pokemon knockouts and switching logic
func handleKnockouts(bs *BattleState) []string {
	logEntries := []string{}

	playerCard := bs.GetActivePlayerCard()
	aiCard := bs.GetActiveAICard()

	if playerCard == nil || aiCard == nil {
		return logEntries
	}

	// Check for knockouts
	playerKO := playerCard.HP <= 0
	aiKO := aiCard.HP <= 0

	if playerKO {
		logEntries = append(logEntries, fmt.Sprintf("Player's %s was knocked out!", playerCard.Name))
		playerCard.IsKnockedOut = true
	}

	if aiKO {
		logEntries = append(logEntries, fmt.Sprintf("AI's %s was knocked out!", aiCard.Name))
		aiCard.IsKnockedOut = true
	}

	// For 1v1, battle ends on knockout
	if bs.Mode == "1v1" {
		return logEntries
	}

	// For 5v5, handle Pokemon switching
	if playerKO && bs.HasPlayerPokemonAlive() {
		// Find next available Pokemon
		for i, card := range bs.PlayerDeck {
			if card.HP > 0 && i != bs.PlayerActiveIdx {
				bs.PlayerActiveIdx = i
				bs.RoundNumber++
				logEntries = append(logEntries, fmt.Sprintf("Player must switch to another Pokemon. Round %d begins.", bs.RoundNumber))
				break
			}
		}
	}

	if aiKO && bs.HasAIPokemonAlive() {
		// AI switches to next available Pokemon
		for i, card := range bs.AIDeck {
			if card.HP > 0 && i != bs.AIActiveIdx {
				bs.AIActiveIdx = i
				logEntries = append(logEntries, fmt.Sprintf("AI switched to %s.", bs.AIDeck[i].Name))
				break
			}
		}
	}

	return logEntries
}

// SwitchPokemon allows the player to switch to a different Pokemon
func SwitchPokemon(bs *BattleState, newIdx int) error {
	if bs.BattleOver {
		return fmt.Errorf("battle is already over")
	}

	if bs.Mode != "5v5" {
		return fmt.Errorf("switching is only allowed in 5v5 battles")
	}

	if newIdx < 0 || newIdx >= len(bs.PlayerDeck) {
		return fmt.Errorf("invalid Pokemon index")
	}

	if newIdx == bs.PlayerActiveIdx {
		return fmt.Errorf("pokemon is already active")
	}

	if bs.PlayerDeck[newIdx].HP <= 0 {
		return fmt.Errorf("cannot switch to a knocked out Pokemon")
	}

	bs.PlayerActiveIdx = newIdx
	bs.RoundNumber++

	return nil
}
