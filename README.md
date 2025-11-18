# ◓ PokeTacTix

A turn-based Pokémon card battle game 💥  
**Frontend** (React + Vite) hosted on **Netlify**  
**Backend** (Go + Fiber API) live on **Railway**  
**Database** (PostgreSQL) on **Neon**

---

## 🧭 Live Demo

- **Website**: [Coming Soon - Deploy with Netlify]
- **API**: [Coming Soon - Deploy with Railway]

## 🚀 Quick Deploy

Want to deploy your own instance? See [Quick Deploy Guide](./docs/QUICK_DEPLOY.md) (30 minutes)

For detailed deployment instructions, see [Deployment Guide](./docs/DEPLOYMENT_GUIDE.md)

---

## ⚙️ Project Structure

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
│   └── database/        # Database + migrations
├── frontend/            # React frontend
├── game/                # CLI battle logic
├── docs/                # Documentation
├── scripts/             # Utility scripts
└── pkg/                 # Shared utilities
```

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

Interactive API documentation available at: **http://localhost:3000/api/docs** (when running)

The API includes endpoints for authentication, cards, battles, shop, and player stats.

---

## 🚀 Getting Started

Want to run PokeTacTix locally? It's easy:

```bash
make dev
```

That's it! Everything starts automatically.

**See the docs for more:**
- **[Quick Start](docs/quick-start.md)** - Get running in 60 seconds
- **[Get Started Guide](docs/get-started.md)** - Step-by-step for beginners
- **[Development Guide](docs/development.md)** - Full development docs
- **[All Documentation](docs/)** - Complete documentation index

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
