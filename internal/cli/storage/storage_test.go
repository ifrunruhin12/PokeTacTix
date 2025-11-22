package storage

import (
	"os"
	"testing"
	"time"

	"pokemon-cli/internal/pokemon"
)

// setupTestEnvironment creates a temporary directory for testing
func setupTestEnvironment(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "poketactix_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	cleanup := func() {
		os.Setenv("HOME", originalHome)
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// Note: TestGetSaveFilePath is in paths_test.go

func TestCreateNewGameState(t *testing.T) {
	playerName := "TestPlayer"
	state := CreateNewGameState(playerName)

	if state.PlayerName != playerName {
		t.Errorf("Expected player name %s, got: %s", playerName, state.PlayerName)
	}

	if state.Coins != 500 {
		t.Errorf("Expected 500 starting coins, got: %d", state.Coins)
	}

	if state.Version != CurrentVersion {
		t.Errorf("Expected version %s, got: %s", CurrentVersion, state.Version)
	}

	if len(state.Collection) != 0 {
		t.Errorf("Expected empty collection, got: %d items", len(state.Collection))
	}
}

func TestSaveAndLoadGameState(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a test game state
	originalState := CreateNewGameState("TestPlayer")
	originalState.Coins = 1000
	originalState.Collection = []PlayerCard{
		{
			ID:          1,
			PokemonID:   25,
			Name:        "Pikachu",
			Level:       5,
			XP:          100,
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
	}
	originalState.Deck = []int{0}

	// Save the state
	err := SaveGameState(originalState)
	if err != nil {
		t.Fatalf("SaveGameState failed: %v", err)
	}

	// Verify file exists
	exists, err := SaveFileExists()
	if err != nil {
		t.Fatalf("SaveFileExists failed: %v", err)
	}
	if !exists {
		t.Fatal("Save file should exist after saving")
	}

	// Load the state
	loadedState, err := LoadGameState()
	if err != nil {
		t.Fatalf("LoadGameState failed: %v", err)
	}

	// Verify loaded state matches original
	if loadedState.PlayerName != originalState.PlayerName {
		t.Errorf("Player name mismatch: expected %s, got %s", originalState.PlayerName, loadedState.PlayerName)
	}

	if loadedState.Coins != originalState.Coins {
		t.Errorf("Coins mismatch: expected %d, got %d", originalState.Coins, loadedState.Coins)
	}

	if len(loadedState.Collection) != len(originalState.Collection) {
		t.Errorf("Collection size mismatch: expected %d, got %d", len(originalState.Collection), len(loadedState.Collection))
	}

	if len(loadedState.Collection) > 0 {
		if loadedState.Collection[0].Name != "Pikachu" {
			t.Errorf("Expected Pikachu in collection, got: %s", loadedState.Collection[0].Name)
		}
	}
}

func TestLoadGameStateNoFile(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Try to load when no save file exists
	state, err := LoadGameState()
	if err != nil {
		t.Fatalf("LoadGameState should not error when file doesn't exist: %v", err)
	}

	if state != nil {
		t.Error("Expected nil state when no save file exists")
	}
}

func TestCreateBackup(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	state := CreateNewGameState("TestPlayer")
	state.Coins = 2000

	// Create backup
	err := CreateBackup(state)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Verify backup was created
	count, err := GetBackupCount()
	if err != nil {
		t.Fatalf("GetBackupCount failed: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 backup, got: %d", count)
	}
}

func TestRestoreFromBackup(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create and save original state
	originalState := CreateNewGameState("TestPlayer")
	originalState.Coins = 3000

	err := CreateBackup(originalState)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Restore from backup
	restoredState, err := RestoreFromBackup()
	if err != nil {
		t.Fatalf("RestoreFromBackup failed: %v", err)
	}

	if restoredState.Coins != originalState.Coins {
		t.Errorf("Coins mismatch after restore: expected %d, got %d", originalState.Coins, restoredState.Coins)
	}
}

func TestBackupCleanup(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	state := CreateNewGameState("TestPlayer")

	// Create more than MaxBackups backups
	for i := 0; i < MaxBackups+2; i++ {
		state.Coins = 1000 + i
		err := CreateBackup(state)
		if err != nil {
			t.Fatalf("CreateBackup %d failed: %v", i, err)
		}
		// Delay to ensure different timestamps (backup uses seconds precision)
		time.Sleep(1100 * time.Millisecond)
	}

	// Verify only MaxBackups remain
	count, err := GetBackupCount()
	if err != nil {
		t.Fatalf("GetBackupCount failed: %v", err)
	}

	if count != MaxBackups {
		t.Errorf("Expected %d backups after cleanup, got: %d", MaxBackups, count)
	}
}

func TestAutoSave(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	state := CreateNewGameState("TestPlayer")
	state.Coins = 1500

	opts := AutoSaveOptions{
		CreateBackup:     false,
		ShowConfirmation: false,
	}

	err := AutoSave(state, opts)
	if err != nil {
		t.Fatalf("AutoSave failed: %v", err)
	}

	// Verify save was created
	exists, err := SaveFileExists()
	if err != nil {
		t.Fatalf("SaveFileExists failed: %v", err)
	}

	if !exists {
		t.Error("Save file should exist after AutoSave")
	}
}

func TestAutoSaveWithBackup(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	state := CreateNewGameState("TestPlayer")
	state.Coins = 2500

	err := AutoSaveWithBackup(state, "Test save")
	if err != nil {
		t.Fatalf("AutoSaveWithBackup failed: %v", err)
	}

	// Verify backup was created
	count, err := GetBackupCount()
	if err != nil {
		t.Fatalf("GetBackupCount failed: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 backup after AutoSaveWithBackup, got: %d", count)
	}
}

func TestPlayerCardGetCurrentStats(t *testing.T) {
	card := PlayerCard{
		BaseHP:      100,
		BaseAttack:  50,
		BaseDefense: 40,
		BaseSpeed:   60,
		Level:       1,
	}

	stats := card.GetCurrentStats()

	// At level 1, stats should equal base stats
	if stats.HP != 100 {
		t.Errorf("Expected HP 100 at level 1, got: %d", stats.HP)
	}

	// Test level scaling
	card.Level = 10
	stats = card.GetCurrentStats()

	// HP should increase by 3% per level
	expectedHP := int(float64(100) * (1.0 + float64(9)*0.03))
	if stats.HP != expectedHP {
		t.Errorf("Expected HP %d at level 10, got: %d", expectedHP, stats.HP)
	}
}

func TestPlayerCardToCard(t *testing.T) {
	playerCard := PlayerCard{
		ID:          1,
		PokemonID:   25,
		Name:        "Pikachu",
		Level:       5,
		XP:          100,
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
	}

	card := playerCard.ToCard()

	if card.Name != playerCard.Name {
		t.Errorf("Name mismatch: expected %s, got %s", playerCard.Name, card.Name)
	}

	if card.Level != playerCard.Level {
		t.Errorf("Level mismatch: expected %d, got %d", playerCard.Level, card.Level)
	}

	if len(card.Moves) != len(playerCard.Moves) {
		t.Errorf("Moves count mismatch: expected %d, got %d", len(playerCard.Moves), len(card.Moves))
	}

	// Verify stats are calculated correctly
	stats := playerCard.GetCurrentStats()
	if card.HP != stats.HP {
		t.Errorf("HP mismatch: expected %d, got %d", stats.HP, card.HP)
	}
}
