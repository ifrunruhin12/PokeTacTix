package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// Renderer handles terminal UI rendering with color support and dimension detection
type Renderer struct {
	Width        int
	Height       int
	ColorSupport bool
}

// NewRenderer initializes a new renderer with terminal capabilities detection
func NewRenderer() *Renderer {
	width, height := getTerminalSize()
	colorSupport := DetectColorSupport()

	// Warn if terminal is too small
	if width < 80 || height < 24 {
		fmt.Printf("Warning: Terminal size (%dx%d) is smaller than recommended minimum (80x24).\n", width, height)
		fmt.Println("Some UI elements may not display correctly.")
	}

	// Notify about color support
	if !colorSupport {
		fmt.Println("Note: Color support not detected. Using plain text mode.")
	}

	return &Renderer{
		Width:        width,
		Height:       height,
		ColorSupport: colorSupport,
	}
}

// Clear clears the terminal screen
func (r *Renderer) Clear() {
	// Use ANSI escape code to clear screen and move cursor to top-left
	fmt.Print("\033[2J\033[H")
}

// getTerminalSize detects the current terminal dimensions
func getTerminalSize() (width, height int) {
	// Try to get terminal size from file descriptor
	fd := int(os.Stdout.Fd())
	if term.IsTerminal(fd) {
		w, h, err := term.GetSize(fd)
		if err == nil && w > 0 && h > 0 {
			return w, h
		}
	}

	// Fallback: try environment variables
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if w, err := strconv.Atoi(cols); err == nil && w > 0 {
			width = w
		}
	}
	if lines := os.Getenv("LINES"); lines != "" {
		if h, err := strconv.Atoi(lines); err == nil && h > 0 {
			height = h
		}
	}

	// If we got valid dimensions from env vars, return them
	if width > 0 && height > 0 {
		return width, height
	}

	// Fallback: try tput command (Unix-like systems)
	if w, h := getTputSize(); w > 0 && h > 0 {
		return w, h
	}

	// Final fallback: use default dimensions
	return 80, 24
}

// getTputSize attempts to get terminal size using tput command
func getTputSize() (width, height int) {
	// Try to get columns
	if cmd := exec.Command("tput", "cols"); cmd.Err == nil {
		if output, err := cmd.Output(); err == nil {
			if w, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil {
				width = w
			}
		}
	}

	// Try to get lines
	if cmd := exec.Command("tput", "lines"); cmd.Err == nil {
		if output, err := cmd.Output(); err == nil {
			if h, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil {
				height = h
			}
		}
	}

	return width, height
}
