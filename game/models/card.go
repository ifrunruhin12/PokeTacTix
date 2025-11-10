// This file is part of the pokemon-cli project
// It defines the PrintCard function to display a Pokémon card's details.
// It is used to print the card information in a readable format.

package models

import (
	"fmt"
	"pokemon-cli/internal/pokemon"
)

func PrintCard(card pokemon.Card) {
	fmt.Println("====================")
	fmt.Println("Name:", card.Name)
	fmt.Println("HP:", card.HP)
	fmt.Println("Stamina:", card.Stamina)
	fmt.Println("Attack:", card.Attack)
	fmt.Println("Defense:", card.Defense)
	fmt.Println("Types:", card.Types)
	fmt.Println("Sprite URL:", card.Sprite)
	fmt.Println("Moves:")
	for i, move := range card.Moves {
		fmt.Printf("  %d. %s (Power: %d, Stamina Cost: %d, Type: %s)\n", i+1, move.Name, move.Power, move.StaminaCost, move.Type)
	}
	fmt.Println("====================")
}
