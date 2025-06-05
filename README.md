# Pokemon CLI Card Game (Golang)

A command-line Pokémon card game built in Go. Fetches real Pokémon data from the PokéAPI, builds game-ready cards, and displays them in the terminal. Future updates will add full battle logic and gameplay!

---

## 🚀 Features (Current State)
- Fetches Pokémon data from the PokéAPI by name.
- Builds a game card with:
  - Name, HP, Stamina, Attack, Defense
  - Types, Sprite URL
  - Up to 3 real damaging moves (with power, type, stamina cost)
- Pretty-prints the card in the terminal for easy testing.

---

## 🗂️ Project Structure
```
pokemon-cli/
├── go.mod
├── main.go                      # Entry point: user input, card display
├── game/
│   ├── card.go                  # Converts pokemon.Pokemon → game.Card, pretty-print helper
│   └── ...                      # (logic.go, player.go, utils.go coming soon)
├── pokemon/
│   ├── fetch.go                 # Fetches API data, returns usable moves
│   ├── types.go                 # Raw API structs, Card struct
│   └── display.go               # (for raw info/debug)
└── README.md
```

---

## 🛠️ Usage
1. **Run the program:**
   ```sh
   go run main.go
   ```
2. **Enter a Pokémon name** (e.g. `pikachu`, `bulbasaur`, `charizard`).
3. **See the card!**
   - All stats, types, sprite URL, and real moves are shown in a pretty format.

---

## 🧭 Roadmap / Next Steps
- [ ] Add deck generation (5 random Pokémon per player/AI, with legendary odds)
- [ ] Implement battle logic (turns, attack/defend, stamina, HP, type effectiveness)
- [ ] Add player/AI structs and state
- [ ] CLI menu for viewing cards, starting battles, and more
- [ ] Full match system (win/lose, surrender, next battle)
- [ ] **Web Version:** Build a Go web server using [Fiber](https://gofiber.io/) to serve the game over HTTP
- [ ] **Frontend:** Use [Alpine.js](https://alpinejs.dev/) to create a modern, interactive UI for cards and gameplay
- [ ] Make cards look awesome with custom styles and dynamic effects

---

## 🤝 Contributing
PRs and suggestions welcome! This is a learning/fun project.

---

## 📦 Credits
- [PokéAPI](https://pokeapi.co/) for all Pokémon data
- Built with Go
