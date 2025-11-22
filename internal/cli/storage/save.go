package storage

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	// SaveFileName is the name of the save file
	SaveFileName = "save.json.gz"
	// SaveDirName is the directory name for game data
	SaveDirName = ".poketactix"
	// CurrentVersion is the current save file version
	CurrentVersion = "1.0.0"
	// MaxBattleHistory is the maximum number of battle records to keep
	MaxBattleHistory = 20
)

// GetSaveFilePath returns the full path to the save file
// Returns ~/.poketactix/save.json.gz
func GetSaveFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	saveDir := filepath.Join(homeDir, SaveDirName)
	savePath := filepath.Join(saveDir, SaveFileName)

	return savePath, nil
}

// GetLegacySaveFilePath returns the path to the old uncompressed save file
// Returns ~/.poketactix/save.json
func GetLegacySaveFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	saveDir := filepath.Join(homeDir, SaveDirName)
	savePath := filepath.Join(saveDir, "save.json")

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
// Optimized with gzip compression and battle history limiting
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

	// Limit battle history to last 20 entries
	if len(state.BattleHistory) > MaxBattleHistory {
		state.BattleHistory = state.BattleHistory[len(state.BattleHistory)-MaxBattleHistory:]
	}

	// Marshal to JSON without indentation (more compact)
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	// Get save file path
	savePath, err := GetSaveFilePath()
	if err != nil {
		return err
	}

	// Compress with gzip
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(data); err != nil {
		gzWriter.Close()
		return fmt.Errorf("failed to compress save data: %w", err)
	}
	if err := gzWriter.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// Write compressed data to file with read/write permissions for user only
	if err := os.WriteFile(savePath, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// LoadGameState reads and parses the save file
// Returns a new game state if the file doesn't exist
// Optimized with gzip decompression and legacy format support
func LoadGameState() (*GameState, error) {
	savePath, err := GetSaveFilePath()
	if err != nil {
		return nil, err
	}

	// Check if compressed save file exists
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		// Try legacy uncompressed format
		return loadLegacySaveFile()
	}

	// Read compressed save file
	compressedData, err := os.ReadFile(savePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read save file: %w", err)
	}

	// Decompress with gzip
	gzReader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		// If decompression fails, try legacy format
		return loadLegacySaveFile()
	}
	defer gzReader.Close()

	// Read decompressed data
	data, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress save file: %w", err)
	}

	// Parse JSON
	var state GameState
	if err := json.Unmarshal(data, &state); err != nil {
		// Save file is corrupted, try to restore from backup
		return nil, fmt.Errorf("corrupted save file: %w", err)
	}

	return &state, nil
}

// loadLegacySaveFile loads the old uncompressed save format
func loadLegacySaveFile() (*GameState, error) {
	legacyPath, err := GetLegacySaveFilePath()
	if err != nil {
		return nil, err
	}

	// Check if legacy save file exists
	if _, err := os.Stat(legacyPath); os.IsNotExist(err) {
		// No save file exists at all
		return nil, nil
	}

	// Read legacy save file
	data, err := os.ReadFile(legacyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read legacy save file: %w", err)
	}

	// Parse JSON
	var state GameState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("corrupted legacy save file: %w", err)
	}

	// Migrate to new compressed format
	if err := SaveGameState(&state); err != nil {
		// Log warning but don't fail - we still have the data
		fmt.Printf("Warning: Failed to migrate save file to compressed format: %v\n", err)
	} else {
		// Remove old file after successful migration
		os.Remove(legacyPath)
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
