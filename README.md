# Pokemon Card Game - PokeTacTix (Golang)

A command-line Pokémon card game built in Go. Fetches real Pokémon data from the PokéAPI, builds game-ready cards, and displays them in the terminal. Now features a full command system, player-vs-AI battles, advanced battle logic, and improved command handling!

---

## 🚀 Features (Current State)
- Interactive command system: type commands to play and control the game
- Fetches Pokémon data from the PokéAPI by name
- Builds a game card with:
  - Name, HP (with bonus), Stamina (based on HP, will be based on Speed in future), Attack, Defense
  - Types, Sprite URL
  - Up to 4 real damaging moves (with power, type, stamina cost)
- Pretty-prints the card in the terminal for easy testing
- **Battle system:**
  - Player and AI each get 5 random Pokémon cards (with rare legendary/mythical odds)
  - Choose your active Pokémon for each round
  - Turn-based battle: attack, defend, pass, switch, surrender, and more
  - Stamina and HP management, type effectiveness (coming soon), and more
  - **Sacrifice mechanic:** Sacrifice HP to regain stamina, up to 3 times per round, with increasing HP cost and decreasing stamina gain
  - **Switch mechanic:** Switch Pokémon at the start of a round if you won the previous round
  - **Pass mechanic:** Skip your turn; AI can also pass under certain conditions
  - **Surrender and Surrender All:** Surrender a round or the entire battle
  - **AI logic:** AI makes strategic decisions about attacking, defending, sacrificing, passing, and switching
  - **Damage calculation:** Attack stat influences the probability of dealing higher damage; legendary/mythical Pokémon will be even stronger in future updates
  - **View your current card at any time with the 'card' command**
  - **Exit logic:** You cannot exit during a battle; you must finish or use 'surrender all' first
  - **Post-battle state reset:** After a battle ends, you can use all normal commands again
- Modular codebase for easy extension
- Improved command handling and in-battle command system
- **Assets directory** for future UI and type logic reference

---

## 🗂️ Project Structure
```
PokeTacTix/
├── go.mod
├── main.go                      # Entry point: minimal, runs command loop
├── game/
│   ├── command.go               # Command parsing, user input, game state
│   ├── commandHandler.go        # Central command dispatcher
│   ├── commandInGame.go         # In-battle command logic
│   ├── logic.go                 # Battle/turn/round logic, state transitions
│   ├── player.go                # Player struct, deck logic
│   ├── card.go                  # Pretty-print helper for cards
│   ├── utils.go                 # Helpers: random deck, type multipliers, damage percent, etc.
│   └── welcome.go               # Welcome message
├── pokemon/
│   ├── fetch.go                 # Fetches API data, returns usable moves
│   ├── types.go                 # Raw API structs, Card struct
│   ├── cardbuilder.go           # Builds Card from Pokemon+Moves (shared logic)
│   └── display.go               # (for raw info/debug)
├── assets/
│   └── type-logic.jpg           # Type effectiveness chart (for reference/UI)
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
   - `exit`     — Exit the game (only allowed outside of battle)
   - `command --in-battle` — Show in-battle commands
4. **To battle:**
   - Type `battle` and follow the prompts (enter your name, confirm to start)
   - Choose your active Pokémon for the round (by number)
   - AI will choose one randomly
   - The game will prompt you for your move each turn (e.g., `attack`, `defend`, `sacrifice`, `pass`, `switch`, `surrender`, `surrender all`)
   - Use `card all` to view your cards at any time during battle
   - Use `card` to view your currently selected card
   - Use `switch` at the start of a round if you won the previous round
   - Use `sacrifice` to trade HP for stamina (up to 3 times per round)
   - Use `pass` to skip your turn (AI may also pass)
   - Use `surrender` to lose the current round, or `surrender all` to lose the entire battle
   - You cannot exit during a battle; you must finish or use `surrender all` first

---

## ⚔️ In-Battle Commands
- `card all`   — Show all your cards
- `card`       — Show your currently selected card
- `attack`     — Choose a move to attack with
- `defend`     — Defend against an attack
- `choose`     — Choose a card to play (after KO or surrender)
- `switch`     — Switch your active Pokémon at the start of a round (if you won the previous round)
- `sacrifice`  — Sacrifice HP to regain stamina (up to 3 times per round)
- `pass`       — Skip your turn
- `surrender`  — Surrender the current round
- `surrender all` — Surrender the entire battle
- `command --in-battle` — Show this command list
- `exit`       — Exit the game (only allowed outside of battle)

---

## 🧭 Roadmap / Next Steps
- [x] Add deck generation (5 random Pokémon per player/AI, with legendary/mythical odds)
- [x] Implement full battle logic (turns, attack/defend/pass/sacrifice/switch, stamina, HP, type effectiveness coming soon)
- [x] Add player/AI structs and state
- [x] CLI menu for viewing cards, starting battles, and more
- [x] Full match system (win/lose, surrender, next battle, post-battle state reset)
- [x] Advanced AI logic for attacking, defending, passing, and sacrificing
- [x] Damage calculation based on attack stat, with probability tables
- [x] In-battle and out-of-battle command separation
- [x] Prevent exit during battle
- [x] View current card at any time with 'card' command
- [x] **Type Multiplier/Weakness:** Implement real Pokémon type effectiveness (planned)
- [x] **Legendary/Mythical Power:** Legendary/mythical Pokémon will do 2x damage to normal Pokémon (planned)
- [x] **Stamina Based on Speed:** Stamina will be based on Pokémon speed stat from PokéAPI (planned)
- [ ] **1v1 Battles:** Add option for short 1v1 battles against AI (planned)
- [ ] Build a Go web server using Fiber to serve the game over HTTP
- [ ] Create a modern frontend using Alpine.js to make the cards look cooler and give the game life
- [ ] Make cards look awesome with custom styles and dynamic effects
- [ ] Use assets/type-logic.jpg for type effectiveness UI or documentation
- [ ] More polish, bug fixes, and UI improvements

---

## 🏗️ Architecture
- **main.go:** Only runs the command loop and delegates all logic
- **game/command.go:** Handles all command parsing, user input, and game state
- **game/commandHandler.go:** Central command dispatcher
- **game/commandInGame.go:** Handles in-battle commands and logic
- **game/logic.go:** Contains all battle, round, and turn logic
- **game/utils.go:** Helpers for random deck, type multipliers, damage percent, etc.
- **pokemon/cardbuilder.go:** Shared logic for building a Card from API data
- **pokemon/fetch.go:** Fetches Pokémon and move data from PokéAPI
- **assets/type-logic.jpg:** Type effectiveness chart (for reference/UI)

---

## 🤝 Contributing
PRs and suggestions welcome! This is a learning/fun project.

---

## 📦 Credits
- [PokéAPI](https://pokeapi.co/) for all Pokémon data
- Built with Go

---

## © 2024 Ifran Kader Ruhin ([ifrunruhin12](https://github.com/ifrunruhin12))
This project is licensed under the Creative Commons Attribution-NonCommercial 4.0 International License (CC BY-NC 4.0).
It is free to use for personal learning and non-commercial purposes only. Commercial use is strictly prohibited.
See the LICENSE file for details.

**Note:** This project is currently open source, but may become closed source in the future at the author's discretion.
