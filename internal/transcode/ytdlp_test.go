package transcode

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestYtDLP_Execute_ExtractsMuxedURL(t *testing.T) {
	y := NewYtDLP("yt-dlp", zerolog.Nop(), false)
	y.Exec = func(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
		_ = ctx
		_ = name
		_ = args
		return []byte(`{
  "id": "abc",
  "title": "hello",
  "duration": 12,
  "thumbnail": "https://example/thumb.jpg",
  "url": "https://fallback/url",
  "format_id": "fallback",
  "formats": [
    {"format_id": "v", "url": "https://video-only", "vcodec": "avc1", "acodec": "none", "tbr": 500},
    {"format_id": "mux1", "format_note": "360p", "url": "https://mux-low", "vcodec": "avc1", "acodec": "mp4a", "tbr": 200},
    {"format_id": "mux2", "format_note": "720p", "url": "https://mux-high", "vcodec": "avc1", "acodec": "mp4a", "tbr": 800}
  ]
}`), nil, nil
	}

	info, err := y.Execute(context.Background(), "https://youtube.example/watch?v=abc")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if info.StreamURL != "https://mux-high" {
		t.Fatalf("expected mux-high, got %q", info.StreamURL)
	}
	if info.FormatID != "mux2" {
		t.Fatalf("expected mux2 format id, got %q", info.FormatID)
	}
}

func TestYtDLP_Execute_CachesInDevMode(t *testing.T) {
	calls := 0
	y := NewYtDLP("yt-dlp", zerolog.Nop(), true)
	y.Exec = func(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
		calls++
		return []byte(`{"id":"abc","title":"hello","duration":12,"thumbnail":"t","url":"https://u"}`), nil, nil
	}

	ctx := context.Background()
	_, err := y.Execute(ctx, "https://youtube.example/watch?v=abc")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	_, err = y.Execute(ctx, "https://youtube.example/watch?v=abc")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 exec call, got %d", calls)
	}
}

func TestYtDLP_Execute_ClassifiesErrors(t *testing.T) {
	y := NewYtDLP("yt-dlp", zerolog.Nop(), false)
	y.Exec = func(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
		return nil, []byte("ERROR: Unsupported URL: https://nope"), errors.New("exit status 1")
	}

	_, err := y.Execute(context.Background(), "https://nope")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrUnsupportedURL) {
		t.Fatalf("expected ErrUnsupportedURL, got %v", err)
	}
}

func TestYtDLP_CacheExpires(t *testing.T) {
	y := NewYtDLP("yt-dlp", zerolog.Nop(), true)
	y.setCached("u", StreamInfo{StreamURL: "x"}, 10*time.Millisecond)
	if _, ok := y.getCached("u"); !ok {
		t.Fatalf("expected cache hit")
	}
	time.Sleep(20 * time.Millisecond)
	if _, ok := y.getCached("u"); ok {
		t.Fatalf("expected cache miss")
	}
}
