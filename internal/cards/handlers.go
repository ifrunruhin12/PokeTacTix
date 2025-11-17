package cards

import (
	"context"
	"pokemon-cli/internal/auth"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles card-related HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new card handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

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
	cards, err := h.service.GetUserCards(ctx, userID)
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
	deck, err := h.service.GetUserDeck(ctx, userID)
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

	// Validate deck has 1-5 cards
	if len(req.CardIDs) < 1 || len(req.CardIDs) > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_DECK",
				"message": "Deck must contain between 1 and 5 cards",
			},
		})
	}

	ctx := context.Background()
	err := h.service.UpdateDeck(ctx, userID, req.CardIDs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_DECK",
				"message": err.Error(),
			},
		})
	}

	// Return updated deck
	deck, err := h.service.GetUserDeck(ctx, userID)
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
	card, err := h.service.repository.GetByID(ctx, cardID)
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
