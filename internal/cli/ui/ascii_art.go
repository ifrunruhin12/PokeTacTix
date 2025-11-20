package ui

import (
	"fmt"
	"strings"
)

// LogoASCII is the PokeTacTix game logo
const LogoASCII = `
╔═══════════════════════════════════════════════════════════════════════════╗
║  ██████╗  ██████╗ ██╗  ██╗███████╗████████╗ █████╗  ██████╗████████╗██╗██╗ ██╗ ║
║  ██╔══██╗██╔═══██╗██║ ██╔╝██╔════╝╚══██╔══╝██╔══██╗██╔════╝╚══██╔══╝██║╚██╗██╔╝ ║
║  ██████╔╝██║   ██║█████╔╝ █████╗     ██║   ███████║██║        ██║   ██║ ╚███╔╝  ║
║  ██╔═══╝ ██║   ██║██╔═██╗ ██╔══╝     ██║   ██╔══██║██║        ██║   ██║ ██╔██╗  ║
║  ██║     ╚██████╔╝██║  ██╗███████╗   ██║   ██║  ██║╚██████╗   ██║   ██║██╔╝ ██╗ ║
║  ╚═╝      ╚═════╝ ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝   ╚═╝   ╚═╝╚═╝  ╚═╝ ║
╚═══════════════════════════════════════════════════════════════════════════╝
`

// PokeballASCII is a simple Pokeball representation for face-down cards
const PokeballASCII = `
     ___
   /     \
  |  ◯ ◯  |
  |   ─   |
   \_____/
`

// SmallPokeballASCII is a compact Pokeball for inline use
const SmallPokeballASCII = `
  ___
 / ◯ \
 \___/
`

// RenderLogo displays the game logo with optional color
func RenderLogo() string {
	if GetColorSupport() {
		return Colorize(LogoASCII, ColorBrightCyan)
	}
	return LogoASCII
}

// RenderHPBar creates an ASCII HP bar visualization
// current: current HP value
// max: maximum HP value
// width: width of the bar in characters (default 20 if <= 0)
func RenderHPBar(current, max int, width int) string {
	if width <= 0 {
		width = 20
	}

	// Ensure current doesn't exceed max
	if current > max {
		current = max
	}
	if current < 0 {
		current = 0
	}

	// Calculate percentage and filled blocks
	percentage := float64(current) / float64(max)
	filled := int(percentage * float64(width))
	empty := width - filled

	// Build the bar
	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	// Color based on HP percentage
	var color string
	if percentage > 0.6 {
		color = ColorGreen
	} else if percentage > 0.3 {
		color = ColorYellow
	} else {
		color = ColorRed
	}

	// Format with color if supported
	if GetColorSupport() {
		bar = Colorize(bar, color)
	}

	return fmt.Sprintf("%s %d/%d", bar, current, max)
}

// RenderStaminaBar creates an ASCII stamina bar visualization
// current: current stamina value
// max: maximum stamina value
// width: width of the bar in characters (default 20 if <= 0)
func RenderStaminaBar(current, max int, width int) string {
	if width <= 0 {
		width = 20
	}

	// Ensure current doesn't exceed max
	if current > max {
		current = max
	}
	if current < 0 {
		current = 0
	}

	// Calculate percentage and filled blocks
	percentage := float64(current) / float64(max)
	filled := int(percentage * float64(width))
	empty := width - filled

	// Build the bar
	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	// Color stamina bar cyan/blue
	if GetColorSupport() {
		bar = Colorize(bar, ColorCyan)
	}

	return fmt.Sprintf("%s %d/%d", bar, current, max)
}

// RenderTypeBadge creates a colored type badge
func RenderTypeBadge(pokemonType string) string {
	typeName := strings.ToUpper(pokemonType)
	
	// Pad to consistent width (10 characters)
	if len(typeName) < 10 {
		padding := (10 - len(typeName)) / 2
		typeName = strings.Repeat(" ", padding) + typeName + strings.Repeat(" ", 10-len(typeName)-padding)
	}

	badge := fmt.Sprintf("[%s]", typeName)

	if GetColorSupport() {
		return ColorizeType(badge, pokemonType)
	}
	return badge
}

// RenderPokemonSprite returns a simple ASCII representation of a Pokemon
// This is a placeholder - in a full implementation, you might have
// specific sprites for different Pokemon
func RenderPokemonSprite(name string) string {
	// Simple generic Pokemon sprite
	sprite := `
    ∧＿∧
   (｡･ω･｡)
   /　 　 つ
  (＿＿_つ
`
	return sprite
}

// RenderActiveIndicator returns an indicator for active Pokemon
func RenderActiveIndicator() string {
	indicator := "●"
	if GetColorSupport() {
		return Colorize(indicator, ColorBrightGreen)
	}
	return indicator
}

// RenderInactiveIndicator returns an indicator for inactive Pokemon
func RenderInactiveIndicator() string {
	indicator := "●"
	if GetColorSupport() {
		return Colorize(indicator, ColorGray)
	}
	return indicator
}

// RenderKOIndicator returns an indicator for knocked out Pokemon
func RenderKOIndicator() string {
	indicator := "○"
	if GetColorSupport() {
		return Colorize(indicator, ColorRed)
	}
	return indicator
}

// RenderBox creates a simple box around text
func RenderBox(title string, content []string, width int) string {
	if width < 10 {
		width = 40
	}

	var result strings.Builder

	// Top border with title
	titleLen := len(title)
	leftPad := (width - titleLen - 2) / 2
	rightPad := width - titleLen - 2 - leftPad
	
	result.WriteString("┌")
	result.WriteString(strings.Repeat("─", leftPad))
	result.WriteString(" " + title + " ")
	result.WriteString(strings.Repeat("─", rightPad))
	result.WriteString("┐\n")

	// Content lines
	for _, line := range content {
		// Truncate or pad line to fit width
		if len(line) > width-2 {
			line = line[:width-5] + "..."
		} else {
			line = line + strings.Repeat(" ", width-2-len(line))
		}
		result.WriteString("│ " + line + " │\n")
	}

	// Bottom border
	result.WriteString("└")
	result.WriteString(strings.Repeat("─", width))
	result.WriteString("┘")

	return result.String()
}

// RenderDivider creates a horizontal divider line
func RenderDivider(width int, char string) string {
	if char == "" {
		char = "─"
	}
	return strings.Repeat(char, width)
}

// RenderProgressBar creates a generic progress bar
func RenderProgressBar(current, max int, width int, label string) string {
	if width <= 0 {
		width = 20
	}

	percentage := float64(current) / float64(max)
	filled := int(percentage * float64(width))
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	
	return fmt.Sprintf("%s: %s %d/%d (%.0f%%)", label, bar, current, max, percentage*100)
}
