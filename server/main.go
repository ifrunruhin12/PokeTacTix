package main

import (
	"fmt"
	"log"
	"pokemon-cli/pokemon"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Create a new engine with template functions
	engine := html.New("../client/views", ".html")

	// Add helper functions
	engine.AddFunc("percentage", func(value, max int) float64 {
		return float64(value) / float64(max) * 100
	})

	// Create new Fiber app with template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Enable CORS
	app.Use(cors.New())

	// Serve static files
	app.Static("/", "../client/public")
	app.Static("/assets", "../assets")

	// Home route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "PokeTacTix",
		})
	})

	// Battle setup route
	app.Get("/battle-setup", func(c *fiber.Ctx) error {
		return c.Render("battle", fiber.Map{
			"Title": "Battle Setup",
		})
	})

	// Battle route - handles both 1v1 and 5v5
	app.Get("/battle", func(c *fiber.Ctx) error {
		mode := c.Query("mode")
		playerName := c.Query("player")

		// Validate parameters
		if mode == "" || playerName == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Mode and player name are required",
			})
		}

		if mode != "1v1" && mode != "5v5" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid battle mode. Must be '1v1' or '5v5'",
			})
		}

		// For now, we'll render a placeholder battle page
		// You can implement the actual battle logic later
		return c.Render("battle-arena", fiber.Map{
			"Title":      fmt.Sprintf("%s Battle", mode),
			"Mode":       mode,
			"PlayerName": playerName,
		})
	})

	// Search Pokemon API endpoint
	app.Get("/pokemon", func(c *fiber.Ctx) error {
		name := c.Query("name")
		if name == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Pokemon name is required",
			})
		}

		// Use the existing FetchPokemon function
		poke, moves, err := pokemon.FetchPokemon(name)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Convert to Card for display
		card := pokemon.BuildCardFromPokemon(poke, moves)

		return c.Render("pokemon", fiber.Map{
			"Pokemon": card,
		})
	})

	// Start server
	port := ":3000"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(app.Listen(port))
}

