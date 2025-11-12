# Quick Start

Get PokeTacTix running in 60 seconds.

## 🐳 Docker (Recommended)

```bash
make dev
```

That's it! Everything starts automatically.

**Access:**
- Frontend: http://localhost:5173
- Backend: http://localhost:3000
- API Docs: http://localhost:3000/api/docs

**Commands:**
```bash
make logs      # View logs
make stop      # Stop services
make restart   # Restart
make help      # Show all commands
```

## 💻 Local Setup

```bash
make local
```

Follow the prompts, then:

```bash
# Terminal 1
export DATABASE_URL="postgresql://pokemon:pokemon123@localhost:5432/poketactix?sslmode=disable"
go run cmd/api/main.go

# Terminal 2
cd frontend && npm run dev
```

## 🎮 Test It

1. Open http://localhost:5173
2. Register an account
3. Start battling!

## ❓ Issues?

```bash
make clean    # Remove all data
make dev      # Start fresh
```

---

**Need more help?** See [get-started.md](get-started.md) or [development.md](development.md)
