package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// SaveFileName is the name of the save file
	SaveFileName = "save.json"
	// SaveDirName is the directory name for game data
	SaveDirName = ".poketactix"
	// CurrentVersion is the current save file version
	CurrentVersion = "1.0.0"
)

// GetSaveFilePath returns the full path to the save file
// Returns ~/.poketactix/save.json
func GetSaveFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	saveDir := filepath.Join(homeDir, SaveDirName)
	savePath := filepath.Join(saveDir, SaveFileName)

	return savePath, nil
}

// GetSaveDirectory returns the directory path for save files
// Returns ~/.poketactix/
func GetSaveDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, SaveDirName), nil
}

// ensureSaveDirectory creates the save directory if it doesn't exist
func ensureSaveDirectory() error {
	saveDir, err := GetSaveDirectory()
	if err != nil {
		return err
	}

	// Create directory with read/write/execute permissions for user only
	if err := os.MkdirAll(saveDir, 0700); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	return nil
}

// SaveGameState serializes and writes the game state to the save file
func SaveGameState(state *GameState) error {
	if state == nil {
		return fmt.Errorf("cannot save nil game state")
	}

	// Ensure save directory exists
	if err := ensureSaveDirectory(); err != nil {
		return err
	}

	// Update last saved timestamp and version
	state.LastSaved = time.Now()
	state.Version = CurrentVersion

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	// Get save file path
	savePath, err := GetSaveFilePath()
	if err != nil {
		return err
	}

	// Write to file with read/write permissions for user only
	if err := os.WriteFile(savePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// LoadGameState reads and parses the save file
// Returns a new game state if the file doesn't exist
func LoadGameState() (*GameState, error) {
	savePath, err := GetSaveFilePath()
	if err != nil {
		return nil, err
	}

	// Check if save file exists
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		// Return nil to indicate no save file exists (not an error)
		return nil, nil
	}

	// Read save file
	data, err := os.ReadFile(savePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read save file: %w", err)
	}

	// Parse JSON
	var state GameState
	if err := json.Unmarshal(data, &state); err != nil {
		// Save file is corrupted, try to restore from backup
		return nil, fmt.Errorf("corrupted save file: %w", err)
	}

	return &state, nil
}

// CreateNewGameState creates a fresh game state for a new player
func CreateNewGameState(playerName string) *GameState {
	return &GameState{
		PlayerName:    playerName,
		Coins:         500, // Starting coins
		Collection:    []PlayerCard{},
		Deck:          []int{},
		Stats:         PlayerStats{},
		ShopState:     ShopState{
			Inventory:           []ShopItem{},
			LastRefresh:         time.Now(),
			BattlesSinceRefresh: 0,
		},
		BattleHistory: []BattleRecord{},
		LastSaved:     time.Now(),
		Version:       CurrentVersion,
	}
}

// SaveFileExists checks if a save file exists
func SaveFileExists() (bool, error) {
	savePath, err := GetSaveFilePath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(savePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
