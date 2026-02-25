package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sixfeetup/blobtube/internal/config"
	"github.com/sixfeetup/blobtube/internal/stream"
)

func NewHandler(cfg config.Config, streams *stream.Manager) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(RequestLogger)

	r.Get("/health", healthHandler)

	r.Route("/api/stream/{id}", func(r chi.Router) {
		r.Options("/*", corsPreflight)
		r.Get("/status", serveStreamStatus(streams))
		r.Get("/master.m3u8", serveMasterPlaylist(cfg, streams))
		r.Get("/{quality}/index.m3u8", serveMediaPlaylist(cfg, streams))
		r.Get("/{quality}/{segment}", serveSegment(cfg, streams))
	})

	r.Handle("/*", http.FileServer(http.Dir(cfg.StaticDir)))

	return r, nil
}
