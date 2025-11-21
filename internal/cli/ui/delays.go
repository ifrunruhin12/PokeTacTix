package ui

import (
	"time"
)

// GetBattleDelay returns the delay duration based on battle speed setting
// delayType can be: "short", "medium", "long"
func GetBattleDelay(battleSpeed string, delayType string) time.Duration {
	// Default to normal if not set
	if battleSpeed == "" {
		battleSpeed = "normal"
	}

	// Define base delays for each type
	var baseDuration time.Duration
	switch delayType {
	case "short":
		baseDuration = 500 * time.Millisecond
	case "medium":
		baseDuration = 1000 * time.Millisecond
	case "long":
		baseDuration = 2000 * time.Millisecond
	default:
		baseDuration = 1000 * time.Millisecond
	}

	// Apply speed multiplier
	switch battleSpeed {
	case "slow":
		return time.Duration(float64(baseDuration) * 1.5)
	case "fast":
		return time.Duration(float64(baseDuration) * 0.5)
	default: // "normal"
		return baseDuration
	}
}

// Sleep pauses execution based on battle speed setting
func Sleep(battleSpeed string, delayType string) {
	time.Sleep(GetBattleDelay(battleSpeed, delayType))
}
