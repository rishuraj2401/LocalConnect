package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"localconnect/internal/auth"

	"github.com/go-chi/chi/v5"
)

func createTestToken(h *Handler, userID, role string) string {
	token, _ := auth.GenerateToken(h.Config.JWTSecret, userID, role)
	return token
}

func TestListCategories(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	rec := httptest.NewRecorder()

	h.ListCategories(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var categories []categoryResponse
	if err := json.NewDecoder(rec.Body).Decode(&categories); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(categories) == 0 {
		t.Error("Expected at least one category")
	}

	// Check for expected categories
	expectedCategories := []string{"labour", "cook", "painter", "carpenter", "home tution", "teacher", "househelp"}
	categoryNames := make(map[string]bool)
	for _, cat := range categories {
		categoryNames[cat.Name] = true
	}

	for _, expected := range expectedCategories {
		if !categoryNames[expected] {
			t.Errorf("Expected category '%s' not found", expected)
		}
	}
}

func TestCreateProfile(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// Create a test worker user
	var userID string
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Test Worker", "worker@test.com", "1234567890", "hash", "worker",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	token := createTestToken(h, userID, "worker")

	tests := []struct {
		name           string
		payload        map[string]interface{}
		token          string
		expectedStatus int
	}{
		{
			name: "Valid profile creation",
			payload: map[string]interface{}{
				"category_id":      1,
				"location":         "New York",
				"rate":             100.50,
				"experience_years": 5,
				"bio":              "Experienced worker",
			},
			token:          token,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing category_id",
			payload: map[string]interface{}{
				"location":         "NewYork",
				"rate":             100.50,
				"experience_years": 5,
				"bio":              "Experienced worker",
			},
			token:          token,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing location",
			payload: map[string]interface{}{
				"category_id":      1,
				"rate":             100.50,
				"experience_years": 5,
				"bio":              "Experienced worker",
			},
			token:          token,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No authentication token",
			payload: map[string]interface{}{
				"category_id":      1,
				"location":         "New York",
				"rate":             100.50,
				"experience_years": 5,
				"bio":              "Experienced worker",
			},
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/profiles", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Add auth context if token is present
			if tt.token != "" {
				req = WithTestAuth(req, userID, "worker")
			}

			rec := httptest.NewRecorder()
			h.CreateProfile(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestListProfiles(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "List all profiles",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var profiles []profileResponse
				if err := json.NewDecoder(rec.Body).Decode(&profiles); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
			},
		},
		{
			name:           "Filter by category",
			queryParams:    "?category=cook",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var profiles []profileResponse
				if err := json.NewDecoder(rec.Body).Decode(&profiles); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				for _, profile := range profiles {
					if profile.CategoryName != "cook" {
						t.Errorf("Expected category 'cook', got '%s'", profile.CategoryName)
					}
				}
			},
		},
		{
			name:           "Filter by location",
			queryParams:    "?location=NewYork",
			expectedStatus: http.StatusOK,
			checkResponse:  nil,
		},
		{
			name:           "Filter by both category and location",
			queryParams:    "?category=painter&location=Boston",
			expectedStatus: http.StatusOK,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/profiles"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			h.ListProfiles(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// Create a test profile
	var userID string
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Test Worker 2", "worker2@test.com", "1234567890", "hash", "worker",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var profileID string
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO worker_profiles (user_id, category_id, location, rate, experience_years, bio)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		userID, 1, "Test City", 100, 5, "Test bio",
	).Scan(&profileID)
	if err != nil {
		t.Fatalf("Failed to create test profile: %v", err)
	}

	tests := []struct {
		name           string
		profileID      string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "Valid profile ID",
			profileID:      profileID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var profile profileResponse
				if err := json.NewDecoder(rec.Body).Decode(&profile); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if profile.ID != profileID {
					t.Errorf("Expected profile ID %s, got %s", profileID, profile.ID)
				}
				if profile.Location != "Test City" {
					t.Errorf("Expected location 'Test City', got '%s'", profile.Location)
				}
			},
		},
		{
			name:           "Invalid profile ID",
			profileID:      "00000000-0000-0000-0000-000000000000",
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/profiles/"+tt.profileID, nil)
			rec := httptest.NewRecorder()

			// Add chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.profileID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			h.GetProfile(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// Create a test profile
	var userID string
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Test Worker 3", "worker3@test.com", "1234567890", "hash", "worker",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var profileID string
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO worker_profiles (user_id, category_id, location, rate, experience_years, bio)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		userID, 1, "Old City", 50, 3, "Old bio",
	).Scan(&profileID)
	if err != nil {
		t.Fatalf("Failed to create test profile: %v", err)
	}

	token := createTestToken(h, userID, "worker")

	updatePayload := map[string]interface{}{
		"category_id":      2,
		"location":         "New City",
		"rate":             150.0,
		"experience_years": 7,
		"bio":              "Updated bio",
	}

	body, _ := json.Marshal(updatePayload)
	req := httptest.NewRequest(http.MethodPut, "/profiles/"+profileID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Add context
	req = WithTestAuth(req, userID, "worker")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", profileID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.UpdateProfile(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	// Verify update
	var updatedProfile profileResponse
	req2 := httptest.NewRequest(http.MethodGet, "/profiles/"+profileID, nil)
	rctx2 := chi.NewRouteContext()
	rctx2.URLParams.Add("id", profileID)
	req2 = req2.WithContext(context.WithValue(req2.Context(), chi.RouteCtxKey, rctx2))
	rec2 := httptest.NewRecorder()
	h.GetProfile(rec2, req2)

	if err := json.NewDecoder(rec2.Body).Decode(&updatedProfile); err != nil {
		t.Fatalf("Failed to decode updated profile: %v", err)
	}

	if updatedProfile.Location != "New City" {
		t.Errorf("Expected location 'New City', got '%s'", updatedProfile.Location)
	}
	if updatedProfile.Rate != 150.0 {
		t.Errorf("Expected rate 150.0, got %f", updatedProfile.Rate)
	}
}
