package battle

import "testing"

func TestSacrificeCountsAreSeparatedBySide(t *testing.T) {
	state := &BattleState{
		Mode:            "1v1",
		PlayerDeck:      []BattleCard{{HP: 100, HPMax: 100, Speed: 20}},
		AIDeck:          []BattleCard{{HP: 100, HPMax: 100, Speed: 20}},
		PlayerActiveIdx: 0,
		AIActiveIdx:     0,
		WhoseTurn:       "player",
		SacrificeCount:  make(map[int]int),
	}

	if _, err := ProcessMove(state, "sacrifice", nil); err != nil {
		t.Fatalf("player sacrifice failed: %v", err)
	}
	processAIMove(state)

	if got := state.SacrificeCount[playerSacrificeKey(0)]; got != 1 {
		t.Errorf("player sacrifice count = %d, want 1", got)
	}
	if got := state.SacrificeCount[aiSacrificeKey(0)]; got != 1 {
		t.Errorf("AI sacrifice count = %d, want 1", got)
	}
}
