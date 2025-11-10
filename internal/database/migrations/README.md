# Database Migrations

This directory contains SQL migration files for the PokeTacTix database schema.

## Migration Files

### 000001 - Users Table
- Creates `users` table with authentication and coin balance
- Adds indexes on username and email
- Includes auto-update trigger for `updated_at`

### 000002 - Player Cards Table
- Creates `player_cards` table for Pokemon cards owned by players
- Includes leveling system (level, XP)
- Tracks legendary and mythical Pokemon
- Manages deck configuration (5 cards per user)
- Enforces deck size constraint with trigger

### 000003 - Battle History Table
- Creates `battle_history` table for tracking battles
- Records mode (1v1 or 5v5), result, coins earned, duration
- Indexed for efficient querying by user and date

### 000004 - Player Stats Table
- Creates `player_stats` table for player statistics
- Tracks separate stats for 1v1 and 5v5 battles
- Records total coins earned and highest level achieved
- Includes consistency checks via trigger

### 000005 - Achievements Tables
- Creates `achievements` table with achievement definitions
- Creates `user_achievements` table for tracking unlocked achievements
- Seeds 8 default achievements

## Running Migrations

### Using Docker Compose

```bash
# Run all migrations
make migrate

# Rollback all migrations
make migrate-down
```

### Manual Execution

```bash
# Connect to database
psql -U pokemon -d poketactix

# Run migrations in order
\i migrations/000001_create_users_table.up.sql
\i migrations/000002_create_player_cards_table.up.sql
\i migrations/000003_create_battle_history_table.up.sql
\i migrations/000004_create_player_stats_table.up.sql
\i migrations/000005_create_achievements_tables.up.sql
```

### Rollback

```bash
# Rollback in reverse order
\i migrations/000005_create_achievements_tables.down.sql
\i migrations/000004_create_player_stats_table.down.sql
\i migrations/000003_create_battle_history_table.down.sql
\i migrations/000002_create_player_cards_table.down.sql
\i migrations/000001_create_users_table.down.sql
```

## Schema Overview

```
users (id, username, email, password_hash, coins, created_at, updated_at)
  ↓ CASCADE DELETE
  ├── player_cards (id, user_id, pokemon_name, level, xp, stats, ...)
  ├── battle_history (id, user_id, mode, result, coins_earned, ...)
  ├── player_stats (user_id, wins, losses, total_coins_earned, ...)
  └── user_achievements (user_id, achievement_id, unlocked_at)
        ↓
      achievements (id, name, description, requirement_type, ...)
```

## Key Features

### Automatic Timestamps
All tables with `updated_at` columns use a trigger to automatically update the timestamp on row updates.

### Cascade Deletes
When a user is deleted, all related data (cards, battles, stats, achievements) are automatically deleted.

### Data Integrity
- Check constraints ensure valid data (e.g., level between 1-50, coins >= 0)
- Foreign key constraints maintain referential integrity
- Triggers enforce business rules (e.g., max 5 cards in deck)

### Indexes
Strategic indexes on frequently queried columns for optimal performance:
- User lookups by username/email
- Card queries by user_id and deck status
- Battle history by user and date

## Default Achievements

The following achievements are seeded automatically:

1. **First Victory** - Win your first battle
2. **Veteran Trainer** - Win 10 battles
3. **Elite Trainer** - Win 50 battles
4. **Champion** - Win 100 battles
5. **Legendary Collector** - Obtain a legendary Pokemon
6. **Mythical Master** - Obtain a mythical Pokemon
7. **Max Level** - Get a Pokemon to level 50
8. **Coin Hoarder** - Accumulate 5000 coins

## Adding New Migrations

When adding new migrations:

1. Create two files:
   - `NNNNNN_description.up.sql` - Apply changes
   - `NNNNNN_description.down.sql` - Revert changes

2. Use sequential numbering (000006, 000007, etc.)

3. Always test both up and down migrations

4. Include appropriate indexes and constraints

5. Update this README with migration details

## Production Deployment

For production (Neon):

1. Create Neon PostgreSQL database
2. Get connection string from Neon dashboard
3. Set DATABASE_URL environment variable
4. Run migrations using the Makefile or manual execution
5. Verify all tables and indexes are created

## Troubleshooting

### Migration Fails

```bash
# Check current database state
\dt  # List tables
\di  # List indexes

# View error details in logs
make logs-db
```

### Rollback Issues

If a rollback fails, you may need to manually drop objects:

```sql
DROP TABLE IF EXISTS table_name CASCADE;
DROP FUNCTION IF EXISTS function_name CASCADE;
```

### Duplicate Migrations

If migrations run multiple times, some may fail due to existing objects. Use `IF EXISTS` and `IF NOT EXISTS` clauses to make migrations idempotent.

## Best Practices

1. **Always test migrations** in development before production
2. **Backup data** before running migrations in production
3. **Use transactions** for complex migrations
4. **Document changes** in this README
5. **Version control** all migration files
6. **Test rollbacks** to ensure they work correctly

## References

- PostgreSQL Documentation: https://www.postgresql.org/docs/
- Neon Documentation: https://neon.tech/docs
- Database Design: See `design.md` in specs directory
