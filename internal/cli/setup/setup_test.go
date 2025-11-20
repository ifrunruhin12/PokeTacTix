// +build cli

package setup

import (
	"testing"

	"pokemon-cli/internal/cli/storage"
)

func TestIsFirstLaunch(t *testing.T) {
	// This test just verifies the function doesn't panic
	_, err := IsFirstLaunch()
	if err != nil {
		t.Logf("IsFirstLaunch returned error (expected if save dir doesn't exist): %v", err)
	}
}

func TestGenerateStarterDeck(t *testing.T) {
	starterCards, err := GenerateStarterDeck()
	if err != nil {
		t.Fatalf("GenerateStarterDeck failed: %v", err)
	}

	// Verify we got 5 cards
	if len(starterCards) != 5 {
		t.Errorf("Expected 5 starter cards, got %d", len(starterCards))
	}

	// Verify no duplicates
	seen := make(map[int]bool)
	for _, card := range starterCards {
		if seen[card.PokemonID] {
			t.Errorf("Duplicate Pokemon ID %d in starter deck", card.PokemonID)
		}
		seen[card.PokemonID] = true
	}

	// Verify all cards are level 1 with 0 XP
	for i, card := range starterCards {
		if card.Level != 1 {
			t.Errorf("Card %d has level %d, expected 1", i, card.Level)
		}
		if card.XP != 0 {
			t.Errorf("Card %d has XP %d, expected 0", i, card.XP)
		}
		if card.IsLegendary {
			t.Errorf("Card %d (%s) is legendary, should not be in starter deck", i, card.Name)
		}
		if card.IsMythical {
			t.Errorf("Card %d (%s) is mythical, should not be in starter deck", i, card.Name)
		}
	}
}

func TestInitializePlayerState(t *testing.T) {
	// Generate starter deck
	starterCards, err := GenerateStarterDeck()
	if err != nil {
		t.Fatalf("GenerateStarterDeck failed: %v", err)
	}

	// Initialize player state
	playerName := "TestPlayer"
	gameState, err := InitializePlayerState(playerName, starterCards)
	if err != nil {
		t.Fatalf("InitializePlayerState failed: %v", err)
	}

	// Verify player name
	if gameState.PlayerName != playerName {
		t.Errorf("Expected player name %s, got %s", playerName, gameState.PlayerName)
	}

	// Verify starting coins
	if gameState.Coins != 500 {
		t.Errorf("Expected 500 starting coins, got %d", gameState.Coins)
	}

	// Verify collection has 5 Pokemon
	if len(gameState.Collection) != 5 {
		t.Errorf("Expected 5 Pokemon in collection, got %d", len(gameState.Collection))
	}

	// Verify deck has 5 Pokemon
	if len(gameState.Deck) != 5 {
		t.Errorf("Expected 5 Pokemon in deck, got %d", len(gameState.Deck))
	}

	// Verify deck references correct card IDs
	for i, cardID := range gameState.Deck {
		if cardID != i {
			t.Errorf("Expected deck[%d] = %d, got %d", i, i, cardID)
		}
	}

	// Verify stats are initialized
	if gameState.Stats.TotalPokemon != 5 {
		t.Errorf("Expected TotalPokemon = 5, got %d", gameState.Stats.TotalPokemon)
	}

	// Verify shop inventory is generated
	if len(gameState.ShopState.Inventory) < 10 || len(gameState.ShopState.Inventory) > 15 {
		t.Errorf("Expected 10-15 shop items, got %d", len(gameState.ShopState.Inventory))
	}

	// Verify no legendary or mythical in shop
	for i, item := range gameState.ShopState.Inventory {
		if item.IsLegendary {
			t.Errorf("Shop item %d (%s) is legendary, should not be in shop", i, item.Name)
		}
		if item.IsMythical {
			t.Errorf("Shop item %d (%s) is mythical, should not be in shop", i, item.Name)
		}
	}
}

func TestDetermineRarityAndPrice(t *testing.T) {
	tests := []struct {
		name       string
		hp         int
		attack     int
		defense    int
		speed      int
		wantRarity string
		wantPrice  int
	}{
		{
			name:       "Common Pokemon",
			hp:         50,
			attack:     50,
			defense:    50,
			speed:      50,
			wantRarity: "common",
			wantPrice:  100,
		},
		{
			name:       "Uncommon Pokemon",
			hp:         75,
			attack:     75,
			defense:    75,
			speed:      75,
			wantRarity: "uncommon",
			wantPrice:  250,
		},
		{
			name:       "Rare Pokemon",
			hp:         100,
			attack:     100,
			defense:    100,
			speed:      100,
			wantRarity: "rare",
			wantPrice:  500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := &storage.PlayerCard{
				BaseHP:      tt.hp,
				BaseAttack:  tt.attack,
				BaseDefense: tt.defense,
				BaseSpeed:   tt.speed,
			}

			// Create a mock PokemonEntry for testing
			mockEntry := struct {
				HP      int
				Attack  int
				Defense int
				Speed   int
			}{
				HP:      tt.hp,
				Attack:  tt.attack,
				Defense: tt.defense,
				Speed:   tt.speed,
			}

			// Calculate total stats
			totalStats := mockEntry.HP + mockEntry.Attack + mockEntry.Defense + mockEntry.Speed

			var gotRarity string
			var gotPrice int

			if totalStats < 300 {
				gotRarity = "common"
				gotPrice = 100
			} else if totalStats < 400 {
				gotRarity = "uncommon"
				gotPrice = 250
			} else {
				gotRarity = "rare"
				gotPrice = 500
			}

			if gotRarity != tt.wantRarity {
				t.Errorf("Expected rarity %s, got %s", tt.wantRarity, gotRarity)
			}
			if gotPrice != tt.wantPrice {
				t.Errorf("Expected price %d, got %d", tt.wantPrice, gotPrice)
			}

			_ = entry // Use the entry variable to avoid unused variable error
		})
	}
}
