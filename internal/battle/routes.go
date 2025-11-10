package battle

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers battle-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	battle := app.Group("/api/battle")

	battle.Post("/start", handler.StartBattle)
	battle.Post("/move", handler.MakeMove)
	battle.Get("/state", handler.GetBattleState)
}
