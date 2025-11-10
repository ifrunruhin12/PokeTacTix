package main

import (
	"errors"
	"pokemon-cli/pkg/logger"
)

func main() {
	// Example 1: Text format (development)
	println("=== Development Mode (Text Format) ===")
	devLogger := logger.NewText(logger.DEBUG)
	
	devLogger.Debug("Debug message", "user_id", 123, "action", "login")
	devLogger.Info("Server starting", "port", 3000, "env", "development")
	devLogger.Warn("Rate limit approaching", "current", 95, "limit", 100)
	devLogger.Error("Database error", "error", errors.New("connection timeout"), "retry", 3)
	
	println("\n=== Production Mode (JSON Format) ===")
	// Example 2: JSON format (production)
	prodLogger := logger.New(logger.INFO)
	
	prodLogger.Info("Server starting", "port", 3000, "env", "production")
	prodLogger.Warn("Rate limit approaching", "current", 95, "limit", 100)
	prodLogger.Error("Database error", "error", errors.New("connection timeout"), "retry", 3)
	
	println("\n=== With Context ===")
	// Example 3: Logger with context
	userLogger := devLogger.With("user_id", 456, "username", "pikachu_trainer")
	
	userLogger.Info("User logged in")
	userLogger.Debug("Fetching user deck")
	userLogger.Info("Battle started", "mode", "5v5", "opponent", "ai")
	
	println("\n=== Structured Fields ===")
	// Example 4: Complex structured data
	battleLogger := devLogger.With("battle_id", "abc-123", "mode", "5v5")
	
	battleLogger.Info("Turn processed",
		"turn_number", 5,
		"player_move", "attack",
		"ai_move", "defend",
		"damage_dealt", 45,
		"hp_remaining", 78,
	)
	
	battleLogger.Error("Move validation failed",
		"error", errors.New("insufficient stamina"),
		"required", 30,
		"available", 15,
		"pokemon", "Charizard",
	)
}
