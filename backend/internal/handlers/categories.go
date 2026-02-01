package handlers

import (
	"context"
	"net/http"
)

type categoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(context.Background(), `SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch categories")
		return
	}
	defer rows.Close()

	var categories []categoryResponse
	for rows.Next() {
		var c categoryResponse
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to parse categories")
			return
		}
		categories = append(categories, c)
	}
	writeJSON(w, http.StatusOK, categories)
}
