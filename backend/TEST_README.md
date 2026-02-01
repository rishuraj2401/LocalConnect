# LocalConnect Backend Testing Guide

This guide explains how to run and work with tests for the LocalConnect backend.

## Quick Start

### 1. Setup Test Database

```bash
make test-setup
```

This creates a separate `localconnect_test` database for testing.

### 2. Run All Tests

```bash
make test
```

This will:
- Check if PostgreSQL is running (start if needed)
- Setup the test database
- Run all tests with coverage
- Generate an HTML coverage report

### 3. View Coverage Report

Open `coverage.html` in your browser to see detailed code coverage.

## Test Organization

Tests are organized by handler:

```
backend/internal/handlers/
├── auth_test.go          # Authentication tests
├── profiles_test.go      # Profile CRUD tests
├── reviews_test.go       # Review tests
├── upvotes_test.go       # Upvote tests
├── contact_test.go       # Contact request tests
└── messaging_test.go     # Messaging tests
```

## Running Specific Tests

### Run a specific test file

```bash
go test -v ./internal/handlers/auth_test.go ./internal/handlers/handler.go ./internal/handlers/auth.go ./internal/handlers/response.go
```

### Run a specific test function

```bash
go test -v -run TestRegister ./internal/handlers/
```

### Run tests matching a pattern

```bash
go test -v -run TestCreate ./internal/handlers/
```

## Test Coverage Commands

### Quick coverage check

```bash
make test-quick
```

### Full coverage report

```bash
make test-coverage
```

### View coverage in terminal

```bash
go test -cover ./internal/handlers/...
```

### View coverage by function

```bash
go tool cover -func=coverage.out
```

## Test Database Management

### Setup test database

```bash
make test-setup
```

### Access test database

```bash
make db-shell-test
```

### Clean test database

```bash
make test-clean
```

### Reset everything (including main database)

```bash
make docker-reset
```

## Writing New Tests

### Test Structure

```go
func TestFeatureName(t *testing.T) {
    // 1. Setup handler
    h := setupTestHandler(t)
    if h == nil {
        t.Skip("Skipping test - database not available")
        return
    }
    defer h.DB.Close()
    defer h.Worker.Shutdown()

    // 2. Create test data
    // Insert test users, profiles, etc.

    // 3. Define test cases
    tests := []struct {
        name           string
        payload        map[string]interface{}
        expectedStatus int
        checkResponse  func(*testing.T, *httptest.ResponseRecorder)
    }{
        {
            name: "Valid case",
            payload: map[string]interface{}{
                "field": "value",
            },
            expectedStatus: http.StatusOK,
            checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
                // Verify response
            },
        },
        // More test cases...
    }

    // 4. Run test cases
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Make request
            // Check response
        })
    }
}
```

### Creating Authenticated Requests

```go
token := createTestToken(h, userID, "worker")

req := httptest.NewRequest(http.MethodPost, "/profiles", body)
req.Header.Set("Authorization", "Bearer "+token)

// Add context for auth middleware
ctx := context.WithValue(req.Context(), "userID", userID)
ctx = context.WithValue(ctx, "userRole", "worker")
req = req.WithContext(ctx)
```

### Testing with URL Parameters

```go
// Using chi router context
rctx := chi.NewRouteContext()
rctx.URLParams.Add("id", profileID)
req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
```

## Test Coverage Goals

- **Minimum**: 70% coverage across all handlers
- **Target**: 80%+ coverage
- **Critical paths**: 90%+ coverage (auth, payments, etc.)

## Current Test Coverage

Run `make test-coverage` to see current coverage statistics.

### Covered Endpoints

✅ **Authentication**
- POST /auth/register
- POST /auth/login

✅ **Categories**
- GET /categories

✅ **Profiles**
- GET /profiles (with filters)
- GET /profiles/:id
- POST /profiles
- PUT /profiles/:id

✅ **Reviews**
- GET /profiles/:id/reviews
- POST /profiles/:id/reviews

✅ **Upvotes**
- POST /profiles/:id/upvote
- DELETE /profiles/:id/upvote

✅ **Contact Requests**
- POST /profiles/:id/contact-requests
- GET /contact-requests

✅ **Messaging**
- POST /messages
- GET /conversations
- GET /conversations/:id/messages

### Not Covered

- Media upload (requires file handling)
- Edge cases and error scenarios (some)
- Integration tests
- Load tests

## Continuous Integration

### Pre-commit checks

```bash
make all
```

This runs:
1. Code formatting
2. Go vet
3. Linter (if installed)
4. All tests
5. Build

### GitHub Actions / CI Pipeline

Add this to your CI configuration:

```yaml
- name: Run tests
  run: |
    cd backend
    make test-setup
    make test
```

## Troubleshooting

### "Test database not available"

**Solution**: Start PostgreSQL and setup test database

```bash
make docker-start
sleep 3
make test-setup
```

### "Port 5432 already in use"

**Solution**: Either stop other PostgreSQL instances or configure a different port

```bash
# Stop other instances
brew services stop postgresql

# Or check what's using the port
lsof -i :5432
```

### "Tests are flaky"

**Common causes**:
- Race conditions in async operations
- Shared test data
- Database state not cleaned between tests

**Solution**: Ensure test isolation

```go
// Clean up test data
defer func() {
    h.DB.Exec(context.Background(), "DELETE FROM table WHERE id = $1", id)
}()
```

### "Coverage report not generating"

**Solution**: Ensure you have write permissions

```bash
chmod +x run_tests.sh
make test
```

## Best Practices

1. **Test Isolation**: Each test should be independent
2. **Cleanup**: Always defer cleanup of test data
3. **Meaningful Names**: Use descriptive test names
4. **Table-Driven Tests**: Use test tables for multiple scenarios
5. **Error Cases**: Test both success and failure paths
6. **Mock External Services**: Don't depend on external APIs
7. **Fast Tests**: Keep tests fast (< 100ms per test)
8. **Deterministic**: Tests should always produce same results

## Performance Testing

### Benchmark tests

```go
func BenchmarkListProfiles(b *testing.B) {
    h := setupTestHandler(&testing.T{})
    // ... setup

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Run operation
    }
}
```

Run benchmarks:

```bash
go test -bench=. ./internal/handlers/
```

### Load testing

Use tools like:
- `ab` (Apache Bench)
- `hey`
- `wrk`
- `k6`

Example:

```bash
# Install hey
go install github.com/rakyll/hey@latest

# Run load test
hey -n 1000 -c 10 http://localhost:8080/profiles
```

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Go Testing Best Practices](https://golang.org/doc/effective_go#testing)
- [httptest Package](https://golang.org/pkg/net/http/httptest/)
