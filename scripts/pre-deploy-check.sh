#!/bin/bash

# Pre-Deployment Check Script
# Run this before deploying to catch common issues

set -e

echo "🚀 PokeTacTix Pre-Deployment Check"
echo "===================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check functions
check_pass() {
    echo -e "${GREEN}✓${NC} $1"
}

check_fail() {
    echo -e "${RED}✗${NC} $1"
    FAILED=1
}

check_warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

FAILED=0

echo "📦 Checking Repository..."
echo ""

# Check if git repo
if [ -d .git ]; then
    check_pass "Git repository initialized"
else
    check_fail "Not a git repository"
fi

# Check for uncommitted changes
if [ -z "$(git status --porcelain)" ]; then
    check_pass "No uncommitted changes"
else
    check_warn "You have uncommitted changes"
fi

# Check if on main/master branch
BRANCH=$(git branch --show-current)
if [ "$BRANCH" = "main" ] || [ "$BRANCH" = "master" ]; then
    check_pass "On main/master branch"
else
    check_warn "Not on main/master branch (current: $BRANCH)"
fi

echo ""
echo "🔧 Checking Configuration Files..."
echo ""

# Check for required files
if [ -f "netlify.toml" ]; then
    check_pass "netlify.toml exists"
else
    check_fail "netlify.toml missing"
fi

if [ -f "frontend/_redirects" ]; then
    check_pass "frontend/_redirects exists"
else
    check_fail "frontend/_redirects missing"
fi

if [ -f "frontend/.env.production" ]; then
    check_pass "frontend/.env.production exists"
else
    check_fail "frontend/.env.production missing"
fi

if [ -f ".env.production.example" ]; then
    check_pass ".env.production.example exists"
else
    check_warn ".env.production.example missing"
fi

echo ""
echo "🎨 Checking Assets..."
echo ""

if [ -f "frontend/public/assets/wallpaper.jpg" ]; then
    check_pass "wallpaper.jpg exists"
else
    check_fail "wallpaper.jpg missing"
fi

if [ -f "frontend/public/assets/pokeball.png" ]; then
    check_pass "pokeball.png exists"
else
    check_fail "pokeball.png missing"
fi

echo ""
echo "🗄️ Checking Database Migrations..."
echo ""

MIGRATION_COUNT=$(ls -1 internal/database/migrations/*.up.sql 2>/dev/null | wc -l)
if [ "$MIGRATION_COUNT" -gt 0 ]; then
    check_pass "Found $MIGRATION_COUNT migration files"
else
    check_fail "No migration files found"
fi

echo ""
echo "🔨 Checking Build..."
echo ""

# Check if Go is installed
if command -v go &> /dev/null; then
    check_pass "Go is installed ($(go version | awk '{print $3}'))"
    
    # Try to build backend
    echo "   Building backend..."
    if go build -o /tmp/poketactix-test ./cmd/api > /dev/null 2>&1; then
        check_pass "Backend builds successfully"
        rm -f /tmp/poketactix-test
    else
        check_fail "Backend build failed"
    fi
else
    check_warn "Go not installed (can't test backend build)"
fi

# Check if Node is installed
if command -v node &> /dev/null; then
    check_pass "Node is installed ($(node --version))"
    
    # Check if frontend dependencies are installed
    if [ -d "frontend/node_modules" ]; then
        check_pass "Frontend dependencies installed"
        
        # Try to build frontend
        echo "   Building frontend..."
        cd frontend
        if npm run build > /dev/null 2>&1; then
            check_pass "Frontend builds successfully"
        else
            check_fail "Frontend build failed"
        fi
        cd ..
    else
        check_warn "Frontend dependencies not installed (run: cd frontend && npm install)"
    fi
else
    check_warn "Node not installed (can't test frontend build)"
fi

echo ""
echo "🔐 Checking Security..."
echo ""

# Check if .env is in .gitignore
if grep -q "^\.env$" .gitignore 2>/dev/null; then
    check_pass ".env is in .gitignore"
else
    check_fail ".env not in .gitignore"
fi

# Check if .env exists (shouldn't be committed)
if [ -f ".env" ]; then
    if git ls-files --error-unmatch .env > /dev/null 2>&1; then
        check_fail ".env is tracked by git (should be ignored!)"
    else
        check_pass ".env exists but not tracked by git"
    fi
fi

echo ""
echo "📚 Checking Documentation..."
echo ""

if [ -f "docs/DEPLOYMENT_GUIDE.md" ]; then
    check_pass "Deployment guide exists"
else
    check_warn "Deployment guide missing"
fi

if [ -f "docs/QUICK_DEPLOY.md" ]; then
    check_pass "Quick deploy guide exists"
else
    check_warn "Quick deploy guide missing"
fi

if [ -f "README.md" ]; then
    check_pass "README.md exists"
else
    check_warn "README.md missing"
fi

echo ""
echo "===================================="

if [ $FAILED -eq 1 ]; then
    echo -e "${RED}❌ Pre-deployment check FAILED${NC}"
    echo "Please fix the issues above before deploying."
    exit 1
else
    echo -e "${GREEN}✅ Pre-deployment check PASSED${NC}"
    echo ""
    echo "You're ready to deploy! 🚀"
    echo ""
    echo "Next steps:"
    echo "1. Push to GitHub: git push origin main"
    echo "2. Follow the deployment guide: docs/DEPLOYMENT_GUIDE.md"
    echo "3. Or use quick deploy: docs/QUICK_DEPLOY.md"
    exit 0
fi
