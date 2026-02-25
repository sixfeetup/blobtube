package stream

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type Resources struct {
	mu     sync.Mutex
	procs  map[string][]*exec.Cmd
	logger zerolog.Logger
}

func NewResources(logger zerolog.Logger) *Resources {
	return &Resources{procs: map[string][]*exec.Cmd{}, logger: logger}
}

func (r *Resources) RegisterProcess(streamID string, cmd *exec.Cmd) {
	if r == nil || cmd == nil || streamID == "" {
		return
	}
	r.mu.Lock()
	r.procs[streamID] = append(r.procs[streamID], cmd)
	r.mu.Unlock()
}

func (r *Resources) CleanupStream(ctx context.Context, streamID string) {
	if r == nil || streamID == "" {
		return
	}

	r.mu.Lock()
	cmds := append([]*exec.Cmd(nil), r.procs[streamID]...)
	delete(r.procs, streamID)
	r.mu.Unlock()

	for _, cmd := range cmds {
		r.stopCmd(ctx, streamID, cmd)
	}
}

func (r *Resources) CleanupAll(ctx context.Context) {
	if r == nil {
		return
	}

	r.mu.Lock()
	all := r.procs
	r.procs = map[string][]*exec.Cmd{}
	r.mu.Unlock()

	for id, cmds := range all {
		for _, cmd := range cmds {
			r.stopCmd(ctx, id, cmd)
		}
	}
}

func (r *Resources) stopCmd(ctx context.Context, streamID string, cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	proc := cmd.Process
	r.logger.Info().Str("stream_id", streamID).Int("pid", proc.Pid).Msg("stopping process")

	_ = proc.Signal(syscall.SIGTERM)

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		r.logWaitResult(streamID, proc.Pid, err)
		return
	case <-time.After(2 * time.Second):
		r.logger.Warn().Str("stream_id", streamID).Int("pid", proc.Pid).Msg("process did not exit after SIGTERM; killing")
		_ = proc.Kill()
		select {
		case err := <-done:
			r.logWaitResult(streamID, proc.Pid, err)
		case <-time.After(1 * time.Second):
			r.logger.Warn().Str("stream_id", streamID).Int("pid", proc.Pid).Msg("process did not exit after SIGKILL")
		}
		return
	case <-ctx.Done():
		_ = proc.Kill()
		select {
		case <-done:
		case <-time.After(1 * time.Second):
		}
		return
	}
}

func (r *Resources) logWaitResult(streamID string, pid int, err error) {
	if err == nil {
		r.logger.Info().Str("stream_id", streamID).Int("pid", pid).Msg("process stopped")
		return
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		r.logger.Info().Str("stream_id", streamID).Int("pid", pid).Msg("process exited")
		return
	}
	r.logger.Warn().Str("stream_id", streamID).Int("pid", pid).Err(err).Msg("process wait failed")
}

func (r *Resources) DebugCounts() string {
	if r == nil {
		return ""
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return fmt.Sprintf("streams=%d", len(r.procs))
}
