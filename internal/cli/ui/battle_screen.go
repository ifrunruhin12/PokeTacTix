package ui

import (
	"fmt"
	"strings"

	"pokemon-cli/internal/battle"
)

// RenderBattleScreen displays the complete battle state
func (r *Renderer) RenderBattleScreen(bs *battle.BattleState) string {
	var result strings.Builder

	// Top border
	width := r.Width
	if width > 100 {
		width = 100
	}

	result.WriteString("╔")
	result.WriteString(strings.Repeat("═", width-2))
	result.WriteString("╗\n")

	// Title with mode and turn info
	title := fmt.Sprintf("BATTLE - %s MODE", strings.ToUpper(bs.Mode))
	turnInfo := fmt.Sprintf("Turn %d | %s's Turn", bs.TurnNumber, strings.Title(bs.WhoseTurn))
	
	titlePadding := (width - len(title) - len(turnInfo) - 4) / 2
	result.WriteString("║ ")
	if r.ColorSupport {
		result.WriteString(Colorize(title, Bold+ColorBrightCyan))
	} else {
		result.WriteString(title)
	}
	result.WriteString(strings.Repeat(" ", titlePadding))
	if r.ColorSupport {
		result.WriteString(Colorize(turnInfo, ColorYellow))
	} else {
		result.WriteString(turnInfo)
	}
	result.WriteString(strings.Repeat(" ", width-len(title)-len(turnInfo)-titlePadding-4))
	result.WriteString(" ║\n")

	// Separator
	result.WriteString("╠")
	result.WriteString(strings.Repeat("═", width-2))
	result.WriteString("╣\n")

	// Battle arena - show active Pokemon
	result.WriteString(r.renderBattleArena(bs, width))

	// Separator
	result.WriteString("╠")
	result.WriteString(strings.Repeat("═", width-2))
	result.WriteString("╣\n")

	// Bottom border
	result.WriteString("╚")
	result.WriteString(strings.Repeat("═", width-2))
	result.WriteString("╝\n")

	return result.String()
}

// renderBattleArena renders the main battle area with Pokemon
func (r *Renderer) renderBattleArena(bs *battle.BattleState, width int) string {
	var result strings.Builder

	// Get active Pokemon
	playerActive := bs.GetActivePlayerCard()
	aiActive := bs.GetActiveAICard()

	// Deck status line
	result.WriteString("║ ")
	result.WriteString(r.renderDeckStatus("PLAYER", bs.PlayerDeck, bs.PlayerActiveIdx, width/2-4))
	result.WriteString("  VS  ")
	result.WriteString(r.renderDeckStatus("AI", bs.AIDeck, bs.AIActiveIdx, width/2-4))
	result.WriteString(" ║\n")

	result.WriteString("║")
	result.WriteString(strings.Repeat(" ", width-2))
	result.WriteString("║\n")

	// Active Pokemon cards side by side
	playerCard := r.renderPokemonCard(playerActive, true)
	aiCard := r.renderPokemonCard(aiActive, false)

	// Split cards into lines and display side by side
	playerLines := strings.Split(playerCard, "\n")
	aiLines := strings.Split(aiCard, "\n")

	maxLines := len(playerLines)
	if len(aiLines) > maxLines {
		maxLines = len(aiLines)
	}

	cardWidth := 30
	spacing := width - (cardWidth * 2) - 6
	if spacing < 2 {
		spacing = 2
	}

	for i := 0; i < maxLines; i++ {
		result.WriteString("║ ")

		// Player card line
		if i < len(playerLines) {
			line := playerLines[i]
			// Remove ANSI codes for length calculation
			displayLen := len(stripANSI(line))
			result.WriteString(line)
			padding := cardWidth - displayLen
			if padding < 0 {
				padding = 0
			}
			result.WriteString(strings.Repeat(" ", padding))
		} else {
			result.WriteString(strings.Repeat(" ", cardWidth))
		}

		result.WriteString(strings.Repeat(" ", spacing))

		// AI card line
		if i < len(aiLines) {
			line := aiLines[i]
			displayLen := len(stripANSI(line))
			result.WriteString(line)
			padding := cardWidth - displayLen
			if padding < 0 {
				padding = 0
			}
			result.WriteString(strings.Repeat(" ", padding))
		} else {
			result.WriteString(strings.Repeat(" ", cardWidth))
		}

		result.WriteString(" ║\n")
	}

	result.WriteString("║")
	result.WriteString(strings.Repeat(" ", width-2))
	result.WriteString("║\n")

	return result.String()
}

// renderDeckStatus shows the deck status with indicators
func (r *Renderer) renderDeckStatus(label string, deck []battle.BattleCard, activeIdx int, maxWidth int) string {
	var result strings.Builder

	if r.ColorSupport {
		result.WriteString(Colorize(label+" DECK: ", Bold))
	} else {
		result.WriteString(label + " DECK: ")
	}

	for i, card := range deck {
		var indicator string
		if card.IsKnockedOut || card.HP <= 0 {
			indicator = RenderKOIndicator()
		} else if i == activeIdx {
			indicator = RenderActiveIndicator()
		} else {
			indicator = RenderInactiveIndicator()
		}

		result.WriteString("[")
		result.WriteString(indicator)
		result.WriteString("]")
	}

	return result.String()
}

// renderPokemonCard renders a single Pokemon card
func (r *Renderer) renderPokemonCard(card *battle.BattleCard, isPlayer bool) string {
	if card == nil {
		return "No Pokemon"
	}

	var result strings.Builder

	// Card border top
	result.WriteString("┌────────────────────────┐\n")

	// Pokemon name
	nameLabel := "│ "
	if isPlayer {
		nameLabel += "[ACTIVE] "
	} else {
		nameLabel += "[ENEMY]  "
	}
	
	name := card.Name
	if len(name) > 14 {
		name = name[:14]
	}
	
	if r.ColorSupport {
		nameLabel += Colorize(name, Bold+ColorBrightWhite)
	} else {
		nameLabel += name
	}
	
	// Pad to card width
	displayLen := 10 + len(card.Name)
	if len(card.Name) > 14 {
		displayLen = 24
	}
	nameLabel += strings.Repeat(" ", 24-displayLen) + " │"
	result.WriteString(nameLabel + "\n")

	// Level
	levelLine := fmt.Sprintf("│ Lv %-19d │\n", card.Level)
	result.WriteString(levelLine)

	// Types
	typesStr := "│ "
	for i, t := range card.Types {
		if i > 0 {
			typesStr += "/"
		}
		if r.ColorSupport {
			typesStr += ColorizeType(strings.ToUpper(t), t)
		} else {
			typesStr += strings.ToUpper(t)
		}
	}
	typesStr += strings.Repeat(" ", 22-len(strings.Join(card.Types, "/"))) + " │\n"
	result.WriteString(typesStr)

	// Separator
	result.WriteString("├────────────────────────┤\n")

	// HP Bar - render with smaller bar width to fit numeric values
	hpBar := RenderHPBar(card.HP, card.HPMax, 8)
	hpStripped := stripANSI(hpBar)
	hpPadding := 18 - len(hpStripped)
	if hpPadding < 0 {
		hpPadding = 0
	}
	result.WriteString("│ HP:  ")
	result.WriteString(hpBar)
	result.WriteString(strings.Repeat(" ", hpPadding))
	result.WriteString(" │\n")

	// Stamina Bar - render with smaller bar width to fit numeric values
	staminaBar := RenderStaminaBar(card.Stamina, card.StaminaMax, 8)
	staStripped := stripANSI(staminaBar)
	staPadding := 18 - len(staStripped)
	if staPadding < 0 {
		staPadding = 0
	}
	result.WriteString("│ STA: ")
	result.WriteString(staminaBar)
	result.WriteString(strings.Repeat(" ", staPadding))
	result.WriteString(" │\n")

	// Stats
	result.WriteString(fmt.Sprintf("│ ATK: %-3d  DEF: %-3d    │\n", card.Attack, card.Defense))
	result.WriteString(fmt.Sprintf("│ SPD: %-17d │\n", card.Speed))

	// Card border bottom
	result.WriteString("└────────────────────────┘")

	return result.String()
}

// stripANSI removes ANSI escape codes from a string
func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteByte(s[i])
	}

	return result.String()
}

// stripANSIToLength strips ANSI codes and pads/truncates to exact length
func stripANSIToLength(s string, length int) string {
	stripped := stripANSI(s)
	if len(stripped) > length {
		return stripped[:length]
	}
	return stripped + strings.Repeat(" ", length-len(stripped))
}

// RenderBattleActions renders the action menu for battle
func (r *Renderer) RenderBattleActions(actions []string, selected int) string {
	var result strings.Builder

	result.WriteString("\n")
	result.WriteString(r.RenderActionMenu(actions, selected))
	result.WriteString("\n")

	return result.String()
}

// RenderMoveSelection renders the move selection menu
func (r *Renderer) RenderMoveSelection(card *battle.BattleCard, selected int) string {
	var result strings.Builder

	result.WriteString("\n")
	if r.ColorSupport {
		result.WriteString(Colorize("SELECT A MOVE:", Bold+ColorBrightCyan))
	} else {
		result.WriteString("SELECT A MOVE:")
	}
	result.WriteString("\n")
	result.WriteString(strings.Repeat("═", 60))
	result.WriteString("\n\n")

	for i, move := range card.Moves {
		isSelected := i == selected
		canUse := card.Stamina >= move.StaminaCost

		var line string
		if isSelected {
			line = fmt.Sprintf("  ▶ [%d] ", i+1)
		} else {
			line = fmt.Sprintf("    [%d] ", i+1)
		}

		// Move name
		moveName := strings.Title(strings.ReplaceAll(move.Name, "-", " "))
		if r.ColorSupport {
			if isSelected {
				moveName = Colorize(moveName, Bold+ColorBrightYellow)
			}
			moveName = ColorizeType(moveName, move.Type)
		}
		line += moveName

		// Move details
		details := fmt.Sprintf(" | Power: %d | Stamina: %d", move.Power, move.StaminaCost)
		if !canUse {
			if r.ColorSupport {
				details += Colorize(" [NOT ENOUGH STAMINA]", ColorRed)
			} else {
				details += " [NOT ENOUGH STAMINA]"
			}
		}
		line += details

		result.WriteString(line)
		result.WriteString("\n")
	}

	result.WriteString("\n")
	result.WriteString("    [0] Back to Actions\n")

	return result.String()
}

// RenderPokemonSwitchMenu renders the Pokemon switching menu for 5v5
func (r *Renderer) RenderPokemonSwitchMenu(deck []battle.BattleCard, activeIdx int, selected int) string {
	var result strings.Builder

	result.WriteString("\n")
	if r.ColorSupport {
		result.WriteString(Colorize("SELECT NEXT POKEMON:", Bold+ColorBrightCyan))
	} else {
		result.WriteString("SELECT NEXT POKEMON:")
	}
	result.WriteString("\n")
	result.WriteString(strings.Repeat("═", 70))
	result.WriteString("\n\n")

	validIdx := 0
	for i, card := range deck {
		if i == activeIdx || card.IsKnockedOut || card.HP <= 0 {
			continue // Skip active and KO'd Pokemon
		}

		isSelected := validIdx == selected
		validIdx++

		var line string
		if isSelected {
			line = fmt.Sprintf("  ▶ [%d] ", validIdx)
		} else {
			line = fmt.Sprintf("    [%d] ", validIdx)
		}

		// Pokemon info
		info := fmt.Sprintf("%-12s Lv%-2d | HP: %3d/%3d | Types: %s",
			card.Name, card.Level, card.HP, card.HPMax, strings.Join(card.Types, "/"))

		if r.ColorSupport && isSelected {
			info = Colorize(info, Bold+ColorBrightYellow)
		}

		line += info
		result.WriteString(line)
		result.WriteString("\n")
	}

	return result.String()
}
