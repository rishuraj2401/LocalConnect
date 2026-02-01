package handlers

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"localconnect/internal/auth"
	"localconnect/internal/media"
)

func (h *Handler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	profileID := chi.URLParam(r, "id")
	var ownerID string
	err = h.DB.QueryRow(context.Background(),
		`SELECT user_id FROM worker_profiles WHERE id = $1`, profileID,
	).Scan(&ownerID)
	if err != nil || ownerID != userID {
		writeError(w, http.StatusForbidden, "not allowed")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid upload")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing file")
		return
	}
	defer file.Close()

	saved, err := media.SaveFile(r.Context(), h.Config.MediaDir, file, header)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	publicPath := filepath.ToSlash(saved.Path)
	if strings.HasPrefix(publicPath, ".") {
		publicPath = strings.TrimPrefix(publicPath, ".")
	}

	var mediaID string
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO profile_media (profile_id, media_type, url)
		 VALUES ($1, $2, $3) RETURNING id`,
		profileID, saved.MediaType, publicPath,
	).Scan(&mediaID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not save media")
		return
	}

	h.Worker.Submit(func(ctx context.Context) {
		_ = ctx
	})

	writeJSON(w, http.StatusCreated, map[string]string{"id": mediaID, "url": publicPath})
}

type mediaResponse struct {
	ID        string `json:"id"`
	ProfileID string `json:"profile_id"`
	MediaType string `json:"media_type"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) ListMedia(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "id")
	if profileID == "" {
		writeError(w, http.StatusBadRequest, "missing profile id")
		return
	}

	rows, err := h.DB.Query(context.Background(),
		`SELECT id, profile_id, media_type, url, created_at 
		 FROM profile_media 
		 WHERE profile_id = $1 
		 ORDER BY created_at DESC`,
		profileID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch media")
		return
	}
	defer rows.Close()

	var mediaList []mediaResponse
	for rows.Next() {
		var m mediaResponse
		var createdAt interface{}
		if err := rows.Scan(&m.ID, &m.ProfileID, &m.MediaType, &m.URL, &createdAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to parse media")
			return
		}
		m.CreatedAt = createdAt.(time.Time).Format(time.RFC3339)
		mediaList = append(mediaList, m)
	}

	writeJSON(w, http.StatusOK, mediaList)
}
