package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ExportSave exports the current save file to a specified location
// Returns the path where the file was exported
func ExportSave(destinationPath string) (string, error) {
	// Get current save file path
	savePath, err := GetSaveFilePath()
	if err != nil {
		return "", fmt.Errorf("failed to get save file path: %w", err)
	}

	// Check if save file exists
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		return "", fmt.Errorf("no save file found to export")
	}

	// If destination is a directory, create a filename
	fileInfo, err := os.Stat(destinationPath)
	if err == nil && fileInfo.IsDir() {
		// Create filename with timestamp
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("poketactix_save_%s.json", timestamp)
		destinationPath = filepath.Join(destinationPath, filename)
	}

	// Ensure destination directory exists
	destDir := filepath.Dir(destinationPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Open source file
	sourceFile, err := os.Open(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to open save file: %w", err)
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(destinationPath)
	if err != nil {
		return "", fmt.Errorf("failed to create export file: %w", err)
	}
	defer destFile.Close()

	// Copy file contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return "", fmt.Errorf("failed to copy save file: %w", err)
	}

	return destinationPath, nil
}

// ImportSave imports a save file from a specified location
// Validates the file format before importing
func ImportSave(sourcePath string) error {
	// Check if source file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("import file not found: %s", sourcePath)
	}

	// Read and validate the import file
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	// Validate JSON format
	var testState GameState
	if err := json.Unmarshal(data, &testState); err != nil {
		return fmt.Errorf("invalid save file format: %w", err)
	}

	// Validate required fields
	if testState.PlayerName == "" {
		return fmt.Errorf("invalid save file: missing player name")
	}
	if testState.Version == "" {
		return fmt.Errorf("invalid save file: missing version")
	}

	// Create backup of current save before importing
	savePath, err := GetSaveFilePath()
	if err != nil {
		return fmt.Errorf("failed to get save file path: %w", err)
	}
	
	if _, err := os.Stat(savePath); err == nil {
		// Current save exists, back it up
		backupPath := savePath + ".before_import"
		if err := copyFile(savePath, backupPath); err != nil {
			return fmt.Errorf("failed to backup current save: %w", err)
		}
	}

	// Copy import file to save location
	if err := copyFile(sourcePath, savePath); err != nil {
		return fmt.Errorf("failed to import save file: %w", err)
	}

	return nil
}

// copyFile copies a file from source to destination
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// ValidateSaveFile validates a save file without importing it
func ValidateSaveFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Validate JSON format
	var testState GameState
	if err := json.Unmarshal(data, &testState); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Validate required fields
	if testState.PlayerName == "" {
		return fmt.Errorf("missing player name")
	}
	if testState.Version == "" {
		return fmt.Errorf("missing version")
	}
	if len(testState.Collection) == 0 {
		return fmt.Errorf("empty collection")
	}

	return nil
}
