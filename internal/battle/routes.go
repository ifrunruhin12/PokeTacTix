package battle

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers battle-related routes
func RegisterRoutes(app *fiber.App, handler *Handler, authMiddleware func(*fiber.Ctx) error) {
	battle := app.Group("/api/battle")

	// Legacy routes (no auth required for backward compatibility)
	battle.Post("/start-legacy", handler.StartBattle)
	battle.Post("/move-legacy", handler.MakeMove)
	battle.Get("/state-legacy", handler.GetBattleStateLegacy)

	// Main routes with authentication (enhanced battle system)
	battleAuth := battle.Group("", authMiddleware)
	battleAuth.Post("/start", handler.StartBattleEnhanced)
	battleAuth.Post("/move", handler.MakeMoveEnhanced)
	battleAuth.Get("/state", handler.GetBattleStateEnhanced)
	battleAuth.Post("/switch", handler.SwitchPokemonHandler)
	battleAuth.Post("/select-reward", handler.SelectRewardHandler)

	// Cleanup endpoint (can be called by cron job or admin)
	battle.Post("/cleanup-sessions", handler.CleanupExpiredSessions)
}
