package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Middleware creates a Fiber middleware for JWT authentication
// Requirement: 2.3, 2.4 - Implement token validation middleware
func Middleware(jwtService *JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "UNAUTHORIZED",
					"message": "Missing authorization header",
				},
			})
		}
		
		// Check for Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "UNAUTHORIZED",
					"message": "Invalid authorization header format. Expected: Bearer <token>",
				},
			})
		}
		
		tokenString := parts[1]
		
		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			var code string
			var message string
			
			switch err {
			case ErrTokenExpired:
				code = "TOKEN_EXPIRED"
				message = "Token has expired"
			case ErrInvalidToken:
				code = "TOKEN_INVALID"
				message = "Token is invalid or malformed"
			default:
				code = "UNAUTHORIZED"
				message = err.Error()
			}
			
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    code,
					"message": message,
				},
			})
		}
		
		// Store user information in context for use in handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		
		return c.Next()
	}
}

// OptionalMiddleware creates a middleware that validates JWT if present but doesn't require it
func OptionalMiddleware(jwtService *JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}
		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}
		
		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err == nil {
			c.Locals("user_id", claims.UserID)
			c.Locals("username", claims.Username)
		}
		
		return c.Next()
	}
}

// GetUserID extracts the user ID from the request context
func GetUserID(c *fiber.Ctx) (int, bool) {
	userID, ok := c.Locals("user_id").(int)
	return userID, ok
}

// GetUsername extracts the username from the request context
func GetUsername(c *fiber.Ctx) (string, bool) {
	username, ok := c.Locals("username").(string)
	return username, ok
}
