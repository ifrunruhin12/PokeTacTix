package battle

import (
	"context"
	"encoding/json"
	"fmt"
	"pokemon-cli/internal/database"
	"time"

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
