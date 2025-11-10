package cards

// This file provides an example of how to integrate the card management system
// into the main server application.

/*
Example integration in server/main.go:

package main

import (
	"log"
	"os"
	"pokemon-cli/auth"
	"pokemon-cli/cards"
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

	// Initialize services
	authService := auth.NewService()
	jwtService, err := auth.NewJWTService()
	if err != nil {
		log.Fatalf("Failed to initialize JWT service: %v", err)
	}

	// Initialize repositories
	userRepo := database.NewUserRepository(database.GetDB())
	cardRepo := database.NewCardRepository(database.GetDB())

	// Initialize card service
	cardService := database.NewCardService(cardRepo)

	// Initialize handlers
	authHandler := auth.NewHandler(authService, jwtService, userRepo, cardService)
	cardHandler := cards.NewHandler(cardService)

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

	// Register routes
	auth.RegisterRoutes(app, authHandler, jwtService)
	cards.RegisterRoutes(app, cardHandler, jwtService)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"service": "poketactix-api",
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

// Example of using card endpoints from frontend:
/*
// Get all cards
fetch('/api/cards', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
})
.then(res => res.json())
.then(data => console.log(data.cards));

// Get current deck
fetch('/api/cards/deck', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
})
.then(res => res.json())
.then(data => console.log(data.deck));

// Update deck
fetch('/api/cards/deck', {
  method: 'PUT',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    card_ids: [1, 2, 3, 4, 5]
  })
})
.then(res => res.json())
.then(data => console.log(data.message, data.deck));
*/

// Example of card leveling after battle:
/*
func awardBattleRewards(ctx context.Context, cardService *database.CardService, cardIDs []int, xpAmount int) error {
	for _, cardID := range cardIDs {
		updatedCard, err := cardService.AddXP(ctx, cardID, xpAmount)
		if err != nil {
			return err
		}
		
		// Check if card leveled up
		if updatedCard.Level > previousLevel {
			log.Printf("Card %s leveled up to %d!", updatedCard.PokemonName, updatedCard.Level)
		}
	}
	return nil
}
*/

// Example of checking card stats at current level:
/*
func displayCardStats(card *database.PlayerCard) {
	stats := card.GetCurrentStats()
	fmt.Printf("Card: %s (Level %d)\n", card.PokemonName, card.Level)
	fmt.Printf("HP: %d\n", stats.HP)
	fmt.Printf("Attack: %d\n", stats.Attack)
	fmt.Printf("Defense: %d\n", stats.Defense)
	fmt.Printf("Speed: %d\n", stats.Speed)
	fmt.Printf("Stamina: %d\n", stats.Stamina)
	fmt.Printf("XP: %d / %d\n", card.XP, 100*card.Level)
}
*/
