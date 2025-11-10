package cards

import (
	"context"
	"pokemon-cli/auth"
	"pokemon-cli/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles card-related HTTP requests
type Handler struct {
	cardService *database.CardService
}

// NewHandler creates a new card handler
func NewHandler(cardService *database.CardService) *Handler {
	return &Handler{
		cardService: cardService,
	}
}

// GetUserCards retrieves all cards for the authenticated user
// Requirement: 12.1 - Implement GET /api/cards endpoint
func (h *Handler) GetUserCards(c *fiber.Ctx) error {
	userID, ok := auth.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}

	ctx := context.Background()
	cards, err := h.cardService.GetUserCards(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve cards",
			},
		})
	}

	return c.JSON(fiber.Map{
		"cards": cards,
	})
}

// GetUserDeck retrieves the user's current deck (5 cards)
// Requirement: 12.2 - Implement GET /api/cards/deck endpoint
func (h *Handler) GetUserDeck(c *fiber.Ctx) error {
	userID, ok := auth.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}

	ctx := context.Background()
	deck, err := h.cardService.GetUserDeck(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve deck",
			},
		})
	}

	return c.JSON(fiber.Map{
		"deck": deck,
	})
}

// UpdateDeckRequest represents the request body for updating a deck
type UpdateDeckRequest struct {
	CardIDs []int `json:"card_ids"`
}

// UpdateDeck updates the user's deck configuration
// Requirement: 12.3 - Implement PUT /api/cards/deck endpoint (must have exactly 5 cards)
func (h *Handler) UpdateDeck(c *fiber.Ctx) error {
	userID, ok := auth.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}

	var req UpdateDeckRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
	}

	// Validate deck has exactly 5 cards
	if len(req.CardIDs) != 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_DECK",
				"message": "Deck must contain exactly 5 cards",
			},
		})
	}

	ctx := context.Background()
	err := h.cardService.UpdateDeck(ctx, userID, req.CardIDs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_DECK",
				"message": err.Error(),
			},
		})
	}

	// Return updated deck
	deck, err := h.cardService.GetUserDeck(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve updated deck",
			},
		})
	}

	return c.JSON(fiber.Map{
		"message": "Deck updated successfully",
		"deck":    deck,
	})
}

// GetCardByID retrieves a specific card by ID
func (h *Handler) GetCardByID(c *fiber.Ctx) error {
	userID, ok := auth.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}

	cardIDStr := c.Params("id")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid card ID",
			},
		})
	}

	ctx := context.Background()
	
	// Get the card
	cardRepo := database.NewCardRepository(database.GetDB())
	card, err := cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "CARD_NOT_FOUND",
				"message": "Card not found",
			},
		})
	}

	// Verify the card belongs to the user
	if card.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "FORBIDDEN",
				"message": "You don't have access to this card",
			},
		})
	}

	return c.JSON(fiber.Map{
		"card": card,
	})
}
