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
├── frontend/       # Static HTML, CSS, JS + assets (deployed via GH Pages)
├── server/         # Go Fiber backend serving REST API (deployed on Railway)
├── pokemon/        # PokeAPI fetch + raw data handling
├── game/           # Core battle logic, card/move handling
├── go.mod
└── README.md       # ← You're looking at it!
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

### CLI version 1.0.0 (alpha) (coming soon)

1. Download the CLI from https://github.com/IfrunRuhin12/PokeTacTix/releases
2. Run `go run main.go` and follow the prompts
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
