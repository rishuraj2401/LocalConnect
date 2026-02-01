package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"localconnect/internal/analytics"
	"localconnect/internal/auth"
	"localconnect/internal/cache"
	"localconnect/internal/config"
	"localconnect/internal/notifications"
	"localconnect/internal/worker"
)

type Handler struct {
	DB       *pgxpool.Pool
	Config   config.Config
	Worker   *worker.Pool
	Auth     func(http.Handler) http.Handler
	Cache    *cache.Cache
	Notifier *notifications.Notifier
	Tracker  *analytics.Tracker
}

func New(db *pgxpool.Pool, cfg config.Config, pool *worker.Pool) *Handler {
	return &Handler{
		DB:       db,
		Config:   cfg,
		Worker:   pool,
		Auth:     auth.Middleware(cfg.JWTSecret),
		Cache:    cache.New(5 * 60), // 5 minute TTL
		Notifier: notifications.NewNotifier(1000, 50, 2),
		Tracker:  analytics.NewTracker(5000),
	}
}

// Shutdown gracefully shuts down all services
func (h *Handler) Shutdown() {
	h.Cache.Close()
	h.Notifier.Shutdown()
	h.Tracker.Shutdown()
}
