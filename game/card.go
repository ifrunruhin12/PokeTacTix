package game

import (
	"fmt"
	"pokemon-cli/pokemon"
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
