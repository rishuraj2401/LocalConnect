# Quick Test Fix Guide

## Issues Found and Fixed

### 1. Database Connection
- ✅ Added connection ping to verify database is accessible
- ✅ Better error messages for database connection issues

### 2. Authentication Context
- ✅ Created `WithTestAuth()` helper function
- ✅ Properly sets user ID and role in context
- ✅ Works with all authenticated endpoints

### 3. URL Encoding
- ✅ Removed spaces from test URLs (NewYork instead of "New York")
- ✅ Fixed malformed HTTP version error

## Running Tests Now

### Step 1: Check PostgreSQL Port

Your PostgreSQL is running on port **543** not **5432**. Update the connection string:

```bash
# Check your docker-compose.yml
cd /Users/rishurajsinha/Desktop/Concepts/LocalConnect

# It shows: "0.0.0.0:543->5432/tcp"
# This means Docker maps host port 543 to container port 5432
```

### Step 2: Update Test Database URL

```bash
cd backend
```

Edit the test database URL in `auth_test.go`:

```go
// Change from:
dbURL := "postgres://postgres:postgres@localhost:5432/localconnect_test?sslmode=disable"

// To:
dbURL := "postgres://postgres:postgres@localhost:543/localconnect_test?sslmode=disable"
```

Or update your `docker-compose.yml` to use standard port:

```yaml
ports:
  - "5432:5432"  # Change from 543:5432
```

### Step 3: Run Tests

```bash
# Restart PostgreSQL with correct port
cd ..
docker compose down
docker compose up -d

# Wait for it to start
sleep 3

# Setup test database  
cd backend
make test-setup

# Run tests
make test
```

## Alternative: Fix Docker Port

If you want to use standard PostgreSQL port 5432:

```bash
# Stop current container
docker compose down

# Edit docker-compose.yml and change:
ports:
  - "5432:5432"  # Standard port

# Start again
docker compose up -d

# Now tests will work with default connection string
cd backend
make test-setup
make test
```

## Verifying the Fix

### Check PostgreSQL Connection

```bash
# Test connection manually
psql -h localhost -p 543 -U postgres -d localconnect_test

# Or with standard port
psql -h localhost -p 5432 -U postgres -d localconnect_test
```

### Run Single Test

```bash
# Test one endpoint
go test -v -run TestListCategories ./internal/handlers/
```

## Summary of Changes Made

1. **test_helpers.go** - New helper file with `WithTestAuth()` function
2. **auth/jwt.go** - Exported context keys for testing
3. **All test files** - Updated to use `WithTestAuth()` instead of manual context setting
4. **URL fixes** - Removed spaces from query parameters

## Next Steps

1. Fix the PostgreSQL port configuration
2. Run `make test-setup`
3. Run `make test`
4. All tests should now pass! ✅

## If Tests Still Fail

Check these:

```bash
# 1. Is PostgreSQL running?
docker ps | grep postgres

# 2. Can you connect?
psql -h localhost -p 543 -U postgres -l

# 3. Does test database exist?
psql -h localhost -p 543 -U postgres -c "\l" | grep localconnect_test

# 4. Re-create test database
docker exec -i localconnect-postgres psql -U postgres -c "DROP DATABASE IF EXISTS localconnect_test;"
make test-setup
```
