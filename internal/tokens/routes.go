package tokens

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers token routes
func RegisterRoutes(app *fiber.App, handler *Handler, authMiddleware fiber.Handler) {
	tokens := app.Group("/api/tokens")

	// Apply authentication middleware to all token routes
	tokens.Use(authMiddleware)

	// GET /api/tokens/balance - Get current token balance and info
	tokens.Get("/balance", handler.GetTokenBalanceHandler)
}
