#!/bin/bash

# PokeTacTix CLI Installation Script
# Supports macOS and Linux

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="yourusername/poketactix"  # Update this with actual GitHub repo
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="poketactix-cli"

echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║                                                        ║"
echo "║         PokeTacTix CLI Installation Script            ║"
echo "║                                                        ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}✗ Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

case "$OS" in
    darwin)
        PLATFORM="darwin"
        ;;
    linux)
        PLATFORM="linux"
        ;;
    *)
        echo -e "${RED}✗ Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

BINARY_FILE="${BINARY_NAME}-${PLATFORM}-${ARCH}"
echo -e "${BLUE}Detected platform: ${PLATFORM}/${ARCH}${NC}"
echo ""

# Check if running with sudo for system-wide install
if [ "$EUID" -ne 0 ] && [ -w "$INSTALL_DIR" ] || [ "$EUID" -eq 0 ]; then
    SYSTEM_INSTALL=true
    TARGET_DIR="$INSTALL_DIR"
else
    SYSTEM_INSTALL=false
    TARGET_DIR="$HOME/.local/bin"
    echo -e "${YELLOW}Note: Installing to user directory ($TARGET_DIR)${NC}"
    echo -e "${YELLOW}Run with sudo for system-wide installation${NC}"
    echo ""
fi

# Create target directory if it doesn't exist
if [ ! -d "$TARGET_DIR" ]; then
    echo -e "${BLUE}Creating directory: $TARGET_DIR${NC}"
    mkdir -p "$TARGET_DIR"
fi

# Download latest release
echo -e "${BLUE}Fetching latest release...${NC}"
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo -e "${RED}✗ Failed to fetch latest release${NC}"
    echo -e "${YELLOW}Please check your internet connection and repository name${NC}"
    exit 1
fi

echo -e "${GREEN}Latest version: $LATEST_RELEASE${NC}"
echo ""

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$BINARY_FILE"
TEMP_FILE="/tmp/$BINARY_FILE"

echo -e "${BLUE}Downloading from: $DOWNLOAD_URL${NC}"
if curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"; then
    echo -e "${GREEN}✓ Download complete${NC}"
else
    echo -e "${RED}✗ Download failed${NC}"
    exit 1
fi

# Make binary executable
chmod +x "$TEMP_FILE"

# Move to target directory
echo -e "${BLUE}Installing to: $TARGET_DIR/$BINARY_NAME${NC}"
if mv "$TEMP_FILE" "$TARGET_DIR/$BINARY_NAME"; then
    echo -e "${GREEN}✓ Installation complete${NC}"
else
    echo -e "${RED}✗ Installation failed${NC}"
    exit 1
fi

# Check if target directory is in PATH
if [[ ":$PATH:" != *":$TARGET_DIR:"* ]]; then
    echo ""
    echo -e "${YELLOW}⚠ Warning: $TARGET_DIR is not in your PATH${NC}"
    echo ""
    echo "Add the following line to your shell configuration file:"
    echo "  (~/.bashrc, ~/.zshrc, or ~/.profile)"
    echo ""
    echo -e "${BLUE}export PATH=\"\$PATH:$TARGET_DIR\"${NC}"
    echo ""
fi

# Display success message
echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║                                                        ║"
echo "║         ✓ PokeTacTix CLI installed successfully!      ║"
echo "║                                                        ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
echo -e "${GREEN}Run the game with: ${BINARY_NAME}${NC}"
echo -e "${GREEN}Check version with: ${BINARY_NAME} --version${NC}"
echo ""
