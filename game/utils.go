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

// RollDamagePercent returns a random percent (0.0-1.0) based on attack stat and move power.
func RollDamagePercent(attackStat int, movePower int) float64 {
	// TODO: Implement probability table logic
	return rand.Float64() // placeholder
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
