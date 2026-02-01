#!/bin/bash

# LocalConnect Test Runner
set -e

echo "🧪 Running LocalConnect Tests..."

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if PostgreSQL is running
echo -e "${BLUE}Checking PostgreSQL...${NC}"
if ! docker ps | grep -q localconnect-postgres; then
    echo -e "${RED}❌ PostgreSQL container is not running${NC}"
    echo -e "${YELLOW}Starting PostgreSQL...${NC}"
    cd ..
    docker compose up -d
    cd backend
    sleep 3
fi

echo -e "${GREEN}✅ PostgreSQL is running${NC}"

# Setup test database
echo -e "${BLUE}Setting up test database...${NC}"
docker exec -i localconnect-postgres psql -U postgres < test_setup.sql 2>/dev/null || {
    echo -e "${YELLOW}Note: Test database setup may have warnings (this is normal)${NC}"
}
echo -e "${GREEN}✅ Test database ready${NC}"

# Run tests
echo -e "${BLUE}Running tests...${NC}"
echo ""

# Run tests with coverage
go test -v -coverprofile=coverage.out ./internal/handlers/... || {
    echo -e "${RED}❌ Tests failed${NC}"
    exit 1
}

echo ""
echo -e "${GREEN}✅ All tests passed!${NC}"
echo ""

# Generate coverage report
echo -e "${BLUE}Generating coverage report...${NC}"
go tool cover -html=coverage.out -o coverage.html
echo -e "${GREEN}✅ Coverage report: coverage.html${NC}"

# Show coverage summary
echo ""
echo -e "${BLUE}Coverage Summary:${NC}"
go tool cover -func=coverage.out | tail -n 1

echo ""
echo -e "${GREEN}🎉 Testing complete!${NC}"
