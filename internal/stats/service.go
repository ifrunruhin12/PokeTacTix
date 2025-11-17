package stats

import (
	"context"
	"pokemon-cli/internal/database"
)

// Service handles business logic for statistics
type Service struct {
	repo *Repository
}

// NewService creates a new stats service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RecordBattle(ctx context.Context, userID int, mode, result string, coinsEarned, duration int) error {
	return s.repo.RecordBattle(ctx, userID, mode, result, coinsEarned, duration)
}

func (s *Service) GetPlayerStats(ctx context.Context, userID int) (*database.PlayerStats, error) {
	return s.repo.GetPlayerStats(ctx, userID)
}

func (s *Service) GetBattleHistory(ctx context.Context, userID int, limit int) ([]database.BattleHistory, error) {
	if limit <= 0 {
		limit = 20 // Default to 20 battles
	}
	return s.repo.GetBattleHistory(ctx, userID, limit)
}

// UpdateHighestLevel updates the highest level achieved by a player
func (s *Service) UpdateHighestLevel(ctx context.Context, userID int, level int) error {
	return s.repo.UpdateHighestLevel(ctx, userID, level)
}

func (s *Service) GetAchievements(ctx context.Context, userID int) ([]database.AchievementWithStatus, error) {
	return s.repo.GetAchievements(ctx, userID)
}

func (s *Service) CheckAndUnlockAchievements(ctx context.Context, userID int) ([]database.AchievementWithStatus, error) {
	// Get current player stats
	stats, err := s.repo.GetPlayerStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get all achievements
	achievements, err := s.repo.GetAchievements(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate total wins
	totalWins := stats.Wins1v1 + stats.Wins5v5
	totalBattles := stats.TotalBattles1v1 + stats.TotalBattles5v5

	var newlyUnlocked []database.AchievementWithStatus

	// Check each achievement
	for _, ach := range achievements {
		// Skip already unlocked achievements
		if ach.Unlocked {
			continue
		}

		shouldUnlock := false

		switch ach.RequirementType {
		case "battles_1v1":
			shouldUnlock = stats.TotalBattles1v1 >= ach.RequirementValue
		case "battles_5v5":
			shouldUnlock = stats.TotalBattles5v5 >= ach.RequirementValue
		case "total_wins":
			shouldUnlock = totalWins >= ach.RequirementValue
		case "total_battles":
			shouldUnlock = totalBattles >= ach.RequirementValue
		case "wins_1v1":
			shouldUnlock = stats.Wins1v1 >= ach.RequirementValue
		case "wins_5v5":
			shouldUnlock = stats.Wins5v5 >= ach.RequirementValue
		case "total_coins", "coins_total":
			shouldUnlock = stats.TotalCoinsEarned >= ach.RequirementValue
		case "max_level":
			shouldUnlock = stats.HighestLevel >= ach.RequirementValue
		case "legendary_owned":
			// Check if user owns any legendary Pokemon
			hasLegendary, err := s.hasLegendaryPokemon(ctx, userID)
			if err == nil && hasLegendary {
				shouldUnlock = true
			}
		case "mythical_owned":
			// Check if user owns any mythical Pokemon
			hasMythical, err := s.hasMythicalPokemon(ctx, userID)
			if err == nil && hasMythical {
				shouldUnlock = true
			}
		}

		if shouldUnlock {
			err := s.repo.UnlockAchievement(ctx, userID, ach.ID)
			if err != nil {
				// Log error but continue checking other achievements
				continue
			}
			ach.Unlocked = true
			newlyUnlocked = append(newlyUnlocked, ach)
		}
	}

	return newlyUnlocked, nil
}

// hasLegendaryPokemon checks if user owns any legendary Pokemon
func (s *Service) hasLegendaryPokemon(ctx context.Context, userID int) (bool, error) {
	var count int
	err := s.repo.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM player_cards 
		WHERE user_id = $1 AND is_legendary = true
	`, userID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// hasMythicalPokemon checks if user owns any mythical Pokemon
func (s *Service) hasMythicalPokemon(ctx context.Context, userID int) (bool, error) {
	var count int
	err := s.repo.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM player_cards 
		WHERE user_id = $1 AND is_mythical = true
	`, userID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// InitializeAchievements initializes the default achievement definitions
func (s *Service) InitializeAchievements(ctx context.Context) error {
	return s.repo.InitializeAchievements(ctx)
}
