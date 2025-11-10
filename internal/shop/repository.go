package shop

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"pokemon-cli/internal/database"
	"pokemon-cli/internal/pokemon"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles shop data access
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new shop repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// PurchaseCard handles the purchase transaction
func (r *Repository) PurchaseCard(ctx context.Context, userID int, pokemonName string, price int) (*database.PlayerCard, error) {
	// Start transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get user's current coins
	var currentCoins int
	err = tx.QueryRow(ctx, `SELECT coins FROM users WHERE id = $1`, userID).Scan(&currentCoins)
	if err != nil {
		return nil, fmt.Errorf("failed to get user coins: %w", err)
	}

	// Check if user has enough coins
	if currentCoins < price {
		return nil, fmt.Errorf("insufficient coins: have %d, need %d", currentCoins, price)
	}

	// Deduct coins from user
	_, err = tx.Exec(ctx, `
		UPDATE users
		SET coins = coins - $1, updated_at = $2
		WHERE id = $3
	`, price, time.Now(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to deduct coins: %w", err)
	}

	// Fetch Pokemon data from API
	poke, moves, err := pokemon.FetchPokemon(pokemonName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemon: %w", err)
	}

	// Build card from Pokemon data
	card := pokemon.BuildCardFromPokemon(poke, moves)
	
	// Check if legendary or mythical
	isLegendary, isMythical := pokemon.IsLegendaryOrMythical(card.Name)

	// Convert types and moves to JSON
	typesJSON, err := json.Marshal(card.Types)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal types: %w", err)
	}

	movesJSON, err := json.Marshal(card.Moves)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal moves: %w", err)
	}

	// Create player card at level 1
	playerCard := &database.PlayerCard{
		UserID:      userID,
		PokemonName: card.Name,
		Level:       1,
		XP:          0,
		BaseHP:      card.HPMax,
		BaseAttack:  card.Attack,
		BaseDefense: card.Defense,
		BaseSpeed:   card.Speed,
		Types:       typesJSON,
		Moves:       movesJSON,
		Sprite:      card.Sprite,
		IsLegendary: isLegendary,
		IsMythical:  isMythical,
		InDeck:      false,
		DeckPosition: nil,
	}

	// Insert card into database
	query := `
		INSERT INTO player_cards (
			user_id, pokemon_name, level, xp, base_hp, base_attack, base_defense, base_speed,
			types, moves, sprite, is_legendary, is_mythical, in_deck, deck_position
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query,
		playerCard.UserID, playerCard.PokemonName, playerCard.Level, playerCard.XP,
		playerCard.BaseHP, playerCard.BaseAttack, playerCard.BaseDefense, playerCard.BaseSpeed,
		playerCard.Types, playerCard.Moves, playerCard.Sprite,
		playerCard.IsLegendary, playerCard.IsMythical, playerCard.InDeck, playerCard.DeckPosition,
	).Scan(&playerCard.ID, &playerCard.CreatedAt, &playerCard.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return playerCard, nil
}

// GetUserCoins retrieves the user's current coin balance
func (r *Repository) GetUserCoins(ctx context.Context, userID int) (int, error) {
	var coins int
	err := r.db.QueryRow(ctx, `SELECT coins FROM users WHERE id = $1`, userID).Scan(&coins)
	if err != nil {
		return 0, fmt.Errorf("failed to get user coins: %w", err)
	}
	return coins, nil
}
