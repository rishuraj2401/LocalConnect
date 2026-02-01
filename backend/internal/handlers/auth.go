package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"localconnect/internal/auth"
)

type authRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type authResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" || req.Role == "" {
		writeError(w, http.StatusBadRequest, "missing required fields")
		return
	}
	if req.Role != "worker" && req.Role != "client" {
		writeError(w, http.StatusBadRequest, "invalid role")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}
	var userID string
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO users (name, email, phone, password_hash, role)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Name, req.Email, req.Phone, string(hash), req.Role,
	).Scan(&userID)
	if err != nil {
		// Check for specific error types
		errMsg := "could not create user"
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "already exists") {
			errMsg = "email already registered"
		} else if strings.Contains(err.Error(), "relation") && strings.Contains(err.Error(), "does not exist") {
			errMsg = "database not initialized - run migrations first"
		} else {
			// Include actual error for debugging
			errMsg = fmt.Sprintf("could not create user: %v", err)
		}
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}
	token, err := auth.GenerateToken(h.Config.JWTSecret, userID, req.Email, req.Name, req.Role)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	writeJSON(w, http.StatusCreated, authResponse{Token: token, Role: req.Role})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "missing required fields")
		return
	}
	var userID string
	var name string
	var email string
	var role string
	var passwordHash string
	err := h.DB.QueryRow(context.Background(),
		`SELECT id, name, email, role, password_hash FROM users WHERE email = $1`, req.Email,
	).Scan(&userID, &name, &email, &role, &passwordHash)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	token, err := auth.GenerateToken(h.Config.JWTSecret, userID, email, name, role)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	writeJSON(w, http.StatusOK, authResponse{Token: token, Role: role})
}
