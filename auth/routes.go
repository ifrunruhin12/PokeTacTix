package auth

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all authentication routes
// Requirements: 2.1, 2.2, 2.3 - Create authentication API endpoints
// Requirement: 14.5 - Add rate limiting to auth endpoints
func RegisterRoutes(app *fiber.App, handler *Handler, jwtService *JWTService) {
	// Create auth group
	auth := app.Group("/api/auth")
	
	// Public routes (no authentication required)
	// Requirement: 14.5 - 3 requests/hour for register
	auth.Post("/register", RegisterRateLimiter(), handler.Register)
	
	// Requirement: 14.5 - 5 requests/minute for login
	auth.Post("/login", LoginRateLimiter(), handler.Login)
	
	// Protected routes (authentication required)
	auth.Get("/me", Middleware(jwtService), handler.GetCurrentUser)
	auth.Post("/logout", Middleware(jwtService), handler.Logout)
}
