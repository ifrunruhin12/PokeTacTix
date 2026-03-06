package tokens

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// Default values for token system configuration
	DefaultDailyTokenLimit    = 10
	DefaultTokenPrice         = 100
	DefaultDailyPurchaseLimit = 10
	DefaultResetTime          = "00:00"
	DefaultTokenCost1v1       = 1
	DefaultTokenCost5v5       = 2
)

var (
	// ErrInsufficientTokens is returned when a user has no tokens available
	ErrInsufficientTokens = errors.New("insufficient tokens")

	// ErrInsufficientCoins is returned when a user doesn't have enough coins
	ErrInsufficientCoins = errors.New("insufficient coins")

	// ErrDailyLimitExceeded is returned when daily purchase limit is reached
	ErrDailyLimitExceeded = errors.New("daily purchase limit exceeded")

	// ErrInvalidQuantity is returned when purchase quantity is invalid
	ErrInvalidQuantity = errors.New("invalid quantity: must be between 1 and 10")

	// ErrUserNotFound is returned when user doesn't exist
	ErrUserNotFound = errors.New("user not found")
)

// Service handles token business logic
type Service struct {
	repo               *Repository
	db                 *pgxpool.Pool
	resetTime          string // Daily reset time in HH:MM format (UTC)
	dailyTokenLimit    int    // Number of free tokens per day
	tokenPrice         int    // Cost of one token in coins
	dailyPurchaseLimit int    // Maximum tokens that can be purchased per day
	tokenCost1v1       int    // Token cost for 1v1 battles
	tokenCost5v5       int    // Token cost for 5v5 battles
}

// NewService creates a new token service
func NewService(db *pgxpool.Pool) *Service {
	// Load reset time from environment variable
	resetTime := os.Getenv("TOKEN_RESET_TIME")
	if resetTime == "" {
		resetTime = DefaultResetTime
	}

	// Load daily token limit from environment variable
	dailyTokenLimit := DefaultDailyTokenLimit
	if envLimit := os.Getenv("DAILY_TOKEN_LIMIT"); envLimit != "" {
		if limit, err := strconv.Atoi(envLimit); err == nil && limit > 0 {
			dailyTokenLimit = limit
		}
	}

	// Load token price from environment variable
	tokenPrice := DefaultTokenPrice
	if envPrice := os.Getenv("TOKEN_PRICE"); envPrice != "" {
		if price, err := strconv.Atoi(envPrice); err == nil && price > 0 {
			tokenPrice = price
		}
	}

	// Load daily purchase limit from environment variable
	dailyPurchaseLimit := DefaultDailyPurchaseLimit
	if envPurchaseLimit := os.Getenv("DAILY_PURCHASE_LIMIT"); envPurchaseLimit != "" {
		if limit, err := strconv.Atoi(envPurchaseLimit); err == nil && limit > 0 {
			dailyPurchaseLimit = limit
		}
	}

	// Load token cost for 1v1 battles from environment variable
	tokenCost1v1 := DefaultTokenCost1v1
	if envCost := os.Getenv("TOKEN_COST_1V1"); envCost != "" {
		if cost, err := strconv.Atoi(envCost); err == nil && cost >= 0 {
			tokenCost1v1 = cost
		}
	}

	// Load token cost for 5v5 battles from environment variable
	tokenCost5v5 := DefaultTokenCost5v5
	if envCost := os.Getenv("TOKEN_COST_5V5"); envCost != "" {
		if cost, err := strconv.Atoi(envCost); err == nil && cost >= 0 {
			tokenCost5v5 = cost
		}
	}

	return &Service{
		repo:               NewRepository(db),
		db:                 db,
		resetTime:          resetTime,
		dailyTokenLimit:    dailyTokenLimit,
		tokenPrice:         tokenPrice,
		dailyPurchaseLimit: dailyPurchaseLimit,
		tokenCost1v1:       tokenCost1v1,
		tokenCost5v5:       tokenCost5v5,
	}
}

// GetResetTime returns the configured daily reset time
func (s *Service) GetResetTime() string {
	return s.resetTime + " UTC"
}

// GetDailyTokenLimit returns the configured daily token limit
func (s *Service) GetDailyTokenLimit() int {
	return s.dailyTokenLimit
}

// GetTokenPrice returns the configured token price
func (s *Service) GetTokenPrice() int {
	return s.tokenPrice
}

// GetDailyPurchaseLimit returns the configured daily purchase limit
func (s *Service) GetDailyPurchaseLimit() int {
	return s.dailyPurchaseLimit
}

// GetTokenCost1v1 returns the configured token cost for 1v1 battles
func (s *Service) GetTokenCost1v1() int {
	return s.tokenCost1v1
}

// GetTokenCost5v5 returns the configured token cost for 5v5 battles
func (s *Service) GetTokenCost5v5() int {
	return s.tokenCost5v5
}

// GetTokenCostForMode returns the token cost for a specific battle mode
func (s *Service) GetTokenCostForMode(mode string) int {
	if mode == "5v5" {
		return s.tokenCost5v5
	}
	return s.tokenCost1v1
}

// GetTokenBalance retrieves user's current token count
func (s *Service) GetTokenBalance(ctx context.Context, userID int) (int, error) {
	data, err := s.repo.GetUserTokenData(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get token balance: %w", err)
	}

	return data.GameTokens, nil
}

// ConsumeToken deducts tokens from user's balance based on battle mode
func (s *Service) ConsumeToken(ctx context.Context, userID int, mode string) error {
	// Get current token data
	data, err := s.repo.GetUserTokenData(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get token data: %w", err)
	}

	// Determine token cost based on mode
	tokenCost := s.GetTokenCostForMode(mode)

	// Validate user has enough tokens
	if data.GameTokens < tokenCost {
		return ErrInsufficientTokens
	}

	// Deduct tokens
	newBalance := data.GameTokens - tokenCost
	err = s.repo.UpdateTokenBalance(ctx, userID, newBalance)
	if err != nil {
		return fmt.Errorf("failed to consume token: %w", err)
	}

	return nil
}

// CheckAndResetTokens checks if new day has started and resets tokens if needed
func (s *Service) CheckAndResetTokens(ctx context.Context, userID int) error {
	// Get current token data
	data, err := s.repo.GetUserTokenData(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get token data: %w", err)
	}

	// Get current date in UTC
	now := time.Now().UTC()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	
	// Get last reset date (normalize to midnight UTC)
	lastResetDate := time.Date(
		data.LastTokenReset.Year(),
		data.LastTokenReset.Month(),
		data.LastTokenReset.Day(),
		0, 0, 0, 0, time.UTC,
	)

	// Check if current date is after last reset date
	if currentDate.After(lastResetDate) {
		// Only reset if user has less than the daily limit
		// This preserves purchased tokens above the limit
		if data.GameTokens < s.dailyTokenLimit {
			err = s.repo.ResetTokens(ctx, userID, s.dailyTokenLimit)
			if err != nil {
				return fmt.Errorf("failed to reset tokens: %w", err)
			}
		} else {
			// Just update the date without changing token count
			// User has more than daily limit from purchases
			err = s.updateResetDate(ctx, userID)
			if err != nil {
				return fmt.Errorf("failed to update reset date: %w", err)
			}
		}
		
		// Also reset the purchase counter for the new day
		err = s.resetPurchaseCounter(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to reset purchase counter: %w", err)
		}
	}

	return nil
}

// updateResetDate updates only the reset date without changing token count
func (s *Service) updateResetDate(ctx context.Context, userID int) error {
	_, err := s.db.Exec(ctx, `
		UPDATE users
		SET last_token_reset = CURRENT_DATE, updated_at = NOW()
		WHERE id = $1
	`, userID)

	if err != nil {
		return fmt.Errorf("failed to update reset date: %w", err)
	}

	return nil
}

// resetPurchaseCounter resets the daily purchase counter
func (s *Service) resetPurchaseCounter(ctx context.Context, userID int) error {
	_, err := s.db.Exec(ctx, `
		UPDATE users
		SET tokens_purchased_today = 0, updated_at = NOW()
		WHERE id = $1
	`, userID)

	if err != nil {
		return fmt.Errorf("failed to reset purchase counter: %w", err)
	}

	return nil
}

// PurchaseTokens adds tokens to user's balance with coin and limit validation
func (s *Service) PurchaseTokens(ctx context.Context, userID int, quantity int) error {
	// Validate quantity
	if quantity < 1 || quantity > s.dailyPurchaseLimit {
		return ErrInvalidQuantity
	}

	// Start a transaction to ensure atomicity
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get current token data and user coins
	var coins int
	var tokensPurchasedToday int
	err = tx.QueryRow(ctx, `
		SELECT coins, tokens_purchased_today
		FROM users
		WHERE id = $1
	`, userID).Scan(&coins, &tokensPurchasedToday)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user data: %w", err)
	}

	// Check daily purchase limit
	if tokensPurchasedToday+quantity > s.dailyPurchaseLimit {
		return ErrDailyLimitExceeded
	}

	// Calculate total cost
	totalCost := quantity * s.tokenPrice

	// Check if user has enough coins
	if coins < totalCost {
		return ErrInsufficientCoins
	}

	// Deduct coins and add tokens in a single transaction
	_, err = tx.Exec(ctx, `
		UPDATE users
		SET coins = coins - $1,
		    game_tokens = game_tokens + $2,
		    tokens_purchased_today = tokens_purchased_today + $2,
		    updated_at = NOW()
		WHERE id = $3
	`, totalCost, quantity, userID)

	if err != nil {
		return fmt.Errorf("failed to purchase tokens: %w", err)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTimeUntilReset calculates time remaining until next daily reset
func (s *Service) GetTimeUntilReset(ctx context.Context) (time.Duration, error) {
	now := time.Now().UTC()
	
	// Parse the reset time (HH:MM format)
	parts := strings.Split(s.resetTime, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid reset time format: %s", s.resetTime)
	}
	
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour in reset time: %w", err)
	}
	
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minute in reset time: %w", err)
	}
	
	// Create today's reset time
	todayReset := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, time.UTC)
	
	// If we've passed today's reset time, calculate time until tomorrow's reset
	var nextReset time.Time
	if now.After(todayReset) || now.Equal(todayReset) {
		// Next reset is tomorrow
		nextReset = todayReset.Add(24 * time.Hour)
	} else {
		// Next reset is today
		nextReset = todayReset
	}
	
	return nextReset.Sub(now), nil
}

// ResetAllPurchaseCounters resets purchase counters for all users who need a daily reset
// This should be called as part of a scheduled job or during the daily reset process
func (s *Service) ResetAllPurchaseCounters(ctx context.Context) (int64, error) {
	rowsAffected, err := s.repo.ResetAllPurchaseCounters(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to reset all purchase counters: %w", err)
	}
	return rowsAffected, nil
}
