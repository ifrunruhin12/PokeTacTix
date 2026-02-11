package battle

import (
	"context"
	"encoding/json"
	"fmt"
	"pokemon-cli/internal/database"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for battle sessions
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new battle repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// SaveBattleSession saves a battle state to the database
func (r *Repository) SaveBattleSession(ctx context.Context, state *BattleState) error {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal battle state: %w", err)
	}

	query := `
		INSERT INTO battle_sessions (session_id, user_id, state_json, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (session_id) 
		DO UPDATE SET 
			state_json = EXCLUDED.state_json,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.db.Exec(ctx, query, state.ID, state.UserID, stateJSON, state.CreatedAt, state.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save battle session: %w", err)
	}

	return nil
}

// GetBattleSession retrieves a battle state from the database
func (r *Repository) GetBattleSession(ctx context.Context, sessionID string) (*BattleState, error) {
	query := `
		SELECT state_json 
		FROM battle_sessions 
		WHERE session_id = $1
	`

	var stateJSON []byte
	err := r.db.QueryRow(ctx, query, sessionID).Scan(&stateJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to get battle session: %w", err)
	}

	var state BattleState
	err = json.Unmarshal(stateJSON, &state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal battle state: %w", err)
	}

	return &state, nil
}

// DeleteBattleSession removes a battle session from the database
func (r *Repository) DeleteBattleSession(ctx context.Context, sessionID string) error {
	query := `DELETE FROM battle_sessions WHERE session_id = $1`

	_, err := r.db.Exec(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete battle session: %w", err)
	}

	return nil
}

// GetUserBattleSessions retrieves all battle sessions for a user
func (r *Repository) GetUserBattleSessions(ctx context.Context, userID int) ([]*BattleState, error) {
	query := `
		SELECT state_json 
		FROM battle_sessions 
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user battle sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*BattleState
	for rows.Next() {
		var stateJSON []byte
		if err := rows.Scan(&stateJSON); err != nil {
			return nil, fmt.Errorf("failed to scan battle session: %w", err)
		}

		var state BattleState
		if err := json.Unmarshal(stateJSON, &state); err != nil {
			return nil, fmt.Errorf("failed to unmarshal battle state: %w", err)
		}

		sessions = append(sessions, &state)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating battle sessions: %w", err)
	}

	return sessions, nil
}

// CleanupExpiredSessions removes battle sessions older than the specified duration
func (r *Repository) CleanupExpiredSessions(ctx context.Context, expiryDuration time.Duration) (int64, error) {
	expiryTime := time.Now().Add(-expiryDuration)

	query := `
		DELETE FROM battle_sessions 
		WHERE updated_at < $1
	`

	result, err := r.db.Exec(ctx, query, expiryTime)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	rowsAffected := result.RowsAffected()
	return rowsAffected, nil
}

// GetUserDeck retrieves the user's current deck from player_cards table
func (r *Repository) GetUserDeck(ctx context.Context, userID int) ([]database.PlayerCard, error) {
	query := `
		SELECT id, user_id, pokemon_name, level, xp, base_hp, base_attack, base_defense, base_speed,
			types, moves, sprite, is_legendary, is_mythical, in_deck, deck_position, created_at, updated_at
		FROM player_cards
		WHERE user_id = $1 AND in_deck = TRUE
		ORDER BY deck_position ASC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user deck: %w", err)
	}
	defer rows.Close()

	var cards []database.PlayerCard
	for rows.Next() {
		var card database.PlayerCard
		err := rows.Scan(
			&card.ID, &card.UserID, &card.PokemonName, &card.Level, &card.XP,
			&card.BaseHP, &card.BaseAttack, &card.BaseDefense, &card.BaseSpeed,
			&card.Types, &card.Moves, &card.Sprite,
			&card.IsLegendary, &card.IsMythical, &card.InDeck, &card.DeckPosition,
			&card.CreatedAt, &card.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cards: %w", err)
	}

	return cards, nil
}

func (r *Repository) RecordBattleHistory(ctx context.Context, userID int, mode, result string, coinsEarned, duration int) error {
	query := `
		INSERT INTO battle_history (user_id, mode, result, coins_earned, duration)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query, userID, mode, result, coinsEarned, duration)
	if err != nil {
		return fmt.Errorf("failed to record battle history: %w", err)
	}

	return nil
}

func (r *Repository) UpdatePlayerStats(ctx context.Context, userID int, mode, result string, coinsEarned int) error {
	// Validate mode to prevent SQL injection
	if mode != "1v1" && mode != "5v5" {
		return fmt.Errorf("invalid battle mode: %s", mode)
	}

	// Validate result to prevent SQL injection
	if result != "win" && result != "loss" && result != "draw" {
		return fmt.Errorf("invalid battle result: %s", result)
	}

	// Use separate parameterized queries for each mode and result combination
	var query string

	if mode == "1v1" {
		switch result {
		case "win":
			query = `
				INSERT INTO player_stats (user_id, total_battles_1v1, wins_1v1, total_coins_earned, updated_at)
				VALUES ($1, 1, 1, $2, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_1v1 = player_stats.total_battles_1v1 + 1,
				    wins_1v1 = player_stats.wins_1v1 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
				    updated_at = NOW()
			`
		case "loss":
			query = `
				INSERT INTO player_stats (user_id, total_battles_1v1, losses_1v1, total_coins_earned, updated_at)
				VALUES ($1, 1, 1, $2, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_1v1 = player_stats.total_battles_1v1 + 1,
				    losses_1v1 = player_stats.losses_1v1 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
				    updated_at = NOW()
			`
		default:
			query = `
				INSERT INTO player_stats (user_id, total_battles_1v1, total_coins_earned, updated_at)
				VALUES ($1, 1, $2, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_1v1 = player_stats.total_battles_1v1 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
				    updated_at = NOW()
			`
		}
	} else { // mode == "5v5"
		switch result {
		case "win":
			query = `
				INSERT INTO player_stats (user_id, total_battles_5v5, wins_5v5, total_coins_earned, updated_at)
				VALUES ($1, 1, 1, $2, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_5v5 = player_stats.total_battles_5v5 + 1,
				    wins_5v5 = player_stats.wins_5v5 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
				    updated_at = NOW()
			`
		case "loss":
			query = `
				INSERT INTO player_stats (user_id, total_battles_5v5, losses_5v5, total_coins_earned, updated_at)
				VALUES ($1, 1, 1, $2, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_5v5 = player_stats.total_battles_5v5 + 1,
				    losses_5v5 = player_stats.losses_5v5 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
				    updated_at = NOW()
			`
		default:
			query = `
				INSERT INTO player_stats (user_id, total_battles_5v5, total_coins_earned, updated_at)
				VALUES ($1, 1, $2, NOW())
				ON CONFLICT (user_id) DO UPDATE
				SET total_battles_5v5 = player_stats.total_battles_5v5 + 1,
				    total_coins_earned = player_stats.total_coins_earned + $2,
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

func (r *Repository) UpdateHighestLevel(ctx context.Context, userID int, level int) error {
	query := `
		INSERT INTO player_stats (user_id, highest_level, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id) DO UPDATE
		SET highest_level = GREATEST(player_stats.highest_level, $2),
		    updated_at = NOW()
	`

	_, err := r.db.Exec(ctx, query, userID, level)
	if err != nil {
		return fmt.Errorf("failed to update highest level: %w", err)
	}

	return nil
}

func (r *Repository) UpdatePlayerStatsInTx(ctx context.Context, tx pgx.Tx, userID int, mode, result string, coinsEarned int) error {
	// Validate inputs
	if mode != "1v1" && mode != "5v5" {
		return fmt.Errorf("invalid battle mode: %s", mode)
	}
	if result != "win" && result != "loss" && result != "draw" {
		return fmt.Errorf("invalid battle result: %s", result)
	}

	// Build query based on mode and result
	var query string
	if mode == "1v1" {
		switch result {
		case "win":
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
		case "loss":
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
		default: // draw
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
	} else { // 5v5
		switch result {
		case "win":
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
		case "loss":
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
		default: // draw
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

	_, err := tx.Exec(ctx, query, userID, coinsEarned)
	if err != nil {
		return fmt.Errorf("failed to update player stats: %w", err)
	}

	return nil
}
