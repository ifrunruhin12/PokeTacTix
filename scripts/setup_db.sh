#!/bin/bash

# PokeTacTix Database Setup Script
# This script sets up a PostgreSQL database using Docker

set -e

echo "🚀 Setting up PokeTacTix database..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if container already exists
if docker ps -a | grep -q poketactix-db; then
    echo "📦 Container 'poketactix-db' already exists."
    
    # Check if it's running
    if docker ps | grep -q poketactix-db; then
        echo "✅ Database is already running!"
    else
        echo "▶️  Starting existing container..."
        docker start poketactix-db
        echo "✅ Database started!"
    fi
else
    echo "📦 Creating new PostgreSQL container..."
    docker run -d \
      --name poketactix-db \
      -e POSTGRES_USER=pokemon \
      -e POSTGRES_PASSWORD=pokemon123 \
      -e POSTGRES_DB=poketactix \
      -p 5432:5432 \
      postgres:15-alpine
    
    echo "⏳ Waiting for PostgreSQL to be ready..."
    sleep 5
    
    echo "✅ Database container created!"
fi

# Wait for database to be ready
echo "⏳ Waiting for database to accept connections..."
for i in {1..30}; do
    if docker exec poketactix-db pg_isready -U pokemon > /dev/null 2>&1; then
        echo "✅ Database is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ Database failed to start"
        exit 1
    fi
    sleep 1
done

# Run migrations
echo "📝 Running database migrations..."

MIGRATION_DIR="internal/database/migrations"

if [ ! -d "$MIGRATION_DIR" ]; then
    echo "❌ Migration directory not found: $MIGRATION_DIR"
    exit 1
fi

# Run migrations in order
for migration in $(ls $MIGRATION_DIR/*up.sql | sort); do
    echo "  Running: $(basename $migration)"
    docker exec -i poketactix-db psql -U pokemon -d poketactix < "$migration"
done

echo "✅ All migrations completed!"

# Test connection
echo "🧪 Testing database connection..."
docker exec poketactix-db psql -U pokemon -d poketactix -c "SELECT COUNT(*) as table_count FROM information_schema.tables WHERE table_schema = 'public';" | grep -q "6" && echo "✅ All tables created successfully!"

echo ""
echo "🎉 Database setup complete!"
echo ""
echo "📋 Connection details:"
echo "   Host: localhost"
echo "   Port: 5432"
echo "   Database: poketactix"
echo "   Username: pokemon"
echo "   Password: pokemon123"
echo ""
echo "🔗 Connection string:"
echo "   postgresql://pokemon:pokemon123@localhost:5432/poketactix?sslmode=disable"
echo ""
echo "💡 Useful commands:"
echo "   Stop database:    docker stop poketactix-db"
echo "   Start database:   docker start poketactix-db"
echo "   Remove database:  docker rm -f poketactix-db"
echo "   View logs:        docker logs poketactix-db"
echo "   Connect to DB:    docker exec -it poketactix-db psql -U pokemon -d poketactix"
