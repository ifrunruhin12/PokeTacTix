package stats

import (
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for statistics and profile
type Handler struct {
	service *Service
}

// NewHandler creates a new stats handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetPlayerStats handles GET /api/profile/stats
// Requirements: 8.1, 8.5
func (h *Handler) GetPlayerStats(c *fiber.Ctx) error {
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

	// Get player stats
	stats, err := h.service.GetPlayerStats(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve player statistics",
				"details": err.Error(),
			},
		})
	}

	return c.JSON(stats)
}

// GetBattleHistory handles GET /api/profile/history
// Requirements: 8.3
func (h *Handler) GetBattleHistory(c *fiber.Ctx) error {
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

	// Get limit from query params (default 20)
	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100 // Cap at 100 to prevent excessive queries
	}

	// Get battle history
	history, err := h.service.GetBattleHistory(c.Context(), userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve battle history",
				"details": err.Error(),
			},
		})
	}

	// Return empty array if no history
	if history == nil {
		return c.JSON(fiber.Map{
			"history": []interface{}{},
			"count":   0,
		})
	}

	return c.JSON(fiber.Map{
		"history": history,
		"count":   len(history),
	})
}

// GetAchievements handles GET /api/profile/achievements
// Requirements: 8.4
func (h *Handler) GetAchievements(c *fiber.Ctx) error {
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

	// Get achievements
	achievements, err := h.service.GetAchievements(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve achievements",
				"details": err.Error(),
			},
		})
	}

	// Return empty array if no achievements
	if achievements == nil {
		return c.JSON(fiber.Map{
			"achievements": []interface{}{},
			"total":        0,
			"unlocked":     0,
			"locked":       0,
		})
	}

	// Count unlocked achievements
	unlockedCount := 0
	for _, ach := range achievements {
		if ach.Unlocked {
			unlockedCount++
		}
	}

	return c.JSON(fiber.Map{
		"achievements": achievements,
		"total":        len(achievements),
		"unlocked":     unlockedCount,
		"locked":       len(achievements) - unlockedCount,
	})
}

// CheckAchievements handles POST /api/profile/achievements/check
// This endpoint checks and unlocks any newly earned achievements
// Requirements: 8.4
func (h *Handler) CheckAchievements(c *fiber.Ctx) error {
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

	// Check and unlock achievements
	newlyUnlocked, err := h.service.CheckAndUnlockAchievements(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to check achievements",
				"details": err.Error(),
			},
		})
	}

	// Return empty array if no new achievements
	if newlyUnlocked == nil {
		return c.JSON(fiber.Map{
			"newly_unlocked": []interface{}{},
			"count":          0,
		})
	}

	return c.JSON(fiber.Map{
		"newly_unlocked": newlyUnlocked,
		"count":          len(newlyUnlocked),
	})
}
