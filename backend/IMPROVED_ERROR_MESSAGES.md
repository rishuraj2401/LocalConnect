# Improved Error Messages

## Changes Made

Updated error messages across all handlers to provide better debugging information.

### Auth Handler

**Before:**
```json
{"error": "could not create user"}
```

**After:**
```json
{"error": "email already registered"}
{"error": "database not initialized - run migrations first"}
{"error": "could not create user: ERROR: relation 'users' does not exist"}
```

### Profile Handler

**Before:**
```json
{"error": "could not create profile"}
```

**After:**
```json
{"error": "profile already exists for this user"}
{"error": "invalid category or user"}
{"error": "could not create profile: [actual database error]"}
```

## Common Errors and Solutions

### 1. "database not initialized - run migrations first"

**Problem:** Tables don't exist in the database

**Solution:**
```bash
cd backend
make migrate-up
```

Or manually:
```bash
docker exec -i localconnect-postgres psql -U postgres -d localconnect < migrations/001_init.sql
docker exec -i localconnect-postgres psql -U postgres -d localconnect < migrations/002_messaging.sql
```

### 2. "email already registered"

**Problem:** Email is already in use

**Solution:** Use a different email or delete the existing user

### 3. "relation 'users' does not exist"

**Problem:** Database tables haven't been created

**Solution:** Run migrations (see solution #1)

### 4. "profile already exists for this user"

**Problem:** User already has a worker profile (each user can only have one)

**Solution:** Update the existing profile instead of creating a new one

## Testing Error Messages

### Test Registration

```bash
# Test with missing database
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@test.com","phone":"123","password":"pass","role":"worker"}'

# Test duplicate email
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test2","email":"test@test.com","phone":"456","password":"pass","role":"client"}'
```

### Test Profile Creation

```bash
# Get token first
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"pass"}' | jq -r .token)

# Test profile creation
curl -X POST http://localhost:8080/profiles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"category_id":1,"location":"NYC","rate":100,"experience_years":5,"bio":"Expert"}'
```

## Benefits

1. **Better Debugging** - See actual database errors
2. **User-Friendly** - Specific messages for common issues
3. **Faster Development** - Identify problems quickly
4. **Production Ready** - Can filter sensitive errors in production

## Future Improvements

### Production Mode

In production, you may want to hide detailed error messages:

```go
func (h *Handler) isProduction() bool {
    return h.Config.Environment == "production"
}

// Then in handlers:
if err != nil {
    if h.isProduction() {
        writeError(w, http.StatusBadRequest, "could not create user")
    } else {
        writeError(w, http.StatusBadRequest, fmt.Sprintf("could not create user: %v", err))
    }
    return
}
```

### Structured Logging

Add structured logging for all errors:

```go
import "log/slog"

if err != nil {
    slog.Error("failed to create user", 
        "error", err,
        "email", req.Email,
        "role", req.Role,
    )
    writeError(w, http.StatusBadRequest, "could not create user")
    return
}
```

## Summary

✅ **Registration errors** now show specific messages  
✅ **Database errors** are visible for debugging  
✅ **Duplicate checks** have user-friendly messages  
✅ **Foreign key errors** indicate invalid references

Restart your backend to see the improved error messages!

```bash
# Stop backend (Ctrl+C)
# Start again
cd backend
go run cmd/api/main.go
```

Now when you try to register, you'll see exactly what's wrong! 🎯
