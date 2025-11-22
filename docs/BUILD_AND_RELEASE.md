# Build and Release Guide

This document provides a quick reference for building and releasing PokeTacTix CLI.

## Quick Commands

### Build for Current Platform
```bash
make build-cli
```

### Build for All Platforms
```bash
make build-cli-all
```

### Create Release Package
```bash
make release
# or
./scripts/create-release.sh
```

### Check Version
```bash
./bin/poketactix-cli --version
```

## Build Targets

The Makefile includes the following CLI-related targets:

- `build-cli` - Build CLI binary for current platform with version info
- `build-cli-all` - Build CLI binaries for all platforms (Windows, macOS, Linux)
- `run-cli` - Build and run CLI locally
- `release` - Create complete release package with checksums

## Supported Platforms

| Platform | Architecture | Binary Name |
|----------|-------------|-------------|
| Windows | amd64 | `poketactix-cli-windows-amd64.exe` |
| macOS | amd64 (Intel) | `poketactix-cli-darwin-amd64` |
| macOS | arm64 (Apple Silicon) | `poketactix-cli-darwin-arm64` |
| Linux | amd64 | `poketactix-cli-linux-amd64` |
| Linux | arm64 | `poketactix-cli-linux-arm64` |

## Build Features

### Version Information

Binaries are built with embedded version information:
- Version tag from git
- Build date (UTC)
- Commit hash
- Go version
- Platform and architecture

This information is displayed with `--version` flag.

### Optimizations

Binaries are built with:
- `-s` flag: Strip symbol table
- `-w` flag: Strip DWARF debugging info
- Result: Smaller binary size (~5-6 MB per binary)

### Embedded Data

Pokemon data is embedded in the binary using Go's `//go:embed` directive:
- No external data files needed
- Single binary distribution
- ~5 MB of Pokemon data included

## Installation Scripts

### Unix/Linux/macOS
```bash
./scripts/install.sh
```

Features:
- Auto-detects platform and architecture
- Downloads latest release from GitHub
- Installs to `/usr/local/bin` (with sudo) or `~/.local/bin` (without)
- Checks PATH and provides instructions if needed

### Windows
```powershell
./scripts/install.ps1
```

Features:
- Auto-detects architecture
- Downloads latest release from GitHub
- Installs to `%LOCALAPPDATA%\PokeTacTix`
- Automatically adds to user PATH
- Provides instructions for terminal restart

## Release Process

### Automated (GitHub Actions)

1. Create and push a version tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. GitHub Actions automatically:
   - Builds all platform binaries
   - Creates release archives
   - Generates checksums
   - Creates GitHub release
   - Uploads all assets

### Manual

1. Build binaries:
   ```bash
   make build-cli-all
   ```

2. Create release package:
   ```bash
   ./scripts/create-release.sh
   ```

3. Follow prompts and upload to GitHub

See [RELEASE_PROCESS.md](RELEASE_PROCESS.md) for detailed instructions.

## File Structure

```
bin/                                    # Built binaries
  poketactix-cli-windows-amd64.exe
  poketactix-cli-darwin-amd64
  poketactix-cli-darwin-arm64
  poketactix-cli-linux-amd64
  poketactix-cli-linux-arm64

scripts/
  install.sh                           # Unix installation script
  install.ps1                          # Windows installation script
  create-release.sh                    # Release packaging script

.github/workflows/
  release-cli.yml                      # GitHub Actions workflow

docs/
  RELEASE_PROCESS.md                   # Detailed release guide
  BUILD_AND_RELEASE.md                 # This file

CLI_README.md                          # User documentation
CHANGELOG.md                           # Version history
```

## Testing Builds

### Test Current Platform
```bash
make build-cli
./bin/poketactix-cli --version
./bin/poketactix-cli
```

### Test Cross-Platform Builds
```bash
make build-cli-all
ls -lh bin/poketactix-cli-*
```

### Test Installation Script
```bash
# Unix/Linux/macOS
./scripts/install.sh

# Windows (in PowerShell)
./scripts/install.ps1
```

## Troubleshooting

### Build Errors

**Problem**: `build constraints exclude all Go files`
- **Solution**: Remove `//go:build` tags from source files

**Problem**: `undefined: syscall.SIGWINCH` on Windows
- **Solution**: Use platform-specific files with build tags

**Problem**: Binary too large
- **Solution**: Ensure `-s -w` flags are used in ldflags

### Cross-Compilation Issues

**Problem**: CGO errors when cross-compiling
- **Solution**: Ensure `CGO_ENABLED=0` or use pure Go code

**Problem**: Missing dependencies
- **Solution**: Run `go mod download` and `go mod tidy`

### Version Information

**Problem**: Version shows as "dev"
- **Solution**: Ensure git tags exist or set VERSION manually

**Problem**: Build date incorrect
- **Solution**: Check system time and timezone

## Development Workflow

1. Make changes to code
2. Test locally: `make build-cli && ./bin/poketactix-cli`
3. Test all platforms: `make build-cli-all`
4. Update CHANGELOG.md
5. Commit changes
6. Create version tag
7. Push tag to trigger release

## CI/CD

GitHub Actions workflow (`.github/workflows/release-cli.yml`) runs on:
- Push of version tags (`v*`)
- Builds all platforms
- Creates GitHub release
- Uploads binaries and archives

## Binary Size

Typical binary sizes:
- Windows (amd64): ~5.7 MB
- macOS (amd64): ~5.7 MB
- macOS (arm64): ~5.4 MB
- Linux (amd64): ~5.6 MB
- Linux (arm64): ~5.4 MB

Size includes:
- Go runtime
- Application code
- Embedded Pokemon data (~5 MB)
- No external dependencies

## Distribution

Binaries are distributed via:
1. GitHub Releases (primary)
2. Direct download links
3. Installation scripts

No package managers yet, but planned for future:
- Homebrew (macOS/Linux)
- Chocolatey (Windows)
- Snap (Linux)
- AUR (Arch Linux)

## Security

### Checksums

All releases include:
- SHA256SUMS - SHA256 checksums for verification
- MD5SUMS - MD5 checksums for compatibility

Users can verify downloads:
```bash
# Unix/Linux/macOS
sha256sum -c SHA256SUMS

# Windows (PowerShell)
Get-FileHash poketactix-cli-windows-amd64.exe -Algorithm SHA256
```

### Signing

Future enhancement: Code signing for binaries
- Windows: Authenticode signing
- macOS: Apple Developer ID signing
- Linux: GPG signatures

## Resources

- [Go Cross Compilation](https://golang.org/doc/install/source#environment)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)

---

**Last Updated**: 2025-11-22
