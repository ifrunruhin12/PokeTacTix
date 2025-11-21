package commands

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"pokemon-cli/internal/battle"
	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
	"pokemon-cli/internal/pokemon"
)

// BattleCommand handles battle-related commands
type BattleCommand struct {
	gameState *storage.GameState
	renderer  *ui.Renderer
	scanner   *bufio.Scanner
}

// NewBattleCommand creates a new battle command handler
func NewBattleCommand(gameState *storage.GameState, renderer *ui.Renderer, scanner *bufio.Scanner) *BattleCommand {
	return &BattleCommand{
		gameState: gameState,
		renderer:  renderer,
		scanner:   scanner,
	}
}

// StartBattle initiates a battle with mode selection (1v1 or 5v5)
func (bc *BattleCommand) StartBattle() error {
	// Check if player has a deck
	if len(bc.gameState.Deck) == 0 {
		return fmt.Errorf("you don't have any Pokemon in your deck. Use 'deck edit' to create a deck")
	}

	if len(bc.gameState.Deck) < 5 {
		return fmt.Errorf("your deck must have exactly 5 Pokemon. Use 'deck edit' to complete your deck")
	}

	// Clear screen and show battle mode selection
	bc.renderer.Clear()
	fmt.Println(ui.RenderLogo())
	fmt.Println()

	// Prompt for battle mode
	modeOptions := []ui.MenuOption{
		{
			Label:       "1v1 Battle",
			Description: "Quick battle with one random Pokemon from your deck",
			Value:       "1v1",
		},
		{
			Label:       "5v5 Battle",
			Description: "Full battle with all 5 Pokemon in your deck",
			Value:       "5v5",
		},
		{
			Label:       "Cancel",
			Description: "Return to main menu",
			Value:       "cancel",
		},
	}

	fmt.Println(bc.renderer.RenderBorderedMenu(modeOptions, 0, "SELECT BATTLE MODE"))
	fmt.Print("Enter your choice (1-3): ")

	var mode string
	for {
		if !bc.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(bc.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 3 {
			fmt.Print("Invalid choice. Enter 1, 2, or 3: ")
			continue
		}

		if choice == 3 {
			fmt.Println("Battle cancelled.")
			return nil
		}

		if choice == 1 {
			mode = "1v1"
		} else {
			mode = "5v5"
		}
		break
	}

	// Load player deck from game state
	playerDeck, err := bc.loadPlayerDeck(mode)
	if err != nil {
		return fmt.Errorf("failed to load player deck: %w", err)
	}

	// Generate AI opponent deck using offline data
	aiDeck, err := bc.generateAIDeck(mode)
	if err != nil {
		return fmt.Errorf("failed to generate AI deck: %w", err)
	}

	// Initialize battle state
	battleState, err := battle.StartBattle(0, mode, playerDeck, aiDeck)
	if err != nil {
		return fmt.Errorf("failed to start battle: %w", err)
	}

	fmt.Println()
	if mode == "1v1" {
		fmt.Printf("Starting 1v1 battle with %s!\n", playerDeck[0].Name)
	} else {
		fmt.Println("Starting 5v5 battle with your full deck!")
	}
	fmt.Println("Press Enter to begin...")
	bc.scanner.Scan()

	// Start the battle loop
	return bc.runBattleLoop(battleState, mode)
}

// loadPlayerDeck loads the player's deck from game state
func (bc *BattleCommand) loadPlayerDeck(mode string) ([]pokemon.Card, error) {
	var playerDeck []pokemon.Card

	if mode == "1v1" {
		// For 1v1, select one random Pokemon from the deck
		if len(bc.gameState.Deck) == 0 {
			return nil, fmt.Errorf("no Pokemon in deck")
		}

		// Use first Pokemon for simplicity (could randomize)
		cardIdx := bc.gameState.Deck[0]
		if cardIdx < 0 || cardIdx >= len(bc.gameState.Collection) {
			return nil, fmt.Errorf("invalid card index in deck")
		}

		playerCard := bc.gameState.Collection[cardIdx]
		playerDeck = []pokemon.Card{playerCard.ToCard()}
	} else {
		// For 5v5, load all 5 Pokemon from deck
		if len(bc.gameState.Deck) != 5 {
			return nil, fmt.Errorf("deck must have exactly 5 Pokemon for 5v5 battles")
		}

		playerDeck = make([]pokemon.Card, 5)
		for i, cardIdx := range bc.gameState.Deck {
			if cardIdx < 0 || cardIdx >= len(bc.gameState.Collection) {
				return nil, fmt.Errorf("invalid card index in deck at position %d", i)
			}

			playerCard := bc.gameState.Collection[cardIdx]
			playerDeck[i] = playerCard.ToCard()
		}
	}

	return playerDeck, nil
}

// generateAIDeck generates an AI opponent deck using offline data
func (bc *BattleCommand) generateAIDeck(mode string) ([]pokemon.Card, error) {
	var aiDeck []pokemon.Card
	var count int

	if mode == "1v1" {
		count = 1
	} else {
		count = 5
	}

	// Generate random Pokemon for AI deck
	for i := 0; i < count; i++ {
		card := pokemon.FetchRandomPokemonCardOffline()
		aiDeck = append(aiDeck, card)
	}

	return aiDeck, nil
}

// runBattleLoop runs the main battle loop
func (bc *BattleCommand) runBattleLoop(bs *battle.BattleState, mode string) error {
	// Main battle loop - continues until battle is over
	for !bs.BattleOver {
		// Clear screen and display battle state
		bc.renderer.Clear()
		fmt.Println(bc.renderer.RenderBattleScreen(bs))
		fmt.Println()

		// Check if it's player's turn
		if bs.WhoseTurn != "player" {
			// This shouldn't happen in our implementation, but handle it
			fmt.Println("Waiting for AI...")
			continue
		}

		// Prompt player for action
		action, moveIdx, err := bc.promptPlayerAction(bs)
		if err != nil {
			return err
		}

		// Handle surrender
		if action == "surrender" {
			if mode == "1v1" {
				fmt.Println("\nYou surrendered the battle!")
			} else {
				fmt.Println("\nYou surrendered this round!")
			}
			fmt.Println("Press Enter to continue...")
			bc.scanner.Scan()
		}

		// Process the move
		logEntries, err := battle.ProcessMove(bs, action, moveIdx)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Press Enter to continue...")
			bc.scanner.Scan()
			continue
		}

		// Display battle log
		bc.renderer.Clear()
		fmt.Println(bc.renderer.RenderBattleScreen(bs))
		fmt.Println()
		fmt.Println(bc.renderer.RenderBattleLogSimple(logEntries, 10))
		fmt.Println()
		fmt.Println("Press Enter to continue...")
		bc.scanner.Scan()

		// Check for Pokemon switching in 5v5 mode
		if mode == "5v5" && !bs.BattleOver {
			playerCard := bs.GetActivePlayerCard()
			if playerCard != nil && playerCard.HP <= 0 && bs.HasPlayerPokemonAlive() {
				// Player's Pokemon was knocked out, must switch
				err := bc.handlePokemonSwitch(bs, true)
				if err != nil {
					return err
				}
			}
		}
	}

	// Battle is over, handle end of battle
	return bc.handleBattleEnd(bs, mode)
}


// promptPlayerAction prompts the player to select an action
func (bc *BattleCommand) promptPlayerAction(bs *battle.BattleState) (string, *int, error) {
	actions := []string{"Attack", "Defend", "Pass", "Sacrifice", "Surrender"}
	
	fmt.Println(bc.renderer.RenderBattleActions(actions, -1))
	fmt.Print("Select action (1-5): ")

	for {
		if !bc.scanner.Scan() {
			return "", nil, fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(bc.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 5 {
			fmt.Print("Invalid choice. Enter 1-5: ")
			continue
		}

		switch choice {
		case 1:
			// Attack - need to select move
			return bc.promptMoveSelection(bs)
		case 2:
			return "defend", nil, nil
		case 3:
			return "pass", nil, nil
		case 4:
			return "sacrifice", nil, nil
		case 5:
			return "surrender", nil, nil
		}
	}
}

// promptMoveSelection prompts the player to select a move for attack
func (bc *BattleCommand) promptMoveSelection(bs *battle.BattleState) (string, *int, error) {
	playerCard := bs.GetActivePlayerCard()
	if playerCard == nil {
		return "", nil, fmt.Errorf("no active Pokemon")
	}

	// Display move selection menu
	bc.renderer.Clear()
	fmt.Println(bc.renderer.RenderBattleScreen(bs))
	fmt.Println()
	fmt.Println(bc.renderer.RenderMoveSelection(playerCard, -1))
	fmt.Print("\nSelect move (1-4) or 0 to go back: ")

	for {
		if !bc.scanner.Scan() {
			return "", nil, fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(bc.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 0 || choice > 4 {
			fmt.Print("Invalid choice. Enter 0-4: ")
			continue
		}

		// Allow canceling back to action menu
		if choice == 0 {
			return bc.promptPlayerAction(bs)
		}

		moveIdx := choice - 1

		// Validate move index
		if moveIdx < 0 || moveIdx >= len(playerCard.Moves) {
			fmt.Print("Invalid move. Enter 1-4 or 0 to go back: ")
			continue
		}

		// Validate stamina availability
		move := playerCard.Moves[moveIdx]
		if playerCard.Stamina < move.StaminaCost {
			fmt.Printf("Not enough stamina! Need %d, have %d. Choose another move: ", move.StaminaCost, playerCard.Stamina)
			continue
		}

		// Return attack action with move index
		return "attack", &moveIdx, nil
	}
}

// handlePokemonSwitch handles Pokemon switching in 5v5 battles
func (bc *BattleCommand) handlePokemonSwitch(bs *battle.BattleState, forced bool) error {
	// Check if player has any Pokemon left
	if !bs.HasPlayerPokemonAlive() {
		return nil // Battle will end
	}

	// Display Pokemon switch menu
	bc.renderer.Clear()
	fmt.Println(bc.renderer.RenderBattleScreen(bs))
	fmt.Println()

	if forced {
		fmt.Println(ui.Colorize("Your Pokemon was knocked out! You must switch to another Pokemon.", ui.ColorRed))
	} else {
		fmt.Println(ui.Colorize("You can switch to another Pokemon.", ui.ColorYellow))
	}
	fmt.Println()

	// Show available Pokemon
	fmt.Println(bc.renderer.RenderPokemonSwitchMenu(bs.PlayerDeck, bs.PlayerActiveIdx, -1))
	fmt.Print("\nSelect Pokemon (1-4): ")

	for {
		if !bc.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(bc.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 {
			fmt.Print("Invalid choice. Enter a valid Pokemon number: ")
			continue
		}

		// Map choice to actual deck index (skipping active and KO'd Pokemon)
		validIdx := 0
		newIdx := -1
		for i, card := range bs.PlayerDeck {
			if i == bs.PlayerActiveIdx || card.IsKnockedOut || card.HP <= 0 {
				continue
			}
			validIdx++
			if validIdx == choice {
				newIdx = i
				break
			}
		}

		if newIdx == -1 {
			fmt.Print("Invalid Pokemon. Choose an available Pokemon: ")
			continue
		}

		// Prevent selection of knocked out Pokemon (double check)
		if bs.PlayerDeck[newIdx].HP <= 0 {
			fmt.Print("That Pokemon is knocked out. Choose another: ")
			continue
		}

		// Switch to the selected Pokemon
		err = battle.SwitchPokemon(bs, newIdx)
		if err != nil {
			fmt.Printf("Error switching Pokemon: %v\n", err)
			fmt.Print("Choose another Pokemon: ")
			continue
		}

		fmt.Printf("\nSwitched to %s!\n", bs.PlayerDeck[newIdx].Name)
		fmt.Println("Press Enter to continue...")
		bc.scanner.Scan()
		return nil
	}
}

// handleBattleEnd handles the end of battle, rewards, and stats updates
func (bc *BattleCommand) handleBattleEnd(bs *battle.BattleState, mode string) error {
	// Clear screen and display final battle state
	bc.renderer.Clear()
	fmt.Println(bc.renderer.RenderBattleScreen(bs))
	fmt.Println()

	// Display battle result
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println()

	var result string
	var coinsEarned int
	var xpPerPokemon int

	switch bs.Winner {
	case "player":
		result = "VICTORY"
		if mode == "1v1" {
			coinsEarned = 50
			xpPerPokemon = 20
		} else {
			coinsEarned = 150
			xpPerPokemon = 15
		}
		fmt.Println(ui.Colorize("🎉 VICTORY! 🎉", ui.Bold+ui.ColorBrightGreen))
	case "ai":
		result = "DEFEAT"
		if mode == "1v1" {
			coinsEarned = 10
		} else {
			coinsEarned = 25
		}
		xpPerPokemon = 0
		fmt.Println(ui.Colorize("💀 DEFEAT 💀", ui.Bold+ui.ColorRed))
	case "draw":
		result = "DRAW"
		if mode == "1v1" {
			coinsEarned = 25
			xpPerPokemon = 10
		} else {
			coinsEarned = 75
			xpPerPokemon = 8
		}
		fmt.Println(ui.Colorize("⚖️  DRAW ⚖️", ui.Bold+ui.ColorYellow))
	}

	fmt.Println()
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println()

	// Award coins
	bc.gameState.Coins += coinsEarned
	fmt.Printf("Coins earned: +%d (Total: %d)\n", coinsEarned, bc.gameState.Coins)
	fmt.Println()

	// Award XP and check for level-ups
	if xpPerPokemon > 0 {
		fmt.Println("Experience gained:")
		leveledUp := false

		for _, deckIdx := range bc.gameState.Deck {
			if deckIdx < 0 || deckIdx >= len(bc.gameState.Collection) {
				continue
			}

			card := &bc.gameState.Collection[deckIdx]
			oldLevel := card.Level
			card.XP += xpPerPokemon

			// Check for level-up (100 XP per level)
			xpNeeded := 100
			if card.XP >= xpNeeded {
				card.Level++
				card.XP -= xpNeeded
				leveledUp = true

				// Display level-up notification
				fmt.Printf("  %s: +%d XP → ", card.Name, xpPerPokemon)
				fmt.Println(ui.Colorize(fmt.Sprintf("LEVEL UP! %d → %d", oldLevel, card.Level), ui.Bold+ui.ColorBrightYellow))

				// Show new stats
				newStats := card.GetCurrentStats()
				fmt.Printf("    New stats: HP: %d, ATK: %d, DEF: %d, SPD: %d\n",
					newStats.HP, newStats.Attack, newStats.Defense, newStats.Speed)
			} else {
				fmt.Printf("  %s: +%d XP (%d/%d to next level)\n", card.Name, xpPerPokemon, card.XP, xpNeeded)
			}

			// Update highest level
			if card.Level > bc.gameState.Stats.HighestLevel {
				bc.gameState.Stats.HighestLevel = card.Level
			}
		}

		if leveledUp {
			fmt.Println()
		}
	}

	// Update player stats
	if mode == "1v1" {
		bc.gameState.Stats.TotalBattles1v1++
		switch result {
		case "VICTORY":
			bc.gameState.Stats.Wins1v1++
		case "DEFEAT":
			bc.gameState.Stats.Losses1v1++
		case "DRAW":
			bc.gameState.Stats.Draws1v1++
		}
	} else {
		bc.gameState.Stats.TotalBattles5v5++
		switch result {
		case "VICTORY":
			bc.gameState.Stats.Wins5v5++
		case "DEFEAT":
			bc.gameState.Stats.Losses5v5++
		case "DRAW":
			bc.gameState.Stats.Draws5v5++
		}
	}

	bc.gameState.Stats.TotalCoinsEarned += coinsEarned

	// Add battle to history
	battleRecord := storage.BattleRecord{
		Mode:        mode,
		Result:      strings.ToLower(result),
		CoinsEarned: coinsEarned,
		Timestamp:   bs.CreatedAt,
	}
	bc.gameState.BattleHistory = append(bc.gameState.BattleHistory, battleRecord)

	// Keep only last 20 battles in history
	if len(bc.gameState.BattleHistory) > 20 {
		bc.gameState.BattleHistory = bc.gameState.BattleHistory[len(bc.gameState.BattleHistory)-20:]
	}

	// Update shop refresh counter
	bc.gameState.ShopState.BattlesSinceRefresh++

	// Save updated game state
	err := storage.SaveGameState(bc.gameState)
	if err != nil {
		fmt.Printf("Warning: Failed to save game state: %v\n", err)
	} else {
		fmt.Println()
		fmt.Println(ui.Colorize("✓ Game saved successfully", ui.ColorGreen))
	}

	fmt.Println()
	fmt.Println("Press Enter to continue...")
	bc.scanner.Scan()

	// Handle post-battle Pokemon selection for 5v5 victories
	if mode == "5v5" && bs.Winner == "player" {
		return bc.handlePostBattlePokemonSelection(bs)
	}

	return nil
}

// handlePostBattlePokemonSelection handles Pokemon selection after 5v5 victory
func (bc *BattleCommand) handlePostBattlePokemonSelection(bs *battle.BattleState) error {
	// Clear screen
	bc.renderer.Clear()

	// Display victory bonus message
	fmt.Println()
	fmt.Println(ui.Colorize("🎁 VICTORY BONUS! 🎁", ui.Bold+ui.ColorBrightGreen))
	fmt.Println()
	fmt.Println("You can select ONE Pokemon from the AI's team to add to your collection!")
	fmt.Println()
	fmt.Println(strings.Repeat("═", 70))
	fmt.Println()

	// Display all AI Pokemon with full details
	for i, aiCard := range bs.AIDeck {
		fmt.Printf("[%d] ", i+1)
		if bc.renderer.ColorSupport {
			fmt.Print(ui.Colorize(aiCard.Name, ui.Bold+ui.ColorBrightCyan))
		} else {
			fmt.Print(aiCard.Name)
		}
		fmt.Printf(" (Lv %d)\n", aiCard.Level)

		// Types
		fmt.Print("    Types: ")
		for j, t := range aiCard.Types {
			if j > 0 {
				fmt.Print("/")
			}
			if bc.renderer.ColorSupport {
				fmt.Print(ui.ColorizeType(strings.ToUpper(t), t))
			} else {
				fmt.Print(strings.ToUpper(t))
			}
		}
		fmt.Println()

		// Stats
		fmt.Printf("    HP: %d | ATK: %d | DEF: %d | SPD: %d\n",
			aiCard.HPMax, aiCard.Attack, aiCard.Defense, aiCard.Speed)

		// Moves
		fmt.Print("    Moves: ")
		for j, move := range aiCard.Moves {
			if j > 0 {
				fmt.Print(", ")
			}
			moveName := strings.Title(strings.ReplaceAll(move.Name, "-", " "))
			if bc.renderer.ColorSupport {
				fmt.Print(ui.ColorizeType(moveName, move.Type))
			} else {
				fmt.Print(moveName)
			}
		}
		fmt.Println()
		fmt.Println()
	}

	fmt.Println(strings.Repeat("═", 70))
	fmt.Print("\nSelect Pokemon to add to your collection (1-5): ")

	// Prompt for selection
	for {
		if !bc.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(bc.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 5 {
			fmt.Print("Invalid choice. Enter 1-5: ")
			continue
		}

		selectedIdx := choice - 1
		selectedCard := bs.AIDeck[selectedIdx]

		// Convert BattleCard to PlayerCard
		newCard := storage.PlayerCard{
			ID:          len(bc.gameState.Collection), // Assign new ID
			PokemonID:   selectedCard.CardID,
			Name:        selectedCard.Name,
			Level:       1, // Add at level 1
			XP:          0,
			BaseHP:      selectedCard.HPMax,
			BaseAttack:  selectedCard.Attack,
			BaseDefense: selectedCard.Defense,
			BaseSpeed:   selectedCard.Speed,
			Types:       selectedCard.Types,
			Moves:       selectedCard.Moves,
			Sprite:      selectedCard.Sprite,
			IsLegendary: false, // Will be set correctly if needed
			IsMythical:  false,
			AcquiredAt:  bs.CreatedAt,
		}

		// Add to collection
		bc.gameState.Collection = append(bc.gameState.Collection, newCard)
		bc.gameState.Stats.TotalPokemon = len(bc.gameState.Collection)

		// Save game state
		err = storage.SaveGameState(bc.gameState)
		if err != nil {
			fmt.Printf("Warning: Failed to save game state: %v\n", err)
		}

		// Display confirmation
		fmt.Println()
		fmt.Println(ui.Colorize(fmt.Sprintf("✓ %s has been added to your collection!", selectedCard.Name), ui.Bold+ui.ColorBrightGreen))
		fmt.Println()
		fmt.Println("Press Enter to continue...")
		bc.scanner.Scan()

		return nil
	}
}
