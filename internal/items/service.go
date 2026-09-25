package items

import (
	"context"
	"encoding/json"
	"fmt"
)

// Quantity bounds for a single shop item purchase.
const (
	MinPurchaseQuantity = 1
	MaxPurchaseQuantity = 99
)

// catalogRepository is the persistence surface the service depends on.
// *Repository implements it; tests provide their own.
type catalogRepository interface {
	ListItems(ctx context.Context) ([]Item, error)
	GetItem(ctx context.Context, itemID string) (*Item, error)
	GetInventory(ctx context.Context, userID int) ([]InventoryEntry, error)
	GetQuantity(ctx context.Context, userID int, itemID string) (int, error)
	Purchase(ctx context.Context, userID int, item *Item, quantity int) (int, error)
	ActivateBooster(ctx context.Context, userID int, item *Item, effect BoosterEffect) (*ActiveBoost, error)
	ActiveBoosts(ctx context.Context, userID int) ([]ActiveBoost, error)
	ConsumeBattleBoosts(ctx context.Context, userID int) error
}

// Service handles item catalog and inventory business logic
type Service struct {
	repository catalogRepository
}

// NewService creates a new items service
func NewService(repository catalogRepository) *Service {
	return &Service{repository: repository}
}

// ListItems returns the full item catalog
func (s *Service) ListItems(ctx context.Context) ([]Item, error) {
	return s.repository.ListItems(ctx)
}

// GetItem returns a catalog item by id
func (s *Service) GetItem(ctx context.Context, itemID string) (*Item, error) {
	return s.repository.GetItem(ctx, itemID)
}

// GetInventory returns the player's inventory
func (s *Service) GetInventory(ctx context.Context, userID int) ([]InventoryEntry, error) {
	return s.repository.GetInventory(ctx, userID)
}

// GetQuantity returns the player's quantity of an item
func (s *Service) GetQuantity(ctx context.Context, userID int, itemID string) (int, error) {
	return s.repository.GetQuantity(ctx, userID, itemID)
}

// Purchase buys the given quantity of an item. The item's effect is NOT
// applied here — items sit in the inventory until they are used.
func (s *Service) Purchase(ctx context.Context, userID int, itemID string, quantity int) (*Item, int, error) {
	if quantity == 0 {
		quantity = 1 // omitted quantity in the request body defaults to one
	}
	if quantity < MinPurchaseQuantity || quantity > MaxPurchaseQuantity {
		return nil, 0, fmt.Errorf("invalid quantity: must be between %d and %d", MinPurchaseQuantity, MaxPurchaseQuantity)
	}

	item, err := s.repository.GetItem(ctx, itemID)
	if err != nil {
		return nil, 0, err
	}

	newQuantity, err := s.repository.Purchase(ctx, userID, item, quantity)
	if err != nil {
		return nil, 0, err
	}

	return item, newQuantity, nil
}

func (s *Service) PurchaseWithKey(ctx context.Context, userID int, itemID string, quantity int, key string) (*Item, int, int, error) {
	if quantity == 0 {
		quantity = 1
	}
	if quantity < MinPurchaseQuantity || quantity > MaxPurchaseQuantity {
		return nil, 0, 0, fmt.Errorf("invalid quantity: must be between %d and %d", MinPurchaseQuantity, MaxPurchaseQuantity)
	}
	item, err := s.repository.GetItem(ctx, itemID)
	if err != nil {
		return nil, 0, 0, err
	}
	keyed, ok := s.repository.(interface {
		PurchaseIdempotent(context.Context, int, *Item, int, string) (int, int, error)
	})
	if !ok {
		return nil, 0, 0, fmt.Errorf("idempotent purchases unavailable")
	}
	newQuantity, coins, err := keyed.PurchaseIdempotent(ctx, userID, item, quantity, key)
	return item, newQuantity, coins, err
}

// ActiveBoosts returns the player's currently applied deck buffs.
func (s *Service) ActiveBoosts(ctx context.Context, userID int) ([]ActiveBoost, error) {
	return s.repository.ActiveBoosts(ctx, userID)
}

// ConsumeBattleBoosts is a standalone legacy tick. New battle sessions reserve
// boost duration atomically with the session in the battle repository.
func (s *Service) ConsumeBattleBoosts(ctx context.Context, userID int) error {
	return s.repository.ConsumeBattleBoosts(ctx, userID)
}

// UseItem activates a booster from the player's inventory. Evolution items
// are rejected here — they are used through the deck's evolution flow, which
// knows the evolution rules. The inventory decrement and the boost insertion
// happen atomically in the repository.
func (s *Service) UseItem(ctx context.Context, userID int, itemID string) (*ActiveBoost, error) {
	item, err := s.repository.GetItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if item.ItemType != TypeBooster {
		return nil, fmt.Errorf("%w: %s is not a booster", ErrNotBooster, item.Name)
	}

	effect, err := ParseBoosterEffect(item)
	if err != nil {
		return nil, err
	}

	// Pre-check for a clear error; the activation transaction re-checks with
	// a guarded decrement so a race cannot over-consume.
	owned, err := s.repository.GetQuantity(ctx, userID, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory for item %s: %w", itemID, err)
	}
	if owned < 1 {
		return nil, fmt.Errorf("%w: you own 0 of %s", ErrInsufficientItem, item.Name)
	}

	return s.repository.ActivateBooster(ctx, userID, item, effect)
}

// ParseBoosterEffect decodes and validates a booster item's effect payload.
func ParseBoosterEffect(item *Item) (BoosterEffect, error) {
	if len(item.Effect) == 0 {
		return BoosterEffect{}, fmt.Errorf("booster %s has no effect configured", item.ID)
	}
	var effect BoosterEffect
	if err := json.Unmarshal(item.Effect, &effect); err != nil {
		return BoosterEffect{}, fmt.Errorf("invalid effect on booster %s: %w", item.ID, err)
	}
	if effect.Stat != StatAttack && effect.Stat != StatHP {
		return BoosterEffect{}, fmt.Errorf("booster %s has unsupported stat %q", item.ID, effect.Stat)
	}
	if effect.Bonus <= 0 || effect.Battles <= 0 {
		return BoosterEffect{}, fmt.Errorf("booster %s must have positive bonus and battles", item.ID)
	}
	return effect, nil
}
