package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"localconnect/internal/analytics"
	"localconnect/internal/auth"
	"localconnect/internal/notifications"
)

type contactRequestPayload struct {
	Message     string `json:"message"`
	PhoneShared bool   `json:"phone_shared"`
}

type contactResponse struct {
	ID        string `json:"id"`
	ProfileID string `json:"profile_id"`
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) CreateContactRequest(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	profileID := chi.URLParam(r, "id")
	var payload contactRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	var phone string
	err = h.DB.QueryRow(context.Background(),
		`SELECT phone FROM users WHERE id = $1`, userID,
	).Scan(&phone)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not fetch phone")
		return
	}
	if !payload.PhoneShared {
		phone = ""
	}
	var contactID string
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO contact_requests (profile_id, user_id, message, phone_shared, phone)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		profileID, userID, payload.Message, payload.PhoneShared, phone,
	).Scan(&contactID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not create request")
		return
	}
	
	// Track analytics
	h.Tracker.Track(analytics.Event{
		Type:      "contact_request",
		UserID:    userID,
		ProfileID: profileID,
	})
	
	// Send notification to profile owner asynchronously
	h.Worker.Submit(func(ctx context.Context) {
		var ownerID string
		if err := h.DB.QueryRow(ctx, `SELECT user_id FROM worker_profiles WHERE id = $1`, profileID).Scan(&ownerID); err == nil {
			h.Notifier.Send(notifications.Notification{
				UserID:  ownerID,
				Type:    "contact_request",
				Message: "Someone wants to contact you",
				Data: map[string]interface{}{
					"profile_id": profileID,
					"user_id":    userID,
				},
			})
		}
	})
	
	writeJSON(w, http.StatusCreated, map[string]string{"id": contactID})
}

func (h *Handler) ListContactRequests(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	rows, err := h.DB.Query(context.Background(),
		`SELECT cr.id, cr.profile_id, cr.user_id, cr.message, cr.phone, cr.created_at
		 FROM contact_requests cr
		 JOIN worker_profiles wp ON wp.id = cr.profile_id
		 WHERE wp.user_id = $1
		 ORDER BY cr.created_at DESC`, userID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch requests")
		return
	}
	defer rows.Close()

	var requests []contactResponse
	for rows.Next() {
		var item contactResponse
		var createdAt time.Time
		if err := rows.Scan(&item.ID, &item.ProfileID, &item.UserID, &item.Message, &item.Phone, &createdAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to parse requests")
			return
		}
		item.CreatedAt = createdAt.Format(time.RFC3339)
		requests = append(requests, item)
	}
	writeJSON(w, http.StatusOK, requests)
}
