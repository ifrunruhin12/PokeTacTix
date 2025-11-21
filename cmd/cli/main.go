//go:build cli
// +build cli

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"pokemon-cli/internal/cli/commands"
	"pokemon-cli/internal/cli/setup"
	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
)

const version = "0.1.0-alpha"

func main() {
	// Check if this is first launch
	isFirst, err := setup.IsFirstLaunch()
	if err != nil {
		log.Fatalf("Error checking first launch: %v", err)
	}

	var gameState *storage.GameState

	if isFirst {
		// Run first-time setup
		gameState, err = setup.RunCompleteSetup()
		if err != nil {
			log.Fatalf("Setup failed: %v", err)
		}
	} else {
		// Load existing game state
		gameState, err = storage.LoadGameState()
		if err != nil {
			log.Fatalf("Failed to load game state: %v", err)
		}

		// Welcome back message
		renderer := ui.NewRenderer()
		renderer.Clear()
		fmt.Println(ui.RenderLogo())
		fmt.Println()
		if ui.GetColorSupport() {
			fmt.Printf("%s\n\n", ui.Colorize(fmt.Sprintf("Welcome back, %s!", gameState.PlayerName), ui.ColorBrightCyan))
		} else {
			fmt.Printf("Welcome back, %s!\n\n", gameState.PlayerName)
		}
	}

	// Display player info
	displayPlayerInfo(gameState)

	// Start command loop
	runCommandLoop(gameState)
}

func displayPlayerInfo(state *storage.GameState) {
	fmt.Println(ui.RenderDivider(75, "═"))
	fmt.Printf("Coins: %d | Pokemon: %d | Deck: %d\n",
		state.Coins,
		len(state.Collection),
		len(state.Deck))
	fmt.Printf("Battles: %d (1v1: %d W/%d L | 5v5: %d W/%d L)\n",
		state.Stats.TotalBattles1v1+state.Stats.TotalBattles5v5,
		state.Stats.Wins1v1, state.Stats.Losses1v1,
		state.Stats.Wins5v5, state.Stats.Losses5v5)
	fmt.Println(ui.RenderDivider(75, "═"))
	fmt.Println()
}

func runCommandLoop(state *storage.GameState) {
	scanner := bufio.NewScanner(os.Stdin)
	renderer := ui.NewRenderer()

	// Create command handler
	cmdHandler := commands.NewCommandHandler(state, renderer, scanner)

	fmt.Println("Type 'help' for available commands or 'quit' to exit.")
	fmt.Println()

	// Show initial hint for new players
	totalBattles := state.Stats.TotalBattles1v1 + state.Stats.TotalBattles5v5
	if totalBattles < 3 {
		cmdHandler.DisplayCommandHints()
		fmt.Println()
	}

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Parse command
		parts := strings.Fields(input)
		command := parts[0]
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}

		// Handle special commands that aren't in the command handler
		if strings.ToLower(command) == "info" || strings.ToLower(command) == "i" {
			displayPlayerInfo(state)
			fmt.Println()
			continue
		}

		if strings.ToLower(command) == "reset" {
			if confirmReset() {
				resetGame()
				return
			}
			fmt.Println()
			continue
		}

		// Use the command handler for all other commands
		err := cmdHandler.HandleCommand(command, args)
		if err != nil {
			// Check if it's a quit signal
			if err.Error() == "QUIT" {
				return
			}

			// Display error if it's not handled by the command handler
			if ui.GetColorSupport() {
				fmt.Printf("%s\n", ui.Colorize(fmt.Sprintf("Error: %v", err), ui.ColorRed))
			} else {
				fmt.Printf("Error: %v\n", err)
			}
		}

		// Reload game state after commands that might modify it
		cmdLower := strings.ToLower(command)
		if cmdLower == "battle" || cmdLower == "b" || 
		   cmdLower == "shop" || cmdLower == "s" ||
		   (cmdLower == "deck" && len(args) > 0 && strings.ToLower(args[0]) == "edit") {
			var reloadErr error
			state, reloadErr = storage.LoadGameState()
			if reloadErr != nil {
				log.Printf("Warning: Failed to reload game state: %v", reloadErr)
			}
		}

		fmt.Println()
	}
}

func confirmReset() bool {
	fmt.Println()
	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("WARNING: This will delete your save file and all progress!", ui.ColorRed))
	} else {
		fmt.Println("WARNING: This will delete your save file and all progress!")
	}
	fmt.Print("Are you sure? Type 'yes' to confirm: ")

	var response string
	fmt.Scanln(&response)

	return strings.ToLower(response) == "yes"
}

func resetGame() {
	savePath, err := storage.GetSaveFilePath()
	if err != nil {
		fmt.Printf("Error getting save path: %v\n", err)
		return
	}

	if err := os.Remove(savePath); err != nil {
		fmt.Printf("Error deleting save file: %v\n", err)
		return
	}

	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("\n✓ Save file deleted. Restart the game to begin fresh!", ui.ColorBrightGreen))
	} else {
		fmt.Println("\n✓ Save file deleted. Restart the game to begin fresh!")
	}
}
