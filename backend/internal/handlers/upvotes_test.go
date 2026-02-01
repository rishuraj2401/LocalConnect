package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestUpvoteProfile(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// Create test users and profile
	var workerID, clientID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Worker Upvote", "workerupvote@test.com", "1234567890", "hash", "worker",
	).Scan(&workerID)

	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Client Upvote", "clientupvote@test.com", "0987654321", "hash", "client",
	).Scan(&clientID)

	var profileID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO worker_profiles (user_id, category_id, location, rate, experience_years, bio)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		workerID, 1, "Test City", 100, 5, "Test bio",
	).Scan(&profileID)

	token := createTestToken(h, clientID, "client")

	// Test upvote
	req := httptest.NewRequest(http.MethodPost, "/profiles/"+profileID+"/upvote", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	req = WithTestAuth(req, clientID, "client")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", profileID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.UpvoteProfile(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	// Test duplicate upvote (should not fail, just be idempotent)
	req2 := httptest.NewRequest(http.MethodPost, "/profiles/"+profileID+"/upvote", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	req2 = WithTestAuth(req2, clientID, "client")
	req2 = req2.WithContext(context.WithValue(req2.Context(), chi.RouteCtxKey, rctx))
	rec2 := httptest.NewRecorder()
	h.UpvoteProfile(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("Expected status %d for duplicate upvote, got %d", http.StatusOK, rec2.Code)
	}
}

func TestRemoveUpvote(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// Create test users and profile
	var workerID, clientID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Worker Remove", "workerremove@test.com", "1234567890", "hash", "worker",
	).Scan(&workerID)

	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Client Remove", "clientremove@test.com", "0987654321", "hash", "client",
	).Scan(&clientID)

	var profileID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO worker_profiles (user_id, category_id, location, rate, experience_years, bio)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		workerID, 1, "Test City", 100, 5, "Test bio",
	).Scan(&profileID)

	// First upvote
	h.DB.Exec(context.Background(),
		`INSERT INTO upvotes (profile_id, user_id) VALUES ($1, $2)`,
		profileID, clientID,
	)

	token := createTestToken(h, clientID, "client")

	// Test remove upvote
	req := httptest.NewRequest(http.MethodDelete, "/profiles/"+profileID+"/upvote", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	req = WithTestAuth(req, clientID, "client")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", profileID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.RemoveUpvote(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	// Verify upvote was removed
	var count int
	h.DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM upvotes WHERE profile_id = $1 AND user_id = $2`,
		profileID, clientID,
	).Scan(&count)

	if count != 0 {
		t.Errorf("Expected 0 upvotes, got %d", count)
	}
}
