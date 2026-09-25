package items

import (
	"encoding/json"
	"time"
)

// Item type values. The catalog is data-driven: evolution items are consumed
// by the item-based evolution flow, boosters carry their parameters in
// Effect and buff the player's deck for a number of battles when used.
const (
	TypeEvolution = "evolution"
	TypeBooster   = "booster"
)

// Booster stat values supported by the battle system.
const (
	StatAttack = "attack"
	StatHP     = "hp"
)

// BoosterEffect is the JSON payload stored in items.effect for booster items.
// Magnitudes and durations are data, so balancing is a seed change.
type BoosterEffect struct {
	Stat    string `json:"stat"`    // "attack" or "hp"
	Bonus   int    `json:"bonus"`   // flat bonus applied to every deck Pokemon
	Battles int    `json:"battles"` // number of battles the boost lasts
}

// ActiveBoost is one currently applied deck buff, ticking down per battle.
type ActiveBoost struct {
	ID               int       `json:"id"`
	ItemID           string    `json:"item_id"`
	ItemName         string    `json:"item_name"`
	Stat             string    `json:"stat"`
	Bonus            int       `json:"bonus"`
	BattlesRemaining int       `json:"battles_remaining"`
	CreatedAt        time.Time `json:"created_at"`
}

// Item represents a purchasable game item in the shop catalog
type Item struct {
	ID          string          `json:"id"` // slug matching PokéAPI item names, e.g. "thunder-stone"
	Name        string          `json:"name"`
	Description string          `json:"description"`
	ItemType    string          `json:"item_type"` // "evolution" or "booster"
	Price       int             `json:"price"`     // cost in coins
	Icon        string          `json:"icon"`
	Effect      json.RawMessage `json:"effect,omitempty"` // future: booster parameters (e.g. {"attack": 5, "battles": 3})
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// InventoryEntry represents an item owned by a player with its quantity
type InventoryEntry struct {
	Item
	Quantity int `json:"quantity"`
}

// PurchaseRequest represents a shop item purchase request
type PurchaseRequest struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"` // 1-99, defaults to 1
}

// PurchaseResponse represents a shop item purchase response
type PurchaseResponse struct {
	Success        bool   `json:"success"`
	ItemID         string `json:"item_id"`
	ItemName       string `json:"item_name"`
	QuantityAdded  int    `json:"quantity_added"`
	NewQuantity    int    `json:"new_quantity"`
	CoinsSpent     int    `json:"coins_spent"`
	RemainingCoins int    `json:"remaining_coins"`
}
