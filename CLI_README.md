# PokeTacTix CLI

A terminal-based Pokemon battle game that works completely offline. Battle, collect, and build your ultimate Pokemon deck - all from your command line!

```
╔══════════════════════════════════════════════════════════════════════════════════╗
║  ██████╗  ██████╗ ██╗  ██╗███████╗████████╗ █████╗  ██████╗████████╗██╗██╗  ██╗  ║
║  ██╔══██╗██╔═══██╗██║ ██╔╝██╔════╝╚══██╔══╝██╔══██╗██╔════╝╚══██╔══╝██║╚██╗██╔╝  ║
║  ██████╔╝██║   ██║█████╔╝ █████╗     ██║   ███████║██║        ██║   ██║ ╚███╔╝   ║
║  ██╔═══╝ ██║   ██║██╔═██╗ ██╔══╝     ██║   ██╔══██║██║        ██║   ██║ ██╔██╗   ║
║  ██║     ╚██████╔╝██║  ██╗███████╗   ██║   ██║  ██║╚██████╗   ██║   ██║██╔╝ ██╗  ║
║  ╚═╝      ╚═════╝ ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝   ╚═╝   ╚═╝╚═╝  ╚═╝  ║
╚══════════════════════════════════════════════════════════════════════════════════╝
```

## Features

- 🎮 **Offline Gameplay** - No internet required! All Pokemon data is embedded
- ⚔️ **Two Battle Modes** - Quick 1v1 battles or strategic 5v5 tournaments
- 📦 **Pokemon Collection** - Collect and level up Pokemon from Gen 1-5 (649 total)
- 🎨 **ASCII Art** - Beautiful terminal graphics with color support
- 💾 **Local Save System** - Your progress is saved locally with automatic backups
- 🛒 **Shop System** - Buy Pokemon with coins earned from battles
- 📊 **Statistics Tracking** - Track your wins, losses, and progress
- 🔧 **Deck Management** - Build and customize your battle deck
- 🌈 **Cross-Platform** - Works on Windows, macOS, and Linux

## Installation

### Quick Install (macOS/Linux)

```bash
curl -L https://github.com/ifrunruhin12/poketactix/raw/main/scripts/install.sh | bash
```

### Quick Install (Windows)

```powershell
powershell -ExecutionPolicy Bypass -Command "iwr https://github.com/ifrunruhin12/poketactix/raw/main/scripts/install.ps1 | iex"
```

### Manual Installation

1. Download the appropriate binary for your platform from the [releases page](https://github.com/yourusername/poketactix/releases/latest):
   - **Windows**: `poketactix-cli-windows-amd64.exe`
   - **macOS (Intel)**: `poketactix-cli-darwin-amd64`
   - **macOS (Apple Silicon)**: `poketactix-cli-darwin-arm64`
   - **Linux (x64)**: `poketactix-cli-linux-amd64`
   - **Linux (ARM)**: `poketactix-cli-linux-arm64`

2. Make it executable (macOS/Linux):
   ```bash
   chmod +x poketactix-cli-*
   ```

3. Move to a directory in your PATH:
   ```bash
   # macOS/Linux
   sudo mv poketactix-cli-* /usr/local/bin/poketactix-cli
   
   # Or for user-only install
   mv poketactix-cli-* ~/.local/bin/poketactix-cli
   ```

4. Run the game:
   ```bash
   poketactix-cli
   ```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/ifrunruhin12/poketactix.git
cd poketactix

# Build for your platform
make build-cli

# Or build for all platforms
make build-cli-all

# Run the game
./bin/poketactix-cli
```

## Getting Started

### First Launch

When you first run PokeTacTix CLI, you'll be guided through a quick setup:

1. **Enter your name** - Choose your trainer name
2. **Receive starter Pokemon** - Get 5 random Pokemon to start your journey
3. **Get starting coins** - Begin with 500 coins to spend in the shop

### Basic Commands

```
battle (b)      - Start a battle (choose 1v1 or 5v5)
collection (c)  - View your Pokemon collection
deck (d)        - View or edit your battle deck
shop (s)        - Visit the shop to buy Pokemon
stats (st)      - View your battle statistics
help (h)        - Show all available commands
quit (q)        - Exit the game
```

## Gameplay Guide

### Battle System

#### 1v1 Mode
- Quick battles using one random Pokemon from your deck
- Win: 50 coins + 20 XP
- Lose: 10 coins

#### 5v5 Mode
- Strategic battles using all 5 Pokemon in your deck
- Switch Pokemon when one is knocked out
- Win: 150 coins + 15 XP per Pokemon + choose 1 AI Pokemon to keep
- Lose: 25 coins

### Battle Actions

During battle, you can choose from these actions:

- **Attack** - Select one of your Pokemon's 4 moves
  - Each move has different power and stamina cost
  - Type effectiveness matters! (Fire > Grass > Water > Fire)
  
- **Defend** - Reduce incoming damage by 50%
  - Costs no stamina
  - Good for recovering or waiting for stamina
  
- **Pass** - Skip your turn and restore 20 stamina
  - Use when low on stamina
  - Strategic for setting up powerful moves
  
- **Sacrifice** - Deal 30 damage but lose 20 HP
  - Risky but doesn't cost stamina
  - Good for finishing off weak opponents
  
- **Surrender** - Give up the battle
  - Receive reduced rewards
  - Use if battle is unwinnable

### Pokemon Collection

View and manage your Pokemon:

```bash
# View all Pokemon
collection

# Filter by type
collection filter fire

# Sort by level
collection sort level

# Search by name
collection search pikachu
```

### Deck Management

Your deck consists of exactly 5 Pokemon used in battles:

```bash
# View current deck
deck

# Edit deck
deck edit

# Tips:
# - Balance types for coverage
# - Include high-level Pokemon
# - Consider move diversity
```

### Shop System

Buy Pokemon with coins earned from battles:

```bash
# Open shop
shop

# Prices:
# - Common: 100 coins
# - Uncommon: 250 coins
# - Rare: 500 coins

# Shop refreshes every 10 battles
```

### Leveling System

Pokemon gain XP from battles and level up automatically:

- Each level increases all stats (HP, Attack, Defense, Speed)
- Higher level Pokemon are stronger in battle
- XP is awarded to all participating Pokemon in 5v5 battles

### Statistics

Track your progress:

```bash
stats

# View:
# - Total battles and win rate
# - Separate stats for 1v1 and 5v5
# - Coins earned
# - Highest level Pokemon
# - Recent battle history
```

## Advanced Features

### Save System

- **Auto-save**: Game saves automatically after battles and major actions
- **Save location**: `~/.poketactix/save.json`
- **Backups**: Last 3 saves are kept automatically
- **Export/Import**: Backup your save file manually

### Battle Speed Settings

Adjust battle animation speed:

```bash
settings speed slow    # Slower animations
settings speed normal  # Default speed
settings speed fast    # Faster animations
```

### Quick Battle Mode

Skip animations for faster gameplay:

```bash
settings quickbattle on
```

## Terminal Requirements

### Minimum Requirements

- **Size**: 80 columns × 24 rows (larger recommended)
- **Colors**: ANSI color support (optional but recommended)
- **OS**: Windows 10+, macOS 10.12+, or Linux with modern terminal

### Recommended Terminals

- **macOS**: Terminal.app, iTerm2
- **Linux**: GNOME Terminal, Konsole, Alacritty
- **Windows**: Windows Terminal, PowerShell, Git Bash

### Color Support

The game automatically detects color support. If colors don't work:

```bash
# Force disable colors
export NO_COLOR=1
poketactix-cli
```

## Troubleshooting

### Game won't start

**Problem**: Binary won't execute
```bash
# macOS: Remove quarantine attribute
xattr -d com.apple.quarantine poketactix-cli

# Linux: Make executable
chmod +x poketactix-cli
```

### Save file issues

**Problem**: Corrupted save file
```bash
# The game automatically tries to restore from backup
# If that fails, delete save file to start fresh:
rm ~/.poketactix/save.json
```

### Display issues

**Problem**: UI looks broken or misaligned
- Ensure terminal is at least 80×24
- Try resizing terminal window
- Check terminal font supports box-drawing characters

**Problem**: Colors not showing
- Check terminal supports ANSI colors
- Try a different terminal emulator
- Use `export TERM=xterm-256color`

### Performance issues

**Problem**: Game is slow
- Close other terminal applications
- Reduce battle speed: `settings speed fast`
- Enable quick battle mode: `settings quickbattle on`

### Platform-specific issues

**Windows**: If you see "cannot be loaded because running scripts is disabled"
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

**macOS**: If you see "cannot be opened because the developer cannot be verified"
```bash
xattr -d com.apple.quarantine poketactix-cli
```

**Linux**: If binary won't run
```bash
# Install required libraries
sudo apt-get install libc6  # Debian/Ubuntu
sudo yum install glibc      # RHEL/CentOS
```

## File Locations

- **Save file**: `~/.poketactix/save.json`
- **Backups**: `~/.poketactix/backups/`
- **Config**: `~/.poketactix/config.json` (future feature)

## Tips & Tricks

1. **Type Advantage**: Learn the type chart for strategic battles
   - Fire beats Grass, Ice, Bug, Steel
   - Water beats Fire, Ground, Rock
   - Grass beats Water, Ground, Rock
   - Electric beats Water, Flying
   - And more!

2. **Stamina Management**: Don't spam your strongest moves
   - Use Pass to recover stamina
   - Mix in weaker moves to conserve stamina
   - Defend when low on stamina

3. **Deck Building**: Balance is key
   - Include multiple types for coverage
   - Mix high-attack and high-defense Pokemon
   - Keep your deck leveled up

4. **Shop Strategy**: Save coins for rare Pokemon
   - Common Pokemon are easy to get from battles
   - Rare Pokemon are worth the investment
   - Shop refreshes every 10 battles

5. **5v5 Battles**: Order matters
   - Put your strongest Pokemon first
   - Save a strong Pokemon for last
   - Consider type matchups when switching

## FAQ

**Q: Is internet required?**
A: No! The game works completely offline. All Pokemon data is embedded in the binary.

**Q: How many Pokemon are available?**
A: 649 Pokemon from Generations 1-5 are included.

**Q: Can I transfer my save to another computer?**
A: Yes! Copy `~/.poketactix/save.json` to the same location on another computer.

**Q: Does the game update automatically?**
A: No, you need to manually download new versions. Check the releases page.

**Q: Can I play on multiple computers?**
A: Yes, but you'll need to manually sync your save file between computers.

**Q: Is there multiplayer?**
A: Not currently. The CLI version is single-player only.

**Q: How do I reset my progress?**
A: Use the `reset` command in-game, or delete `~/.poketactix/save.json`

**Q: Can I customize the colors?**
A: Not currently, but this is planned for a future update.

## Contributing

Found a bug or have a feature request? Please open an issue on GitHub!

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Credits

- Pokemon data from [PokeAPI](https://pokeapi.co/)
- Built with Go and love for Pokemon

## Version

Check your version:
```bash
poketactix-cli --version
```

---

**Enjoy your Pokemon journey in the terminal! 🎮⚡**
