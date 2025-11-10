# PokeTacTix API Documentation

## Overview

The PokeTacTix API provides comprehensive endpoints for user authentication, Pokemon card management, battles, shop operations, and player statistics. All endpoints follow RESTful conventions and return JSON responses.

## Accessing the Documentation

### Swagger UI (Interactive)

Once the server is running, you can access the interactive API documentation at:

```
http://localhost:8080/api/docs/
```

In production:
```
https://your-domain.com/api/docs/
```

The Swagger UI provides:
- Interactive API testing
- Request/response examples
- Schema definitions
- Authentication testing

### OpenAPI Specification

The raw OpenAPI 3.0 specification file is available at:
```
http://localhost:8080/api/docs/swagger.yaml
```

You can import this file into tools like:
- Postman
- Insomnia
- API testing frameworks
- Code generators

## Quick Start

### 1. Register a New User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "trainer_ash",
    "email": "ash@pokemon.com",
    "password": "SecurePass123!"
  }'
```

Response includes:
- User information
- JWT token (valid for 24 hours)
- 5 starter Pokemon cards

### 2. Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "trainer_ash",
    "password": "SecurePass123!"
  }'
```

Save the returned token for authenticated requests.

### 3. Make Authenticated Requests

Include the JWT token in the Authorization header:

```bash
curl -X GET http://localhost:8080/api/cards \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## API Endpoints Overview

### Authentication
- `POST /api/auth/register` - Create new account
- `POST /api/auth/login` - Login and get JWT token
- `GET /api/auth/me` - Get current user info

### Cards
- `GET /api/cards` - Get all owned cards
- `GET /api/cards/deck` - Get current battle deck
- `PUT /api/cards/deck` - Update battle deck (5 cards)

### Battle
- `POST /api/battle/start` - Start 1v1 or 5v5 battle
- `GET /api/battle/state` - Get current battle state
- `POST /api/battle/move` - Submit battle move
- `POST /api/battle/switch` - Switch Pokemon (5v5)
- `POST /api/battle/select-reward` - Select AI Pokemon after 5v5 win

### Shop
- `GET /api/shop/inventory` - View available Pokemon
- `POST /api/shop/purchase` - Buy Pokemon with coins

### Profile
- `GET /api/profile/stats` - Get player statistics
- `GET /api/profile/history` - Get battle history
- `GET /api/profile/achievements` - Get achievements

## Authentication

### JWT Token

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Token Details:**
- Algorithm: HS256
- Expiration: 24 hours
- Payload includes: user_id, username

### Password Requirements

Passwords must meet these criteria:
- Minimum 8 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 number
- At least 1 special character (!@#$%^&*)

## Rate Limiting

Certain endpoints have rate limits to prevent abuse:

| Endpoint | Limit |
|----------|-------|
| `/api/auth/login` | 5 requests/minute |
| `/api/auth/register` | 3 requests/hour |
| `/api/shop/purchase` | 10 requests/minute |
| `/api/battle/move` | 100 requests/minute |

When rate limited, you'll receive a 429 status code.

## Error Responses

All errors follow a consistent format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {}
  }
}
```

### Common Error Codes

| Code | Status | Description |
|------|--------|-------------|
| `INVALID_CREDENTIALS` | 401 | Wrong username/password |
| `TOKEN_EXPIRED` | 401 | JWT token expired |
| `TOKEN_INVALID` | 401 | Malformed JWT token |
| `UNAUTHORIZED` | 403 | Missing authentication |
| `USER_EXISTS` | 409 | Username/email taken |
| `WEAK_PASSWORD` | 400 | Password doesn't meet requirements |
| `BATTLE_NOT_FOUND` | 404 | Battle session not found |
| `INVALID_MOVE` | 400 | Move not allowed |
| `NOT_YOUR_TURN` | 400 | Not player's turn |
| `INSUFFICIENT_STAMINA` | 400 | Not enough stamina |
| `INSUFFICIENT_COINS` | 402 | Not enough coins |
| `ITEM_NOT_FOUND` | 404 | Pokemon not in shop |
| `INVALID_DECK` | 400 | Deck must have 5 cards |
| `CARD_NOT_FOUND` | 404 | Card doesn't exist |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |

## Battle System

### Battle Modes

**1v1 Mode:**
- Single Pokemon battle
- Quick matches (~1-2 minutes)
- Rewards: 50 coins (win), 10 coins (loss)
- XP: 20 XP to winning Pokemon

**5v5 Mode:**
- Full deck battle (5 Pokemon each)
- Strategic gameplay (~3-5 minutes)
- Rewards: 150 coins (win), 25 coins (loss)
- XP: 15 XP to each participating Pokemon
- Bonus: Select one AI Pokemon after victory

### Available Moves

| Move | Description | Stamina Cost |
|------|-------------|--------------|
| `attack` | Use a Pokemon move | Varies by move |
| `defend` | Reduce incoming damage | (HP_max + 1) / 2 |
| `pass` | Skip turn, restore stamina | 0 (restores stamina) |
| `sacrifice` | KO current Pokemon, heal next | 0 |
| `surrender` | End battle (counts as loss) | 0 |

### Card Visibility

In 5v5 battles:
- You can see all 5 of your Pokemon
- You can only see the AI's active Pokemon
- Inactive AI Pokemon appear as face-down Pokeball images
- AI Pokemon are revealed when they become active

## Leveling System

Pokemon cards gain XP and level up through battles:

- **XP per Battle:** 20 (1v1), 15 (5v5)
- **XP to Level Up:** 100 × current_level
- **Level Cap:** 50
- **Stat Increases per Level:**
  - HP: +3%
  - Attack: +2%
  - Defense: +2%
  - Speed: +1%

## Shop System

### Pricing

| Rarity | Price | Discount Price |
|--------|-------|----------------|
| Common | 100 | N/A |
| Uncommon | 250 | N/A |
| Rare | 500 | N/A |
| Legendary | 2500 | 1500 (40% off) |
| Mythical | 5000 | 3500 (30% off) |

### Inventory

- Refreshes every 24 hours
- 10-15 common/uncommon Pokemon
- 5-8 rare Pokemon
- 15% chance for 1-2 legendary/mythical Pokemon

## Testing with Postman

1. Import the OpenAPI spec from `/api/docs/swagger.yaml`
2. Create an environment with:
   - `base_url`: http://localhost:8080
   - `token`: (will be set after login)
3. Register/login to get a token
4. Set the token in your environment
5. Use `{{token}}` in Authorization headers

## Testing with cURL

See the examples in the Swagger UI or use the quick start examples above.

## Support

For issues or questions:
- Check the interactive documentation at `/api/docs/`
- Review error messages for specific guidance
- Ensure JWT tokens haven't expired (24-hour validity)

## Version

Current API Version: 1.0.0
