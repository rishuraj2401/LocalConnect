package handlers

import (
	"context"
	"net/http"
	
	"localconnect/internal/auth"
)

// authContextKey is used to mark requests that should bypass auth middleware in tests
type authContextKey string

const testAuthKey authContextKey = "test_auth_bypass"

// WithTestAuth adds authentication context to a request for testing
func WithTestAuth(req *http.Request, userID, role string) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, testAuthKey, true)
	ctx = context.WithValue(ctx, auth.UserIDKey, userID)
	ctx = context.WithValue(ctx, auth.UserRoleKey, role)
	return req.WithContext(ctx)
}
