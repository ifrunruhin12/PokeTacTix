# Pokemon Card Game - PokeTacTix (Golang)

A command-line Pokémon card game built in Go. Fetches real Pokémon data from the PokéAPI, builds game-ready cards, and displays them in the terminal. Now features a full command system, player-vs-AI battles, and modular game logic!

---

## 🚀 Features (Current State)
- Interactive command system: type commands to play and control the game
- Fetches Pokémon data from the PokéAPI by name
- Builds a game card with:
  - Name, HP (with bonus), Stamina, Attack, Defense
  - Types, Sprite URL
  - Up to 4 real damaging moves (with power, type, stamina cost)
- Pretty-prints the card in the terminal for easy testing
- **Battle system:**
  - Player and AI each get 5 random Pokémon cards (with very rare legendary/mythical odds)
  - Choose your active Pokémon for each round
  - Turn-based battle: attack, defend, surrender, and more
  - Stamina and HP management, type effectiveness, and more (in progress)
- Modular codebase for easy extension

---

## 🗂️ Project Structure
```
PokeTacTix/
├── go.mod
├── main.go                      # Entry point: minimal, runs command loop
├── game/
│   ├── command.go               # Command parsing, user input, game state
│   ├── logic.go                 # Battle/turn/round logic, state transitions
│   ├── player.go                # Player struct, deck logic
│   ├── card.go                  # Pretty-print helper for cards
│   ├── utils.go                 # Helpers: random deck, type multipliers, etc.
│   └── welcome.go               # Welcome message
├── pokemon/
│   ├── fetch.go                 # Fetches API data, returns usable moves
│   ├── types.go                 # Raw API structs, Card struct
│   ├── cardbuilder.go           # Builds Card from Pokemon+Moves (shared logic)
│   └── display.go               # (for raw info/debug)
└── README.md
```

---

## 🛠️ Usage
1. **Run the program:**
   ```sh
   go run main.go
   ```
2. **You will see a welcome message and a prompt.**
3. **Type `command` to see all available commands:**
   - `search`   — Search for a Pokémon by name and see its card
   - `battle`   — Start a 5v5 match against AI (get 5 random cards)
   - `card all` — Show all your cards (only after starting battle)
   - `exit`     — Exit the game
4. **To battle:**
   - Type `battle` and follow the prompts (enter your name, confirm to start)
   - Choose your active Pokémon for the round (by number or name)
   - AI will choose one randomly
   - The game will prompt you for your move each turn (e.g., `attack 1`, `defend`, `surrender`)
   - Use `card all` to view your cards at any time during battle

---

## ⚔️ In-Battle Commands
- `attack N`   — Use move N (1-4) of your active Pokémon
- `defend`     — Defend this turn (reduces or blocks damage, drains stamina)
- `surrender`  — Surrender the round
- `switch`     — (After a round ends) Switch to another Pokémon for the next round
- `card all`   — View all your cards
- `exit`       — Exit the game

---

## 🧭 Roadmap / Next Steps
- [x] Add deck generation (5 random Pokémon per player/AI, with legendary/mythical odds)
- [x] Implement battle logic (turns, attack/defend, stamina, HP, type effectiveness)
- [x] Add player/AI structs and state
- [x] CLI menu for viewing cards, starting battles, and more
- [ ] Full match system (win/lose, surrender, next battle)
- [ ] **Web Version:** Build a Go web server using [Fiber](https://gofiber.io/) to serve the game over HTTP
- [ ] **Frontend:** Use [Alpine.js](https://alpinejs.dev/) to create a modern, interactive UI for cards and gameplay
- [ ] Make cards look awesome with custom styles and dynamic effects

---

## 🏗️ Architecture
- **main.go:** Only runs the command loop and delegates all logic
- **game/command.go:** Handles all command parsing, user input, and game state
- **game/logic.go:** Contains all battle, round, and turn logic
- **game/utils.go:** Helpers for random deck, type multipliers, dice rolls, etc.
- **pokemon/cardbuilder.go:** Shared logic for building a Card from API data
- **pokemon/fetch.go:** Fetches Pokémon and move data from PokéAPI

---

## 🤝 Contributing
PRs and suggestions welcome! This is a learning/fun project.

---

## 📦 Credits
- [PokéAPI](https://pokeapi.co/) for all Pokémon data
- Built with Go
