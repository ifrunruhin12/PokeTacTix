# Database Layer

This directory contains the database connection, models, and repository layer for PokeTacTix.

## Structure

- `connection.go` - Database connection pool management with pgx
- `models.go` - Data models and structs
- `user_repository.go` - User CRUD operations
- `card_repository.go` - Player card management
- `battle_repository.go` - Battle history tracking
- `stats_repository.go` - Player statistics and achievements

## Usage

### Initialize Database Connection

```go
import "your-module/database"

func main() {
    // Initialize database connection
    if err := database.InitDB(); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer database.CloseDB()
    
    // Get database pool
    db := database.GetDB()
}
```

### Environment Variables

Required environment variables:

```bash
DATABASE_URL=postgresql://user:password@host:5432/dbname
DB_MAX_CONNECTIONS=20        # Optional, default: 20
DB_MIN_CONNECTIONS=2         # Optional, default: 2
DB_IDLE_TIMEOUT=300          # Optional, default: 300 seconds
DB_MAX_LIFETIME=1800         # Optional, default: 1800 seconds
```

### Using Repositories

#### User Repository

```go
userRepo := database.NewUserRepository(db)

// Create user
user, err := userRepo.Create(ctx, "username", "email@example.com", "hashedPassword")

// Get user by ID
user, err := userRepo.GetByID(ctx, 1)

// Get user by username
user, err := userRepo.GetByUsername(ctx, "username")

// Update coins
err := userRepo.AddCoins(ctx, userID, 100)

// Check if username exists
exists, err := userRepo.UsernameExists(ctx, "username")
```

#### Card Repository

```go
cardRepo := database.NewCardRepository(db)

// Create card
card := &database.PlayerCard{
    UserID:      1,
    PokemonName: "Pikachu",
    Level:       1,
    XP:          0,
    BaseHP:      35,
    BaseAttack:  55,
    BaseDefense: 40,
    BaseSpeed:   90,
    Types:       json.RawMessage(`["electric"]`),
    Moves:       json.RawMessage(`[...]`),
    Sprite:      "https://...",
}
newCard, err := cardRepo.Create(ctx, card)

// Get user's cards
cards, err := cardRepo.GetUserCards(ctx, userID)

// Get user's deck
deck, err := cardRepo.GetUserDeck(ctx, userID)

// Update deck
cardIDs := []int{1, 2, 3, 4, 5}
err := cardRepo.UpdateDeck(ctx, userID, cardIDs)

// Add XP and handle level-ups
card, err := cardRepo.AddXP(ctx, cardID, 50)

// Get current stats based on level
stats := card.GetCurrentStats()
```

#### Battle Repository

```go
battleRepo := database.NewBattleRepository(db)

// Create battle history
battle := &database.BattleHistory{
    UserID:      1,
    Mode:        "5v5",
    Result:      "win",
    CoinsEarned: 150,
    Duration:    &duration,
}
newBattle, err := battleRepo.Create(ctx, battle)

// Get user's battle history
battles, err := battleRepo.GetUserHistory(ctx, userID, 20)

// Get history by mode
battles, err := battleRepo.GetUserHistoryByMode(ctx, userID, "1v1", 10)
```

#### Stats Repository

```go
statsRepo := database.NewStatsRepository(db)

// Get or create stats
stats, err := statsRepo.GetOrCreate(ctx, userID)

// Record battle
err := statsRepo.RecordBattle(ctx, userID, "5v5", "win", 150)

// Update highest level
err := statsRepo.UpdateHighestLevel(ctx, userID, 25)

// Get achievements
achievements, err := statsRepo.GetAchievements(ctx, userID)

// Check and unlock achievements
newAchievements, err := statsRepo.CheckAndUnlockAchievements(ctx, userID)
```

## Database Schema

### Tables

1. **users** - User accounts
2. **player_cards** - Pokemon cards owned by players
3. **battle_history** - Battle records
4. **player_stats** - Player statistics
5. **achievements** - Achievement definitions
6. **user_achievements** - Unlocked achievements

### Relationships

- `player_cards.user_id` → `users.id` (CASCADE DELETE)
- `battle_history.user_id` → `users.id` (CASCADE DELETE)
- `player_stats.user_id` → `users.id` (CASCADE DELETE)
- `user_achievements.user_id` → `users.id` (CASCADE DELETE)
- `user_achievements.achievement_id` → `achievements.id` (CASCADE DELETE)

### Indexes

- `users`: username, email
- `player_cards`: user_id, (user_id, in_deck), (user_id, deck_position)
- `battle_history`: user_id, created_at, (user_id, created_at)

## Migrations

Migration files are located in the `../migrations` directory. See the main README for migration instructions.

## Connection Pooling

The connection pool is configured with:
- Max connections: 20 (configurable)
- Min connections: 2 (configurable)
- Idle timeout: 300 seconds (configurable)
- Max lifetime: 1800 seconds (configurable)

## Error Handling

All repository methods return errors that should be handled appropriately:

```go
user, err := userRepo.GetByID(ctx, id)
if err != nil {
    if err.Error() == "user not found" {
        // Handle not found
    } else {
        // Handle other errors
    }
}
```

## Testing

To test database operations:

```bash
# Run all tests
make test

# Run with coverage
make test-coverage
```

## Production Deployment

For production deployment with Neon:

1. Create a Neon PostgreSQL database
2. Set the `DATABASE_URL` environment variable
3. Run migrations using the Makefile
4. Ensure SSL/TLS is enabled in production

See the deployment documentation for detailed instructions.
