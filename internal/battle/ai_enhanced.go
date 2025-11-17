package battle

import (
	"math/rand"
	"pokemon-cli/game/core"
	"pokemon-cli/game/utils"
	"pokemon-cli/internal/pokemon"
	"strings"
)

type EnhancedAIDecision struct {
	Move    string
	MoveIdx int
	Score   float64
	Reason  string
}

func GetEnhancedAIMove(bs *BattleState, playerMove string) (string, int) {
	aiCard := bs.GetActiveAICard()
	playerCard := bs.GetActivePlayerCard()

	if aiCard == nil || playerCard == nil {
		return "pass", 0
	}

	aCard := ConvertFromBattleCard(*aiCard)
	pCard := ConvertFromBattleCard(*playerCard)

	maxStamina := aCard.Speed * 2

	hpPercent := float64(aCard.HP) / float64(aCard.HPMax)

	canAttack := false
	attackMoves := []int{}
	for i, move := range aCard.Moves {
		if aCard.Stamina >= move.StaminaCost {
			canAttack = true
			attackMoves = append(attackMoves, i)
		}
	}

	defendCost := core.GetDefendCost(aCard.HPMax)
	canDefend := aCard.Stamina >= defendCost

	sacrificeCount := bs.SacrificeCount[bs.AIActiveIdx]
	canSacrifice := false
	var hpCost int
	switch sacrificeCount {
	case 0:
		hpCost = 10
	case 1:
		hpCost = 15
	case 2:
		hpCost = 20
	default:
		hpCost = 9999
	}
	if float64(aCard.Stamina) < 0.5*float64(maxStamina) && aCard.HP > hpCost && sacrificeCount < 3 {
		canSacrifice = true
	}

	if !canAttack && !canDefend {
		if canSacrifice {
			return "sacrifice", 0
		}

		if bs.Mode == "1v1" {
			if hpPercent < 0.1 && aCard.Stamina == 0 {
				if rand.Float64() < 0.3 { // 30% chance to surrender even in dire situation
					return "surrender", 0
				}
			}
		} else if bs.Mode == "5v5" {
			hasStrongerPokemon := false
			for i, card := range bs.AIDeck {
				if i != bs.AIActiveIdx && card.HP > 0 {
					cardHPPercent := float64(card.HP) / float64(card.HPMax)
					if cardHPPercent > hpPercent+0.3 {
						hasStrongerPokemon = true
						break
					}
				}
			}

			if hpPercent < 0.25 && hasStrongerPokemon && aCard.Stamina < maxStamina/4 {
				if rand.Float64() < 0.4 { // 40% chance for strategic retreat
					return "surrender", 0
				}
			}
		}

		return "pass", 0
	}

	// Evaluate all possible moves
	decisions := []EnhancedAIDecision{}

	// Evaluate attack moves
	if canAttack {
		for _, moveIdx := range attackMoves {
			score := evaluateAttackMove(&aCard, &pCard, moveIdx, playerMove, hpPercent)
			decisions = append(decisions, EnhancedAIDecision{
				Move:    "attack",
				MoveIdx: moveIdx,
				Score:   score,
				Reason:  "attack evaluation",
			})
		}
	}

	// Evaluate defend
	if canDefend {
		score := evaluateDefend(&aCard, &pCard, playerMove, hpPercent)
		decisions = append(decisions, EnhancedAIDecision{
			Move:    "defend",
			MoveIdx: 0,
			Score:   score,
			Reason:  "defend evaluation",
		})
	}

	// Evaluate pass (usually low score)
	decisions = append(decisions, EnhancedAIDecision{
		Move:    "pass",
		MoveIdx: 0,
		Score:   0.1,
		Reason:  "pass evaluation",
	})

	// Find best decision
	bestDecision := decisions[0]
	for _, decision := range decisions {
		if decision.Score > bestDecision.Score {
			bestDecision = decision
		}
	}

	return bestDecision.Move, bestDecision.MoveIdx
}

// evaluateAttackMove scores an attack move based on various factors
func evaluateAttackMove(aiCard, playerCard *pokemon.Card, moveIdx int, playerMove string, hpPercent float64) float64 {
	move := aiCard.Moves[moveIdx]
	score := 0.0

	// Base score from move power
	score += float64(move.Power) / 100.0

	// Type effectiveness (60% weight)
	typeMultiplier := getTypeEffectiveness(move.Type, playerCard.Types)
	if typeMultiplier > 1.0 {
		score += 0.6 * (typeMultiplier - 1.0) // Super effective bonus
	} else if typeMultiplier < 1.0 {
		score -= 0.3 * (1.0 - typeMultiplier) // Not very effective penalty
	}

	if playerMove == "defend" {
		if move.Power >= 80 {
			score += 0.4 // High power bonus when opponent defending
		}
	}

	if playerMove == "attack" {
		score -= 0.2 // Slight penalty for attacking when opponent attacks
	}

	if playerMove == "pass" {
		score += 0.5
	}

	staminaEfficiency := float64(move.Power) / float64(move.StaminaCost)
	score += staminaEfficiency * 0.1

	if hpPercent < 0.3 {
		score += 0.3
	}

	return score
}

func evaluateDefend(aiCard, playerCard *pokemon.Card, playerMove string, hpPercent float64) float64 {
	score := 0.0

	// If player is attacking, defending is valuable
	if playerMove == "attack" {
		score += 0.7

		// Check if we're at type disadvantage
		for _, pMove := range playerCard.Moves {
			typeMultiplier := getTypeEffectiveness(pMove.Type, aiCard.Types)
			if typeMultiplier > 1.0 {
				score += 0.3 // Extra bonus for defending against super effective
				break
			}
		}
	}

	// If player is defending or passing, defending is less valuable
	if playerMove == "defend" || playerMove == "pass" {
		score -= 0.5
	}

	// If low HP, defending is more important
	if hpPercent < 0.3 {
		score += 0.4
	}

	// If high HP, defending is less important
	if hpPercent > 0.7 {
		score -= 0.2
	}

	return score
}

// getTypeEffectiveness calculates type effectiveness multiplier
func getTypeEffectiveness(moveType string, defenderTypes []string) float64 {
	multiplier := 1.0
	moveType = strings.ToLower(moveType)

	if effectiveness, exists := utils.TypeChart[moveType]; exists {
		for _, defenderType := range defenderTypes {
			defenderType = strings.ToLower(defenderType)
			if typeMultiplier, typeExists := effectiveness[defenderType]; typeExists {
				multiplier *= typeMultiplier
			}
		}
	}

	return multiplier
}

// ShouldAISwitch determines if AI should switch Pokemon in 5v5
func ShouldAISwitch(bs *BattleState) (bool, int) {
	if bs.Mode != "5v5" {
		return false, -1
	}

	aiCard := bs.GetActiveAICard()
	if aiCard == nil || aiCard.HP <= 0 {
		// Must switch if knocked out
		for i, card := range bs.AIDeck {
			if card.HP > 0 && i != bs.AIActiveIdx {
				return true, i
			}
		}
		return false, -1
	}

	// Check if current Pokemon is in bad shape
	hpPercent := float64(aiCard.HP) / float64(aiCard.HPMax)
	staminaPercent := float64(aiCard.Stamina) / float64(aiCard.StaminaMax)

	// Switch if HP < 30% or stamina < 30%
	if hpPercent < 0.3 || staminaPercent < 0.3 {
		// Find best alternative
		bestIdx := -1
		bestScore := hpPercent + staminaPercent

		for i, card := range bs.AIDeck {
			if card.HP > 0 && i != bs.AIActiveIdx {
				cardScore := float64(card.HP)/float64(card.HPMax) + float64(card.Stamina)/float64(card.StaminaMax)
				if cardScore > bestScore {
					bestIdx = i
					bestScore = cardScore
				}
			}
		}

		if bestIdx != -1 {
			return true, bestIdx
		}
	}

	return false, -1
}
