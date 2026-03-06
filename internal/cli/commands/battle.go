package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"pokemon-cli/internal/battle"
	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
	"pokemon-cli/internal/pokemon"
)

type BattleCommand struct {
	gameState *storage.GameState
	renderer  *ui.Renderer
	scanner   *bufio.Scanner
}

func NewBattleCommand(gameState *storage.GameState, renderer *ui.Renderer, scanner *bufio.Scanner) *BattleCommand {
	return &BattleCommand{
		gameState: gameState,
		renderer:  renderer,
		scanner:   scanner,
	}
}

func (bc *BattleCommand) StartBattle() error {
	if len(bc.gameState.Deck) == 0 {
		return fmt.Errorf("you don't have any Pokemon in your deck. Use 'deck edit' to create a deck")
	}

	if len(bc.gameState.Deck) < 5 {
		return fmt.Errorf("your deck must have exactly 5 Pokemon. Use 'deck edit' to complete your deck")
	}

	bc.renderer.Clear()
	fmt.Println(ui.RenderLogo())
	fmt.Println()

	modeOptions := []ui.MenuOption{
		{
			Label:       "1v1 Battle (costs 1 token)",
			Description: "Quick battle with one random Pokemon from your deck",
			Value:       "1v1",
		},
		{
			Label:       "5v5 Battle (costs 1 token)",
			Description: "Full battle with all 5 Pokemon in your deck",
			Value:       "5v5",
		},
		{
			Label:       "View Token Info",
			Description: "Show current tokens, reset time, and time until reset",
			Value:       "token_info",
		},
		{
			Label:       "Cancel",
			Description: "Return to main menu",
			Value:       "cancel",
		},
	}

	fmt.Println(bc.renderer.RenderBorderedMenu(modeOptions, 0, "SELECT BATTLE MODE"))
	fmt.Print("Enter your choice (1-4): ")

	var mode string
	for {
		if !bc.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(bc.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 4 {
			fmt.Print("Invalid choice. Enter 1, 2, 3, or 4: ")
			continue
		}

		if choice == 4 {
			fmt.Println("Battle cancelled.")
			return nil
		}

		if choice == 3 {
			// Display token info
			bc.displayTokenInfo()
			fmt.Println()
			fmt.Println("Press Enter to continue...")
			bc.scanner.Scan()
			// Return to start of battle menu
			return bc.StartBattle()
		}

		if choice == 1 {
			mode = "1v1"
		} else {
			mode = "5v5"
		}
		break
	}

	// Check if user has enough tokens
	if bc.gameState.GameTokens < 1 {
		fmt.Println()
		fmt.Println(ui.Colorize("❌ Insufficient Tokens!", ui.ColorRed))
		fmt.Println()
		fmt.Println("You don't have enough tokens to start a battle.")
		fmt.Printf("Current tokens: %d/5\n", bc.gameState.GameTokens)
		fmt.Println()
		bc.displayTokenInfo()
		fmt.Println()
		fmt.Println("You can purchase more tokens from the shop or wait for the daily reset.")
		fmt.Println()
		fmt.Println("Press Enter to return to main menu...")
		bc.scanner.Scan()
		return nil
	}

	playerDeck, err := bc.loadPlayerDeck(mode)
	if err != nil {
		return fmt.Errorf("failed to load player deck: %w", err)
	}

	aiDeck, err := bc.generateAIDeck(mode)
	if err != nil {
		return fmt.Errorf("failed to generate AI deck: %w", err)
	}

	// Consume 1 token before starting the battle
	bc.gameState.GameTokens--
	
	// Save the game state to persist token consumption
	err = storage.SaveGameState(bc.gameState)
	if err != nil {
		// Rollback token consumption if save fails
		bc.gameState.GameTokens++
		return fmt.Errorf("failed to save game state: %w", err)
	}

	battleState, err := battle.StartBattle(0, mode, playerDeck, aiDeck)
	if err != nil {
		// Rollback token consumption if battle start fails
		bc.gameState.GameTokens++
		storage.SaveGameState(bc.gameState)
		return fmt.Errorf("failed to start battle: %w", err)
	}

	fmt.Println()
	fmt.Println(ui.Colorize("✓ 1 token consumed", ui.ColorYellow))
	if mode == "1v1" {
		fmt.Printf("Starting 1v1 battle with %s!\n", playerDeck[0].Name)
	} else {
		fmt.Println("Starting 5v5 battle with your full deck!")
	}
	fmt.Printf("Tokens remaining: %d\n", bc.gameState.GameTokens)
	fmt.Println("Press Enter to begin...")
	bc.scanner.Scan()

	return bc.runBattleLoop(battleState, mode)
}

func (bc *BattleCommand) loadPlayerDeck(mode string) ([]pokemon.Card, error) {
	var playerDeck []pokemon.Card

	if mode == "1v1" {
		if len(bc.gameState.Deck) == 0 {
			return nil, fmt.Errorf("no Pokemon in deck")
		}

		cardIdx := bc.gameState.Deck[0]
		if cardIdx < 0 || cardIdx >= len(bc.gameState.Collection) {
			return nil, fmt.Errorf("invalid card index in deck")
		}

		playerCard := bc.gameState.Collection[cardIdx]
		playerDeck = []pokemon.Card{playerCard.ToCard()}
	} else {
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

// displayTokenInfo shows detailed token information
func (bc *BattleCommand) displayTokenInfo() {
	// Get configured reset time (default 00:00 UTC)
	resetTimeStr := os.Getenv("TOKEN_RESET_TIME")
	if resetTimeStr == "" {
		resetTimeStr = "00:00"
	}
	
	// Parse reset time
	resetHour := 0
	resetMinute := 0
	fmt.Sscanf(resetTimeStr, "%d:%d", &resetHour, &resetMinute)
	
	// Get current time in UTC
	now := time.Now().UTC()
	
	// Calculate next reset time
	nextReset := time.Date(now.Year(), now.Month(), now.Day(), resetHour, resetMinute, 0, 0, time.UTC)
	
	// If reset time has passed today, move to tomorrow
	if now.After(nextReset) || now.Equal(nextReset) {
		nextReset = nextReset.Add(24 * time.Hour)
	}
	
	// Calculate duration until reset
	duration := nextReset.Sub(now)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	
	fmt.Println(ui.RenderDivider(60, "═"))
	fmt.Println(ui.Colorize("🎫 TOKEN INFORMATION", ui.Bold+ui.ColorBrightCyan))
	fmt.Println(ui.RenderDivider(60, "═"))
	fmt.Println()
	fmt.Printf("Current Tokens: %d/5\n", bc.gameState.GameTokens)
	fmt.Printf("Daily Reset: %s UTC (in %dh %dm)\n", resetTimeStr, hours, minutes)
	fmt.Printf("Tokens Purchased Today: %d/10\n", bc.gameState.TokensPurchasedToday)
	fmt.Println()
	fmt.Println("Each battle (1v1 or 5v5) costs 1 token.")
	fmt.Println("Tokens reset to 5 daily at the configured time.")
	fmt.Println("You can purchase additional tokens from the shop (100 coins each).")
	fmt.Println(ui.RenderDivider(60, "═"))
}

func (bc *BattleCommand) generateAIDeck(mode string) ([]pokemon.Card, error) {
	var aiDeck []pokemon.Card
	var count int

	if mode == "1v1" {
		count = 1
	} else {
		count = 5
	}

	for i := 0; i < count; i++ {
		card := pokemon.FetchRandomPokemonCardOffline()
		aiDeck = append(aiDeck, card)
	}

	return aiDeck, nil
}

func (bc *BattleCommand) runBattleLoop(bs *battle.BattleState, mode string) error {
	quickBattle := bc.gameState.Settings.QuickBattle

	for !bs.BattleOver {
		bc.renderer.Clear()

		if quickBattle {
			fmt.Println(bc.renderer.RenderBattleScreenCondensed(bs))
		} else {
			fmt.Println(bc.renderer.RenderBattleScreen(bs))
		}
		fmt.Println()

		if bs.WhoseTurn != "player" {
			fmt.Println("Waiting for AI...")
			continue
		}

		action, moveIdx, err := bc.promptPlayerAction(bs)
		if err != nil {
			return err
		}

		if action == "surrender" {
			if mode == "1v1" {
				fmt.Println("\nYou surrendered the battle!")
			} else {
				fmt.Println("\nYou surrendered this round!")
			}
			if !quickBattle {
				fmt.Println("Press Enter to continue...")
				bc.scanner.Scan()
			}
		}

		logEntries, err := battle.ProcessMove(bs, action, moveIdx)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			if !quickBattle {
				fmt.Println("Press Enter to continue...")
				bc.scanner.Scan()
			}
			continue
		}

		bc.renderer.Clear()

		if quickBattle {
			fmt.Println(bc.renderer.RenderBattleScreenCondensed(bs))
		} else {
			fmt.Println(bc.renderer.RenderBattleScreen(bs))
		}
		fmt.Println()

		if quickBattle {
			fmt.Println(bc.renderer.RenderBattleLogSimple(logEntries, 5))
			ui.Sleep(bc.gameState.Settings.BattleSpeed, "short")
		} else {
			fmt.Println(bc.renderer.RenderBattleLogSimple(logEntries, 10))
			fmt.Println("Press Enter to continue...")
			bc.scanner.Scan()
		}

		if mode == "5v5" && !bs.BattleOver {
			playerCard := bs.GetActivePlayerCard()
			if playerCard != nil && playerCard.HP <= 0 && bs.HasPlayerPokemonAlive() {
				err := bc.handlePokemonSwitch(bs, true)
				if err != nil {
					return err
				}
			}
		}
	}

	return bc.handleBattleEnd(bs, mode)
}

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

func (bc *BattleCommand) promptMoveSelection(bs *battle.BattleState) (string, *int, error) {
	playerCard := bs.GetActivePlayerCard()
	if playerCard == nil {
		return "", nil, fmt.Errorf("no active Pokemon")
	}

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

		if choice == 0 {
			return bc.promptPlayerAction(bs)
		}

		moveIdx := choice - 1

		if moveIdx < 0 || moveIdx >= len(playerCard.Moves) {
			fmt.Print("Invalid move. Enter 1-4 or 0 to go back: ")
			continue
		}

		move := playerCard.Moves[moveIdx]
		if playerCard.Stamina < move.StaminaCost {
			fmt.Printf("Not enough stamina! Need %d, have %d. Choose another move: ", move.StaminaCost, playerCard.Stamina)
			continue
		}

		return "attack", &moveIdx, nil
	}
}

func (bc *BattleCommand) handlePokemonSwitch(bs *battle.BattleState, forced bool) error {
	if !bs.HasPlayerPokemonAlive() {
		return nil
	}

	bc.renderer.Clear()
	fmt.Println(bc.renderer.RenderBattleScreen(bs))
	fmt.Println()

	if forced {
		fmt.Println(ui.Colorize("Your Pokemon was knocked out! You must switch to another Pokemon.", ui.ColorRed))
	} else {
		fmt.Println(ui.Colorize("You can switch to another Pokemon.", ui.ColorYellow))
	}
	fmt.Println()

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

		if bs.PlayerDeck[newIdx].HP <= 0 {
			fmt.Print("That Pokemon is knocked out. Choose another: ")
			continue
		}

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

func (bc *BattleCommand) handleBattleEnd(bs *battle.BattleState, mode string) error {
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

	bc.gameState.Coins += coinsEarned
	fmt.Printf("Coins earned: +%d (Total: %d)\n", coinsEarned, bc.gameState.Coins)
	fmt.Println()

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

			xpNeeded := 100
			if card.XP >= xpNeeded {
				card.Level++
				card.XP -= xpNeeded
				leveledUp = true

				fmt.Printf("  %s: +%d XP → ", card.Name, xpPerPokemon)
				fmt.Println(ui.Colorize(fmt.Sprintf("LEVEL UP! %d → %d", oldLevel, card.Level), ui.Bold+ui.ColorBrightYellow))

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

	battleRecord := storage.BattleRecord{
		Mode:        mode,
		Result:      strings.ToLower(result),
		CoinsEarned: coinsEarned,
		Timestamp:   bs.CreatedAt,
	}
	bc.gameState.BattleHistory = append(bc.gameState.BattleHistory, battleRecord)

	if len(bc.gameState.BattleHistory) > 20 {
		bc.gameState.BattleHistory = bc.gameState.BattleHistory[len(bc.gameState.BattleHistory)-20:]
	}

	bc.gameState.ShopState.BattlesSinceRefresh++

	err := storage.SaveGameState(bc.gameState)
	if err != nil {
		fmt.Printf("Warning: Failed to save game state: %v\n", err)
	} else {
		fmt.Println()
		fmt.Println(ui.Colorize("✓ Game saved successfully", ui.ColorGreen))
	}

	if bc.gameState.ShopState.BattlesSinceRefresh >= 10 {
		shopCmd := NewShopCommand(bc.gameState, bc.renderer, bc.scanner)
		if err := shopCmd.CheckAndRefreshShop(); err != nil {
			fmt.Printf("Warning: Failed to refresh shop: %v\n", err)
		}
		if err := storage.SaveGameState(bc.gameState); err != nil {
			fmt.Printf("Warning: Failed to save shop refresh: %v\n", err)
		}
	}

	fmt.Println()
	fmt.Println("Press Enter to continue...")
	bc.scanner.Scan()

	if mode == "5v5" && bs.Winner == "player" {
		return bc.handlePostBattlePokemonSelection(bs)
	}

	return nil
}

func (bc *BattleCommand) handlePostBattlePokemonSelection(bs *battle.BattleState) error {
	bc.renderer.Clear()

	fmt.Println()
	fmt.Println(ui.Colorize("🎁 VICTORY BONUS! 🎁", ui.Bold+ui.ColorBrightGreen))
	fmt.Println()
	fmt.Println("You can select ONE Pokemon from the AI's team to add to your collection!")
	fmt.Println()
	fmt.Println(strings.Repeat("═", 70))
	fmt.Println()

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

		bc.gameState.Collection = append(bc.gameState.Collection, newCard)
		bc.gameState.Stats.TotalPokemon = len(bc.gameState.Collection)

		err = storage.SaveGameState(bc.gameState)
		if err != nil {
			fmt.Printf("Warning: Failed to save game state: %v\n", err)
		}

		fmt.Println()
		fmt.Println(ui.Colorize(fmt.Sprintf("✓ %s has been added to your collection!", selectedCard.Name), ui.Bold+ui.ColorBrightGreen))
		fmt.Println()
		fmt.Println("Press Enter to continue...")
		bc.scanner.Scan()

		return nil
	}
}
