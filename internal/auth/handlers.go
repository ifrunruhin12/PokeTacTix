package auth

import (
	"context"
	"errors"
	"pokemon-cli/internal/database"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

// CardService interface for generating starter decks
type CardService interface {
	GenerateStarterDeck(ctx context.Context, userID int) ([]database.PlayerCard, error)
}

// Handler handles authentication HTTP requests
type Handler struct {
	authService *Service
	jwtService  *JWTService
	repository  *Repository
	cardService CardService
}

// NewHandler creates a new auth handler
func NewHandler(authService *Service, jwtService *JWTService, repository *Repository, cardService CardService) *Handler {
	return &Handler{
		authService: authService,
		jwtService:  jwtService,
		repository:  repository,
		cardService: cardService,
	}
}

// Register handles user registration
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
				"details": err.Error(),
			},
		})
	}

	// Sanitize inputs
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	
	// Additional XSS protection - ensure no HTML tags
	if strings.Contains(req.Username, "<") || strings.Contains(req.Username, ">") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_USERNAME",
				"message": "Username contains invalid characters",
			},
		})
	}
	
	// Validate username
	if err := h.authService.ValidateUsername(req.Username); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_USERNAME",
				"message": err.Error(),
			},
		})
	}

	if err := h.authService.ValidateEmail(req.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_EMAIL",
				"message": err.Error(),
			},
		})
	}
	
	if err := h.authService.ValidatePassword(req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "WEAK_PASSWORD",
				"message": err.Error(),
			},
		})
	}
	
	ctx := context.Background()
	
	usernameExists, err := h.repository.UsernameExists(ctx, req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to check username availability",
			},
		})
	}
	if usernameExists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "USER_EXISTS",
				"message": "Username already exists",
			},
		})
	}
	
	emailExists, err := h.repository.EmailExists(ctx, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to check email availability",
			},
		})
	}
	if emailExists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "USER_EXISTS",
				"message": "Email already exists",
			},
		})
	}

	passwordHash, err := h.authService.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to process password",
			},
		})
	}
	
	user, err := h.repository.Create(ctx, req.Username, req.Email, passwordHash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create user",
			},
		})
	}

	// Generate starter deck for new user (Requirement: 3.1, 3.2, 3.3, 3.4)
	if h.cardService != nil {
		_, err = h.cardService.GenerateStarterDeck(ctx, user.ID)
		if err != nil {
			// Log error but don't fail registration
			// User can still use the account, just without starter cards
			_ = err
		}
	}

	token, err := h.jwtService.GenerateToken(user.ID, user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to generate authentication token",
			},
		})
	}
	
	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login handles user login
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
				"details": err.Error(),
			},
		})
	}

	req.Username = strings.TrimSpace(req.Username)
	
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Username and password are required",
			},
		})
	}
	
	ctx := context.Background()
	
	// Get user by username
	user, err := h.repository.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "INVALID_CREDENTIALS",
					"message": "Invalid username or password",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to authenticate",
			},
		})
	}

	if err := h.authService.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid username or password",
			},
		})
	}

	token, err := h.jwtService.GenerateToken(user.ID, user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to generate authentication token",
			},
		})
	}
	
	return c.JSON(AuthResponse{
		Token: token,
		User:  user,
	})
}

// GetCurrentUser handles getting current user info
func (h *Handler) GetCurrentUser(c *fiber.Ctx) error {
	// Extract user ID from context (set by auth middleware)
	userID, ok := GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
	}
	
	ctx := context.Background()

	user, err := h.repository.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "USER_NOT_FOUND",
					"message": "User not found",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve user information",
			},
		})
	}
	
	return c.JSON(fiber.Map{
		"user": user,
	})
}

// Logout handles user logout
func (h *Handler) Logout(c *fiber.Ctx) error {
	// In a stateless JWT system, logout is primarily handled client-side
	// by removing the token. This endpoint exists for consistency and
	// could be extended to maintain a token blacklist if needed.
	
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
