//go:build cli
// +build cli

package commands

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
	"pokemon-cli/internal/pokemon"
)

func createTestGameState() *storage.GameState {
	state := storage.CreateNewGameState("TestPlayer")
	
	// Add some test Pokemon to collection
	state.Collection = []storage.PlayerCard{
		{
			ID:          0,
			PokemonID:   25,
			Name:        "Pikachu",
			Level:       10,
			XP:          50,
			BaseHP:      35,
			BaseAttack:  55,
			BaseDefense: 40,
			BaseSpeed:   90,
			Types:       []string{"electric"},
			Moves: []pokemon.Move{
				{Name: "Thunderbolt", Power: 90, StaminaCost: 15, Type: "electric"},
			},
			Sprite:      "pikachu.png",
			IsLegendary: false,
			IsMythical:  false,
			AcquiredAt:  time.Now(),
		},
		{
			ID:          1,
			PokemonID:   6,
			Name:        "Charizard",
			Level:       15,
			XP:          75,
			BaseHP:      78,
			BaseAttack:  84,
			BaseDefense: 78,
			BaseSpeed:   100,
			Types:       []string{"fire", "flying"},
			Moves: []pokemon.Move{
				{Name: "Flamethrower", Power: 90, StaminaCost: 15, Type: "fire"},
			},
			Sprite:      "charizard.png",
			IsLegendary: false,
			IsMythical:  false,
			AcquiredAt:  time.Now(),
		},
		{
			ID:          2,
			PokemonID:   9,
			Name:        "Blastoise",
			Level:       12,
			XP:          60,
			BaseHP:      79,
			BaseAttack:  83,
			BaseDefense: 100,
			BaseSpeed:   78,
			Types:       []string{"water"},
			Moves: []pokemon.Move{
				{Name: "Hydro Pump", Power: 110, StaminaCost: 20, Type: "water"},
			},
			Sprite:      "blastoise.png",
			IsLegendary: false,
			IsMythical:  false,
			AcquiredAt:  time.Now(),
		},
		{
			ID:          3,
			PokemonID:   3,
			Name:        "Venusaur",
			Level:       11,
			XP:          55,
			BaseHP:      80,
			BaseAttack:  82,
			BaseDefense: 83,
			BaseSpeed:   80,
			Types:       []string{"grass", "poison"},
			Moves: []pokemon.Move{
				{Name: "Solar Beam", Power: 120, StaminaCost: 25, Type: "grass"},
			},
			Sprite:      "venusaur.png",
			IsLegendary: false,
			IsMythical:  false,
			AcquiredAt:  time.Now(),
		},
		{
			ID:          4,
			PokemonID:   94,
			Name:        "Gengar",
			Level:       13,
			XP:          65,
			BaseHP:      60,
			BaseAttack:  65,
			BaseDefense: 60,
			BaseSpeed:   110,
			Types:       []string{"ghost", "poison"},
			Moves: []pokemon.Move{
				{Name: "Shadow Ball", Power: 80, StaminaCost: 15, Type: "ghost"},
			},
			Sprite:      "gengar.png",
			IsLegendary: false,
			IsMythical:  false,
			AcquiredAt:  time.Now(),
		},
	}
	
	// Set up a valid deck with 5 Pokemon
	state.Deck = []int{0, 1, 2, 3, 4}
	
	return state
}

func TestValidateDeck(t *testing.T) {
	state := createTestGameState()
	renderer := ui.NewRenderer()
	scanner := bufio.NewScanner(strings.NewReader(""))
	
	deckCmd := NewDeckCommand(state, renderer, scanner)
	
	// Test valid deck
	errors := deckCmd.validateDeck()
	if len(errors) != 0 {
		t.Errorf("Expected no validation errors for valid deck, got: %v", errors)
	}
	
	// Test deck with wrong number of Pokemon
	state.Deck = []int{0, 1, 2}
	errors = deckCmd.validateDeck()
	if len(errors) == 0 {
		t.Error("Expected validation error for deck with 3 Pokemon")
	}
	
	// Test deck with invalid card index
	state.Deck = []int{0, 1, 2, 3, 99}
	errors = deckCmd.validateDeck()
	if len(errors) == 0 {
		t.Error("Expected validation error for invalid card index")
	}
	
	// Test deck with duplicate Pokemon
	state.Deck = []int{0, 1, 2, 3, 0}
	errors = deckCmd.validateDeck()
	if len(errors) == 0 {
		t.Error("Expected validation error for duplicate Pokemon")
	}
}

func TestGetAvailablePokemon(t *testing.T) {
	state := createTestGameState()
	renderer := ui.NewRenderer()
	scanner := bufio.NewScanner(strings.NewReader(""))
	
	deckCmd := NewDeckCommand(state, renderer, scanner)
	
	// All Pokemon are in deck, so available should be empty
	available := deckCmd.getAvailablePokemon()
	if len(available) != 0 {
		t.Errorf("Expected 0 available Pokemon when all are in deck, got: %d", len(available))
	}
	
	// Remove one Pokemon from deck
	state.Deck = []int{0, 1, 2, 3}
	available = deckCmd.getAvailablePokemon()
	if len(available) != 1 {
		t.Errorf("Expected 1 available Pokemon, got: %d", len(available))
	}
	if len(available) > 0 && available[0] != 4 {
		t.Errorf("Expected available Pokemon to be index 4, got: %d", available[0])
	}
	
	// Empty deck
	state.Deck = []int{}
	available = deckCmd.getAvailablePokemon()
	if len(available) != 5 {
		t.Errorf("Expected 5 available Pokemon with empty deck, got: %d", len(available))
	}
}

func TestUndoDeckChanges(t *testing.T) {
	state := createTestGameState()
	renderer := ui.NewRenderer()
	scanner := bufio.NewScanner(strings.NewReader(""))
	
	deckCmd := NewDeckCommand(state, renderer, scanner)
	
	// Save original deck
	originalDeck := make([]int, len(state.Deck))
	copy(originalDeck, state.Deck)
	
	// Modify deck
	state.Deck = []int{4, 3, 2, 1, 0}
	
	// Undo changes
	deckCmd.UndoDeckChanges(originalDeck)
	
	// Verify deck is restored
	if len(state.Deck) != len(originalDeck) {
		t.Errorf("Deck length mismatch after undo: expected %d, got %d", len(originalDeck), len(state.Deck))
	}
	
	for i := range state.Deck {
		if state.Deck[i] != originalDeck[i] {
			t.Errorf("Deck mismatch at position %d: expected %d, got %d", i, originalDeck[i], state.Deck[i])
		}
	}
}
