package main

import (
	"log"
	"os"
	"pokemon-cli/game/models"
	"pokemon-cli/game/web"
	"pokemon-cli/pokemon"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
)

var (
	sessions   = make(map[string]*Session)
	sessionsMu sync.Mutex
)

// Session holds core game state plus web-only turn state
type Session struct {
	State *models.GameState
	Turn  *web.TurnState
}

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

	// POST /battle/start
	app.Post("/battle/start", func(c *fiber.Ctx) error {
		var req struct {
			PlayerName string `json:"playerName"`
			AIName     string `json:"aiName"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		// Use random cards for both player and AI
		playerCard := pokemon.FetchRandomPokemonCard(false)
		aiCard := pokemon.FetchRandomPokemonCard(false)
		player := models.NewPlayer(req.PlayerName, []pokemon.Card{playerCard})
		ai := models.NewPlayer(req.AIName, []pokemon.Card{aiCard})
		state := &models.GameState{
			Player:          player,
			AI:              ai,
			BattleMode:      "1v1",
			PlayerActiveIdx: 0,
			AIActiveIdx:     0,
			BattleStarted:   true,
			InBattle:        true,
			TurnNumber:      1,
		}
		turn := &web.TurnState{WhoseTurn: "player"}
		id := uuid.New().String()
		sessionsMu.Lock()
		sessions[id] = &Session{State: state, Turn: turn}
		sessionsMu.Unlock()
		return c.JSON(fiber.Map{"session": id, "state": state, "turn": turn})
	})

	// POST /battle/move
	app.Post("/battle/move", func(c *fiber.Ctx) error {
		var req struct {
			Session string `json:"session"`
			Move    string `json:"move"`
			MoveIdx *int   `json:"moveIdx"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		sessionsMu.Lock()
		sess, ok := sessions[req.Session]
		sessionsMu.Unlock()
		if !ok {
			return c.Status(404).JSON(fiber.Map{"error": "Session not found"})
		}
		// Call a function to process the move and update state
		result, err := web.ProcessWebMove(sess.State, sess.Turn, req.Move, req.MoveIdx)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	// GET /battle/state?session=...
	app.Get("/battle/state", func(c *fiber.Ctx) error {
		session := c.Query("session")
		sessionsMu.Lock()
		sess, ok := sessions[session]
		sessionsMu.Unlock()
		if !ok {
			return c.Status(404).JSON(fiber.Map{"error": "Session not found"})
		}
		return c.JSON(fiber.Map{"state": sess.State, "turn": sess.Turn})
	})

	log.Printf("Serving on :%s...", port)
	log.Fatal(app.Listen(":" + port))
}
