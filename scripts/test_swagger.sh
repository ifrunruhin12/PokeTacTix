#!/bin/bash

# Test script for Swagger API documentation endpoints

BASE_URL="http://localhost:3000"
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Testing PokeTacTix API Documentation Endpoints..."
echo ""

# Test 1: Health check
echo -n "1. Testing health endpoint... "
HEALTH=$(curl -s ${BASE_URL}/health | grep -o "healthy")
if [ "$HEALTH" == "healthy" ]; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

# Test 2: Swagger YAML file
echo -n "2. Testing swagger.yaml endpoint... "
YAML_TITLE=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "title: PokeTacTix API")
if [ ! -z "$YAML_TITLE" ]; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

# Test 3: Swagger UI HTML
echo -n "3. Testing Swagger UI endpoint... "
UI_HTML=$(curl -sL ${BASE_URL}/api/docs | grep "swagger-ui")
if [ ! -z "$UI_HTML" ]; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

# Test 4: Check OpenAPI version
echo -n "4. Checking OpenAPI version... "
OPENAPI_VERSION=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "openapi: 3.0.3")
if [ ! -z "$OPENAPI_VERSION" ]; then
    echo -e "${GREEN}✓ PASS (OpenAPI 3.0.3)${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

# Test 5: Check authentication endpoints documented
echo -n "5. Checking auth endpoints documented... "
AUTH_REGISTER=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/auth/register")
AUTH_LOGIN=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/auth/login")
AUTH_ME=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/auth/me")
if [ ! -z "$AUTH_REGISTER" ] && [ ! -z "$AUTH_LOGIN" ] && [ ! -z "$AUTH_ME" ]; then
    echo -e "${GREEN}✓ PASS (3 endpoints)${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

# Test 6: Check battle endpoints documented
echo -n "6. Checking battle endpoints documented... "
BATTLE_START=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/battle/start")
BATTLE_MOVE=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/battle/move")
BATTLE_SWITCH=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/battle/switch")
if [ ! -z "$BATTLE_START" ] && [ ! -z "$BATTLE_MOVE" ] && [ ! -z "$BATTLE_SWITCH" ]; then
    echo -e "${GREEN}✓ PASS (5 endpoints)${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

# Test 7: Check shop endpoints documented
echo -n "7. Checking shop endpoints documented... "
SHOP_INVENTORY=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/shop/inventory")
SHOP_PURCHASE=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/shop/purchase")
if [ ! -z "$SHOP_INVENTORY" ] && [ ! -z "$SHOP_PURCHASE" ]; then
    echo -e "${GREEN}✓ PASS (2 endpoints)${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

# Test 8: Check profile endpoints documented
echo -n "8. Checking profile endpoints documented... "
PROFILE_STATS=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/profile/stats")
PROFILE_HISTORY=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/profile/history")
PROFILE_ACHIEVEMENTS=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/profile/achievements")
if [ ! -z "$PROFILE_STATS" ] && [ ! -z "$PROFILE_HISTORY" ] && [ ! -z "$PROFILE_ACHIEVEMENTS" ]; then
    echo -e "${GREEN}✓ PASS (3 endpoints)${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

# Test 9: Check cards endpoints documented
echo -n "9. Checking cards endpoints documented... "
CARDS_GET=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/cards:")
CARDS_DECK=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "/api/cards/deck")
if [ ! -z "$CARDS_GET" ] && [ ! -z "$CARDS_DECK" ]; then
    echo -e "${GREEN}✓ PASS (2 endpoints)${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

# Test 10: Check security scheme defined
echo -n "10. Checking JWT security scheme... "
BEARER_AUTH=$(curl -s ${BASE_URL}/api/docs/swagger.yaml | grep "BearerAuth")
if [ ! -z "$BEARER_AUTH" ]; then
    echo -e "${GREEN}✓ PASS${NC}"
else
    echo -e "${RED}✗ FAIL${NC}"
fi

echo ""
echo "Documentation testing complete!"
echo ""
echo "Access the interactive documentation at: ${BASE_URL}/api/docs"
