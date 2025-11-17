package middleware

import (
	"html"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// SanitizeString removes potentially dangerous characters and HTML tags
func SanitizeString(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Escape HTML to prevent XSS
	input = html.EscapeString(input)

	return input
}

// ValidateUsername validates username format
func ValidateUsername(username string) *ValidationError {
	username = strings.TrimSpace(username)

	if username == "" {
		return &ValidationError{
			Field:   "username",
			Message: "Username is required",
		}
	}

	if len(username) < 3 {
		return &ValidationError{
			Field:   "username",
			Message: "Username must be at least 3 characters long",
		}
	}

	if len(username) > 50 {
		return &ValidationError{
			Field:   "username",
			Message: "Username must not exceed 50 characters",
		}
	}

	// Only allow alphanumeric characters and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	if !matched {
		return &ValidationError{
			Field:   "username",
			Message: "Username can only contain letters, numbers, and underscores",
		}
	}

	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) *ValidationError {
	email = strings.TrimSpace(email)

	if email == "" {
		return &ValidationError{
			Field:   "email",
			Message: "Email is required",
		}
	}

	// Basic email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		}
	}

	if len(email) > 255 {
		return &ValidationError{
			Field:   "email",
			Message: "Email must not exceed 255 characters",
		}
	}

	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) *ValidationError {
	if password == "" {
		return &ValidationError{
			Field:   "password",
			Message: "Password is required",
		}
	}

	if len(password) < 8 {
		return &ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
		}
	}

	if len(password) > 100 {
		return &ValidationError{
			Field:   "password",
			Message: "Password must not exceed 100 characters",
		}
	}

	// Check for uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return &ValidationError{
			Field:   "password",
			Message: "Password must contain at least one uppercase letter",
		}
	}

	// Check for lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return &ValidationError{
			Field:   "password",
			Message: "Password must contain at least one lowercase letter",
		}
	}

	// Check for number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return &ValidationError{
			Field:   "password",
			Message: "Password must contain at least one number",
		}
	}

	// Check for special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		return &ValidationError{
			Field:   "password",
			Message: "Password must contain at least one special character",
		}
	}

	return nil
}

// ValidatePokemonName validates Pokemon name input
func ValidatePokemonName(name string) *ValidationError {
	name = strings.TrimSpace(name)

	if name == "" {
		return &ValidationError{
			Field:   "pokemon_name",
			Message: "Pokemon name is required",
		}
	}

	// Pokemon names should only contain letters, hyphens, and spaces
	matched, _ := regexp.MatchString(`^[a-zA-Z\-\s]+$`, name)
	if !matched {
		return &ValidationError{
			Field:   "pokemon_name",
			Message: "Invalid Pokemon name format",
		}
	}

	if len(name) > 100 {
		return &ValidationError{
			Field:   "pokemon_name",
			Message: "Pokemon name is too long",
		}
	}

	return nil
}

// ValidateBattleMode validates battle mode input
func ValidateBattleMode(mode string) *ValidationError {
	mode = strings.TrimSpace(mode)

	if mode == "" {
		return &ValidationError{
			Field:   "mode",
			Message: "Battle mode is required",
		}
	}

	if mode != "1v1" && mode != "5v5" {
		return &ValidationError{
			Field:   "mode",
			Message: "Battle mode must be either '1v1' or '5v5'",
		}
	}

	return nil
}

// ValidateBattleMove validates battle move input
func ValidateBattleMove(move string) *ValidationError {
	move = strings.TrimSpace(move)

	if move == "" {
		return &ValidationError{
			Field:   "move",
			Message: "Move is required",
		}
	}

	validMoves := map[string]bool{
		"attack":    true,
		"defend":    true,
		"pass":      true,
		"sacrifice": true,
		"surrender": true,
	}

	if !validMoves[move] {
		return &ValidationError{
			Field:   "move",
			Message: "Invalid move. Must be one of: attack, defend, pass, sacrifice, surrender",
		}
	}

	return nil
}

// ValidateCardIDs validates an array of card IDs
func ValidateCardIDs(cardIDs []int, expectedCount int) *ValidationError {
	if len(cardIDs) != expectedCount {
		return &ValidationError{
			Field:   "card_ids",
			Message: "Invalid number of cards",
		}
	}

	// Check for duplicates
	seen := make(map[int]bool)
	for _, id := range cardIDs {
		if id <= 0 {
			return &ValidationError{
				Field:   "card_ids",
				Message: "Invalid card ID",
			}
		}
		if seen[id] {
			return &ValidationError{
				Field:   "card_ids",
				Message: "Duplicate card IDs are not allowed",
			}
		}
		seen[id] = true
	}

	return nil
}

// ValidatePositiveInteger validates that a value is a positive integer
func ValidatePositiveInteger(value int, fieldName string) *ValidationError {
	if value <= 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be a positive integer",
		}
	}
	return nil
}

// ValidateNonNegativeInteger validates that a value is non-negative
func ValidateNonNegativeInteger(value int, fieldName string) *ValidationError {
	if value < 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be non-negative",
		}
	}
	return nil
}

// InputSanitizer middleware sanitizes all string inputs in request body
// Note: This is a basic implementation. For production use, consider using
// a dedicated sanitization library or implementing sanitization at the handler level
func InputSanitizer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Sanitization is now handled at the handler level for better control
		// This middleware can be extended in the future if needed
		return c.Next()
	}
}

// sanitizeMap recursively sanitizes all string values in a map
func sanitizeMap(data map[string]any) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			data[key] = SanitizeString(v)
		case map[string]any:
			sanitizeMap(v)
		case []any:
			sanitizeSlice(v)
		}
	}
}

// sanitizeSlice recursively sanitizes all string values in a slice
func sanitizeSlice(data []any) {
	for i, value := range data {
		switch v := value.(type) {
		case string:
			data[i] = SanitizeString(v)
		case map[string]any:
			sanitizeMap(v)
		case []any:
			sanitizeSlice(v)
		}
	}
}

// ValidationErrorResponse returns a standardized validation error response
func ValidationErrorResponse(c *fiber.Ctx, errors []*ValidationError) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    "VALIDATION_ERROR",
			"message": "Input validation failed",
			"details": errors,
		},
	})
}
