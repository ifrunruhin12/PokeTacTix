package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid or malformed token")
	ErrTokenExpired = errors.New("token has expired")
	ErrMissingToken = errors.New("missing authentication token")
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey       []byte
	tokenExpiration time.Duration
}

func NewJWTService(secret string, expiration time.Duration) (*JWTService, error) {
	if secret == "" {
		return nil, errors.New("JWT secret is not set")
	}

	if len(secret) < 32 {
		return nil, errors.New("JWT secret must be at least 32 characters (256 bits)")
	}

	return &JWTService{
		secretKey:       []byte(secret),
		tokenExpiration: expiration,
	}, nil
}

// GenerateToken generates a new JWT token
func (s *JWTService) GenerateToken(userID int, username string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(s.tokenExpiration)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, ErrMissingToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		if !errors.Is(err, ErrTokenExpired) {
			return "", err
		}
	}

	return s.GenerateToken(claims.UserID, claims.Username)
}

func (s *JWTService) ExtractUserID(tokenString string) (int, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
