package transcode

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/rs/zerolog"
)

func TestTranscodeMultiQualityHLS_SpawnsAllVariants(t *testing.T) {
	ff := NewFFmpeg("ffmpeg", zerolog.Nop())

	mu := sync.Mutex{}
	calls := 0
	ff.Exec = func(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
		_ = ctx
		_ = name
		_ = args
		mu.Lock()
		calls++
		mu.Unlock()
		return nil, nil, nil
	}

	outDir := t.TempDir()
	res, err := TranscodeMultiQualityHLS(context.Background(), zerolog.Nop(), ff, "input", outDir, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 ffmpeg calls, got %d", calls)
	}
	if len(res.Errors) != 0 {
		t.Fatalf("expected no errors, got %v", res.Errors)
	}
	if len(res.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res.Results))
	}
}

func TestTranscodeMultiQualityHLS_ContinuesOnFailure(t *testing.T) {
	ff := NewFFmpeg("ffmpeg", zerolog.Nop())

	ff.Exec = func(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
		_ = ctx
		_ = name
		playlist := args[len(args)-1]
		if strings.Contains(playlist, string(filepath.Separator)+string(Quality128)+string(filepath.Separator)) {
			return nil, []byte("nope"), errors.New("exit status 1")
		}
		return nil, nil, nil
	}

	outDir := t.TempDir()
	res, err := TranscodeMultiQualityHLS(context.Background(), zerolog.Nop(), ff, "input", outDir, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if _, ok := res.Errors[Quality128]; !ok {
		t.Fatalf("expected error for %s", Quality128)
	}
	if _, ok := res.Results[Quality64]; !ok {
		t.Fatalf("expected result for %s", Quality64)
	}
	if _, ok := res.Results[Quality256]; !ok {
		t.Fatalf("expected result for %s", Quality256)
	}
}

func TestTranscodeMultiQualityHLS_RunsVariantsConcurrently(t *testing.T) {
	ff := NewFFmpeg("ffmpeg", zerolog.Nop())

	started := make(chan struct{}, 3)
	release := make(chan struct{})
	ff.Exec = func(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
		_ = ctx
		_ = name
		_ = args
		started <- struct{}{}
		<-release
		return nil, nil, nil
	}

	outDir := t.TempDir()
	resCh := make(chan MultiQualityResult, 1)
	errCh := make(chan error, 1)
	go func() {
		res, err := TranscodeMultiQualityHLS(context.Background(), zerolog.Nop(), ff, "input", outDir, nil)
		resCh <- res
		errCh <- err
	}()

	for i := 0; i < 3; i++ {
		<-started
	}
	close(release)

	res := <-resCh
	err := <-errCh
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(res.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res.Results))
	}
}
