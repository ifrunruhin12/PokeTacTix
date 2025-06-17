package main

import (
	"log"
	"os"
	"pokemon-cli/pokemon"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/pokemon", func(c *fiber.Ctx) error {
		// Allow CORS
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusOK)
		}

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

	log.Printf("Serving on :%s...", port)
	log.Fatal(app.Listen(":" + port))
}