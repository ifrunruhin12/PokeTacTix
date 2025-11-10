# CLI Enhancement Tasks

## Overview

The CLI version is a **local, offline, single-player** experience that mirrors the web app's game logic but with:
- ✅ No authentication (download and play)
- ✅ No multiplayer (always Player vs AI)
- ✅ Local file persistence (save game state)
- ✅ Starter deck (5 random non-legendary cards)
- ✅ Card collection and progression
- ✅ Enhanced ASCII art and UI
- ✅ Same battle logic as web (1v1 and 5v5)
- ❌ No online features (multiplayer, leaderboards)

## Task 10: CLI Enhancement and Polish

### 10.1 Implement Local File Persistence
- [ ] Create `cli/storage/` package for local file operations
- [ ] Implement `SaveGameState()` function to save player data to JSON file
- [ ] Implement `LoadGameState()` function to load player data from JSON file
- [ ] Store game state in `~/.poketactix/save.json` (cross-platform)
- [ ] Save player data structure:
  - Player name
  - Card collection (all owned cards with levels and XP)
  - Current deck (5 cards)
  - Statistics (wins, losses, total battles)
  - Coins earned
  - Last played timestamp
- [ ] Auto-save after every battle
- [ ] Add manual save command
- [ ] Handle corrupted save files gracefully (backup and reset)
- _Requirements: Local persistence, user experience_

### 10.2 Implement Starter Deck Generation for CLI
- [ ] Create `cli/setup/` package for first-time setup
- [ ] Implement `GenerateStarterDeck()` for CLI (reuse logic from web)
- [ ] Generate 5 random non-legendary, non-mythical Pokemon
- [ ] Initialize all cards at level 1 with 0 XP
- [ ] Prompt user for player name on first launch
- [ ] Save starter deck to local file
- [ ] Display welcome message with starter Pokemon
- _Requirements: First-time user experience_

### 10.3 Implement Card Collection System
- [ ] Create `cli/collection/` package for card management
- [ ] Implement `ViewCollection()` command to display all owned cards
- [ ] Show card details: name, level, XP, stats, types, moves
- [ ] Add sorting options (by level, name, type)
- [ ] Add filtering options (by type, level range)
- [ ] Implement `ViewDeck()` command to show current deck
- [ ] Implement `EditDeck()` command to modify deck (select 5 from collection)
- [ ] Validate deck changes (must have exactly 5 cards)
- [ ] Save deck changes to local file
- _Requirements: Card management, user experience_

### 10.4 Implement Battle Rewards System
- [ ] Award coins after battles:
  - 1v1 win: 50 coins
  - 1v1 loss: 10 coins
  - 5v5 win: 150 coins
  - 5v5 loss: 25 coins
- [ ] Award XP to participating Pokemon:
  - 1v1: 20 XP to active Pokemon
  - 5v5: 15 XP to each Pokemon that participated
- [ ] Trigger level-up checks after XP distribution
- [ ] Display level-up notifications with stat increases
- [ ] Update and save player statistics
- [ ] Save updated card data and coins to local file
- _Requirements: Progression system, rewards_

### 10.5 Implement Post-Battle Pokemon Selection (5v5)
- [ ] After winning 5v5 battle, display all 5 AI Pokemon
- [ ] Show Pokemon details: name, level, types, stats
- [ ] Allow player to select one Pokemon to add to collection
- [ ] Add selected Pokemon at level 1 with 0 XP
- [ ] Allow duplicates (for deck building flexibility)
- [ ] Save new Pokemon to collection
- [ ] Display confirmation message
- _Requirements: Rewards, card collection_

### 10.6 Implement Simple Shop System
- [ ] Create `cli/shop/` package for shop functionality
- [ ] Generate shop inventory with 10-15 Pokemon
- [ ] Pricing: common 100, uncommon 250, rare 500
- [ ] Exclude legendary and mythical from shop (too expensive for CLI)
- [ ] Display shop inventory with prices
- [ ] Implement `BuyCard()` function with coin validation
- [ ] Deduct coins and add Pokemon to collection
- [ ] Refresh shop inventory every 10 battles
- [ ] Save shop state and player coins
- _Requirements: Card acquisition, progression_

### 10.7 Enhance CLI UI with ASCII Art
- [ ] Create `cli/ui/` package for UI components
- [ ] Add colorized output using ANSI color codes
- [ ] Design ASCII art for:
  - Game logo/title screen
  - Pokemon type badges (fire, water, grass, etc.)
  - HP bars (visual representation)
  - Stamina bars
  - Battle arena layout
  - Victory/defeat banners
  - Level-up animations
- [ ] Implement `PrintCard()` with enhanced formatting
- [ ] Add type-based color coding (red for fire, blue for water, etc.)
- [ ] Create battle log with better formatting
- [ ] Add loading animations for Pokemon fetching
- [ ] Implement progress bars for XP
- _Requirements: User experience, visual appeal_

### 10.8 Implement 5v5 Battle Mode for CLI
- [ ] Extend battle system to support 5v5 mode
- [ ] Allow player to select battle mode (1v1 or 5v5)
- [ ] Display all 5 Pokemon in deck before battle
- [ ] Show active Pokemon prominently
- [ ] Display inactive Pokemon as "benched"
- [ ] Implement Pokemon switching when one is knocked out
- [ ] Show round progression (Round 1/5, 2/5, etc.)
- [ ] Display battle summary at end (which Pokemon won/lost)
- [ ] Award appropriate rewards for 5v5 battles
- _Requirements: Battle modes, game logic parity with web_

### 10.9 Improve Battle UI and Flow
- [ ] Redesign battle screen layout:
  ```
  ╔════════════════════════════════════════════════════════════╗
  ║                    POKEMON BATTLE                          ║
  ╠════════════════════════════════════════════════════════════╣
  ║  Player: Pikachu (Lv 15)          AI: Charizard (Lv 12)  ║
  ║  HP: ████████░░ 80/100            HP: ██████░░░░ 60/100   ║
  ║  Stamina: ████░░ 40/60            Stamina: ███░░░ 30/50   ║
  ║  Type: Electric                   Type: Fire, Flying      ║
  ╠════════════════════════════════════════════════════════════╣
  ║  Available Moves:                                          ║
  ║  1. Thunderbolt (Power: 90, Stamina: 30) [Electric]      ║
  ║  2. Quick Attack (Power: 40, Stamina: 13) [Normal]       ║
  ║  3. Thunder Wave (Power: 20, Stamina: 7) [Electric]      ║
  ║  4. Iron Tail (Power: 100, Stamina: 33) [Steel]          ║
  ╠════════════════════════════════════════════════════════════╣
  ║  Actions: [A]ttack [D]efend [P]ass [S]acrifice [Q]uit    ║
  ╚════════════════════════════════════════════════════════════╝
  ```
- [ ] Add turn-by-turn battle log
- [ ] Show damage calculations
- [ ] Display type effectiveness indicators
- [ ] Add battle statistics (turns taken, damage dealt)
- [ ] Implement quick rematch option
- _Requirements: User experience, visual clarity_

### 10.10 Add Statistics and Progress Tracking
- [ ] Create `cli/stats/` package for statistics
- [ ] Track and display:
  - Total battles (1v1 and 5v5 separately)
  - Win/loss record
  - Win rate percentage
  - Total coins earned
  - Highest level Pokemon
  - Total Pokemon collected
  - Favorite Pokemon (most used)
- [ ] Implement `ViewStats()` command
- [ ] Display statistics in formatted table
- [ ] Add battle history (last 10 battles)
- [ ] Save statistics to local file
- _Requirements: Progression tracking, player engagement_

### 10.11 Implement Help System and Commands
- [ ] Create comprehensive help command
- [ ] Document all available commands:
  - `battle` - Start a battle (1v1 or 5v5)
  - `collection` - View all owned Pokemon
  - `deck` - View/edit current deck
  - `shop` - Browse and buy Pokemon
  - `stats` - View player statistics
  - `save` - Manually save game
  - `help` - Show help menu
  - `quit` - Exit game
- [ ] Add command aliases (e.g., `b` for battle, `c` for collection)
- [ ] Implement tab completion for commands
- [ ] Add contextual help during battles
- [ ] Create tutorial for first-time players
- _Requirements: User experience, accessibility_

### 10.12 Add Quality of Life Features
- [ ] Implement auto-save after every action
- [ ] Add confirmation prompts for important actions
- [ ] Implement undo for deck changes
- [ ] Add quick battle mode (skip animations)
- [ ] Implement battle speed settings (slow, normal, fast)
- [ ] Add sound effects (optional, using system beep)
- [ ] Create backup save files (keep last 3 saves)
- [ ] Add export/import save file feature
- [ ] Implement reset progress option
- _Requirements: User experience, data safety_

### 10.13 Optimize CLI Performance
- [ ] Cache Pokemon data locally to reduce API calls
- [ ] Implement lazy loading for card collection
- [ ] Optimize battle rendering (only redraw changed elements)
- [ ] Add loading indicators for slow operations
- [ ] Implement async Pokemon fetching
- [ ] Reduce startup time to < 1 second
- [ ] Optimize save file size (compress if needed)
- _Requirements: Performance, user experience_

### 10.14 Cross-Platform Compatibility
- [ ] Test on Windows, macOS, and Linux
- [ ] Handle different terminal sizes gracefully
- [ ] Use cross-platform file paths
- [ ] Detect terminal color support
- [ ] Fallback to plain text if colors not supported
- [ ] Handle different line endings (CRLF vs LF)
- [ ] Test with different terminal emulators
- _Requirements: Compatibility, accessibility_

### 10.15 Create CLI Distribution Package
- [ ] Build binaries for all platforms:
  - Windows (amd64, arm64)
  - macOS (amd64, arm64)
  - Linux (amd64, arm64)
- [ ] Create installation scripts
- [ ] Write CLI-specific README
- [ ] Add version command
- [ ] Implement auto-update checker (optional)
- [ ] Create release notes
- [ ] Package with assets (if any)
- _Requirements: Distribution, user onboarding_

## Implementation Priority

### Phase 1: Core Functionality (Must Have)
1. ✅ Local file persistence (10.1)
2. ✅ Starter deck generation (10.2)
3. ✅ Card collection system (10.3)
4. ✅ Battle rewards (10.4)
5. ✅ 5v5 battle mode (10.8)

### Phase 2: Enhanced Experience (Should Have)
6. ✅ Enhanced UI with ASCII art (10.7)
7. ✅ Improved battle UI (10.9)
8. ✅ Statistics tracking (10.10)
9. ✅ Help system (10.11)
10. ✅ Post-battle Pokemon selection (10.5)

### Phase 3: Polish (Nice to Have)
11. ✅ Simple shop system (10.6)
12. ✅ Quality of life features (10.12)
13. ✅ Performance optimization (10.13)
14. ✅ Cross-platform compatibility (10.14)
15. ✅ Distribution package (10.15)

## Estimated Time

- **Phase 1**: 2-3 days
- **Phase 2**: 2-3 days
- **Phase 3**: 1-2 days
- **Total**: 5-8 days

## Success Criteria

✅ CLI is fully playable offline
✅ Game state persists between sessions
✅ Battle logic matches web app
✅ UI is visually appealing with colors and ASCII art
✅ Card collection and progression work smoothly
✅ Performance is smooth (< 1s startup, instant commands)
✅ Works on Windows, macOS, and Linux
✅ Easy to download and play (no setup required)
