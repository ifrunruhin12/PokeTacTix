package cards

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"pokemon-cli/internal/database"
	"pokemon-cli/internal/items"
	"pokemon-cli/internal/pokemon"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles card data access
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new card repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create creates a new player card
func (r *Repository) Create(ctx context.Context, card *database.PlayerCard) (*database.PlayerCard, error) {
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
func (r *Repository) GetByID(ctx context.Context, id int) (*database.PlayerCard, error) {
	query := `
		SELECT id, user_id, pokemon_name, level, xp, base_hp, base_attack, base_defense, base_speed,
			types, moves, sprite, is_legendary, is_mythical, in_deck, deck_position, created_at, updated_at
		FROM player_cards
		WHERE id = $1
	`

	card := &database.PlayerCard{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&card.ID, &card.UserID, &card.PokemonName, &card.Level, &card.XP,
		&card.BaseHP, &card.BaseAttack, &card.BaseDefense, &card.BaseSpeed,
		&card.Types, &card.Moves, &card.Sprite,
		&card.IsLegendary, &card.IsMythical, &card.InDeck, &card.DeckPosition,
		&card.CreatedAt, &card.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCardNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	return card, nil
}

// GetUserCards retrieves all cards for a user
func (r *Repository) GetUserCards(ctx context.Context, userID int) ([]database.PlayerCard, error) {
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

	return cards, nil
}

// GetUserDeck retrieves the user's current deck
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

	return cards, nil
}

// UpdateDeck updates the user's deck configuration
func (r *Repository) UpdateDeck(ctx context.Context, userID int, cardIDs []int) error {
	if len(cardIDs) < 1 || len(cardIDs) > 5 {
		return fmt.Errorf("deck must contain between 1 and 5 cards")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE player_cards
		SET in_deck = FALSE, deck_position = NULL, updated_at = $1
		WHERE user_id = $2
	`, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to clear deck: %w", err)
	}

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

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// AddXP adds XP to a card and handles level-ups
func (r *Repository) AddXP(ctx context.Context, cardID int, xp int) (*database.PlayerCard, error) {
	card, err := r.GetByID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	card.XP += xp

	for card.Level < 50 {
		xpRequired := 100 * card.Level
		if card.XP >= xpRequired {
			card.XP -= xpRequired
			card.Level++
		} else {
			break
		}
	}

	if card.Level >= 50 {
		card.XP = 0
	}

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
func (r *Repository) Update(ctx context.Context, card *database.PlayerCard) error {
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
func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM player_cards WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	return nil
}

// GetHighestLevel gets the highest level card for a user
func (r *Repository) GetHighestLevel(ctx context.Context, userID int) (int, error) {
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

// CountLegendary counts legendary cards for a user
func (r *Repository) CountLegendary(ctx context.Context, userID int) (int, error) {
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

// CountMythical counts mythical cards for a user
func (r *Repository) CountMythical(ctx context.Context, userID int) (int, error) {
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

// GetPokemonID returns the card's canonical pokemon row id, or nil for legacy
// cards created before the link column existed.
func (r *Repository) GetPokemonID(ctx context.Context, cardID, userID int) (*int, error) {
	var pokemonID *int
	err := r.db.QueryRow(ctx, `
		SELECT pokemon_id FROM player_cards WHERE id = $1 AND user_id = $2
	`, cardID, userID).Scan(&pokemonID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCardNotFound
		}
		return nil, fmt.Errorf("failed to get pokemon_id: %w", err)
	}
	return pokemonID, nil
}

// SetPokemonID backfills the canonical pokemon link on a legacy card
func (r *Repository) SetPokemonID(ctx context.Context, cardID, userID, pokemonID int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE player_cards SET pokemon_id = $1 WHERE id = $2 AND user_id = $3
	`, pokemonID, cardID, userID)
	if err != nil {
		return fmt.Errorf("failed to set pokemon_id: %w", err)
	}
	return nil
}

// EvolveWithItem applies an item-based evolution and consumes one unit of the
// item in a single transaction:
//
//  1. The card row is locked (FOR UPDATE) and its pokemon_id re-verified, so a
//     concurrent evolution of the same card cannot double-apply.
//  2. The item is consumed with a guarded UPDATE (quantity > 0), so inventory
//     can never go negative and a failed evolution consumes nothing (any later
//     failure rolls the whole transaction back).
//  3. The card's species is updated with the target's base stats, exactly like
//     the level-up path in the battle package (moves carry over).
func (r *Repository) EvolveWithItem(ctx context.Context, userID, cardID, expectedPokemonID int, target *pokemon.Pokemon, itemID string) (int, error) {
	baseHP, baseAttack, baseDefense, baseSpeed := target.CardBaseStats()

	typesJSON, err := json.Marshal(target.Types)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal evolved types: %w", err)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Lock the card and re-verify it is still the species the evolution rule
	// was resolved for.
	var currentPokemonID *int
	err = tx.QueryRow(ctx, `
		SELECT pokemon_id FROM player_cards WHERE id = $1 AND user_id = $2 FOR UPDATE
	`, cardID, userID).Scan(&currentPokemonID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrCardNotFound
		}
		return 0, fmt.Errorf("failed to lock card %d: %w", cardID, err)
	}
	if currentPokemonID == nil || *currentPokemonID != expectedPokemonID {
		return 0, ErrConcurrentEvolution
	}

	// Consume one unit of the item; the quantity > 0 guard makes this safe
	// under concurrent use.
	var remaining int
	err = tx.QueryRow(ctx, `
		UPDATE player_inventory
		SET quantity = quantity - 1, updated_at = $1
		WHERE user_id = $2 AND item_id = $3 AND quantity > 0
		RETURNING quantity
	`, time.Now(), userID, itemID).Scan(&remaining)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%w: %s", items.ErrInsufficientItem, itemID)
		}
		return 0, fmt.Errorf("failed to consume item %s: %w", itemID, err)
	}

	// Apply the species change. Moves intentionally carry over: evolution is a
	// species change with new base stats, not a move-set reset (same policy as
	// level-up evolution).
	_, err = tx.Exec(ctx, `
		UPDATE player_cards
		SET pokemon_name = $1, pokemon_id = $2, sprite = $3, types = $4,
		    base_hp = $5, base_attack = $6, base_defense = $7, base_speed = $8
		WHERE id = $9 AND user_id = $10
	`,
		target.Name, target.ID, target.SpriteURL, typesJSON,
		baseHP, baseAttack, baseDefense, baseSpeed,
		cardID, userID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to evolve card %d into %s: %w", cardID, target.Name, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return remaining, nil
}
