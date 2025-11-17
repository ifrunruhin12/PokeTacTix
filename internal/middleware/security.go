package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prevent clickjacking attacks
		c.Set("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// Control referrer information
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// XSS Protection (legacy browsers)
		c.Set("X-XSS-Protection", "1; mode=block")

		// Content Security Policy
		c.Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data: https:; "+
				"font-src 'self' data:; "+
				"connect-src 'self'; "+
				"frame-ancestors 'none'")

		// Strict Transport Security (HTTPS only)
		if c.Protocol() == "https" {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		return c.Next()
	}
}

// HTTPSRedirect redirects HTTP requests to HTTPS in production
func HTTPSRedirect(env string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only redirect in production
		if env == "production" {
			// Check if request is not HTTPS
			if c.Protocol() != "https" {
				// Check X-Forwarded-Proto header (for proxies/load balancers)
				proto := c.Get("X-Forwarded-Proto")
				if proto != "https" {
					// Redirect to HTTPS
					return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusMovedPermanently)
				}
			}
		}
		return c.Next()
	}
}
