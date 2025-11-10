package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Middleware(jwtService *JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "UNAUTHORIZED",
					"message": "Missing authorization header",
				},
			})
		}

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

		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		
		return c.Next()
	}
}

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

func GetUserID(c *fiber.Ctx) (int, bool) {
	userID, ok := c.Locals("user_id").(int)
	return userID, ok
}

func GetUsername(c *fiber.Ctx) (string, bool) {
	username, ok := c.Locals("username").(string)
	return username, ok
}
