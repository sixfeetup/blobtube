package hls

import (
	"strings"
	"testing"
)

func TestParseBitrate(t *testing.T) {
	got, err := ParseBitrate("50k")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != 50000 {
		t.Fatalf("expected 50000, got %d", got)
	}
}

func TestBuildMasterPlaylist(t *testing.T) {
	b, err := BuildMasterPlaylist([]Variant{
		{URI: "64x64/index.m3u8", Bandwidth: 50000, Resolution: "64x64"},
		{URI: "128x128/index.m3u8", Bandwidth: 100000, Resolution: "128x128"},
		{URI: "256x256/index.m3u8", Bandwidth: 200000, Resolution: "256x256"},
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	s := string(b)
	if !strings.HasPrefix(s, "#EXTM3U") {
		t.Fatalf("expected playlist header")
	}
	if !strings.Contains(s, "#EXT-X-STREAM-INF") {
		t.Fatalf("expected variant stream info")
	}
	if !strings.Contains(s, "64x64/index.m3u8") {
		t.Fatalf("expected 64x64 uri")
	}
}
