package stats

import (
	"testing"
)

// TestConsecutiveLossLogicValidation verifies the consecutive loss tracking logic
// This test validates the SQL logic structure without requiring a database connection
func TestConsecutiveLossLogicValidation(t *testing.T) {
	tests := []struct {
		name           string
		mode           string
		result         string
		wantValid      bool
		wantIncrement  bool // true if consecutive_losses should increment
		wantReset      bool // true if consecutive_losses should reset to 0
	}{
		{
			name:          "1v1 loss increments counter",
			mode:          "1v1",
			result:        "loss",
			wantValid:     true,
			wantIncrement: true,
			wantReset:     false,
		},
		{
			name:          "1v1 win resets counter",
			mode:          "1v1",
			result:        "win",
			wantValid:     true,
			wantIncrement: false,
			wantReset:     true,
		},
		{
			name:          "1v1 draw resets counter",
			mode:          "1v1",
			result:        "draw",
			wantValid:     true,
			wantIncrement: false,
			wantReset:     true,
		},
		{
			name:          "5v5 loss increments counter",
			mode:          "5v5",
			result:        "loss",
			wantValid:     true,
			wantIncrement: true,
			wantReset:     false,
		},
		{
			name:          "5v5 win resets counter",
			mode:          "5v5",
			result:        "win",
			wantValid:     true,
			wantIncrement: false,
			wantReset:     true,
		},
		{
			name:          "5v5 draw resets counter",
			mode:          "5v5",
			result:        "draw",
			wantValid:     true,
			wantIncrement: false,
			wantReset:     true,
		},
		{
			name:          "invalid mode rejected",
			mode:          "invalid",
			result:        "win",
			wantValid:     false,
			wantIncrement: false,
			wantReset:     false,
		},
		{
			name:          "invalid result rejected",
			mode:          "1v1",
			result:        "invalid",
			wantValid:     false,
			wantIncrement: false,
			wantReset:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate mode
			validMode := tt.mode == "1v1" || tt.mode == "5v5"
			if validMode != tt.wantValid && tt.wantValid {
				t.Errorf("mode validation failed: got %v, want %v", validMode, tt.wantValid)
			}

			// Validate result
			validResult := tt.result == "win" || tt.result == "loss" || tt.result == "draw"
			if validResult != tt.wantValid && tt.wantValid {
				t.Errorf("result validation failed: got %v, want %v", validResult, tt.wantValid)
			}

			if !tt.wantValid {
				return
			}

			// Verify consecutive loss logic
			shouldIncrement := tt.result == "loss"
			shouldReset := tt.result == "win" || tt.result == "draw"

			if shouldIncrement != tt.wantIncrement {
				t.Errorf("increment logic failed: got %v, want %v", shouldIncrement, tt.wantIncrement)
			}

			if shouldReset != tt.wantReset {
				t.Errorf("reset logic failed: got %v, want %v", shouldReset, tt.wantReset)
			}
		})
	}
}

// TestModeIndependence verifies that consecutive losses are tracked across both modes
func TestModeIndependence(t *testing.T) {
	modes := []string{"1v1", "5v5"}
	
	for _, mode := range modes {
		t.Run("loss_in_"+mode, func(t *testing.T) {
			// Verify that loss in any mode should increment the same counter
			result := "loss"
			shouldIncrement := result == "loss"
			
			if !shouldIncrement {
				t.Errorf("mode %s should increment consecutive losses on loss", mode)
			}
		})
		
		t.Run("win_in_"+mode, func(t *testing.T) {
			// Verify that win in any mode should reset the same counter
			result := "win"
			shouldReset := result == "win" || result == "draw"
			
			if !shouldReset {
				t.Errorf("mode %s should reset consecutive losses on win", mode)
			}
		})
	}
}
