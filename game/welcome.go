package game

import "fmt"

func PrintWelcome() {
	fmt.Println("Welcome to the game PokeTacTix made by Ifran Ruhin aka popcycle!")
	fmt.Println("This is version 1.0.0(alpha) of the game. You can also check the version using the command \"version\"")
	fmt.Println("This is a text-based stategic game where you can play against an AI in a pokemon battle.\nYou can also search for a pokemon name to see its card")
	fmt.Println("Enter a command to explore and do stuff (you can type \"command\" to see all the commands available)")
}

func GetWelcomeMessage() string {
	return `🔥PokeTacTix is a game made by Ifran Ruhin aka popcycle!
This is version 1.0.0(alpha) of the game. 
This is a text-based strategic game where you can play against an AI in a Pokémon battle.
You can also search for a Pokémon name to see its card by clicking the "search" button`
}
