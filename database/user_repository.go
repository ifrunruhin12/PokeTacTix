package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, username, email, passwordHash string) (*User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, coins)
		VALUES ($1, $2, $3, 0)
		RETURNING id, username, email, password_hash, coins, created_at, updated_at
	`
	
	user := &User{}
	err := r.db.QueryRow(ctx, query, username, email, passwordHash).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Coins,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, coins, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	user := &User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Coins,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, coins, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	
	user := &User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Coins,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, coins, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	
	user := &User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Coins,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

// Update updates user information
func (r *UserRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, coins = $3, updated_at = $4
		WHERE id = $5
	`
	
	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query, user.Username, user.Email, user.Coins, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	return nil
}

// UpdateCoins updates user's coin balance
func (r *UserRepository) UpdateCoins(ctx context.Context, userID int, coins int) error {
	query := `
		UPDATE users
		SET coins = $1, updated_at = $2
		WHERE id = $3
	`
	
	_, err := r.db.Exec(ctx, query, coins, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update coins: %w", err)
	}
	
	return nil
}

// AddCoins adds coins to user's balance
func (r *UserRepository) AddCoins(ctx context.Context, userID int, amount int) error {
	query := `
		UPDATE users
		SET coins = coins + $1, updated_at = $2
		WHERE id = $3
	`
	
	_, err := r.db.Exec(ctx, query, amount, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to add coins: %w", err)
	}
	
	return nil
}

// Delete deletes a user (cascade will delete related data)
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	return nil
}

// UsernameExists checks if a username already exists
func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	
	var exists bool
	err := r.db.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username: %w", err)
	}
	
	return exists, nil
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email: %w", err)
	}
	
	return exists, nil
}
