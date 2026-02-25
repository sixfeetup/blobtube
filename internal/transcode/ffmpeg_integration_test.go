//go:build integration

package transcode

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestFFmpeg_TranscodeHLS_GeneratesPlaylistAndSegments(t *testing.T) {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		t.Skip("ffmpeg not found")
	}

	// Best-effort check that the encoder exists, otherwise this test is noisy.
	encOut, _ := exec.Command(ffmpegPath, "-hide_banner", "-encoders").CombinedOutput()
	if !strings.Contains(string(encOut), "libsvtav1") {
		t.Skip("ffmpeg missing libsvtav1 encoder")
	}

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime.Caller failed")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	input := filepath.Join(repoRoot, "beach.mp4")
	if _, statErr := os.Stat(input); statErr != nil {
		t.Fatalf("expected test input file %q to exist: %v", input, statErr)
	}

	f := NewFFmpeg(ffmpegPath, zerolog.Nop())
	f.MaxDurationSeconds = 2

	outDir := t.TempDir()
	res, terr := f.TranscodeHLS(context.Background(), HLSRequest{InputURL: input, OutputDir: outDir, DisableAudio: true})
	if terr != nil {
		t.Fatalf("transcode failed: %v\nstderr:\n%s", terr, string(res.Stderr))
	}

	if _, statErr := os.Stat(res.PlaylistPath); statErr != nil {
		t.Fatalf("expected playlist to exist: %v", statErr)
	}
	segments, gerr := filepath.Glob(filepath.Join(outDir, "segment_*.ts"))
	if gerr != nil {
		t.Fatalf("glob segments: %v", gerr)
	}
	if len(segments) == 0 {
		t.Fatalf("expected at least 1 segment in %q", outDir)
	}
}
