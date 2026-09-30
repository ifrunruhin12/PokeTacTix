package battle

import "testing"

func TestStreakBonusCoins(t *testing.T) {
	tests := []struct {
		name   string
		streak int
		want   int
	}{
		{name: "first win has no bonus", streak: 1, want: 0},
		{name: "second win adds 10", streak: 2, want: 10},
		{name: "third win adds 20", streak: 3, want: 20},
		{name: "fifth win adds 40", streak: 5, want: 40},
		{name: "sixth win caps at 50", streak: 6, want: 50},
		{name: "long streak stays capped", streak: 12, want: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := streakBonusCoins(tt.streak); got != tt.want {
				t.Errorf("streakBonusCoins(%d) = %d, want %d", tt.streak, got, tt.want)
			}
		})
	}
}

func TestStreakTier(t *testing.T) {
	tests := []struct {
		streak int
		want   string
	}{
		{streak: 2, want: "2"},
		{streak: 4, want: "4"},
		{streak: 5, want: "5+"},
		{streak: 9, want: "5+"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := streakTier(tt.streak); got != tt.want {
				t.Errorf("streakTier(%d) = %q, want %q", tt.streak, got, tt.want)
			}
		})
	}
}
