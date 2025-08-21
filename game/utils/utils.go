// Package utils provides utility functions for the Pokémon CLI application
package utils

import (
	"fmt"
	"pokemon-cli/pokemon"
)

// FetchRandomDeck returns a slice of 5 random Pokémon cards, with legendary/mythical odds handled in FetchRandomPokemonCard
func FetchRandomDeck() []pokemon.Card {
	deck := make([]pokemon.Card, 0, 5)
	for range [5]int{} {
		card := pokemon.FetchRandomPokemonCard(false)
		deck = append(deck, card)
	}
	return deck
}

// Ordinal suffix for round numbers
func Ordinal(n int) string {
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
