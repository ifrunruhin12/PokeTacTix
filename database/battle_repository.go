package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BattleRepository handles battle history data operations
type BattleRepository struct {
	db *pgxpool.Pool
}

// NewBattleRepository creates a new BattleRepository
func NewBattleRepository(db *pgxpool.Pool) *BattleRepository {
	return &BattleRepository{db: db}
}

// Create creates a new battle history record
func (r *BattleRepository) Create(ctx context.Context, battle *BattleHistory) (*BattleHistory, error) {
	query := `
		INSERT INTO battle_history (user_id, mode, result, coins_earned, duration)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	
	err := r.db.QueryRow(ctx, query,
		battle.UserID, battle.Mode, battle.Result, battle.CoinsEarned, battle.Duration,
	).Scan(&battle.ID, &battle.CreatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create battle history: %w", err)
	}
	
	return battle, nil
}

// GetByID retrieves a battle history record by ID
func (r *BattleRepository) GetByID(ctx context.Context, id int) (*BattleHistory, error) {
	query := `
		SELECT id, user_id, mode, result, coins_earned, duration, created_at
		FROM battle_history
		WHERE id = $1
	`
	
	battle := &BattleHistory{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&battle.ID, &battle.UserID, &battle.Mode, &battle.Result,
		&battle.CoinsEarned, &battle.Duration, &battle.CreatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get battle history: %w", err)
	}
	
	return battle, nil
}

// GetUserHistory retrieves battle history for a user
func (r *BattleRepository) GetUserHistory(ctx context.Context, userID int, limit int) ([]*BattleHistory, error) {
	query := `
		SELECT id, user_id, mode, result, coins_earned, duration, created_at
		FROM battle_history
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	
	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user history: %w", err)
	}
	defer rows.Close()
	
	var battles []*BattleHistory
	for rows.Next() {
		battle := &BattleHistory{}
		err := rows.Scan(
			&battle.ID, &battle.UserID, &battle.Mode, &battle.Result,
			&battle.CoinsEarned, &battle.Duration, &battle.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan battle history: %w", err)
		}
		battles = append(battles, battle)
	}
	
	return battles, nil
}

// GetUserHistoryByMode retrieves battle history for a user filtered by mode
func (r *BattleRepository) GetUserHistoryByMode(ctx context.Context, userID int, mode string, limit int) ([]*BattleHistory, error) {
	query := `
		SELECT id, user_id, mode, result, coins_earned, duration, created_at
		FROM battle_history
		WHERE user_id = $1 AND mode = $2
		ORDER BY created_at DESC
		LIMIT $3
	`
	
	rows, err := r.db.Query(ctx, query, userID, mode, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user history by mode: %w", err)
	}
	defer rows.Close()
	
	var battles []*BattleHistory
	for rows.Next() {
		battle := &BattleHistory{}
		err := rows.Scan(
			&battle.ID, &battle.UserID, &battle.Mode, &battle.Result,
			&battle.CoinsEarned, &battle.Duration, &battle.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan battle history: %w", err)
		}
		battles = append(battles, battle)
	}
	
	return battles, nil
}

// Delete deletes a battle history record
func (r *BattleRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM battle_history WHERE id = $1`
	
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete battle history: %w", err)
	}
	
	return nil
}
