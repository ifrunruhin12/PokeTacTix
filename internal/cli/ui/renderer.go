package ui

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"golang.org/x/term"
)

// Renderer handles terminal UI rendering with color support and dimension detection
type Renderer struct {
	Width        int
	Height       int
	ColorSupport bool
	lastScreen   string // Cache of last rendered screen for diff-based updates
	mu           sync.RWMutex // Protects Width and Height during resize
}

var (
	cachedRenderer *Renderer
	rendererOnce   sync.Once
)

// NewRenderer initializes a new renderer with terminal capabilities detection
// Uses sync.Once to cache the renderer and avoid repeated terminal checks
func NewRenderer() *Renderer {
	rendererOnce.Do(func() {
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

		cachedRenderer = &Renderer{
			Width:        width,
			Height:       height,
			ColorSupport: colorSupport,
		}

		// Set up terminal resize handler
		cachedRenderer.setupResizeHandler()
	})

	return cachedRenderer
}

// handleResize updates the renderer dimensions when terminal is resized
func (r *Renderer) handleResize() {
	width, height := getTerminalSize()
	
	r.mu.Lock()
	r.Width = width
	r.Height = height
	r.mu.Unlock()

	// Clear cache to force redraw with new dimensions
	r.ClearCache()
}

// GetDimensions returns the current terminal dimensions (thread-safe)
func (r *Renderer) GetDimensions() (width, height int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Width, r.Height
}

// IsTerminalTooSmall checks if the terminal meets minimum size requirements
func (r *Renderer) IsTerminalTooSmall() bool {
	width, height := r.GetDimensions()
	return width < 80 || height < 24
}

// GetAdaptiveWidth returns a width value adapted to terminal size
// Returns a smaller width if terminal is constrained
func (r *Renderer) GetAdaptiveWidth() int {
	width, _ := r.GetDimensions()
	if width < 80 {
		return width - 2 // Leave margin
	}
	return 80 // Use standard width
}

// GetAdaptiveHeight returns a height value adapted to terminal size
func (r *Renderer) GetAdaptiveHeight() int {
	_, height := r.GetDimensions()
	if height < 24 {
		return height - 2 // Leave margin
	}
	return 24 // Use standard height
}

// Clear clears the terminal screen
func (r *Renderer) Clear() {
	// Use ANSI escape code to clear screen and move cursor to top-left
	fmt.Print("\033[2J\033[H")
	r.lastScreen = "" // Reset cache
}

// ClearLine clears the current line (more efficient than full clear)
func (r *Renderer) ClearLine() {
	fmt.Print("\033[2K\r")
}

// MoveCursor moves the cursor to the specified position (1-indexed)
func (r *Renderer) MoveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

// RenderBuffered renders content with buffering for better performance
// Only clears and redraws if content has changed
func (r *Renderer) RenderBuffered(content string) {
	if content == r.lastScreen {
		return // No change, skip redraw
	}
	
	// Move cursor to home position instead of clearing (faster)
	fmt.Print("\033[H")
	fmt.Print(content)
	r.lastScreen = content
}

// RenderBufferedWithClear renders content with full clear (use sparingly)
func (r *Renderer) RenderBufferedWithClear(content string) {
	r.Clear()
	fmt.Print(content)
	r.lastScreen = content
}

// ClearCache clears the cached screen data to free memory
func (r *Renderer) ClearCache() {
	r.lastScreen = ""
}

// getTerminalSize detects the current terminal dimensions
// Optimized to try fastest methods first
func getTerminalSize() (width, height int) {
	// Try to get terminal size from file descriptor (fastest method)
	fd := int(os.Stdout.Fd())
	if term.IsTerminal(fd) {
		w, h, err := term.GetSize(fd)
		if err == nil && w > 0 && h > 0 {
			return w, h
		}
	}

	// Fallback: try environment variables (fast)
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

	// Final fallback: use default dimensions (skip slow tput command)
	// tput is slow and rarely needed with modern terminals
	return 80, 24
}


