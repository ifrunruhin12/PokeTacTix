# Database Migrations

How database migrations work in PokeTacTix.

## 🎯 Overview

Migrations run automatically when you start with `make dev`.

## 📁 Migration Files

Located in `internal/database/migrations/`:

```
000001_create_users_table.up.sql       # Creates users table
000001_create_users_table.down.sql     # Drops users table
000002_create_player_cards_table.up.sql
000002_create_player_cards_table.down.sql
...
```

- `.up.sql` - Apply changes
- `.down.sql` - Revert changes

## 🚀 Running Migrations

**Automatic (Docker):**
```bash
make dev    # Runs migrations automatically
```

**Manual:**
```bash
make migrate         # Run migrations
make migrate-down    # Rollback migrations
make migrate-reset   # Reset database
```

## 📝 Creating New Migrations

1. **Create files:**
```bash
touch internal/database/migrations/000006_add_feature.up.sql
touch internal/database/migrations/000006_add_feature.down.sql
```

2. **Write SQL:**

`000006_add_feature.up.sql`:
```sql
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);
```

`000006_add_feature.down.sql`:
```sql
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
```

3. **Run:**
```bash
make migrate
```

## 🔍 Check Status

```bash
make db-shell    # Open database
\dt              # List tables
\q               # Exit
```

## ❓ Troubleshooting

**Migration failed?**
```bash
make migrate-down    # Rollback
# Fix the SQL
make migrate         # Try again
```

**Start fresh?**
```bash
make clean
make dev
```

---

**Need more?** See [development.md](development.md)
