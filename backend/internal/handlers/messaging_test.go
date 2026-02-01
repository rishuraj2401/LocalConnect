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

func TestSendMessage(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// Create test users
	var user1ID, user2ID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"User 1", "user1@test.com", "1234567890", "hash", "client",
	).Scan(&user1ID)

	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"User 2", "user2@test.com", "0987654321", "hash", "worker",
	).Scan(&user2ID)

	token := createTestToken(h, user1ID, "client")

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid message",
			payload: map[string]interface{}{
				"receiver_id": user2ID,
				"content":     "Hello, I need your services",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing receiver_id",
			payload: map[string]interface{}{
				"content": "Hello",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing content",
			payload: map[string]interface{}{
				"receiver_id": user2ID,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Message to self",
			payload: map[string]interface{}{
				"receiver_id": user1ID,
				"content":     "Hello myself",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			req = WithTestAuth(req, user1ID, "client")

			rec := httptest.NewRecorder()
			h.SendMessage(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestGetConversations(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// Create test users
	var user1ID, user2ID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Conv User 1", "conv1@test.com", "1234567890", "hash", "client",
	).Scan(&user1ID)

	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Conv User 2", "conv2@test.com", "0987654321", "hash", "worker",
	).Scan(&user2ID)

	// Create conversation
	var convID string
	users := []string{user1ID, user2ID}
	if users[0] > users[1] {
		users[0], users[1] = users[1], users[0]
	}
	h.DB.QueryRow(context.Background(),
		`INSERT INTO conversations (user1_id, user2_id)
		 VALUES ($1, $2) RETURNING id`,
		users[0], users[1],
	).Scan(&convID)

	// Add a message
	var msgID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO messages (conversation_id, sender_id, receiver_id, content)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		convID, user1ID, user2ID, "Test message",
	).Scan(&msgID)

	token := createTestToken(h, user1ID, "client")

	req := httptest.NewRequest(http.MethodGet, "/conversations", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	req = WithTestAuth(req, user1ID, "client")

	rec := httptest.NewRecorder()
	h.GetConversations(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var conversations []conversationResponse
	if err := json.NewDecoder(rec.Body).Decode(&conversations); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(conversations) == 0 {
		t.Error("Expected at least one conversation")
	}

	if conversations[0].OtherUserID != user2ID {
		t.Errorf("Expected other user ID %s, got %s", user2ID, conversations[0].OtherUserID)
	}
}

func TestGetMessages(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// Create test users
	var user1ID, user2ID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Msg User 1", "msg1@test.com", "1234567890", "hash", "client",
	).Scan(&user1ID)

	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Msg User 2", "msg2@test.com", "0987654321", "hash", "worker",
	).Scan(&user2ID)

	// Create conversation
	var convID string
	users := []string{user1ID, user2ID}
	if users[0] > users[1] {
		users[0], users[1] = users[1], users[0]
	}
	h.DB.QueryRow(context.Background(),
		`INSERT INTO conversations (user1_id, user2_id)
		 VALUES ($1, $2) RETURNING id`,
		users[0], users[1],
	).Scan(&convID)

	// Add messages
	h.DB.Exec(context.Background(),
		`INSERT INTO messages (conversation_id, sender_id, receiver_id, content)
		 VALUES ($1, $2, $3, $4)`,
		convID, user1ID, user2ID, "Hello",
	)
	h.DB.Exec(context.Background(),
		`INSERT INTO messages (conversation_id, sender_id, receiver_id, content)
		 VALUES ($1, $2, $3, $4)`,
		convID, user2ID, user1ID, "Hi there",
	)

	token := createTestToken(h, user1ID, "client")

	req := httptest.NewRequest(http.MethodGet, "/conversations/"+convID+"/messages", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	req = WithTestAuth(req, user1ID, "client")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", convID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.GetMessages(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var messages []messageResponse
	if err := json.NewDecoder(rec.Body).Decode(&messages); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}

	if messages[0].Content != "Hello" {
		t.Errorf("Expected first message 'Hello', got '%s'", messages[0].Content)
	}
	if messages[1].Content != "Hi there" {
		t.Errorf("Expected second message 'Hi there', got '%s'", messages[1].Content)
	}
}

func TestGetMessagesUnauthorized(t *testing.T) {
	h := setupTestHandler(t)
	if h == nil {
		t.Skip("Skipping test - database not available")
		return
	}
	defer h.DB.Close()
	defer h.Worker.Shutdown()

	// Create test users
	var user1ID, user2ID, user3ID string
	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Unauth User 1", "unauth1@test.com", "1234567890", "hash", "client",
	).Scan(&user1ID)

	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Unauth User 2", "unauth2@test.com", "0987654321", "hash", "worker",
	).Scan(&user2ID)

	h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"Unauth User 3", "unauth3@test.com", "1111111111", "hash", "client",
	).Scan(&user3ID)

	// Create conversation between user1 and user2
	var convID string
	users := []string{user1ID, user2ID}
	if users[0] > users[1] {
		users[0], users[1] = users[1], users[0]
	}
	h.DB.QueryRow(context.Background(),
		`INSERT INTO conversations (user1_id, user2_id)
		 VALUES ($1, $2) RETURNING id`,
		users[0], users[1],
	).Scan(&convID)

	// Try to access as user3 (should fail)
	token := createTestToken(h, user3ID, "client")

	req := httptest.NewRequest(http.MethodGet, "/conversations/"+convID+"/messages", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	req = WithTestAuth(req, user3ID, "client")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", convID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.GetMessages(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusForbidden, rec.Code, rec.Body.String())
	}
}
