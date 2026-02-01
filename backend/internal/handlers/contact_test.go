package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestCreateContactRequest(t *testing.T) {
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
		"Worker Contact", "workercontact@test.com", "1234567890", "hash", "worker",
	).Scan(&workerID)

	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Client Contact", "clientcontact@test.com", "0987654321", "hash", "client",
	).Scan(&clientID)

	var profileID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO worker_profiles (user_id, category_id, location, rate, experience_years, bio)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		workerID, 1, "Test City", 100, 5, "Test bio",
	).Scan(&profileID)

	token := createTestToken(h, clientID, "client")

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid contact request with phone",
			payload: map[string]interface{}{
				"message":      "I need your services",
				"phone_shared": true,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Valid contact request without phone",
			payload: map[string]interface{}{
				"message":      "Please contact me",
				"phone_shared": false,
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/profiles/"+profileID+"/contact-requests", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			req = WithTestAuth(req, clientID, "client")
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", profileID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()
			h.CreateContactRequest(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestListContactRequests(t *testing.T) {
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
		"Worker List", "workerlist@test.com", "1234567890", "hash", "worker",
	).Scan(&workerID)

	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Client List", "clientlist@test.com", "0987654321", "hash", "client",
	).Scan(&clientID)

	var profileID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO worker_profiles (user_id, category_id, location, rate, experience_years, bio)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		workerID, 1, "Test City", 100, 5, "Test bio",
	).Scan(&profileID)

	// Add contact request
	h.DB.Exec(context.Background(),
		`INSERT INTO contact_requests (profile_id, user_id, message, phone_shared, phone)
		 VALUES ($1, $2, $3, $4, $5)`,
		profileID, clientID, "Need services", true, "0987654321",
	)

	token := createTestToken(h, workerID, "worker")

	req := httptest.NewRequest(http.MethodGet, "/contact-requests", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	req = WithTestAuth(req, workerID, "worker")

	rec := httptest.NewRecorder()
	h.ListContactRequests(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var requests []contactResponse
	if err := json.NewDecoder(rec.Body).Decode(&requests); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(requests) == 0 {
		t.Error("Expected at least one contact request")
	}
}
