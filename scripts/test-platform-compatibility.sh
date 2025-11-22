#!/bin/bash

# Platform Compatibility Test Script
# Tests cross-platform features of PokeTacTix CLI

set -e

echo "=========================================="
echo "PokeTacTix CLI Platform Compatibility Test"
echo "=========================================="
echo ""

# Detect platform
OS=$(uname -s)
ARCH=$(uname -m)

echo "Platform Information:"
echo "  OS: $OS"
echo "  Architecture: $ARCH"
echo "  Shell: $SHELL"
echo ""

# Detect terminal
echo "Terminal Information:"
echo "  TERM: ${TERM:-not set}"
echo "  COLORTERM: ${COLORTERM:-not set}"
echo "  TERM_PROGRAM: ${TERM_PROGRAM:-not set}"
echo "  Terminal Size: ${COLUMNS:-?}x${LINES:-?}"
echo ""

# Test color support
echo "Color Support Tests:"
echo "  NO_COLOR: ${NO_COLOR:-not set}"
echo "  CLICOLOR: ${CLICOLOR:-not set}"
echo "  CLICOLOR_FORCE: ${CLICOLOR_FORCE:-not set}"

# Test ANSI colors
if [ -z "$NO_COLOR" ]; then
    echo -e "  \033[31mRed\033[0m \033[32mGreen\033[0m \033[33mYellow\033[0m \033[34mBlue\033[0m - Colors working!"
else
    echo "  Colors disabled (NO_COLOR set)"
fi
echo ""

# Test file paths
echo "File Path Tests:"
HOME_DIR="${HOME:-$USERPROFILE}"
echo "  Home Directory: $HOME_DIR"

SAVE_DIR="$HOME_DIR/.poketactix"
echo "  Save Directory: $SAVE_DIR"

if [ -d "$SAVE_DIR" ]; then
    echo "  ✓ Save directory exists"
    
    # Check permissions (Unix-like only)
    if [ "$OS" != "MINGW"* ] && [ "$OS" != "MSYS"* ]; then
        PERMS=$(stat -c "%a" "$SAVE_DIR" 2>/dev/null || stat -f "%Lp" "$SAVE_DIR" 2>/dev/null || echo "unknown")
        echo "  Directory permissions: $PERMS"
    fi
else
    echo "  ℹ Save directory does not exist yet (will be created on first run)"
fi
echo ""

# Test Go environment
echo "Go Environment:"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    echo "  ✓ Go installed: $GO_VERSION"
    
    # Test build for current platform
    echo ""
    echo "Testing build for current platform..."
    
    BUILD_OUTPUT="bin/poketactix-cli-test"
    if [ "$OS" = "MINGW"* ] || [ "$OS" = "MSYS"* ]; then
        BUILD_OUTPUT="bin/poketactix-cli-test.exe"
    fi
    
    if go build -o "$BUILD_OUTPUT" cmd/cli/main.go 2>&1; then
        echo "  ✓ Build successful: $BUILD_OUTPUT"
        
        # Check binary size
        if [ -f "$BUILD_OUTPUT" ]; then
            SIZE=$(du -h "$BUILD_OUTPUT" | cut -f1)
            echo "  Binary size: $SIZE"
        fi
        
        # Clean up test binary
        rm -f "$BUILD_OUTPUT"
    else
        echo "  ✗ Build failed"
        exit 1
    fi
else
    echo "  ✗ Go not installed"
    exit 1
fi
echo ""

# Run Go tests
echo "Running Go Tests:"
echo "  Testing UI package..."
if go test -v ./internal/cli/ui/... 2>&1 | grep -E "(PASS|FAIL|ok|FAIL)"; then
    echo "  ✓ UI tests completed"
else
    echo "  ✗ UI tests failed"
fi

echo ""
echo "  Testing storage package..."
if go test -v ./internal/cli/storage/... 2>&1 | grep -E "(PASS|FAIL|ok|FAIL)"; then
    echo "  ✓ Storage tests completed"
else
    echo "  ✗ Storage tests failed"
fi

echo ""
echo "  Testing commands package..."
if go test -v ./internal/cli/commands/... 2>&1 | grep -E "(PASS|FAIL|ok|FAIL)"; then
    echo "  ✓ Commands tests completed"
else
    echo "  ✗ Commands tests failed"
fi

echo ""
echo "=========================================="
echo "Platform Compatibility Test Complete"
echo "=========================================="
echo ""

# Platform-specific notes
case "$OS" in
    Linux*)
        echo "Linux Notes:"
        echo "  - Ensure terminal supports ANSI colors"
        echo "  - Save directory: ~/.poketactix/"
        echo "  - Recommended terminals: GNOME Terminal, Konsole, Alacritty"
        ;;
    Darwin*)
        echo "macOS Notes:"
        echo "  - Save directory: ~/.poketactix/"
        echo "  - Recommended terminals: iTerm2, Terminal.app"
        echo "  - On first run, you may need to allow the binary in System Preferences"
        ;;
    MINGW*|MSYS*|CYGWIN*)
        echo "Windows Notes:"
        echo "  - Save directory: %USERPROFILE%\\.poketactix\\"
        echo "  - Recommended terminals: Windows Terminal, PowerShell"
        echo "  - cmd.exe has limited color support"
        ;;
    *)
        echo "Unknown platform: $OS"
        ;;
esac

echo ""
echo "For detailed platform testing guide, see: docs/cli-platform-testing.md"
