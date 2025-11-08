package auth

// This file provides an example of how to integrate the authentication system
// into the main server application.

/*
Example integration in server/main.go:

package main

import (
	"log"
	"os"
	"pokemon-cli/auth"
	"pokemon-cli/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Initialize auth services
	authService := auth.NewService()
	jwtService, err := auth.NewJWTService()
	if err != nil {
		log.Fatalf("Failed to initialize JWT service: %v", err)
	}

	// Initialize repositories
	userRepo := database.NewUserRepository(database.GetDB())

	// Initialize handlers
	authHandler := auth.NewHandler(authService, jwtService, userRepo)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ORIGINS"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Register auth routes
	auth.RegisterRoutes(app, authHandler, jwtService)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"service": "poketactix-api",
		})
	})

	// Protected route example
	app.Get("/api/protected", auth.Middleware(jwtService), func(c *fiber.Ctx) error {
		userID, _ := auth.GetUserID(c)
		username, _ := auth.GetUsername(c)
		
		return c.JSON(fiber.Map{
			"message": "This is a protected route",
			"user_id": userID,
			"username": username,
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
*/

// Example of using auth middleware in route groups:
/*
func setupRoutes(app *fiber.App, jwtService *auth.JWTService) {
	api := app.Group("/api")
	
	// Public routes
	public := api.Group("/public")
	public.Get("/pokemon", getPokemonHandler)
	
	// Protected routes
	protected := api.Group("/protected", auth.Middleware(jwtService))
	protected.Get("/profile", getProfileHandler)
	protected.Post("/battle/start", startBattleHandler)
	protected.Get("/cards", getCardsHandler)
}
*/

// Example of extracting user info in handlers:
/*
func someProtectedHandler(c *fiber.Ctx) error {
	userID, ok := auth.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}
	
	username, _ := auth.GetUsername(c)
	
	// Use userID and username in your handler logic
	return c.JSON(fiber.Map{
		"user_id": userID,
		"username": username,
	})
}
*/
