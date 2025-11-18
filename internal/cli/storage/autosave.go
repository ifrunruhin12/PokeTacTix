package storage

import (
	"fmt"
)

// AutoSaveOptions configures auto-save behavior
type AutoSaveOptions struct {
	CreateBackup      bool
	ShowConfirmation  bool
	ConfirmationMsg   string
}

// DefaultAutoSaveOptions returns the default auto-save configuration
func DefaultAutoSaveOptions() AutoSaveOptions {
	return AutoSaveOptions{
		CreateBackup:     false,
		ShowConfirmation: true,
		ConfirmationMsg:  "Game saved successfully.",
	}
}

// AutoSave saves the game state with optional backup and confirmation message
func AutoSave(state *GameState, opts AutoSaveOptions) error {
	if state == nil {
		return fmt.Errorf("cannot auto-save nil game state")
	}

	// Create backup if requested
	if opts.CreateBackup {
		if err := CreateBackup(state); err != nil {
			// Log warning but continue with save
			fmt.Printf("Warning: failed to create backup: %v\n", err)
		}
	}

	// Save game state
	if err := SaveGameState(state); err != nil {
		return fmt.Errorf("auto-save failed: %w", err)
	}

	// Show confirmation if requested
	if opts.ShowConfirmation && opts.ConfirmationMsg != "" {
		fmt.Println(opts.ConfirmationMsg)
	}

	return nil
}

// AutoSaveAfterBattle saves the game state after a battle completes
func AutoSaveAfterBattle(state *GameState, mode string, result string) error {
	opts := AutoSaveOptions{
		CreateBackup:     false,
		ShowConfirmation: true,
		ConfirmationMsg:  fmt.Sprintf("Battle complete (%s - %s). Game saved.", mode, result),
	}
	return AutoSave(state, opts)
}

// AutoSaveAfterDeckChange saves the game state after deck modifications
func AutoSaveAfterDeckChange(state *GameState) error {
	opts := AutoSaveOptions{
		CreateBackup:     true, // Create backup before deck changes
		ShowConfirmation: true,
		ConfirmationMsg:  "Deck updated and saved.",
	}
	return AutoSave(state, opts)
}

// AutoSaveAfterPurchase saves the game state after a shop purchase
func AutoSaveAfterPurchase(state *GameState, pokemonName string) error {
	opts := AutoSaveOptions{
		CreateBackup:     true, // Create backup before purchases
		ShowConfirmation: true,
		ConfirmationMsg:  fmt.Sprintf("Purchased %s. Game saved.", pokemonName),
	}
	return AutoSave(state, opts)
}

// AutoSaveAfterPokemonSelection saves the game state after selecting a Pokemon reward
func AutoSaveAfterPokemonSelection(state *GameState, pokemonName string) error {
	opts := AutoSaveOptions{
		CreateBackup:     false,
		ShowConfirmation: true,
		ConfirmationMsg:  fmt.Sprintf("Added %s to collection. Game saved.", pokemonName),
	}
	return AutoSave(state, opts)
}

// AutoSaveWithBackup saves the game state and creates a backup
func AutoSaveWithBackup(state *GameState, message string) error {
	opts := AutoSaveOptions{
		CreateBackup:     true,
		ShowConfirmation: true,
		ConfirmationMsg:  message,
	}
	return AutoSave(state, opts)
}

// QuietAutoSave saves the game state without showing a confirmation message
func QuietAutoSave(state *GameState) error {
	opts := AutoSaveOptions{
		CreateBackup:     false,
		ShowConfirmation: false,
	}
	return AutoSave(state, opts)
}
