.PHONY: help dev local stop logs status db-shell clean restart build

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo ''
	@echo '╔════════════════════════════════════════════════════════╗'
	@echo '║                                                        ║'
	@echo '║              PokeTacTix - Make Commands                ║'
	@echo '║                                                        ║'
	@echo '╚════════════════════════════════════════════════════════╝'
	@echo ''
	@echo 'Usage: make [target]'
	@echo ''
	@echo '🚀 Quick Start:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "dev|local"
	@echo ''
	@echo '📊 Management:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "stop|logs|status|restart|clean"
	@echo ''
	@echo '🗄️  Database:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "db-"
	@echo ''
	@echo '🔧 Advanced:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "build|migrate|seed|test"
	@echo ''

# ============================================================================
# Quick Start Commands
# ============================================================================

dev: ## 🐳 Start everything with Docker (RECOMMENDED)
	@./scripts/docker-dev.sh

local: ## 💻 Setup for local development (no Docker)
	@./scripts/local-dev.sh

# ============================================================================
# Management Commands
# ============================================================================

stop: ## ⏹️  Stop all Docker services
	@echo "Stopping all services..."
	@docker-compose down
	@echo "✅ All services stopped"

logs: ## 📋 View logs from all services
	@docker-compose logs -f

status: ## 📊 Show status of all services
	@echo ""
	@echo "╔════════════════════════════════════════════════════════╗"
	@echo "║              Service Status                            ║"
	@echo "╚════════════════════════════════════════════════════════╝"
	@echo ""
	@docker-compose ps
	@echo ""

restart: ## 🔄 Restart all services
	@echo "Restarting services..."
	@docker-compose restart
	@echo "✅ Services restarted"

clean: ## 🧹 Stop services and remove all data
	@echo "⚠️  This will delete all database data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		rm -rf tmp/; \
		echo "✅ Cleaned up successfully"; \
	else \
		echo "❌ Cancelled"; \
	fi

# ============================================================================
# Database Commands
# ============================================================================

db-shell: ## 🗄️  Open PostgreSQL shell
	@docker-compose exec postgres psql -U pokemon -d pokemon

db-logs: ## 📋 View database logs
	@docker-compose logs -f postgres

db-only: ## 🗄️  Start only the database
	@echo "Starting database..."
	@docker-compose up -d postgres
	@echo "✅ Database started"
	@echo ""
	@echo "Connection: postgresql://pokemon:pokemon123@localhost:5432/pokemon"

# ============================================================================
# Advanced Commands
# ============================================================================

build: ## 🔨 Rebuild all Docker containers
	@echo "Building containers..."
	@docker-compose build
	@echo "✅ Build complete"

rebuild: ## 🔨 Rebuild and restart everything
	@echo "Rebuilding and restarting..."
	@docker-compose down
	@docker-compose build
	@docker-compose up -d
	@echo "✅ Rebuild complete"

migrate: ## 🔄 Run database migrations (up)
	@./scripts/migrate.sh up

migrate-down: ## ⬇️  Rollback database migrations
	@./scripts/migrate.sh down

migrate-reset: ## 🔄 Reset database (down + up)
	@./scripts/migrate.sh reset

seed: ## 🌱 Seed database with initial data
	@echo "Seeding database..."
	@docker-compose exec backend go run migrations/seed.go 2>/dev/null || echo "⚠️  Seed script not found"

test: ## 🧪 Run tests
	@docker-compose exec backend go test ./...

# ============================================================================
# Individual Service Commands
# ============================================================================

backend-logs: ## 📋 View backend logs
	@docker-compose logs -f backend

frontend-logs: ## 📋 View frontend logs
	@docker-compose logs -f frontend

backend-only: ## 🔧 Start database and backend only
	@docker-compose up -d postgres backend

frontend-only: ## 🔧 Start frontend only
	@docker-compose up -d frontend
