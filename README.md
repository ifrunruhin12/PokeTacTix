<div align="center">

# <img src="https://raw.githubusercontent.com/ifrunruhin12/PokeTacTix/refs/heads/main/frontend/public/assets/pokeball-throwing.gif"/> PokeTacTix

### *Strategic Turn-Based Pokémon Card Battles*

<img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
<img src="https://img.shields.io/badge/React-61DAFB?style=for-the-badge&logo=react&logoColor=black" alt="React" />
<img src="https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />

**[🎮 Play Now](https://poketactix.netlify.app/) • [📖 Docs](./docs/)**

</div>

---

## 🎯 What is PokeTacTix?

PokeTactix is a strategic game where you battle, collect, and level up Pokémon cards in strategic turn-based combat. Build your perfect deck, master type advantages, and climb the ranks.

<table>
<tr>
<td width="50%">

### 🎮 **Play Anywhere**
- 🌐 **Web** - Play in browser
- 📱 **Mobile** - Coming soon

</td>
<td width="50%">

### ⚔️ **Battle Modes**
- 🥊 **1v1** - Quick battles
- 👥 **5v5** - Team strategy
- 🤖 **AI** - Smart opponents

</td>
</tr>
</table>

---

## ✨ Features

<table>
<tr>
<td>

### 🃏 **Card Collection & Evolution**
Collect Pokémon, level them up through battles, and evolve them — by level or with evolution stones

</td>
<td>

### 🎲 **Strategic Combat**
Type advantages, stamina management, and tactical decisions

</td>
</tr>
<tr>
<td>

### 🛍️ **Item System**
Evolution stones and battle boosters, purchased in the shop and used from your inventory

</td>
<td>

### 🏆 **Progression**
Earn coins, track stats, unlock achievements

</td>
</tr>
<tr>
<td colspan="2">

### 🛒 **Shop**
Three categories — Pokémon cards, game tokens, and items — with daily inventory rotation and discounts

</td>
</tr>
</table>

---

## 🧬 Item & Evolution System

Pokémon evolve through **data-driven evolution rules** sourced from PokéAPI evolution chains — no Pokémon-specific logic in the code:

| Method | How it works | Example |
|--------|--------------|--------|
| **Level** | Applied automatically from battle rewards once the required level is reached | Charmander → Charmeleon at level 16 |
| **Item** | Player uses the required evolution stone from their deck | Pikachu + Thunder Stone → Raichu |
| **Friendship** | Friendship-based evolutions become level-up evolutions at a fixed level (the game has no friendship mechanic) | Pichu → Pikachu at level 20 |

**Items** live in a shared catalog (`internal/items`) with two types:

- **Evolution items** (Thunder/Fire/Water Stone) — consumed by item-based evolution, in the same transaction
- **Booster items** (Attack Booster: +5 attack for 3 battles; HP Booster: +10 HP for 4 battles) — activated from the shop, buff your whole deck, and tick down one battle at a time

Booster magnitudes and durations are data (in the `items.effect` JSONB column), so balancing is a seed change — no code required.

---

## 🚀 Quick Start

### 🌐 Web Version

Click the link and start playing [poketactix.netlify.app](https://poketactix.netlify.app/) 🎉

---

## 📚 Documentation

<table>
<tr>
<td align="center" width="25%">

### ⚡ [Quick Start](docs/quick-start.md)
*60 seconds to running*

</td>
<td align="center" width="25%">

### 🛠️ [Development](docs/development.md)
*Full dev setup*

</td>
<td align="center" width="25%">

### 🚀 [Get Started](docs/get-started.md)
*How to kick off things*

</td>
<td align="center" width="25%">

### 📡 [API Docs](http://localhost:3000/api/docs)
*Interactive API*

</td>
</tr>
</table>

---

## 🏗️ Tech Stack

<div align="center">

| Layer | Technology |
|-------|-----------|
| **Frontend** | React + Vite + Tailwind CSS |
| **Backend** | Go + Fiber |
| **Database** | PostgreSQL (Neon) |
| **Hosting** | Railway + Netlify |
| **Auth** | JWT |

</div>

---

## 🎯 Roadmap

- [x] 1v1 Battles
- [x] Card Collection
- [x] Shop System
- [x] Player Stats
- [x] 5v5 Team Battles
- [x] Item System (evolution stones & battle boosters)
- [x] Item-based Pokémon evolution
- [ ] Multiplayer PvP
- [ ] Trading System
- [ ] Mobile App
- [ ] Tournaments

---

## 🤝 Contributing

Contributions welcome! Feel free to open issues or submit PRs.

---

<div align="center">

### 👨‍💻 Built by **Ifrun Kader Ruhin**

*Self-taught engineer & constant learner*

---

### 📄 License

**CC BY-NC 4.0** - Free for non-commercial use  
[View License](https://creativecommons.org/licenses/by-nc/4.0/)

---

**⭐ Star this repo if you like it!**

</div>
