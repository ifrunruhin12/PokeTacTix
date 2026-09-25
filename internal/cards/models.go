package cards

import (
	"pokemon-cli/internal/database"
)

// UpdateDeckRequest represents the request body for updating a deck
type UpdateDeckRequest struct {
	CardIDs []int `json:"card_ids"`
}

// EvolveRequest represents the request body for evolving a card with an item
type EvolveRequest struct {
	ItemID string `json:"item_id"`
}

// EvolutionOptionResponse describes one way a card can evolve, enriched with
// the player's inventory so the deck UI can show what is actionable now.
type EvolutionOptionResponse struct {
	Method           string `json:"method"` // "level" or "item"
	TargetName       string `json:"target_name"`
	TargetSprite     string `json:"target_sprite"`
	MinLevel         int    `json:"min_level,omitempty"`          // method "level"
	RequiredItem     string `json:"required_item,omitempty"`      // method "item": item id (slug)
	RequiredItemName string `json:"required_item_name,omitempty"` // method "item": display name
	OwnedQuantity    int    `json:"owned_quantity,omitempty"`     // method "item"
	Eligible         bool   `json:"eligible"`
}

// EvolutionInfo is the deck-facing view of a card's evolution possibilities
type EvolutionInfo struct {
	CardID       int                       `json:"card_id"`
	PokemonName  string                    `json:"pokemon_name"`
	Level        int                       `json:"level"`
	Options      []EvolutionOptionResponse `json:"options"`
	CanEvolveNow bool                      `json:"can_evolve_now"` // an item evolution is actionable right now
}

// EvolutionSummary tells the deck UI whether a card's species has any
// evolution path at all, so the Evolve action can be hidden for final forms.
type EvolutionSummary struct {
	CardID       int  `json:"card_id"`
	HasEvolution bool `json:"has_evolution"`
}

// EvolutionResult reports a completed item-based evolution
type EvolutionResult struct {
	CardID                int                  `json:"card_id"`
	EvolvedFrom           string               `json:"evolved_from"`
	EvolvedInto           string               `json:"evolved_into"`
	ItemConsumed          string               `json:"item_consumed"`
	RemainingItemQuantity int                  `json:"remaining_item_quantity"`
	Card                  *database.PlayerCard `json:"card"`
}
