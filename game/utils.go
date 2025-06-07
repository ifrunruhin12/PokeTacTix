package game

import (
	"math/rand"
	"pokemon-cli/pokemon"
)

// TypeMultiplier returns the type effectiveness multiplier for a move type against defender types.
func TypeMultiplier(moveType string, defenderTypes []string) float64 {
	// TODO: Implement type chart logic
	return 1.0 // placeholder
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

// FetchRandomDeck returns a slice of 5 random Pokémon cards, with legendary/mythical odds handled in FetchRandomPokemonCard
func FetchRandomDeck() []pokemon.Card {
	deck := make([]pokemon.Card, 0, 5)
	for range [5]int{} {
		card := pokemon.FetchRandomPokemonCard(false)
		deck = append(deck, card)
	}
	return deck
}
