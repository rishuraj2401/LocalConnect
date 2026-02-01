# Environment Setup

## Quick Start

### 1. Create `.env` file

```bash
cd backend
cat > .env << 'EOF'
DATABASE_URL=postgres://postgres:postgres@localhost:5430/localconnect?sslmode=disable
JWT_SECRET=dev-secret-change-in-production
MEDIA_DIR=./media
PORT=8080
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
docker exec -i localconnect-postgres psql -U postgres -d localconnect < migrations/001_init.sql
docker exec -i localconnect-postgres psql -U postgres -d localconnect < migrations/002_messaging.sql

# Start server
go run cmd/api/main.go
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/localconnect?sslmode=disable` |
| `JWT_SECRET` | Secret key for JWT tokens | `your-secret-key` |
| `MEDIA_DIR` | Directory for media uploads | `./media` |
| `PORT` | Server port | `8080` |

## Important Notes

- **Never commit `.env` to git** - it's already in `.gitignore`
- Change `JWT_SECRET` before deploying to production
- Make sure PostgreSQL is running on the correct port (check with `docker ps`)
