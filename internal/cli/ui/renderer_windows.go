//go:build windows
// +build windows

package ui

// setupResizeHandler is a no-op on Windows
// Windows doesn't support SIGWINCH signal for terminal resize events
func (r *Renderer) setupResizeHandler() {
	// No-op on Windows
	// Terminal resize detection would require Windows-specific APIs
	// which are not critical for basic functionality
}
