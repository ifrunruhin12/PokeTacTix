package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"pokemon-cli/internal/cli/commands"
	"pokemon-cli/internal/cli/setup"
	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
)

var (
	Version    = "dev"
	BuildDate  = "unknown"
	CommitHash = "unknown"
	GoVersion  = runtime.Version()
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v" || os.Args[1] == "version") {
		displayVersion()
		return
	}
	isFirst, err := setup.IsFirstLaunch()
	if err != nil {
		log.Fatalf("Error checking first launch: %v", err)
	}

	var gameState *storage.GameState

	if isFirst {
		gameState, err = setup.RunCompleteSetup()
		if err != nil {
			log.Fatalf("Setup failed: %v", err)
		}
	} else {
		gameState, err = storage.LoadGameState()
		if err != nil {
			log.Fatalf("Failed to load game state: %v", err)
		}

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

	displayPlayerInfo(gameState)

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

	cmdHandler := commands.NewCommandHandler(state, renderer, scanner)

	fmt.Println("Type 'help' for available commands or 'quit' to exit.")
	fmt.Println()

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

		parts := strings.Fields(input)
		command := parts[0]
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}

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

		err := cmdHandler.HandleCommand(command, args)
		if err != nil {
			if err.Error() == "QUIT" {
				return
			}

			if ui.GetColorSupport() {
				fmt.Printf("%s\n", ui.Colorize(fmt.Sprintf("Error: %v", err), ui.ColorRed))
			} else {
				fmt.Printf("Error: %v\n", err)
			}
		}

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

func displayVersion() {
	fmt.Println()
	fmt.Println(ui.RenderLogo())
	fmt.Println()
	fmt.Printf("Version:      %s\n", Version)
	fmt.Printf("Build Date:   %s\n", BuildDate)
	fmt.Printf("Commit:       %s\n", CommitHash)
	fmt.Printf("Go Version:   %s\n", GoVersion)
	fmt.Printf("Platform:     %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()
}
