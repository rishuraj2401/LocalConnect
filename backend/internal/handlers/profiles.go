package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"localconnect/internal/analytics"
	"localconnect/internal/auth"
)

type profileResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	CategoryID      int     `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	Location        string  `json:"location"`
	Rate            float64 `json:"rate"`
	ExperienceYears int     `json:"experience_years"`
	Bio             string  `json:"bio"`
	Upvotes         int     `json:"upvotes"`
	ReviewCount     int     `json:"review_count"`
	AverageRating   float64 `json:"average_rating"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type profileRequest struct {
	CategoryID      int     `json:"category_id"`
	Location        string  `json:"location"`
	Rate            float64 `json:"rate"`
	ExperienceYears int     `json:"experience_years"`
	Bio             string  `json:"bio"`
}

func (h *Handler) ListProfiles(w http.ResponseWriter, r *http.Request) {
	category := strings.TrimSpace(r.URL.Query().Get("category"))
	location := strings.TrimSpace(r.URL.Query().Get("location"))

	// Track analytics
	h.Tracker.Track(analytics.Event{
		Type: "category_search",
		Metadata: map[string]interface{}{
			"category": category,
			"location": location,
		},
	})

	// Try cache first
	cacheKey := "profiles:" + category + ":" + location
	if cached, found := h.Cache.Get(cacheKey); found {
		if profiles, ok := cached.([]profileResponse); ok {
			writeJSON(w, http.StatusOK, profiles)
			return
		}
	}

	query := `
		SELECT p.id, p.user_id, p.category_id, c.name, p.location, p.rate, p.experience_years, p.bio,
			p.upvote_count, p.review_count, p.average_rating,
			p.created_at, p.updated_at
		FROM worker_profiles p
		JOIN categories c ON c.id = p.category_id
		WHERE ($1 = '' OR c.name ILIKE $1 || '%')
		AND ($2 = '' OR p.location ILIKE $2 || '%')
		ORDER BY p.upvote_count DESC, p.review_count DESC`

	rows, err := h.DB.Query(context.Background(), query, category, location)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch profiles")
		return
	}
	defer rows.Close()

	var profiles []profileResponse
	for rows.Next() {
		var p profileResponse
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(&p.ID, &p.UserID, &p.CategoryID, &p.CategoryName, &p.Location, &p.Rate, &p.ExperienceYears, &p.Bio, &p.Upvotes, &p.ReviewCount, &p.AverageRating, &createdAt, &updatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to parse profiles")
			return
		}
		p.CreatedAt = createdAt.Format(time.RFC3339)
		p.UpdatedAt = updatedAt.Format(time.RFC3339)
		profiles = append(profiles, p)
	}
	
	// Cache the result
	h.Cache.Set(cacheKey, profiles)
	
	writeJSON(w, http.StatusOK, profiles)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing profile id")
		return
	}
	
	// Track profile view
	h.Tracker.Track(analytics.Event{
		Type:      "profile_view",
		ProfileID: id,
	})
	
	// Try cache first
	cacheKey := "profile:" + id
	if cached, found := h.Cache.Get(cacheKey); found {
		if profile, ok := cached.(profileResponse); ok {
			writeJSON(w, http.StatusOK, profile)
			return
		}
	}
	
	query := `
		SELECT p.id, p.user_id, p.category_id, c.name, p.location, p.rate, p.experience_years, p.bio,
			p.upvote_count, p.review_count, p.average_rating,
			p.created_at, p.updated_at
		FROM worker_profiles p
		JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1`

	var p profileResponse
	var createdAt time.Time
	var updatedAt time.Time
	err := h.DB.QueryRow(context.Background(), query, id).Scan(&p.ID, &p.UserID, &p.CategoryID, &p.CategoryName, &p.Location, &p.Rate, &p.ExperienceYears, &p.Bio, &p.Upvotes, &p.ReviewCount, &p.AverageRating, &createdAt, &updatedAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "profile not found")
		return
	}
	p.CreatedAt = createdAt.Format(time.RFC3339)
	p.UpdatedAt = updatedAt.Format(time.RFC3339)
	
	// Cache the result
	h.Cache.Set(cacheKey, p)
	
	writeJSON(w, http.StatusOK, p)
}

// GetMyProfile returns the current logged-in worker's profile
func (h *Handler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	role, err := auth.RoleFromContext(r.Context())
	if err != nil || role != "worker" {
		writeError(w, http.StatusForbidden, "workers only")
		return
	}
	
	query := `
		SELECT p.id, p.user_id, p.category_id, c.name, p.location, p.rate, p.experience_years, p.bio,
			p.upvote_count, p.review_count, p.average_rating,
			p.created_at, p.updated_at
		FROM worker_profiles p
		JOIN categories c ON c.id = p.category_id
		WHERE p.user_id = $1`

	var p profileResponse
	var createdAt time.Time
	var updatedAt time.Time
	err = h.DB.QueryRow(context.Background(), query, userID).Scan(&p.ID, &p.UserID, &p.CategoryID, &p.CategoryName, &p.Location, &p.Rate, &p.ExperienceYears, &p.Bio, &p.Upvotes, &p.ReviewCount, &p.AverageRating, &createdAt, &updatedAt)
	if err != nil {
		// Profile doesn't exist yet
		writeJSON(w, http.StatusOK, map[string]interface{}{"exists": false})
		return
	}
	p.CreatedAt = createdAt.Format(time.RFC3339)
	p.UpdatedAt = updatedAt.Format(time.RFC3339)
	
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	role, err := auth.RoleFromContext(r.Context())
	if err != nil || role != "worker" {
		writeError(w, http.StatusForbidden, "workers only")
		return
	}

	var req profileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if req.CategoryID == 0 || req.Location == "" {
		writeError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	var profileID string
	err = h.DB.QueryRow(context.Background(), `
		INSERT INTO worker_profiles (user_id, category_id, location, rate, experience_years, bio)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`, userID, req.CategoryID, req.Location, req.Rate, req.ExperienceYears, req.Bio,
	).Scan(&profileID)
	if err != nil {
		errMsg := "could not create profile"
		if strings.Contains(err.Error(), "duplicate") {
			errMsg = "profile already exists for this user"
		} else if strings.Contains(err.Error(), "foreign key") {
			errMsg = "invalid category or user"
		} else {
			errMsg = fmt.Sprintf("could not create profile: %v", err)
		}
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}
	
	// Invalidate profile cache
	h.Cache.InvalidatePattern("profiles:")
	
	writeJSON(w, http.StatusCreated, map[string]string{"id": profileID})
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	role, err := auth.RoleFromContext(r.Context())
	if err != nil || role != "worker" {
		writeError(w, http.StatusForbidden, "workers only")
		return
	}
	profileID := chi.URLParam(r, "id")
	if profileID == "" {
		writeError(w, http.StatusBadRequest, "missing profile id")
		return
	}
	var req profileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	result, err := h.DB.Exec(context.Background(), `
		UPDATE worker_profiles
		SET category_id = $1, location = $2, rate = $3, experience_years = $4, bio = $5, updated_at = NOW()
		WHERE id = $6 AND user_id = $7`,
		req.CategoryID, req.Location, req.Rate, req.ExperienceYears, req.Bio, profileID, userID,
	)
	if err != nil || result.RowsAffected() == 0 {
		writeError(w, http.StatusBadRequest, "could not update profile")
		return
	}
	
	// Invalidate caches
	h.Cache.Delete("profile:" + profileID)
	h.Cache.InvalidatePattern("profiles:")
	
	writeJSON(w, http.StatusOK, map[string]string{"id": profileID})
}

