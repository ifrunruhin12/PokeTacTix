package cards

// UpdateDeckRequest represents the request body for updating a deck
type UpdateDeckRequest struct {
	CardIDs []int `json:"card_ids"`
}
