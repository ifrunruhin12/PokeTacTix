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

	cards.Get("/:id", handler.GetCardByID)
}
