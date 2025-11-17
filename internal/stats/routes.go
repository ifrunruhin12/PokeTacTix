package stats

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *Handler, authMiddleware fiber.Handler) {
	profile := app.Group("/api/profile")

	profile.Use(authMiddleware)

	profile.Get("/stats", handler.GetPlayerStats)

	profile.Get("/history", handler.GetBattleHistory)

	profile.Get("/achievements", handler.GetAchievements)

	profile.Post("/achievements/check", handler.CheckAchievements)
}
