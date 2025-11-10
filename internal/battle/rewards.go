package battle

import (
	"context"
	"fmt"
	"pokemon-cli/internal/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BattleRewards represents rewards earned from a battle
type BattleRewards struct {
	CoinsEarned int            `json:"coins_earned"`
	XPGained    map[int]int    `json:"xp_gained"` // cardID -> XP amount
	LevelUps    []LevelUpInfo  `json:"level_ups"` // Cards that leveled up
}

// LevelUpInfo contains information about a Pokemon that leveled up
type LevelUpInfo struct {
	CardID   int    `json:"card_id"`
	Name     string `json:"name"`
	OldLevel int    `json:"old_level"`
	NewLevel int    `json:"new_level"`
	NewHP    int    `json:"new_hp"`
	NewAttack int   `json:"new_attack"`
	NewDefense int  `json:"new_defense"`
	NewSpeed  int   `json:"new_speed"`
}

// CalculateRewards calculates coins and XP rewards based on battle outcome
func CalculateRewards(bs *BattleState) *BattleRewards {
	rewards := &BattleRewards{
		XPGained: make(map[int]int),
		LevelUps: []LevelUpInfo{},
	}
	
	// Calculate coins based on mode and outcome
	if bs.Winner == "player" {
		if bs.Mode == "1v1" {
			rewards.CoinsEarned = 50
		} else if bs.Mode == "5v5" {
			rewards.CoinsEarned = 150
		}
	} else {
		// Consolation coins for losing
		rewards.CoinsEarned = 10
	}
	
	// Calculate XP for participating Pokemon
	if bs.Winner == "player" {
		if bs.Mode == "1v1" {
			// Award 20 XP to the winning Pokemon
			playerCard := bs.GetActivePlayerCard()
			if playerCard != nil {
				rewards.XPGained[playerCard.CardID] = 20
			}
		} else if bs.Mode == "5v5" {
			// Award 15 XP to each Pokemon that participated (had HP > 0 at some point)
			for _, card := range bs.PlayerDeck {
				// Award XP to all cards (they all participated in 5v5)
				rewards.XPGained[card.CardID] = 15
			}
		}
	}
	
	return rewards
}

// ApplyRewards applies rewards to the database (coins and XP)
func ApplyRewards(ctx context.Context, db *pgxpool.Pool, userID int, bs *BattleState, rewards *BattleRewards) error {
	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	// Update user coins
	_, err = tx.Exec(ctx, `
		UPDATE users 
		SET coins = coins + $1 
		WHERE id = $2
	`, rewards.CoinsEarned, userID)
	if err != nil {
		return fmt.Errorf("failed to update coins: %w", err)
	}
	
	// Apply XP to each card and check for level ups
	for cardID, xp := range rewards.XPGained {
		// Get current card stats
		var currentLevel, currentXP, baseHP, baseAttack, baseDefense, baseSpeed int
		var pokemonName string
		err = tx.QueryRow(ctx, `
			SELECT level, xp, base_hp, base_attack, base_defense, base_speed, pokemon_name
			FROM player_cards
			WHERE id = $1 AND user_id = $2
		`, cardID, userID).Scan(&currentLevel, &currentXP, &baseHP, &baseAttack, &baseDefense, &baseSpeed, &pokemonName)
		if err != nil {
			return fmt.Errorf("failed to get card stats: %w", err)
		}
		
		// Add XP
		newXP := currentXP + xp
		newLevel := currentLevel
		
		// Check for level ups (max level 50)
		for newLevel < 50 {
			xpRequired := 100 * newLevel
			if newXP >= xpRequired {
				newXP -= xpRequired
				newLevel++
			} else {
				break
			}
		}
		
		// If leveled up, calculate new stats
		if newLevel > currentLevel {
			// Calculate new stats based on level
			multiplier := 1.0 + (float64(newLevel-1) * 0.03) // 3% per level for HP
			newHP := int(float64(baseHP) * multiplier)
			newAttack := int(float64(baseAttack) * (1.0 + float64(newLevel-1)*0.02))
			newDefense := int(float64(baseDefense) * (1.0 + float64(newLevel-1)*0.02))
			newSpeed := int(float64(baseSpeed) * (1.0 + float64(newLevel-1)*0.01))
			
			// Update card in database
			_, err = tx.Exec(ctx, `
				UPDATE player_cards
				SET level = $1, xp = $2
				WHERE id = $3 AND user_id = $4
			`, newLevel, newXP, cardID, userID)
			if err != nil {
				return fmt.Errorf("failed to update card level: %w", err)
			}
			
			// Add to level up info
			rewards.LevelUps = append(rewards.LevelUps, LevelUpInfo{
				CardID:     cardID,
				Name:       pokemonName,
				OldLevel:   currentLevel,
				NewLevel:   newLevel,
				NewHP:      newHP,
				NewAttack:  newAttack,
				NewDefense: newDefense,
				NewSpeed:   newSpeed,
			})
		} else {
			// Just update XP
			_, err = tx.Exec(ctx, `
				UPDATE player_cards
				SET xp = $1
				WHERE id = $2 AND user_id = $3
			`, newXP, cardID, userID)
			if err != nil {
				return fmt.Errorf("failed to update card XP: %w", err)
			}
		}
	}
	
	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// RecordBattleHistory records the battle outcome in the database
func RecordBattleHistory(ctx context.Context, db *pgxpool.Pool, userID int, bs *BattleState, coinsEarned int) error {
	result := "loss"
	if bs.Winner == "player" {
		result = "win"
	} else if bs.Winner == "draw" {
		result = "draw"
	}
	
	// Insert battle history
	_, err := db.Exec(ctx, `
		INSERT INTO battle_history (user_id, mode, result, coins_earned, duration)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, bs.Mode, result, coinsEarned, 0) // Duration can be calculated if needed
	if err != nil {
		return fmt.Errorf("failed to record battle history: %w", err)
	}
	
	// Update player stats
	var column string
	if bs.Mode == "1v1" {
		if result == "win" {
			column = "wins_1v1"
		} else {
			column = "losses_1v1"
		}
	} else if bs.Mode == "5v5" {
		if result == "win" {
			column = "wins_5v5"
		} else {
			column = "losses_5v5"
		}
	}
	
	if column != "" {
		_, err = db.Exec(ctx, fmt.Sprintf(`
			INSERT INTO player_stats (user_id, %s, total_coins_earned)
			VALUES ($1, 1, $2)
			ON CONFLICT (user_id) DO UPDATE
			SET %s = player_stats.%s + 1,
			    total_coins_earned = player_stats.total_coins_earned + $2
		`, column, column, column), userID, coinsEarned)
		if err != nil {
			return fmt.Errorf("failed to update player stats: %w", err)
		}
	}
	
	return nil
}

// GetCurrentStats calculates current stats for a card based on level
func GetCurrentStats(baseHP, baseAttack, baseDefense, baseSpeed, level int) database.CardStats {
	multiplier := 1.0 + (float64(level-1) * 0.03) // 3% per level for HP
	return database.CardStats{
		HP:      int(float64(baseHP) * multiplier),
		Attack:  int(float64(baseAttack) * (1.0 + float64(level-1)*0.02)),
		Defense: int(float64(baseDefense) * (1.0 + float64(level-1)*0.02)),
		Speed:   int(float64(baseSpeed) * (1.0 + float64(level-1)*0.01)),
		Stamina: int(float64(baseSpeed) * 2 * (1.0 + float64(level-1)*0.01)),
	}
}
