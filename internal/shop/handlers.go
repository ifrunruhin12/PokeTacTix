package shop

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"pokemon-cli/internal/tokens"

	"github.com/gofiber/fiber/v2"
)

// Handler handles shop HTTP requests
type Handler struct {
	service      *Service
	repository   *Repository
	tokenService *tokens.Service
}

// NewHandler creates a new shop handler
func NewHandler(service *Service, repository *Repository, tokenService *tokens.Service) *Handler {
	return &Handler{
		service:      service,
		repository:   repository,
		tokenService: tokenService,
	}
}

// isValidPokemonName validates Pokemon name format
func isValidPokemonName(name string) bool {
	// Pokemon names should only contain letters, hyphens, spaces, and apostrophes
	matched, _ := regexp.MatchString(`^[a-zA-Z\-\s']+$`, name)
	return matched
}

// GetInventory handles GET /api/shop/inventory
func (h *Handler) GetInventory(c *fiber.Ctx) error {
	inventory := h.service.GetInventory()

	// Apply current prices with discounts
	for i := range inventory.Items {
		inventory.Items[i].Price = h.service.GetItemPrice(inventory.Items[i])
	}

	// Add game token item if user is authenticated
	userID, ok := c.Locals("user_id").(int)
	if ok {
		// Get tokens purchased today for this user
		tokensPurchasedToday, err := h.repository.GetTokensPurchasedToday(c.Context(), userID)
		if err != nil {
			// If we can't get the data, default to 0 (allow full purchase)
			tokensPurchasedToday = 0
		}

		// Get configuration from token service
		dailyPurchaseLimit := h.tokenService.GetDailyPurchaseLimit()
		tokenPrice := h.tokenService.GetTokenPrice()
		resetTime := h.tokenService.GetResetTime()

		// Calculate available quantity (dailyPurchaseLimit - tokens_purchased_today)
		availableQuantity := dailyPurchaseLimit - tokensPurchasedToday
		if availableQuantity < 0 {
			availableQuantity = 0
		}

		// Add game token item to inventory
		inventory.GameTokens = &GameTokenItem{
			ItemType:          "game_token",
			Price:             tokenPrice,
			AvailableQuantity: availableQuantity,
			MaxDailyPurchase:  dailyPurchaseLimit,
			Description:       fmt.Sprintf("Game tokens allow you to participate in battles. Each battle costs 1 token. Tokens reset daily at %s.", resetTime),
		}
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

	// Sanitize and validate Pokemon name
	req.PokemonName = strings.TrimSpace(strings.ToLower(req.PokemonName))

	if req.PokemonName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Pokemon name is required",
			},
		})
	}

	// Validate Pokemon name format (only letters, hyphens, and spaces)
	if !isValidPokemonName(req.PokemonName) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid Pokemon name format",
			},
		})
	}

	if len(req.PokemonName) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Pokemon name is too long",
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

// PurchaseTokens handles POST /api/shop/tokens/purchase
func (h *Handler) PurchaseTokens(c *fiber.Ctx) error {
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
	var req TokenPurchaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
			},
		})
	}

	// Get configuration from token service
	dailyPurchaseLimit := h.tokenService.GetDailyPurchaseLimit()
	tokenPrice := h.tokenService.GetTokenPrice()

	// Validate quantity range (1 to dailyPurchaseLimit)
	if req.Quantity < 1 || req.Quantity > dailyPurchaseLimit {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_QUANTITY",
				"message": fmt.Sprintf("Invalid quantity. Must be between 1 and %d.", dailyPurchaseLimit),
			},
		})
	}

	// Get user's current token and coin data before purchase
	tokenData, err := h.tokenService.GetTokenBalance(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "TOKEN_SYSTEM_ERROR",
				"message": "Failed to retrieve token balance",
			},
		})
	}

	// Get user's current coins
	currentCoins, err := h.repository.GetUserCoins(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "SYSTEM_ERROR",
				"message": "Failed to retrieve user coins",
			},
		})
	}

	// Call token service PurchaseTokens method
	err = h.tokenService.PurchaseTokens(c.Context(), userID, req.Quantity)
	if err != nil {
		// Handle specific error types
		if errors.Is(err, tokens.ErrInsufficientCoins) {
			totalCost := req.Quantity * tokenPrice
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "INSUFFICIENT_COINS",
					"message": "You don't have enough coins to purchase tokens",
					"details": fiber.Map{
						"required": totalCost,
						"current":  currentCoins,
					},
				},
			})
		}

		if errors.Is(err, tokens.ErrDailyLimitExceeded) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "DAILY_LIMIT_EXCEEDED",
					"message": fmt.Sprintf("Daily token purchase limit reached. Maximum %d tokens per day.", dailyPurchaseLimit),
				},
			})
		}

		if errors.Is(err, tokens.ErrInvalidQuantity) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "INVALID_QUANTITY",
					"message": fmt.Sprintf("Invalid quantity. Must be between 1 and %d.", dailyPurchaseLimit),
				},
			})
		}

		// Generic error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "PURCHASE_FAILED",
				"message": "Failed to complete token purchase",
				"details": err.Error(),
			},
		})
	}

	// Get updated balances after purchase
	newTokenBalance, err := h.tokenService.GetTokenBalance(c.Context(), userID)
	if err != nil {
		newTokenBalance = tokenData + req.Quantity // Fallback calculation
	}

	newCoinBalance, err := h.repository.GetUserCoins(c.Context(), userID)
	if err != nil {
		newCoinBalance = currentCoins - (req.Quantity * tokenPrice) // Fallback calculation
	}

	// Get tokens purchased today (need to query this from the database)
	tokensPurchasedToday, err := h.repository.GetTokensPurchasedToday(c.Context(), userID)
	if err != nil {
		tokensPurchasedToday = req.Quantity // Fallback - at least the current purchase
	}

	// Calculate coins spent
	coinsSpent := req.Quantity * tokenPrice

	// Return success response with updated balances
	return c.Status(fiber.StatusOK).JSON(TokenPurchaseResponse{
		Success:              true,
		TokensAdded:          req.Quantity,
		NewTokenBalance:      newTokenBalance,
		CoinsSpent:           coinsSpent,
		RemainingCoins:       newCoinBalance,
		TokensPurchasedToday: tokensPurchasedToday,
		DailyLimitRemaining:  dailyPurchaseLimit - tokensPurchasedToday,
	})
}
