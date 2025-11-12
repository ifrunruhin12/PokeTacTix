#!/bin/bash

# Database Migration Script
# Usage: ./migrate.sh [up|down|reset]

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

MIGRATION_DIR="internal/database/migrations"
ACTION="${1:-up}"

# Check if migrations directory exists
if [ ! -d "$MIGRATION_DIR" ]; then
    echo -e "${RED}❌ Migration directory not found: $MIGRATION_DIR${NC}"
    exit 1
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
        
        migration_count=$((migration_count + 1))
        echo ""
    done
    
    echo -e "${GREEN}✅ Processed $migration_count migrations${NC}"
}

# Function to reset database
migrate_reset() {
    echo -e "${RED}⚠️  WARNING: This will drop all tables and re-run migrations!${NC}"
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
        echo "Usage: $0 [up|down|reset]"
        echo ""
        echo "Actions:"
        echo "  up     - Run all migrations (default)"
        echo "  down   - Rollback all migrations"
        echo "  reset  - Rollback and re-run all migrations"
        echo ""
        exit 1
        ;;
esac
