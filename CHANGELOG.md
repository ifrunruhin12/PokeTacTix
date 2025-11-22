# Changelog

All notable changes to PokeTacTix CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial CLI implementation with offline gameplay
- 649 Pokemon from Gen 1-5
- 1v1 and 5v5 battle modes
- Local save system with automatic backups
- Shop system for purchasing Pokemon
- Statistics tracking
- Deck management
- ASCII art and color support
- Cross-platform support (Windows, macOS, Linux)

### Changed
- N/A

### Fixed
- N/A

## [1.0.0] - YYYY-MM-DD

### Added
- Initial release of PokeTacTix CLI
- Complete offline Pokemon battle game
- 649 Pokemon from Generations 1-5
- Two battle modes: 1v1 and 5v5
- Pokemon collection and deck management
- Shop system with dynamic inventory
- Statistics and battle history tracking
- Local save system with automatic backups
- ASCII art graphics with ANSI color support
- Cross-platform binaries for Windows, macOS, and Linux
- Comprehensive documentation and installation scripts

### Features
- **Offline Gameplay**: No internet required, all data embedded
- **Battle System**: Strategic turn-based battles with type effectiveness
- **Progression**: Level up Pokemon through battles
- **Collection**: Collect Pokemon from battles and shop purchases
- **Customization**: Build and customize your 5-Pokemon deck
- **Statistics**: Track wins, losses, and progress
- **Quality of Life**: Auto-save, backups, battle speed settings

### Technical
- Built with Go for performance and portability
- Single binary with no external dependencies
- Embedded Pokemon data (~5MB)
- Local JSON save files with compression
- Terminal resize handling (Unix)
- Color detection and fallback support

---

## Version History

- **1.0.0** - Initial release

---

## How to Update

### Automatic (Recommended)
Run the installation script again:

**macOS/Linux:**
```bash
curl -L https://github.com/yourusername/poketactix/raw/main/scripts/install.sh | bash
```

**Windows:**
```powershell
powershell -ExecutionPolicy Bypass -Command "iwr https://github.com/yourusername/poketactix/raw/main/scripts/install.ps1 | iex"
```

### Manual
1. Download the latest binary from [releases](https://github.com/yourusername/poketactix/releases/latest)
2. Replace your existing binary
3. Your save file will be preserved

---

## Breaking Changes

None yet! We'll document any breaking changes here in future releases.

---

## Migration Guide

### From Pre-1.0 Versions
If you were using a development version, your save file should be compatible. If you encounter issues:

1. Backup your save file: `~/.poketactix/save.json`
2. Delete the save file
3. Restart the game
4. If needed, manually edit the backup to match the new format

---

## Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/poketactix/issues)
- **Documentation**: [CLI_README.md](CLI_README.md)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/poketactix/discussions)
