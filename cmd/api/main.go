package main

import (
	"os"
	"pokemon-cli/internal/auth"
	"pokemon-cli/internal/battle"
	"pokemon-cli/internal/cards"
	"pokemon-cli/internal/database"
	"pokemon-cli/internal/pokemon"
	"pokemon-cli/internal/shop"
	"pokemon-cli/pkg/config"
	"pokemon-cli/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger with slog
	logLevel := logger.INFO
	if cfg.Server.Env == "development" {
		logLevel = logger.DEBUG
	}
	
	// Use text handler for development, JSON for production
	var appLogger *logger.Logger
	if cfg.Server.Env == "development" {
		appLogger = logger.NewText(logLevel)
	} else {
		appLogger = logger.New(logLevel)
	}

	appLogger.Info("Starting PokeTacTix API", "env", cfg.Server.Env, "port", cfg.Server.Port)

	// Initialize database
	if cfg.Database.URL != "" {
		if err := database.InitDB(&cfg.Database); err != nil {
			appLogger.Error("Failed to initialize database", "error", err)
			os.Exit(1)
		}
		defer database.CloseDB()
		appLogger.Info("Database connection established")
	} else {
		appLogger.Warn("Database URL not set, running without database")
	}

	// Initialize JWT service
	jwtService, err := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)
	if err != nil {
		appLogger.Error("Failed to initialize JWT service", "error", err)
		os.Exit(1)
	}

	// Initialize repositories
	authRepo := auth.NewRepository(database.GetDB())
	cardsRepo := cards.NewRepository(database.GetDB())
	shopRepo := shop.NewRepository(database.GetDB())

	// Initialize services
	authService := auth.NewService()
	cardsService := cards.NewService(cardsRepo)
	shopService := shop.NewService()

	// Initialize handlers
	authHandler := auth.NewHandler(authService, jwtService, authRepo, cardsService)
	cardsHandler := cards.NewHandler(cardsService)
	battleHandler := battle.NewHandler()
	shopHandler := shop.NewHandler(shopService, shopRepo)

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
	app.Use(cors.New(cors.Config{
		AllowOrigins: func() string {
			if len(cfg.CORS.AllowedOrigins) > 0 {
				origins := ""
				for i, origin := range cfg.CORS.AllowedOrigins {
					if i > 0 {
						origins += ","
					}
					origins += origin
				}
				return origins
			}
			return "*"
		}(),
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"env":    cfg.Server.Env,
		})
	})

	// Pokemon endpoint (legacy support)
	app.Get("/pokemon", func(c *fiber.Ctx) error {
		name := c.Query("name")
		if name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing name parameter",
			})
		}

		poke, moves, err := pokemon.FetchPokemon(name)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		card := pokemon.BuildCardFromPokemon(poke, moves)
		return c.JSON(card)
	})

	// Register routes
	auth.RegisterRoutes(app, authHandler, jwtService)
	cards.RegisterRoutes(app, cardsHandler, jwtService)
	
	// Create auth middleware for protected routes
	authMiddleware := auth.Middleware(jwtService)
	battle.RegisterRoutes(app, battleHandler, authMiddleware)
	shop.RegisterRoutes(app, shopHandler, authMiddleware)

	// Start server
	port := cfg.Server.Port
	appLogger.Info("Server starting", "port", port)
	if err := app.Listen(":" + port); err != nil {
		appLogger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
