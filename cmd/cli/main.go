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

	fmt.Println("Type 'help' for available commands or 'quit' to exit.")
	fmt.Println()

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
		command := strings.ToLower(parts[0])

		// Handle commands
		switch command {
		case "help", "h":
			showHelp()

		case "info", "i":
			displayPlayerInfo(state)

		case "collection", "c":
			collectionCmd := commands.NewCollectionCommand(state, renderer, scanner)
			if err := collectionCmd.ViewCollection(); err != nil {
				if ui.GetColorSupport() {
					fmt.Printf("%s\n", ui.Colorize(fmt.Sprintf("Collection error: %v", err), ui.ColorRed))
				} else {
					fmt.Printf("Collection error: %v\n", err)
				}
			}

		case "deck", "d":
			handleDeckCommand(state, renderer, scanner, parts)

		case "stats", "st":
			showStats(state)

		case "battle", "b":
			battleCmd := commands.NewBattleCommand(state, renderer, scanner)
			if err := battleCmd.StartBattle(); err != nil {
				if ui.GetColorSupport() {
					fmt.Printf("%s\n", ui.Colorize(fmt.Sprintf("Battle error: %v", err), ui.ColorRed))
				} else {
					fmt.Printf("Battle error: %v\n", err)
				}
			}
			// Reload game state after battle
			var err error
			state, err = storage.LoadGameState()
			if err != nil {
				log.Printf("Warning: Failed to reload game state: %v", err)
			}

		case "quit", "q", "exit":
			fmt.Println("\nThanks for playing PokeTacTix!")
			fmt.Println("Your progress has been saved.")
			return

		case "reset":
			if confirmReset() {
				resetGame()
				return
			}

		default:
			if ui.GetColorSupport() {
				fmt.Printf("%s\n", ui.Colorize(fmt.Sprintf("Unknown command: %s", command), ui.ColorRed))
			} else {
				fmt.Printf("Unknown command: %s\n", command)
			}
			fmt.Println("Type 'help' for available commands.")
		}

		fmt.Println()
	}
}

func showHelp() {
	fmt.Println()
	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("=== AVAILABLE COMMANDS ===", ui.ColorBrightYellow))
	} else {
		fmt.Println("=== AVAILABLE COMMANDS ===")
	}
	fmt.Println()

	commands := []struct {
		name    string
		aliases string
		desc    string
	}{
		{"help", "h", "Show this help message"},
		{"info", "i", "Display player information"},
		{"collection", "c", "View your Pokemon collection"},
		{"deck", "d", "View your battle deck"},
		{"deck edit", "", "Edit your battle deck"},
		{"battle", "b", "Start a battle (1v1 or 5v5)"},
		{"stats", "st", "View your battle statistics"},
		{"quit", "q, exit", "Exit the game"},
		{"reset", "", "Reset game (delete save file)"},
	}

	for _, cmd := range commands {
		if cmd.aliases != "" {
			fmt.Printf("  %-15s (%s) - %s\n", cmd.name, cmd.aliases, cmd.desc)
		} else {
			fmt.Printf("  %-15s - %s\n", cmd.name, cmd.desc)
		}
	}

	fmt.Println()
	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("Note: Shop and other commands are coming soon!", ui.ColorYellow))
	} else {
		fmt.Println("Note: Shop and other commands are coming soon!")
	}
}

func handleDeckCommand(state *storage.GameState, renderer *ui.Renderer, scanner *bufio.Scanner, parts []string) {
	deckCmd := commands.NewDeckCommand(state, renderer, scanner)

	// Check if subcommand is provided
	if len(parts) > 1 {
		subcommand := strings.ToLower(parts[1])
		switch subcommand {
		case "edit", "e":
			if err := deckCmd.EditDeck(); err != nil {
				if ui.GetColorSupport() {
					fmt.Printf("%s\n", ui.Colorize(fmt.Sprintf("Deck edit error: %v", err), ui.ColorRed))
				} else {
					fmt.Printf("Deck edit error: %v\n", err)
				}
			}
			// Reload game state after editing
			var err error
			state, err = storage.LoadGameState()
			if err != nil {
				log.Printf("Warning: Failed to reload game state: %v", err)
			}
		default:
			if ui.GetColorSupport() {
				fmt.Printf("%s\n", ui.Colorize(fmt.Sprintf("Unknown deck subcommand: %s", subcommand), ui.ColorRed))
			} else {
				fmt.Printf("Unknown deck subcommand: %s\n", subcommand)
			}
			fmt.Println("Available subcommands: edit")
		}
	} else {
		// No subcommand, just view the deck
		if err := deckCmd.ViewDeck(); err != nil {
			if ui.GetColorSupport() {
				fmt.Printf("%s\n", ui.Colorize(fmt.Sprintf("Deck view error: %v", err), ui.ColorRed))
			} else {
				fmt.Printf("Deck view error: %v\n", err)
			}
		}
	}
}

func showStats(state *storage.GameState) {
	fmt.Println()
	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("=== BATTLE STATISTICS ===", ui.ColorBrightYellow))
	} else {
		fmt.Println("=== BATTLE STATISTICS ===")
	}
	fmt.Println()

	// 1v1 Stats
	fmt.Println("1v1 Battles:")
	fmt.Printf("  Total: %d\n", state.Stats.TotalBattles1v1)
	fmt.Printf("  Wins: %d\n", state.Stats.Wins1v1)
	fmt.Printf("  Losses: %d\n", state.Stats.Losses1v1)
	fmt.Printf("  Draws: %d\n", state.Stats.Draws1v1)
	if state.Stats.TotalBattles1v1 > 0 {
		winRate := float64(state.Stats.Wins1v1) / float64(state.Stats.TotalBattles1v1) * 100
		fmt.Printf("  Win Rate: %.1f%%\n", winRate)
	}
	fmt.Println()

	// 5v5 Stats
	fmt.Println("5v5 Battles:")
	fmt.Printf("  Total: %d\n", state.Stats.TotalBattles5v5)
	fmt.Printf("  Wins: %d\n", state.Stats.Wins5v5)
	fmt.Printf("  Losses: %d\n", state.Stats.Losses5v5)
	fmt.Printf("  Draws: %d\n", state.Stats.Draws5v5)
	if state.Stats.TotalBattles5v5 > 0 {
		winRate := float64(state.Stats.Wins5v5) / float64(state.Stats.TotalBattles5v5) * 100
		fmt.Printf("  Win Rate: %.1f%%\n", winRate)
	}
	fmt.Println()

	// Overall Stats
	fmt.Println("Overall:")
	fmt.Printf("  Total Pokemon: %d\n", state.Stats.TotalPokemon)
	fmt.Printf("  Highest Level: %d\n", state.Stats.HighestLevel)
	fmt.Printf("  Total Coins Earned: %d\n", state.Stats.TotalCoinsEarned)
	fmt.Printf("  Current Coins: %d\n", state.Coins)
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
