package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CardRepository handles player card data operations
type CardRepository struct {
	db *pgxpool.Pool
}

// NewCardRepository creates a new CardRepository
func NewCardRepository(db *pgxpool.Pool) *CardRepository {
	return &CardRepository{db: db}
}

// Create creates a new player card
func (r *CardRepository) Create(ctx context.Context, card *PlayerCard) (*PlayerCard, error) {
	query := `
		INSERT INTO player_cards (
			user_id, pokemon_name, level, xp, base_hp, base_attack, base_defense, base_speed,
			types, moves, sprite, is_legendary, is_mythical, in_deck, deck_position
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(ctx, query,
		card.UserID, card.PokemonName, card.Level, card.XP,
		card.BaseHP, card.BaseAttack, card.BaseDefense, card.BaseSpeed,
		card.Types, card.Moves, card.Sprite,
		card.IsLegendary, card.IsMythical, card.InDeck, card.DeckPosition,
	).Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}
	
	return card, nil
}

// GetByID retrieves a card by ID
func (r *CardRepository) GetByID(ctx context.Context, id int) (*PlayerCard, error) {
	query := `
		SELECT id, user_id, pokemon_name, level, xp, base_hp, base_attack, base_defense, base_speed,
			types, moves, sprite, is_legendary, is_mythical, in_deck, deck_position, created_at, updated_at
		FROM player_cards
		WHERE id = $1
	`
	
	card := &PlayerCard{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&card.ID, &card.UserID, &card.PokemonName, &card.Level, &card.XP,
		&card.BaseHP, &card.BaseAttack, &card.BaseDefense, &card.BaseSpeed,
		&card.Types, &card.Moves, &card.Sprite,
		&card.IsLegendary, &card.IsMythical, &card.InDeck, &card.DeckPosition,
		&card.CreatedAt, &card.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("card not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}
	
	return card, nil
}

// GetUserCards retrieves all cards for a user
func (r *CardRepository) GetUserCards(ctx context.Context, userID int) ([]*PlayerCard, error) {
	query := `
		SELECT id, user_id, pokemon_name, level, xp, base_hp, base_attack, base_defense, base_speed,
			types, moves, sprite, is_legendary, is_mythical, in_deck, deck_position, created_at, updated_at
		FROM player_cards
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user cards: %w", err)
	}
	defer rows.Close()
	
	var cards []*PlayerCard
	for rows.Next() {
		card := &PlayerCard{}
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
	
	return cards, nil
}

// GetUserDeck retrieves the user's current deck (5 cards)
func (r *CardRepository) GetUserDeck(ctx context.Context, userID int) ([]*PlayerCard, error) {
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
	
	var cards []*PlayerCard
	for rows.Next() {
		card := &PlayerCard{}
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
	
	return cards, nil
}

// UpdateDeck updates the user's deck configuration
func (r *CardRepository) UpdateDeck(ctx context.Context, userID int, cardIDs []int) error {
	if len(cardIDs) != 5 {
		return fmt.Errorf("deck must contain exactly 5 cards")
	}
	
	// Start transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	// Clear current deck
	_, err = tx.Exec(ctx, `
		UPDATE player_cards
		SET in_deck = FALSE, deck_position = NULL, updated_at = $1
		WHERE user_id = $2
	`, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to clear deck: %w", err)
	}
	
	// Set new deck
	for i, cardID := range cardIDs {
		_, err = tx.Exec(ctx, `
			UPDATE player_cards
			SET in_deck = TRUE, deck_position = $1, updated_at = $2
			WHERE id = $3 AND user_id = $4
		`, i+1, time.Now(), cardID, userID)
		if err != nil {
			return fmt.Errorf("failed to update deck card: %w", err)
		}
	}
	
	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// AddXP adds XP to a card and handles level-ups
func (r *CardRepository) AddXP(ctx context.Context, cardID int, xp int) (*PlayerCard, error) {
	// Get current card
	card, err := r.GetByID(ctx, cardID)
	if err != nil {
		return nil, err
	}
	
	// Add XP
	card.XP += xp
	
	// Check for level-ups
	for card.Level < 50 {
		xpRequired := 100 * card.Level
		if card.XP >= xpRequired {
			card.XP -= xpRequired
			card.Level++
		} else {
			break
		}
	}
	
	// Cap XP at max level
	if card.Level >= 50 {
		card.XP = 0
	}
	
	// Update card
	query := `
		UPDATE player_cards
		SET level = $1, xp = $2, updated_at = $3
		WHERE id = $4
	`
	
	_, err = r.db.Exec(ctx, query, card.Level, card.XP, time.Now(), cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to update card XP: %w", err)
	}
	
	return card, nil
}

// Update updates a card
func (r *CardRepository) Update(ctx context.Context, card *PlayerCard) error {
	query := `
		UPDATE player_cards
		SET pokemon_name = $1, level = $2, xp = $3, base_hp = $4, base_attack = $5,
			base_defense = $6, base_speed = $7, types = $8, moves = $9, sprite = $10,
			is_legendary = $11, is_mythical = $12, in_deck = $13, deck_position = $14, updated_at = $15
		WHERE id = $16
	`
	
	card.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		card.PokemonName, card.Level, card.XP,
		card.BaseHP, card.BaseAttack, card.BaseDefense, card.BaseSpeed,
		card.Types, card.Moves, card.Sprite,
		card.IsLegendary, card.IsMythical, card.InDeck, card.DeckPosition,
		card.UpdatedAt, card.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}
	
	return nil
}

// Delete deletes a card
func (r *CardRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM player_cards WHERE id = $1`
	
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}
	
	return nil
}

// GetHighestLevel returns the highest level card for a user
func (r *CardRepository) GetHighestLevel(ctx context.Context, userID int) (int, error) {
	query := `
		SELECT COALESCE(MAX(level), 1)
		FROM player_cards
		WHERE user_id = $1
	`
	
	var maxLevel int
	err := r.db.QueryRow(ctx, query, userID).Scan(&maxLevel)
	if err != nil {
		return 1, fmt.Errorf("failed to get highest level: %w", err)
	}
	
	return maxLevel, nil
}

// CountLegendary counts legendary cards owned by user
func (r *CardRepository) CountLegendary(ctx context.Context, userID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM player_cards
		WHERE user_id = $1 AND is_legendary = TRUE
	`
	
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count legendary cards: %w", err)
	}
	
	return count, nil
}

// CountMythical counts mythical cards owned by user
func (r *CardRepository) CountMythical(ctx context.Context, userID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM player_cards
		WHERE user_id = $1 AND is_mythical = TRUE
	`
	
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count mythical cards: %w", err)
	}
	
	return count, nil
}
