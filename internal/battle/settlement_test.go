package battle

import (
	"encoding/json"
	"testing"
)

// TestNeedsSettlement verifies when a finished battle still owes its payout.
func TestNeedsSettlement(t *testing.T) {
	tests := []struct {
		name  string
		state *BattleState
		want  bool
	}{
		{
			name:  "ongoing battle never settles",
			state: &BattleState{BattleOver: false, RewardsSettled: false},
			want:  false,
		},
		{
			name:  "finished unsettled battle settles",
			state: &BattleState{BattleOver: true, RewardsSettled: false},
			want:  true,
		},
		{
			name:  "finished settled battle does not settle again",
			state: &BattleState{BattleOver: true, RewardsSettled: true},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.needsSettlement(); got != tt.want {
				t.Errorf("needsSettlement() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRewardsSettledSurvivesRoundtrip verifies the settled flag persists
// through the same marshal/unmarshal cycle battle sessions go through.
func TestRewardsSettledSurvivesRoundtrip(t *testing.T) {
	state := &BattleState{
		ID:             "test-battle",
		Mode:           "1v1",
		BattleOver:     true,
		Winner:         "player",
		RewardsSettled: true,
		SacrificeCount: map[int]int{},
	}

	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("failed to marshal battle state: %v", err)
	}

	var restored BattleState
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal battle state: %v", err)
	}

	if !restored.RewardsSettled {
		t.Error("RewardsSettled did not survive JSON roundtrip")
	}
	if restored.needsSettlement() {
		t.Error("settled battle must not be settleable again after roundtrip")
	}
}
