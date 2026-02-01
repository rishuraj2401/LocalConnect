package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"localconnect/internal/auth"
	"localconnect/internal/notifications"
)

type messageRequest struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

type messageResponse struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	Read       bool   `json:"read"`
	CreatedAt  string `json:"created_at"`
}

type conversationResponse struct {
	ID              string `json:"id"`
	OtherUserID     string `json:"other_user_id"`
	OtherUserName   string `json:"other_user_name"`
	LastMessage     string `json:"last_message"`
	LastMessageAt   string `json:"last_message_at"`
	UnreadCount     int    `json:"unread_count"`
	CreatedAt       string `json:"created_at"`
}

// SendMessage sends a message to another user
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req messageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	if req.ReceiverID == "" || req.Content == "" {
		writeError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	if req.ReceiverID == userID {
		writeError(w, http.StatusBadRequest, "cannot message yourself")
		return
	}

	// Process message sending asynchronously
	h.Worker.Submit(func(ctx context.Context) {
		if err := h.createMessage(ctx, userID, req.ReceiverID, req.Content); err != nil {
			return
		}
		
		// Send notification
		h.Notifier.Send(notifications.Notification{
			UserID:  req.ReceiverID,
			Type:    "new_message",
			Message: "You have a new message",
			Data: map[string]interface{}{
				"sender_id": userID,
			},
		})
	})

	writeJSON(w, http.StatusCreated, map[string]string{"status": "message sent"})
}

func (h *Handler) createMessage(ctx context.Context, senderID, receiverID, content string) error {
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Get or create conversation
	var conversationID string
	user1, user2 := senderID, receiverID
	if user1 > user2 {
		user1, user2 = user2, user1
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO conversations (user1_id, user2_id)
		 VALUES ($1, $2)
		 ON CONFLICT (user1_id, user2_id) DO UPDATE SET user1_id = conversations.user1_id
		 RETURNING id`,
		user1, user2,
	).Scan(&conversationID)
	if err != nil {
		return err
	}

	// Insert message
	_, err = tx.Exec(ctx,
		`INSERT INTO messages (conversation_id, sender_id, receiver_id, content)
		 VALUES ($1, $2, $3, $4)`,
		conversationID, senderID, receiverID, content,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetConversations returns all conversations for the current user
func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	rows, err := h.DB.Query(context.Background(), `
		SELECT 
			c.id,
			CASE WHEN c.user1_id = $1 THEN c.user2_id ELSE c.user1_id END as other_user_id,
			u.name as other_user_name,
			COALESCE(m.content, '') as last_message,
			c.last_message_at,
			(SELECT COUNT(*) FROM messages 
			 WHERE conversation_id = c.id 
			 AND receiver_id = $1 
			 AND read = false) as unread_count,
			c.created_at
		FROM conversations c
		LEFT JOIN users u ON u.id = (CASE WHEN c.user1_id = $1 THEN c.user2_id ELSE c.user1_id END)
		LEFT JOIN messages m ON m.id = c.last_message_id
		WHERE c.user1_id = $1 OR c.user2_id = $1
		ORDER BY c.last_message_at DESC
	`, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch conversations")
		return
	}
	defer rows.Close()

	var conversations []conversationResponse
	for rows.Next() {
		var conv conversationResponse
		var lastMessageAt, createdAt time.Time
		if err := rows.Scan(
			&conv.ID,
			&conv.OtherUserID,
			&conv.OtherUserName,
			&conv.LastMessage,
			&lastMessageAt,
			&conv.UnreadCount,
			&createdAt,
		); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to parse conversations")
			return
		}
		conv.LastMessageAt = lastMessageAt.Format(time.RFC3339)
		conv.CreatedAt = createdAt.Format(time.RFC3339)
		conversations = append(conversations, conv)
	}

	writeJSON(w, http.StatusOK, conversations)
}

// GetMessages returns all messages in a conversation
func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	conversationID := chi.URLParam(r, "id")
	if conversationID == "" {
		writeError(w, http.StatusBadRequest, "missing conversation id")
		return
	}

	// Verify user is part of conversation
	var user1ID, user2ID string
	err = h.DB.QueryRow(context.Background(),
		`SELECT user1_id, user2_id FROM conversations WHERE id = $1`,
		conversationID,
	).Scan(&user1ID, &user2ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "conversation not found")
		} else {
			writeError(w, http.StatusInternalServerError, "failed to fetch conversation")
		}
		return
	}

	if user1ID != userID && user2ID != userID {
		writeError(w, http.StatusForbidden, "not authorized for this conversation")
		return
	}

	// Mark messages as read asynchronously
	h.Worker.Submit(func(ctx context.Context) {
		_, _ = h.DB.Exec(ctx,
			`UPDATE messages SET read = true 
			 WHERE conversation_id = $1 AND receiver_id = $2 AND read = false`,
			conversationID, userID,
		)
	})

	// Fetch messages
	rows, err := h.DB.Query(context.Background(), `
		SELECT id, sender_id, receiver_id, content, read, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`, conversationID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch messages")
		return
	}
	defer rows.Close()

	var messages []messageResponse
	for rows.Next() {
		var msg messageResponse
		var createdAt time.Time
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Read, &createdAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to parse messages")
			return
		}
		msg.CreatedAt = createdAt.Format(time.RFC3339)
		messages = append(messages, msg)
	}

	writeJSON(w, http.StatusOK, messages)
}
