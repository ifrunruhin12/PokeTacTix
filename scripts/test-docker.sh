#!/bin/bash

# Test script to verify Docker setup
echo "Testing Docker setup..."
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running"
    exit 1
fi
echo "✅ Docker is running"

# Check docker-compose file
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ docker-compose.yml not found"
    exit 1
fi
echo "✅ docker-compose.yml found"

# Validate docker-compose file
if docker-compose config > /dev/null 2>&1; then
    echo "✅ docker-compose.yml is valid"
else
    echo "❌ docker-compose.yml has errors"
    exit 1
fi

echo ""
echo "🎉 Docker setup is ready!"
echo ""
echo "To start the database:"
echo "  ./docker-start.sh"
echo "  or"
echo "  make db-only"
