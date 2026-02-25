package stream

import (
	"context"
	"os/exec"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestResources_CleanupStream_StopsProcess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses sleep")
	}

	cmd := exec.Command("sleep", "60")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	r := NewResources(zerolog.Nop())
	r.RegisterProcess("s", cmd)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r.CleanupStream(ctx, "s")

	if err := cmd.Process.Signal(syscall.Signal(0)); err == nil {
		t.Fatalf("expected process to be stopped")
	}
}
