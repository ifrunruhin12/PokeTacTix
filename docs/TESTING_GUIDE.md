# PokeTacTix API Testing Guide

This guide will help you set up and test the PokeTacTix API.

## Prerequisites

- Docker installed and running
- Go 1.21+ installed
- curl or Postman for API testing

## Quick Start

### 1. Set Up Database

Run the setup script to create and configure the PostgreSQL database:

```bash
./scripts/setup_db.sh
```

This will:
- Create a Docker container with PostgreSQL 15
- Set up the database with proper credentials
- Run all migrations to create tables
- Verify the setup

### 2. Start the API Server

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:3000`

### 3. Verify Server is Running

```bash
curl http://localhost:3000/health
```

Expected response:
```json
{
  "status": "healthy",
  "env": "development"
}
```

## API Testing

### Authentication Endpoints

#### Register a New User

```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test123!@#"
  }'
```

Expected response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "coins": 0,
    "created_at": "2025-11-11T17:12:03.770603Z",
    "updated_at": "2025-11-11T17:12:03.770603Z"
  }
}
```

**Password Requirements:**
- Minimum 8 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 number
- At least 1 special character

#### Login

```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123!@#"
  }'
```

Expected response: Same as registration

#### Save the Token

```bash
# Save token for subsequent requests
export TOKEN="your_token_here"
```

### Protected Endpoints

All endpoints below require authentication. Include the token in the Authorization header:

```bash
-H "Authorization: Bearer $TOKEN"
```

#### Get Player Stats

```bash
curl -X GET http://localhost:3000/api/profile/stats \
  -H "Authorization: Bearer $TOKEN"
```

Expected response:
```json
{
  "user_id": 1,
  "total_battles_1v1": 0,
  "wins_1v1": 0,
  "losses_1v1": 0,
  "total_battles_5v5": 0,
  "wins_5v5": 0,
  "losses_5v5": 0,
  "total_coins_earned": 0,
  "highest_level": 0,
  "updated_at": "0001-01-01T00:00:00Z"
}
```

#### Get Battle History

```bash
curl -X GET http://localhost:3000/api/profile/history \
  -H "Authorization: Bearer $TOKEN"
```

#### Get Achievements

```bash
curl -X GET http://localhost:3000/api/profile/achievements \
  -H "Authorization: Bearer $TOKEN"
```

#### Get Player Cards

```bash
curl -X GET http://localhost:3000/api/cards \
  -H "Authorization: Bearer $TOKEN"
```

#### Get Shop Inventory

```bash
curl -X GET http://localhost:3000/api/shop/inventory \
  -H "Authorization: Bearer $TOKEN"
```

#### Purchase a Pokemon

```bash
curl -X POST http://localhost:3000/api/shop/purchase \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "pokemon_name": "pikachu"
  }'
```

## Frontend Testing

### 1. Install Frontend Dependencies

```bash
cd frontend
npm install
```

### 2. Create Environment File

```bash
cp .env.example .env
```

The default configuration should work:
```
VITE_API_URL=http://localhost:3000
```

### 3. Start Frontend Dev Server

```bash
npm run dev
```

The frontend will be available at `http://localhost:5173`

### 4. Test the Flow

1. Open `http://localhost:5173` in your browser
2. Click "Get Started" or "Login / Register"
3. Register a new account
4. You should be redirected to the dashboard
5. Explore the different pages (Battle, Shop, Deck, Profile)

## Database Management

### Connect to Database

```bash
docker exec -it poketactix-db psql -U pokemon -d poketactix
```

### View Tables

```sql
\dt
```

### View Users

```sql
SELECT * FROM users;
```

### View Player Cards

```sql
SELECT * FROM player_cards;
```

### Reset Database

To completely reset the database:

```bash
# Stop and remove container
docker rm -f poketactix-db

# Run setup script again
./scripts/setup_db.sh
```

## Troubleshooting

### Database Connection Failed

1. Check if Docker is running:
   ```bash
   docker ps
   ```

2. Check if database container is running:
   ```bash
   docker ps | grep poketactix-db
   ```

3. Start the container if stopped:
   ```bash
   docker start poketactix-db
   ```

### API Server Won't Start

1. Check if port 3000 is already in use:
   ```bash
   lsof -i :3000
   ```

2. Check the .env file exists and has correct values

3. Verify database connection string in .env:
   ```
   DATABASE_URL=postgresql://pokemon:pokemon123@localhost:5432/poketactix?sslmode=disable
   ```

### Frontend Can't Connect to API

1. Verify API server is running on port 3000
2. Check CORS configuration in .env:
   ```
   CORS_ORIGINS=http://localhost:5173,http://localhost:3000
   ```
3. Check browser console for CORS errors

### Token Expired

JWT tokens expire after 24 hours by default. If you get 401 errors:

1. Login again to get a new token
2. Or adjust JWT_EXPIRATION in .env

## API Documentation

Full API documentation is available via Swagger UI:

```
http://localhost:3000/api/docs/
```

## Testing with Postman

Import the following base URL and environment variables:

- Base URL: `http://localhost:3000`
- Token: `{{token}}` (set after login/register)

Create a collection with the endpoints listed above.

## Next Steps

- Implement battle functionality
- Add more Pokemon to the shop
- Create deck management UI
- Add real-time battle updates
- Implement achievements system
