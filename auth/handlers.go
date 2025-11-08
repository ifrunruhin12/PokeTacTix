package auth

import (
	"context"
	"errors"
	"pokemon-cli/database"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

// Handler handles authentication HTTP requests
type Handler struct {
	authService *Service
	jwtService  *JWTService
	userRepo    *database.UserRepository
}

// NewHandler creates a new authentication handler
func NewHandler(authService *Service, jwtService *JWTService, userRepo *database.UserRepository) *Handler {
	return &Handler{
		authService: authService,
		jwtService:  jwtService,
		userRepo:    userRepo,
	}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string         `json:"token"`
	User  *database.User `json:"user"`
}

// Register handles user registration
// Requirement: 2.1 - POST /api/auth/register endpoint with user creation
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
	
	// Trim whitespace
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	
	// Validate username
	if err := h.authService.ValidateUsername(req.Username); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_USERNAME",
				"message": err.Error(),
			},
		})
	}
	
	// Validate email
	if err := h.authService.ValidateEmail(req.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_EMAIL",
				"message": err.Error(),
			},
		})
	}
	
	// Validate password
	// Requirement: 2.6, 14.1 - Password validation with security requirements
	if err := h.authService.ValidatePassword(req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "WEAK_PASSWORD",
				"message": err.Error(),
			},
		})
	}
	
	ctx := context.Background()
	
	// Check if username already exists
	// Requirement: 2.2 - Reject registration with existing username
	usernameExists, err := h.userRepo.UsernameExists(ctx, req.Username)
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
	
	// Check if email already exists
	// Requirement: 2.2 - Reject registration with existing email
	emailExists, err := h.userRepo.EmailExists(ctx, req.Email)
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
	
	// Hash password
	// Requirement: 14.1 - Use bcrypt hashing with cost factor 12
	passwordHash, err := h.authService.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to process password",
			},
		})
	}
	
	// Create user
	user, err := h.userRepo.Create(ctx, req.Username, req.Email, passwordHash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create user",
			},
		})
	}
	
	// Generate JWT token
	// Requirement: 2.3 - Generate JWT token on successful registration
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
// Requirement: 2.2 - POST /api/auth/login endpoint with credential validation
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
	
	// Trim whitespace
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
	user, err := h.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "not found") {
			// Don't reveal whether username exists
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
	
	// Compare password
	if err := h.authService.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid username or password",
			},
		})
	}
	
	// Generate JWT token
	// Requirement: 2.3 - Generate JWT token on successful login
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

// GetCurrentUser returns the currently authenticated user's information
// Requirement: 2.3 - GET /api/auth/me endpoint for current user info
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
	
	// Get user from database
	user, err := h.userRepo.GetByID(ctx, userID)
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

// Logout handles user logout (optional - mainly for client-side token removal)
// Requirement: 2.5 - Invalidate session on logout
func (h *Handler) Logout(c *fiber.Ctx) error {
	// In a stateless JWT system, logout is primarily handled client-side
	// by removing the token. This endpoint exists for consistency and
	// could be extended to maintain a token blacklist if needed.
	
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
