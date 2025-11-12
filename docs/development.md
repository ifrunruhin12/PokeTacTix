# Development Guide

This guide explains how to set up and run PokeTacTix for development.

## 🎯 Choose Your Setup

### 🐳 Docker (Recommended)

**Best for:**
- Quick setup
- Consistent environment
- No local installation needed
- Testing and development

**Pros:**
- ✅ One command to start everything
- ✅ No need to install PostgreSQL locally
- ✅ Consistent across all machines
- ✅ Easy to clean up and reset

**Cons:**
- ❌ Requires Docker Desktop
- ❌ Slightly slower than native (minimal difference)

### 💻 Local Development

**Best for:**
- Faster iteration
- Debugging with IDE
- Lower resource usage
- Production-like environment

**Pros:**
- ✅ Faster rebuild times
- ✅ Better IDE integration
- ✅ Lower memory usage
- ✅ Direct access to all tools

**Cons:**
- ❌ Requires manual setup
- ❌ Need to install PostgreSQL, Go, Node.js
- ❌ Environment differences between machines

---

## 🚀 Quick Start

### Docker Setup (Easiest)

```bash
make dev
```

**Access:**
- Frontend: http://localhost:5173
- Backend: http://localhost:3000
- API Docs: http://localhost:3000/api/docs

**Management:**
```bash
make logs      # View logs
make stop      # Stop services
make restart   # Restart
make clean     # Remove all data
```

### Local Setup

```bash
make local
```

Then start services in separate terminals:

**Terminal 1 (Backend):**
```bash
export DATABASE_URL="postgresql://pokemon:pokemon123@localhost:5432/poketactix?sslmode=disable"
go run cmd/api/main.go
```

**Terminal 2 (Frontend):**
```bash
cd frontend
npm run dev
```

---

## 📁 Project Structure

```
PokeTacTix/
├── cmd/api/              # Application entry point
├── internal/             # Private application code
│   ├── auth/            # Authentication
│   ├── battle/          # Battle system
│   ├── cards/           # Card management
│   ├── pokemon/         # Pokemon data
│   ├── shop/            # Shop system
│   ├── stats/           # Statistics
│   └── database/        # Database layer
├── frontend/            # React frontend
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── pages/       # Page components
│   │   ├── services/    # API services
│   │   └── contexts/    # React contexts
│   └── public/          # Static assets
├── pkg/                 # Shared utilities
├── docs/                # Documentation
└── scripts/             # Utility scripts
```

---

## 🔧 Configuration

### Environment Variables

**Backend (.env or export):**
```bash
DATABASE_URL=postgresql://pokemon:pokemon123@localhost:5432/poketactix?sslmode=disable
JWT_SECRET=dev-secret-key-change-in-production-min-256-bits
JWT_EXPIRATION=24h
PORT=3000
ENV=development
CORS_ORIGINS=http://localhost:5173
```

**Frontend (frontend/.env):**
```bash
VITE_API_URL=http://localhost:3000
```

### Database Connection

**Docker:**
- Host: `localhost`
- Port: `5432`
- Database: `poketactix`
- User: `pokemon`
- Password: `pokemon123`

**Connection String:**
```
postgresql://pokemon:pokemon123@localhost:5432/poketactix?sslmode=disable
```

---

## 🛠️ Development Workflow

### Docker Workflow

```bash
# Start development
make dev

# View logs (in another terminal)
make logs

# Make code changes (auto-reload enabled)
# - Backend: Air watches for Go file changes
# - Frontend: Vite watches for React changes

# Restart if needed
make restart

# Stop when done
make stop
```

### Local Workflow

```bash
# Start PostgreSQL
sudo systemctl start postgresql  # Linux
brew services start postgresql   # macOS

# Start backend (Terminal 1)
export DATABASE_URL="postgresql://pokemon:pokemon123@localhost:5432/poketactix?sslmode=disable"
go run cmd/api/main.go

# Start frontend (Terminal 2)
cd frontend
npm run dev

# Make changes (both have hot reload)

# Stop services
# Ctrl+C in each terminal
```

---

## 🧪 Testing

### Run Tests

**Docker:**
```bash
make test
```

**Local:**
```bash
go test ./...
```

### Test API with curl

```bash
# Register
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"player1","email":"player1@example.com","password":"Test123!@#"}'

# Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"player1","password":"Test123!@#"}'

# Get profile (replace TOKEN)
curl -X GET http://localhost:3000/api/profile/stats \
  -H "Authorization: Bearer TOKEN"
```

### Interactive API Testing

Visit http://localhost:3000/api/docs for Swagger UI

---

## 🗄️ Database Management

### Docker

```bash
make db-shell          # Open psql shell
make db-logs           # View database logs
make db-only           # Start only database
```

### Local

```bash
# Connect to database
psql -U pokemon -d poketactix

# View tables
\dt

# Query data
SELECT * FROM users;

# Exit
\q
```

### Reset Database

**Docker:**
```bash
make clean  # Removes all data
make dev    # Start fresh
```

**Local:**
```bash
sudo -u postgres psql
DROP DATABASE poketactix;
CREATE DATABASE poketactix OWNER pokemon;
\q
```

---

## 📊 Monitoring

### View Logs

**Docker:**
```bash
make logs              # All services
make backend-logs      # Backend only
make frontend-logs     # Frontend only
make db-logs           # Database only
```

**Local:**
- Backend: Check terminal output
- Frontend: Check terminal output
- Database: `sudo journalctl -u postgresql`

### Check Service Status

**Docker:**
```bash
make status
```

**Local:**
```bash
# Check if services are running
lsof -i :3000  # Backend
lsof -i :5173  # Frontend
lsof -i :5432  # Database
```

---

## 🐛 Troubleshooting

### Docker Issues

**Services won't start:**
```bash
make clean
make dev
```

**Port already in use:**
```bash
# Stop local PostgreSQL
sudo systemctl stop postgresql

# Or change ports in docker-compose.yml
```

**View detailed logs:**
```bash
docker-compose logs backend
docker-compose logs frontend
docker-compose logs postgres
```

### Local Issues

**Database connection refused:**
```bash
# Start PostgreSQL
sudo systemctl start postgresql

# Check if running
sudo systemctl status postgresql
```

**Port 3000 already in use:**
```bash
lsof -i :3000
kill -9 <PID>
```

**Frontend won't start:**
```bash
cd frontend
rm -rf node_modules
npm install
npm run dev
```

---

## 💡 Tips

1. **Use Docker for quick testing** - It's the fastest way to get started
2. **Use local for active development** - Better IDE integration and faster rebuilds
3. **Check logs first** - Most issues are visible in logs
4. **Use make commands** - They're shortcuts for common tasks
5. **Keep Docker Desktop running** - Required for Docker setup

---

## 🆘 Getting Help

1. Check logs: `make logs`
2. Check status: `make status`
3. Try restarting: `make restart`
4. Clean and rebuild: `make clean && make dev`
5. Check documentation in `docs/` folder

---

**Happy Coding! 🚀**
