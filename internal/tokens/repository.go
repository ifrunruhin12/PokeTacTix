package tokens

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserTokenData represents token-related fields for a user
type UserTokenData struct {
	UserID               int
	GameTokens           int
	LastTokenReset       time.Time
	TokensPurchasedToday int
}

// Repository handles database operations for tokens
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new token repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetUserTokenData retrieves token-related fields for a user
func (r *Repository) GetUserTokenData(ctx context.Context, userID int) (*UserTokenData, error) {
	data := &UserTokenData{UserID: userID}

	err := r.db.QueryRow(ctx, `
		SELECT game_tokens, last_token_reset, tokens_purchased_today
		FROM users
		WHERE id = $1
	`, userID).Scan(&data.GameTokens, &data.LastTokenReset, &data.TokensPurchasedToday)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user token data: %w", err)
	}

	return data, nil
}

// UpdateTokenBalance updates user's token count
func (r *Repository) UpdateTokenBalance(ctx context.Context, userID int, newBalance int) error {
	result, err := r.db.Exec(ctx, `
		UPDATE users
		SET game_tokens = $1, updated_at = NOW()
		WHERE id = $2
	`, newBalance, userID)

	if err != nil {
		return fmt.Errorf("failed to update token balance: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ResetTokens sets tokens to the daily limit and updates reset date to current date
func (r *Repository) ResetTokens(ctx context.Context, userID int, dailyLimit int) error {
	result, err := r.db.Exec(ctx, `
		UPDATE users
		SET game_tokens = $1, last_token_reset = CURRENT_DATE, updated_at = NOW()
		WHERE id = $2
	`, dailyLimit, userID)

	if err != nil {
		return fmt.Errorf("failed to reset tokens: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// IncrementTokens adds tokens to user's balance
func (r *Repository) IncrementTokens(ctx context.Context, userID int, amount int) error {
	result, err := r.db.Exec(ctx, `
		UPDATE users
		SET game_tokens = game_tokens + $1, updated_at = NOW()
		WHERE id = $2
	`, amount, userID)

	if err != nil {
		return fmt.Errorf("failed to increment tokens: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// RecordTokenPurchase increments purchase counter
func (r *Repository) RecordTokenPurchase(ctx context.Context, userID int, quantity int) error {
	result, err := r.db.Exec(ctx, `
		UPDATE users
		SET tokens_purchased_today = tokens_purchased_today + $1, updated_at = NOW()
		WHERE id = $2
	`, quantity, userID)

	if err != nil {
		return fmt.Errorf("failed to record token purchase: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ResetAllPurchaseCounters resets purchase counters for all users whose last reset date is before current date
func (r *Repository) ResetAllPurchaseCounters(ctx context.Context) (int64, error) {
	result, err := r.db.Exec(ctx, `
		UPDATE users
		SET tokens_purchased_today = 0, updated_at = NOW()
		WHERE last_token_reset < CURRENT_DATE
	`)

	if err != nil {
		return 0, fmt.Errorf("failed to reset purchase counters: %w", err)
	}

	return result.RowsAffected(), nil
}
