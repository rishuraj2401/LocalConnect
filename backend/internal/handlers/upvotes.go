package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"localconnect/internal/auth"
)

func (h *Handler) UpvoteProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	profileID := chi.URLParam(r, "id")
	_, err = h.DB.Exec(context.Background(),
		`INSERT INTO upvotes (profile_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		profileID, userID,
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not upvote")
		return
	}
	
	// Update count asynchronously
	h.Worker.Submit(func(ctx context.Context) {
		_, _ = h.DB.Exec(ctx,
			`UPDATE worker_profiles
			 SET upvote_count = (SELECT COUNT(*) FROM upvotes WHERE profile_id = $1)
			 WHERE id = $1`, profileID,
		)
	})
	
	// Invalidate cache
	h.Cache.Delete("profile:" + profileID)
	h.Cache.InvalidatePattern("profiles:")
	
	writeJSON(w, http.StatusOK, map[string]string{"status": "upvoted"})
}

func (h *Handler) RemoveUpvote(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	profileID := chi.URLParam(r, "id")
	_, err = h.DB.Exec(context.Background(),
		`DELETE FROM upvotes WHERE profile_id = $1 AND user_id = $2`,
		profileID, userID,
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not remove upvote")
		return
	}
	
	// Update count asynchronously
	h.Worker.Submit(func(ctx context.Context) {
		_, _ = h.DB.Exec(ctx,
			`UPDATE worker_profiles
			 SET upvote_count = (SELECT COUNT(*) FROM upvotes WHERE profile_id = $1)
			 WHERE id = $1`, profileID,
		)
	})
	
	// Invalidate cache
	h.Cache.Delete("profile:" + profileID)
	h.Cache.InvalidatePattern("profiles:")
	
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}
