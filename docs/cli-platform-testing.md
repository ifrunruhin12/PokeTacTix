# CLI Platform Testing Guide

This document outlines the cross-platform testing strategy for PokeTacTix CLI and documents platform-specific issues and workarounds.

## Supported Platforms

### Operating Systems
- **Windows**: 10, 11 (amd64)
- **macOS**: 10.15+ (amd64, arm64/Apple Silicon)
- **Linux**: Ubuntu 20.04+, Fedora 35+, Arch Linux (amd64, arm64)

### Terminal Emulators

#### Windows
- **Windows Terminal** (Recommended) - Full ANSI color support
- **PowerShell** - Full support with Windows 10+
- **Command Prompt (cmd.exe)** - Limited color support, use Windows Terminal instead
- **ConEmu** - Full support with ANSI enabled
- **Git Bash** - Full support

#### macOS
- **Terminal.app** (Default) - Full support
- **iTerm2** (Recommended) - Full support with enhanced features
- **Alacritty** - Full support
- **Kitty** - Full support

#### Linux
- **GNOME Terminal** - Full support
- **Konsole** (KDE) - Full support
- **Alacritty** - Full support
- **Kitty** - Full support
- **xterm** - Full support
- **tmux/screen** - Full support

## Testing Checklist

### Basic Functionality Tests

Run these tests on each platform:

- [ ] **Installation**
  - Binary downloads and runs without errors
  - No missing dependencies
  - Correct architecture (amd64/arm64)

- [ ] **First Launch**
  - Welcome screen displays correctly
  - Player name prompt works
  - Starter deck generation succeeds
  - Save file created in correct location

- [ ] **Terminal Compatibility**
  - Minimum size warning (< 80x24) displays
  - Terminal resize handled gracefully
  - Color support detected correctly
  - NO_COLOR environment variable respected

- [ ] **File Operations**
  - Save file created at `~/.poketactix/save.json.gz`
  - Backups created successfully
  - Export/import works correctly
  - Path separators correct for OS

- [ ] **Battle System**
  - 1v1 battles work correctly
  - 5v5 battles work correctly
  - ASCII art renders properly
  - Colors display correctly (if supported)
  - Battle log updates correctly

- [ ] **Commands**
  - All commands execute without errors
  - Help system displays correctly
  - Collection viewing works
  - Deck management works
  - Shop system works
  - Stats display correctly

### Performance Tests

- [ ] **Startup Time**
  - Cold start < 1 second
  - Warm start < 500ms

- [ ] **Save/Load**
  - Save operation < 100ms
  - Load operation < 100ms

- [ ] **Rendering**
  - Battle screen render < 50ms
  - No visible lag or flicker

- [ ] **Memory Usage**
  - Baseline < 50MB
  - During battle < 100MB

## Platform-Specific Issues

### Windows

#### Issue: Color Support in cmd.exe
**Status**: Known Limitation  
**Description**: Windows Command Prompt (cmd.exe) has limited ANSI color support before Windows 10.  
**Workaround**: Use Windows Terminal, PowerShell, or set `NO_COLOR=1` environment variable.  
**Detection**: Automatically falls back to plain text if colors not supported.

#### Issue: Path Separators
**Status**: Resolved  
**Description**: Windows uses backslash (`\`) instead of forward slash (`/`).  
**Solution**: All code uses `filepath.Join()` for cross-platform compatibility.

#### Issue: Home Directory
**Status**: Resolved  
**Description**: Windows home directory is typically `C:\Users\<username>`.  
**Solution**: Using `os.UserHomeDir()` which handles Windows correctly.

#### Issue: File Permissions
**Status**: Resolved  
**Description**: Windows doesn't use Unix-style file permissions.  
**Solution**: Permission checks are skipped on Windows in tests.

### macOS

#### Issue: Apple Silicon (M1/M2) Compatibility
**Status**: Resolved  
**Description**: Need separate arm64 binary for Apple Silicon Macs.  
**Solution**: Building separate binaries for amd64 and arm64 architectures.

#### Issue: Gatekeeper Security
**Status**: Known Limitation  
**Description**: macOS may block unsigned binaries from running.  
**Workaround**: Users need to right-click and select "Open" on first launch, or run:
```bash
xattr -d com.apple.quarantine poketactix-cli
```

#### Issue: Terminal.app Color Support
**Status**: Resolved  
**Description**: Default Terminal.app supports colors but may need profile configuration.  
**Solution**: Color detection works automatically. Users can set profile to "Pro" or "Homebrew" for better colors.

### Linux

#### Issue: Terminal Emulator Variety
**Status**: Resolved  
**Description**: Many different terminal emulators with varying capabilities.  
**Solution**: Comprehensive terminal detection in `DetectColorSupport()`.

#### Issue: Save Directory Permissions
**Status**: Resolved  
**Description**: Some distributions have strict home directory permissions.  
**Solution**: Save directory created with 0700 permissions (user-only access).

#### Issue: Missing Dependencies
**Status**: Resolved  
**Description**: Binary is statically linked, no external dependencies required.  
**Solution**: Go builds static binaries by default.

## Testing Procedure

### Automated Testing

Run the test suite on each platform:

```bash
# Run all tests
go test ./internal/cli/...

# Run with verbose output
go test -v ./internal/cli/...

# Run specific test suites
go test ./internal/cli/ui/...
go test ./internal/cli/storage/...
go test ./internal/cli/commands/...
```

### Manual Testing

1. **Build for target platform**:
   ```bash
   # Windows
   GOOS=windows GOARCH=amd64 go build -o bin/poketactix-cli-windows-amd64.exe cmd/cli/main.go
   
   # macOS Intel
   GOOS=darwin GOARCH=amd64 go build -o bin/poketactix-cli-darwin-amd64 cmd/cli/main.go
   
   # macOS Apple Silicon
   GOOS=darwin GOARCH=arm64 go build -o bin/poketactix-cli-darwin-arm64 cmd/cli/main.go
   
   # Linux
   GOOS=linux GOARCH=amd64 go build -o bin/poketactix-cli-linux-amd64 cmd/cli/main.go
   ```

2. **Test on target platform**:
   - Copy binary to target system
   - Run binary in various terminal emulators
   - Complete the testing checklist above
   - Document any issues found

3. **Test color support**:
   ```bash
   # Test with colors (default)
   ./poketactix-cli
   
   # Test without colors
   NO_COLOR=1 ./poketactix-cli
   
   # Force colors
   CLICOLOR_FORCE=1 ./poketactix-cli
   ```

4. **Test terminal sizes**:
   - Resize terminal to minimum (80x24)
   - Resize to smaller than minimum (< 80x24)
   - Resize to large size (> 120x40)
   - Verify warning messages and layout adaptation

5. **Test file operations**:
   ```bash
   # Check save file location
   ls -la ~/.poketactix/
   
   # Verify file permissions (Unix-like)
   stat ~/.poketactix/save.json.gz
   
   # Test backup creation
   ls -la ~/.poketactix/save_backup_*
   ```

## Continuous Integration

### GitHub Actions

The project uses GitHub Actions for automated cross-platform testing:

```yaml
# .github/workflows/test.yml
name: Cross-Platform Tests

on: [push, pull_request]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: [1.21]
    
    runs-on: ${{ matrix.os }}
    
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}
      
      - name: Run tests
        run: go test -v ./internal/cli/...
      
      - name: Build binary
        run: go build -o poketactix-cli cmd/cli/main.go
```

## Reporting Issues

When reporting platform-specific issues, please include:

1. **Platform Information**:
   - Operating System and version
   - Architecture (amd64/arm64)
   - Terminal emulator and version

2. **Environment**:
   - `TERM` environment variable value
   - `COLORTERM` environment variable value
   - Terminal size (from `echo $COLUMNS x $LINES`)

3. **Steps to Reproduce**:
   - Exact commands run
   - Expected behavior
   - Actual behavior

4. **Logs**:
   - Any error messages
   - Screenshots if visual issue

## Known Working Configurations

### Windows
- ✅ Windows 11 + Windows Terminal + PowerShell 7
- ✅ Windows 10 + Windows Terminal + PowerShell 5.1
- ✅ Windows 10 + Git Bash
- ⚠️ Windows 10 + cmd.exe (limited colors)

### macOS
- ✅ macOS 13 (Ventura) + Terminal.app
- ✅ macOS 13 (Ventura) + iTerm2
- ✅ macOS 12 (Monterey) + Terminal.app (Apple Silicon)
- ✅ macOS 11 (Big Sur) + iTerm2 (Intel)

### Linux
- ✅ Ubuntu 22.04 + GNOME Terminal
- ✅ Ubuntu 20.04 + GNOME Terminal
- ✅ Fedora 38 + Konsole
- ✅ Arch Linux + Alacritty
- ✅ Debian 11 + xterm

## Future Improvements

- [ ] Add automated visual regression testing
- [ ] Create Docker containers for testing different Linux distributions
- [ ] Add telemetry for platform usage statistics
- [ ] Improve Windows cmd.exe compatibility
- [ ] Add support for more terminal emulators
