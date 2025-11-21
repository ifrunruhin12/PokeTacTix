package commands

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
)

// SettingsCommand handles settings-related commands
type SettingsCommand struct {
	gameState *storage.GameState
	renderer  *ui.Renderer
	scanner   *bufio.Scanner
}

// NewSettingsCommand creates a new settings command handler
func NewSettingsCommand(gameState *storage.GameState, renderer *ui.Renderer, scanner *bufio.Scanner) *SettingsCommand {
	return &SettingsCommand{
		gameState: gameState,
		renderer:  renderer,
		scanner:   scanner,
	}
}

// ViewSettings displays the settings menu
func (sc *SettingsCommand) ViewSettings() error {
	for {
		// Clear screen
		sc.renderer.Clear()

		// Display header
		fmt.Println(ui.RenderLogo())
		fmt.Println()
		fmt.Println(strings.Repeat("═", 80))
		fmt.Println(ui.Colorize("GAME SETTINGS", ui.Bold+ui.ColorBrightCyan))
		fmt.Println(strings.Repeat("═", 80))
		fmt.Println()

		// Display current settings
		fmt.Println("Current Settings:")
		fmt.Println()

		// Quick Battle Mode
		quickBattleStatus := "OFF"
		if sc.gameState.Settings.QuickBattle {
			quickBattleStatus = ui.Colorize("ON", ui.ColorGreen)
		} else {
			quickBattleStatus = ui.Colorize("OFF", ui.ColorRed)
		}
		fmt.Printf("  Quick Battle Mode: %s\n", quickBattleStatus)
		fmt.Println("    Skip animations and reduce delays for faster battles")
		fmt.Println()

		// Battle Speed
		speedDisplay := strings.ToUpper(sc.gameState.Settings.BattleSpeed)
		if sc.gameState.Settings.BattleSpeed == "" {
			speedDisplay = "NORMAL"
			sc.gameState.Settings.BattleSpeed = "normal"
		}
		
		speedColor := ui.ColorYellow
		switch sc.gameState.Settings.BattleSpeed {
		case "slow":
			speedColor = ui.ColorCyan
		case "fast":
			speedColor = ui.ColorRed
		}
		
		fmt.Printf("  Battle Speed: %s\n", ui.Colorize(speedDisplay, speedColor))
		fmt.Println("    Adjust animation and text display speed")
		fmt.Println()

		fmt.Println(strings.Repeat("─", 80))
		fmt.Println()

		// Display menu options
		options := []ui.MenuOption{
			{Label: "Toggle Quick Battle", Description: "Enable/disable quick battle mode", Value: "quick"},
			{Label: "Change Battle Speed", Description: "Set battle speed (slow/normal/fast)", Value: "speed"},
			{Label: "Export Save", Description: "Export save file to a location", Value: "export"},
			{Label: "Import Save", Description: "Import save file from a location", Value: "import"},
			{Label: "Save & Exit", Description: "Save settings and return to menu", Value: "save"},
			{Label: "Cancel", Description: "Discard changes and exit", Value: "cancel"},
		}

		fmt.Println(sc.renderer.RenderBorderedMenu(options, -1, "SETTINGS OPTIONS"))
		fmt.Print("Enter your choice (1-6): ")

		// Get user input
		if !sc.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(sc.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 6 {
			fmt.Println(ui.Colorize("Invalid choice. Press Enter to continue...", ui.ColorRed))
			sc.scanner.Scan()
			continue
		}

		switch choice {
		case 1:
			// Toggle quick battle
			sc.gameState.Settings.QuickBattle = !sc.gameState.Settings.QuickBattle
			status := "disabled"
			if sc.gameState.Settings.QuickBattle {
				status = "enabled"
			}
			fmt.Println()
			fmt.Println(ui.Colorize(fmt.Sprintf("Quick Battle Mode %s!", status), ui.ColorGreen))
			fmt.Println("Press Enter to continue...")
			sc.scanner.Scan()

		case 2:
			// Change battle speed
			err := sc.changeBattleSpeed()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				fmt.Println("Press Enter to continue...")
				sc.scanner.Scan()
			}

		case 3:
			// Export save
			err := sc.exportSave()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				fmt.Println("Press Enter to continue...")
				sc.scanner.Scan()
			}

		case 4:
			// Import save
			err := sc.importSave()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				fmt.Println("Press Enter to continue...")
				sc.scanner.Scan()
			}

		case 5:
			// Save and exit
			err := storage.SaveGameState(sc.gameState)
			if err != nil {
				fmt.Println()
				fmt.Println(ui.Colorize(fmt.Sprintf("✗ Failed to save settings: %v", err), ui.ColorRed))
				fmt.Println("Press Enter to continue...")
				sc.scanner.Scan()
				continue
			}

			fmt.Println()
			fmt.Println(ui.Colorize("✓ Settings saved!", ui.ColorGreen))
			fmt.Println("Press Enter to continue...")
			sc.scanner.Scan()
			return nil

		case 6:
			// Cancel
			fmt.Println()
			fmt.Println(ui.Colorize("Settings changes discarded.", ui.ColorYellow))
			fmt.Println("Press Enter to continue...")
			sc.scanner.Scan()
			return nil
		}
	}
}

// changeBattleSpeed allows the user to change battle speed
func (sc *SettingsCommand) changeBattleSpeed() error {
	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println(ui.Colorize("BATTLE SPEED SETTINGS", ui.Bold+ui.ColorBrightYellow))
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()

	speedOptions := []ui.MenuOption{
		{Label: "Slow", Description: "Longer delays, more time to read", Value: "slow"},
		{Label: "Normal", Description: "Default speed (recommended)", Value: "normal"},
		{Label: "Fast", Description: "Minimal delays, quick battles", Value: "fast"},
		{Label: "Cancel", Description: "Keep current speed", Value: "cancel"},
	}

	fmt.Println(sc.renderer.RenderBorderedMenu(speedOptions, -1, "SELECT BATTLE SPEED"))
	fmt.Print("Enter your choice (1-4): ")

	if !sc.scanner.Scan() {
		return fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(sc.scanner.Text())
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > 4 {
		fmt.Println(ui.Colorize("Invalid choice.", ui.ColorRed))
		fmt.Println("Press Enter to continue...")
		sc.scanner.Scan()
		return nil
	}

	if choice == 4 {
		return nil
	}

	speeds := []string{"slow", "normal", "fast"}
	sc.gameState.Settings.BattleSpeed = speeds[choice-1]

	fmt.Println()
	fmt.Println(ui.Colorize(fmt.Sprintf("Battle speed set to %s!", strings.ToUpper(speeds[choice-1])), ui.ColorGreen))
	fmt.Println("Press Enter to continue...")
	sc.scanner.Scan()

	return nil
}


// exportSave exports the save file to a specified location
func (sc *SettingsCommand) exportSave() error {
	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println(ui.Colorize("EXPORT SAVE FILE", ui.Bold+ui.ColorBrightYellow))
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()

	fmt.Println("Enter the destination path (or press Enter for current directory):")
	fmt.Print("> ")

	if !sc.scanner.Scan() {
		return fmt.Errorf("failed to read input")
	}

	destPath := strings.TrimSpace(sc.scanner.Text())
	if destPath == "" {
		// Use current directory
		destPath = "."
	}

	// Export the save file
	exportedPath, err := storage.ExportSave(destPath)
	if err != nil {
		fmt.Println()
		fmt.Println(ui.Colorize(fmt.Sprintf("✗ Export failed: %v", err), ui.ColorRed))
		fmt.Println("Press Enter to continue...")
		sc.scanner.Scan()
		return err
	}

	// Display success message
	fmt.Println()
	fmt.Println(ui.Colorize("═══════════════════════════════════════", ui.ColorGreen))
	fmt.Println(ui.Colorize("  EXPORT SUCCESSFUL!", ui.Bold+ui.ColorGreen))
	fmt.Println(ui.Colorize("═══════════════════════════════════════", ui.ColorGreen))
	fmt.Println()
	fmt.Printf("Save file exported to:\n%s\n", ui.Colorize(exportedPath, ui.Bold))
	fmt.Println()
	fmt.Println("You can use this file to:")
	fmt.Println("  • Backup your progress")
	fmt.Println("  • Transfer to another computer")
	fmt.Println("  • Share with friends")
	fmt.Println()
	fmt.Println("Press Enter to continue...")
	sc.scanner.Scan()

	return nil
}

// importSave imports a save file from a specified location
func (sc *SettingsCommand) importSave() error {
	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println(ui.Colorize("IMPORT SAVE FILE", ui.Bold+ui.ColorBrightYellow))
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()

	fmt.Println(ui.Colorize("⚠ WARNING: This will replace your current save file!", ui.Bold+ui.ColorRed))
	fmt.Println("Your current save will be backed up before importing.")
	fmt.Println()

	// Confirm action
	if !ui.ConfirmDestructiveAction(sc.scanner, "import a save file") {
		fmt.Println()
		fmt.Println(ui.Colorize("Import cancelled.", ui.ColorYellow))
		fmt.Println("Press Enter to continue...")
		sc.scanner.Scan()
		return nil
	}

	fmt.Println()
	fmt.Println("Enter the path to the save file to import:")
	fmt.Print("> ")

	if !sc.scanner.Scan() {
		return fmt.Errorf("failed to read input")
	}

	sourcePath := strings.TrimSpace(sc.scanner.Text())
	if sourcePath == "" {
		fmt.Println()
		fmt.Println(ui.Colorize("Import cancelled - no path provided.", ui.ColorYellow))
		fmt.Println("Press Enter to continue...")
		sc.scanner.Scan()
		return nil
	}

	// Validate the save file first
	fmt.Println()
	fmt.Println("Validating save file...")
	err := storage.ValidateSaveFile(sourcePath)
	if err != nil {
		fmt.Println()
		fmt.Println(ui.Colorize(fmt.Sprintf("✗ Invalid save file: %v", err), ui.ColorRed))
		fmt.Println()
		fmt.Println("The file may be corrupted or not a valid PokeTacTix save file.")
		fmt.Println("Press Enter to continue...")
		sc.scanner.Scan()
		return err
	}

	fmt.Println(ui.Colorize("✓ Save file is valid!", ui.ColorGreen))
	fmt.Println()

	// Final confirmation
	if !ui.ConfirmationPrompt(sc.scanner, "Proceed with import?", true) {
		fmt.Println()
		fmt.Println(ui.Colorize("Import cancelled.", ui.ColorYellow))
		fmt.Println("Press Enter to continue...")
		sc.scanner.Scan()
		return nil
	}

	// Import the save file
	err = storage.ImportSave(sourcePath)
	if err != nil {
		fmt.Println()
		fmt.Println(ui.Colorize(fmt.Sprintf("✗ Import failed: %v", err), ui.ColorRed))
		fmt.Println("Press Enter to continue...")
		sc.scanner.Scan()
		return err
	}

	// Display success message
	fmt.Println()
	fmt.Println(ui.Colorize("═══════════════════════════════════════", ui.ColorGreen))
	fmt.Println(ui.Colorize("  IMPORT SUCCESSFUL!", ui.Bold+ui.ColorGreen))
	fmt.Println(ui.Colorize("═══════════════════════════════════════", ui.ColorGreen))
	fmt.Println()
	fmt.Println("Save file imported successfully!")
	fmt.Println()
	fmt.Println(ui.Colorize("⚠ IMPORTANT: You must restart the game for changes to take effect.", ui.Bold+ui.ColorYellow))
	fmt.Println()
	fmt.Println("Press Enter to continue...")
	sc.scanner.Scan()

	return nil
}
