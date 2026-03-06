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

// GameTokenItem represents game tokens available for purchase
type GameTokenItem struct {
	ItemType           string `json:"item_type"`            // "game_token"
	Price              int    `json:"price"`                // 100 coins per token
	AvailableQuantity  int    `json:"available_quantity"`   // 10 - tokens_purchased_today
	MaxDailyPurchase   int    `json:"max_daily_purchase"`   // 10
	Description        string `json:"description"`
}

// ShopInventory represents the current shop state
type ShopInventory struct {
	Items           []ShopItem     `json:"items"`
	GameTokens      *GameTokenItem `json:"game_tokens,omitempty"` // Token purchase option
	DiscountActive  bool           `json:"discount_active"`
	DiscountPercent int            `json:"discount_percent"`
	RefreshTime     time.Time      `json:"refresh_time"`
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

// TokenPurchaseRequest represents a token purchase request
type TokenPurchaseRequest struct {
	Quantity int `json:"quantity"` // 1-10
}

// TokenPurchaseResponse represents a token purchase response
type TokenPurchaseResponse struct {
	Success              bool `json:"success"`
	TokensAdded          int  `json:"tokens_added"`
	NewTokenBalance      int  `json:"new_token_balance"`
	CoinsSpent           int  `json:"coins_spent"`
	RemainingCoins       int  `json:"remaining_coins"`
	TokensPurchasedToday int  `json:"tokens_purchased_today"`
	DailyLimitRemaining  int  `json:"daily_limit_remaining"`
}
