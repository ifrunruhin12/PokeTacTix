# PokeTacTix Server

This is the Go Fiber backend for the PokeTacTix Pokémon card battle game.

---

## 🚀 Current Features

- **REST API** for all game actions (start battle, make move, get state)
- **1v1 Battle Mode**: Full turn-based logic, random Pokémon, player vs AI
- **All Game Logic in Go**: Turns, moves, stamina, type multipliers, legendary bonuses, etc.
- **Battle Log**: Grouped, turn-by-turn log matching CLI
- **Surrender & Draw**: Surrender ends the battle instantly; draws are detected
- **Result Reporting**: API returns winner (player, ai, draw) for frontend display
- **Frontend/Backend Sync**: All rules and flow match between web and CLI

---


## 🔧 Setup & Development (this section will be removed in the future)

### Frontend
- `frontend/`: HTML + Alpine.js + Tailwind CSS  
- Deploys to GitHub Pages via `gh-pages` branch  
- **Make sure**: `script.js` points to your Railway domain

### Backend
- `server/`: Go + Fiber REST API  
- Endpoints:
  - `GET /pokemon?name=<name>` → returns JSON card
  - `POST /battle/start` → starts session
  - `POST /battle/turn` → progress turn via JSON
- Deploy on Railway: setup build = `go build -o main ./server`, run = `./main`
- **Enable CORS** (`app.Use(cors.New())`) so frontend can fetch

### Shared Logic
- `pokemon/`: Fetch and parse PokeAPI
- `game/`: Battle engine, moves, card definitions, turn logic
- Two frontends supported: CLI (via `main.go`) and Web via `server/main.go` and Alpine.js

---

## 🛠️ Planned Features

- **5v5 Team Battles**: Full multi-round, multi-card support
- **Account System**: User login, persistent stats, profiles
- **Multiplayer**: PvP, matchmaking, live battles (WebSocket support)
- **In-Game Store**: Buy/sell cards, cosmetics, upgrades
- **Card Reveal/Hide**: Fog of war, secret moves, hidden info
- **Deck Building**: Save, export, and share custom decks
- **Session History**: Battle logs, replays, stats
- **More polish**: Animations, accessibility, and performance improvements

---

## 📂 Structure

- `main.go` — Fiber server, API endpoints
- Uses `../game/` for all battle logic
- Uses `../pokemon/` for card data and helpers

---

**Contributions and feedback welcome!** 
