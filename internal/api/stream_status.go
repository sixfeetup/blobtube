package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/sixfeetup/blobtube/internal/stream"
)

type streamStatusResponse struct {
	ID                       string    `json:"id"`
	Qualities                []string  `json:"qualities"`
	State                    string    `json:"state"`
	CreatedAt                time.Time `json:"created_at"`
	LastAccess               time.Time `json:"last_access"`
	InactivityTimeoutSeconds int       `json:"inactivity_timeout_seconds"`
	Error                    string    `json:"error,omitempty"`
}

func serveStreamStatus(streams *stream.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if streams == nil {
			http.Error(w, "stream manager not configured", http.StatusServiceUnavailable)
			return
		}

		id := chi.URLParam(r, "id")
		if !streamIDRe.MatchString(id) {
			http.Error(w, "invalid stream id", http.StatusBadRequest)
			return
		}

		s, ok := streams.Get(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		// Status checks count as activity.
		_ = streams.Touch(id, time.Now())

		resp := streamStatusResponse{
			ID:                       s.ID,
			Qualities:                append([]string(nil), s.Qualities...),
			State:                    string(s.State),
			CreatedAt:                s.CreatedAt,
			LastAccess:               s.LastAccess,
			InactivityTimeoutSeconds: int(streams.InactivityTimeout().Seconds()),
			Error:                    s.Error,
		}

		w.Header().Set("Content-Type", "application/json")
		setCORSHeaders(w)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
