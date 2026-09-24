package cards

import (
	"context"
	"errors"
	"fmt"
	"pokemon-cli/internal/auth"
	"pokemon-cli/internal/items"
	"strconv"
	"strings"

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

// parseCardID extracts and validates the :id route parameter
func parseCardID(c *fiber.Ctx) (int, bool) {
	cardID, err := strconv.Atoi(c.Params("id"))
	if err != nil || cardID <= 0 {
		return 0, false
	}
	return cardID, true
}

// GetEvolution handles GET /api/cards/:id/evolution
func (h *Handler) GetEvolution(c *fiber.Ctx) error {
	userID, ok := auth.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}

	cardID, ok := parseCardID(c)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid card ID",
			},
		})
	}

	info, err := h.service.GetEvolutionInfo(c.Context(), userID, cardID)
	if err != nil {
		return evolutionErrorResponse(c, err)
	}

	return c.JSON(fiber.Map{
		"evolution": info,
	})
}

// GetEvolutionSummaries handles GET /api/cards/evolution-summary
func (h *Handler) GetEvolutionSummaries(c *fiber.Ctx) error {
	userID, ok := auth.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}

	summaries, err := h.service.GetEvolutionSummaries(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve evolution summaries",
			},
		})
	}

	return c.JSON(fiber.Map{
		"summaries": summaries,
	})
}

// EvolveCard handles POST /api/cards/:id/evolve
func (h *Handler) EvolveCard(c *fiber.Ctx) error {
	userID, ok := auth.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}

	cardID, ok := parseCardID(c)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid card ID",
			},
		})
	}

	var req EvolveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
	}

	req.ItemID = strings.TrimSpace(strings.ToLower(req.ItemID))
	if req.ItemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Item ID is required",
			},
		})
	}

	result, err := h.service.EvolveCardWithItem(c.Context(), userID, cardID, req.ItemID)
	if err != nil {
		return evolutionErrorResponse(c, err)
	}

	return c.JSON(fiber.Map{
		"message":   fmt.Sprintf("%s evolved into %s!", result.EvolvedFrom, result.EvolvedInto),
		"evolution": result,
	})
}

// evolutionErrorResponse maps service errors onto the project's error shape
func evolutionErrorResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrCardNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "CARD_NOT_FOUND",
				"message": "Card not found",
			},
		})
	case errors.Is(err, ErrNotCardOwner):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "FORBIDDEN",
				"message": "You don't have access to this card",
			},
		})
	case errors.Is(err, items.ErrItemNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "ITEM_NOT_FOUND",
				"message": "Item not found",
			},
		})
	case errors.Is(err, ErrNotEligible):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "EVOLUTION_NOT_ELIGIBLE",
				"message": err.Error(),
			},
		})
	case errors.Is(err, items.ErrInsufficientItem):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INSUFFICIENT_ITEM",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrConcurrentEvolution):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "CONCURRENT_EVOLUTION",
				"message": "This Pokemon is already evolving or was just evolved. Reload and try again.",
			},
		})
	default:
		fmt.Printf("Evolution failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Evolution failed",
			},
		})
	}
}
