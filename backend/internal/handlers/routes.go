package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(corsMiddleware)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Handle("/media/*", http.StripPrefix("/media/", http.FileServer(http.Dir(h.Config.MediaDir))))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})

	r.Get("/categories", h.ListCategories)
	r.Get("/profiles", h.ListProfiles)
	r.With(h.Auth).Get("/profiles/me", h.GetMyProfile)
	r.Get("/profiles/{id}", h.GetProfile)
	r.With(h.Auth).Post("/profiles", h.CreateProfile)
	r.With(h.Auth).Put("/profiles/{id}", h.UpdateProfile)
	r.Get("/profiles/{id}/media", h.ListMedia)
	r.With(h.Auth).Post("/profiles/{id}/media", h.UploadMedia)
	r.With(h.Auth).Post("/profiles/{id}/reviews", h.CreateReview)
	r.Get("/profiles/{id}/reviews", h.ListReviews)
	r.With(h.Auth).Post("/profiles/{id}/upvote", h.UpvoteProfile)
	r.With(h.Auth).Delete("/profiles/{id}/upvote", h.RemoveUpvote)
	r.With(h.Auth).Post("/profiles/{id}/contact-requests", h.CreateContactRequest)
	r.With(h.Auth).Get("/contact-requests", h.ListContactRequests)
	
	// Messaging routes
	r.With(h.Auth).Post("/messages", h.SendMessage)
	r.With(h.Auth).Get("/conversations", h.GetConversations)
	r.With(h.Auth).Get("/conversations/{id}/messages", h.GetMessages)

	return r
}
