package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sixfeetup/blobtube/internal/config"
	"github.com/sixfeetup/blobtube/internal/stream"
)

func TestServeSegment_Returns202WhenPlaylistExistsButSegmentMissing(t *testing.T) {
	root := t.TempDir()
	streamID := "abc123"
	quality := "128x128"

	qDir := filepath.Join(root, streamID, quality)
	if err := os.MkdirAll(qDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(qDir, "index.m3u8"), []byte("#EXTM3U\n"), 0o644); err != nil {
		t.Fatalf("write playlist: %v", err)
	}

	h, err := NewHandler(config.Config{StreamsDir: root, StaticDir: root}, stream.NewManager(5*time.Minute))
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/stream/"+streamID+"/"+quality+"/segment_00001.ts", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("expected CORS header")
	}
}

func TestServeSegment_ServesSegmentFile(t *testing.T) {
	root := t.TempDir()
	streamID := "abc123"
	quality := "64x64"

	qDir := filepath.Join(root, streamID, quality)
	if err := os.MkdirAll(qDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	segPath := filepath.Join(qDir, "segment_00001.ts")
	if err := os.WriteFile(segPath, []byte("abc"), 0o644); err != nil {
		t.Fatalf("write segment: %v", err)
	}

	h, err := NewHandler(config.Config{StreamsDir: root, StaticDir: root}, stream.NewManager(5*time.Minute))
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/stream/"+streamID+"/"+quality+"/segment_00001.ts", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Header().Get("Content-Type") != "video/MP2T" {
		t.Fatalf("unexpected content-type %q", rr.Header().Get("Content-Type"))
	}
}
