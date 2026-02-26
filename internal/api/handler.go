package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/sixfeetup/blobtube/internal/config"
	"github.com/sixfeetup/blobtube/internal/stream"
	"github.com/sixfeetup/blobtube/internal/transcode"
)

func NewHandler(cfg config.Config, streams *stream.Manager) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(RequestLogger)

	r.Get("/health", healthHandler)

	// Initialize transcoding components
	ytdlp := transcode.NewYtDLP(cfg.YtDLPPath, log.Logger, cfg.DevMode)
	ffmpeg := transcode.NewFFmpeg("ffmpeg", log.Logger)
	resources := stream.NewResources(log.Logger)

	orch := &StreamOrchestrator{
		cfg:      cfg,
		streams:  streams,
		ytdlp:    ytdlp,
		ffmpeg:   ffmpeg,
		resource: resources,
	}

	r.Route("/api/stream", func(r chi.Router) {
		r.Options("/*", corsPreflight)
		r.Post("/", serveCreateStream(orch))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/status", serveStreamStatus(streams))
			r.Get("/master.m3u8", serveMasterPlaylist(cfg, streams))
			r.Get("/{quality}/index.m3u8", serveMediaPlaylist(cfg, streams))
			r.Get("/{quality}/{segment}", serveSegment(cfg, streams))
		})
	})

	r.Handle("/*", http.FileServer(http.Dir(cfg.StaticDir)))

	return r, nil
}
