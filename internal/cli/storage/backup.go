package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	// BackupPrefix is the prefix for backup files
	BackupPrefix = "save_backup_"
	// BackupExtension is the file extension for backups
	BackupExtension = ".json"
	// MaxBackups is the maximum number of backups to keep
	MaxBackups = 3
)

// CreateBackup creates a timestamped backup of the current game state
// Maintains only the last 3 backups
func CreateBackup(state *GameState) error {
	if state == nil {
		return fmt.Errorf("cannot backup nil game state")
	}

	// Ensure save directory exists
	if err := ensureSaveDirectory(); err != nil {
		return err
	}

	saveDir, err := GetSaveDirectory()
	if err != nil {
		return err
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s%s%s", BackupPrefix, timestamp, BackupExtension)
	backupPath := filepath.Join(saveDir, backupName)

	// Marshal game state to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal game state for backup: %w", err)
	}

	// Write backup file
	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	// Clean up old backups
	if err := cleanupOldBackups(); err != nil {
		// Log error but don't fail the backup operation
		fmt.Printf("Warning: failed to cleanup old backups: %v\n", err)
	}

	return nil
}

// RestoreFromBackup attempts to restore from the most recent backup
// Returns the restored game state or an error if no backups exist
func RestoreFromBackup() (*GameState, error) {
	backups, err := listBackups()
	if err != nil {
		return nil, err
	}

	if len(backups) == 0 {
		return nil, fmt.Errorf("no backups available")
	}

	// Get the most recent backup (last in sorted list)
	mostRecentBackup := backups[len(backups)-1]

	saveDir, err := GetSaveDirectory()
	if err != nil {
		return nil, err
	}

	backupPath := filepath.Join(saveDir, mostRecentBackup)

	// Read backup file
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %w", err)
	}

	// Parse JSON
	var state GameState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse backup file: %w", err)
	}

	return &state, nil
}

// listBackups returns a sorted list of backup filenames (oldest to newest)
func listBackups() ([]string, error) {
	saveDir, err := GetSaveDirectory()
	if err != nil {
		return nil, err
	}

	// Read directory contents
	entries, err := os.ReadDir(saveDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read save directory: %w", err)
	}

	// Filter backup files
	var backups []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, BackupPrefix) && strings.HasSuffix(name, BackupExtension) {
			backups = append(backups, name)
		}
	}

	// Sort by filename (which includes timestamp)
	sort.Strings(backups)

	return backups, nil
}

// cleanupOldBackups removes old backups, keeping only the most recent MaxBackups
func cleanupOldBackups() error {
	backups, err := listBackups()
	if err != nil {
		return err
	}

	// If we have more than MaxBackups, delete the oldest ones
	if len(backups) <= MaxBackups {
		return nil
	}

	saveDir, err := GetSaveDirectory()
	if err != nil {
		return err
	}

	// Calculate how many to delete
	toDelete := len(backups) - MaxBackups

	// Delete oldest backups
	for i := 0; i < toDelete; i++ {
		backupPath := filepath.Join(saveDir, backups[i])
		if err := os.Remove(backupPath); err != nil {
			return fmt.Errorf("failed to delete old backup %s: %w", backups[i], err)
		}
	}

	return nil
}

// GetBackupCount returns the number of available backups
func GetBackupCount() (int, error) {
	backups, err := listBackups()
	if err != nil {
		return 0, err
	}
	return len(backups), nil
}

// ListAllBackups returns information about all available backups
func ListAllBackups() ([]BackupInfo, error) {
	backups, err := listBackups()
	if err != nil {
		return nil, err
	}

	saveDir, err := GetSaveDirectory()
	if err != nil {
		return nil, err
	}

	var backupInfos []BackupInfo
	for _, backup := range backups {
		backupPath := filepath.Join(saveDir, backup)
		info, err := os.Stat(backupPath)
		if err != nil {
			continue
		}

		backupInfos = append(backupInfos, BackupInfo{
			Filename:  backup,
			Size:      info.Size(),
			Timestamp: info.ModTime(),
		})
	}

	return backupInfos, nil
}

// BackupInfo contains information about a backup file
type BackupInfo struct {
	Filename  string
	Size      int64
	Timestamp time.Time
}
