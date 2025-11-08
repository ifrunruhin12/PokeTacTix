.PHONY: help up down logs restart migrate migrate-down migrate-up seed test clean build

# Default target
help:
	@echo "PokeTacTix Development Commands"
	@echo "================================"
	@echo "make up          - Start all services (postgres, backend, frontend)"
	@echo "make down        - Stop all services"
	@echo "make logs        - View logs from all services"
	@echo "make restart     - Restart backend and frontend services"
	@echo "make migrate     - Run database migrations"
	@echo "make migrate-up  - Run migrations up"
	@echo "make migrate-down- Rollback last migration"
	@echo "make seed        - Seed database with initial data"
	@echo "make test        - Run backend tests"
	@echo "make clean       - Stop services and remove volumes"
	@echo "make build       - Build production binaries"

# Start all services
up:
	@echo "Starting all services..."
	docker-compose up -d
	@echo "Services started! Backend: http://localhost:3000, Frontend: http://localhost:5173"

# Stop all services
down:
	@echo "Stopping all services..."
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# View logs for specific service
logs-backend:
	docker-compose logs -f backend

logs-frontend:
	docker-compose logs -f frontend

logs-db:
	docker-compose logs -f postgres

# Restart services
restart:
	@echo "Restarting backend and frontend..."
	docker-compose restart backend frontend

# Run database migrations
migrate:
	@echo "Running database migrations..."
	docker-compose exec postgres psql -U pokemon -d poketactix -c "\i /docker-entrypoint-initdb.d/000001_create_users_table.up.sql"
	docker-compose exec postgres psql -U pokemon -d poketactix -c "\i /docker-entrypoint-initdb.d/000002_create_player_cards_table.up.sql"
	docker-compose exec postgres psql -U pokemon -d poketactix -c "\i /docker-entrypoint-initdb.d/000003_create_battle_history_table.up.sql"
	docker-compose exec postgres psql -U pokemon -d poketactix -c "\i /docker-entrypoint-initdb.d/000004_create_player_stats_table.up.sql"
	docker-compose exec postgres psql -U pokemon -d poketactix -c "\i /docker-entrypoint-initdb.d/000005_create_achievements_tables.up.sql"
	@echo "Migrations completed!"

# Run migrations up
migrate-up:
	@echo "Running migrations up..."
	@make migrate

# Rollback migrations
migrate-down:
	@echo "Rolling back migrations..."
	docker-compose exec postgres psql -U pokemon -d poketactix -c "\i /docker-entrypoint-initdb.d/000005_create_achievements_tables.down.sql"
	docker-compose exec postgres psql -U pokemon -d poketactix -c "\i /docker-entrypoint-initdb.d/000004_create_player_stats_table.down.sql"
	docker-compose exec postgres psql -U pokemon -d poketactix -c "\i /docker-entrypoint-initdb.d/000003_create_battle_history_table.down.sql"
	docker-compose exec postgres psql -U pokemon -d poketactix -c "\i /docker-entrypoint-initdb.d/000002_create_player_cards_table.down.sql"
	docker-compose exec postgres psql -U pokemon -d poketactix -c "\i /docker-entrypoint-initdb.d/000001_create_users_table.down.sql"
	@echo "Migrations rolled back!"

# Seed database with initial data
seed:
	@echo "Seeding database..."
	@echo "Note: Implement seed script in migrations/seed.sql"
	# docker-compose exec postgres psql -U pokemon -d poketactix -f /docker-entrypoint-initdb.d/seed.sql

# Run tests
test:
	@echo "Running backend tests..."
	docker-compose exec backend go test ./... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	docker-compose exec backend go test ./... -coverprofile=coverage.out
	docker-compose exec backend go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean everything (including volumes)
clean:
	@echo "Cleaning up..."
	docker-compose down -v
	rm -rf tmp/
	rm -rf coverage.out coverage.html
	@echo "Cleanup complete!"

# Build production binaries
build:
	@echo "Building production binary..."
	CGO_ENABLED=0 GOOS=linux go build -o bin/server ./server
	@echo "Build complete! Binary: bin/server"

# Database shell
db-shell:
	docker-compose exec postgres psql -U pokemon -d poketactix

# Backend shell
backend-shell:
	docker-compose exec backend sh

# Check service status
status:
	docker-compose ps

# Pull latest images
pull:
	docker-compose pull

# Rebuild containers
rebuild:
	docker-compose up -d --build

# View backend logs in real-time
watch-backend:
	docker-compose logs -f backend

# View frontend logs in real-time
watch-frontend:
	docker-compose logs -f frontend
