package ui

import (
	"os"
	"testing"
)

func TestGetTerminalSize(t *testing.T) {
	width, height := getTerminalSize()
	
	// Should return valid dimensions
	if width <= 0 || height <= 0 {
		t.Errorf("Expected positive dimensions, got %dx%d", width, height)
	}
	
	// Should at least return default values
	if width < 80 || height < 24 {
		t.Logf("Terminal size %dx%d is smaller than default 80x24", width, height)
	}
}

func TestGetTerminalSizeWithEnvVars(t *testing.T) {
	// Save original env vars
	origCols := os.Getenv("COLUMNS")
	origLines := os.Getenv("LINES")
	defer func() {
		os.Setenv("COLUMNS", origCols)
		os.Setenv("LINES", origLines)
	}()
	
	// Set custom dimensions
	os.Setenv("COLUMNS", "100")
	os.Setenv("LINES", "30")
	
	width, height := getTerminalSize()
	
	// Should respect environment variables or return defaults
	if width <= 0 || height <= 0 {
		t.Errorf("Expected positive dimensions, got %dx%d", width, height)
	}
}

func TestRendererDimensions(t *testing.T) {
	renderer := NewRenderer()
	
	width, height := renderer.GetDimensions()
	
	if width <= 0 || height <= 0 {
		t.Errorf("Expected positive dimensions, got %dx%d", width, height)
	}
}

func TestIsTerminalTooSmall(t *testing.T) {
	renderer := NewRenderer()
	
	// Just verify the method works without error
	_ = renderer.IsTerminalTooSmall()
}

func TestGetAdaptiveWidth(t *testing.T) {
	renderer := NewRenderer()
	
	width := renderer.GetAdaptiveWidth()
	
	if width <= 0 {
		t.Errorf("Expected positive adaptive width, got %d", width)
	}
	
	// Should not exceed actual terminal width
	actualWidth, _ := renderer.GetDimensions()
	if width > actualWidth {
		t.Errorf("Adaptive width %d exceeds actual width %d", width, actualWidth)
	}
}

func TestGetAdaptiveHeight(t *testing.T) {
	renderer := NewRenderer()
	
	height := renderer.GetAdaptiveHeight()
	
	if height <= 0 {
		t.Errorf("Expected positive adaptive height, got %d", height)
	}
	
	// Should not exceed actual terminal height
	_, actualHeight := renderer.GetDimensions()
	if height > actualHeight {
		t.Errorf("Adaptive height %d exceeds actual height %d", height, actualHeight)
	}
}
