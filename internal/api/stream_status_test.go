package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sixfeetup/blobtube/internal/config"
	"github.com/sixfeetup/blobtube/internal/stream"
)

func TestServeStreamStatus_404WhenUnknown(t *testing.T) {
	h, err := NewHandler(config.Config{StaticDir: t.TempDir()}, stream.NewManager(5*time.Minute))
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/stream/abc/status", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestServeStreamStatus_200WhenKnown(t *testing.T) {
	mgr := stream.NewManager(5 * time.Minute)
	_, _ = mgr.Register("abc", time.Unix(0, 0))

	h, err := NewHandler(config.Config{StaticDir: t.TempDir()}, mgr)
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/stream/abc/status", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
