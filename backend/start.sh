#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    echo "Loading environment variables from .env..."
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "Warning: .env file not found!"
    exit 1
fi

# Run migrations
echo "Running database migrations..."
docker exec -i localconnect-postgres psql -U postgres -d localconnect < migrations/001_init.sql 2>&1 | grep -v "already exists" || true
docker exec -i localconnect-postgres psql -U postgres -d localconnect < migrations/002_messaging.sql 2>&1 | grep -v "already exists" || true
docker exec -i localconnect-postgres psql -U postgres -d localconnect < migrations/003_seed_data.sql 2>&1 | grep -v "already exists" || true
echo "✅ Migrations completed!"

echo ""
echo "Starting LocalConnect API server..."
echo "API URL: http://localhost:${PORT}"
echo "Database: ${DATABASE_URL}"
echo ""

# Start the backend
go run cmd/api/main.go
