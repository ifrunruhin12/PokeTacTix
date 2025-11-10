package shop

import (
	"github.com/gofiber/fiber/v2"
)

// Handler handles shop HTTP requests
type Handler struct {
	service    *Service
	repository *Repository
}

// NewHandler creates a new shop handler
func NewHandler(service *Service, repository *Repository) *Handler {
	return &Handler{
		service:    service,
		repository: repository,
	}
}

// GetInventory handles GET /api/shop/inventory
func (h *Handler) GetInventory(c *fiber.Ctx) error {
	inventory := h.service.GetInventory()
	
	// Apply current prices with discounts
	for i := range inventory.Items {
		inventory.Items[i].Price = h.service.GetItemPrice(inventory.Items[i])
	}
	
	return c.JSON(inventory)
}

// Purchase handles POST /api/shop/purchase
func (h *Handler) Purchase(c *fiber.Ctx) error {
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

	// Parse request body
	var req PurchaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
	}

	// Validate Pokemon name
	if req.PokemonName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Pokemon name is required",
			},
		})
	}

	// Find item in shop inventory
	item, err := h.service.FindItem(req.PokemonName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "ITEM_NOT_FOUND",
				"message": "Pokemon not found in shop inventory",
			},
		})
	}

	// Get current price (with discount if active)
	price := h.service.GetItemPrice(*item)

	// Purchase the card
	card, err := h.repository.PurchaseCard(c.Context(), userID, req.PokemonName, price)
	if err != nil {
		// Check for insufficient coins error
		if err.Error()[:len("insufficient coins")] == "insufficient coins" {
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "INSUFFICIENT_COINS",
					"message": err.Error(),
				},
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "PURCHASE_FAILED",
				"message": "Failed to complete purchase",
				"details": err.Error(),
			},
		})
	}

	// Get remaining coins
	remainingCoins, err := h.repository.GetUserCoins(c.Context(), userID)
	if err != nil {
		remainingCoins = 0 // Fallback
	}

	return c.Status(fiber.StatusOK).JSON(PurchaseResponse{
		Card:           card,
		RemainingCoins: remainingCoins,
	})
}
