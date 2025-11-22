#!/bin/bash

# PokeTacTix CLI Release Packaging Script
# Creates release packages with checksums for all platforms

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║                                                        ║"
echo "║         PokeTacTix CLI Release Packager                ║"
echo "║                                                        ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""

# Get version from git or prompt
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "")
if [ -z "$VERSION" ]; then
    echo -e "${YELLOW}No git tags found. Please enter version (e.g., v1.0.0):${NC}"
    read -r VERSION
fi

echo -e "${BLUE}Creating release for version: $VERSION${NC}"
echo ""

# Create release directory
RELEASE_DIR="release/$VERSION"
mkdir -p "$RELEASE_DIR"

# Build all binaries
echo -e "${BLUE}Building binaries for all platforms...${NC}"
make build-cli-all

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Build complete${NC}"
echo ""

# Copy binaries to release directory
echo -e "${BLUE}Packaging binaries...${NC}"
cp bin/poketactix-cli-windows-amd64.exe "$RELEASE_DIR/"
cp bin/poketactix-cli-darwin-amd64 "$RELEASE_DIR/"
cp bin/poketactix-cli-darwin-arm64 "$RELEASE_DIR/"
cp bin/poketactix-cli-linux-amd64 "$RELEASE_DIR/"
cp bin/poketactix-cli-linux-arm64 "$RELEASE_DIR/"

# Copy documentation
echo -e "${BLUE}Adding documentation...${NC}"
cp CLI_README.md "$RELEASE_DIR/README.md"
cp LICENSE "$RELEASE_DIR/" 2>/dev/null || echo "Note: LICENSE file not found"

# Create installation scripts in release
cp scripts/install.sh "$RELEASE_DIR/"
cp scripts/install.ps1 "$RELEASE_DIR/"

# Generate checksums
echo -e "${BLUE}Generating checksums...${NC}"
cd "$RELEASE_DIR"

# SHA256 checksums
if command -v sha256sum &> /dev/null; then
    sha256sum poketactix-cli-* > SHA256SUMS
elif command -v shasum &> /dev/null; then
    shasum -a 256 poketactix-cli-* > SHA256SUMS
else
    echo -e "${YELLOW}Warning: sha256sum not found, skipping checksums${NC}"
fi

# MD5 checksums (for compatibility)
if command -v md5sum &> /dev/null; then
    md5sum poketactix-cli-* > MD5SUMS
elif command -v md5 &> /dev/null; then
    md5 poketactix-cli-* > MD5SUMS
fi

cd - > /dev/null

echo -e "${GREEN}✓ Checksums generated${NC}"
echo ""

# Create release notes template
RELEASE_NOTES="$RELEASE_DIR/RELEASE_NOTES.md"
cat > "$RELEASE_NOTES" << EOF
# PokeTacTix CLI $VERSION

## What's New

<!-- Add your release notes here -->

- Feature 1
- Feature 2
- Bug fix 1

## Installation

### Quick Install

**macOS/Linux:**
\`\`\`bash
curl -L https://github.com/yourusername/poketactix/raw/main/scripts/install.sh | bash
\`\`\`

**Windows:**
\`\`\`powershell
powershell -ExecutionPolicy Bypass -Command "iwr https://github.com/yourusername/poketactix/raw/main/scripts/install.ps1 | iex"
\`\`\`

### Manual Download

Download the appropriate binary for your platform:

- **Windows (64-bit)**: \`poketactix-cli-windows-amd64.exe\`
- **macOS (Intel)**: \`poketactix-cli-darwin-amd64\`
- **macOS (Apple Silicon)**: \`poketactix-cli-darwin-arm64\`
- **Linux (64-bit)**: \`poketactix-cli-linux-amd64\`
- **Linux (ARM64)**: \`poketactix-cli-linux-arm64\`

## Checksums

Verify your download with SHA256:

\`\`\`
$(cat "$RELEASE_DIR/SHA256SUMS" 2>/dev/null || echo "See SHA256SUMS file")
\`\`\`

## Full Changelog

See [CHANGELOG.md](CHANGELOG.md) for complete details.

## System Requirements

- **OS**: Windows 10+, macOS 10.12+, or Linux
- **Terminal**: 80x24 minimum, ANSI color support recommended
- **Disk Space**: ~10 MB

## Known Issues

<!-- List any known issues -->

- None

## Credits

Built with ❤️ for Pokemon fans everywhere!
EOF

echo -e "${BLUE}Release notes template created at: $RELEASE_NOTES${NC}"
echo ""

# Create archive for each platform
echo -e "${BLUE}Creating platform archives...${NC}"

cd "$RELEASE_DIR"

# Windows
zip -q "poketactix-cli-windows-amd64-$VERSION.zip" \
    poketactix-cli-windows-amd64.exe \
    README.md \
    install.ps1 \
    LICENSE 2>/dev/null || true

# macOS Intel
tar -czf "poketactix-cli-darwin-amd64-$VERSION.tar.gz" \
    poketactix-cli-darwin-amd64 \
    README.md \
    install.sh \
    LICENSE 2>/dev/null || true

# macOS ARM
tar -czf "poketactix-cli-darwin-arm64-$VERSION.tar.gz" \
    poketactix-cli-darwin-arm64 \
    README.md \
    install.sh \
    LICENSE 2>/dev/null || true

# Linux AMD64
tar -czf "poketactix-cli-linux-amd64-$VERSION.tar.gz" \
    poketactix-cli-linux-amd64 \
    README.md \
    install.sh \
    LICENSE 2>/dev/null || true

# Linux ARM64
tar -czf "poketactix-cli-linux-arm64-$VERSION.tar.gz" \
    poketactix-cli-linux-arm64 \
    README.md \
    install.sh \
    LICENSE 2>/dev/null || true

cd - > /dev/null

echo -e "${GREEN}✓ Archives created${NC}"
echo ""

# Display summary
echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║                                                        ║"
echo "║         ✓ Release package created successfully!       ║"
echo "║                                                        ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
echo -e "${GREEN}Release directory: $RELEASE_DIR${NC}"
echo ""
echo "Contents:"
ls -lh "$RELEASE_DIR"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Edit release notes: $RELEASE_NOTES"
echo "2. Test binaries on each platform"
echo "3. Create GitHub release with tag: $VERSION"
echo "4. Upload archives and checksums to GitHub"
echo ""
echo -e "${BLUE}GitHub Release Command:${NC}"
echo "gh release create $VERSION \\"
echo "  $RELEASE_DIR/*.zip \\"
echo "  $RELEASE_DIR/*.tar.gz \\"
echo "  $RELEASE_DIR/SHA256SUMS \\"
echo "  $RELEASE_DIR/MD5SUMS \\"
echo "  --title \"PokeTacTix CLI $VERSION\" \\"
echo "  --notes-file $RELEASE_NOTES"
echo ""
