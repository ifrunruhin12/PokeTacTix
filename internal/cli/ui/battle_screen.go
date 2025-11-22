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
	playerStatus := r.renderDeckStatus("PLAYER", bs.PlayerDeck, bs.PlayerActiveIdx)
	aiStatus := r.renderDeckStatus("AI", bs.AIDeck, bs.AIActiveIdx)
	vsText := "  VS  "
	
	// Calculate total length without ANSI codes
	statusLen := len(stripANSI(playerStatus)) + len(vsText) + len(stripANSI(aiStatus))
	padding := width - statusLen - 4 // 4 for borders and spaces
	if padding < 0 {
		padding = 0
	}
	
	result.WriteString("║ ")
	result.WriteString(playerStatus)
	result.WriteString(vsText)
	result.WriteString(aiStatus)
	result.WriteString(strings.Repeat(" ", padding))
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

	// Card width is 26 characters (including borders)
	cardWidth := 26
	// Calculate spacing between cards
	totalCardWidth := cardWidth * 2
	availableSpace := width - 4 // 4 for outer borders and padding
	spacing := availableSpace - totalCardWidth
	if spacing < 2 {
		spacing = 2
	}

	for i := 0; i < maxLines; i++ {
		result.WriteString("║ ")

		// Player card line
		if i < len(playerLines) && playerLines[i] != "" {
			line := playerLines[i]
			result.WriteString(line)
		} else {
			result.WriteString(strings.Repeat(" ", cardWidth))
		}

		result.WriteString(strings.Repeat(" ", spacing))

		// AI card line
		if i < len(aiLines) && aiLines[i] != "" {
			line := aiLines[i]
			result.WriteString(line)
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
func (r *Renderer) renderDeckStatus(label string, deck []battle.BattleCard, activeIdx int) string {
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

// renderPokemonCard renders a single Pokemon card (26 chars wide)
func (r *Renderer) renderPokemonCard(card *battle.BattleCard, isPlayer bool) string {
	if card == nil {
		// Return empty card placeholder
		var result strings.Builder
		for i := 0; i < 11; i++ {
			result.WriteString(strings.Repeat(" ", 26))
			if i < 10 {
				result.WriteString("\n")
			}
		}
		return result.String()
	}

	var result strings.Builder

	// Card border top (26 chars)
	result.WriteString("┌────────────────────────┐")

	// Pokemon name
	result.WriteString("\n│ ")
	if isPlayer {
		result.WriteString("[ACTIVE] ")
	} else {
		result.WriteString("[ENEMY]  ")
	}
	
	name := card.Name
	nameLen := len(name)
	if nameLen > 13 {
		name = name[:13]
		nameLen = 13
	}
	
	if r.ColorSupport {
		result.WriteString(Colorize(name, Bold+ColorBrightWhite))
	} else {
		result.WriteString(name)
	}
	
	// Pad to 22 chars inside (24 - 2 for borders)
	// 9 for label + name + padding = 22
	padding := 22 - 9 - nameLen
	if padding < 0 {
		padding = 0
	}
	result.WriteString(strings.Repeat(" ", padding))
	result.WriteString(" │")

	// Level
	result.WriteString(fmt.Sprintf("\n│ Lv %-19d │", card.Level))

	// Types
	result.WriteString("\n│ ")
	typesText := ""
	for i, t := range card.Types {
		if i > 0 {
			typesText += "/"
		}
		typesText += strings.ToUpper(t)
	}
	
	if r.ColorSupport {
		// Apply color to each type
		coloredTypes := ""
		for i, t := range card.Types {
			if i > 0 {
				coloredTypes += "/"
			}
			coloredTypes += ColorizeType(strings.ToUpper(t), t)
		}
		result.WriteString(coloredTypes)
	} else {
		result.WriteString(typesText)
	}
	
	typePadding := 22 - len(typesText)
	if typePadding < 0 {
		typePadding = 0
	}
	result.WriteString(strings.Repeat(" ", typePadding))
	result.WriteString(" │")

	// Separator
	result.WriteString("\n├────────────────────────┤")

	// HP Bar - use smaller bar width since it includes numbers
	result.WriteString("\n│ HP:  ")
	hpBar := RenderHPBar(card.HP, card.HPMax, 6)
	result.WriteString(hpBar)
	hpStripped := stripANSI(hpBar)
	// Total width: 26 chars
	// "│ HP:  " = 7 chars, "│" = 1 char, so bar + padding = 18 chars
	hpPadding := 18 - len([]rune(hpStripped))
	if hpPadding < 0 {
		hpPadding = 0
	}
	result.WriteString(strings.Repeat(" ", hpPadding))
	result.WriteString("│")

	// Stamina Bar - use smaller bar width since it includes numbers
	result.WriteString("\n│ STA: ")
	staminaBar := RenderStaminaBar(card.Stamina, card.StaminaMax, 6)
	result.WriteString(staminaBar)
	staStripped := stripANSI(staminaBar)
	// "│ STA: " = 7 chars, "│" = 1 char, so bar + padding = 18 chars
	staPadding := 18 - len([]rune(staStripped))
	if staPadding < 0 {
		staPadding = 0
	}
	result.WriteString(strings.Repeat(" ", staPadding))
	result.WriteString("│")

	// Stats - ensure proper spacing (24 chars inside borders)
	result.WriteString(fmt.Sprintf("\n│ ATK: %-3d  DEF: %-3d     │", card.Attack, card.Defense))
	result.WriteString(fmt.Sprintf("\n│ SPD: %-18d│", card.Speed))

	// Card border bottom
	result.WriteString("\n└────────────────────────┘")

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

// RenderBattleScreenCondensed displays a condensed battle state for quick battle mode
func (r *Renderer) RenderBattleScreenCondensed(bs *battle.BattleState) string {
	var result strings.Builder

	// Simple header
	result.WriteString(strings.Repeat("═", 60))
	result.WriteString("\n")
	
	title := fmt.Sprintf("BATTLE %s | Turn %d", strings.ToUpper(bs.Mode), bs.TurnNumber)
	if r.ColorSupport {
		result.WriteString(Colorize(title, Bold+ColorBrightCyan))
	} else {
		result.WriteString(title)
	}
	result.WriteString("\n")
	result.WriteString(strings.Repeat("═", 60))
	result.WriteString("\n\n")

	// Get active Pokemon
	playerActive := bs.GetActivePlayerCard()
	aiActive := bs.GetActiveAICard()

	// Player Pokemon (condensed)
	if playerActive != nil {
		result.WriteString("PLAYER: ")
		if r.ColorSupport {
			result.WriteString(Colorize(playerActive.Name, Bold+ColorBrightGreen))
		} else {
			result.WriteString(playerActive.Name)
		}
		result.WriteString(fmt.Sprintf(" (Lv %d)\n", playerActive.Level))
		
		// HP bar
		result.WriteString("  HP:  ")
		result.WriteString(RenderHPBar(playerActive.HP, playerActive.HPMax, 30))
		result.WriteString(fmt.Sprintf(" %d/%d\n", playerActive.HP, playerActive.HPMax))
		
		// Stamina bar
		result.WriteString("  STA: ")
		result.WriteString(RenderStaminaBar(playerActive.Stamina, playerActive.Speed*2, 30))
		result.WriteString(fmt.Sprintf(" %d/%d\n", playerActive.Stamina, playerActive.Speed*2))
		
		// Deck status
		aliveCount := 0
		for _, card := range bs.PlayerDeck {
			if card.HP > 0 {
				aliveCount++
			}
		}
		result.WriteString(fmt.Sprintf("  Deck: %d/%d alive\n", aliveCount, len(bs.PlayerDeck)))
	}

	result.WriteString("\n")

	// AI Pokemon (condensed)
	if aiActive != nil {
		result.WriteString("AI: ")
		if r.ColorSupport {
			result.WriteString(Colorize(aiActive.Name, Bold+ColorBrightRed))
		} else {
			result.WriteString(aiActive.Name)
		}
		result.WriteString(fmt.Sprintf(" (Lv %d)\n", aiActive.Level))
		
		// HP bar
		result.WriteString("  HP:  ")
		result.WriteString(RenderHPBar(aiActive.HP, aiActive.HPMax, 30))
		result.WriteString(fmt.Sprintf(" %d/%d\n", aiActive.HP, aiActive.HPMax))
		
		// Stamina bar
		result.WriteString("  STA: ")
		result.WriteString(RenderStaminaBar(aiActive.Stamina, aiActive.Speed*2, 30))
		result.WriteString(fmt.Sprintf(" %d/%d\n", aiActive.Stamina, aiActive.Speed*2))
		
		// Deck status
		aliveCount := 0
		for _, card := range bs.AIDeck {
			if card.HP > 0 {
				aliveCount++
			}
		}
		result.WriteString(fmt.Sprintf("  Deck: %d/%d alive\n", aliveCount, len(bs.AIDeck)))
	}

	result.WriteString("\n")
	result.WriteString(strings.Repeat("═", 60))
	result.WriteString("\n")

	return result.String()
}
