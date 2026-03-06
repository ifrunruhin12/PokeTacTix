package tokens

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Handler handles token HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new token handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// TokenBalanceResponse represents the response for token balance endpoint
type TokenBalanceResponse struct {
	GameTokens           int    `json:"game_tokens"`
	TimeUntilReset       string `json:"time_until_reset"`
	ResetTime            string `json:"reset_time"`
	TokensPurchasedToday int    `json:"tokens_purchased_today"`
	DailyPurchaseLimit   int    `json:"daily_purchase_limit"`
	DailyTokenLimit      int    `json:"daily_token_limit"`
	TokenCost1v1         int    `json:"token_cost_1v1"`
	TokenCost5v5         int    `json:"token_cost_5v5"`
}

// GetTokenBalanceHandler handles GET /api/tokens/balance
func (h *Handler) GetTokenBalanceHandler(c *fiber.Ctx) error {
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

	// Check and reset tokens if needed
	err := h.service.CheckAndResetTokens(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "TOKEN_SYSTEM_ERROR",
				"message": "Failed to check token reset",
				"details": err.Error(),
			},
		})
	}

	// Get token balance
	tokenBalance, err := h.service.GetTokenBalance(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "TOKEN_SYSTEM_ERROR",
				"message": "Failed to retrieve token balance",
				"details": err.Error(),
			},
		})
	}

	// Get time until reset (no longer needs userID)
	timeUntilReset, err := h.service.GetTimeUntilReset(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "TOKEN_SYSTEM_ERROR",
				"message": "Failed to calculate time until reset",
				"details": err.Error(),
			},
		})
	}

	// Format time until reset as human-readable string
	hours := int(timeUntilReset.Hours())
	minutes := int(timeUntilReset.Minutes()) % 60
	timeUntilResetStr := fmt.Sprintf("%dh %dm", hours, minutes)

	// Get reset time from service
	resetTime := h.service.GetResetTime()

	// Get tokens purchased today from repository
	tokenData, err := h.service.repo.GetUserTokenData(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "TOKEN_SYSTEM_ERROR",
				"message": "Failed to retrieve token purchase data",
				"details": err.Error(),
			},
		})
	}

	// Return token balance response
	return c.Status(fiber.StatusOK).JSON(TokenBalanceResponse{
		GameTokens:           tokenBalance,
		TimeUntilReset:       timeUntilResetStr,
		ResetTime:            resetTime,
		TokensPurchasedToday: tokenData.TokensPurchasedToday,
		DailyPurchaseLimit:   h.service.GetDailyPurchaseLimit(),
		DailyTokenLimit:      h.service.GetDailyTokenLimit(),
		TokenCost1v1:         h.service.GetTokenCost1v1(),
		TokenCost5v5:         h.service.GetTokenCost5v5(),
	})
}
