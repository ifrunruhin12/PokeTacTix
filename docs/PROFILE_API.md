# Profile and Statistics API Documentation

This document describes the profile and statistics API endpoints for the PokeTacTix web application.

## Authentication

All profile endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### Get Player Statistics

Retrieves comprehensive statistics for the authenticated player.

**Endpoint:** `GET /api/profile/stats`

**Response:**
```json
{
  "user_id": 1,
  "total_battles_1v1": 25,
  "wins_1v1": 18,
  "losses_1v1": 7,
  "total_battles_5v5": 10,
  "wins_5v5": 6,
  "losses_5v5": 4,
  "total_coins_earned": 1250,
  "highest_level": 23,
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Status Codes:**
- `200 OK` - Statistics retrieved successfully
- `401 Unauthorized` - Missing or invalid authentication token
- `500 Internal Server Error` - Server error

---

### Get Battle History

Retrieves the battle history for the authenticated player.

**Endpoint:** `GET /api/profile/history`

**Query Parameters:**
- `limit` (optional, default: 20, max: 100) - Number of battles to retrieve

**Example:** `GET /api/profile/history?limit=10`

**Response:**
```json
{
  "history": [
    {
      "id": 123,
      "user_id": 1,
      "mode": "5v5",
      "result": "win",
      "coins_earned": 150,
      "duration": 180,
      "created_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": 122,
      "user_id": 1,
      "mode": "1v1",
      "result": "loss",
      "coins_earned": 10,
      "duration": 45,
      "created_at": "2024-01-15T09:15:00Z"
    }
  ],
  "count": 2
}
```

**Battle Result Values:**
- `win` - Player won the battle
- `loss` - Player lost the battle
- `draw` - Battle ended in a draw

**Status Codes:**
- `200 OK` - History retrieved successfully
- `401 Unauthorized` - Missing or invalid authentication token
- `500 Internal Server Error` - Server error

---

### Get Achievements

Retrieves all achievements with their unlock status for the authenticated player.

**Endpoint:** `GET /api/profile/achievements`

**Response:**
```json
{
  "achievements": [
    {
      "id": 1,
      "name": "First Victory",
      "description": "Win your first battle",
      "icon": "🏆",
      "requirement_type": "total_wins",
      "requirement_value": 1,
      "unlocked": true,
      "unlocked_at": "2024-01-10T14:20:00Z"
    },
    {
      "id": 2,
      "name": "Veteran Trainer",
      "description": "Win 10 battles",
      "icon": "⭐",
      "requirement_type": "total_wins",
      "requirement_value": 10,
      "unlocked": false,
      "unlocked_at": null
    }
  ],
  "total": 10,
  "unlocked": 3,
  "locked": 7
}
```

**Achievement Requirement Types:**
- `total_wins` - Total number of wins across all battle modes
- `total_battles` - Total number of battles completed
- `wins_5v5` - Number of 5v5 battle wins
- `total_coins` - Total coins earned
- `max_level` - Highest Pokemon level achieved
- `legendary_owned` - Owns at least one legendary Pokemon
- `mythical_owned` - Owns at least one mythical Pokemon

**Status Codes:**
- `200 OK` - Achievements retrieved successfully
- `401 Unauthorized` - Missing or invalid authentication token
- `500 Internal Server Error` - Server error

---

### Check and Unlock Achievements

Checks the player's progress and unlocks any newly earned achievements.

**Endpoint:** `POST /api/profile/achievements/check`

**Response:**
```json
{
  "newly_unlocked": [
    {
      "id": 2,
      "name": "Veteran Trainer",
      "description": "Win 10 battles",
      "icon": "⭐",
      "requirement_type": "total_wins",
      "requirement_value": 10,
      "unlocked": true,
      "unlocked_at": "2024-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

**Status Codes:**
- `200 OK` - Achievements checked successfully
- `401 Unauthorized` - Missing or invalid authentication token
- `500 Internal Server Error` - Server error

---

## Default Achievements

The system includes the following default achievements:

1. **First Victory** 🏆 - Win your first battle
2. **Veteran Trainer** ⭐ - Win 10 battles
3. **Elite Trainer** 💫 - Win 50 battles
4. **Champion** 👑 - Win 100 battles
5. **Legendary Collector** 🌟 - Obtain a legendary Pokemon
6. **Mythical Master** ✨ - Obtain a mythical Pokemon
7. **Max Level** 📈 - Get a Pokemon to level 50
8. **Coin Hoarder** 💰 - Accumulate 5000 coins
9. **Battle Enthusiast** ⚔️ - Complete 25 battles
10. **5v5 Specialist** 🎯 - Win 20 5v5 battles

---

## Integration with Battle System

The statistics system automatically tracks battle outcomes when battles are completed. The battle system should call the stats service to record battles:

```go
// After battle completion
err := statsService.RecordBattle(ctx, userID, mode, result, coinsEarned, duration)
if err != nil {
    // Handle error
}

// Check for newly unlocked achievements
newAchievements, err := statsService.CheckAndUnlockAchievements(ctx, userID)
if err != nil {
    // Handle error
}
```

---

## Error Responses

All endpoints follow a consistent error format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": "Additional error details (optional)"
  }
}
```

**Common Error Codes:**
- `UNAUTHORIZED` - Missing or invalid authentication
- `INTERNAL_ERROR` - Server-side error
- `DATABASE_ERROR` - Database operation failed

---

## Notes

- Battle history is ordered by creation date (most recent first)
- Statistics are updated in real-time after each battle
- Achievement checks can be triggered manually or automatically after battles
- The highest level is automatically updated when Pokemon level up
- All timestamps are in ISO 8601 format (UTC)
