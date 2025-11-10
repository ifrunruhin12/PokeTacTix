package battle

import (
	"pokemon-cli/game/models"
	"pokemon-cli/internal/pokemon"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Handler handles battle-related HTTP requests
type Handler struct {
	sessions   map[string]*Session
}

// Session holds core game state plus web-only turn state
type Session struct {
	State *models.GameState
	Turn  *TurnState
}

// NewHandler creates a new battle handler
func NewHandler() *Handler {
	return &Handler{
		sessions: make(map[string]*Session),
	}
}

// StartBattle handles POST /api/battle/start
func (h *Handler) StartBattle(c *fiber.Ctx) error {
	var req struct {
		PlayerName string `json:"playerName"`
		AIName     string `json:"aiName"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
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

	turn := &TurnState{WhoseTurn: "player"}
	id := uuid.New().String()

	h.sessions[id] = &Session{State: state, Turn: turn}

	return c.JSON(fiber.Map{
		"session": id,
		"state":   state,
		"turn":    turn,
	})
}

// MakeMove handles POST /api/battle/move
func (h *Handler) MakeMove(c *fiber.Ctx) error {
	var req struct {
		Session string `json:"session"`
		Move    string `json:"move"`
		MoveIdx *int   `json:"moveIdx"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	sess, ok := h.sessions[req.Session]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Session not found"})
	}

	// Call ProcessWebMove to process the move and update state
	result, err := ProcessWebMove(sess.State, sess.Turn, req.Move, req.MoveIdx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// GetBattleState handles GET /api/battle/state
func (h *Handler) GetBattleState(c *fiber.Ctx) error {
	session := c.Query("session")
	sess, ok := h.sessions[session]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Session not found"})
	}

	return c.JSON(fiber.Map{
		"state": sess.State,
		"turn":  sess.Turn,
	})
}
