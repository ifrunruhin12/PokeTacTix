package battle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"pokemon-cli/internal/database"
	"pokemon-cli/internal/items"
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

// SaveBattleSessionWithBoosts reserves the player's active boosts and persists
// the boosted state together. Repeating a battle ID reuses its saved state.
func (r *Repository) SaveBattleSessionWithBoosts(ctx context.Context, state *BattleState) (bool, bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, false, fmt.Errorf("failed to start battle transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var reservedID string
	err = tx.QueryRow(ctx, `
		INSERT INTO battle_boost_reservations (battle_id, user_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING battle_id
	`, state.ID, state.UserID).Scan(&reservedID)
	if errors.Is(err, pgx.ErrNoRows) {
		var stored []byte
		var owner int
		if err := tx.QueryRow(ctx, `SELECT user_id, state_json FROM battle_sessions WHERE session_id = $1`, state.ID).Scan(&owner, &stored); err != nil {
			return false, false, fmt.Errorf("failed to load reserved battle: %w", err)
		}
		if owner != state.UserID {
			return false, false, fmt.Errorf("battle belongs to another user")
		}
		if err := json.Unmarshal(stored, state); err != nil {
			return false, false, fmt.Errorf("failed to decode reserved battle: %w", err)
		}
		var count int
		if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM battle_boost_reservation_items WHERE battle_id = $1`, state.ID).Scan(&count); err != nil {
			return false, false, fmt.Errorf("failed to count reserved boosts: %w", err)
		}
		return count > 0, false, nil
	}
	if err != nil {
		return false, false, fmt.Errorf("failed to reserve battle: %w", err)
	}

	rows, err := tx.Query(ctx, `
		SELECT ab.id, ab.item_id, i.name, ab.stat, ab.bonus, ab.battles_remaining, ab.created_at
		FROM active_boosts ab JOIN items i ON i.id = ab.item_id
		WHERE ab.user_id = $1 ORDER BY ab.id FOR UPDATE OF ab
	`, state.UserID)
	if err != nil {
		return false, false, fmt.Errorf("failed to lock active boosts: %w", err)
	}
	var boosts []items.ActiveBoost
	for rows.Next() {
		var boost items.ActiveBoost
		if err := rows.Scan(&boost.ID, &boost.ItemID, &boost.ItemName, &boost.Stat, &boost.Bonus, &boost.BattlesRemaining, &boost.CreatedAt); err != nil {
			rows.Close()
			return false, false, fmt.Errorf("failed to scan active boost: %w", err)
		}
		boosts = append(boosts, boost)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return false, false, fmt.Errorf("failed to load active boosts: %w", err)
	}

	for _, boost := range boosts {
		_, err := tx.Exec(ctx, `
			INSERT INTO battle_boost_reservation_items (battle_id, boost_id, item_id, stat, bonus, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, state.ID, boost.ID, boost.ItemID, boost.Stat, boost.Bonus, boost.CreatedAt)
		if err != nil {
			return false, false, fmt.Errorf("failed to record reserved boost: %w", err)
		}
		if boost.BattlesRemaining == 1 {
			_, err = tx.Exec(ctx, `DELETE FROM active_boosts WHERE id = $1 AND battles_remaining = 1`, boost.ID)
		} else {
			_, err = tx.Exec(ctx, `UPDATE active_boosts SET battles_remaining = battles_remaining - 1 WHERE id = $1 AND battles_remaining > 1`, boost.ID)
		}
		if err != nil {
			return false, false, fmt.Errorf("failed to reserve boost duration: %w", err)
		}
	}
	ApplyBoosts(state.PlayerDeck, boosts)
	state.UpdatedAt = time.Now()
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return false, false, fmt.Errorf("failed to marshal boosted battle: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO battle_sessions (session_id, user_id, state_json, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, state.ID, state.UserID, stateJSON, state.CreatedAt, state.UpdatedAt); err != nil {
		return false, false, fmt.Errorf("failed to save boosted battle: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return false, false, fmt.Errorf("failed to commit boosted battle: %w", err)
	}
	return len(boosts) > 0, true, nil
}

// CancelBattleBoostReservation restores reserved duration when a saved battle
// cannot start because its token charge failed.
func (r *Repository) CancelBattleBoostReservation(ctx context.Context, battleID string, userID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var finalized bool
	err = tx.QueryRow(ctx, `SELECT finalized FROM battle_boost_reservations WHERE battle_id = $1 AND user_id = $2 FOR UPDATE`, battleID, userID).Scan(&finalized)
	if errors.Is(err, pgx.ErrNoRows) || finalized {
		return nil
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO active_boosts (id, user_id, item_id, stat, bonus, battles_remaining, created_at)
		SELECT boost_id, $2, item_id, stat, bonus, 1, created_at
		FROM battle_boost_reservation_items WHERE battle_id = $1
		ON CONFLICT (id) DO UPDATE SET battles_remaining = active_boosts.battles_remaining + 1
	`, battleID, userID)
	if err != nil {
		return fmt.Errorf("failed to restore boost duration: %w", err)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM battle_sessions WHERE session_id = $1 AND user_id = $2`, battleID, userID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM battle_boost_reservations WHERE battle_id = $1 AND user_id = $2`, battleID, userID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) FinalizeBattleBoostReservation(ctx context.Context, battleID string, userID int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE battle_boost_reservations SET finalized = TRUE
		WHERE battle_id = $1 AND user_id = $2 AND finalized = FALSE
	`, battleID, userID)
	return err
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
