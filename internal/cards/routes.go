package cards

import (
	"pokemon-cli/internal/auth"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *Handler, jwtService *auth.JWTService) {
	cards := app.Group("/api/cards", auth.Middleware(jwtService))

	cards.Get("/", handler.GetUserCards)

	cards.Get("/deck", handler.GetUserDeck)

	cards.Put("/deck", handler.UpdateDeck)

	// GET /api/cards/evolution-summary - Which cards have an evolution path
	// (registered before /:id so the static segment is not captured as a card id)
	cards.Get("/evolution-summary", handler.GetEvolutionSummaries)

	cards.Get("/:id", handler.GetCardByID)

	// GET /api/cards/:id/evolution - View how a card can evolve (method, required item/level)
	cards.Get("/:id/evolution", handler.GetEvolution)

	// POST /api/cards/:id/evolve - Evolve a card using an evolution item
	cards.Post("/:id/evolve", handler.EvolveCard)
}
