package cards

import (
	"pokemon-cli/internal/auth"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all card management routes
// Requirements: 12.1, 12.2, 12.3 - Create card management API endpoints with authentication
func RegisterRoutes(app *fiber.App, handler *Handler, jwtService *auth.JWTService) {
	// Create cards group with authentication middleware
	cards := app.Group("/api/cards", auth.Middleware(jwtService))
	
	// GET /api/cards - Retrieve user's card collection
	// Requirement: 12.1
	cards.Get("/", handler.GetUserCards)
	
	// GET /api/cards/deck - Get current deck
	// Requirement: 12.2
	cards.Get("/deck", handler.GetUserDeck)
	
	// PUT /api/cards/deck - Update deck (must have exactly 5 cards)
	// Requirement: 12.3
	cards.Put("/deck", handler.UpdateDeck)
	
	// GET /api/cards/:id - Get specific card by ID
	cards.Get("/:id", handler.GetCardByID)
}
