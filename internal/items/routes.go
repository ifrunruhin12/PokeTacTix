package items

import (
	"pokemon-cli/internal/auth"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers item routes
func RegisterRoutes(app *fiber.App, handler *Handler, jwtService *auth.JWTService) {
	items := app.Group("/api/items", auth.Middleware(jwtService))

	// GET /api/items/inventory - Get the player's item inventory
	items.Get("/inventory", handler.GetInventory)

	// POST /api/items/:id/use - Activate a booster from the inventory
	items.Post("/:id/use", handler.UseItem)
}
