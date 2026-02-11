package stats

import (
	"context"
	"errors"
	"fmt"
	"pokemon-cli/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for statistics
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new stats repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) RecordBattle(ctx context.Context, userID int, mode, result string, coinsEarned, duration int) error {
	// Insert battle history
	_, err := r.db.Exec(ctx, `
		INSERT INTO battle_history (user_id, mode, result, coins_earned, duration)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, mode, result, coinsEarned, duration)
	if err != nil {
		return fmt.Errorf("failed to record battle history: %w", err)
	}

	// Update player stats based on mode and result
	return r.updatePlayerStats(ctx, userID, mode, result, coinsEarned)
}

// updatePlayerStats updates the player_stats table with win/loss counts
// This function uses separate queries for each mode to prevent SQL injection
func (r *Repository) updatePlayerStats(ctx context.Context, userID int, mode, result string, coinsEarned int) error {
	// Validate mode to prevent SQL injection
	if mode != "1v1" && mode != "5v5" {
		return fmt.Errorf("invalid battle mode: %s", mode)
	}

	// Validate result to prevent SQL injection
	if result != "win" && result != "loss" && result != "draw" {
		return fmt.Errorf("invalid battle result: %s", result)
	}

	// Use separate parameterized queries for each mode and result combination
	// This prevents SQL injection by never concatenating user input into queries
	var query string

	if mode == "1v1" {
		if result == "win" {
			query = `
				INSERT INTO player_stats (user_id, total_battles_1v1, wins_1v1, total_coins_earned, consecutive_losses, updated_at)
				VALUES ($1, 1, 1, $2, 0, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_1v1 = player_stats.total_battles_1v1 + 1,
				    wins_1v1 = player_stats.wins_1v1 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
				    consecutive_losses = 0,
				    updated_at = NOW()
			`
		} else if result == "loss" {
			query = `
				INSERT INTO player_stats (user_id, total_battles_1v1, losses_1v1, total_coins_earned, consecutive_losses, updated_at)
				VALUES ($1, 1, 1, $2, 1, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_1v1 = player_stats.total_battles_1v1 + 1,
				    losses_1v1 = player_stats.losses_1v1 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
				    consecutive_losses = player_stats.consecutive_losses + 1,
				    updated_at = NOW()
			`
		} else { // result == "draw"
			query = `
				INSERT INTO player_stats (user_id, total_battles_1v1, draws_1v1, total_coins_earned, consecutive_losses, updated_at)
				VALUES ($1, 1, 1, $2, 0, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_1v1 = player_stats.total_battles_1v1 + 1,
				    draws_1v1 = player_stats.draws_1v1 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
				    consecutive_losses = 0,
				    updated_at = NOW()
			`
		}
	} else { // mode == "5v5"
		if result == "win" {
			query = `
				INSERT INTO player_stats (user_id, total_battles_5v5, wins_5v5, total_coins_earned, consecutive_losses, updated_at)
				VALUES ($1, 1, 1, $2, 0, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_5v5 = player_stats.total_battles_5v5 + 1,
				    wins_5v5 = player_stats.wins_5v5 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
				    consecutive_losses = 0,
				    updated_at = NOW()
			`
		} else if result == "loss" {
			query = `
				INSERT INTO player_stats (user_id, total_battles_5v5, losses_5v5, total_coins_earned, consecutive_losses, updated_at)
				VALUES ($1, 1, 1, $2, 1, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_5v5 = player_stats.total_battles_5v5 + 1,
				    losses_5v5 = player_stats.losses_5v5 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
				    consecutive_losses = player_stats.consecutive_losses + 1,
				    updated_at = NOW()
			`
		} else { // result == "draw"
			query = `
				INSERT INTO player_stats (user_id, total_battles_5v5, draws_5v5, total_coins_earned, consecutive_losses, updated_at)
				VALUES ($1, 1, 1, $2, 0, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_5v5 = player_stats.total_battles_5v5 + 1,
				    draws_5v5 = player_stats.draws_5v5 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
				    consecutive_losses = 0,
				    updated_at = NOW()
			`
		}
	}

	_, err := r.db.Exec(ctx, query, userID, coinsEarned)
	if err != nil {
		return fmt.Errorf("failed to update player stats: %w", err)
	}

	return nil
}

// GetPlayerStats retrieves player statistics
func (r *Repository) GetPlayerStats(ctx context.Context, userID int) (*database.PlayerStats, error) {
	stats := &database.PlayerStats{UserID: userID}

	err := r.db.QueryRow(ctx, `
		SELECT 
			COALESCE(total_battles_1v1, 0),
			COALESCE(wins_1v1, 0),
			COALESCE(losses_1v1, 0),
			COALESCE(draws_1v1, 0),
			COALESCE(total_battles_5v5, 0),
			COALESCE(wins_5v5, 0),
			COALESCE(losses_5v5, 0),
			COALESCE(draws_5v5, 0),
			COALESCE(total_coins_earned, 0),
			COALESCE(highest_level, 1),
			COALESCE(consecutive_losses, 0),
			updated_at
		FROM player_stats
		WHERE user_id = $1
	`, userID).Scan(
		&stats.TotalBattles1v1,
		&stats.Wins1v1,
		&stats.Losses1v1,
		&stats.Draws1v1,
		&stats.TotalBattles5v5,
		&stats.Wins5v5,
		&stats.Losses5v5,
		&stats.Draws5v5,
		&stats.TotalCoinsEarned,
		&stats.HighestLevel,
		&stats.ConsecutiveLosses,
		&stats.UpdatedAt,
	)

	if err != nil {
		// If no stats exist yet, return empty stats with default values
		if errors.Is(err, pgx.ErrNoRows) {
			// Return default stats for new users
			stats.TotalBattles1v1 = 0
			stats.Wins1v1 = 0
			stats.Losses1v1 = 0
			stats.Draws1v1 = 0
			stats.TotalBattles5v5 = 0
			stats.Wins5v5 = 0
			stats.Losses5v5 = 0
			stats.Draws5v5 = 0
			stats.TotalCoinsEarned = 0
			stats.HighestLevel = 1
			stats.ConsecutiveLosses = 0
			return stats, nil
		}
		return nil, fmt.Errorf("failed to get player stats: %w", err)
	}

	return stats, nil
}

// GetBattleHistory retrieves the last N battles for a user
func (r *Repository) GetBattleHistory(ctx context.Context, userID int, limit int) ([]database.BattleHistory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, mode, result, coins_earned, duration, created_at
		FROM battle_history
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query battle history: %w", err)
	}
	defer rows.Close()

	var history []database.BattleHistory
	for rows.Next() {
		var battle database.BattleHistory
		err := rows.Scan(
			&battle.ID,
			&battle.UserID,
			&battle.Mode,
			&battle.Result,
			&battle.CoinsEarned,
			&battle.Duration,
			&battle.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan battle history: %w", err)
		}
		history = append(history, battle)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating battle history: %w", err)
	}

	return history, nil
}

// UpdateHighestLevel updates the highest level achieved by a player
func (r *Repository) UpdateHighestLevel(ctx context.Context, userID int, level int) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO player_stats (user_id, highest_level, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id) DO UPDATE
		SET highest_level = GREATEST(player_stats.highest_level, $2),
		    updated_at = NOW()
	`, userID, level)

	if err != nil {
		return fmt.Errorf("failed to update highest level: %w", err)
	}

	return nil
}

// GetAchievements retrieves all achievements with unlock status for a user
// Requirements: 8.4
func (r *Repository) GetAchievements(ctx context.Context, userID int) ([]database.AchievementWithStatus, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
			a.id, a.name, a.description, a.icon, a.requirement_type, a.requirement_value,
			ua.user_id IS NOT NULL as unlocked,
			ua.unlocked_at
		FROM achievements a
		LEFT JOIN user_achievements ua ON a.id = ua.achievement_id AND ua.user_id = $1
		ORDER BY a.id
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query achievements: %w", err)
	}
	defer rows.Close()

	var achievements []database.AchievementWithStatus
	for rows.Next() {
		var ach database.AchievementWithStatus
		err := rows.Scan(
			&ach.ID,
			&ach.Name,
			&ach.Description,
			&ach.Icon,
			&ach.RequirementType,
			&ach.RequirementValue,
			&ach.Unlocked,
			&ach.UnlockedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan achievement: %w", err)
		}
		achievements = append(achievements, ach)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating achievements: %w", err)
	}

	return achievements, nil
}

// UnlockAchievement unlocks an achievement for a user
func (r *Repository) UnlockAchievement(ctx context.Context, userID, achievementID int) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_achievements (user_id, achievement_id, unlocked_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id, achievement_id) DO NOTHING
	`, userID, achievementID)

	if err != nil {
		return fmt.Errorf("failed to unlock achievement: %w", err)
	}

	return nil
}

// GetUserAchievementCount returns the count of unlocked achievements for a user
func (r *Repository) GetUserAchievementCount(ctx context.Context, userID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM user_achievements WHERE user_id = $1
	`, userID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get achievement count: %w", err)
	}

	return count, nil
}

// InitializeAchievements creates the default achievement definitions
// This should be called during database setup/migration
func (r *Repository) InitializeAchievements(ctx context.Context) error {
	achievements := []struct {
		name             string
		description      string
		icon             string
		requirementType  string
		requirementValue int
	}{
		{"First 1v1 Battle", "Complete your first 1v1 battle", "🎯", "battles_1v1", 1},
		{"First 5v5 Battle", "Complete your first 5v5 battle", "🎲", "battles_5v5", 1},
		{"First Victory", "Win your first battle", "🏆", "total_wins", 1},
		{"Veteran Trainer", "Win 10 battles", "⭐", "total_wins", 10},
		{"Elite Trainer", "Win 50 battles", "💫", "total_wins", 50},
		{"Champion", "Win 100 battles", "👑", "total_wins", 100},
		{"Legendary Collector", "Obtain a legendary Pokemon", "🌟", "legendary_owned", 1},
		{"Mythical Master", "Obtain a mythical Pokemon", "✨", "mythical_owned", 1},
		{"Max Level", "Get a Pokemon to level 50", "📈", "max_level", 50},
		{"Coin Hoarder", "Accumulate 5000 coins", "💰", "total_coins", 5000},
		{"Battle Enthusiast", "Complete 25 battles", "⚔️", "total_battles", 25},
		{"5v5 Specialist", "Win 20 5v5 battles", "🎯", "wins_5v5", 20},
	}

	for _, ach := range achievements {
		_, err := r.db.Exec(ctx, `
			INSERT INTO achievements (name, description, icon, requirement_type, requirement_value)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (name) DO NOTHING
		`, ach.name, ach.description, ach.icon, ach.requirementType, ach.requirementValue)

		if err != nil {
			return fmt.Errorf("failed to initialize achievement %s: %w", ach.name, err)
		}
	}

	return nil
}
