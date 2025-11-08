package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Max        int
	Expiration time.Duration
}

// CreateRateLimiter creates a rate limiter middleware with custom configuration
// Requirement: 14.5 - Implement rate limiting for auth endpoints
func CreateRateLimiter(config RateLimitConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        config.Max,
		Expiration: config.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as the key for rate limiting
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests. Please try again later.",
					"details": fiber.Map{
						"max":        config.Max,
						"window":     config.Expiration.String(),
						"retry_after": config.Expiration.Seconds(),
					},
				},
			})
		},
		Storage: nil, // Uses in-memory storage by default
	})
}

// LoginRateLimiter creates a rate limiter for login endpoint
// Requirement: 14.5 - 5 requests/minute for login
func LoginRateLimiter() fiber.Handler {
	return CreateRateLimiter(RateLimitConfig{
		Max:        5,
		Expiration: 1 * time.Minute,
	})
}

// RegisterRateLimiter creates a rate limiter for register endpoint
// Requirement: 14.5 - 3 requests/hour for register
func RegisterRateLimiter() fiber.Handler {
	return CreateRateLimiter(RateLimitConfig{
		Max:        3,
		Expiration: 1 * time.Hour,
	})
}

// DefaultRateLimiter creates a default rate limiter for general auth endpoints
func DefaultRateLimiter() fiber.Handler {
	return CreateRateLimiter(RateLimitConfig{
		Max:        10,
		Expiration: 1 * time.Minute,
	})
}
