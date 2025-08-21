package main

import (
	"bufio"
	"os"
	"pokemon-cli/game"
)

func main() {
	game.PrintWelcome()

	scanner := bufio.NewScanner(os.Stdin)
	state := &game.GameState{}

	for {
		print("\n> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		game.HandleCommand(input, scanner, state)
	}
}
