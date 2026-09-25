package items

import (
	"errors"
	"fmt"
	"pokemon-cli/internal/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Handler handles item-related HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new items handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetInventory handles GET /api/items/inventory
func (h *Handler) GetInventory(c *fiber.Ctx) error {
	userID, ok := auth.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}

	inventory, err := h.service.GetInventory(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve inventory",
			},
		})
	}

	boosts, err := h.service.ActiveBoosts(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve active boosts",
			},
		})
	}

	return c.JSON(fiber.Map{
		"inventory":     inventory,
		"active_boosts": boosts,
	})
}

// UseItem handles POST /api/items/:id/use — activate a booster from the
// player's inventory.
func (h *Handler) UseItem(c *fiber.Ctx) error {
	userID, ok := auth.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}

	itemID := strings.TrimSpace(strings.ToLower(c.Params("id")))
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Item id is required",
			},
		})
	}

	boost, err := h.service.UseItem(c.Context(), userID, itemID)
	if err != nil {
		switch {
		case errors.Is(err, ErrItemNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "ITEM_NOT_FOUND",
					"message": "Item not found",
				},
			})
		case errors.Is(err, ErrNotBooster):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "NOT_A_BOOSTER",
					"message": err.Error(),
				},
			})
		case errors.Is(err, ErrInsufficientItem):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "INSUFFICIENT_ITEM",
					"message": err.Error(),
				},
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to use item",
				},
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("%s activated! It applies to your next battles.", boost.ItemName),
		"boost":   boost,
	})
}
