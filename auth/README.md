# Authentication Package

This package provides a complete authentication system for the PokeTacTix web application, including user registration, login, JWT token management, and rate limiting.

## Features

- **Password Security**: bcrypt hashing with cost factor 12
- **Password Validation**: Enforces strong password requirements (min 8 chars, uppercase, lowercase, number, special character)
- **JWT Tokens**: HS256 signing with 24-hour expiration
- **Rate Limiting**: Configurable rate limits for auth endpoints
- **Input Validation**: Username and email validation
- **Middleware**: Easy-to-use authentication middleware for protected routes

## Requirements Met

- **Requirement 2.1**: User registration endpoint
- **Requirement 2.2**: User login endpoint with credential validation
- **Requirement 2.3**: JWT token generation with 24-hour expiration
- **Requirement 2.4**: Token validation middleware
- **Requirement 2.6**: Password validation with security requirements
- **Requirement 14.1**: bcrypt password hashing with cost factor 12
- **Requirement 14.3**: JWT secret stored in environment variable
- **Requirement 14.5**: Rate limiting (5 req/min for login, 3 req/hour for register)

## Components

### Service (`service.go`)
Core authentication logic including:
- Password hashing and comparison
- Password validation
- Username validation
- Email validation

### JWT Service (`jwt.go`)
JWT token operations:
- Token generation
- Token validation
- Token refresh
- User ID extraction

### Middleware (`middleware.go`)
Fiber middleware for:
- Required authentication
- Optional authentication
- User context extraction

### Handlers (`handlers.go`)
HTTP request handlers:
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/me` - Get current user info
- `POST /api/auth/logout` - Logout (client-side token removal)

### Rate Limiting (`ratelimit.go`)
Rate limiter configurations:
- Login: 5 requests per minute
- Register: 3 requests per hour
- Default: 10 requests per minute

## Usage

### 1. Initialize Services

```go
import (
    "pokemon-cli/auth"
    "pokemon-cli/database"
)

// Initialize auth services
authService := auth.NewService()
jwtService, err := auth.NewJWTService()
if err != nil {
    log.Fatal(err)
}

// Initialize repositories
userRepo := database.NewUserRepository(database.GetDB())

// Initialize handlers
authHandler := auth.NewHandler(authService, jwtService, userRepo)
```

### 2. Register Routes

```go
// Register all auth routes with rate limiting
auth.RegisterRoutes(app, authHandler, jwtService)
```

### 3. Protect Routes

```go
// Protect a single route
app.Get("/api/profile", auth.Middleware(jwtService), profileHandler)

// Protect a route group
protected := app.Group("/api/protected", auth.Middleware(jwtService))
protected.Get("/cards", getCardsHandler)
protected.Post("/battle/start", startBattleHandler)
```

### 4. Extract User Info in Handlers

```go
func myHandler(c *fiber.Ctx) error {
    userID, ok := auth.GetUserID(c)
    if !ok {
        return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
    }
    
    username, _ := auth.GetUsername(c)
    
    // Use userID and username
    return c.JSON(fiber.Map{
        "user_id": userID,
        "username": username,
    })
}
```

## Environment Variables

Required environment variables:

```bash
# JWT Configuration (REQUIRED)
JWT_SECRET=your-secret-key-min-32-chars-256-bits

# Optional
JWT_EXPIRATION=24h
```

**Important**: The JWT_SECRET must be at least 32 characters (256 bits) for security.

## API Endpoints

### POST /api/auth/register

Register a new user account.

**Request Body:**
```json
{
  "username": "trainer_ash",
  "email": "ash@pokemon.com",
  "password": "SecurePass123!"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "trainer_ash",
    "email": "ash@pokemon.com",
    "coins": 0,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Error Responses:**
- `400 WEAK_PASSWORD`: Password doesn't meet requirements
- `400 INVALID_USERNAME`: Username format invalid
- `400 INVALID_EMAIL`: Email format invalid
- `409 USER_EXISTS`: Username or email already exists
- `429 RATE_LIMIT_EXCEEDED`: Too many registration attempts (3/hour limit)

### POST /api/auth/login

Authenticate and receive a JWT token.

**Request Body:**
```json
{
  "username": "trainer_ash",
  "password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "trainer_ash",
    "email": "ash@pokemon.com",
    "coins": 150,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Error Responses:**
- `401 INVALID_CREDENTIALS`: Invalid username or password
- `429 RATE_LIMIT_EXCEEDED`: Too many login attempts (5/minute limit)

### GET /api/auth/me

Get current authenticated user's information.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "user": {
    "id": 1,
    "username": "trainer_ash",
    "email": "ash@pokemon.com",
    "coins": 150,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Error Responses:**
- `401 UNAUTHORIZED`: Missing or invalid token
- `401 TOKEN_EXPIRED`: Token has expired
- `404 USER_NOT_FOUND`: User no longer exists

### POST /api/auth/logout

Logout (primarily for client-side token removal).

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "message": "Logged out successfully"
}
```

## Password Requirements

Passwords must meet the following criteria:
- Minimum 8 characters
- Maximum 100 characters
- At least one uppercase letter (A-Z)
- At least one lowercase letter (a-z)
- At least one digit (0-9)
- At least one special character (!@#$%^&*()_+-=[]{};':"\\|,.<>/?)

## Username Requirements

Usernames must meet the following criteria:
- 3-50 characters
- Alphanumeric characters and underscores only
- No spaces or special characters (except underscore)

## Email Requirements

Emails must:
- Be a valid email format
- Maximum 255 characters

## Security Features

1. **Password Hashing**: All passwords are hashed using bcrypt with cost factor 12
2. **JWT Tokens**: Signed with HS256 algorithm, 24-hour expiration
3. **Rate Limiting**: Prevents brute force attacks
4. **Input Validation**: All inputs are validated and sanitized
5. **Error Messages**: Generic error messages to prevent user enumeration
6. **Secure Headers**: Middleware sets appropriate security headers

## Testing

See `example_integration.go` for complete integration examples.

## Error Codes

| Code | Description |
|------|-------------|
| `INVALID_REQUEST` | Malformed request body |
| `INVALID_USERNAME` | Username doesn't meet requirements |
| `INVALID_EMAIL` | Email format is invalid |
| `WEAK_PASSWORD` | Password doesn't meet security requirements |
| `USER_EXISTS` | Username or email already registered |
| `INVALID_CREDENTIALS` | Wrong username or password |
| `UNAUTHORIZED` | Missing or invalid authentication |
| `TOKEN_EXPIRED` | JWT token has expired |
| `TOKEN_INVALID` | JWT token is malformed |
| `RATE_LIMIT_EXCEEDED` | Too many requests |
| `USER_NOT_FOUND` | User doesn't exist |
| `INTERNAL_ERROR` | Server error |

## Dependencies

- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/golang-jwt/jwt/v5` - JWT token handling
- `github.com/gofiber/fiber/v2` - HTTP framework
- `github.com/jackc/pgx/v5` - PostgreSQL driver
