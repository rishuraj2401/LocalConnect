# Environment Setup

## Quick Start

### 1. Create `.env` file

```bash
cd backend
cat > .env << 'EOF'
# Local Postgres (docker-compose.yml maps host 5430 -> container 5432)
DATABASE_URL=postgres://postgres:postgres@localhost:5430/localconnect?sslmode=disable
JWT_SECRET=dev-secret-change-in-production
MEDIA_DIR=./media
PORT=8080
# Optional: set to false in production
SEED_DEMO_DATA=true
EOF
```

### 2. Start the backend

**Option A: Use the startup script (Recommended)**
```bash
./start.sh
```

**Option B: Manual start**
```bash
# Load environment variables
source .env  # or: export $(cat .env | xargs)

# Run migrations
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f migrations/001_init.sql
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f migrations/002_messaging.sql
if [ "${SEED_DEMO_DATA:-true}" = "true" ]; then
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f migrations/003_seed_data.sql
fi

# Start server
go run cmd/api/main.go
```

## Using Supabase instead of local Postgres

1. Create a project in Supabase.
2. Copy the **Postgres connection string** from Supabase (Project Settings → Database).
3. Put it in `backend/.env` as `DATABASE_URL`, and ensure SSL is enabled (usually `sslmode=require`).

Example:

```bash
DATABASE_URL=postgresql://postgres:<password>@db.<project-ref>.supabase.co:5432/postgres?sslmode=require
```

Then run `./start.sh` (it will run migrations using `psql "$DATABASE_URL"` if `psql` is installed).

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/localconnect?sslmode=disable` |
| `JWT_SECRET` | Secret key for JWT tokens | `your-secret-key` |
| `MEDIA_DIR` | Directory for media uploads | `./media` |
| `PORT` | Server port | `8080` |
| `SEED_DEMO_DATA` | Seed demo data on startup | `true` |

## Important Notes

- **Never commit `.env` to git** - it's already in `.gitignore`
- Change `JWT_SECRET` before deploying to production
- Make sure PostgreSQL is running on the correct port (check with `docker ps`)
