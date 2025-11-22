package ui

import (
	"os"
	"strings"
)

// ANSI color codes
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorGray    = "\033[90m"

	// Bright colors
	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightBlue    = "\033[94m"
	ColorBrightMagenta = "\033[95m"
	ColorBrightCyan    = "\033[96m"
	ColorBrightWhite   = "\033[97m"

	// Background colors
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"

	// Text styles
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
)

// TypeColors maps Pokemon types to their corresponding colors
var TypeColors = map[string]string{
	"normal":   ColorWhite,
	"fire":     ColorRed,
	"water":    ColorBlue,
	"electric": ColorYellow,
	"grass":    ColorGreen,
	"ice":      ColorCyan,
	"fighting": ColorBrightRed,
	"poison":   ColorMagenta,
	"ground":   ColorYellow,
	"flying":   ColorBrightCyan,
	"psychic":  ColorBrightMagenta,
	"bug":      ColorGreen,
	"rock":     ColorYellow,
	"ghost":    ColorMagenta,
	"dragon":   ColorBrightBlue,
	"dark":     ColorGray,
	"steel":    ColorWhite,
	"fairy":    ColorBrightMagenta,
}

// Colorize wraps text with the specified color code
// If color support is disabled, returns the text unchanged
func Colorize(text, color string) string {
	if !globalColorSupport {
		return text
	}
	return color + text + ColorReset
}

// ColorizeType wraps text with the color for a specific Pokemon type
func ColorizeType(text, pokemonType string) string {
	if !globalColorSupport {
		return text
	}
	color, exists := TypeColors[strings.ToLower(pokemonType)]
	if !exists {
		color = ColorWhite
	}
	return color + text + ColorReset
}

// globalColorSupport stores the detected color support status
var globalColorSupport bool

// DetectColorSupport checks if the terminal supports ANSI colors
// Supports NO_COLOR environment variable and various terminal emulators
func DetectColorSupport() bool {
	// Check NO_COLOR environment variable (universal standard)
	// https://no-color.org/
	if os.Getenv("NO_COLOR") != "" {
		globalColorSupport = false
		return false
	}

	// Check CLICOLOR_FORCE for forced color output
	if os.Getenv("CLICOLOR_FORCE") != "" && os.Getenv("CLICOLOR_FORCE") != "0" {
		globalColorSupport = true
		return true
	}

	// Check CLICOLOR for color preference
	if os.Getenv("CLICOLOR") == "0" {
		globalColorSupport = false
		return false
	}

	// Check TERM environment variable
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		globalColorSupport = false
		return false
	}

	// Check for common color-supporting terminals
	colorTerms := []string{
		"xterm",
		"xterm-256",
		"xterm-256color",
		"screen",
		"screen-256color",
		"tmux",
		"tmux-256color",
		"rxvt",
		"color",
		"ansi",
		"cygwin",
		"linux",
		"vt100",
		"vt220",
		"konsole",
		"gnome",
		"alacritty",
		"kitty",
	}

	termLower := strings.ToLower(term)
	for _, colorTerm := range colorTerms {
		if strings.Contains(termLower, colorTerm) {
			globalColorSupport = true
			return true
		}
	}

	// Check COLORTERM environment variable (modern terminals)
	if os.Getenv("COLORTERM") != "" {
		globalColorSupport = true
		return true
	}

	// On Windows, check for modern terminal support
	// Windows Terminal, ConEmu, and other modern emulators
	if os.Getenv("WT_SESSION") != "" || 
	   os.Getenv("WT_PROFILE_ID") != "" ||
	   os.Getenv("ConEmuANSI") == "ON" ||
	   os.Getenv("ANSICON") != "" {
		globalColorSupport = true
		return true
	}

	// Check for iTerm2 on macOS
	if strings.Contains(os.Getenv("TERM_PROGRAM"), "iTerm") {
		globalColorSupport = true
		return true
	}

	// Check for VS Code integrated terminal
	if os.Getenv("TERM_PROGRAM") == "vscode" {
		globalColorSupport = true
		return true
	}

	// Default to no color support if we can't determine
	globalColorSupport = false
	return false
}

// SetColorSupport manually sets the color support status
// Useful for testing or user preferences
func SetColorSupport(enabled bool) {
	globalColorSupport = enabled
}

// GetColorSupport returns the current color support status
func GetColorSupport() bool {
	return globalColorSupport
}
