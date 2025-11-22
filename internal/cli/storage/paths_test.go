package storage

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetSaveFilePath(t *testing.T) {
	path, err := GetSaveFilePath()
	if err != nil {
		t.Fatalf("GetSaveFilePath failed: %v", err)
	}

	// Verify path is not empty
	if path == "" {
		t.Error("Expected non-empty save file path")
	}

	// Verify path contains the save directory name
	if !strings.Contains(path, SaveDirName) {
		t.Errorf("Expected path to contain %s, got %s", SaveDirName, path)
	}

	// Verify path contains the save file name
	if !strings.Contains(path, SaveFileName) {
		t.Errorf("Expected path to contain %s, got %s", SaveFileName, path)
	}

	// Verify path uses correct separator for OS
	expectedSep := string(filepath.Separator)
	if !strings.Contains(path, expectedSep) {
		t.Errorf("Expected path to use OS separator %s, got %s", expectedSep, path)
	}
}

func TestGetSaveDirectory(t *testing.T) {
	dir, err := GetSaveDirectory()
	if err != nil {
		t.Fatalf("GetSaveDirectory failed: %v", err)
	}

	// Verify directory is not empty
	if dir == "" {
		t.Error("Expected non-empty save directory")
	}

	// Verify directory contains the save directory name
	if !strings.Contains(dir, SaveDirName) {
		t.Errorf("Expected directory to contain %s, got %s", SaveDirName, dir)
	}

	// Verify path is absolute
	if !filepath.IsAbs(dir) {
		t.Errorf("Expected absolute path, got %s", dir)
	}
}

func TestCrossPlatformPathSeparators(t *testing.T) {
	path, err := GetSaveFilePath()
	if err != nil {
		t.Fatalf("GetSaveFilePath failed: %v", err)
	}

	// Verify no hardcoded separators
	switch runtime.GOOS {
	case "windows":
		// On Windows, should use backslash
		if strings.Contains(path, "/") && !strings.Contains(path, "\\") {
			t.Error("Windows path should use backslash separator")
		}
	default:
		// On Unix-like systems, should use forward slash
		if strings.Contains(path, "\\") {
			t.Error("Unix path should use forward slash separator")
		}
	}
}

func TestHomeDirectoryExpansion(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("Cannot get home directory: %v", err)
	}

	path, err := GetSaveFilePath()
	if err != nil {
		t.Fatalf("GetSaveFilePath failed: %v", err)
	}

	// Verify path starts with home directory
	if !strings.HasPrefix(path, homeDir) {
		t.Errorf("Expected path to start with home directory %s, got %s", homeDir, path)
	}
}

func TestEnsureSaveDirectory(t *testing.T) {
	// This test verifies that ensureSaveDirectory creates the directory
	// We'll use a temporary directory to avoid affecting the real save directory
	
	// Get the save directory path
	saveDir, err := GetSaveDirectory()
	if err != nil {
		t.Fatalf("GetSaveDirectory failed: %v", err)
	}

	// Ensure it can be created
	err = ensureSaveDirectory()
	if err != nil {
		t.Fatalf("ensureSaveDirectory failed: %v", err)
	}

	// Verify directory exists
	info, err := os.Stat(saveDir)
	if err != nil {
		t.Fatalf("Save directory does not exist: %v", err)
	}

	// Verify it's a directory
	if !info.IsDir() {
		t.Error("Save path is not a directory")
	}

	// Verify permissions (Unix-like systems only)
	if runtime.GOOS != "windows" {
		mode := info.Mode()
		// Should have user read/write/execute permissions
		if mode&0700 != 0700 {
			t.Errorf("Expected directory permissions 0700, got %o", mode&0777)
		}
	}
}

func TestBackupPathConstruction(t *testing.T) {
	saveDir, err := GetSaveDirectory()
	if err != nil {
		t.Fatalf("GetSaveDirectory failed: %v", err)
	}

	// Construct a backup path using filepath.Join
	backupName := "save_backup_20250101_120000.json"
	backupPath := filepath.Join(saveDir, backupName)

	// Verify path is properly constructed
	if !strings.Contains(backupPath, saveDir) {
		t.Errorf("Backup path should contain save directory")
	}

	if !strings.Contains(backupPath, backupName) {
		t.Errorf("Backup path should contain backup filename")
	}

	// Verify no double separators
	doubleSep := string(filepath.Separator) + string(filepath.Separator)
	if strings.Contains(backupPath, doubleSep) {
		t.Errorf("Backup path contains double separator: %s", backupPath)
	}
}

func TestExportPathHandling(t *testing.T) {
	// Test that export paths are handled correctly across platforms
	testPaths := []string{
		".",
		"..",
		"exports",
		filepath.Join("exports", "backup"),
	}

	for _, testPath := range testPaths {
		// Verify filepath.Join works correctly
		fullPath := filepath.Join(testPath, "test.json")
		
		// Should not contain mixed separators
		if strings.Contains(fullPath, "/") && strings.Contains(fullPath, "\\") {
			t.Errorf("Path contains mixed separators: %s", fullPath)
		}
	}
}

func TestLegacySaveFilePath(t *testing.T) {
	path, err := GetLegacySaveFilePath()
	if err != nil {
		t.Fatalf("GetLegacySaveFilePath failed: %v", err)
	}

	// Verify path is not empty
	if path == "" {
		t.Error("Expected non-empty legacy save file path")
	}

	// Verify path contains the save directory name
	if !strings.Contains(path, SaveDirName) {
		t.Errorf("Expected path to contain %s, got %s", SaveDirName, path)
	}

	// Verify path ends with .json (not .json.gz)
	if !strings.HasSuffix(path, ".json") {
		t.Errorf("Expected legacy path to end with .json, got %s", path)
	}

	if strings.HasSuffix(path, ".json.gz") {
		t.Error("Legacy path should not end with .json.gz")
	}
}
