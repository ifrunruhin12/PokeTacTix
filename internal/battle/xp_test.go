package battle

import "testing"

func TestCalculateXPForBattle1v1(t *testing.T) {
	tests := []struct {
		name string
		bs   *BattleState
		want int // expected number of cards earning XP
	}{
		{
			name: "damaged winner earns XP",
			bs: &BattleState{
				Mode:       "1v1",
				Winner:     "player",
				PlayerDeck: []BattleCard{{CardID: 11, HP: 40, HPMax: 100}},
			},
			want: 1,
		},
		{
			name: "undamaged winner earns XP via recorded action",
			bs: &BattleState{
				Mode:         "1v1",
				Winner:       "player",
				PlayerDeck:   []BattleCard{{CardID: 12, HP: 100, HPMax: 100}},
				Participants: map[int]bool{12: true},
			},
			want: 1,
		},
		{
			name: "undamaged winner without recorded action falls back to active card",
			bs: &BattleState{
				Mode:            "1v1",
				Winner:          "player",
				PlayerDeck:      []BattleCard{{CardID: 13, HP: 100, HPMax: 100}},
				PlayerActiveIdx: 0,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xpMap := CalculateXPForBattle(tt.bs)
			if len(xpMap) != tt.want {
				t.Errorf("cards earning XP = %d, want %d (map: %v)", len(xpMap), tt.want, xpMap)
			}
			for _, card := range tt.bs.PlayerDeck {
				if xp, ok := xpMap[card.CardID]; !ok || xp <= 0 {
					t.Errorf("card %d earned %v, want positive XP", card.CardID, xp)
				}
			}
		})
	}
}

func TestCalculateXPForBattle5v5(t *testing.T) {
	tests := []struct {
		name string
		bs   *BattleState
		want int
	}{
		{
			name: "all damaged cards earn XP",
			bs: &BattleState{
				Mode:   "5v5",
				Winner: "player",
				PlayerDeck: []BattleCard{
					{CardID: 21, HP: 30, HPMax: 100},
					{CardID: 22, HP: 0, HPMax: 100, IsKnockedOut: true},
					{CardID: 23, HP: 90, HPMax: 100},
					{CardID: 24, HP: 55, HPMax: 100},
					{CardID: 25, HP: 100, HPMax: 100},
				},
			},
			want: 4,
		},
		{
			name: "loss still awards consolation XP to damaged cards",
			bs: &BattleState{
				Mode:   "5v5",
				Winner: "ai",
				PlayerDeck: []BattleCard{
					{CardID: 31, HP: 0, HPMax: 100, IsKnockedOut: true},
					{CardID: 32, HP: 100, HPMax: 100},
					{CardID: 33, HP: 100, HPMax: 100},
					{CardID: 34, HP: 100, HPMax: 100},
					{CardID: 35, HP: 100, HPMax: 100},
				},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xpMap := CalculateXPForBattle(tt.bs)
			if len(xpMap) != tt.want {
				t.Errorf("cards earning XP = %d, want %d (map: %v)", len(xpMap), tt.want, xpMap)
			}
		})
	}
}

func TestProcessMoveRecordsParticipation(t *testing.T) {
	state := &BattleState{
		Mode:            "1v1",
		PlayerDeck:      []BattleCard{{CardID: 41, HP: 100, HPMax: 100, Stamina: 50, StaminaMax: 40, Speed: 20}},
		AIDeck:          []BattleCard{{CardID: 99, HP: 100, HPMax: 100, Stamina: 50, StaminaMax: 40, Speed: 20}},
		PlayerActiveIdx: 0,
		AIActiveIdx:     0,
		WhoseTurn:       "player",
		SacrificeCount:  make(map[int]int),
	}

	// Defend is enough to count as acting — no move list needed
	if _, err := ProcessMove(state, "defend", nil); err != nil {
		t.Fatalf("defend move failed: %v", err)
	}

	if !state.Participants[41] {
		t.Errorf("active card 41 not recorded as participant after defending (map: %v)", state.Participants)
	}
}
