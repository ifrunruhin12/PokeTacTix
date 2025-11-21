package commands

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
)

// TestNewCommandHandler tests the creation of a new command handler
func TestNewCommandHandler(t *testing.T) {
	gameState := createTestGameState()
	renderer := ui.NewRenderer()
	scanner := bufio.NewScanner(strings.NewReader(""))

	handler := NewCommandHandler(gameState, renderer, scanner)

	if handler == nil {
		t.Fatal("NewCommandHandler returned nil")
	}

	if handler.gameState != gameState {
		t.Error("CommandHandler gameState not set correctly")
	}

	if handler.renderer != renderer {
		t.Error("CommandHandler renderer not set correctly")
	}

	if handler.scanner != scanner {
		t.Error("CommandHandler scanner not set correctly")
	}

	// Check that sub-handlers are initialized
	if handler.battleCmd == nil {
		t.Error("BattleCommand not initialized")
	}

	if handler.collectionCmd == nil {
		t.Error("CollectionCommand not initialized")
	}

	if handler.deckCmd == nil {
		t.Error("DeckCommand not initialized")
	}

	if handler.shopCmd == nil {
		t.Error("ShopCommand not initialized")
	}

	if handler.statsCmd == nil {
		t.Error("StatsCommand not initialized")
	}
}

// TestHandleUnknownCommand tests handling of unknown commands
func TestHandleUnknownCommand(t *testing.T) {
	gameState := createTestGameState()
	renderer := ui.NewRenderer()
	scanner := bufio.NewScanner(strings.NewReader("\n"))

	handler := NewCommandHandler(gameState, renderer, scanner)

	// Test unknown command
	err := handler.handleUnknownCommand("invalidcommand")
	if err != nil {
		t.Errorf("handleUnknownCommand returned error: %v", err)
	}
}

// TestGetContextualHints tests the contextual hints system
func TestGetContextualHints(t *testing.T) {
	gameState := createTestGameState()
	renderer := ui.NewRenderer()
	scanner := bufio.NewScanner(strings.NewReader(""))

	handler := NewCommandHandler(gameState, renderer, scanner)

	// Test with new player (no battles)
	gameState.Stats.TotalBattles1v1 = 0
	gameState.Stats.TotalBattles5v5 = 0
	hints := handler.getContextualHints()

	if len(hints) == 0 {
		t.Error("Expected hints for new player, got none")
	}

	// Test with incomplete deck
	gameState.Deck = []int{0, 1} // Only 2 Pokemon
	hints = handler.getContextualHints()

	foundDeckHint := false
	for _, hint := range hints {
		if strings.Contains(hint, "deck") {
			foundDeckHint = true
			break
		}
	}

	if !foundDeckHint {
		t.Error("Expected deck hint for incomplete deck")
	}

	// Test with enough coins
	gameState.Coins = 500
	hints = handler.getContextualHints()

	foundCoinsHint := false
	for _, hint := range hints {
		if strings.Contains(hint, "coins") || strings.Contains(hint, "shop") {
			foundCoinsHint = true
			break
		}
	}

	if !foundCoinsHint {
		t.Error("Expected coins/shop hint when player has 500+ coins")
	}
}

// TestShowHint tests the hint display functionality
func TestShowHint(t *testing.T) {
	gameState := createTestGameState()
	renderer := ui.NewRenderer()
	scanner := bufio.NewScanner(strings.NewReader(""))

	handler := NewCommandHandler(gameState, renderer, scanner)

	// Test with new player
	gameState.Stats.TotalBattles1v1 = 0
	gameState.Stats.TotalBattles5v5 = 0

	hint := handler.ShowHint()

	if hint == "" {
		t.Error("Expected hint for new player, got empty string")
	}

	if !strings.Contains(hint, "Hint") && !strings.Contains(hint, "💡") {
		t.Error("Hint should contain 'Hint' or hint emoji")
	}
}

// createTestGameState creates a test game state for testing
func createTestGameState() *storage.GameState {
	return &storage.GameState{
		PlayerName: "TestPlayer",
		Coins:      100,
		Collection: []storage.PlayerCard{
			{
				ID:          0,
				PokemonID:   25,
				Name:        "Pikachu",
				Level:       5,
				XP:          0,
				BaseHP:      35,
				BaseAttack:  55,
				BaseDefense: 40,
				BaseSpeed:   90,
				Types:       []string{"electric"},
				AcquiredAt:  time.Now(),
			},
			{
				ID:          1,
				PokemonID:   1,
				Name:        "Bulbasaur",
				Level:       5,
				XP:          0,
				BaseHP:      45,
				BaseAttack:  49,
				BaseDefense: 49,
				BaseSpeed:   45,
				Types:       []string{"grass", "poison"},
				AcquiredAt:  time.Now(),
			},
			{
				ID:          2,
				PokemonID:   4,
				Name:        "Charmander",
				Level:       5,
				XP:          0,
				BaseHP:      39,
				BaseAttack:  52,
				BaseDefense: 43,
				BaseSpeed:   65,
				Types:       []string{"fire"},
				AcquiredAt:  time.Now(),
			},
			{
				ID:          3,
				PokemonID:   7,
				Name:        "Squirtle",
				Level:       5,
				XP:          0,
				BaseHP:      44,
				BaseAttack:  48,
				BaseDefense: 65,
				BaseSpeed:   43,
				Types:       []string{"water"},
				AcquiredAt:  time.Now(),
			},
			{
				ID:          4,
				PokemonID:   133,
				Name:        "Eevee",
				Level:       5,
				XP:          0,
				BaseHP:      55,
				BaseAttack:  55,
				BaseDefense: 50,
				BaseSpeed:   55,
				Types:       []string{"normal"},
				AcquiredAt:  time.Now(),
			},
		},
		Deck: []int{0, 1, 2, 3, 4},
		Stats: storage.PlayerStats{
			TotalBattles1v1:  0,
			Wins1v1:          0,
			Losses1v1:        0,
			Draws1v1:         0,
			TotalBattles5v5:  0,
			Wins5v5:          0,
			Losses5v5:        0,
			Draws5v5:         0,
			HighestLevel:     5,
			TotalCoinsEarned: 0,
			TotalPokemon:     5,
		},
		ShopState: storage.ShopState{
			Inventory:           []storage.ShopItem{},
			LastRefresh:         time.Now(),
			BattlesSinceRefresh: 0,
		},
		BattleHistory: []storage.BattleRecord{},
		LastSaved:     time.Now(),
		Version:       "0.1.0-test",
	}
}
