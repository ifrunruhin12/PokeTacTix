#!/bin/bash

# Database Migration Script
# Usage: ./scripts/migrate.sh [up|down|reset] [local|neon]
# Examples:
#   ./scripts/migrate.sh up          # Run migrations on local DB
#   ./scripts/migrate.sh up neon     # Run migrations on Neon
#   ./scripts/migrate.sh reset neon  # Reset Neon database

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

MIGRATION_DIR="internal/database/migrations"
ACTION="${1:-up}"
TARGET="${2:-local}"

# Neon database URL (from your connection string)
NEON_URL="postgresql://neondb_owner:npg_vk0HZgI1euyz@ep-billowing-fog-a16m0zo7-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require"

# Check if migrations directory exists
if [ ! -d "$MIGRATION_DIR" ]; then
    echo -e "${RED}❌ Migration directory not found: $MIGRATION_DIR${NC}"
    exit 1
fi

# Display target database
if [ "$TARGET" = "neon" ]; then
    echo -e "${BLUE}🎯 Target: Neon Production Database${NC}"
    echo ""
else
    echo -e "${BLUE}🎯 Target: Local Database${NC}"
    echo ""
fi

# Function to run migrations up
migrate_up() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}📝 Running migrations UP${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    migration_count=0
    for migration in $(ls $MIGRATION_DIR/*up.sql 2>/dev/null | sort); do
        migration_name=$(basename $migration)
        echo -e "${YELLOW}→ Running: $migration_name${NC}"
        
        if [ "$TARGET" = "neon" ]; then
            # Run on Neon
            if psql "$NEON_URL" < "$migration" 2>/dev/null; then
                echo -e "${GREEN}  ✓ Success${NC}"
            else
                echo -e "${RED}  ✗ Failed (might already be applied)${NC}"
            fi
        else
            # Try Docker first, then local psql
            if docker-compose exec -T postgres psql -U pokemon -d pokemon < "$migration" 2>/dev/null; then
                echo -e "${GREEN}  ✓ Success${NC}"
            elif PGPASSWORD=pokemon123 psql -h postgres -U pokemon -d pokemon < "$migration" 2>/dev/null; then
                echo -e "${GREEN}  ✓ Success${NC}"
            elif psql -U pokemon -d pokemon < "$migration" 2>/dev/null; then
                echo -e "${GREEN}  ✓ Success${NC}"
            else
                echo -e "${RED}  ✗ Failed (might already be applied)${NC}"
            fi
        fi
        
        migration_count=$((migration_count + 1))
        echo ""
    done
    
    echo -e "${GREEN}✅ Processed $migration_count migrations${NC}"
}

# Function to run migrations down
migrate_down() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}📝 Running migrations DOWN${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    migration_count=0
    for migration in $(ls $MIGRATION_DIR/*down.sql 2>/dev/null | sort -r); do
        migration_name=$(basename $migration)
        echo -e "${YELLOW}→ Running: $migration_name${NC}"
        
        if [ "$TARGET" = "neon" ]; then
            # Run on Neon
            if psql "$NEON_URL" < "$migration" 2>/dev/null; then
                echo -e "${GREEN}  ✓ Success${NC}"
            else
                echo -e "${RED}  ✗ Failed${NC}"
            fi
        else
            # Try Docker first, then local psql
            if docker-compose exec -T postgres psql -U pokemon -d pokemon < "$migration" 2>/dev/null; then
                echo -e "${GREEN}  ✓ Success${NC}"
            elif PGPASSWORD=pokemon123 psql -h postgres -U pokemon -d pokemon < "$migration" 2>/dev/null; then
                echo -e "${GREEN}  ✓ Success${NC}"
            elif psql -U pokemon -d pokemon < "$migration" 2>/dev/null; then
                echo -e "${GREEN}  ✓ Success${NC}"
            else
                echo -e "${RED}  ✗ Failed${NC}"
            fi
        fi
        
        migration_count=$((migration_count + 1))
        echo ""
    done
    
    echo -e "${GREEN}✅ Processed $migration_count migrations${NC}"
}

# Function to reset database
migrate_reset() {
    if [ "$TARGET" = "neon" ]; then
        echo -e "${RED}⚠️  WARNING: This will drop all tables in PRODUCTION (Neon) and re-run migrations!${NC}"
    else
        echo -e "${RED}⚠️  WARNING: This will drop all tables and re-run migrations!${NC}"
    fi
    read -p "Are you sure? [y/N] " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo ""
        migrate_down
        echo ""
        migrate_up
        echo ""
        echo -e "${GREEN}✅ Database reset complete!${NC}"
    else
        echo -e "${YELLOW}❌ Cancelled${NC}"
    fi
}

# Main
case "$ACTION" in
    up)
        migrate_up
        ;;
    down)
        migrate_down
        ;;
    reset)
        migrate_reset
        ;;
    *)
        echo -e "${RED}❌ Invalid action: $ACTION${NC}"
        echo ""
        echo "Usage: $0 [up|down|reset] [local|neon]"
        echo ""
        echo "Actions:"
        echo "  up     - Run all migrations (default)"
        echo "  down   - Rollback all migrations"
        echo "  reset  - Rollback and re-run all migrations"
        echo ""
        echo "Targets:"
        echo "  local  - Local database (default)"
        echo "  neon   - Neon production database"
        echo ""
        echo "Examples:"
        echo "  $0 up          # Run migrations on local DB"
        echo "  $0 up neon     # Run migrations on Neon"
        echo "  $0 reset neon  # Reset Neon database"
        echo ""
        exit 1
        ;;
esac
