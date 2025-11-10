# CLI vs Web App - Feature Comparison

## Overview

PokeTacTix includes **two versions** of the game:
1. **Web App** - Online, multiplayer, account-based
2. **CLI** - Offline, single-player, local save

Both share the same **core battle logic** but differ in features and infrastructure.

---

## Feature Comparison

| Feature | Web App | CLI |
|---------|---------|-----|
| **Platform** | Browser (any device) | Desktop (Windows/Mac/Linux) |
| **Installation** | None (visit URL) | Download binary |
| **Authentication** | ✅ Required (register/login) | ❌ Not needed |
| **Multiplayer** | ✅ PvP (future) | ❌ AI only |
| **Data Storage** | ☁️ Cloud (PostgreSQL) | 💾 Local file (~/.poketactix/) |
| **Internet Required** | ✅ Yes | ❌ No (offline) |
| **Battle Modes** | 1v1, 5v5 | 1v1, 5v5 |
| **Battle Logic** | ✅ Same | ✅ Same |
| **Starter Deck** | ✅ 5 random cards | ✅ 5 random cards |
| **Card Collection** | ✅ Persistent | ✅ Local save |
| **Leveling System** | ✅ XP and levels | ✅ XP and levels |
| **Shop System** | ✅ Full shop | ✅ Simple shop |
| **Battle Rewards** | ✅ Coins + XP | ✅ Coins + XP |
| **Post-Battle Selection** | ✅ Pick AI Pokemon | ✅ Pick AI Pokemon |
| **Statistics** | ✅ Cloud-synced | ✅ Local only |
| **Achievements** | ✅ Yes | ❌ Not planned |
| **Leaderboards** | ✅ Global (future) | ❌ No |
| **Profile System** | ✅ Yes | ❌ No |
| **UI Style** | 🎨 Modern web UI | 🖥️ ASCII art + colors |
| **Updates** | 🔄 Automatic | 📦 Manual download |

---

## Shared Components

Both versions share:

### ✅ Core Battle Engine
- `game/core/engine.go` - Turn processing
- `game/core/damage_calculator.go` - Damage calculation
- `game/core/ai_logic.go` - AI decision making
- `game/core/helper.go` - Helper functions

### ✅ Game Models
- `game/models/game_state.go` - Battle state
- `game/models/player.go` - Player structure
- `game/models/card.go` - Card display

### ✅ Pokemon Fetching
- `internal/pokemon/fetcher.go` - PokeAPI integration
- `internal/pokemon/builder.go` - Card building
- `internal/pokemon/types.go` - Pokemon types

### ✅ Game Utilities
- `game/utils/typechart.go` - Type effectiveness
- `game/utils/utils.go` - Utility functions

---

## Architecture Differences

### Web App Architecture

```
Web App (Online)
├── Frontend (Browser)
│   ├── HTML/CSS/JS
│   └── API calls
├── Backend (Server)
│   ├── cmd/api/
│   └── internal/
│       ├── auth/        # Authentication
│       ├── battle/      # Battle sessions
│       ├── cards/       # Card management
│       └── database/    # PostgreSQL
└── Database (Cloud)
    └── PostgreSQL
```

### CLI Architecture

```
CLI (Offline)
├── Binary (Executable)
│   ├── main.go
│   └── game/
│       ├── commands/    # CLI commands
│       ├── core/        # Battle engine
│       └── models/      # Game models
├── Local Storage
│   └── ~/.poketactix/
│       ├── save.json    # Game state
│       └── cache/       # Pokemon cache
└── PokeAPI (Internet)
    └── Fetch Pokemon data
```

---

## User Experience

### Web App Journey

1. **Visit website** → Register/Login
2. **Get starter deck** → 5 random Pokemon
3. **Battle AI or Players** → Earn coins and XP
4. **Visit shop** → Buy new Pokemon
5. **Build deck** → Select 5 cards
6. **Track progress** → View stats and achievements
7. **Compete** → Leaderboards (future)

### CLI Journey

1. **Download binary** → Run executable
2. **First launch** → Enter name, get starter deck
3. **Battle AI** → Earn coins and XP
4. **Collect Pokemon** → Win battles, buy from shop
5. **Build deck** → Edit deck from collection
6. **Track progress** → View local stats
7. **Play offline** → No internet needed

---

## Development Roadmap

### Web App (Current Tasks 1-9)

- [x] Task 1: Fix logic errors
- [x] Task 2: Database setup
- [x] Task 3: Authentication
- [x] Task 4: Card system
- [x] Task 4.6: Architecture refactoring
- [ ] Task 5: Enhanced battle system (5v5)
- [ ] Task 6: Shop system
- [ ] Task 7: Post-battle selection
- [ ] Task 8: Statistics and profile
- [ ] Task 9: Security hardening
- [ ] **Future**: Multiplayer PvP

### CLI (New Task 10)

- [ ] Task 10.1: Local file persistence
- [ ] Task 10.2: Starter deck generation
- [ ] Task 10.3: Card collection system
- [ ] Task 10.4: Battle rewards
- [ ] Task 10.5: Post-battle selection
- [ ] Task 10.6: Simple shop
- [ ] Task 10.7: ASCII art UI
- [ ] Task 10.8: 5v5 battle mode
- [ ] Task 10.9: Improved battle UI
- [ ] Task 10.10: Statistics tracking
- [ ] Task 10.11: Help system
- [ ] Task 10.12: Quality of life
- [ ] Task 10.13: Performance optimization
- [ ] Task 10.14: Cross-platform compatibility
- [ ] Task 10.15: Distribution package

---

## Why Two Versions?

### Web App Benefits
✅ **Accessibility** - Play anywhere, any device
✅ **Social** - Multiplayer, leaderboards, community
✅ **Always updated** - No manual downloads
✅ **Cloud saves** - Access from anywhere
✅ **Rich UI** - Modern web interface

### CLI Benefits
✅ **Offline play** - No internet required
✅ **Privacy** - No account needed
✅ **Performance** - Fast, lightweight
✅ **Nostalgia** - Classic terminal experience
✅ **Portability** - Single binary, no dependencies

---

## Target Audiences

### Web App
- 🌐 Casual players who want quick access
- 👥 Players who enjoy multiplayer
- 📱 Mobile users
- 🏆 Competitive players (leaderboards)

### CLI
- 💻 Terminal enthusiasts
- 🔒 Privacy-conscious users
- ✈️ Offline players (travel, no internet)
- 🎮 Retro gaming fans
- 🚀 Power users who prefer keyboard

---

## Summary

Both versions offer the **same core gameplay** but cater to different preferences:

- **Web App**: Modern, social, online experience
- **CLI**: Classic, private, offline experience

Players can choose based on their needs, and both versions will be maintained and updated with new features!

---

## Next Steps

1. ✅ Complete web app tasks (5-9)
2. ✅ Implement CLI enhancements (task 10)
3. ✅ Test both versions thoroughly
4. ✅ Create distribution packages
5. ✅ Launch both versions simultaneously
6. 🚀 Gather feedback and iterate

**Goal**: Provide the best Pokemon battle experience in both web and terminal! 🎮⚡
