# Scripts

Utility scripts for PokeTacTix development.

## 🚀 Main Scripts

### `docker-dev.sh`
Start everything with Docker (database, backend, frontend).

```bash
./scripts/docker-dev.sh
# Or
make dev
```

### `local-dev.sh`
Setup guide for local development (without Docker).

```bash
./scripts/local-dev.sh
# Or
make local
```

### `migrate.sh`
Database migration management.

```bash
./scripts/migrate.sh up      # Run migrations
./scripts/migrate.sh down    # Rollback migrations
./scripts/migrate.sh reset   # Reset database

# Or use make commands
make migrate
make migrate-down
make migrate-reset
```

### `setup_db.sh`
Legacy database setup script (use docker-dev.sh instead).

### `test_swagger.sh`
Test Swagger API documentation.

```bash
./scripts/test_swagger.sh
```

### `test-docker.sh`
Verify Docker setup.

```bash
./scripts/test-docker.sh
```

---

## 📝 Usage

All scripts should be run from the project root directory:

```bash
# Good
./scripts/docker-dev.sh

# Bad (don't do this)
cd scripts && ./docker-dev.sh
```

---

## 💡 Recommended Workflow

1. **First time setup**: `make dev`
2. **Daily development**: `make dev`
3. **Run migrations**: `make migrate`
4. **Stop services**: `make stop`

See the [Makefile](../Makefile) for all available commands.
