package ui

import (
	"strings"
	"testing"

	"pokemon-cli/internal/battle"
)

func TestPokemonCardWidth(t *testing.T) {
	renderer := &Renderer{
		Width:        100,
		Height:       30,
		ColorSupport: false,
	}

	testCard := &battle.BattleCard{
		Name:       "Pikachu",
		Level:      5,
		HP:         100,
		HPMax:      100,
		Stamina:    180,
		StaminaMax: 180,
		Attack:     55,
		Defense:    40,
		Speed:      90,
		Types:      []string{"electric"},
	}

	cardStr := renderer.renderPokemonCard(testCard, true)
	lines := strings.Split(cardStr, "\n")

	// Each line should be exactly 26 characters (visual width, not bytes)
	expectedWidth := 26
	for i, line := range lines {
		// Count runes (characters) not bytes
		stripped := stripANSI(line)
		actualWidth := len([]rune(stripped))
		if actualWidth != expectedWidth {
			// Debug: show each character
			runes := []rune(stripped)
			t.Logf("Line %d runes: %v", i, runes)
			t.Errorf("Line %d has width %d, expected %d: %q", i, actualWidth, expectedWidth, stripped)
		}
	}

	// Should have 11 lines (10 lines + 1 empty from final split)
	// Actually, without trailing newline, we get 10 lines which is correct
	expectedLines := 10
	if len(lines) != expectedLines {
		t.Errorf("Card has %d lines, expected %d", len(lines), expectedLines)
	}
}

func TestDeckStatusRendering(t *testing.T) {
	renderer := &Renderer{
		Width:        100,
		Height:       30,
		ColorSupport: false,
	}

	deck := []battle.BattleCard{
		{HP: 100, HPMax: 100, IsKnockedOut: false},
		{HP: 50, HPMax: 100, IsKnockedOut: false},
		{HP: 0, HPMax: 100, IsKnockedOut: true},
	}

	status := renderer.renderDeckStatus("PLAYER", deck, 0)

	// Should contain deck label and indicators
	if !strings.Contains(status, "PLAYER DECK:") {
		t.Error("Deck status should contain label")
	}

	// Should have 3 indicators (one per card)
	indicatorCount := strings.Count(status, "[")
	if indicatorCount != 3 {
		t.Errorf("Expected 3 indicators, got %d", indicatorCount)
	}
}

func TestBattleScreenWidth(t *testing.T) {
	renderer := &Renderer{
		Width:        100,
		Height:       30,
		ColorSupport: false,
	}

	bs := &battle.BattleState{
		Mode:       "1v1",
		TurnNumber: 1,
		WhoseTurn:  "player",
		PlayerDeck: []battle.BattleCard{
			{
				Name:       "Pikachu",
				Level:      5,
				HP:         100,
				HPMax:      100,
				Stamina:    180,
				StaminaMax: 180,
				Attack:     55,
				Defense:    40,
				Speed:      90,
				Types:      []string{"electric"},
			},
		},
		AIDeck: []battle.BattleCard{
			{
				Name:       "Charmander",
				Level:      5,
				HP:         90,
				HPMax:      100,
				Stamina:    150,
				StaminaMax: 160,
				Attack:     52,
				Defense:    43,
				Speed:      65,
				Types:      []string{"fire"},
			},
		},
		PlayerActiveIdx: 0,
		AIActiveIdx:     0,
	}

	screen := renderer.RenderBattleScreen(bs)
	lines := strings.Split(screen, "\n")

	// Check that no line exceeds the width (count runes, not bytes)
	maxWidth := renderer.Width
	for i, line := range lines {
		stripped := stripANSI(line)
		actualWidth := len([]rune(stripped))
		if actualWidth > maxWidth {
			t.Errorf("Line %d exceeds max width: %d > %d\n%q", i, actualWidth, maxWidth, stripped)
		}
	}

	// Check for proper borders
	if !strings.HasPrefix(lines[0], "╔") {
		t.Error("Screen should start with top border")
	}

	lastLine := lines[len(lines)-1]
	if lastLine != "" && !strings.HasPrefix(lastLine, "╚") {
		t.Error("Screen should end with bottom border")
	}
}
