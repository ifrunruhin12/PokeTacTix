package shop

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RegisterRoutes registers shop routes
func RegisterRoutes(app *fiber.App, handler *Handler, authMiddleware fiber.Handler) {
	shop := app.Group("/api/shop")

	// Apply authentication middleware to all shop routes
	shop.Use(authMiddleware)

	// GET /api/shop/inventory - Get current shop inventory
	shop.Get("/inventory", handler.GetInventory)

	// POST /api/shop/purchase - Purchase a Pokemon card
	// Rate limit: 10 purchases per minute
	shop.Post("/purchase", createPurchaseRateLimiter(), handler.Purchase)
}

// createPurchaseRateLimiter creates a rate limiter for purchase endpoint
func createPurchaseRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Rate limit per user ID if available, otherwise per IP
			if userID, ok := c.Locals("user_id").(int); ok {
				return fmt.Sprintf("user:%d", userID)
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many purchase requests. Please try again later.",
					"details": fiber.Map{
						"max":         10,
						"window":      "1 minute",
						"retry_after": 60,
					},
				},
			})
		},
		Storage: nil,
	})
}
