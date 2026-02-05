#!/bin/bash
set -euo pipefail

# Load environment variables from .env file
if [ -f .env ]; then
    echo "Loading environment variables from .env..."
    set -a
    # shellcheck disable=SC1091
    source .env
    set +a
else
    echo "Warning: .env file not found!"
    exit 1
fi

SEED_DEMO_DATA="${SEED_DEMO_DATA:-true}"

run_migrations_with_psql() {
    echo "Running database migrations via DATABASE_URL..."
    psql "${DATABASE_URL}" -v ON_ERROR_STOP=1 -f migrations/001_init.sql
    psql "${DATABASE_URL}" -v ON_ERROR_STOP=1 -f migrations/002_messaging.sql
    if [ "${SEED_DEMO_DATA}" = "true" ]; then
        psql "${DATABASE_URL}" -v ON_ERROR_STOP=1 -f migrations/003_seed_data.sql
    else
        echo "Skipping demo seed data (SEED_DEMO_DATA=${SEED_DEMO_DATA})"
    fi
    echo "✅ Migrations completed!"
}

run_migrations_with_docker() {
    echo "Running database migrations via Docker container (localconnect-postgres)..."
    docker exec -i localconnect-postgres psql -U postgres -d localconnect -v ON_ERROR_STOP=1 < migrations/001_init.sql
    docker exec -i localconnect-postgres psql -U postgres -d localconnect -v ON_ERROR_STOP=1 < migrations/002_messaging.sql
    if [ "${SEED_DEMO_DATA}" = "true" ]; then
        docker exec -i localconnect-postgres psql -U postgres -d localconnect -v ON_ERROR_STOP=1 < migrations/003_seed_data.sql
    else
        echo "Skipping demo seed data (SEED_DEMO_DATA=${SEED_DEMO_DATA})"
    fi
    echo "✅ Migrations completed!"
}

# Prefer running migrations against DATABASE_URL (works for Supabase and local DB),
# but fall back to docker exec if psql isn't installed.
if command -v psql >/dev/null 2>&1; then
    run_migrations_with_psql
else
    echo "Warning: psql not found on PATH; falling back to Docker migrations."
    run_migrations_with_docker
fi

echo ""
echo "Starting LocalConnect API server..."
echo "API URL: http://localhost:${PORT}"
echo "Database: ${DATABASE_URL}"
echo ""

# Start the backend
go run cmd/api/main.go
