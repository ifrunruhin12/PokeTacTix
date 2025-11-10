# Cards Package

This package handles Pokemon card management functionality including card collection retrieval, deck management, and card operations.

## API Endpoints

All endpoints require JWT authentication via the `Authorization: Bearer <token>` header.

### GET /api/cards

Retrieves all Pokemon cards owned by the authenticated user.

**Response:**
```json
{
  "cards": [
    {
      "id": 1,
      "user_id": 123,
      "pokemon_name": "pikachu",
      "level": 15,
      "xp": 450,
      "base_hp": 35,
      "base_attack": 55,
      "base_defense": 40,
      "base_speed": 90,
      "types": ["electric"],
      "moves": [...],
      "sprite": "https://...",
      "is_legendary": false,
      "is_mythical": false,
      "in_deck": true,
      "deck_position": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### GET /api/cards/deck

Retrieves the user's current battle deck (5 cards).

**Response:**
```json
{
  "deck": [
    {
      "id": 1,
      "pokemon_name": "pikachu",
      "level": 15,
      ...
    },
    ...
  ]
}
```

### PUT /api/cards/deck

Updates the user's battle deck. Must contain exactly 5 card IDs.

**Request:**
```json
{
  "card_ids": [1, 2, 3, 4, 5]
}
```

**Response:**
```json
{
  "message": "Deck updated successfully",
  "deck": [...]
}
```

**Errors:**
- `INVALID_DECK` (400): Deck doesn't have exactly 5 cards
- `CARD_NOT_FOUND` (404): One or more card IDs don't exist or don't belong to user

### GET /api/cards/:id

Retrieves a specific card by ID. The card must belong to the authenticated user.

**Response:**
```json
{
  "card": {
    "id": 1,
    "pokemon_name": "pikachu",
    ...
  }
}
```

**Errors:**
- `CARD_NOT_FOUND` (404): Card doesn't exist
- `FORBIDDEN` (403): Card doesn't belong to user

## Usage Example

```go
// In server/main.go
import (
    "pokemon-cli/cards"
    "pokemon-cli/database"
)

// Initialize services
cardRepo := database.NewCardRepository(database.GetDB())
cardService := database.NewCardService(cardRepo)
cardHandler := cards.NewHandler(cardService)

// Register routes
cards.RegisterRoutes(app, cardHandler, jwtService)
```

## Requirements Fulfilled

- **Requirement 12.1**: GET /api/cards endpoint to retrieve user's card collection
- **Requirement 12.2**: GET /api/cards/deck endpoint to get current deck
- **Requirement 12.3**: PUT /api/cards/deck endpoint to update deck (must have exactly 5 cards)
- All endpoints include authentication middleware as required
