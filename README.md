
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

---

## 🔧 Setup & Development

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

## 🔮 What’s Next

- Add 5v5 **team battles**
- Create battle logs and session history
- Future-proof with WebSockets for live play!
- Polish UI/UX, theme based on card types
- Add custom deck building + save/export features

---

## 👤 About

**PokeTacTix** is built by **Ifrun Kader Ruhin**, a student and dev leveling up full-stack real-time strategy games in Golang.

Contributions are welcome — but watch my README evolve as the app does 😉  
Expect new features and fresh rewrites soon.

---

## 📄 License

This project is licensed under the **Creative Commons Attribution-NonCommercial 4.0 International (CC BY-NC 4.0)** license.  
You’re free to remix, adapt, and build upon it — just give credit and keep it non-commercial.  
Full license text: https://creativecommons.org/licenses/by-nc/4.0/
