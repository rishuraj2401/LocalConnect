package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"localconnect/internal/config"
	"localconnect/internal/worker"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	// Use test database URL - using port 543 as shown in docker ps
	dbURL := "postgres://postgres:postgres@localhost:543/localconnect_test?sslmode=disable"
	
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Skip("Test database not available:", err)
		return nil
	}
	
	// Test connection
	if err := pool.Ping(ctx); err != nil {
		t.Skip("Cannot connect to test database:", err)
		return nil
	}
	
	return pool
}

func setupTestHandler(t *testing.T) *Handler {
	pool := setupTestDB(t)
	if pool == nil {
		return nil
	}

	cfg := config.Config{
		Port:        "8080",
		DatabaseURL: "postgres://postgres:postgres@localhost:543/localconnect_test?sslmode=disable",
		JWTSecret:   "test-secret",
		MediaDir:    "./test-media",
	}

	workerPool := worker.NewPool(10, 2)

	return New(pool, cfg, workerPool)
}

func TestRegister(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Valid worker registration",
			payload: map[string]interface{}{
				"name":     "John Doe",
				"email":    "john@example.com",
				"phone":    "1234567890",
				"password": "password123",
				"role":     "worker",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp authResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Token == "" {
					t.Error("Expected token in response")
				}
				if resp.Role != "worker" {
					t.Errorf("Expected role 'worker', got '%s'", resp.Role)
				}
			},
		},
		{
			name: "Valid client registration",
			payload: map[string]interface{}{
				"name":     "Jane Client",
				"email":    "jane@example.com",
				"phone":    "0987654321",
				"password": "password123",
				"role":     "client",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp authResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Role != "client" {
					t.Errorf("Expected role 'client', got '%s'", resp.Role)
				}
			},
		},
		{
			name: "Missing email",
			payload: map[string]interface{}{
				"name":     "Test User",
				"phone":    "1234567890",
				"password": "password123",
				"role":     "worker",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp map[string]string
				json.NewDecoder(rec.Body).Decode(&resp)
				if _, ok := resp["error"]; !ok {
					t.Error("Expected error in response")
				}
			},
		},
		{
			name: "Invalid role",
			payload: map[string]interface{}{
				"name":     "Test User",
				"email":    "test@example.com",
				"phone":    "1234567890",
				"password": "password123",
				"role":     "admin",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "Missing password",
			payload: map[string]interface{}{
				"name":  "Test User",
				"email": "test2@example.com",
				"phone": "1234567890",
				"role":  "worker",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.Register(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// First register a user
	registerPayload := map[string]interface{}{
		"name":     "Login Test User",
		"email":    "logintest@example.com",
		"phone":    "1234567890",
		"password": "testpassword",
		"role":     "worker",
	}
	body, _ := json.Marshal(registerPayload)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Failed to register test user: %s", rec.Body.String())
	}

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Valid login",
			payload: map[string]interface{}{
				"email":    "logintest@example.com",
				"password": "testpassword",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp authResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Token == "" {
					t.Error("Expected token in response")
				}
				if resp.Role != "worker" {
					t.Errorf("Expected role 'worker', got '%s'", resp.Role)
				}
			},
		},
		{
			name: "Invalid password",
			payload: map[string]interface{}{
				"email":    "logintest@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  nil,
		},
		{
			name: "Invalid email",
			payload: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "testpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  nil,
		},
		{
			name: "Missing email",
			payload: map[string]interface{}{
				"password": "testpassword",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "Missing password",
			payload: map[string]interface{}{
				"email": "logintest@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.Login(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}
