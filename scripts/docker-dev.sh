#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                                                        ║${NC}"
echo -e "${BLUE}║        🐳 PokeTacTix Development Environment           ║${NC}"
echo -e "${BLUE}║                                                        ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running!${NC}"
    echo -e "${YELLOW}Please start Docker Desktop and try again.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Docker is running${NC}"
echo ""

# Stop any existing containers
echo -e "${YELLOW}🧹 Cleaning up old containers...${NC}"
docker-compose down > /dev/null 2>&1

# Start all services
echo -e "${GREEN}🚀 Starting all services...${NC}"
echo ""

# Start services with output
docker-compose up -d

# Wait a moment for services to initialize
sleep 2

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✨ Services Status:${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
docker-compose ps
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Wait for database to be healthy
echo -e "${YELLOW}⏳ Waiting for database to be ready...${NC}"
max_attempts=30
attempt=0
until docker-compose exec -T postgres pg_isready -U pokemon > /dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ $attempt -eq $max_attempts ]; then
        echo -e "${RED}❌ Database failed to start after ${max_attempts} attempts${NC}"
        echo -e "${YELLOW}Run 'docker-compose logs postgres' to see what went wrong${NC}"
        exit 1
    fi
    echo -e "${YELLOW}   Attempt $attempt/$max_attempts...${NC}"
    sleep 2
done

echo -e "${GREEN}✅ Database is ready!${NC}"
echo ""

# Run database migrations
echo -e "${YELLOW}📝 Running database migrations...${NC}"
MIGRATION_DIR="internal/database/migrations"

if [ -d "$MIGRATION_DIR" ]; then
    migration_count=0
    for migration in $(ls $MIGRATION_DIR/*up.sql 2>/dev/null | sort); do
        migration_name=$(basename $migration)
        echo -e "${BLUE}   → $migration_name${NC}"
        docker-compose exec -T postgres psql -U pokemon -d poketactix < "$migration" > /dev/null 2>&1 || true
        migration_count=$((migration_count + 1))
    done
    
    if [ $migration_count -gt 0 ]; then
        echo -e "${GREEN}✅ Ran $migration_count migrations${NC}"
    else
        echo -e "${YELLOW}⚠️  No migration files found${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Migration directory not found, skipping migrations${NC}"
fi
echo ""

# Wait for backend to be ready
echo -e "${YELLOW}⏳ Waiting for backend to be ready...${NC}"
max_attempts=30
attempt=0
until curl -s http://localhost:3000/health > /dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ $attempt -eq $max_attempts ]; then
        echo -e "${YELLOW}⚠️  Backend is taking longer than expected${NC}"
        echo -e "${YELLOW}   It might still be starting up. Check logs with: make logs${NC}"
        break
    fi
    sleep 2
done

if curl -s http://localhost:3000/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Backend is ready!${NC}"
fi
echo ""

# Wait for frontend to be ready
echo -e "${YELLOW}⏳ Waiting for frontend to be ready...${NC}"
max_attempts=30
attempt=0
until curl -s http://localhost:5173 > /dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ $attempt -eq $max_attempts ]; then
        echo -e "${YELLOW}⚠️  Frontend is taking longer than expected${NC}"
        echo -e "${YELLOW}   It might still be starting up. Check logs with: make logs${NC}"
        break
    fi
    sleep 2
done

if curl -s http://localhost:5173 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Frontend is ready!${NC}"
fi
echo ""

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                                                        ║${NC}"
echo -e "${BLUE}║              🎉 All Services Running! 🎉               ║${NC}"
echo -e "${BLUE}║                                                        ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}📱 Access your application:${NC}"
echo -e "   ${BLUE}Frontend:${NC}     http://localhost:5173"
echo -e "   ${BLUE}Backend API:${NC}  http://localhost:3000"
echo -e "   ${BLUE}API Docs:${NC}     http://localhost:3000/api/docs"
echo ""
echo -e "${GREEN}🗄️  Database Connection:${NC}"
echo -e "   ${BLUE}Host:${NC}         localhost"
echo -e "   ${BLUE}Port:${NC}         5432"
echo -e "   ${BLUE}Database:${NC}     poketactix"
echo -e "   ${BLUE}User:${NC}         pokemon"
echo -e "   ${BLUE}Password:${NC}     pokemon123"
echo ""
echo -e "${GREEN}📊 Useful Commands:${NC}"
echo -e "   ${BLUE}View logs:${NC}           make logs"
echo -e "   ${BLUE}Stop services:${NC}       make stop"
echo -e "   ${BLUE}Restart services:${NC}    make restart"
echo -e "   ${BLUE}Database shell:${NC}      make db-shell"
echo -e "   ${BLUE}Clean everything:${NC}    make clean"
echo ""
echo -e "${YELLOW}💡 Tip: Run 'make logs' in another terminal to see live logs${NC}"
echo ""
