package auth

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *Handler, jwtService *JWTService) {
	auth := app.Group("/api/auth")

	auth.Post("/register", RegisterRateLimiter(), handler.Register)

	auth.Post("/login", LoginRateLimiter(), handler.Login)

	auth.Get("/me", Middleware(jwtService), handler.GetCurrentUser)
	auth.Post("/logout", Middleware(jwtService), handler.Logout)
}
