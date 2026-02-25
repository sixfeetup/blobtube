package api

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/sixfeetup/blobtube/internal/config"
	"github.com/sixfeetup/blobtube/internal/stream"
)

var streamIDRe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

var segmentRe = regexp.MustCompile(`^segment_\d+\.ts$`)

var allowedQualities = map[string]struct{}{
	"64x64":   {},
	"128x128": {},
	"256x256": {},
}

func serveMasterPlaylist(cfg config.Config, streams *stream.Manager) http.HandlerFunc {
	base := cfg.StreamsDir
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if !streamIDRe.MatchString(id) {
			http.Error(w, "invalid stream id", http.StatusBadRequest)
			return
		}

		p := filepath.Join(base, id, "master.m3u8")
		if _, err := os.Stat(p); err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, "failed to read playlist", http.StatusInternalServerError)
			return
		}

		touchOrRegister(streams, id)
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		setCORSHeaders(w)
		http.ServeFile(w, r, p)
	}
}

func serveMediaPlaylist(cfg config.Config, streams *stream.Manager) http.HandlerFunc {
	base := cfg.StreamsDir
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if !streamIDRe.MatchString(id) {
			http.Error(w, "invalid stream id", http.StatusBadRequest)
			return
		}
		quality := chi.URLParam(r, "quality")
		if _, ok := allowedQualities[quality]; !ok {
			http.Error(w, "invalid quality", http.StatusBadRequest)
			return
		}

		p := filepath.Join(base, id, quality, "index.m3u8")
		if _, err := os.Stat(p); err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, "failed to read playlist", http.StatusInternalServerError)
			return
		}

		touchOrRegister(streams, id)
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		setCORSHeaders(w)
		http.ServeFile(w, r, p)
	}
}

func serveSegment(cfg config.Config, streams *stream.Manager) http.HandlerFunc {
	base := cfg.StreamsDir
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if !streamIDRe.MatchString(id) {
			http.Error(w, "invalid stream id", http.StatusBadRequest)
			return
		}
		quality := chi.URLParam(r, "quality")
		if _, ok := allowedQualities[quality]; !ok {
			http.Error(w, "invalid quality", http.StatusBadRequest)
			return
		}
		seg := chi.URLParam(r, "segment")
		if !segmentRe.MatchString(seg) {
			http.Error(w, "invalid segment", http.StatusBadRequest)
			return
		}

		segPath := filepath.Join(base, id, quality, seg)
		if _, err := os.Stat(segPath); err != nil {
			if os.IsNotExist(err) {
				// If a stream exists (playlist present), treat missing segments as "not ready".
				playlistPath := filepath.Join(base, id, quality, "index.m3u8")
				if _, perr := os.Stat(playlistPath); perr == nil {
					touchOrRegister(streams, id)
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					setCORSHeaders(w)
					w.Header().Set("Retry-After", "1")
					w.WriteHeader(http.StatusAccepted)
					_, _ = w.Write([]byte("segment not ready"))
					return
				}
				http.NotFound(w, r)
				return
			}
			http.Error(w, "failed to read segment", http.StatusInternalServerError)
			return
		}

		touchOrRegister(streams, id)
		w.Header().Set("Content-Type", "video/MP2T")
		setCORSHeaders(w)
		http.ServeFile(w, r, segPath)
	}
}

func corsPreflight(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	// Minimal headers for HLS requests from browsers.
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	h := r.Header.Get("Access-Control-Request-Headers")
	if strings.TrimSpace(h) != "" {
		w.Header().Set("Access-Control-Allow-Headers", h)
	}
	w.WriteHeader(http.StatusNoContent)
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func touchOrRegister(streams *stream.Manager, id string) {
	if streams == nil {
		return
	}
	now := time.Now()
	if streams.Touch(id, now) {
		return
	}
	streams.Register(id, now)
}
