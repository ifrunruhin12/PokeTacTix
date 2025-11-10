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

// CreateRateLimiter creates a rate limiter with the given configuration
func CreateRateLimiter(config RateLimitConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        config.Max,
		Expiration: config.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests. Please try again later.",
					"details": fiber.Map{
						"max":         config.Max,
						"window":      config.Expiration.String(),
						"retry_after": config.Expiration.Seconds(),
					},
				},
			})
		},
		Storage: nil,
	})
}

// LoginRateLimiter creates a rate limiter for login endpoint
func LoginRateLimiter() fiber.Handler {
	return CreateRateLimiter(RateLimitConfig{
		Max:        5,
		Expiration: 1 * time.Minute,
	})
}

// RegisterRateLimiter creates a rate limiter for register endpoint
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
