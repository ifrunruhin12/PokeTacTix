# PokeTacTix Quick Start Guide

Get up and running with PokeTacTix in 3 simple steps!

## 🚀 Quick Setup (5 minutes)

### Step 1: Set Up Database

```bash
./scripts/setup_db.sh
```

This creates a PostgreSQL database in Docker with all tables configured.

### Step 2: Start Backend API

```bash
go run cmd/api/main.go
```

API will be available at `http://localhost:3000`

### Step 3: Start Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend will be available at `http://localhost:5173`

## 🎮 Test It Out

1. Open `http://localhost:5173` in your browser
2. Click "Get Started"
3. Register a new account
4. Start exploring!

## 📚 Common Commands

### Database

```bash
# View database logs
docker logs poketactix-db

# Connect to database
docker exec -it poketactix-db psql -U pokemon -d poketactix

# Stop database
docker stop poketactix-db

# Start database
docker start poketactix-db

# Reset database (WARNING: deletes all data)
docker rm -f poketactix-db && ./scripts/setup_db.sh
```

### Backend

```bash
# Run API server
go run cmd/api/main.go

# Run tests
go test ./...

# Build binary
go build -o poketactix cmd/api/main.go

# View API docs
open http://localhost:3000/api/docs/
```

### Frontend

```bash
# Install dependencies
npm install --prefix frontend

# Start dev server
npm run dev --prefix frontend

# Build for production
npm run build --prefix frontend

# Preview production build
npm run preview --prefix frontend
```

## 🧪 Test API with curl

```bash
# Register
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"player1","email":"player1@example.com","password":"Test123!@#"}'

# Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"player1","password":"Test123!@#"}'

# Get stats (replace TOKEN with your JWT)
curl -X GET http://localhost:3000/api/profile/stats \
  -H "Authorization: Bearer TOKEN"
```

## 🔧 Configuration

### Environment Variables (.env)

```bash
# Database
DATABASE_URL=postgresql://pokemon:pokemon123@localhost:5432/poketactix?sslmode=disable

# JWT
JWT_SECRET=test-secret-key-for-development-only-min-256-bits
JWT_EXPIRATION=24h

# Server
PORT=3000
ENV=development

# CORS
CORS_ORIGINS=http://localhost:5173,http://localhost:3000
```

### Frontend Environment (frontend/.env)

```bash
VITE_API_URL=http://localhost:3000
```

## 📖 Documentation

- **Full Testing Guide**: [docs/TESTING_GUIDE.md](docs/TESTING_GUIDE.md)
- **API Documentation**: http://localhost:3000/api/docs/
- **Swagger Spec**: [docs/swagger.yaml](docs/swagger.yaml)

## ❓ Troubleshooting

### "Database connection failed"
```bash
docker start poketactix-db
```

### "Port 3000 already in use"
```bash
# Find and kill the process
lsof -i :3000
kill -9 <PID>
```

### "CORS error in browser"
Check that CORS_ORIGINS in .env includes `http://localhost:5173`

### "Token expired"
Login again to get a new token (tokens expire after 24 hours)

## 🎯 What's Next?

- Explore the battle system
- Purchase Pokemon from the shop
- Build your deck
- Check out achievements
- View your stats and battle history

## 💡 Tips

- Use the Swagger UI at `/api/docs/` to explore all API endpoints
- Check browser console for frontend errors
- Use `docker logs poketactix-db` to debug database issues
- JWT tokens are stored in localStorage (check Application tab in DevTools)

## 🆘 Need Help?

- Check [docs/TESTING_GUIDE.md](docs/TESTING_GUIDE.md) for detailed instructions
- Review API documentation at http://localhost:3000/api/docs/
- Check the logs for error messages

---

**Happy Gaming! 🎮✨**
