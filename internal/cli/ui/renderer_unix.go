//go:build !windows
// +build !windows

package ui

import (
	"os"
	"os/signal"
	"syscall"
)

// setupResizeHandler sets up a signal handler for terminal resize events (Unix)
func (r *Renderer) setupResizeHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)

	go func() {
		for range sigChan {
			r.handleResize()
		}
	}()
}
