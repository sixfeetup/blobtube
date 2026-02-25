package transcode

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestFFmpeg_TranscodeHLS_BuildsCommandDefaults(t *testing.T) {
	f := NewFFmpeg("ffmpeg", zerolog.Nop())

	var gotArgs []string
	f.Exec = func(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
		_ = ctx
		if name != "ffmpeg" {
			t.Fatalf("expected ffmpeg binary, got %q", name)
		}
		gotArgs = append([]string(nil), args...)
		return []byte("ok"), []byte("warn"), nil
	}

	res, err := f.TranscodeHLS(context.Background(), HLSRequest{InputURL: "https://example/video"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.PlaylistPath == "" || !strings.HasSuffix(res.PlaylistPath, string(filepath.Separator)+"index.m3u8") {
		t.Fatalf("expected playlist path to end with index.m3u8, got %q", res.PlaylistPath)
	}
	if res.OutputDir == "" {
		t.Fatalf("expected output dir")
	}
	if !strings.HasPrefix(res.OutputDir, os.TempDir()) {
		t.Fatalf("expected output dir under %q, got %q", os.TempDir(), res.OutputDir)
	}
	if filepath.Dir(res.PlaylistPath) != res.OutputDir {
		t.Fatalf("expected playlist to be in output dir")
	}

	assertHasArgPair(t, gotArgs, "-c:v", "libsvtav1")
	assertHasArgPair(t, gotArgs, "-preset", "8")
	assertHasArgPair(t, gotArgs, "-crf", "35")
	assertHasArgPair(t, gotArgs, "-t", "3600")
	assertHasArgPair(t, gotArgs, "-vf", "scale=128:128:flags=lanczos")
	assertHasArgPair(t, gotArgs, "-f", "hls")
	assertHasArgPair(t, gotArgs, "-hls_time", "4")
	assertHasArgPair(t, gotArgs, "-c:a", "aac")
}

func TestFFmpeg_TranscodeHLS_UsesProvidedOutputDirAndDisablesAudio(t *testing.T) {
	f := NewFFmpeg("ffmpeg", zerolog.Nop())

	var gotArgs []string
	f.Exec = func(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
		_ = ctx
		_ = name
		gotArgs = append([]string(nil), args...)
		return nil, nil, nil
	}

	outDir := t.TempDir()
	res, err := f.TranscodeHLS(context.Background(), HLSRequest{InputURL: "u", OutputDir: outDir, DisableAudio: true, PlaylistName: "out.m3u8"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.OutputDir != outDir {
		t.Fatalf("expected output dir %q, got %q", outDir, res.OutputDir)
	}
	if res.PlaylistPath != filepath.Join(outDir, "out.m3u8") {
		t.Fatalf("unexpected playlist path %q", res.PlaylistPath)
	}
	assertHasArg(t, gotArgs, "-an")
}

func TestFFmpeg_TranscodeHLS_PropagatesStderrOnError(t *testing.T) {
	f := NewFFmpeg("ffmpeg", zerolog.Nop())
	f.Exec = func(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
		_ = ctx
		_ = name
		_ = args
		return []byte("stdout"), []byte("something bad\n"), errors.New("exit status 1")
	}

	res, err := f.TranscodeHLS(context.Background(), HLSRequest{InputURL: "u", OutputDir: t.TempDir()})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "ffmpeg failed: something bad") {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(res.Stdout) != "stdout" {
		t.Fatalf("expected stdout to be returned")
	}
	if string(res.Stderr) != "something bad\n" {
		t.Fatalf("expected stderr to be returned")
	}
}

func assertHasArgPair(t *testing.T, args []string, k, v string) {
	t.Helper()
	for i := 0; i < len(args)-1; i++ {
		if args[i] == k && args[i+1] == v {
			return
		}
	}
	t.Fatalf("expected args to include %q %q; got %v", k, v, args)
}

func assertHasArg(t *testing.T, args []string, want string) {
	t.Helper()
	for _, a := range args {
		if a == want {
			return
		}
	}
	t.Fatalf("expected args to include %q; got %v", want, args)
}
