package hls

import (
	"strings"
	"testing"
)

func TestMediaPlaylist_AppendSegment(t *testing.T) {
	pl, err := NewMediaPlaylist(10, 10)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if err := pl.AppendSegment("segment_00001.ts", 4.0); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	s := string(pl.Bytes())
	if !strings.Contains(s, "#EXTINF:4") {
		t.Fatalf("expected EXTINF line, got:\n%s", s)
	}
	if !strings.Contains(s, "segment_00001.ts") {
		t.Fatalf("expected segment uri")
	}
}
