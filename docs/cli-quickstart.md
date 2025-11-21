# PokeTacTix CLI - Quick Start

## Building the CLI

```bash
# Build the CLI binary
make build-cli

# Or manually:
go build -tags cli -o bin/poketactix-cli ./cmd/cli
```

## Running the CLI

```bash
# Run the CLI
./bin/poketactix-cli

# Or use make:
make run-cli
```

## First Launch

On first launch, you'll go through the onboarding process:

1. **Welcome Screen** - See the PokeTacTix logo and game introduction
2. **Enter Your Name** - Choose your player name (2-20 characters)
3. **Tutorial** - Optional tutorial explaining game mechanics
4. **Starter Pokemon** - Receive 5 random Pokemon to start your journey
5. **Game Initialized** - Your save file is created at `~/.poketactix/save.json`

## Available Commands

Once in the game, you can use these commands:

- `help` or `h` - Show available commands
- `info` or `i` - Display your player information
- `collection` or `c` - View your Pokemon collection
- `deck` or `d` - View your battle deck
- `battle` or `b` - Start a battle (1v1 or 5v5) ⭐ NEW!
- `stats` or `st` - View your battle statistics
- `quit` or `q` - Exit the game
- `reset` - Delete save file and start fresh

See [CLI_BATTLE_GUIDE.md](CLI_BATTLE_GUIDE.md) for detailed battle instructions.

## Current Status

✅ **Implemented:**
- First-time setup and onboarding
- Player name validation
- Starter deck generation (5 random Pokemon)
- Save/load game state
- Collection viewer
- Deck viewer
- Stats viewer
- **Battle system (1v1 and 5v5)** ⭐ NEW!
- Pokemon leveling and XP system
- Battle rewards and coins
- Post-battle Pokemon selection (5v5 victories)

🚧 **Coming Soon:**
- Shop system
- Deck management/editing
- Battle history viewer
- Tournament mode
- And more!

## Save File Location

Your game progress is saved at: `~/.poketactix/save.json`

To start fresh, either:
- Use the `reset` command in-game
- Delete the file manually: `rm ~/.poketactix/save.json`

## Testing the Setup Flow

To test the onboarding flow again:

```bash
# Delete your save file
rm ~/.poketactix/save.json

# Run the CLI again
./bin/poketactix-cli
```

## Development

The CLI is built with the `cli` build tag to enable offline Pokemon data:

```bash
# Run tests
go test -tags cli -v ./internal/cli/setup

# Build for development
go build -tags cli -o bin/poketactix-cli ./cmd/cli
```

## Architecture

```
cmd/cli/main.go              # CLI entry point
internal/cli/
  setup/                     # Onboarding and initialization
  storage/                   # Save/load game state
  ui/                        # Terminal UI rendering
  commands/                  # Command handlers (coming soon)
internal/pokemon/
  offline_data.go            # Embedded Pokemon data
  data/pokemon_data.json     # Pokemon database
```

## Notes

- The CLI works completely offline - no internet required
- All data is stored locally in `~/.poketactix/`
- The CLI uses ANSI colors if your terminal supports them
- Minimum terminal size: 80x24 (recommended)
