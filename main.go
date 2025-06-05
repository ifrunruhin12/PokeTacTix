package main

import (
	"fmt"
	"pokemon-cli/game"
	"pokemon-cli/pokemon"
)

func main() {
	fmt.Print("Enter the name of the Pokémon: ")
	var name string
	fmt.Scan(&name)

	poke, moves, err := pokemon.FetchPokemon(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	card := game.BuildCardFromPokemon(poke, moves)
	game.PrintCard(card)
}
