# Task 4 Implementation Guide: Enhanced Pokemon Card System with Leveling

## Quick Start

This guide shows how to use the newly implemented Pokemon card system with leveling mechanics.

## 1. Card Model Updates

### Card Structure (pokemon/types.go)
```go
type Card struct {
    Name        string
    HP          int
    HPMax       int
    Stamina     int
    Defense     int
    Attack      int
    Speed       int
    Moves       []Move
    Types       []string
    Sprite      string
    Level       int        // NEW: Current level (1-50)
    XP          int        // NEW: Current experience points
    IsLegendary bool       // NEW: Legendary flag
    IsMythical  bool       // NEW: Mythical flag
}
```

### Get Current Stats Based on Level
```go
card := &Card{Level: 15, HPMax: 100, Attack: 50, Defense: 40, Speed: 60}
stats := card.GetCurrentStats()
// stats.HP = 100 * (1 + 14*0.03) = 142
// stats.Attack = 50 * (1 + 14*0.02) = 64
// stats.Defense = 40 * (1 + 14*0.02) = 51
// stats.Speed = 60 * (1 + 14*0.01) = 68
// stats.Stamina = 68 * 2 = 136
```

## 2. Legendary/Mythical Identification

### Check if Pokemon is Legendary or Mythical
```go
import "pokemon-cli/pokemon"

isLegendary, isMythical := pokemon.IsLegendaryOrMythical("mewtwo")
// isLegendary = true, isMythical = false

isLegendary, isMythical = pokemon.IsLegendaryOrMythical("mew")
// isLegendary = false, isMythical = true

isLegendary, isMythical = pokemon.IsLegendaryOrMythical("pikachu")
// isLegendary = false, isMythical = false
```

## 3. Starter Deck Generation

### Automatic on User Registration
When a user registers, 5 random non-legendary Pokemon are automatically generated:

```go
// In auth/handlers.go Register() function
user, err := h.userRepo.Create(ctx, req.Username, req.Email, passwordHash)
if err != nil {
    return err
}

// Automatically generates 5 starter cards
_, err = h.cardService.GenerateStarterDeck(ctx, user.ID)
```

### Manual Starter Deck Generation
```go
cardService := database.NewCardService(cardRepo)
starterCards, err := cardService.GenerateStarterDeck(ctx, userID)
if err != nil {
    log.Fatal(err)
}

// starterCards contains 5 unique non-legendary Pokemon at Level 1
for _, card := range starterCards {
    fmt.Printf("%s (Level %d, XP: %d)\n", card.PokemonName, card.Level, card.XP)
}
```

## 4. XP and Leveling System

### Award XP to a Card
```go
cardService := database.NewCardService(cardRepo)

// Award 20 XP for 1v1 battle win
updatedCard, err := cardService.AddXP(ctx, cardID, 20)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Card: %s\n", updatedCard.PokemonName)
fmt.Printf("Level: %d\n", updatedCard.Level)
fmt.Printf("XP: %d / %d\n", updatedCard.XP, 100*updatedCard.Level)
```

### Leveling Mechanics
- **XP Required for Next Level:** `100 * current_level`
  - Level 1 → 2: 100 XP
  - Level 2 → 3: 200 XP
  - Level 10 → 11: 1000 XP
- **Level Cap:** 50
- **Stat Increases per Level:**
  - HP: +3%
  - Attack: +2%
  - Defense: +2%
  - Speed: +1%
  - Stamina: +1% (derived from Speed)

### Example Level-Up Flow
```go
card := &database.PlayerCard{
    Level: 1,
    XP: 0,
    BaseHP: 100,
    BaseAttack: 50,
}

// Award 250 XP (enough for 2 level-ups)
updatedCard, _ := cardService.AddXP(ctx, card.ID, 250)
// Level 1: needs 100 XP → Level 2 (150 XP remaining)
// Level 2: needs 200 XP → stays at Level 2 (150 XP)
// Result: Level 2, XP 150/200
```

## 5. Card Management API

### Get All User Cards
```bash
curl -X GET http://localhost:3000/api/cards \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:
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
      "is_legendary": false,
      "is_mythical": false,
      "in_deck": true,
      "deck_position": 1
    }
  ]
}
```

### Get Current Deck
```bash
curl -X GET http://localhost:3000/api/cards/deck \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update Deck
```bash
curl -X PUT http://localhost:3000/api/cards/deck \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "card_ids": [1, 2, 3, 4, 5]
  }'
```

## 6. Integration Example

### Complete Server Setup
```go
package main

import (
    "log"
    "pokemon-cli/auth"
    "pokemon-cli/cards"
    "pokemon-cli/database"
    "github.com/gofiber/fiber/v2"
)

func main() {
    // Initialize database
    if err := database.InitDB(); err != nil {
        log.Fatal(err)
    }
    defer database.CloseDB()

    // Initialize services
    authService := auth.NewService()
    jwtService, _ := auth.NewJWTService()
    
    userRepo := database.NewUserRepository(database.GetDB())
    cardRepo := database.NewCardRepository(database.GetDB())
    cardService := database.NewCardService(cardRepo)

    // Initialize handlers
    authHandler := auth.NewHandler(authService, jwtService, userRepo, cardService)
    cardHandler := cards.NewHandler(cardService)

    // Create app
    app := fiber.New()

    // Register routes
    auth.RegisterRoutes(app, authHandler, jwtService)
    cards.RegisterRoutes(app, cardHandler, jwtService)

    // Start server
    log.Fatal(app.Listen(":3000"))
}
```

## 7. Battle Rewards Integration (Future)

### Award XP After Battle
```go
func awardBattleRewards(ctx context.Context, cardService *database.CardService, 
                        participatingCardIDs []int, battleMode string) error {
    xpAmount := 20 // 1v1
    if battleMode == "5v5" {
        xpAmount = 15 // per card in 5v5
    }

    for _, cardID := range participatingCardIDs {
        updatedCard, err := cardService.AddXP(ctx, cardID, xpAmount)
        if err != nil {
            return err
        }
        
        // Notify if leveled up
        if updatedCard.Level > previousLevel {
            log.Printf("%s leveled up to %d!", 
                      updatedCard.PokemonName, updatedCard.Level)
        }
    }
    return nil
}
```

## 8. Frontend Integration

### Display Card with Level
```javascript
function CardDisplay({ card }) {
  const stats = calculateCurrentStats(card);
  const xpProgress = (card.xp / (100 * card.level)) * 100;
  
  return (
    <div className="pokemon-card">
      <h3>{card.pokemon_name}</h3>
      <div className="level">Level {card.level}</div>
      <div className="xp-bar">
        <div className="xp-fill" style={{ width: `${xpProgress}%` }} />
        <span>{card.xp} / {100 * card.level} XP</span>
      </div>
      <div className="stats">
        <div>HP: {stats.hp}</div>
        <div>Attack: {stats.attack}</div>
        <div>Defense: {stats.defense}</div>
        <div>Speed: {stats.speed}</div>
      </div>
      {card.is_legendary && <span className="badge legendary">Legendary</span>}
      {card.is_mythical && <span className="badge mythical">Mythical</span>}
    </div>
  );
}

function calculateCurrentStats(card) {
  const levelMultiplier = card.level - 1;
  return {
    hp: Math.floor(card.base_hp * (1 + levelMultiplier * 0.03)),
    attack: Math.floor(card.base_attack * (1 + levelMultiplier * 0.02)),
    defense: Math.floor(card.base_defense * (1 + levelMultiplier * 0.02)),
    speed: Math.floor(card.base_speed * (1 + levelMultiplier * 0.01)),
  };
}
```

## 9. Testing Checklist

- [ ] Register new user and verify 5 starter cards created
- [ ] Verify no legendary/mythical in starter deck
- [ ] Verify no duplicates in starter deck
- [ ] Award XP and verify level-up at correct threshold
- [ ] Verify stats increase by correct percentages
- [ ] Verify level cap at 50
- [ ] Test GET /api/cards endpoint
- [ ] Test GET /api/cards/deck endpoint
- [ ] Test PUT /api/cards/deck with valid 5 cards
- [ ] Test PUT /api/cards/deck with invalid count (should fail)
- [ ] Verify authentication required on all card endpoints

## 10. Environment Variables

No new environment variables required for this task. Existing database configuration is sufficient.

## Troubleshooting

### Issue: Starter deck not generated on registration
**Solution:** Ensure cardService is passed to auth.NewHandler() and database is initialized

### Issue: Level-up not occurring
**Solution:** Check XP threshold calculation (100 * current_level)

### Issue: Stats not increasing
**Solution:** Use GetCurrentStats() method, not base stats directly

### Issue: Legendary Pokemon in starter deck
**Solution:** Verify IsLegendaryOrMythical() is called in GenerateStarterDeck()
