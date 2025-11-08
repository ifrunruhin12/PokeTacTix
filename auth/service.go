package auth

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost is the cost factor for bcrypt hashing (requirement: 14.1)
	BcryptCost = 12
	
	// Password requirements (requirement: 2.6)
	MinPasswordLength = 8
	MaxPasswordLength = 100
	
	// Username requirements
	MinUsernameLength = 3
	MaxUsernameLength = 50
)

var (
	// ErrWeakPassword indicates password doesn't meet requirements
	ErrWeakPassword = errors.New("password must be at least 8 characters with uppercase, lowercase, number, and special character")
	
	// ErrInvalidUsername indicates username doesn't meet requirements
	ErrInvalidUsername = errors.New("username must be 3-50 alphanumeric characters")
	
	// ErrInvalidEmail indicates email format is invalid
	ErrInvalidEmail = errors.New("invalid email format")
	
	// Regular expressions for validation
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	uppercaseRegex = regexp.MustCompile(`[A-Z]`)
	lowercaseRegex = regexp.MustCompile(`[a-z]`)
	digitRegex     = regexp.MustCompile(`[0-9]`)
	specialRegex   = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
)

// Service handles authentication operations
type Service struct{}

// NewService creates a new auth service
func NewService() *Service {
	return &Service{}
}

// HashPassword hashes a password using bcrypt with cost factor 12
// Requirement: 14.1 - Use bcrypt hashing with minimum cost factor of 12
func (s *Service) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	
	return string(hash), nil
}

// ComparePassword compares a plaintext password with a hashed password
func (s *Service) ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("invalid password")
		}
		return fmt.Errorf("failed to compare password: %w", err)
	}
	return nil
}

// ValidatePassword validates password meets security requirements
// Requirement: 2.6 - Password must have min 8 chars, uppercase, lowercase, number, special char
func (s *Service) ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ErrWeakPassword
	}
	
	if len(password) > MaxPasswordLength {
		return errors.New("password is too long (max 100 characters)")
	}
	
	// Check for uppercase letter
	if !uppercaseRegex.MatchString(password) {
		return ErrWeakPassword
	}
	
	// Check for lowercase letter
	if !lowercaseRegex.MatchString(password) {
		return ErrWeakPassword
	}
	
	// Check for digit
	if !digitRegex.MatchString(password) {
		return ErrWeakPassword
	}
	
	// Check for special character
	if !specialRegex.MatchString(password) {
		return ErrWeakPassword
	}
	
	return nil
}

// ValidateUsername validates username meets requirements
func (s *Service) ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	
	if len(username) < MinUsernameLength {
		return fmt.Errorf("username must be at least %d characters", MinUsernameLength)
	}
	
	if len(username) > MaxUsernameLength {
		return fmt.Errorf("username must be at most %d characters", MaxUsernameLength)
	}
	
	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}
	
	return nil
}

// ValidateEmail validates email format
func (s *Service) ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	
	if email == "" {
		return ErrInvalidEmail
	}
	
	if len(email) > 255 {
		return errors.New("email is too long (max 255 characters)")
	}
	
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	
	return nil
}
