package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sixfeetup/blobtube/internal/config"
)

func NewHandler(cfg config.Config) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(RequestLogger)

	r.Get("/health", healthHandler)
	r.Handle("/*", http.FileServer(http.Dir(cfg.StaticDir)))

	return r, nil
}
