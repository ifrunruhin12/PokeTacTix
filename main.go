package main

import (
	"bufio"
	"os"
	"pokemon-cli/game"
	"pokemon-cli/game/commands"
	"pokemon-cli/game/models"
)

func main() {
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
