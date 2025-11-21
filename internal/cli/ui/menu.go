package ui

import (
	"fmt"
	"strings"
)

// MenuOption represents a single menu option
type MenuOption struct {
	Label       string
	Description string
	Value       string
}

// RenderMenu displays a menu with options and highlights the selected one
// options: list of menu options to display
// selected: index of the currently selected option (0-based)
// title: optional title for the menu
func (r *Renderer) RenderMenu(options []MenuOption, selected int, title string) string {
	var result strings.Builder

	// Add title if provided
	if title != "" {
		result.WriteString("\n")
		if r.ColorSupport {
			result.WriteString(Colorize(title, Bold+ColorBrightCyan))
		} else {
			result.WriteString(title)
		}
		result.WriteString("\n")
		result.WriteString(strings.Repeat("═", len(title)))
		result.WriteString("\n\n")
	}

	// Render each option
	for i, option := range options {
		isSelected := i == selected

		// Build the option line
		var line string
		if isSelected {
			// Highlight selected option
			if r.ColorSupport {
				line = fmt.Sprintf("  ▶ [%d] %s", i+1, option.Label)
				line = Colorize(line, Bold+ColorBrightYellow)
			} else {
				line = fmt.Sprintf(" >[%d] %s", i+1, option.Label)
			}
		} else {
			line = fmt.Sprintf("   [%d] %s", i+1, option.Label)
		}

		result.WriteString(line)

		// Add description if available
		if option.Description != "" {
			result.WriteString("\n")
			descLine := "       " + option.Description
			if r.ColorSupport {
				descLine = Colorize(descLine, ColorGray)
			}
			result.WriteString(descLine)
		}

		result.WriteString("\n")
	}

	result.WriteString("\n")
	result.WriteString("Use arrow keys or enter a number to select: ")

	return result.String()
}

// RenderSimpleMenu displays a simple menu without descriptions
func (r *Renderer) RenderSimpleMenu(options []string, selected int, title string) string {
	menuOptions := make([]MenuOption, len(options))
	for i, opt := range options {
		menuOptions[i] = MenuOption{Label: opt}
	}
	return r.RenderMenu(menuOptions, selected, title)
}

// RenderActionMenu displays a horizontal action menu (for battle actions)
func (r *Renderer) RenderActionMenu(actions []string, selected int) string {
	var result strings.Builder

	result.WriteString("ACTIONS: ")

	for i, action := range actions {
		isSelected := i == selected

		var actionText string
		if isSelected {
			if r.ColorSupport {
				actionText = fmt.Sprintf("[%d] %s", i+1, action)
				actionText = Colorize(actionText, Bold+BgBlue+ColorBrightWhite)
			} else {
				actionText = fmt.Sprintf(">[%d] %s<", i+1, action)
			}
		} else {
			actionText = fmt.Sprintf("[%d] %s", i+1, action)
		}

		result.WriteString(actionText)
		if i < len(actions)-1 {
			result.WriteString("  ")
		}
	}

	return result.String()
}

// RenderBorderedMenu displays a menu with a border
func (r *Renderer) RenderBorderedMenu(options []MenuOption, selected int, title string) string {
	var result strings.Builder

	// Calculate width based on longest option
	maxWidth := len(title)
	for _, opt := range options {
		optLen := len(opt.Label) + 10 // Account for numbering and padding
		if len(opt.Description) > 0 {
			descLen := len(opt.Description) + 7
			if descLen > optLen {
				optLen = descLen
			}
		}
		if optLen > maxWidth {
			maxWidth = optLen
		}
	}

	if maxWidth < 40 {
		maxWidth = 40
	}
	if maxWidth > r.Width-4 {
		maxWidth = r.Width - 4
	}

	// Top border
	result.WriteString("╔")
	result.WriteString(strings.Repeat("═", maxWidth))
	result.WriteString("╗\n")

	// Title
	if title != "" {
		titlePadding := (maxWidth - len(title)) / 2
		result.WriteString("║")
		result.WriteString(strings.Repeat(" ", titlePadding))
		if r.ColorSupport {
			result.WriteString(Colorize(title, Bold+ColorBrightCyan))
		} else {
			result.WriteString(title)
		}
		result.WriteString(strings.Repeat(" ", maxWidth-len(title)-titlePadding))
		result.WriteString("║\n")

		// Separator
		result.WriteString("╠")
		result.WriteString(strings.Repeat("═", maxWidth))
		result.WriteString("╣\n")
	}

	// Options
	for i, option := range options {
		isSelected := i == selected

		// Build plain text version for length calculation
		plainOption := fmt.Sprintf("   [%d] %s", i+1, option.Label)
		plainLen := len(plainOption)

		// Option line with potential coloring
		var optionLine string
		if isSelected {
			if r.ColorSupport {
				optionLine = fmt.Sprintf(" ▶ [%d] %s", i+1, option.Label)
				optionLine = Colorize(optionLine, Bold+ColorBrightYellow)
			} else {
				optionLine = fmt.Sprintf(" >[%d] %s", i+1, option.Label)
			}
		} else {
			optionLine = plainOption
		}

		// Pad to width (use plain length for calculation)
		padding := maxWidth - plainLen
		if padding < 0 {
			padding = 0
		}

		result.WriteString("║")
		result.WriteString(optionLine)
		result.WriteString(strings.Repeat(" ", padding))
		result.WriteString("║\n")

		// Description line if present
		if option.Description != "" {
			plainDesc := "     " + option.Description
			plainDescLen := len(plainDesc)
			
			// Truncate if too long
			if plainDescLen > maxWidth {
				plainDesc = plainDesc[:maxWidth-3] + "..."
				plainDescLen = maxWidth
			}
			
			descLine := plainDesc
			if r.ColorSupport {
				descLine = Colorize(plainDesc, ColorGray)
			}
			
			descPadding := maxWidth - plainDescLen
			if descPadding < 0 {
				descPadding = 0
			}

			result.WriteString("║")
			result.WriteString(descLine)
			result.WriteString(strings.Repeat(" ", descPadding))
			result.WriteString("║\n")
		}
	}

	// Bottom border
	result.WriteString("╚")
	result.WriteString(strings.Repeat("═", maxWidth))
	result.WriteString("╝\n")

	return result.String()
}

// RenderPrompt displays a prompt for user input
func (r *Renderer) RenderPrompt(prompt string) string {
	if r.ColorSupport {
		return Colorize(prompt+" ", Bold+ColorBrightGreen)
	}
	return prompt + " "
}

// RenderConfirmation displays a yes/no confirmation prompt
func (r *Renderer) RenderConfirmation(message string) string {
	var result strings.Builder
	
	result.WriteString("\n")
	if r.ColorSupport {
		result.WriteString(Colorize(message, ColorBrightYellow))
	} else {
		result.WriteString(message)
	}
	result.WriteString("\n")
	result.WriteString(r.RenderPrompt("Continue? (y/n):"))
	
	return result.String()
}
