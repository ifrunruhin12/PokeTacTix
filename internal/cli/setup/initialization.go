//go:build cli
// +build cli

package setup

import (
	"fmt"
	"math/rand"
	"time"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
	"pokemon-cli/internal/pokemon"
)

// InitializePlayerState creates the initial game state for a new player
// Includes player name, starter deck, starting coins, and initial shop inventory
func InitializePlayerState(playerName string, starterCards []storage.PlayerCard) (*storage.GameState, error) {
	// Create new game state
	gameState := storage.CreateNewGameState(playerName)

	// Add starter cards to collection with proper IDs
	for i := range starterCards {
		starterCards[i].ID = i
	}
	gameState.Collection = starterCards

	// Set deck to use all starter cards (IDs 0-4)
	gameState.Deck = []int{0, 1, 2, 3, 4}

	// Initialize stats
	gameState.Stats = storage.PlayerStats{
		TotalBattles1v1:  0,
		Wins1v1:          0,
		Losses1v1:        0,
		Draws1v1:         0,
		TotalBattles5v5:  0,
		Wins5v5:          0,
		Losses5v5:        0,
		Draws5v5:         0,
		HighestLevel:     1,
		TotalCoinsEarned: 0,
		TotalPokemon:     5, // Starting with 5 Pokemon
	}

	// Generate initial shop inventory
	shopInventory, err := generateInitialShopInventory()
	if err != nil {
		return nil, fmt.Errorf("failed to generate shop inventory: %w", err)
	}

	gameState.ShopState = storage.ShopState{
		Inventory:           shopInventory,
		LastRefresh:         time.Now(),
		BattlesSinceRefresh: 0,
	}

	return gameState, nil
}

// generateInitialShopInventory creates the initial shop inventory
// Generates 10-15 random non-legendary, non-mythical Pokemon with pricing
func generateInitialShopInventory() ([]storage.ShopItem, error) {
	// Generate 10-15 Pokemon
	count := 10 + rand.Intn(6) // 10-15
	inventory := make([]storage.ShopItem, 0, count)
	usedIDs := make(map[int]bool)

	for i := 0; i < count; i++ {
		var pokemonEntry *pokemon.PokemonEntry
		var err error

		// Keep trying until we get a unique Pokemon
		for {
			pokemonEntry, err = pokemon.GetRandomPokemon(true, true) // Exclude legendary and mythical
			if err != nil {
				return nil, fmt.Errorf("failed to generate shop Pokemon: %w", err)
			}

			if !usedIDs[pokemonEntry.ID] {
				usedIDs[pokemonEntry.ID] = true
				break
			}
		}

		// Determine rarity and price
		rarity, price := determineRarityAndPrice(pokemonEntry)

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

	return inventory, nil
}

// determineRarityAndPrice assigns rarity and price based on Pokemon stats
func determineRarityAndPrice(entry *pokemon.PokemonEntry) (string, int) {
	// Calculate total base stats
	totalStats := entry.HP + entry.Attack + entry.Defense + entry.Speed

	// Determine rarity based on total stats
	// Common: < 300, Uncommon: 300-400, Rare: > 400
	if totalStats < 300 {
		return "common", 100
	} else if totalStats < 400 {
		return "uncommon", 250
	} else {
		return "rare", 500
	}
}

// SaveInitialGameState saves the initial game state and displays success message
func SaveInitialGameState(gameState *storage.GameState) error {
	// Save game state
	if err := storage.SaveGameState(gameState); err != nil {
		return fmt.Errorf("failed to save game state: %w", err)
	}

	// Display success message
	fmt.Println()
	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("✓ Game initialized successfully!", ui.ColorBrightGreen))
	} else {
		fmt.Println("✓ Game initialized successfully!")
	}

	// Show save location
	savePath, err := storage.GetSaveFilePath()
	if err == nil {
		fmt.Printf("Your progress is saved at: %s\n", savePath)
	}

	fmt.Println()
	fmt.Printf("Starting coins: %d\n", gameState.Coins)
	fmt.Printf("Pokemon in collection: %d\n", len(gameState.Collection))
	fmt.Printf("Shop items available: %d\n", len(gameState.ShopState.Inventory))
	fmt.Println()

	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("You're all set! Type 'help' to see available commands.", ui.ColorBrightCyan))
		fmt.Println(ui.Colorize("Type 'battle' to start your first battle!", ui.ColorBrightYellow))
	} else {
		fmt.Println("You're all set! Type 'help' to see available commands.")
		fmt.Println("Type 'battle' to start your first battle!")
	}
	fmt.Println()

	return nil
}

// RunCompleteSetup runs the complete setup process for a new player
// Returns the initialized game state
func RunCompleteSetup() (*storage.GameState, error) {
	// Run onboarding to get player name
	playerName, err := RunOnboarding()
	if err != nil {
		return nil, fmt.Errorf("onboarding failed: %w", err)
	}

	// Generate starter deck
	fmt.Println()
	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("Generating your starter Pokemon...", ui.ColorBrightYellow))
	} else {
		fmt.Println("Generating your starter Pokemon...")
	}
	fmt.Println()

	starterCards, err := GenerateStarterDeck()
	if err != nil {
		return nil, fmt.Errorf("failed to generate starter deck: %w", err)
	}

	// Display starter Pokemon
	DisplayStarterPokemon(starterCards)

	// Wait for user confirmation
	if err := ConfirmStarterDeck(); err != nil {
		return nil, err
	}

	// Initialize player state
	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("\nInitializing game state...", ui.ColorBrightYellow))
	} else {
		fmt.Println("\nInitializing game state...")
	}

	gameState, err := InitializePlayerState(playerName, starterCards)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize player state: %w", err)
	}

	// Save initial game state
	if err := SaveInitialGameState(gameState); err != nil {
		return nil, err
	}

	return gameState, nil
}
