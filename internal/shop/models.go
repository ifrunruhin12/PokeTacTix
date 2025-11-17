package shop

import (
	"pokemon-cli/internal/pokemon"
	"time"
)

// ShopItem represents a Pokemon card available for purchase in the shop
type ShopItem struct {
	PokemonName string   `json:"pokemon_name"`
	Price       int      `json:"price"`
	Rarity      string   `json:"rarity"` // "common", "uncommon", "rare", "legendary", "mythical"
	IsLegendary bool     `json:"is_legendary"`
	IsMythical  bool     `json:"is_mythical"`
	Sprite      string   `json:"sprite"`
	Types       []string `json:"types"`
	InStock     bool     `json:"in_stock"`
	// Full card details for display
	BaseHP      int            `json:"base_hp"`
	BaseAttack  int            `json:"base_attack"`
	BaseDefense int            `json:"base_defense"`
	BaseSpeed   int            `json:"base_speed"`
	Moves       []pokemon.Move `json:"moves"`
}

// ShopInventory represents the current shop state
type ShopInventory struct {
	Items           []ShopItem `json:"items"`
	DiscountActive  bool       `json:"discount_active"`
	DiscountPercent int        `json:"discount_percent"`
	RefreshTime     time.Time  `json:"refresh_time"`
}

// PurchaseRequest represents a purchase request
type PurchaseRequest struct {
	PokemonName string `json:"pokemon_name"`
}

// PurchaseResponse represents a purchase response
type PurchaseResponse struct {
	Card           any `json:"card"`
	RemainingCoins int `json:"remaining_coins"`
}
