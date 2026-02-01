package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"localconnect/internal/auth"
	"localconnect/internal/notifications"
)

type reviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type reviewResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	profileID := chi.URLParam(r, "id")
	var req reviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if req.Rating < 1 || req.Rating > 5 {
		writeError(w, http.StatusBadRequest, "rating must be 1-5")
		return
	}
	var reviewID string
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO reviews (profile_id, user_id, rating, comment)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		profileID, userID, req.Rating, req.Comment,
	).Scan(&reviewID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not add review")
		return
	}
	
	// Update profile stats asynchronously
	h.Worker.Submit(func(ctx context.Context) {
		_, _ = h.DB.Exec(ctx,
			`UPDATE worker_profiles
			 SET review_count = (SELECT COUNT(*) FROM reviews WHERE profile_id = $1),
			     average_rating = (SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE profile_id = $1)
			 WHERE id = $1`, profileID,
		)
		
		// Get profile owner to send notification
		var ownerID string
		if err := h.DB.QueryRow(ctx, `SELECT user_id FROM worker_profiles WHERE id = $1`, profileID).Scan(&ownerID); err == nil {
			h.Notifier.Send(notifications.Notification{
				UserID:  ownerID,
				Type:    "new_review",
				Message: "You received a new review",
				Data: map[string]interface{}{
					"profile_id": profileID,
					"rating":     req.Rating,
				},
			})
		}
	})
	
	// Invalidate cache
	h.Cache.Delete("profile:" + profileID)
	h.Cache.InvalidatePattern("profiles:")
	
	writeJSON(w, http.StatusCreated, map[string]string{"id": reviewID})
}

func (h *Handler) ListReviews(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "id")
	rows, err := h.DB.Query(context.Background(),
		`SELECT id, user_id, rating, comment, created_at
		 FROM reviews WHERE profile_id = $1 ORDER BY created_at DESC`, profileID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch reviews")
		return
	}
	defer rows.Close()

	var reviews []reviewResponse
	for rows.Next() {
		var rev reviewResponse
		var createdAt time.Time
		if err := rows.Scan(&rev.ID, &rev.UserID, &rev.Rating, &rev.Comment, &createdAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to parse reviews")
			return
		}
		rev.CreatedAt = createdAt.Format(time.RFC3339)
		reviews = append(reviews, rev)
	}
	writeJSON(w, http.StatusOK, reviews)
}
