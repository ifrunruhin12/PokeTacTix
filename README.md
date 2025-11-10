# ◓ PokeTacTix

A turn-based Pokémon card battle game 💥  
**Frontend** (Alpine/Tailwind) hosted on **GitHub Pages**  
**Backend** (Go + Fiber API) live on **Railway**

---

## 🧭 Live Demo

- **Website**: https://ifrunruhin12.github.io/PokeTacTix/
Hopefully will have it's own domain soon!

---

## ⚙️ Project Structure

```
PokeTacTix/
├── cmd/
│   └── api/              # Application entry point
├── internal/             # Private application code
│   ├── auth/            # Authentication domain (handlers, service, JWT, middleware)
│   ├── battle/          # Battle system for web API (handlers, session management)
│   ├── cards/           # Card management domain (handlers, service, repository)
│   ├── pokemon/         # Pokemon fetching and card building
│   └── database/        # Database connection, models, and migrations
├── pkg/                 # Shared utilities
│   ├── config/          # Configuration management
│   └── logger/          # Structured logging
├── game/                # Core battle logic (CLI only)
│   ├── commands/        # CLI command handlers
│   ├── core/            # Battle engine and game logic
│   ├── models/          # Game state models
│   └── utils/           # Game utilities (type chart, etc.)
├── frontend/            # Static HTML, CSS, JS + assets (deployed via GH Pages)
├── main.go              # CLI entry point
├── go.mod
└── README.md            # ← You're looking at it!
```

### Architecture

The project follows a **feature-based architecture** with clear separation of concerns:

- **cmd/api**: Main application entry point with dependency injection
- **internal/**: Private application code organized by domain (auth, battle, cards, pokemon)
- **pkg/**: Shared packages that could be used by other projects
- **game/**: CLI-only battle logic (separate from web API)
- **frontend/**: Static web UI (HTML, CSS, JS)
- **Database layer**: Centralized in internal/database with domain-specific repositories
- **Clean dependencies**: Each domain is self-contained with its own models, handlers, and business logic

---

## ✅ How to Play

### Web version 1.0.0 (alpha)
1. **Frontend**:  
   Browse to the GitHub Pages URL, which loads the card battlefield.

2. **Search Pokémon**:  
   Enter a name on the home page after clicking search — it fetches a styled card.

3. **1v1 Arena**:  
   Head to the battle arena, choose **1v1** mode, and battle your Pokémon against AI:
   - Select **Attack**, **Defend**, **Sacrifice**, etc.
   - Buttons represent moves with type-based colors.
   - Battle log shows turn progression and damage data.

The frontend uses JS fetch calls to your live backend for everything — no page reloads once loaded.

### CLI version 1.0.0 (alpha)

1. Download the CLI from https://github.com/IfrunRuhin12/PokeTacTix/releases
   - Linux: poketactix_linux_amd64, poketactix_linux_arm64
   - Windows: poketactix_windows_amd64.exe, poketactix_windows_arm64.exe
   - macOS: poketactix_darwin_amd64, poketactix_darwin_arm64
2. On Linux/macOS: `chmod +x ./poketactix_*`
3. Run the binary:
   - Linux/macOS: `./poketactix_linux_amd64` (or your arch file)
   - Windows: double-click or `poketactix_windows_amd64.exe` in cmd/PowerShell

---

## 📚 API Documentation

The PokeTacTix API is fully documented with **OpenAPI 3.0** (Swagger) specification.

### Interactive Documentation

Once the server is running, access the interactive Swagger UI at:
```
http://localhost:3000/api/docs
```

The Swagger UI provides:
- **Interactive API testing** - Try endpoints directly from your browser
- **Request/response examples** - See exactly what to send and expect
- **Schema definitions** - Understand all data models
- **Authentication testing** - Test JWT authentication flows

### API Endpoints

The API includes comprehensive endpoints for:
- **Authentication** (`/api/auth/*`) - Register, login, session management
- **Cards** (`/api/cards/*`) - Collection and deck management
- **Battle** (`/api/battle/*`) - 1v1 and 5v5 battle operations
- **Shop** (`/api/shop/*`) - Pokemon card purchases
- **Profile** (`/api/profile/*`) - Statistics, history, and achievements

### Quick Start

1. **Register a new account:**
   ```bash
   curl -X POST http://localhost:3000/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username":"trainer","email":"trainer@pokemon.com","password":"SecurePass123!"}'
   ```

2. **Login and get JWT token:**
   ```bash
   curl -X POST http://localhost:3000/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"trainer","password":"SecurePass123!"}'
   ```

3. **Use the token for authenticated requests:**
   ```bash
   curl -X GET http://localhost:3000/api/cards \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

For detailed documentation, see [docs/API_DOCUMENTATION.md](docs/API_DOCUMENTATION.md)

---

## 🚀 Latest Features (June, 2025)

- **Modern 1v1 Battle Arena**: Beautiful, card-based UI with colored type badges, responsive layout, and smooth turn flow.
- **Full Backend Logic**: All game rules (turns, moves, AI, sacrifice, surrender, damage, type multipliers) handled by Go backend for perfect consistency.
- **Battle Log**: Grouped, turn-by-turn log matching the CLI, with move names and results.
- **Surrender & Draw**: Surrender ends the battle instantly; draws are detected and shown.
- **Result Banner**: Shows "You won!", "You lost", or "Draw!" based on the true outcome.
- **5v5 Mode**: UI placeholder/under construction (coming soon).
- **Frontend/Backend Sync**: All rules, turn order, and log formatting match between web and CLI.

---

## 🔮 What’s Next

- 5v5 **team battles** (full implementation)
- Account system (login, persistent stats)
- Multiplayer (PvP, matchmaking, live battles)
- In-game store (buy/sell cards, cosmetics)
- Card reveal/hide mechanics (fog of war, secret moves)
- Deck building and export
- More polish, animations, and accessibility improvements

---

## 👤 About

**PokeTacTix** is built by **Ifrun Kader Ruhin**, a student and dev leveling up full-stack real-time strategy games in Golang.

Contributions are welcome — but watch my README evolve as the app does 😉  
Expect new features and fresh rewrites soon.

---

## 📄 License

This project is licensed under the **Creative Commons Attribution-NonCommercial 4.0 International (CC BY-NC 4.0)** license.  
You're free to remix, adapt, and build upon it — just give credit and keep it non-commercial.  
Full license text: https://creativecommons.org/licenses/by-nc/4.0/

README.md
Displaying README.md.
