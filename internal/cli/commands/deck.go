package commands

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
)

// DeckCommand handles deck-related commands
type DeckCommand struct {
	gameState *storage.GameState
	renderer  *ui.Renderer
	scanner   *bufio.Scanner
}

// NewDeckCommand creates a new deck command handler
func NewDeckCommand(gameState *storage.GameState, renderer *ui.Renderer, scanner *bufio.Scanner) *DeckCommand {
	return &DeckCommand{
		gameState: gameState,
		renderer:  renderer,
		scanner:   scanner,
	}
}

// ViewDeck displays the current deck configuration
func (dc *DeckCommand) ViewDeck() error {
	// Clear screen
	dc.renderer.Clear()

	// Display header
	fmt.Println(ui.RenderLogo())
	fmt.Println()
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println(ui.Colorize("CURRENT DECK", ui.Bold+ui.ColorBrightCyan))
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()

	// Check if deck is empty
	if len(dc.gameState.Deck) == 0 {
		fmt.Println(ui.Colorize("Your deck is empty!", ui.ColorYellow))
		fmt.Println("Use 'deck edit' to add Pokemon to your deck.")
		fmt.Println()
		fmt.Println("Press Enter to continue...")
		dc.scanner.Scan()
		return nil
	}

	// Display deck Pokemon
	fmt.Println("Deck Configuration (5 Pokemon):")
	fmt.Println()

	var totalHP, totalAttack, totalDefense, totalSpeed int
	var totalLevel int
	validCards := 0

	for i, cardIdx := range dc.gameState.Deck {
		// Validate card index
		if cardIdx < 0 || cardIdx >= len(dc.gameState.Collection) {
			fmt.Printf("[%d] (Invalid card index)\n", i+1)
			continue
		}

		card := dc.gameState.Collection[cardIdx]
		stats := card.GetCurrentStats()

		// Accumulate totals
		totalHP += stats.HP
		totalAttack += stats.Attack
		totalDefense += stats.Defense
		totalSpeed += stats.Speed
		totalLevel += card.Level
		validCards++

		// Display position number
		fmt.Printf("[%d] ", i+1)

		// Display Pokemon name with color
		name := card.Name
		if dc.renderer.ColorSupport {
			if card.IsMythical {
				name = ui.Colorize(name, ui.Bold+ui.ColorMagenta)
			} else if card.IsLegendary {
				name = ui.Colorize(name, ui.Bold+ui.ColorYellow)
			} else {
				name = ui.Colorize(name, ui.Bold+ui.ColorBrightCyan)
			}
		}
		fmt.Printf("%s ", name)

		// Display level
		fmt.Printf("(Lv %d, XP: %d/100)\n", card.Level, card.XP)

		// Display types
		fmt.Print("    Types: ")
		for j, t := range card.Types {
			if j > 0 {
				fmt.Print("/")
			}
			if dc.renderer.ColorSupport {
				fmt.Print(ui.ColorizeType(strings.ToUpper(t), t))
			} else {
				fmt.Print(strings.ToUpper(t))
			}
		}
		fmt.Println()

		// Display stats
		fmt.Printf("    Stats: HP: %d | ATK: %d | DEF: %d | SPD: %d | STA: %d\n",
			stats.HP, stats.Attack, stats.Defense, stats.Speed, stats.Stamina)

		// Display moves
		fmt.Print("    Moves: ")
		for j, move := range card.Moves {
			if j > 0 {
				fmt.Print(", ")
			}
			moveName := strings.Title(strings.ReplaceAll(move.Name, "-", " "))
			if dc.renderer.ColorSupport {
				fmt.Print(ui.ColorizeType(moveName, move.Type))
			} else {
				fmt.Print(moveName)
			}
		}
		fmt.Println()
		fmt.Println()
	}

	// Display deck summary statistics
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()
	fmt.Println(ui.Colorize("DECK STATISTICS", ui.Bold+ui.ColorBrightYellow))
	fmt.Println()

	if validCards > 0 {
		avgLevel := float64(totalLevel) / float64(validCards)
		avgHP := float64(totalHP) / float64(validCards)
		avgAttack := float64(totalAttack) / float64(validCards)
		avgDefense := float64(totalDefense) / float64(validCards)
		avgSpeed := float64(totalSpeed) / float64(validCards)

		fmt.Printf("Total Pokemon:   %d/5\n", validCards)
		fmt.Printf("Average Level:   %.1f\n", avgLevel)
		fmt.Printf("Total HP:        %d (Avg: %.1f)\n", totalHP, avgHP)
		fmt.Printf("Total Attack:    %d (Avg: %.1f)\n", totalAttack, avgAttack)
		fmt.Printf("Total Defense:   %d (Avg: %.1f)\n", totalDefense, avgDefense)
		fmt.Printf("Total Speed:     %d (Avg: %.1f)\n", totalSpeed, avgSpeed)
	} else {
		fmt.Println("No valid Pokemon in deck.")
	}

	fmt.Println()
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()
	fmt.Println("Press Enter to continue...")
	dc.scanner.Scan()

	return nil
}

// EditDeck allows the user to modify their deck configuration
func (dc *DeckCommand) EditDeck() error {
	// Create a working copy of the deck for undo functionality
	originalDeck := make([]int, len(dc.gameState.Deck))
	copy(originalDeck, dc.gameState.Deck)

	for {
		// Clear screen and display current deck
		dc.renderer.Clear()

		fmt.Println(ui.RenderLogo())
		fmt.Println()
		fmt.Println(strings.Repeat("═", 80))
		fmt.Println(ui.Colorize("DECK EDITOR", ui.Bold+ui.ColorBrightCyan))
		fmt.Println(strings.Repeat("═", 80))
		fmt.Println()

		// Display current deck
		dc.displayDeckEditor()

		// Display menu options
		fmt.Println()
		fmt.Println(strings.Repeat("─", 80))
		fmt.Println()

		options := []ui.MenuOption{
			{Label: "Add Pokemon", Description: "Add a Pokemon to your deck", Value: "add"},
			{Label: "Remove Pokemon", Description: "Remove a Pokemon from your deck", Value: "remove"},
			{Label: "Reorder Pokemon", Description: "Change Pokemon positions in deck", Value: "reorder"},
			{Label: "Save Deck", Description: "Save changes and exit", Value: "save"},
			{Label: "Cancel", Description: "Discard changes and exit", Value: "cancel"},
		}

		fmt.Println(dc.renderer.RenderBorderedMenu(options, -1, "DECK EDITOR OPTIONS"))
		fmt.Print("Enter your choice (1-5): ")

		// Get user input
		if !dc.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(dc.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 5 {
			fmt.Println(ui.Colorize("Invalid choice. Press Enter to continue...", ui.ColorRed))
			dc.scanner.Scan()
			continue
		}

		switch choice {
		case 1:
			// Add Pokemon
			err := dc.addPokemonToDeck()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				fmt.Println("Press Enter to continue...")
				dc.scanner.Scan()
			}
		case 2:
			// Remove Pokemon
			err := dc.removePokemonFromDeck()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				fmt.Println("Press Enter to continue...")
				dc.scanner.Scan()
			}
		case 3:
			// Reorder Pokemon
			err := dc.reorderDeck()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				fmt.Println("Press Enter to continue...")
				dc.scanner.Scan()
			}
		case 4:
			// Save deck - will be handled in task 7.3
			return dc.saveDeck(originalDeck)
		case 5:
			// Cancel - restore original deck
			dc.gameState.Deck = originalDeck
			fmt.Println()
			fmt.Println(ui.Colorize("Changes discarded.", ui.ColorYellow))
			fmt.Println("Press Enter to continue...")
			dc.scanner.Scan()
			return nil
		}
	}
}

// displayDeckEditor displays the current deck in edit mode
func (dc *DeckCommand) displayDeckEditor() {
	fmt.Printf("Current Deck: %d/5 Pokemon\n", len(dc.gameState.Deck))
	fmt.Println()

	if len(dc.gameState.Deck) == 0 {
		fmt.Println(ui.Colorize("Deck is empty. Add Pokemon to get started!", ui.ColorYellow))
		return
	}

	// Display deck Pokemon in compact format
	for i, cardIdx := range dc.gameState.Deck {
		if cardIdx < 0 || cardIdx >= len(dc.gameState.Collection) {
			fmt.Printf("[%d] (Invalid card)\n", i+1)
			continue
		}

		card := dc.gameState.Collection[cardIdx]
		stats := card.GetCurrentStats()

		// Format: [1] Pikachu (Lv 15) - ELECTRIC - HP: 140 | ATK: 65 | DEF: 50 | SPD: 110
		fmt.Printf("[%d] ", i+1)

		name := card.Name
		if dc.renderer.ColorSupport {
			if card.IsMythical {
				name = ui.Colorize(name, ui.ColorMagenta)
			} else if card.IsLegendary {
				name = ui.Colorize(name, ui.ColorYellow)
			}
		}

		fmt.Printf("%s (Lv %d) - ", name, card.Level)

		// Types
		for j, t := range card.Types {
			if j > 0 {
				fmt.Print("/")
			}
			if dc.renderer.ColorSupport {
				fmt.Print(ui.ColorizeType(strings.ToUpper(t), t))
			} else {
				fmt.Print(strings.ToUpper(t))
			}
		}

		fmt.Printf(" - HP: %d | ATK: %d | DEF: %d | SPD: %d\n",
			stats.HP, stats.Attack, stats.Defense, stats.Speed)
	}
}

// addPokemonToDeck adds a Pokemon to the deck
func (dc *DeckCommand) addPokemonToDeck() error {
	// Check if deck is full
	if len(dc.gameState.Deck) >= 5 {
		fmt.Println()
		fmt.Println(ui.Colorize("Deck is full! Remove a Pokemon first.", ui.ColorRed))
		fmt.Println("Press Enter to continue...")
		dc.scanner.Scan()
		return nil
	}

	// Check if collection is empty
	if len(dc.gameState.Collection) == 0 {
		fmt.Println()
		fmt.Println(ui.Colorize("Your collection is empty!", ui.ColorRed))
		fmt.Println("Press Enter to continue...")
		dc.scanner.Scan()
		return nil
	}

	// Display available Pokemon (not in deck)
	dc.renderer.Clear()
	fmt.Println(ui.RenderLogo())
	fmt.Println()
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println(ui.Colorize("ADD POKEMON TO DECK", ui.Bold+ui.ColorBrightCyan))
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()

	// Get available Pokemon
	availablePokemon := dc.getAvailablePokemon()

	if len(availablePokemon) == 0 {
		fmt.Println(ui.Colorize("All your Pokemon are already in the deck!", ui.ColorYellow))
		fmt.Println("Press Enter to continue...")
		dc.scanner.Scan()
		return nil
	}

	// Display available Pokemon
	fmt.Println("Available Pokemon:")
	fmt.Println()

	for i, cardIdx := range availablePokemon {
		card := dc.gameState.Collection[cardIdx]
		stats := card.GetCurrentStats()

		fmt.Printf("[%d] ", i+1)

		name := card.Name
		if dc.renderer.ColorSupport {
			if card.IsMythical {
				name = ui.Colorize(name, ui.ColorMagenta)
			} else if card.IsLegendary {
				name = ui.Colorize(name, ui.ColorYellow)
			}
		}

		fmt.Printf("%s (Lv %d) - ", name, card.Level)

		// Types
		for j, t := range card.Types {
			if j > 0 {
				fmt.Print("/")
			}
			if dc.renderer.ColorSupport {
				fmt.Print(ui.ColorizeType(strings.ToUpper(t), t))
			} else {
				fmt.Print(strings.ToUpper(t))
			}
		}

		fmt.Printf(" - HP: %d | ATK: %d | DEF: %d | SPD: %d\n",
			stats.HP, stats.Attack, stats.Defense, stats.Speed)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("\nSelect Pokemon to add (1-%d) or 0 to cancel: ", len(availablePokemon))

	// Get user input
	if !dc.scanner.Scan() {
		return fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(dc.scanner.Text())
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 0 || choice > len(availablePokemon) {
		fmt.Println(ui.Colorize("Invalid choice.", ui.ColorRed))
		fmt.Println("Press Enter to continue...")
		dc.scanner.Scan()
		return nil
	}

	if choice == 0 {
		return nil
	}

	// Add Pokemon to deck
	selectedCardIdx := availablePokemon[choice-1]
	dc.gameState.Deck = append(dc.gameState.Deck, selectedCardIdx)

	card := dc.gameState.Collection[selectedCardIdx]
	fmt.Println()
	fmt.Println(ui.Colorize(fmt.Sprintf("✓ %s added to deck!", card.Name), ui.ColorGreen))
	fmt.Println("Press Enter to continue...")
	dc.scanner.Scan()

	return nil
}

// removePokemonFromDeck removes a Pokemon from the deck
func (dc *DeckCommand) removePokemonFromDeck() error {
	// Check if deck is empty
	if len(dc.gameState.Deck) == 0 {
		fmt.Println()
		fmt.Println(ui.Colorize("Deck is empty!", ui.ColorYellow))
		fmt.Println("Press Enter to continue...")
		dc.scanner.Scan()
		return nil
	}

	// Display current deck
	dc.renderer.Clear()
	fmt.Println(ui.RenderLogo())
	fmt.Println()
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println(ui.Colorize("REMOVE POKEMON FROM DECK", ui.Bold+ui.ColorBrightCyan))
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()

	fmt.Println("Current Deck:")
	fmt.Println()

	for i, cardIdx := range dc.gameState.Deck {
		if cardIdx < 0 || cardIdx >= len(dc.gameState.Collection) {
			fmt.Printf("[%d] (Invalid card)\n", i+1)
			continue
		}

		card := dc.gameState.Collection[cardIdx]
		stats := card.GetCurrentStats()

		fmt.Printf("[%d] ", i+1)

		name := card.Name
		if dc.renderer.ColorSupport {
			if card.IsMythical {
				name = ui.Colorize(name, ui.ColorMagenta)
			} else if card.IsLegendary {
				name = ui.Colorize(name, ui.ColorYellow)
			}
		}

		fmt.Printf("%s (Lv %d) - ", name, card.Level)

		// Types
		for j, t := range card.Types {
			if j > 0 {
				fmt.Print("/")
			}
			if dc.renderer.ColorSupport {
				fmt.Print(ui.ColorizeType(strings.ToUpper(t), t))
			} else {
				fmt.Print(strings.ToUpper(t))
			}
		}

		fmt.Printf(" - HP: %d | ATK: %d | DEF: %d | SPD: %d\n",
			stats.HP, stats.Attack, stats.Defense, stats.Speed)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("\nSelect Pokemon to remove (1-%d) or 0 to cancel: ", len(dc.gameState.Deck))

	// Get user input
	if !dc.scanner.Scan() {
		return fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(dc.scanner.Text())
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 0 || choice > len(dc.gameState.Deck) {
		fmt.Println(ui.Colorize("Invalid choice.", ui.ColorRed))
		fmt.Println("Press Enter to continue...")
		dc.scanner.Scan()
		return nil
	}

	if choice == 0 {
		return nil
	}

	// Remove Pokemon from deck
	removeIdx := choice - 1
	cardIdx := dc.gameState.Deck[removeIdx]
	card := dc.gameState.Collection[cardIdx]

	// Remove from deck
	dc.gameState.Deck = append(dc.gameState.Deck[:removeIdx], dc.gameState.Deck[removeIdx+1:]...)

	fmt.Println()
	fmt.Println(ui.Colorize(fmt.Sprintf("✓ %s removed from deck!", card.Name), ui.ColorGreen))
	fmt.Println("Press Enter to continue...")
	dc.scanner.Scan()

	return nil
}

// reorderDeck allows the user to reorder Pokemon in the deck
func (dc *DeckCommand) reorderDeck() error {
	// Check if deck has at least 2 Pokemon
	if len(dc.gameState.Deck) < 2 {
		fmt.Println()
		fmt.Println(ui.Colorize("Need at least 2 Pokemon in deck to reorder.", ui.ColorYellow))
		fmt.Println("Press Enter to continue...")
		dc.scanner.Scan()
		return nil
	}

	// Display current deck
	dc.renderer.Clear()
	fmt.Println(ui.RenderLogo())
	fmt.Println()
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println(ui.Colorize("REORDER DECK", ui.Bold+ui.ColorBrightCyan))
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()

	fmt.Println("Current Deck Order:")
	fmt.Println()

	for i, cardIdx := range dc.gameState.Deck {
		if cardIdx < 0 || cardIdx >= len(dc.gameState.Collection) {
			fmt.Printf("[%d] (Invalid card)\n", i+1)
			continue
		}

		card := dc.gameState.Collection[cardIdx]

		fmt.Printf("[%d] ", i+1)

		name := card.Name
		if dc.renderer.ColorSupport {
			if card.IsMythical {
				name = ui.Colorize(name, ui.ColorMagenta)
			} else if card.IsLegendary {
				name = ui.Colorize(name, ui.ColorYellow)
			}
		}

		fmt.Printf("%s (Lv %d)\n", name, card.Level)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("\nSelect Pokemon to move (1-%d) or 0 to cancel: ", len(dc.gameState.Deck))

	// Get Pokemon to move
	if !dc.scanner.Scan() {
		return fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(dc.scanner.Text())
	fromPos, err := strconv.Atoi(input)
	if err != nil || fromPos < 0 || fromPos > len(dc.gameState.Deck) {
		fmt.Println(ui.Colorize("Invalid choice.", ui.ColorRed))
		fmt.Println("Press Enter to continue...")
		dc.scanner.Scan()
		return nil
	}

	if fromPos == 0 {
		return nil
	}

	// Get new position
	fmt.Printf("Move to position (1-%d): ", len(dc.gameState.Deck))

	if !dc.scanner.Scan() {
		return fmt.Errorf("failed to read input")
	}

	input = strings.TrimSpace(dc.scanner.Text())
	toPos, err := strconv.Atoi(input)
	if err != nil || toPos < 1 || toPos > len(dc.gameState.Deck) {
		fmt.Println(ui.Colorize("Invalid position.", ui.ColorRed))
		fmt.Println("Press Enter to continue...")
		dc.scanner.Scan()
		return nil
	}

	// Convert to 0-based indices
	fromIdx := fromPos - 1
	toIdx := toPos - 1

	if fromIdx == toIdx {
		fmt.Println()
		fmt.Println(ui.Colorize("Pokemon is already in that position.", ui.ColorYellow))
		fmt.Println("Press Enter to continue...")
		dc.scanner.Scan()
		return nil
	}

	// Reorder deck
	cardIdx := dc.gameState.Deck[fromIdx]
	dc.gameState.Deck = append(dc.gameState.Deck[:fromIdx], dc.gameState.Deck[fromIdx+1:]...)
	
	// Insert at new position
	newDeck := make([]int, 0, len(dc.gameState.Deck)+1)
	newDeck = append(newDeck, dc.gameState.Deck[:toIdx]...)
	newDeck = append(newDeck, cardIdx)
	newDeck = append(newDeck, dc.gameState.Deck[toIdx:]...)
	dc.gameState.Deck = newDeck

	card := dc.gameState.Collection[cardIdx]
	fmt.Println()
	fmt.Println(ui.Colorize(fmt.Sprintf("✓ %s moved to position %d!", card.Name, toPos), ui.ColorGreen))
	fmt.Println("Press Enter to continue...")
	dc.scanner.Scan()

	return nil
}

// getAvailablePokemon returns indices of Pokemon not in the deck
func (dc *DeckCommand) getAvailablePokemon() []int {
	// Create a map of Pokemon in deck for quick lookup
	inDeck := make(map[int]bool)
	for _, cardIdx := range dc.gameState.Deck {
		inDeck[cardIdx] = true
	}

	// Find Pokemon not in deck
	var available []int
	for i := range dc.gameState.Collection {
		if !inDeck[i] {
			available = append(available, i)
		}
	}

	return available
}

// saveDeck validates and saves the deck configuration
func (dc *DeckCommand) saveDeck(originalDeck []int) error {
	// Validate deck
	validationErrors := dc.validateDeck()

	if len(validationErrors) > 0 {
		// Display validation errors
		fmt.Println()
		fmt.Println(ui.Colorize("⚠ Cannot save deck - validation errors:", ui.Bold+ui.ColorRed))
		fmt.Println()

		for _, err := range validationErrors {
			fmt.Printf("  • %s\n", err)
		}

		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  [1] Continue editing")
		fmt.Println("  [2] Discard changes and exit")
		fmt.Print("\nEnter your choice (1-2): ")

		if !dc.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(dc.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 2 {
			// Default to continue editing
			return nil
		}

		if choice == 2 {
			// Discard changes
			dc.gameState.Deck = originalDeck
			fmt.Println()
			fmt.Println(ui.Colorize("Changes discarded.", ui.ColorYellow))
			fmt.Println("Press Enter to continue...")
			dc.scanner.Scan()
			return nil
		}

		// Continue editing
		return nil
	}

	// Deck is valid, save it
	err := storage.SaveGameState(dc.gameState)
	if err != nil {
		fmt.Println()
		fmt.Println(ui.Colorize(fmt.Sprintf("✗ Failed to save deck: %v", err), ui.ColorRed))
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  [1] Try again")
		fmt.Println("  [2] Continue editing")
		fmt.Println("  [3] Discard changes and exit")
		fmt.Print("\nEnter your choice (1-3): ")

		if !dc.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(dc.scanner.Text())
		choice, parseErr := strconv.Atoi(input)
		if parseErr != nil || choice < 1 || choice > 3 {
			// Default to continue editing
			return nil
		}

		switch choice {
		case 1:
			// Try saving again
			return dc.saveDeck(originalDeck)
		case 2:
			// Continue editing
			return nil
		case 3:
			// Discard changes
			dc.gameState.Deck = originalDeck
			fmt.Println()
			fmt.Println(ui.Colorize("Changes discarded.", ui.ColorYellow))
			fmt.Println("Press Enter to continue...")
			dc.scanner.Scan()
			return nil
		}
	}

	// Success
	fmt.Println()
	fmt.Println(ui.Colorize("✓ Deck saved successfully!", ui.Bold+ui.ColorGreen))
	fmt.Println()
	fmt.Println("Press Enter to continue...")
	dc.scanner.Scan()

	return nil
}

// validateDeck validates the current deck configuration
func (dc *DeckCommand) validateDeck() []string {
	var errors []string

	// Check if deck has exactly 5 Pokemon
	if len(dc.gameState.Deck) != 5 {
		errors = append(errors, fmt.Sprintf("Deck must have exactly 5 Pokemon (currently has %d)", len(dc.gameState.Deck)))
	}

	// Check for invalid card indices
	for i, cardIdx := range dc.gameState.Deck {
		if cardIdx < 0 || cardIdx >= len(dc.gameState.Collection) {
			errors = append(errors, fmt.Sprintf("Invalid card at position %d", i+1))
		}
	}

	// Check for duplicate Pokemon in deck
	seen := make(map[int]bool)
	for i, cardIdx := range dc.gameState.Deck {
		if seen[cardIdx] {
			if cardIdx >= 0 && cardIdx < len(dc.gameState.Collection) {
				card := dc.gameState.Collection[cardIdx]
				errors = append(errors, fmt.Sprintf("Duplicate Pokemon at position %d: %s", i+1, card.Name))
			} else {
				errors = append(errors, fmt.Sprintf("Duplicate card index at position %d", i+1))
			}
		}
		seen[cardIdx] = true
	}

	return errors
}

// UndoDeckChanges restores the deck to its original state (helper for external use)
func (dc *DeckCommand) UndoDeckChanges(originalDeck []int) {
	dc.gameState.Deck = make([]int, len(originalDeck))
	copy(dc.gameState.Deck, originalDeck)
}
