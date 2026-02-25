package api

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type statusCapturingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusCapturingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		cw := &statusCapturingResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(cw, r)

		d := time.Since(start)
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", cw.status).
			Dur("duration", d).
			Msg("request")
	})
}
