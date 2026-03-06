package commands

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
	"pokemon-cli/internal/pokemon"
)

// ShopCommand handles shop-related commands
type ShopCommand struct {
	gameState *storage.GameState
	renderer  *ui.Renderer
	scanner   *bufio.Scanner
}

// NewShopCommand creates a new shop command handler
func NewShopCommand(gameState *storage.GameState, renderer *ui.Renderer, scanner *bufio.Scanner) *ShopCommand {
	return &ShopCommand{
		gameState: gameState,
		renderer:  renderer,
		scanner:   scanner,
	}
}

// GenerateShopInventory generates a new shop inventory using offline Pokemon data
// Generates 10-15 Pokemon with pricing: common 100, uncommon 250, rare 500
// Excludes legendary and mythical Pokemon
func (sc *ShopCommand) GenerateShopInventory() error {
	// Determine number of items (10-15)
	numItems := 10 + rand.Intn(6) // 10-15 items

	inventory := make([]storage.ShopItem, 0, numItems)

	// Generate Pokemon for shop
	for i := 0; i < numItems; i++ {
		// Get random non-legendary, non-mythical Pokemon
		pokemonEntry, err := pokemon.GetRandomPokemon(true, true)
		if err != nil {
			return fmt.Errorf("failed to generate shop inventory: %w", err)
		}

		// Calculate rarity and price based on base stat total
		statTotal := pokemonEntry.HP + pokemonEntry.Attack + pokemonEntry.Defense + pokemonEntry.Speed
		
		var rarity string
		var price int

		if statTotal >= 500 {
			rarity = "rare"
			price = 500
		} else if statTotal >= 400 {
			rarity = "uncommon"
			price = 250
		} else {
			rarity = "common"
			price = 100
		}

		// Create shop item
		shopItem := storage.ShopItem{
			PokemonID:   pokemonEntry.ID,
			Name:        pokemonEntry.Name,
			Types:       pokemonEntry.Types,
			BaseHP:      pokemonEntry.HP,
			BaseAttack:  pokemonEntry.Attack,
			BaseDefense: pokemonEntry.Defense,
			BaseSpeed:   pokemonEntry.Speed,
			Moves:       pokemonEntry.Moves,
			Sprite:      pokemonEntry.Sprite,
			Price:       price,
			Rarity:      rarity,
			IsLegendary: pokemonEntry.IsLegendary,
			IsMythical:  pokemonEntry.IsMythical,
		}

		inventory = append(inventory, shopItem)
	}

	// Update shop state
	sc.gameState.ShopState.Inventory = inventory
	sc.gameState.ShopState.LastRefresh = time.Now()
	sc.gameState.ShopState.BattlesSinceRefresh = 0

	return nil
}

// ViewShop displays the shop inventory
func (sc *ShopCommand) ViewShop() error {
	// Check if shop needs initialization
	if len(sc.gameState.ShopState.Inventory) == 0 {
		fmt.Println(ui.Colorize("Generating shop inventory...", ui.ColorYellow))
		if err := sc.GenerateShopInventory(); err != nil {
			return err
		}
		// Save the new inventory
		if err := storage.SaveGameState(sc.gameState); err != nil {
			fmt.Println(ui.Colorize("Warning: Failed to save shop inventory", ui.ColorRed))
		}
	}

	for {
		// Clear screen
		sc.renderer.Clear()

		// Display header
		fmt.Println(ui.RenderLogo())
		fmt.Println()
		fmt.Println(strings.Repeat("═", 80))
		fmt.Println(ui.Colorize("POKEMON SHOP", ui.Bold+ui.ColorBrightCyan))
		fmt.Println(strings.Repeat("═", 80))
		fmt.Println()

		// Display player coins
		coinsText := fmt.Sprintf("Your Coins: %d", sc.gameState.Coins)
		if sc.renderer.ColorSupport {
			coinsText = ui.Colorize(coinsText, ui.Bold+ui.ColorYellow)
		}
		fmt.Println(coinsText)
		fmt.Println()

		// Display shop refresh info
		battlesUntilRefresh := 10 - sc.gameState.ShopState.BattlesSinceRefresh
		if battlesUntilRefresh < 0 {
			battlesUntilRefresh = 0
		}
		refreshText := fmt.Sprintf("Shop refreshes in %d battles", battlesUntilRefresh)
		if sc.renderer.ColorSupport {
			refreshText = ui.Colorize(refreshText, ui.ColorCyan)
		}
		fmt.Println(refreshText)
		fmt.Println()

		// Display shop inventory in grid format
		sc.displayShopGrid()

		fmt.Println()
		fmt.Println(strings.Repeat("═", 80))
		fmt.Println()

		// Display options
		fmt.Println("Options:")
		fmt.Println("  [1-" + strconv.Itoa(len(sc.gameState.ShopState.Inventory)) + "] Buy Pokemon by number")
		fmt.Println("  [T] Buy Game Tokens (100 coins each, 10/day max, resets at " + getResetTimeString() + " UTC)")
		fmt.Println("  [R] Refresh shop (costs 50 coins)")
		fmt.Println("  [Q] Back to menu")
		fmt.Println()
		fmt.Print("Enter your choice: ")

		// Get user input
		if !sc.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}

		input := strings.ToUpper(strings.TrimSpace(sc.scanner.Text()))

		switch input {
		case "Q":
			return nil
		case "T":
			// Token purchase
			if err := sc.PurchaseTokens(); err != nil {
				fmt.Println()
				fmt.Println(ui.Colorize(fmt.Sprintf("Error: %v", err), ui.ColorRed))
				fmt.Println("Press Enter to continue...")
				sc.scanner.Scan()
			}
		case "R":
			// Manual refresh for 50 coins
			if sc.gameState.Coins < 50 {
				fmt.Println()
				fmt.Println(ui.Colorize("Not enough coins! Manual refresh costs 50 coins.", ui.ColorRed))
				fmt.Println("Press Enter to continue...")
				sc.scanner.Scan()
				continue
			}

			fmt.Println()
			fmt.Print("Refresh shop for 50 coins? (y/n): ")
			if !sc.scanner.Scan() {
				return fmt.Errorf("failed to read input")
			}
			confirm := strings.ToLower(strings.TrimSpace(sc.scanner.Text()))
			if confirm == "y" || confirm == "yes" {
				sc.gameState.Coins -= 50
				if err := sc.GenerateShopInventory(); err != nil {
					return err
				}
				if err := storage.SaveGameState(sc.gameState); err != nil {
					fmt.Println(ui.Colorize("Warning: Failed to save game state", ui.ColorRed))
				}
				fmt.Println(ui.Colorize("Shop refreshed!", ui.ColorGreen))
				time.Sleep(1 * time.Second)
			}
		default:
			// Try to parse as number for purchase
			choice, err := strconv.Atoi(input)
			if err != nil || choice < 1 || choice > len(sc.gameState.ShopState.Inventory) {
				fmt.Println()
				fmt.Println(ui.Colorize("Invalid choice. Please try again.", ui.ColorRed))
				fmt.Println("Press Enter to continue...")
				sc.scanner.Scan()
				continue
			}

			// Buy Pokemon
			if err := sc.BuyPokemon(choice - 1); err != nil {
				fmt.Println()
				fmt.Println(ui.Colorize(fmt.Sprintf("Error: %v", err), ui.ColorRed))
				fmt.Println("Press Enter to continue...")
				sc.scanner.Scan()
			}
		}
	}
}

// displayShopGrid displays the shop inventory in a grid format
func (sc *ShopCommand) displayShopGrid() {
	inventory := sc.gameState.ShopState.Inventory

	// Display table header
	fmt.Printf("%-4s %-15s %-20s %-8s %-8s %-10s\n",
		"#", "NAME", "TYPES", "RARITY", "PRICE", "OWNED")
	fmt.Println(strings.Repeat("-", 80))

	// Display each item
	for i, item := range inventory {
		// Format number
		num := fmt.Sprintf("%d.", i+1)

		// Format name
		name := item.Name

		// Format types
		typeStr := ""
		for j, t := range item.Types {
			if j > 0 {
				typeStr += "/"
			}
			if sc.renderer.ColorSupport {
				typeStr += ui.ColorizeType(strings.ToUpper(t), t)
			} else {
				typeStr += strings.ToUpper(t)
			}
		}

		// Format rarity with color
		rarityStr := strings.ToUpper(item.Rarity)
		if sc.renderer.ColorSupport {
			switch item.Rarity {
			case "rare":
				rarityStr = ui.Colorize(rarityStr, ui.ColorMagenta)
			case "uncommon":
				rarityStr = ui.Colorize(rarityStr, ui.ColorBlue)
			case "common":
				rarityStr = ui.Colorize(rarityStr, ui.ColorWhite)
			}
		}

		// Format price
		priceStr := fmt.Sprintf("%d", item.Price)
		if sc.renderer.ColorSupport {
			priceStr = ui.Colorize(priceStr, ui.ColorYellow)
		}

		// Check if owned
		ownedStr := ""
		ownedCount := sc.countOwnedPokemon(item.PokemonID)
		if ownedCount > 0 {
			ownedStr = fmt.Sprintf("x%d", ownedCount)
			if sc.renderer.ColorSupport {
				ownedStr = ui.Colorize(ownedStr, ui.ColorGreen)
			}
		}

		fmt.Printf("%-4s %-15s %-20s %-8s %-8s %-10s\n",
			num, name, typeStr, rarityStr, priceStr, ownedStr)
	}

	// Display stats summary
	fmt.Println()
	fmt.Println("Stats Preview (at Level 1):")
	fmt.Printf("%-4s %-6s %-6s %-6s %-6s\n", "#", "HP", "ATK", "DEF", "SPD")
	fmt.Println(strings.Repeat("-", 30))
	for i, item := range inventory {
		fmt.Printf("%-4d %-6d %-6d %-6d %-6d\n",
			i+1, item.BaseHP, item.BaseAttack, item.BaseDefense, item.BaseSpeed)
	}
}

// BuyPokemon handles the purchase of a Pokemon from the shop
func (sc *ShopCommand) BuyPokemon(index int) error {
	if index < 0 || index >= len(sc.gameState.ShopState.Inventory) {
		return fmt.Errorf("invalid Pokemon selection")
	}

	item := sc.gameState.ShopState.Inventory[index]

	// Check if player has enough coins
	if sc.gameState.Coins < item.Price {
		return fmt.Errorf("not enough coins! You need %d coins but only have %d", item.Price, sc.gameState.Coins)
	}

	// Display purchase details and get confirmation
	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println(ui.Colorize("PURCHASE CONFIRMATION", ui.Bold+ui.ColorBrightYellow))
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()
	fmt.Printf("Pokemon: %s\n", ui.Colorize(item.Name, ui.Bold))
	
	// Display types
	typeStr := "Types: "
	for i, t := range item.Types {
		if i > 0 {
			typeStr += ", "
		}
		if sc.renderer.ColorSupport {
			typeStr += ui.ColorizeType(strings.ToUpper(t), t)
		} else {
			typeStr += strings.ToUpper(t)
		}
	}
	fmt.Println(typeStr)
	
	fmt.Printf("Stats: HP=%d, ATK=%d, DEF=%d, SPD=%d\n",
		item.BaseHP, item.BaseAttack, item.BaseDefense, item.BaseSpeed)
	fmt.Printf("Rarity: %s\n", strings.ToUpper(item.Rarity))
	fmt.Println()

	// Use confirmation prompt for expensive purchases (>= 250 coins)
	var confirmed bool
	if item.Price >= 250 {
		confirmed = ui.ConfirmExpensivePurchase(sc.scanner, item.Name, item.Price, sc.gameState.Coins)
	} else {
		// For cheaper items, use simple confirmation
		confirmed = ui.ConfirmationPrompt(sc.scanner, "Confirm purchase?", true)
	}

	if !confirmed {
		fmt.Println()
		fmt.Println(ui.Colorize("Purchase cancelled.", ui.ColorYellow))
		time.Sleep(1 * time.Second)
		return nil
	}

	// Deduct coins
	sc.gameState.Coins -= item.Price

	// Add Pokemon to collection at level 1
	newCard := storage.PlayerCard{
		ID:          sc.getNextCardID(),
		PokemonID:   item.PokemonID,
		Name:        item.Name,
		Level:       1,
		XP:          0,
		BaseHP:      item.BaseHP,
		BaseAttack:  item.BaseAttack,
		BaseDefense: item.BaseDefense,
		BaseSpeed:   item.BaseSpeed,
		Types:       item.Types,
		Moves:       item.Moves,
		Sprite:      item.Sprite,
		IsLegendary: item.IsLegendary,
		IsMythical:  item.IsMythical,
		AcquiredAt:  time.Now(),
	}

	sc.gameState.Collection = append(sc.gameState.Collection, newCard)
	sc.gameState.Stats.TotalPokemon = len(sc.gameState.Collection)

	// Save game state
	if err := storage.SaveGameState(sc.gameState); err != nil {
		return fmt.Errorf("failed to save game state: %w", err)
	}

	// Display success message
	fmt.Println()
	fmt.Println(ui.Colorize("═══════════════════════════════════════", ui.ColorGreen))
	fmt.Println(ui.Colorize("  PURCHASE SUCCESSFUL!", ui.Bold+ui.ColorGreen))
	fmt.Println(ui.Colorize("═══════════════════════════════════════", ui.ColorGreen))
	fmt.Println()
	fmt.Printf("%s has been added to your collection!\n", ui.Colorize(item.Name, ui.Bold))
	fmt.Printf("Remaining coins: %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.Coins), ui.ColorYellow))
	fmt.Println()
	fmt.Println("Press Enter to continue...")
	sc.scanner.Scan()

	return nil
}

// countOwnedPokemon counts how many of a specific Pokemon the player owns
func (sc *ShopCommand) countOwnedPokemon(pokemonID int) int {
	count := 0
	for _, card := range sc.gameState.Collection {
		if card.PokemonID == pokemonID {
			count++
		}
	}
	return count
}

// getNextCardID returns the next available card ID
func (sc *ShopCommand) getNextCardID() int {
	maxID := 0
	for _, card := range sc.gameState.Collection {
		if card.ID > maxID {
			maxID = card.ID
		}
	}
	return maxID + 1
}

// CheckAndRefreshShop checks if shop should be refreshed based on battles
// Should be called after each battle
func (sc *ShopCommand) CheckAndRefreshShop() error {
	sc.gameState.ShopState.BattlesSinceRefresh++

	if sc.gameState.ShopState.BattlesSinceRefresh >= 10 {
		fmt.Println()
		fmt.Println(ui.Colorize("═══════════════════════════════════════", ui.ColorBrightCyan))
		fmt.Println(ui.Colorize("  SHOP REFRESHED!", ui.Bold+ui.ColorBrightCyan))
		fmt.Println(ui.Colorize("═══════════════════════════════════════", ui.ColorBrightCyan))
		fmt.Println()
		fmt.Println("New Pokemon are now available in the shop!")
		fmt.Println()

		if err := sc.GenerateShopInventory(); err != nil {
			return fmt.Errorf("failed to refresh shop: %w", err)
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}

// getResetTimeString returns the configured token reset time
func getResetTimeString() string {
	resetTime := os.Getenv("TOKEN_RESET_TIME")
	if resetTime == "" {
		resetTime = "00:00"
	}
	return resetTime
}

// PurchaseTokens handles the token purchase flow
func (sc *ShopCommand) PurchaseTokens() error {
	// Clear screen
	sc.renderer.Clear()

	// Display header
	fmt.Println(ui.RenderLogo())
	fmt.Println()
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println(ui.Colorize("BUY GAME TOKENS", ui.Bold+ui.ColorBrightCyan))
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()

	// Get token price and daily limit from environment or use defaults
	tokenPrice := 100
	if envPrice := os.Getenv("TOKEN_PRICE"); envPrice != "" {
		if price, err := strconv.Atoi(envPrice); err == nil && price > 0 {
			tokenPrice = price
		}
	}

	dailyLimit := 10
	if envLimit := os.Getenv("DAILY_PURCHASE_LIMIT"); envLimit != "" {
		if limit, err := strconv.Atoi(envLimit); err == nil && limit > 0 {
			dailyLimit = limit
		}
	}

	// Display current status
	fmt.Printf("Your Coins: %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.Coins), ui.ColorYellow))
	fmt.Printf("Your Tokens: %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.GameTokens), ui.ColorCyan))
	fmt.Printf("Purchased Today: %d/%d\n", sc.gameState.TokensPurchasedToday, dailyLimit)
	fmt.Printf("Price: %d coins per token\n", tokenPrice)
	
	// Calculate time until reset
	resetTimeStr := getResetTimeString()
	resetHour := 0
	resetMinute := 0
	fmt.Sscanf(resetTimeStr, "%d:%d", &resetHour, &resetMinute)
	
	now := time.Now().UTC()
	nextReset := time.Date(now.Year(), now.Month(), now.Day(), resetHour, resetMinute, 0, 0, time.UTC)
	if now.After(nextReset) || now.Equal(nextReset) {
		nextReset = nextReset.Add(24 * time.Hour)
	}
	
	duration := nextReset.Sub(now)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	
	fmt.Printf("Resets at: %s UTC (in %dh %dm)\n", resetTimeStr, hours, minutes)
	fmt.Println()

	// Check if daily limit reached
	remainingPurchases := dailyLimit - sc.gameState.TokensPurchasedToday
	if remainingPurchases <= 0 {
		fmt.Println(ui.Colorize("Daily purchase limit reached!", ui.ColorRed))
		fmt.Printf("You can purchase more tokens after the daily reset at %s UTC.\n", resetTimeStr)
		fmt.Println()
		fmt.Println("Press Enter to continue...")
		sc.scanner.Scan()
		return nil
	}

	// Get quantity from user
	fmt.Printf("How many tokens would you like to buy? (1-%d): ", remainingPurchases)
	if !sc.scanner.Scan() {
		return fmt.Errorf("failed to read input")
	}

	quantityStr := strings.TrimSpace(sc.scanner.Text())
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 1 || quantity > remainingPurchases {
		return fmt.Errorf("invalid quantity. Must be between 1 and %d", remainingPurchases)
	}

	// Calculate total cost
	totalCost := quantity * tokenPrice

	// Check if user has enough coins
	if sc.gameState.Coins < totalCost {
		return fmt.Errorf("not enough coins! You need %d coins but only have %d", totalCost, sc.gameState.Coins)
	}

	// Display purchase confirmation
	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println(ui.Colorize("PURCHASE CONFIRMATION", ui.Bold+ui.ColorBrightYellow))
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()
	fmt.Printf("Tokens to buy: %d\n", quantity)
	fmt.Printf("Total cost: %s\n", ui.Colorize(fmt.Sprintf("%d coins", totalCost), ui.ColorYellow))
	fmt.Printf("You have: %s\n", ui.Colorize(fmt.Sprintf("%d coins", sc.gameState.Coins), ui.ColorYellow))
	fmt.Printf("Purchased today: %d/%d (resets at %s UTC)\n", sc.gameState.TokensPurchasedToday, dailyLimit, resetTimeStr)
	fmt.Println()

	// Get confirmation
	confirmed := ui.ConfirmationPrompt(sc.scanner, "Confirm purchase?", true)
	if !confirmed {
		fmt.Println()
		fmt.Println(ui.Colorize("Purchase cancelled.", ui.ColorYellow))
		time.Sleep(1 * time.Second)
		return nil
	}

	// Process purchase
	sc.gameState.Coins -= totalCost
	sc.gameState.GameTokens += quantity
	sc.gameState.TokensPurchasedToday += quantity

	// Save game state
	if err := storage.SaveGameState(sc.gameState); err != nil {
		// Rollback on save failure
		sc.gameState.Coins += totalCost
		sc.gameState.GameTokens -= quantity
		sc.gameState.TokensPurchasedToday -= quantity
		return fmt.Errorf("failed to save game state: %w", err)
	}

	// Display success message
	fmt.Println()
	fmt.Println(ui.Colorize("═══════════════════════════════════════", ui.ColorGreen))
	fmt.Println(ui.Colorize("  PURCHASE SUCCESSFUL!", ui.Bold+ui.ColorGreen))
	fmt.Println(ui.Colorize("═══════════════════════════════════════", ui.ColorGreen))
	fmt.Println()
	fmt.Printf("Purchased %d token(s) for %d coins\n", quantity, totalCost)
	fmt.Printf("New balance: %s | %s\n",
		ui.Colorize(fmt.Sprintf("Tokens: %d", sc.gameState.GameTokens), ui.ColorCyan),
		ui.Colorize(fmt.Sprintf("Coins: %d", sc.gameState.Coins), ui.ColorYellow))
	fmt.Println()
	fmt.Println("Press Enter to continue...")
	sc.scanner.Scan()

	return nil
}
