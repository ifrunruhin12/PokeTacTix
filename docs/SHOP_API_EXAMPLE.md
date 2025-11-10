# Shop API Response Examples

## GET /api/shop/inventory

Returns the current shop inventory with full Pokemon card details for display.

### Response Structure

```json
{
  "items": [
    {
      "pokemon_name": "pikachu",
      "price": 250,
      "rarity": "uncommon",
      "is_legendary": false,
      "is_mythical": false,
      "sprite": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png",
      "types": ["electric"],
      "in_stock": true,
      "base_hp": 35,
      "base_attack": 55,
      "base_defense": 40,
      "base_speed": 90,
      "moves": [
        {
          "name": "thunderbolt",
          "power": 90,
          "stamina_cost": 30,
          "attack_type": "electric"
        },
        {
          "name": "quick-attack",
          "power": 40,
          "stamina_cost": 13,
          "attack_type": "normal"
        },
        {
          "name": "iron-tail",
          "power": 100,
          "stamina_cost": 33,
          "attack_type": "steel"
        },
        {
          "name": "electro-ball",
          "power": 80,
          "stamina_cost": 26,
          "attack_type": "electric"
        }
      ]
    },
    {
      "pokemon_name": "charizard",
      "price": 500,
      "rarity": "rare",
      "is_legendary": false,
      "is_mythical": false,
      "sprite": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/6.png",
      "types": ["fire", "flying"],
      "in_stock": true,
      "base_hp": 78,
      "base_attack": 84,
      "base_defense": 78,
      "base_speed": 100,
      "moves": [
        {
          "name": "flamethrower",
          "power": 90,
          "stamina_cost": 30,
          "attack_type": "fire"
        },
        {
          "name": "air-slash",
          "power": 75,
          "stamina_cost": 25,
          "attack_type": "flying"
        },
        {
          "name": "dragon-claw",
          "power": 80,
          "stamina_cost": 26,
          "attack_type": "dragon"
        },
        {
          "name": "fire-blast",
          "power": 110,
          "stamina_cost": 36,
          "attack_type": "fire"
        }
      ]
    },
    {
      "pokemon_name": "mewtwo",
      "price": 1500,
      "rarity": "legendary",
      "is_legendary": true,
      "is_mythical": false,
      "sprite": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/150.png",
      "types": ["psychic"],
      "in_stock": true,
      "base_hp": 106,
      "base_attack": 110,
      "base_defense": 90,
      "base_speed": 130,
      "moves": [
        {
          "name": "psychic",
          "power": 90,
          "stamina_cost": 30,
          "attack_type": "psychic"
        },
        {
          "name": "aura-sphere",
          "power": 80,
          "stamina_cost": 26,
          "attack_type": "fighting"
        },
        {
          "name": "shadow-ball",
          "power": 80,
          "stamina_cost": 26,
          "attack_type": "ghost"
        },
        {
          "name": "psystrike",
          "power": 100,
          "stamina_cost": 33,
          "attack_type": "psychic"
        }
      ]
    }
  ],
  "discount_active": true,
  "discount_percent": 40,
  "refresh_time": "2025-11-12T02:09:00Z"
}
```

### With Active Discount

When a discount event is active:
- Legendary Pokemon get 40% off (e.g., 2500 → 1500)
- Mythical Pokemon get 30% off (e.g., 5000 → 3500)
- Other rarities remain at base price

## POST /api/shop/purchase

Purchase a Pokemon card from the shop.

### Request

```json
{
  "pokemon_name": "pikachu"
}
```

### Success Response (200 OK)

```json
{
  "card": {
    "id": 42,
    "user_id": 1,
    "pokemon_name": "pikachu",
    "level": 1,
    "xp": 0,
    "base_hp": 35,
    "base_attack": 55,
    "base_defense": 40,
    "base_speed": 90,
    "types": ["electric"],
    "moves": [
      {
        "name": "thunderbolt",
        "power": 90,
        "stamina_cost": 30,
        "attack_type": "electric"
      },
      {
        "name": "quick-attack",
        "power": 40,
        "stamina_cost": 13,
        "attack_type": "normal"
      }
    ],
    "sprite": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png",
    "is_legendary": false,
    "is_mythical": false,
    "in_deck": false,
    "deck_position": null,
    "created_at": "2025-11-11T02:09:00Z",
    "updated_at": "2025-11-11T02:09:00Z"
  },
  "remaining_coins": 750
}
```

### Error Responses

#### Insufficient Coins (402 Payment Required)

```json
{
  "error": {
    "code": "INSUFFICIENT_COINS",
    "message": "insufficient coins: have 100, need 250"
  }
}
```

#### Pokemon Not Found (404 Not Found)

```json
{
  "error": {
    "code": "ITEM_NOT_FOUND",
    "message": "Pokemon not found in shop inventory"
  }
}
```

#### Rate Limit Exceeded (429 Too Many Requests)

```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many purchase requests. Please try again later.",
    "details": {
      "max": 10,
      "window": "1 minute",
      "retry_after": 60
    }
  }
}
```

## Frontend Card Display

The shop items include all necessary data to display Pokemon cards:

- **pokemon_name**: Display name
- **sprite**: Card image URL
- **types**: Type badges (fire, water, electric, etc.)
- **base_hp, base_attack, base_defense, base_speed**: Base stats for display
- **moves**: All 4 moves with power and stamina cost
- **rarity**: For card styling (common, uncommon, rare, legendary, mythical)
- **price**: Display price (already includes discount if active)
- **is_legendary/is_mythical**: For special card effects/styling

### Suggested Card Layout

```
┌─────────────────────────┐
│  [Sprite Image]         │
│                         │
│  Pikachu                │
│  [⚡ Electric]          │
│                         │
│  HP: 35  ATK: 55        │
│  DEF: 40 SPD: 90        │
│                         │
│  Moves:                 │
│  • Thunderbolt (90)     │
│  • Quick Attack (40)    │
│  • Iron Tail (100)      │
│  • Electro Ball (80)    │
│                         │
│  💰 250 coins           │
│  [BUY NOW]              │
└─────────────────────────┘
```

## Notes

- Shop inventory refreshes every 24 hours
- Users can purchase duplicate Pokemon for deck building
- All purchased cards start at level 1 with 0 XP
- Rate limit: 10 purchases per minute per user
- Authentication required for all shop endpoints
