# Code Quality Improvements

## SQL Query Management

### Problem
Previously, SQL queries were scattered throughout the codebase, particularly in `rewards.go`. This made the code:
- Hard to maintain and review
- Prone to errors and inconsistencies
- Difficult to test in isolation
- Repetitive with similar patterns

### Solution: Repository Pattern
We've centralized database operations in the repository layer (`internal/battle/repository.go`), following these principles:

#### 1. **Single Responsibility**
Each repository method handles one specific database operation:
```go
// Good: Centralized in repository
func (r *Repository) UpdatePlayerStatsInTx(ctx, tx, userID, mode, result, coins) error

// Bad: SQL scattered in business logic
func ApplyAllRewards(...) {
    // SQL queries mixed with business logic
}
```

#### 2. **Separation of Concerns**
- **Repository Layer**: Handles all SQL queries and database operations
- **Service/Business Layer**: Handles business logic and orchestration
- **Handler Layer**: Handles HTTP requests and responses

#### 3. **Benefits**
- ✅ **Maintainability**: All SQL in one place
- ✅ **Testability**: Can mock repository methods
- ✅ **Reusability**: Same query logic used everywhere
- ✅ **Security**: Centralized validation and parameterization
- ✅ **Code Review**: Easier to spot SQL issues

### Example Refactor

**Before:**
```go
// rewards.go - SQL mixed with business logic
func updatePlayerStatsInTransaction(ctx, tx, userID, mode, result, coins) {
    var query string
    if mode == "1v1" {
        if result == "win" {
            query = `INSERT INTO player_stats ...` // 50+ lines of SQL
        }
    }
    tx.Exec(ctx, query, userID, coins)
}
```

**After:**
```go
// repository.go - SQL centralized
func (r *Repository) UpdatePlayerStatsInTx(ctx, tx, userID, mode, result, coins) error {
    // All SQL logic here
}

// rewards.go - Clean business logic
func ApplyAllRewards(..., repo *Repository) error {
    // Business logic
    err := repo.UpdatePlayerStatsInTx(ctx, tx, userID, mode, result, coins)
    // More business logic
}
```

## Future Improvements

### Option 1: Query Builder (Squirrel)
For complex dynamic queries:
```go
import sq "github.com/Masterminds/squirrel"

query := sq.Insert("player_stats").
    Columns("user_id", "wins_1v1", "total_coins_earned").
    Values(userID, 1, coins).
    Suffix("ON CONFLICT (user_id) DO UPDATE SET wins_1v1 = wins_1v1 + 1").
    PlaceholderFormat(sq.Dollar)

sql, args, _ := query.ToSql()
```

### Option 2: sqlx (Enhanced SQL)
For better struct scanning and named parameters:
```go
import "github.com/jmoiron/sqlx"

type PlayerStats struct {
    UserID     int `db:"user_id"`
    Wins1v1    int `db:"wins_1v1"`
    TotalCoins int `db:"total_coins_earned"`
}

var stats PlayerStats
err := db.Get(&stats, "SELECT * FROM player_stats WHERE user_id = $1", userID)
```

### Option 3: GORM (Full ORM)
For complex relationships and migrations:
```go
type PlayerStats struct {
    UserID     uint
    Wins1v1    int
    TotalCoins int
}

db.Model(&PlayerStats{}).Where("user_id = ?", userID).Update("wins_1v1", gorm.Expr("wins_1v1 + ?", 1))
```

## Recommendation

For this project, the **Repository Pattern** we just implemented is sufficient because:
1. Queries are relatively simple
2. We're already using pgx (excellent driver)
3. No need for complex query building
4. Keeps dependencies minimal

If the project grows and you need:
- Complex dynamic queries → Use **Squirrel**
- Better struct mapping → Use **sqlx**
- Full ORM features → Use **GORM**

## Best Practices

1. **Always use parameterized queries** (we do this)
2. **Validate inputs before building queries** (we do this)
3. **Keep SQL in repository layer** (we now do this)
4. **Use transactions for multi-step operations** (we do this)
5. **Handle errors consistently** (we do this)

## Testing

With the repository pattern, you can now easily mock database operations:

```go
type MockRepository struct {
    UpdatePlayerStatsInTxFunc func(ctx, tx, userID, mode, result, coins) error
}

func (m *MockRepository) UpdatePlayerStatsInTx(ctx, tx, userID, mode, result, coins) error {
    return m.UpdatePlayerStatsInTxFunc(ctx, tx, userID, mode, result, coins)
}

// In tests
mockRepo := &MockRepository{
    UpdatePlayerStatsInTxFunc: func(...) error {
        return nil // or return error to test error handling
    },
}
```

This makes unit testing much easier!
