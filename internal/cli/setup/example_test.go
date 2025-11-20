// +build cli

package setup

import (
	"fmt"
)

// ExampleGenerateStarterDeck demonstrates how to generate a starter deck
func ExampleGenerateStarterDeck() {
	starterCards, err := GenerateStarterDeck()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Generated %d starter Pokemon\n", len(starterCards))
	for i, card := range starterCards {
		fmt.Printf("%d. %s (Level %d)\n", i+1, card.Name, card.Level)
	}
}

// ExampleInitializePlayerState demonstrates how to initialize a new player
func ExampleInitializePlayerState() {
	// Generate starter deck
	starterCards, err := GenerateStarterDeck()
	if err != nil {
		fmt.Printf("Error generating starter deck: %v\n", err)
		return
	}

	// Initialize player state
	gameState, err := InitializePlayerState("Ash", starterCards)
	if err != nil {
		fmt.Printf("Error initializing player state: %v\n", err)
		return
	}

	fmt.Printf("Player: %s\n", gameState.PlayerName)
	fmt.Printf("Starting Coins: %d\n", gameState.Coins)
	fmt.Printf("Pokemon in Collection: %d\n", len(gameState.Collection))
	fmt.Printf("Shop Items: %d\n", len(gameState.ShopState.Inventory))
}
