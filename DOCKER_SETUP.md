# Docker Development Setup

This guide explains how to set up and run PokeTacTix using Docker Compose for local development.

## Prerequisites

- Docker Desktop (or Docker Engine + Docker Compose)
- Make (optional, but recommended)

## Quick Start

### 1. Start All Services

```bash
make up
```

This will start:
- PostgreSQL database on port 5432
- Backend API on port 3000
- Frontend on port 5173

### 2. Run Database Migrations

```bash
make migrate
```

This creates all necessary database tables and seeds initial data.

### 3. Access the Application

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:3000
- **Database**: localhost:5432

## Available Commands

### Service Management

```bash
make up          # Start all services
make down        # Stop all services
make restart     # Restart backend and frontend
make status      # Check service status
make rebuild     # Rebuild and restart containers
```

### Logs

```bash
make logs              # View all logs
make logs-backend      # View backend logs only
make logs-frontend     # View frontend logs only
make logs-db           # View database logs only
make watch-backend     # Follow backend logs in real-time
make watch-frontend    # Follow frontend logs in real-time
```

### Database

```bash
make migrate           # Run all migrations
make migrate-up        # Run migrations up
make migrate-down      # Rollback migrations
make seed              # Seed database with test data
make db-shell          # Open PostgreSQL shell
```

### Testing

```bash
make test              # Run all tests
make test-coverage     # Run tests with coverage report
```

### Cleanup

```bash
make clean             # Stop services and remove volumes
```

## Service Details

### PostgreSQL Database

- **Image**: postgres:15-alpine
- **Port**: 5432
- **Database**: poketactix
- **User**: pokemon
- **Password**: pokemon123
- **Volume**: postgres_data (persists data)

### Backend (Go + Fiber)

- **Port**: 3000
- **Live Reload**: Enabled with Air
- **Hot Reload**: Code changes trigger automatic rebuild
- **Environment**: Development mode

### Frontend (React + Vite)

- **Port**: 5173
- **Hot Reload**: Enabled with Vite HMR
- **API Proxy**: Configured to backend at localhost:3000

## Environment Variables

The following environment variables are configured in `docker-compose.yml`:

### Backend

```bash
DATABASE_URL=postgresql://pokemon:pokemon123@postgres:5432/poketactix?sslmode=disable
JWT_SECRET=dev-secret-key-change-in-production-min-256-bits
JWT_EXPIRATION=24h
PORT=3000
ENV=development
CORS_ORIGINS=http://localhost:5173
DB_MAX_CONNECTIONS=20
DB_MIN_CONNECTIONS=2
DB_IDLE_TIMEOUT=300
DB_MAX_LIFETIME=1800
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
POKEAPI_BASE_URL=https://pokeapi.co/api/v2
POKEAPI_TIMEOUT=10s
```

### Frontend

```bash
VITE_API_URL=http://localhost:3000
```

## Development Workflow

### 1. Make Code Changes

- **Backend**: Edit files in the root directory. Air will automatically rebuild.
- **Frontend**: Edit files in `frontend/`. Vite will hot-reload changes.

### 2. View Logs

```bash
make logs
```

### 3. Test Changes

```bash
make test
```

### 4. Database Changes

If you modify the database schema:

1. Create a new migration file in `migrations/`
2. Run `make migrate` to apply changes
3. Restart backend: `make restart`

## Troubleshooting

### Services Won't Start

```bash
# Check if ports are already in use
lsof -i :3000
lsof -i :5173
lsof -i :5432

# Clean and restart
make clean
make up
```

### Database Connection Issues

```bash
# Check database health
docker-compose ps

# View database logs
make logs-db

# Restart database
docker-compose restart postgres
```

### Backend Build Errors

```bash
# View backend logs
make logs-backend

# Rebuild backend
docker-compose up -d --build backend
```

### Frontend Not Loading

```bash
# View frontend logs
make logs-frontend

# Rebuild frontend
docker-compose up -d --build frontend
```

### Clear All Data

```bash
# This will delete all database data
make clean
make up
make migrate
```

## Live Reload

### Backend (Air)

Air watches for file changes and automatically rebuilds the Go application. Configuration is in `.air.toml`.

**Excluded directories**: assets, tmp, vendor, testdata, frontend, migrations, .git, .github, .kiro

### Frontend (Vite)

Vite provides instant hot module replacement (HMR) for React components.

## Database Access

### Using psql

```bash
make db-shell
```

This opens a PostgreSQL shell where you can run SQL queries:

```sql
-- List all tables
\dt

-- View users
SELECT * FROM users;

-- View player cards
SELECT * FROM player_cards;

-- Exit
\q
```

### Using a GUI Client

Connect with your favorite PostgreSQL client:

- **Host**: localhost
- **Port**: 5432
- **Database**: poketactix
- **User**: pokemon
- **Password**: pokemon123

## Production Differences

This Docker setup is for **development only**. For production:

1. Use production Dockerfile (not Dockerfile.dev)
2. Change JWT_SECRET to a secure random value
3. Enable SSL for database connections
4. Use environment-specific CORS_ORIGINS
5. Deploy to Railway (backend) and Vercel/Netlify (frontend)
6. Use Neon for PostgreSQL database

See deployment documentation for production setup.

## Network

All services run on the `poketactix-network` bridge network, allowing them to communicate using service names:

- Backend connects to database using `postgres:5432`
- Frontend connects to backend using `http://localhost:3000` (from host)

## Volumes

- `postgres_data`: Persists database data between restarts
- Named volumes for node_modules and Go vendor directories

## Tips

1. **First time setup**: Run `make up && make migrate` to get started
2. **Daily development**: Just run `make up` - data persists
3. **Fresh start**: Run `make clean && make up && make migrate`
4. **View all commands**: Run `make help`
5. **Backend shell access**: Run `make backend-shell`

## Next Steps

After setting up Docker:

1. Review the database schema in `migrations/`
2. Check the API documentation (coming soon)
3. Start implementing features from the task list
4. Run tests regularly with `make test`

For questions or issues, refer to the main README or project documentation.
