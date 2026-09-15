package items

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// ErrItemNotFound is returned when an item id is not in the catalog
	ErrItemNotFound = errors.New("item not found")

	// ErrInsufficientCoins is returned when the user cannot afford the purchase
	ErrInsufficientCoins = errors.New("insufficient coins")

	// ErrInsufficientItem is returned when the inventory has none of the item left
	ErrInsufficientItem = errors.New("insufficient item quantity")

	// ErrNotBooster is returned when the item being used is not a booster
	ErrNotBooster = errors.New("item is not usable here")
)

// Repository handles item catalog and player inventory data access
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new items repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const itemColumns = `id, name, description, item_type, price, icon, effect, created_at, updated_at`

func scanItem(row pgx.Row) (*Item, error) {
	var item Item
	var effect []byte
	if err := row.Scan(
		&item.ID, &item.Name, &item.Description, &item.ItemType, &item.Price,
		&item.Icon, &effect, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if len(effect) > 0 {
		item.Effect = effect
	}
	return &item, nil
}

// ListItems returns the full item catalog
func (r *Repository) ListItems(ctx context.Context) ([]Item, error) {
	rows, err := r.db.Query(ctx, `SELECT `+itemColumns+` FROM items ORDER BY price ASC, id ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}
	defer rows.Close()

	var result []Item
	for rows.Next() {
		var item Item
		var effect []byte
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &item.ItemType, &item.Price,
			&item.Icon, &effect, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		if len(effect) > 0 {
			item.Effect = effect
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

// GetItem returns a single catalog item by id
func (r *Repository) GetItem(ctx context.Context, itemID string) (*Item, error) {
	item, err := scanItem(r.db.QueryRow(ctx, `SELECT `+itemColumns+` FROM items WHERE id = $1`, itemID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to get item %s: %w", itemID, err)
	}
	return item, nil
}

// GetInventory returns the player's inventory joined with catalog details.
// Items the player has never owned are not included.
func (r *Repository) GetInventory(ctx context.Context, userID int) ([]InventoryEntry, error) {
	rows, err := r.db.Query(ctx, `
		SELECT i.id, i.name, i.description, i.item_type, i.price, i.icon, i.effect, i.created_at, i.updated_at,
		       pi.quantity
		FROM player_inventory pi
		JOIN items i ON i.id = pi.item_id
		WHERE pi.user_id = $1 AND pi.quantity > 0
		ORDER BY i.item_type ASC, i.name ASC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	defer rows.Close()

	var result []InventoryEntry
	for rows.Next() {
		var entry InventoryEntry
		var effect []byte
		if err := rows.Scan(
			&entry.ID, &entry.Name, &entry.Description, &entry.ItemType, &entry.Price,
			&entry.Icon, &effect, &entry.CreatedAt, &entry.UpdatedAt,
			&entry.Quantity,
		); err != nil {
			return nil, fmt.Errorf("failed to scan inventory entry: %w", err)
		}
		if len(effect) > 0 {
			entry.Effect = effect
		}
		result = append(result, entry)
	}
	return result, rows.Err()
}

// GetQuantity returns the player's current quantity of an item (0 when none)
func (r *Repository) GetQuantity(ctx context.Context, userID int, itemID string) (int, error) {
	var quantity int
	err := r.db.QueryRow(ctx, `
		SELECT quantity FROM player_inventory WHERE user_id = $1 AND item_id = $2
	`, userID, itemID).Scan(&quantity)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get item quantity: %w", err)
	}
	return quantity, nil
}

// Purchase charges the user's coins and adds the items to their inventory in
// one transaction. The coin deduction is guarded (`coins >= cost`) so a
// concurrent purchase can never drive the balance negative, and the inventory
// upsert accumulates quantity.
func (r *Repository) Purchase(ctx context.Context, userID int, item *Item, quantity int) (newQuantity int, err error) {
	if item == nil {
		return 0, fmt.Errorf("item cannot be nil")
	}

	totalCost := item.Price * quantity

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
		UPDATE users
		SET coins = coins - $1, updated_at = $2
		WHERE id = $3 AND coins >= $1
	`, totalCost, time.Now(), userID)
	if err != nil {
		return 0, fmt.Errorf("failed to deduct coins: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return 0, ErrInsufficientCoins
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO player_inventory (user_id, item_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, item_id)
		DO UPDATE SET quantity = player_inventory.quantity + $3, updated_at = $4
		RETURNING quantity
	`, userID, item.ID, quantity, time.Now()).Scan(&newQuantity)
	if err != nil {
		return 0, fmt.Errorf("failed to add item to inventory: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newQuantity, nil
}

// ActivateBooster consumes one booster from the inventory and inserts the
// corresponding active boost in a single transaction. The guarded inventory
// decrement (`quantity > 0`) means a race can never over-consume.
func (r *Repository) ActivateBooster(ctx context.Context, userID int, item *Item, effect BoosterEffect) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
		UPDATE player_inventory
		SET quantity = quantity - 1, updated_at = $3
		WHERE user_id = $1 AND item_id = $2 AND quantity > 0
	`, userID, item.ID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to consume booster: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrInsufficientItem
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO active_boosts (user_id, item_id, stat, bonus, battles_remaining)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, item.ID, effect.Stat, effect.Bonus, effect.Battles)
	if err != nil {
		return fmt.Errorf("failed to insert active boost: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// ActiveBoosts returns the player's currently applied deck buffs, joined
// with the catalog for display names.
func (r *Repository) ActiveBoosts(ctx context.Context, userID int) ([]ActiveBoost, error) {
	rows, err := r.db.Query(ctx, `
		SELECT ab.id, ab.item_id, i.name, ab.stat, ab.bonus, ab.battles_remaining, ab.created_at
		FROM active_boosts ab
		JOIN items i ON i.id = ab.item_id
		WHERE ab.user_id = $1
		ORDER BY ab.created_at ASC, ab.id ASC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list active boosts: %w", err)
	}
	defer rows.Close()

	var result []ActiveBoost
	for rows.Next() {
		var b ActiveBoost
		if err := rows.Scan(&b.ID, &b.ItemID, &b.ItemName, &b.Stat, &b.Bonus, &b.BattlesRemaining, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan active boost: %w", err)
		}
		result = append(result, b)
	}
	return result, rows.Err()
}

// ConsumeBattleBoosts ticks every active boost of the player down by one
// battle and removes exhausted ones, in one transaction. Called when a battle
// ends so boosts always last exactly their configured number of battles.
func (r *Repository) ConsumeBattleBoosts(ctx context.Context, userID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE active_boosts
		SET battles_remaining = battles_remaining - 1
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return fmt.Errorf("failed to tick active boosts: %w", err)
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM active_boosts WHERE user_id = $1 AND battles_remaining <= 0
	`, userID)
	if err != nil {
		return fmt.Errorf("failed to remove exhausted boosts: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
