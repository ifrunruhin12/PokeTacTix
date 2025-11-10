package auth

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptCost = 12

	MinPasswordLength = 8
	MaxPasswordLength = 100

	MinUsernameLength = 3
	MaxUsernameLength = 50
)

var (
	ErrWeakPassword = errors.New("password must be at least 8 characters with uppercase, lowercase, number, and special character")

	ErrInvalidUsername = errors.New("username must be 3-50 alphanumeric characters")

	ErrInvalidEmail = errors.New("invalid email format")
	
	emailRegex     = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	usernameRegex  = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	uppercaseRegex = regexp.MustCompile(`[A-Z]`)
	lowercaseRegex = regexp.MustCompile(`[a-z]`)
	digitRegex     = regexp.MustCompile(`[0-9]`)
	specialRegex   = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
)

// Service handles authentication business logic
type Service struct{}

// NewService creates a new auth service
func NewService() *Service {
	return &Service{}
}

// HashPassword hashes a password using bcrypt
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

// ComparePassword compares a hashed password with a plaintext password
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

// ValidatePassword validates password strength
func (s *Service) ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ErrWeakPassword
	}
	
	if len(password) > MaxPasswordLength {
		return errors.New("password is too long (max 100 characters)")
	}

	if !uppercaseRegex.MatchString(password) {
		return ErrWeakPassword
	}

	if !lowercaseRegex.MatchString(password) {
		return ErrWeakPassword
	}

	if !digitRegex.MatchString(password) {
		return ErrWeakPassword
	}

	if !specialRegex.MatchString(password) {
		return ErrWeakPassword
	}
	
	return nil
}

// ValidateUsername validates username format
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
