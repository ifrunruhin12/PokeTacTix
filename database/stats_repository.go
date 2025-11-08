package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StatsRepository handles player statistics data operations
type StatsRepository struct {
	db *pgxpool.Pool
}

// NewStatsRepository creates a new StatsRepository
func NewStatsRepository(db *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{db: db}
}

// GetOrCreate retrieves player stats or creates if not exists
func (r *StatsRepository) GetOrCreate(ctx context.Context, userID int) (*PlayerStats, error) {
	// Try to get existing stats
	stats, err := r.GetByUserID(ctx, userID)
	if err == nil {
		return stats, nil
	}
	
	// Create new stats if not found
	return r.Create(ctx, userID)
}

// Create creates new player stats
func (r *StatsRepository) Create(ctx context.Context, userID int) (*PlayerStats, error) {
	query := `
		INSERT INTO player_stats (user_id)
		VALUES ($1)
		RETURNING user_id, total_battles_1v1, wins_1v1, losses_1v1,
			total_battles_5v5, wins_5v5, losses_5v5, total_coins_earned, highest_level, updated_at
	`
	
	stats := &PlayerStats{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stats.UserID, &stats.TotalBattles1v1, &stats.Wins1v1, &stats.Losses1v1,
		&stats.TotalBattles5v5, &stats.Wins5v5, &stats.Losses5v5,
		&stats.TotalCoinsEarned, &stats.HighestLevel, &stats.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create player stats: %w", err)
	}
	
	return stats, nil
}

// GetByUserID retrieves player stats by user ID
func (r *StatsRepository) GetByUserID(ctx context.Context, userID int) (*PlayerStats, error) {
	query := `
		SELECT user_id, total_battles_1v1, wins_1v1, losses_1v1,
			total_battles_5v5, wins_5v5, losses_5v5, total_coins_earned, highest_level, updated_at
		FROM player_stats
		WHERE user_id = $1
	`
	
	stats := &PlayerStats{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stats.UserID, &stats.TotalBattles1v1, &stats.Wins1v1, &stats.Losses1v1,
		&stats.TotalBattles5v5, &stats.Wins5v5, &stats.Losses5v5,
		&stats.TotalCoinsEarned, &stats.HighestLevel, &stats.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("player stats not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get player stats: %w", err)
	}
	
	return stats, nil
}

// RecordBattle updates stats after a battle
func (r *StatsRepository) RecordBattle(ctx context.Context, userID int, mode string, result string, coinsEarned int) error {
	// Get or create stats
	stats, err := r.GetOrCreate(ctx, userID)
	if err != nil {
		return err
	}
	
	// Update stats based on mode and result
	if mode == "1v1" {
		stats.TotalBattles1v1++
		if result == "win" {
			stats.Wins1v1++
		} else if result == "loss" {
			stats.Losses1v1++
		}
	} else if mode == "5v5" {
		stats.TotalBattles5v5++
		if result == "win" {
			stats.Wins5v5++
		} else if result == "loss" {
			stats.Losses5v5++
		}
	}
	
	stats.TotalCoinsEarned += coinsEarned
	
	// Update in database
	query := `
		UPDATE player_stats
		SET total_battles_1v1 = $1, wins_1v1 = $2, losses_1v1 = $3,
			total_battles_5v5 = $4, wins_5v5 = $5, losses_5v5 = $6,
			total_coins_earned = $7, updated_at = $8
		WHERE user_id = $9
	`
	
	_, err = r.db.Exec(ctx, query,
		stats.TotalBattles1v1, stats.Wins1v1, stats.Losses1v1,
		stats.TotalBattles5v5, stats.Wins5v5, stats.Losses5v5,
		stats.TotalCoinsEarned, time.Now(), userID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to record battle: %w", err)
	}
	
	return nil
}

// UpdateHighestLevel updates the highest level achieved
func (r *StatsRepository) UpdateHighestLevel(ctx context.Context, userID int, level int) error {
	query := `
		UPDATE player_stats
		SET highest_level = GREATEST(highest_level, $1), updated_at = $2
		WHERE user_id = $3
	`
	
	_, err := r.db.Exec(ctx, query, level, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update highest level: %w", err)
	}
	
	return nil
}

// Update updates player stats
func (r *StatsRepository) Update(ctx context.Context, stats *PlayerStats) error {
	query := `
		UPDATE player_stats
		SET total_battles_1v1 = $1, wins_1v1 = $2, losses_1v1 = $3,
			total_battles_5v5 = $4, wins_5v5 = $5, losses_5v5 = $6,
			total_coins_earned = $7, highest_level = $8, updated_at = $9
		WHERE user_id = $10
	`
	
	stats.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		stats.TotalBattles1v1, stats.Wins1v1, stats.Losses1v1,
		stats.TotalBattles5v5, stats.Wins5v5, stats.Losses5v5,
		stats.TotalCoinsEarned, stats.HighestLevel, stats.UpdatedAt,
		stats.UserID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update player stats: %w", err)
	}
	
	return nil
}

// Delete deletes player stats
func (r *StatsRepository) Delete(ctx context.Context, userID int) error {
	query := `DELETE FROM player_stats WHERE user_id = $1`
	
	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete player stats: %w", err)
	}
	
	return nil
}

// GetAchievements retrieves all achievements with unlock status for a user
func (r *StatsRepository) GetAchievements(ctx context.Context, userID int) ([]*AchievementWithStatus, error) {
	query := `
		SELECT a.id, a.name, a.description, a.icon, a.requirement_type, a.requirement_value,
			CASE WHEN ua.user_id IS NOT NULL THEN TRUE ELSE FALSE END as unlocked,
			ua.unlocked_at
		FROM achievements a
		LEFT JOIN user_achievements ua ON a.id = ua.achievement_id AND ua.user_id = $1
		ORDER BY a.id
	`
	
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get achievements: %w", err)
	}
	defer rows.Close()
	
	var achievements []*AchievementWithStatus
	for rows.Next() {
		ach := &AchievementWithStatus{}
		err := rows.Scan(
			&ach.ID, &ach.Name, &ach.Description, &ach.Icon,
			&ach.RequirementType, &ach.RequirementValue,
			&ach.Unlocked, &ach.UnlockedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan achievement: %w", err)
		}
		achievements = append(achievements, ach)
	}
	
	return achievements, nil
}

// UnlockAchievement unlocks an achievement for a user
func (r *StatsRepository) UnlockAchievement(ctx context.Context, userID int, achievementID int) error {
	query := `
		INSERT INTO user_achievements (user_id, achievement_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, achievement_id) DO NOTHING
	`
	
	_, err := r.db.Exec(ctx, query, userID, achievementID)
	if err != nil {
		return fmt.Errorf("failed to unlock achievement: %w", err)
	}
	
	return nil
}

// CheckAndUnlockAchievements checks and unlocks achievements based on current stats
func (r *StatsRepository) CheckAndUnlockAchievements(ctx context.Context, userID int) ([]*Achievement, error) {
	// Get current stats
	stats, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// Get card stats
	cardRepo := NewCardRepository(r.db)
	highestLevel, _ := cardRepo.GetHighestLevel(ctx, userID)
	legendaryCount, _ := cardRepo.CountLegendary(ctx, userID)
	mythicalCount, _ := cardRepo.CountMythical(ctx, userID)
	
	// Get user
	userRepo := NewUserRepository(r.db)
	user, err := userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// Get all achievements
	allAchievements, err := r.GetAchievements(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	var newlyUnlocked []*Achievement
	totalWins := stats.Wins1v1 + stats.Wins5v5
	
	for _, ach := range allAchievements {
		if ach.Unlocked {
			continue
		}
		
		shouldUnlock := false
		
		switch ach.RequirementType {
		case "total_wins":
			shouldUnlock = totalWins >= ach.RequirementValue
		case "legendary_owned":
			shouldUnlock = legendaryCount >= ach.RequirementValue
		case "mythical_owned":
			shouldUnlock = mythicalCount >= ach.RequirementValue
		case "max_level":
			shouldUnlock = highestLevel >= ach.RequirementValue
		case "coins_total":
			shouldUnlock = user.Coins >= ach.RequirementValue
		}
		
		if shouldUnlock {
			if err := r.UnlockAchievement(ctx, userID, ach.ID); err == nil {
				newlyUnlocked = append(newlyUnlocked, &ach.Achievement)
			}
		}
	}
	
	return newlyUnlocked, nil
}
