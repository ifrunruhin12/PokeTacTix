package setup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
)

// IsFirstLaunch checks if this is the first time the game is being launched
// Returns true if no save file exists
func IsFirstLaunch() (bool, error) {
	exists, err := storage.SaveFileExists()
	if err != nil {
		return false, fmt.Errorf("failed to check save file: %w", err)
	}
	return !exists, nil
}

// PromptPlayerName prompts the user to enter their player name with validation
// Returns the validated player name
func PromptPlayerName() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nEnter your player name: ")
		name, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		// Trim whitespace and newline
		name = strings.TrimSpace(name)

		// Validate name
		if len(name) == 0 {
			fmt.Println("Name cannot be empty. Please try again.")
			continue
		}

		if len(name) < 2 {
			fmt.Println("Name must be at least 2 characters long. Please try again.")
			continue
		}

		if len(name) > 20 {
			fmt.Println("Name must be 20 characters or less. Please try again.")
			continue
		}

		// Check for valid characters (alphanumeric, spaces, hyphens, underscores)
		valid := true
		for _, char := range name {
			if !((char >= 'a' && char <= 'z') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') ||
				char == ' ' || char == '-' || char == '_') {
				valid = false
				break
			}
		}

		if !valid {
			fmt.Println("Name can only contain letters, numbers, spaces, hyphens, and underscores.")
			continue
		}

		return name, nil
	}
}

// ShowWelcomeMessage displays the welcome message with ASCII art logo
func ShowWelcomeMessage() {
	// Clear screen
	renderer := ui.NewRenderer()
	renderer.Clear()

	// Display logo
	fmt.Println(ui.RenderLogo())
	fmt.Println()

	// Welcome message
	welcomeText := `
Welcome to PokeTacTix CLI - The Ultimate Pokemon Battle Experience!

This is an offline Pokemon battle game where you can:
  • Collect and battle with Pokemon from Generations 1-5
  • Build your ultimate deck of 5 Pokemon
  • Battle in 1v1 or 5v5 modes against AI opponents
  • Level up your Pokemon and earn coins
  • Purchase new Pokemon from the shop
  • Track your battle statistics and progress

All data is stored locally on your computer, and the game works
completely offline - no internet connection required!
`

	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize(welcomeText, ui.ColorBrightCyan))
	} else {
		fmt.Println(welcomeText)
	}

	fmt.Println(ui.RenderDivider(75, "═"))
	fmt.Println()
}

// ShowTutorial displays a tutorial explaining game basics
func ShowTutorial() {
	fmt.Println()
	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("=== GAME TUTORIAL ===", ui.ColorBrightYellow))
	} else {
		fmt.Println("=== GAME TUTORIAL ===")
	}
	fmt.Println()

	tutorial := []string{
		"BATTLE MODES:",
		"  • 1v1 Mode: Quick battles with one random Pokemon from your deck",
		"  • 5v5 Mode: Strategic battles using all 5 Pokemon in your deck",
		"",
		"BATTLE ACTIONS:",
		"  • Attack: Use one of your Pokemon's 4 moves to deal damage",
		"  • Defend: Reduce incoming damage by 50% for one turn",
		"  • Pass: Skip your turn to regenerate 20% stamina",
		"  • Sacrifice: In 5v5, sacrifice current Pokemon to fully heal the next one",
		"  • Surrender: Give up the battle (you'll still earn some coins)",
		"",
		"PROGRESSION:",
		"  • Win battles to earn coins and XP for your Pokemon",
		"  • Pokemon level up when they gain enough XP, increasing their stats",
		"  • Use coins to purchase new Pokemon from the shop",
		"  • After winning 5v5 battles, select one AI Pokemon to add to your collection",
		"",
		"COMMANDS:",
		"  • battle (b)     - Start a battle",
		"  • collection (c) - View your Pokemon collection",
		"  • deck (d)       - View and edit your battle deck",
		"  • shop (s)       - Visit the shop to buy Pokemon",
		"  • stats (st)     - View your battle statistics",
		"  • help (h)       - Show all available commands",
		"  • quit (q)       - Exit the game",
		"",
		"TYPE EFFECTIVENESS:",
		"  • Fire beats Grass, Grass beats Water, Water beats Fire",
		"  • Electric beats Water and Flying",
		"  • Fighting beats Normal, Rock, and Steel",
		"  • And many more type matchups to discover!",
	}

	for _, line := range tutorial {
		if strings.HasPrefix(line, "  •") {
			if ui.GetColorSupport() {
				fmt.Println(ui.Colorize(line, ui.ColorBrightWhite))
			} else {
				fmt.Println(line)
			}
		} else if strings.HasSuffix(line, ":") {
			if ui.GetColorSupport() {
				fmt.Println(ui.Colorize(line, ui.ColorBrightGreen))
			} else {
				fmt.Println(line)
			}
		} else {
			fmt.Println(line)
		}
	}

	fmt.Println()
	fmt.Println(ui.RenderDivider(75, "─"))
}

// PromptTutorialSkip asks if the user wants to see the tutorial
func PromptTutorialSkip() (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nWould you like to see the tutorial? (y/n): ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}

// RunOnboarding runs the complete onboarding flow for new players
// Returns the player name if successful
func RunOnboarding() (string, error) {
	// Show welcome message
	ShowWelcomeMessage()

	// Prompt for player name
	playerName, err := PromptPlayerName()
	if err != nil {
		return "", err
	}

	// Confirm name
	fmt.Printf("\nWelcome, %s!\n", playerName)

	// Ask about tutorial
	showTutorial, err := PromptTutorialSkip()
	if err != nil {
		return "", err
	}

	if showTutorial {
		ShowTutorial()
	} else {
		fmt.Println("\nYou can always view the tutorial later by typing 'help'.")
	}

	return playerName, nil
}
