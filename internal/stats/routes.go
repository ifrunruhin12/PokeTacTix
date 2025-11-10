package stats

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all stats and profile routes
// Requirements: 8.1, 8.3, 8.4
func RegisterRoutes(app *fiber.App, handler *Handler, authMiddleware fiber.Handler) {
	// Profile routes group
	profile := app.Group("/api/profile")
	
	// Apply authentication middleware to all profile routes
	profile.Use(authMiddleware)
	
	// GET /api/profile/stats - Get player statistics
	profile.Get("/stats", handler.GetPlayerStats)
	
	// GET /api/profile/history - Get battle history (last 20 battles by default)
	profile.Get("/history", handler.GetBattleHistory)
	
	// GET /api/profile/achievements - Get all achievements with unlock status
	profile.Get("/achievements", handler.GetAchievements)
	
	// POST /api/profile/achievements/check - Check and unlock new achievements
	profile.Post("/achievements/check", handler.CheckAchievements)
}
