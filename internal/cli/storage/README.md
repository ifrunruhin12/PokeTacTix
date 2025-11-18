# CLI Storage Package

This package provides local file persistence for the PokeTacTix CLI game. It handles saving and loading game state, creating backups, and auto-saving functionality.

## Features

- **Local Save Files**: Game state saved to `~/.poketactix/save.json`
- **Automatic Backups**: Maintains last 3 backups with timestamps
- **Auto-Save**: Convenient functions for auto-saving after key game events
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Data Structures

### GameState
The main game state structure containing all player data:
- Player name and coins
- Pokemon collection
- Current deck configuration
- Battle statistics
- Shop state
- Battle history

### PlayerCard
Represents a Pokemon card owned by the player with:
- Base stats and current level
- XP and progression
- Moves and types
- Rarity flags

### PlayerStats
Tracks battle statistics:
- Wins/losses/draws for 1v1 and 5v5 modes
- Total coins earned
- Highest level achieved
- Total Pokemon collected

### ShopState
Manages shop inventory:
- Current inventory items
- Last refresh timestamp
- Battles since last refresh

## Usage Examples

### Creating a New Game

```go
import "pokemon-cli/internal/cli/storage"

// Create new game state
state := storage.CreateNewGameState("PlayerName")

// Save it
err := storage.SaveGameState(state)
if err != nil {
    log.Fatal(err)
}
```

### Loading Existing Game

```go
// Load game state
state, err := storage.LoadGameState()
if err != nil {
    log.Fatal(err)
}

// Check if save file exists
if state == nil {
    // No save file, create new game
    state = storage.CreateNewGameState("PlayerName")
}
```

### Auto-Save After Battle

```go
// After battle completes
err := storage.AutoSaveAfterBattle(state, "5v5", "win")
// Output: "Battle complete (5v5 - win). Game saved."
```

### Auto-Save After Deck Change

```go
// After modifying deck
err := storage.AutoSaveAfterDeckChange(state)
// Creates backup and saves
// Output: "Deck updated and saved."
```

### Auto-Save After Purchase

```go
// After buying a Pokemon
err := storage.AutoSaveAfterPurchase(state, "Pikachu")
// Creates backup and saves
// Output: "Purchased Pikachu. Game saved."
```

### Creating Manual Backup

```go
// Create backup before risky operation
err := storage.CreateBackup(state)
if err != nil {
    log.Printf("Failed to create backup: %v", err)
}
```

### Restoring from Backup

```go
// If save file is corrupted
state, err := storage.RestoreFromBackup()
if err != nil {
    log.Fatal("No backups available")
}

// Save restored state
storage.SaveGameState(state)
```

### Custom Auto-Save

```go
// Custom auto-save with specific options
opts := storage.AutoSaveOptions{
    CreateBackup:     true,
    ShowConfirmation: true,
    ConfirmationMsg:  "Custom save complete!",
}

err := storage.AutoSave(state, opts)
```

### Quiet Save (No Confirmation)

```go
// Save without showing confirmation message
err := storage.QuietAutoSave(state)
```

## File Locations

- **Save File**: `~/.poketactix/save.json`
- **Backups**: `~/.poketactix/save_backup_YYYYMMDD_HHMMSS.json`

## Backup Management

The package automatically:
- Creates timestamped backups when requested
- Maintains only the last 3 backups
- Cleans up old backups automatically
- Provides backup restoration functionality

## Error Handling

All functions return errors that should be checked:

```go
if err := storage.SaveGameState(state); err != nil {
    // Handle error - maybe try backup?
    log.Printf("Save failed: %v", err)
    
    // Try to restore from backup
    if restored, err := storage.RestoreFromBackup(); err == nil {
        state = restored
    }
}
```

## Testing

Run tests with:

```bash
go test -v ./internal/cli/storage
```

Tests cover:
- Save/load operations
- Backup creation and restoration
- Backup cleanup (max 3 backups)
- Auto-save functionality
- PlayerCard stat calculations
- Cross-platform file paths

## Thread Safety

This package is **not thread-safe**. If you need concurrent access to game state, implement your own synchronization.

## Performance

- Save operations: < 100ms
- Load operations: < 100ms
- Backup operations: < 150ms
- Save file size: typically < 1MB

## Version Compatibility

The package includes version tracking in the save file. Future versions can implement migration logic if the save format changes.

Current version: `1.0.0`
