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

func TestCreateReview(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// Create test users and profile
	var workerID, clientID string
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Worker", "reviewworker@test.com", "1234567890", "hash", "worker",
	).Scan(&workerID)
	if err != nil {
		t.Fatalf("Failed to create worker: %v", err)
	}

	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Client", "reviewclient@test.com", "0987654321", "hash", "client",
	).Scan(&clientID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var profileID string
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO worker_profiles (user_id, category_id, location, rate, experience_years, bio)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		workerID, 1, "Test City", 100, 5, "Test bio",
	).Scan(&profileID)
	if err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	token := createTestToken(h, clientID, "client")

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Valid review with rating 5",
			payload: map[string]interface{}{
				"rating":  5,
				"comment": "Excellent work!",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp map[string]string
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp["id"] == "" {
					t.Error("Expected review ID in response")
				}
			},
		},
		{
			name: "Valid review with rating 1",
			payload: map[string]interface{}{
				"rating":  1,
				"comment": "Poor service",
			},
			expectedStatus: http.StatusCreated,
			checkResponse:  nil,
		},
		{
			name: "Invalid rating - too high",
			payload: map[string]interface{}{
				"rating":  6,
				"comment": "Good",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "Invalid rating - too low",
			payload: map[string]interface{}{
				"rating":  0,
				"comment": "Bad",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "Missing rating",
			payload: map[string]interface{}{
				"comment": "Good work",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/profiles/"+profileID+"/reviews", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			// Add context
			req = WithTestAuth(req, clientID, "client")
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", profileID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()
			h.CreateReview(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

func TestListReviews(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// Create test data
	var workerID, clientID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Worker Review", "workerreview@test.com", "1234567890", "hash", "worker",
	).Scan(&workerID)

	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Client Review", "clientreview@test.com", "0987654321", "hash", "client",
	).Scan(&clientID)

	var profileID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO worker_profiles (user_id, category_id, location, rate, experience_years, bio)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		workerID, 1, "Test City", 100, 5, "Test bio",
	).Scan(&profileID)

	// Add a review
	h.DB.Exec(context.Background(),
		`INSERT INTO reviews (profile_id, user_id, rating, comment)
		 VALUES ($1, $2, $3, $4)`,
		profileID, clientID, 5, "Great work!",
	)

	req := httptest.NewRequest(http.MethodGet, "/profiles/"+profileID+"/reviews", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", profileID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.ListReviews(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var reviews []reviewResponse
	if err := json.NewDecoder(rec.Body).Decode(&reviews); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(reviews) == 0 {
		t.Error("Expected at least one review")
	}

	if reviews[0].Rating != 5 {
		t.Errorf("Expected rating 5, got %d", reviews[0].Rating)
	}
	if reviews[0].Comment != "Great work!" {
		t.Errorf("Expected comment 'Great work!', got '%s'", reviews[0].Comment)
	}
}
