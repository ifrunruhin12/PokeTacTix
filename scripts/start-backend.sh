#!/bin/sh
set -e

echo "🚀 Starting PokeTacTix Backend..."
echo ""

# Wait for postgres to be ready
echo "⏳ Waiting for database..."
until PGPASSWORD=pokemon123 psql -h postgres -U pokemon -d poketactix -c '\q' 2>/dev/null; do
  sleep 1
done
echo "✅ Database is ready!"
echo ""

# Run migrations directly
echo "📝 Running database migrations..."
for migration in /app/internal/database/migrations/*up.sql; do
    if [ -f "$migration" ]; then
        echo "  → $(basename $migration)"
        PGPASSWORD=pokemon123 psql -h postgres -U pokemon -d poketactix -f "$migration" 2>&1 | grep -v "already exists" || true
    fi
done
echo "✅ Migrations complete!"
echo ""

# Start the application
echo "🎮 Starting application..."
exec air
