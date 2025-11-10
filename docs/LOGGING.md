# Structured Logging with slog

The PokeTacTix API uses Go's built-in `log/slog` package for structured logging. This provides better observability and makes it easier to parse logs in production.

## Features

- **Structured logging**: All logs are output as JSON in production for easy parsing
- **Human-readable**: Development mode uses text format for easier reading
- **Context fields**: Add structured fields to log entries
- **Log levels**: DEBUG, INFO, WARN, ERROR

## Usage

### Basic Logging

```go
logger := logger.New(logger.INFO)

// Simple log messages
logger.Info("Server starting", "port", 3000)
logger.Error("Database connection failed", "error", err)
logger.Debug("Processing request", "user_id", 123, "endpoint", "/api/users")
logger.Warn("Rate limit approaching", "current", 95, "limit", 100)
```

### Output Format

**Development (Text):**
```
time=2024-01-15T10:30:45.123Z level=INFO msg="Server starting" port=3000
time=2024-01-15T10:30:46.456Z level=ERROR msg="Database connection failed" error="connection timeout"
```

**Production (JSON):**
```json
{"time":"2024-01-15T10:30:45.123Z","level":"INFO","msg":"Server starting","port":3000}
{"time":"2024-01-15T10:30:46.456Z","level":"ERROR","msg":"Database connection failed","error":"connection timeout"}
```

### Adding Context

Create a logger with persistent context fields:

```go
// Create a logger with user context
userLogger := logger.With("user_id", userID, "username", username)

// All subsequent logs will include these fields
userLogger.Info("User logged in")
userLogger.Debug("Fetching user data")
userLogger.Error("Failed to update profile", "error", err)
```

### Log Levels

Set the minimum log level when creating the logger:

```go
// Development: show all logs including DEBUG
devLogger := logger.NewText(logger.DEBUG)

// Production: only INFO and above
prodLogger := logger.New(logger.INFO)

// Critical systems: only WARN and ERROR
criticalLogger := logger.New(logger.WARN)
```

## Best Practices

### 1. Use Structured Fields

❌ **Bad:**
```go
logger.Info(fmt.Sprintf("User %d logged in from %s", userID, ip))
```

✅ **Good:**
```go
logger.Info("User logged in", "user_id", userID, "ip", ip)
```

### 2. Include Context

Always include relevant context fields:

```go
logger.Error("Database query failed",
    "error", err,
    "query", "SELECT * FROM users",
    "duration_ms", duration.Milliseconds(),
    "user_id", userID,
)
```

### 3. Use Appropriate Levels

- **DEBUG**: Detailed information for debugging (not shown in production)
- **INFO**: General informational messages (server started, request processed)
- **WARN**: Warning messages (deprecated API used, rate limit approaching)
- **ERROR**: Error messages (database connection failed, invalid input)

### 4. Don't Log Sensitive Data

❌ **Never log:**
- Passwords
- API keys
- Credit card numbers
- Personal identification numbers

✅ **Safe to log:**
- User IDs (not usernames if sensitive)
- Request IDs
- Error messages
- Performance metrics

## Examples from the Codebase

### Server Startup
```go
appLogger.Info("Starting PokeTacTix API", 
    "env", cfg.Server.Env, 
    "port", cfg.Server.Port,
)
```

### Database Connection
```go
if err := database.InitDB(&cfg.Database); err != nil {
    appLogger.Error("Failed to initialize database", "error", err)
    os.Exit(1)
}
appLogger.Info("Database connection established")
```

### Request Handling
```go
logger.Info("Battle started",
    "user_id", userID,
    "mode", battleMode,
    "battle_id", battleID,
)
```

### Error Handling
```go
if err != nil {
    logger.Error("Failed to process move",
        "error", err,
        "battle_id", battleID,
        "move", move,
        "user_id", userID,
    )
    return err
}
```

## Migration from Old Logger

The new logger is backward compatible with the old API:

```go
// Old way (still works)
logger.Info("message", "key", "value")
logger.Error("error occurred", "error", err)

// New way (same thing)
logger.Info("message", "key", "value")
logger.Error("error occurred", "error", err)
```

## Advanced Usage

### Access Underlying slog.Logger

For advanced slog features:

```go
slogger := logger.GetSlog()
slogger.LogAttrs(context.Background(), slog.LevelInfo,
    "Complex log entry",
    slog.String("key", "value"),
    slog.Int("count", 42),
)
```

### Custom Handlers

Create custom handlers for specific needs:

```go
// File output
file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
handler := slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo})
customLogger := &logger.Logger{GetSlog: slog.New(handler)}
```
