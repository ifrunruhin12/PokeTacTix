package core

import (
	"math/rand"
	"pokemon-cli/game/utils"
	"pokemon-cli/internal/pokemon"
	"strings"
)

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

// Roll damage percent based on attack stat
func rollDamagePercent(attackStat int) float64 {
	// Three tables: low (<=30), high (70), super (>=120)
	low := []struct{ pct, prob float64 }{
		{0.10, 0.07}, {0.20, 0.13}, {0.30, 0.35}, {0.40, 0.25}, {0.60, 0.10}, {0.80, 0.07}, {1.00, 0.03},
	}
	high := []struct{ pct, prob float64 }{
		{0.10, 0.01}, {0.20, 0.04}, {0.30, 0.10}, {0.40, 0.15}, {0.60, 0.25}, {0.80, 0.25}, {1.00, 0.20},
	}
	super := []struct{ pct, prob float64 }{
		{0.10, 0.00}, {0.20, 0.01}, {0.30, 0.04}, {0.40, 0.10}, {0.60, 0.15}, {0.80, 0.30}, {1.00, 0.40},
	}

	var table []struct{ pct, prob float64 }
	if attackStat <= 30 {
		table = low
	} else if attackStat >= 120 {
		table = super
	} else if attackStat >= 70 {
		// Interpolate between high and super
		frac := float64(attackStat-70) / 50.0
		table = make([]struct{ pct, prob float64 }, len(high))
		for i := range high {
			table[i].pct = high[i].pct
			table[i].prob = high[i].prob*(1-frac) + super[i].prob*frac
		}
	} else {
		// Interpolate between low and high
		frac := float64(attackStat-30) / 40.0
		table = make([]struct{ pct, prob float64 }, len(low))
		for i := range low {
			table[i].pct = low[i].pct
			table[i].prob = low[i].prob*(1-frac) + high[i].prob*frac
		}
	}
	// Roll
	roll := rand.Float64()
	cum := 0.0
	for _, entry := range table {
		cum += entry.prob
		if roll < cum {
			return entry.pct
		}
	}
	return 0.10 // fallback
}

// TypeMultiplier returns the type effectiveness multiplier for a move type against defender types.
// Also applies 2x multiplier for legendary/mythical Pokémon attacks.
func TypeMultiplier(moveType string, defenderTypes []string, attackerName string) float64 {
	multiplier := 1.0

	// Apply type effectiveness
	moveType = strings.ToLower(moveType)
	if effectiveness, exists := utils.TypeChart[moveType]; exists {
		for _, defenderType := range defenderTypes {
			defenderType = strings.ToLower(defenderType)
			if typeMultiplier, typeExists := effectiveness[defenderType]; typeExists {
				multiplier *= typeMultiplier
			}
		}
	}

	// Apply legendary/mythical bonus (2x damage for their attacks)
	if utils.IsLegendaryOrMythical(attackerName) {
		multiplier *= 2.0
	}

	return multiplier
}
