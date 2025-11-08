package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// TokenExpiration is the duration for which tokens are valid (24 hours)
	// Requirement: 2.4 - Set token expiration to 24 hours
	TokenExpiration = 24 * time.Hour
)

var (
	// ErrInvalidToken indicates the token is invalid or malformed
	ErrInvalidToken = errors.New("invalid or malformed token")
	
	// ErrTokenExpired indicates the token has expired
	ErrTokenExpired = errors.New("token has expired")
	
	// ErrMissingToken indicates no token was provided
	ErrMissingToken = errors.New("missing authentication token")
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secretKey []byte
}

// NewJWTService creates a new JWT service
// Requirement: 14.3 - Store JWT secret in environment variable
func NewJWTService() (*JWTService, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET environment variable is not set")
	}
	
	// Requirement: 14.3 - Use secure random secrets with minimum 256-bit entropy
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must be at least 32 characters (256 bits)")
	}
	
	return &JWTService{
		secretKey: []byte(secret),
	}, nil
}

// GenerateToken generates a new JWT token for a user
// Requirement: 2.3 - Generate secure JWT token with 24-hour expiration using HS256
func (s *JWTService) GenerateToken(userID int, username string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(TokenExpiration)
	
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	
	// Create token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign the token with the secret key
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	
	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
// Requirement: 2.4 - Implement token validation
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, ErrMissingToken
	}
	
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
	
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	
	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	
	return claims, nil
}

// RefreshToken generates a new token for an existing valid token
// This can be used to extend user sessions
func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		// Allow refresh even if token is expired (but not if invalid)
		if !errors.Is(err, ErrTokenExpired) {
			return "", err
		}
	}
	
	// Generate new token with same user info
	return s.GenerateToken(claims.UserID, claims.Username)
}

// ExtractUserID extracts the user ID from a token without full validation
// Useful for logging or non-critical operations
func (s *JWTService) ExtractUserID(tokenString string) (int, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
