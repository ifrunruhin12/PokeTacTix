package main

import (
	"fmt"
	"pokemon-cli/pokemon"
)

func main() {
	fmt.Print("Enter the name of the Pokémon: ")
	var name string
	fmt.Scan(&name)

	poke, err := pokemon.FetchPokemon(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	pokemon.DisplayPokemon(poke)
}

