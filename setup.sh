#!/bin/bash

# LocalConnect Setup Script
set -e

echo "🚀 Setting up LocalConnect..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check prerequisites
echo -e "${BLUE}Checking prerequisites...${NC}"

if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed${NC}"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    exit 1
fi

if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js is not installed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ All prerequisites found${NC}"

# Start PostgreSQL
echo -e "${BLUE}Starting PostgreSQL...${NC}"
docker compose up -d

# Wait for PostgreSQL to be ready
echo -e "${BLUE}Waiting for PostgreSQL to be ready...${NC}"
sleep 5

# Run migrations
echo -e "${BLUE}Running database migrations...${NC}"
docker exec -i localconnect-postgres psql -U postgres -d localconnect < backend/migrations/001_init.sql || true
docker exec -i localconnect-postgres psql -U postgres -d localconnect < backend/migrations/002_messaging.sql || true

echo -e "${GREEN}✅ Database initialized${NC}"

# Setup backend
echo -e "${BLUE}Setting up backend...${NC}"
cd backend

if [ ! -f "go.mod" ]; then
    echo -e "${RED}go.mod not found${NC}"
    exit 1
fi

go mod download
echo -e "${GREEN}✅ Backend dependencies installed${NC}"

cd ..

# Setup frontend
echo -e "${BLUE}Setting up frontend...${NC}"
cd frontend

npm install
echo -e "${GREEN}✅ Frontend dependencies installed${NC}"

cd ..

# Create media directory
echo -e "${BLUE}Creating media directory...${NC}"
mkdir -p backend/media
echo -e "${GREEN}✅ Media directory created${NC}"

echo ""
echo -e "${GREEN}✅ Setup complete!${NC}"
echo ""
echo -e "${BLUE}To start the application:${NC}"
echo ""
echo -e "  1. Start the backend:"
echo -e "     ${GREEN}cd backend && go run ./cmd/api${NC}"
echo ""
echo -e "  2. In a new terminal, start the frontend:"
echo -e "     ${GREEN}cd frontend && npm run dev${NC}"
echo ""
echo -e "  3. Open your browser to:"
echo -e "     ${GREEN}http://localhost:5173${NC}"
echo ""
echo -e "${BLUE}API will be available at:${NC}"
echo -e "     ${GREEN}http://localhost:8080${NC}"
echo ""
