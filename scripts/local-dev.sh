#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                                                        ║${NC}"
echo -e "${BLUE}║        💻 PokeTacTix Local Development Setup          ║${NC}"
echo -e "${BLUE}║                                                        ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

echo -e "${YELLOW}This script will guide you through setting up PokeTacTix locally${NC}"
echo -e "${YELLOW}(without Docker)${NC}"
echo ""

# Check PostgreSQL
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}Step 1: PostgreSQL Database${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

if command -v psql &> /dev/null; then
    echo -e "${GREEN}✅ PostgreSQL is installed${NC}"
    
    # Check if PostgreSQL is running
    if sudo systemctl is-active --quiet postgresql 2>/dev/null || pg_isready &> /dev/null; then
        echo -e "${GREEN}✅ PostgreSQL is running${NC}"
    else
        echo -e "${YELLOW}⚠️  PostgreSQL is not running${NC}"
        echo -e "${YELLOW}Starting PostgreSQL...${NC}"
        
        # Try to start PostgreSQL
        if command -v systemctl &> /dev/null; then
            sudo systemctl start postgresql
            if [ $? -eq 0 ]; then
                echo -e "${GREEN}✅ PostgreSQL started successfully${NC}"
            else
                echo -e "${RED}❌ Failed to start PostgreSQL${NC}"
                echo -e "${YELLOW}Please start it manually: sudo systemctl start postgresql${NC}"
            fi
        else
            echo -e "${YELLOW}Please start PostgreSQL manually${NC}"
        fi
    fi
else
    echo -e "${RED}❌ PostgreSQL is not installed${NC}"
    echo ""
    echo -e "${YELLOW}Install PostgreSQL:${NC}"
    echo -e "  ${BLUE}Ubuntu/Debian:${NC} sudo apt install postgresql postgresql-contrib"
    echo -e "  ${BLUE}Fedora:${NC}        sudo dnf install postgresql postgresql-server"
    echo -e "  ${BLUE}Arch:${NC}          sudo pacman -S postgresql"
    echo -e "  ${BLUE}macOS:${NC}         brew install postgresql"
    echo ""
    exit 1
fi

echo ""

# Check if database exists
echo -e "${YELLOW}Checking if 'poketactix' database exists...${NC}"
if sudo -u postgres psql -lqt 2>/dev/null | cut -d \| -f 1 | grep -qw poketactix; then
    echo -e "${GREEN}✅ Database 'poketactix' exists${NC}"
else
    echo -e "${YELLOW}⚠️  Database 'poketactix' does not exist${NC}"
    echo -e "${YELLOW}Creating database...${NC}"
    
    # Create user and database
    sudo -u postgres psql -c "CREATE USER pokemon WITH PASSWORD 'pokemon123';" 2>/dev/null
    sudo -u postgres psql -c "CREATE DATABASE poketactix OWNER pokemon;" 2>/dev/null
    sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE poketactix TO pokemon;" 2>/dev/null
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Database created successfully${NC}"
    else
        echo -e "${YELLOW}⚠️  Could not create database automatically${NC}"
        echo ""
        echo -e "${YELLOW}Create it manually:${NC}"
        echo -e "  sudo -u postgres psql"
        echo -e "  CREATE USER pokemon WITH PASSWORD 'pokemon123';"
        echo -e "  CREATE DATABASE poketactix OWNER pokemon;"
        echo -e "  GRANT ALL PRIVILEGES ON DATABASE poketactix TO pokemon;"
        echo -e "  \\q"
        echo ""
    fi
fi

echo ""

# Check Go
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}Step 2: Go Backend${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✅ Go is installed ($GO_VERSION)${NC}"
else
    echo -e "${RED}❌ Go is not installed${NC}"
    echo ""
    echo -e "${YELLOW}Install Go:${NC}"
    echo -e "  Visit: https://go.dev/doc/install"
    echo ""
    exit 1
fi

echo ""

# Check Node.js
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}Step 3: Node.js Frontend${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    echo -e "${GREEN}✅ Node.js is installed ($NODE_VERSION)${NC}"
else
    echo -e "${RED}❌ Node.js is not installed${NC}"
    echo ""
    echo -e "${YELLOW}Install Node.js:${NC}"
    echo -e "  Visit: https://nodejs.org/"
    echo -e "  Or use nvm: https://github.com/nvm-sh/nvm"
    echo ""
    exit 1
fi

# Check if frontend dependencies are installed
if [ -d "frontend/node_modules" ]; then
    echo -e "${GREEN}✅ Frontend dependencies are installed${NC}"
else
    echo -e "${YELLOW}⚠️  Frontend dependencies not installed${NC}"
    echo -e "${YELLOW}Installing dependencies...${NC}"
    cd frontend && npm install && cd ..
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Dependencies installed successfully${NC}"
    else
        echo -e "${RED}❌ Failed to install dependencies${NC}"
        exit 1
    fi
fi

echo ""

# Set environment variables
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}Step 4: Environment Setup${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

export DATABASE_URL="postgresql://pokemon:pokemon123@localhost:5432/poketactix?sslmode=disable"
export JWT_SECRET="dev-secret-key-change-in-production-min-256-bits"
export JWT_EXPIRATION="24h"
export PORT="3000"
export ENV="development"
export CORS_ORIGINS="http://localhost:5173"

echo -e "${GREEN}✅ Environment variables set${NC}"
echo ""

# Summary
echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                                                        ║${NC}"
echo -e "${BLUE}║              ✅ Setup Complete! ✅                     ║${NC}"
echo -e "${BLUE}║                                                        ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}🚀 To start development:${NC}"
echo ""
echo -e "${YELLOW}Terminal 1 - Backend:${NC}"
echo -e "  ${BLUE}export DATABASE_URL=\"postgresql://pokemon:pokemon123@localhost:5432/poketactix?sslmode=disable\"${NC}"
echo -e "  ${BLUE}go run cmd/api/main.go${NC}"
echo ""
echo -e "${YELLOW}Terminal 2 - Frontend:${NC}"
echo -e "  ${BLUE}cd frontend${NC}"
echo -e "  ${BLUE}npm run dev${NC}"
echo ""
echo -e "${GREEN}📱 Access your application:${NC}"
echo -e "   ${BLUE}Frontend:${NC}     http://localhost:5173"
echo -e "   ${BLUE}Backend API:${NC}  http://localhost:3000"
echo -e "   ${BLUE}API Docs:${NC}     http://localhost:3000/api/docs"
echo ""
echo -e "${GREEN}🗄️  Database Connection:${NC}"
echo -e "   ${BLUE}Host:${NC}         localhost"
echo -e "   ${BLUE}Port:${NC}         5432"
echo -e "   ${BLUE}Database:${NC}     poketactix"
echo -e "   ${BLUE}User:${NC}         pokemon"
echo -e "   ${BLUE}Password:${NC}     pokemon123"
echo ""
echo -e "${YELLOW}💡 Tip: Use 'make local' to run this script again${NC}"
echo ""
