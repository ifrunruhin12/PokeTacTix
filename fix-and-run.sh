#!/bin/bash

echo "🔧 Fixing and restarting PokeTacTix..."
echo ""

# Stop everything
echo "⏹️  Stopping services..."
docker-compose down -v

# Clean Docker cache
echo "🧹 Cleaning Docker cache..."
docker system prune -f

# Rebuild
echo "🔨 Rebuilding containers..."
docker-compose build --no-cache

# Start
echo "🚀 Starting services..."
./scripts/docker-dev.sh
