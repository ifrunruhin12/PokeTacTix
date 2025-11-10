package battle

import (
	"fmt"
	"pokemon-cli/game/models"
	"pokemon-cli/internal/pokemon"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler handles battle-related HTTP requests
type Handler struct {
	sessions   map[string]*Session
	battleStates map[string]*BattleState // Enhanced battle state storage
	mu         sync.RWMutex             // Mutex for thread-safe access
}

// Session holds core game state plus web-only turn state (legacy support)
type Session struct {
	State *models.GameState
	Turn  *TurnState
}

// NewHandler creates a new battle handler
func NewHandler() *Handler {
	return &Handler{
		sessions:     make(map[string]*Session),
		battleStates: make(map[string]*BattleState),
	}
}

// GetBattleState retrieves a battle state by ID (thread-safe)
func (h *Handler) GetBattleState(id string) (*BattleState, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	state, ok := h.battleStates[id]
	return state, ok
}

// SaveBattleState stores a battle state (thread-safe)
func (h *Handler) SaveBattleState(state *BattleState) {
	h.mu.Lock()
	defer h.mu.Unlock()
	state.UpdatedAt = time.Now()
	h.battleStates[state.ID] = state
}

// DeleteBattleState removes a battle state (thread-safe)
func (h *Handler) DeleteBattleState(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.battleStates, id)
}

// StartBattleEnhanced handles POST /api/battle/start with mode selection
func (h *Handler) StartBattleEnhanced(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}

	var req struct {
		Mode string `json:"mode"` // "1v1" or "5v5"
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
	}

	// Sanitize and validate mode
	req.Mode = strings.TrimSpace(req.Mode)
	if req.Mode != "1v1" && req.Mode != "5v5" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_MODE",
				"message": "Invalid battle mode. Must be '1v1' or '5v5'",
			},
		})
	}

	// TODO: Fetch player's deck from database
	// For now, use random cards
	var playerDeck []pokemon.Card
	var aiDeck []pokemon.Card

	if req.Mode == "1v1" {
		playerDeck = []pokemon.Card{pokemon.FetchRandomPokemonCard(false)}
		aiDeck = []pokemon.Card{pokemon.FetchRandomPokemonCard(false)}
	} else {
		// Generate 5 random cards for each side
		for i := 0; i < 5; i++ {
			playerDeck = append(playerDeck, pokemon.FetchRandomPokemonCard(false))
			aiDeck = append(aiDeck, pokemon.FetchRandomPokemonCard(false))
		}
	}

	// Start the battle
	battleState, err := StartBattle(userID, req.Mode, playerDeck, aiDeck)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save battle state
	h.SaveBattleState(battleState)

	// Return battle state with card visibility
	response := BuildBattleResponse(battleState, []string{fmt.Sprintf("Battle started! Mode: %s", req.Mode)}, true)

	return c.JSON(response)
}

// StartBattle handles POST /api/battle/start (legacy support)
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

// MakeMoveEnhanced handles POST /api/battle/move with enhanced battle system
func (h *Handler) MakeMoveEnhanced(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}

	var req struct {
		BattleID string `json:"battle_id"`
		Move     string `json:"move"`
		MoveIdx  *int   `json:"move_idx"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
	}

	// Sanitize and validate inputs
	req.BattleID = strings.TrimSpace(req.BattleID)
	req.Move = strings.TrimSpace(req.Move)
	
	if req.BattleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Battle ID is required",
			},
		})
	}
	
	// Validate move
	validMoves := map[string]bool{
		"attack": true, "defend": true, "pass": true, "sacrifice": true, "surrender": true,
	}
	if !validMoves[req.Move] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_MOVE",
				"message": "Invalid move. Must be one of: attack, defend, pass, sacrifice, surrender",
			},
		})
	}
	
	// Validate move_idx for attack moves
	if req.Move == "attack" && (req.MoveIdx == nil || *req.MoveIdx < 0) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_MOVE",
				"message": "Move index is required for attack moves",
			},
		})
	}

	// Get battle state
	battleState, ok := h.GetBattleState(req.BattleID)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "BATTLE_NOT_FOUND",
				"message": "Battle not found",
			},
		})
	}

	// Verify user owns this battle
	if battleState.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "FORBIDDEN",
				"message": "Not your battle",
			},
		})
	}

	// Process the move
	logEntries, err := ProcessMove(battleState, req.Move, req.MoveIdx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save updated battle state
	h.SaveBattleState(battleState)

	// Return updated state with logs
	response := BuildBattleResponse(battleState, logEntries, true)

	return c.JSON(response)
}

// GetBattleStateEnhanced handles GET /api/battle/state with enhanced system
func (h *Handler) GetBattleStateEnhanced(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	battleID := c.Query("battle_id")
	if battleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "battle_id required"})
	}

	// Get battle state
	battleState, ok := h.GetBattleState(battleID)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Battle not found"})
	}

	// Verify user owns this battle
	if battleState.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not your battle"})
	}

	// Return state
	response := BuildBattleResponse(battleState, []string{}, true)

	return c.JSON(response)
}

// SwitchPokemonHandler handles POST /api/battle/switch
func (h *Handler) SwitchPokemonHandler(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		BattleID string `json:"battle_id"`
		NewIdx   int    `json:"new_idx"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Get battle state
	battleState, ok := h.GetBattleState(req.BattleID)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Battle not found"})
	}

	// Verify user owns this battle
	if battleState.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not your battle"})
	}

	// Switch Pokemon
	err := SwitchPokemon(battleState, req.NewIdx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save updated battle state
	h.SaveBattleState(battleState)

	// Return updated state
	response := BuildBattleResponse(battleState, []string{fmt.Sprintf("Switched to %s", battleState.PlayerDeck[req.NewIdx].Name)}, true)

	return c.JSON(response)
}

// GetBattleState handles GET /api/battle/state (legacy support)
func (h *Handler) GetBattleStateLegacy(c *fiber.Ctx) error {
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

// SelectRewardHandler handles POST /api/battle/select-reward
// Requirements: 11.1, 11.2, 11.3
func (h *Handler) SelectRewardHandler(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		BattleID      string `json:"battle_id"`
		PokemonIndex  int    `json:"pokemon_index"` // Index of AI Pokemon to select (0-4)
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Get battle state
	battleState, ok := h.GetBattleState(req.BattleID)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Battle not found"})
	}

	// Verify user owns this battle
	if battleState.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not your battle"})
	}

	// Validate battle was won by player and is 5v5 mode
	if battleState.Mode != "5v5" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Reward selection only available for 5v5 battles"})
	}

	if battleState.Winner != "player" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Can only select reward after winning the battle"})
	}

	if !battleState.BattleOver {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Battle is not over yet"})
	}

	// Check if reward already claimed
	if battleState.RewardClaimed {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Reward already claimed for this battle"})
	}

	// Validate Pokemon index
	if req.PokemonIndex < 0 || req.PokemonIndex >= len(battleState.AIDeck) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Pokemon index"})
	}

	// Get the selected AI Pokemon
	selectedPokemon := battleState.AIDeck[req.PokemonIndex]

	// Get database connection from context
	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection not available"})
	}

	// Add the Pokemon to player's collection
	addedCard, err := AddAIPokemonToCollection(c.Context(), db, userID, selectedPokemon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to add Pokemon to collection: %v", err)})
	}

	// Mark battle as reward claimed
	battleState.RewardClaimed = true
	h.SaveBattleState(battleState)

	// Return success with the added card
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Successfully added %s to your collection!", selectedPokemon.Name),
		"card":    addedCard,
	})
}
