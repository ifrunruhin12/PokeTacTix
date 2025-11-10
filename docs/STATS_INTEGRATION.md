# Statistics System Integration Guide

This document describes how to integrate the statistics tracking system with the battle system.

## Overview

The statistics system has been implemented in `internal/stats/` and provides the following functionality:

1. **Battle History Tracking** - Records all battle outcomes
2. **Player Statistics** - Tracks wins, losses, coins earned, and highest level
3. **Achievement System** - Defines and tracks player achievements
4. **Profile API** - Provides endpoints to retrieve stats, history, and achievements

## Integration Points

### 1. Battle Completion

When a battle ends (when `battleState.BattleOver` becomes `true`), the battle system should:

1. Calculate rewards using `CalculateRewards()`
2. Apply rewards to the database using `ApplyRewards()`
3. Record the battle in history using the stats service
4. Check for newly unlocked achievements

**Example Integration in `internal/battle/handlers.go`:**

```go
// After ProcessMove, check if battle is over
if battleState.BattleOver {
    // Get database connection
    db := c.Locals("db").(*pgxpool.Pool)
    
    // Calculate rewards
    rewards := CalculateRewards(battleState)
    
    // Apply rewards (coins and XP)
    err := ApplyRewards(c.Context(), db, userID, battleState, rewards)
    if err != nil {
        // Log error but don't fail the request
        log.Printf("Failed to apply rewards: %v", err)
    }
    
    // Record battle history (requires stats service)
    // This should be injected into the handler
    duration := int(time.Since(battleState.CreatedAt).Seconds())
    err = statsService.RecordBattle(c.Context(), userID, battleState.Mode, 
        battleState.Winner, rewards.CoinsEarned, duration)
    if err != nil {
        log.Printf("Failed to record battle: %v", err)
    }
    
    // Check for newly unlocked achievements
    newAchievements, err := statsService.CheckAndUnlockAchievements(c.Context(), userID)
    if err != nil {
        log.Printf("Failed to check achievements: %v", err)
    }
    
    // Include achievements in response if any were unlocked
    if len(newAchievements) > 0 {
        response["newly_unlocked_achievements"] = newAchievements
    }
}
```

### 2. Level Up Tracking

When a Pokemon levels up (in `ApplyRewards`), update the highest level stat:

```go
// After level ups are processed
for _, levelUp := range rewards.LevelUps {
    err := statsService.UpdateHighestLevel(ctx, userID, levelUp.NewLevel)
    if err != nil {
        log.Printf("Failed to update highest level: %v", err)
    }
}
```

### 3. Handler Dependency Injection

The battle handler needs access to the stats service. Update the handler constructor:

```go
// Handler handles battle-related HTTP requests
type Handler struct {
    sessions     map[string]*Session
    battleStates map[string]*BattleState
    mu           sync.RWMutex
    statsService *stats.Service  // Add this
}

// NewHandler creates a new battle handler
func NewHandler(statsService *stats.Service) *Handler {
    return &Handler{
        sessions:     make(map[string]*Session),
        battleStates: make(map[string]*BattleState),
        statsService: statsService,
    }
}
```

Then update `cmd/api/main.go`:

```go
battleHandler := battle.NewHandler(statsService)
```

### 4. Achievement Initialization

During application startup, initialize the default achievements:

```go
// In cmd/api/main.go, after database initialization
err = statsService.InitializeAchievements(context.Background())
if err != nil {
    appLogger.Warn("Failed to initialize achievements", "error", err)
}
```

## Battle Result Mapping

The stats system expects battle results as strings:

- `"win"` - Player won the battle
- `"loss"` - Player lost the battle  
- `"draw"` - Battle ended in a draw

Map the battle state winner to result:

```go
result := "loss"
if battleState.Winner == "player" {
    result = "win"
} else if battleState.Winner == "draw" {
    result = "draw"
}
```

## Testing the Integration

After integration, test the following flows:

1. **Win a 1v1 battle** - Verify 50 coins awarded, stats updated, "First Victory" achievement unlocked
2. **Lose a battle** - Verify 10 coins awarded, loss count incremented
3. **Win 10 battles** - Verify "Veteran Trainer" achievement unlocked
4. **Level up a Pokemon to 50** - Verify "Max Level" achievement unlocked and highest level updated
5. **View profile stats** - Verify all stats are accurate
6. **View battle history** - Verify battles are recorded with correct data
7. **View achievements** - Verify achievements show correct unlock status

## API Endpoints

The following endpoints are now available:

- `GET /api/profile/stats` - Get player statistics
- `GET /api/profile/history?limit=20` - Get battle history
- `GET /api/profile/achievements` - Get all achievements with unlock status
- `POST /api/profile/achievements/check` - Manually check and unlock achievements

All endpoints require authentication (JWT token in Authorization header).

## Database Requirements

Ensure the following tables exist:

- `battle_history` - Records all battles
- `player_stats` - Stores aggregated player statistics
- `achievements` - Defines available achievements
- `user_achievements` - Tracks which achievements users have unlocked

See migration files in `internal/database/migrations/` for schema details.

## Future Enhancements

Consider these enhancements for the statistics system:

1. **Leaderboards** - Rank players by wins, coins, or achievements
2. **Battle Replays** - Store detailed battle logs for replay
3. **Daily/Weekly Challenges** - Time-limited achievement goals
4. **Statistics Dashboard** - Visualize stats with charts and graphs
5. **Achievement Notifications** - Real-time notifications when achievements unlock
6. **Battle Analytics** - Track move usage, type effectiveness, etc.
