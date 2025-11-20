package ui

import (
	"testing"

	"pokemon-cli/internal/battle"
	"pokemon-cli/internal/pokemon"
)

func TestNewRenderer(t *testing.T) {
	renderer := NewRenderer()
	if renderer == nil {
		t.Fatal("NewRenderer returned nil")
	}

	if renderer.Width < 80 || renderer.Height < 24 {
		t.Logf("Warning: Terminal size (%dx%d) is smaller than recommended", renderer.Width, renderer.Height)
	}
}

func TestColorize(t *testing.T) {
	// Test with color support enabled
	SetColorSupport(true)
	result := Colorize("test", ColorRed)
	if result != ColorRed+"test"+ColorReset {
		t.Errorf("Colorize with support enabled failed: got %q", result)
	}

	// Test with color support disabled
	SetColorSupport(false)
	result = Colorize("test", ColorRed)
	if result != "test" {
		t.Errorf("Colorize with support disabled failed: got %q", result)
	}
}

func TestColorizeType(t *testing.T) {
	SetColorSupport(true)
	
	tests := []struct {
		text     string
		poketype string
		expected string
	}{
		{"Fire", "fire", ColorRed + "Fire" + ColorReset},
		{"Water", "water", ColorBlue + "Water" + ColorReset},
		{"Grass", "grass", ColorGreen + "Grass" + ColorReset},
	}

	for _, tt := range tests {
		result := ColorizeType(tt.text, tt.poketype)
		if result != tt.expected {
			t.Errorf("ColorizeType(%q, %q) = %q, want %q", tt.text, tt.poketype, result, tt.expected)
		}
	}
}

func TestRenderHPBar(t *testing.T) {
	tests := []struct {
		current int
		max     int
		width   int
	}{
		{100, 100, 20},
		{50, 100, 20},
		{25, 100, 20},
		{0, 100, 20},
	}

	for _, tt := range tests {
		result := RenderHPBar(tt.current, tt.max, tt.width)
		if result == "" {
			t.Errorf("RenderHPBar(%d, %d, %d) returned empty string", tt.current, tt.max, tt.width)
		}
	}
}

func TestRenderStaminaBar(t *testing.T) {
	result := RenderStaminaBar(50, 100, 20)
	if result == "" {
		t.Error("RenderStaminaBar returned empty string")
	}
}

func TestRenderTypeBadge(t *testing.T) {
	types := []string{"fire", "water", "grass", "electric"}
	for _, poketype := range types {
		result := RenderTypeBadge(poketype)
		if result == "" {
			t.Errorf("RenderTypeBadge(%q) returned empty string", poketype)
		}
	}
}

func TestRenderLogo(t *testing.T) {
	logo := RenderLogo()
	if logo == "" {
		t.Error("RenderLogo returned empty string")
	}
}

func TestRenderBox(t *testing.T) {
	content := []string{"Line 1", "Line 2", "Line 3"}
	result := RenderBox("Test Box", content, 40)
	if result == "" {
		t.Error("RenderBox returned empty string")
	}
}

func TestBattleLog(t *testing.T) {
	log := NewBattleLog(10)
	if log == nil {
		t.Fatal("NewBattleLog returned nil")
	}

	log.Add("Test message", LogTypeInfo)
	if len(log.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(log.Entries))
	}

	// Test max entries limit
	for i := 0; i < 15; i++ {
		log.Add("Message", LogTypeInfo)
	}
	if len(log.Entries) > 10 {
		t.Errorf("Expected max 10 entries, got %d", len(log.Entries))
	}
}

func TestDetectLogType(t *testing.T) {
	tests := []struct {
		message  string
		expected LogEntryType
	}{
		{"Pikachu dealt 50 damage!", LogTypeDamage},
		{"Charizard healed 30 HP", LogTypeHealing},
		{"Bulbasaur defended", LogTypeStatus},
		{"Victory!", LogTypeVictory},
		{"Defeat!", LogTypeDefeat},
		{"Warning: Low stamina", LogTypeWarning},
		{"Error occurred", LogTypeError},
		{"Turn 1", LogTypeInfo},
	}

	for _, tt := range tests {
		result := detectLogType(tt.message)
		if result != tt.expected {
			t.Errorf("detectLogType(%q) = %v, want %v", tt.message, result, tt.expected)
		}
	}
}

func TestRenderMenu(t *testing.T) {
	renderer := NewRenderer()
	options := []MenuOption{
		{Label: "Option 1", Description: "First option"},
		{Label: "Option 2", Description: "Second option"},
		{Label: "Option 3", Description: "Third option"},
	}

	result := renderer.RenderMenu(options, 0, "Test Menu")
	if result == "" {
		t.Error("RenderMenu returned empty string")
	}
}

func TestRenderSimpleMenu(t *testing.T) {
	renderer := NewRenderer()
	options := []string{"Start", "Options", "Quit"}

	result := renderer.RenderSimpleMenu(options, 1, "Main Menu")
	if result == "" {
		t.Error("RenderSimpleMenu returned empty string")
	}
}

func TestRenderActionMenu(t *testing.T) {
	renderer := NewRenderer()
	actions := []string{"Attack", "Defend", "Pass", "Surrender"}

	result := renderer.RenderActionMenu(actions, 0)
	if result == "" {
		t.Error("RenderActionMenu returned empty string")
	}
}

func TestRenderBattleScreen(t *testing.T) {
	renderer := NewRenderer()

	// Create a mock battle state
	bs := &battle.BattleState{
		Mode:            "1v1",
		TurnNumber:      1,
		WhoseTurn:       "player",
		PlayerActiveIdx: 0,
		AIActiveIdx:     0,
		PlayerDeck: []battle.BattleCard{
			{
				Name:       "Pikachu",
				HP:         100,
				HPMax:      100,
				Stamina:    80,
				StaminaMax: 100,
				Attack:     55,
				Defense:    40,
				Speed:      90,
				Types:      []string{"electric"},
				Level:      15,
			},
		},
		AIDeck: []battle.BattleCard{
			{
				Name:       "Charmander",
				HP:         80,
				HPMax:      100,
				Stamina:    60,
				StaminaMax: 80,
				Attack:     52,
				Defense:    43,
				Speed:      65,
				Types:      []string{"fire"},
				Level:      12,
			},
		},
	}

	result := renderer.RenderBattleScreen(bs)
	if result == "" {
		t.Error("RenderBattleScreen returned empty string")
	}
}

func TestRenderMoveSelection(t *testing.T) {
	renderer := NewRenderer()

	card := &battle.BattleCard{
		Name:       "Pikachu",
		Stamina:    80,
		StaminaMax: 100,
		Moves: []pokemon.Move{
			{Name: "thunderbolt", Power: 90, StaminaCost: 30, Type: "electric"},
			{Name: "quick-attack", Power: 40, StaminaCost: 15, Type: "normal"},
			{Name: "thunder-wave", Power: 0, StaminaCost: 20, Type: "electric"},
			{Name: "iron-tail", Power: 100, StaminaCost: 40, Type: "steel"},
		},
	}

	result := renderer.RenderMoveSelection(card, 0)
	if result == "" {
		t.Error("RenderMoveSelection returned empty string")
	}
}

func TestStripANSI(t *testing.T) {
	input := ColorRed + "Hello" + ColorReset + " World"
	expected := "Hello World"
	result := stripANSI(input)
	if result != expected {
		t.Errorf("stripANSI(%q) = %q, want %q", input, result, expected)
	}
}

func TestWordWrap(t *testing.T) {
	text := "This is a long message that needs to be wrapped to fit within a certain width"
	lines := wordWrap(text, 20)
	
	if len(lines) == 0 {
		t.Error("wordWrap returned no lines")
	}

	for _, line := range lines {
		if len(line) > 20 {
			t.Errorf("Line exceeds width: %q (len=%d)", line, len(line))
		}
	}
}
