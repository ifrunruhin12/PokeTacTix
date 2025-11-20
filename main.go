package main

import (
	"bufio"
	"log"
	"os"
	"pokemon-cli/game"
	"pokemon-cli/game/commands"
	"pokemon-cli/game/models"
)

var version string

func main() {
	log.Printf("PokeTacTix CLI %s", version)
	game.PrintWelcome()

	scanner := bufio.NewScanner(os.Stdin)
	state := &models.GameState{}

	for {
		print("\n> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		commands.HandleCommand(input, scanner, state)
	}
}
