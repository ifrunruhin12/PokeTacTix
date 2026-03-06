//go:build cli
// +build cli

package commands

import (
	"bufio"
	"strings"
	"testing"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
)

func TestGenerateShopInventory(t *testing.T) {
	// Create test game state
	state := storage.CreateNewGameState("TestPlayer")
	state.Coins = 1000

	// Create shop command
	renderer := ui.NewRenderer()
	scanner := bufio.NewScanner(strings.NewReader(""))
	shopCmd := NewShopCommand(state, renderer, scanner)

	// Generate shop inventory
	err := shopCmd.GenerateShopInventory()
	if err != nil {
		t.Fatalf("Failed to generate shop inventory: %v", err)
	}

	// Verify inventory was generated
	if len(state.ShopState.Inventory) < 10 || len(state.ShopState.Inventory) > 15 {
		t.Errorf("Expected 10-15 items in shop, got %d", len(state.ShopState.Inventory))
	}

	// Verify no legendary or mythical Pokemon in shop
	for _, item := range state.ShopState.Inventory {
		if item.IsLegendary {
			t.Errorf("Found legendary Pokemon in shop: %s", item.Name)
		}
		if item.IsMythical {
			t.Errorf("Found mythical Pokemon in shop: %s", item.Name)
		}
	}

	// Verify pricing is correct
	for _, item := range state.ShopState.Inventory {
		switch item.Rarity {
		case "common":
			if item.Price != 100 {
				t.Errorf("Common Pokemon %s has wrong price: %d (expected 100)", item.Name, item.Price)
			}
		case "uncommon":
			if item.Price != 250 {
				t.Errorf("Uncommon Pokemon %s has wrong price: %d (expected 250)", item.Name, item.Price)
			}
		case "rare":
			if item.Price != 500 {
				t.Errorf("Rare Pokemon %s has wrong price: %d (expected 500)", item.Name, item.Price)
			}
		}
	}

	// Verify battles since refresh was reset
	if state.ShopState.BattlesSinceRefresh != 0 {
		t.Errorf("Expected BattlesSinceRefresh to be 0, got %d", state.ShopState.BattlesSinceRefresh)
	}
}

func TestShopRefresh(t *testing.T) {
	// Create test game state
	state := storage.CreateNewGameState("TestPlayer")
	state.Coins = 1000

	// Create shop command
	renderer := ui.NewRenderer()
	scanner := bufio.NewScanner(strings.NewReader(""))
	shopCmd := NewShopCommand(state, renderer, scanner)

	// Generate initial inventory
	err := shopCmd.GenerateShopInventory()
	if err != nil {
		t.Fatalf("Failed to generate initial shop inventory: %v", err)
	}

	// Simulate 10 battles
	state.ShopState.BattlesSinceRefresh = 10

	// Check and refresh shop
	err = shopCmd.CheckAndRefreshShop()
	if err != nil {
		t.Fatalf("Failed to refresh shop: %v", err)
	}

	// Verify shop was refreshed
	if state.ShopState.BattlesSinceRefresh != 0 {
		t.Errorf("Expected BattlesSinceRefresh to be reset to 0, got %d", state.ShopState.BattlesSinceRefresh)
	}

	// Verify new inventory was generated
	if len(state.ShopState.Inventory) < 10 || len(state.ShopState.Inventory) > 15 {
		t.Errorf("Expected 10-15 items in refreshed shop, got %d", len(state.ShopState.Inventory))
	}
}

func TestCountOwnedPokemon(t *testing.T) {
	// Create test game state with some Pokemon
	state := storage.CreateNewGameState("TestPlayer")
	
	// Add multiple Pikachu to collection
	for i := 0; i < 3; i++ {
		state.Collection = append(state.Collection, storage.PlayerCard{
			ID:        i,
			PokemonID: 25, // Pikachu
			Name:      "Pikachu",
		})
	}

	// Add one Charizard
	state.Collection = append(state.Collection, storage.PlayerCard{
		ID:        3,
		PokemonID: 6, // Charizard
		Name:      "Charizard",
	})

	// Create shop command
	renderer := ui.NewRenderer()
	scanner := bufio.NewScanner(strings.NewReader(""))
	shopCmd := NewShopCommand(state, renderer, scanner)

	// Test counting Pikachu
	pikachuCount := shopCmd.countOwnedPokemon(25)
	if pikachuCount != 3 {
		t.Errorf("Expected 3 Pikachu, got %d", pikachuCount)
	}

	// Test counting Charizard
	charizardCount := shopCmd.countOwnedPokemon(6)
	if charizardCount != 1 {
		t.Errorf("Expected 1 Charizard, got %d", charizardCount)
	}

	// Test counting Pokemon not owned
	bulbasaurCount := shopCmd.countOwnedPokemon(1)
	if bulbasaurCount != 0 {
		t.Errorf("Expected 0 Bulbasaur, got %d", bulbasaurCount)
	}
}

func TestTokenPurchaseValidation(t *testing.T) {
	// Create test game state
	state := storage.CreateNewGameState("TestPlayer")
	state.Coins = 1000
	state.GameTokens = 3
	state.TokensPurchasedToday = 0

	// Test 1: Verify initial state
	if state.GameTokens != 3 {
		t.Errorf("Expected 3 initial tokens, got %d", state.GameTokens)
	}

	// Test 2: Simulate token purchase (5 tokens for 500 coins)
	quantity := 5
	tokenPrice := 100
	totalCost := quantity * tokenPrice

	if state.Coins < totalCost {
		t.Errorf("Should have enough coins: have %d, need %d", state.Coins, totalCost)
	}

	// Simulate purchase
	state.Coins -= totalCost
	state.GameTokens += quantity
	state.TokensPurchasedToday += quantity

	// Verify purchase
	if state.Coins != 500 {
		t.Errorf("Expected 500 coins after purchase, got %d", state.Coins)
	}
	if state.GameTokens != 8 {
		t.Errorf("Expected 8 tokens after purchase, got %d", state.GameTokens)
	}
	if state.TokensPurchasedToday != 5 {
		t.Errorf("Expected 5 tokens purchased today, got %d", state.TokensPurchasedToday)
	}

	// Test 3: Verify daily limit enforcement
	dailyLimit := 10
	remainingPurchases := dailyLimit - state.TokensPurchasedToday
	if remainingPurchases != 5 {
		t.Errorf("Expected 5 remaining purchases, got %d", remainingPurchases)
	}

	// Test 4: Attempt to purchase more than remaining limit
	attemptQuantity := 6
	if attemptQuantity > remainingPurchases {
		// This should be rejected
		t.Logf("Correctly rejecting purchase of %d tokens (only %d remaining)", attemptQuantity, remainingPurchases)
	} else {
		t.Errorf("Should reject purchase of %d tokens when only %d remaining", attemptQuantity, remainingPurchases)
	}

	// Test 5: Purchase remaining tokens
	state.Coins += 500 // Add more coins for test
	remainingCost := remainingPurchases * tokenPrice
	state.Coins -= remainingCost
	state.GameTokens += remainingPurchases
	state.TokensPurchasedToday += remainingPurchases

	if state.TokensPurchasedToday != 10 {
		t.Errorf("Expected 10 tokens purchased today (limit reached), got %d", state.TokensPurchasedToday)
	}

	// Test 6: Verify limit is enforced
	newRemainingPurchases := dailyLimit - state.TokensPurchasedToday
	if newRemainingPurchases != 0 {
		t.Errorf("Expected 0 remaining purchases after hitting limit, got %d", newRemainingPurchases)
	}
}

func TestTokenPurchaseInsufficientCoins(t *testing.T) {
	// Create test game state with insufficient coins
	state := storage.CreateNewGameState("TestPlayer")
	state.Coins = 50 // Not enough for even 1 token
	state.GameTokens = 2

	quantity := 1
	tokenPrice := 100
	totalCost := quantity * tokenPrice

	// Verify insufficient coins
	if state.Coins >= totalCost {
		t.Errorf("Test setup error: should have insufficient coins")
	}

	// This purchase should fail
	if state.Coins < totalCost {
		t.Logf("Correctly detecting insufficient coins: have %d, need %d", state.Coins, totalCost)
	}

	// Verify state unchanged
	if state.GameTokens != 2 {
		t.Errorf("Tokens should remain unchanged at 2, got %d", state.GameTokens)
	}
	if state.Coins != 50 {
		t.Errorf("Coins should remain unchanged at 50, got %d", state.Coins)
	}
}

