package ui

import (
	"os"
	"testing"
)

func TestDetectColorSupport(t *testing.T) {
	// Save original environment
	origNoColor := os.Getenv("NO_COLOR")
	origTerm := os.Getenv("TERM")
	origColorTerm := os.Getenv("COLORTERM")
	origCliColor := os.Getenv("CLICOLOR")
	origCliColorForce := os.Getenv("CLICOLOR_FORCE")
	
	defer func() {
		os.Setenv("NO_COLOR", origNoColor)
		os.Setenv("TERM", origTerm)
		os.Setenv("COLORTERM", origColorTerm)
		os.Setenv("CLICOLOR", origCliColor)
		os.Setenv("CLICOLOR_FORCE", origCliColorForce)
	}()

	tests := []struct {
		name     string
		setup    func()
		expected bool
	}{
		{
			name: "NO_COLOR set",
			setup: func() {
				os.Setenv("NO_COLOR", "1")
				os.Setenv("TERM", "xterm-256color")
			},
			expected: false,
		},
		{
			name: "CLICOLOR_FORCE set",
			setup: func() {
				os.Unsetenv("NO_COLOR")
				os.Setenv("CLICOLOR_FORCE", "1")
				os.Setenv("TERM", "dumb")
			},
			expected: true,
		},
		{
			name: "CLICOLOR disabled",
			setup: func() {
				os.Unsetenv("NO_COLOR")
				os.Unsetenv("CLICOLOR_FORCE")
				os.Setenv("CLICOLOR", "0")
				os.Setenv("TERM", "xterm-256color")
			},
			expected: false,
		},
		{
			name: "xterm-256color",
			setup: func() {
				os.Unsetenv("NO_COLOR")
				os.Unsetenv("CLICOLOR")
				os.Unsetenv("CLICOLOR_FORCE")
				os.Setenv("TERM", "xterm-256color")
			},
			expected: true,
		},
		{
			name: "dumb terminal",
			setup: func() {
				os.Unsetenv("NO_COLOR")
				os.Unsetenv("CLICOLOR")
				os.Unsetenv("CLICOLOR_FORCE")
				os.Setenv("TERM", "dumb")
				os.Unsetenv("COLORTERM")
			},
			expected: false,
		},
		{
			name: "COLORTERM set",
			setup: func() {
				os.Unsetenv("NO_COLOR")
				os.Unsetenv("CLICOLOR")
				os.Unsetenv("CLICOLOR_FORCE")
				os.Setenv("TERM", "unknown")
				os.Setenv("COLORTERM", "truecolor")
			},
			expected: true,
		},
		{
			name: "Windows Terminal",
			setup: func() {
				os.Unsetenv("NO_COLOR")
				os.Unsetenv("CLICOLOR")
				os.Unsetenv("CLICOLOR_FORCE")
				os.Setenv("TERM", "unknown")
				os.Unsetenv("COLORTERM")
				os.Setenv("WT_SESSION", "12345")
			},
			expected: true,
		},
		{
			name: "iTerm2",
			setup: func() {
				os.Unsetenv("NO_COLOR")
				os.Unsetenv("CLICOLOR")
				os.Unsetenv("CLICOLOR_FORCE")
				os.Setenv("TERM", "unknown")
				os.Unsetenv("COLORTERM")
				os.Unsetenv("WT_SESSION")
				os.Setenv("TERM_PROGRAM", "iTerm.app")
			},
			expected: true,
		},
		{
			name: "VS Code",
			setup: func() {
				os.Unsetenv("NO_COLOR")
				os.Unsetenv("CLICOLOR")
				os.Unsetenv("CLICOLOR_FORCE")
				os.Setenv("TERM", "unknown")
				os.Unsetenv("COLORTERM")
				os.Unsetenv("WT_SESSION")
				os.Setenv("TERM_PROGRAM", "vscode")
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result := DetectColorSupport()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestColorize(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		color        string
		colorEnabled bool
		wantColor    bool
	}{
		{
			name:         "with color support",
			text:         "test",
			color:        ColorRed,
			colorEnabled: true,
			wantColor:    true,
		},
		{
			name:         "without color support",
			text:         "test",
			color:        ColorRed,
			colorEnabled: false,
			wantColor:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetColorSupport(tt.colorEnabled)
			result := Colorize(tt.text, tt.color)
			
			if tt.wantColor {
				if result == tt.text {
					t.Error("Expected colored text, got plain text")
				}
				if result != tt.color+tt.text+ColorReset {
					t.Errorf("Expected %q, got %q", tt.color+tt.text+ColorReset, result)
				}
			} else {
				if result != tt.text {
					t.Errorf("Expected plain text %q, got %q", tt.text, result)
				}
			}
		})
	}
}

func TestColorizeType(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		pokemonType  string
		colorEnabled bool
		wantColor    bool
	}{
		{
			name:         "fire type with color",
			text:         "Charmander",
			pokemonType:  "fire",
			colorEnabled: true,
			wantColor:    true,
		},
		{
			name:         "water type with color",
			text:         "Squirtle",
			pokemonType:  "water",
			colorEnabled: true,
			wantColor:    true,
		},
		{
			name:         "unknown type with color",
			text:         "Unknown",
			pokemonType:  "unknown",
			colorEnabled: true,
			wantColor:    true,
		},
		{
			name:         "fire type without color",
			text:         "Charmander",
			pokemonType:  "fire",
			colorEnabled: false,
			wantColor:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetColorSupport(tt.colorEnabled)
			result := ColorizeType(tt.text, tt.pokemonType)
			
			if tt.wantColor {
				if result == tt.text {
					t.Error("Expected colored text, got plain text")
				}
			} else {
				if result != tt.text {
					t.Errorf("Expected plain text %q, got %q", tt.text, result)
				}
			}
		})
	}
}

func TestSetAndGetColorSupport(t *testing.T) {
	// Test setting color support
	SetColorSupport(true)
	if !GetColorSupport() {
		t.Error("Expected color support to be enabled")
	}

	SetColorSupport(false)
	if GetColorSupport() {
		t.Error("Expected color support to be disabled")
	}
}

func TestTypeColors(t *testing.T) {
	// Verify all expected types have colors defined
	expectedTypes := []string{
		"normal", "fire", "water", "electric", "grass", "ice",
		"fighting", "poison", "ground", "flying", "psychic",
		"bug", "rock", "ghost", "dragon", "dark", "steel", "fairy",
	}

	for _, pokemonType := range expectedTypes {
		if _, exists := TypeColors[pokemonType]; !exists {
			t.Errorf("Type %s missing from TypeColors map", pokemonType)
		}
	}
}

func TestColorConstants(t *testing.T) {
	// Verify color constants are not empty
	colors := map[string]string{
		"ColorReset":   ColorReset,
		"ColorRed":     ColorRed,
		"ColorGreen":   ColorGreen,
		"ColorYellow":  ColorYellow,
		"ColorBlue":    ColorBlue,
		"ColorMagenta": ColorMagenta,
		"ColorCyan":    ColorCyan,
		"ColorWhite":   ColorWhite,
	}

	for name, value := range colors {
		if value == "" {
			t.Errorf("Color constant %s is empty", name)
		}
	}
}

func TestNO_COLORStandard(t *testing.T) {
	// Test NO_COLOR environment variable (https://no-color.org/)
	origNoColor := os.Getenv("NO_COLOR")
	origTerm := os.Getenv("TERM")
	
	defer func() {
		if origNoColor != "" {
			os.Setenv("NO_COLOR", origNoColor)
		} else {
			os.Unsetenv("NO_COLOR")
		}
		os.Setenv("TERM", origTerm)
	}()

	// Set up a terminal that would normally support colors
	os.Setenv("TERM", "xterm-256color")
	os.Setenv("NO_COLOR", "1")

	result := DetectColorSupport()
	if result {
		t.Error("NO_COLOR=1 should disable color support regardless of TERM")
	}

	// Test with any non-empty NO_COLOR value (should disable)
	os.Setenv("NO_COLOR", "true")
	result = DetectColorSupport()
	if result {
		t.Error("NO_COLOR=true should disable color support")
	}

	// Test without NO_COLOR (should enable colors)
	os.Unsetenv("NO_COLOR")
	result = DetectColorSupport()
	if !result {
		t.Error("Without NO_COLOR, xterm-256color should support colors")
	}
}
