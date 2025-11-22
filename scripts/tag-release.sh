#!/bin/bash

# PokeTacTix CLI Release Tagging Script
# Creates a git tag which triggers the GitHub Actions release workflow

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
echo "║         PokeTacTix CLI Release Tagger                  ║"
echo "║                                                        ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}✗ Not in a git repository${NC}"
    exit 1
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD -- 2>/dev/null; then
    echo -e "${YELLOW}⚠ Warning: You have uncommitted changes${NC}"
    echo ""
    git status --short
    echo ""
    read -p "Continue anyway? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${RED}Aborted${NC}"
        exit 1
    fi
fi

# Get current tags
echo -e "${BLUE}Recent tags:${NC}"
git tag -l --sort=-v:refname | head -5
echo ""

# Prompt for version
echo -e "${YELLOW}Enter new version (e.g., v1.0.0):${NC}"
read -r VERSION

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
    echo -e "${RED}✗ Invalid version format${NC}"
    echo "Expected format: v1.0.0 or v1.0.0-beta.1"
    exit 1
fi

# Check if tag already exists
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo -e "${RED}✗ Tag $VERSION already exists${NC}"
    exit 1
fi

# Prompt for release notes
echo ""
echo -e "${YELLOW}Enter release notes (press Ctrl+D when done):${NC}"
NOTES=$(cat)

if [ -z "$NOTES" ]; then
    NOTES="Release $VERSION"
fi

# Confirm
echo ""
echo -e "${BLUE}Ready to create and push tag:${NC}"
echo "  Version: $VERSION"
echo "  Notes: $NOTES"
echo ""
read -p "Continue? [y/N] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}Aborted${NC}"
    exit 1
fi

# Create annotated tag
echo ""
echo -e "${BLUE}Creating tag...${NC}"
git tag -a "$VERSION" -m "$NOTES"

# Push tag
echo -e "${BLUE}Pushing tag to origin...${NC}"
git push origin "$VERSION"

echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║                                                        ║"
echo "║         ✓ Tag created and pushed successfully!        ║"
echo "║                                                        ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
echo -e "${GREEN}Tag: $VERSION${NC}"
echo ""
echo -e "${BLUE}GitHub Actions will now:${NC}"
echo "  1. Build binaries for all platforms"
echo "  2. Create release archives"
echo "  3. Generate checksums"
echo "  4. Create GitHub release"
echo "  5. Upload all assets"
echo ""
echo -e "${YELLOW}Monitor progress at:${NC}"
echo "  https://github.com/$(git remote get-url origin | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/actions"
echo ""
echo -e "${BLUE}View release when ready at:${NC}"
echo "  https://github.com/$(git remote get-url origin | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/releases/tag/$VERSION"
echo ""
