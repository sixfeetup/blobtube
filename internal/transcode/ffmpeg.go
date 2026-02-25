package transcode

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type FFmpeg struct {
	Path               string
	Exec               ExecFunc
	Logger             zerolog.Logger
	MaxDurationSeconds int
}

type HLSRequest struct {
	InputURL string

	// OutputDir defaults to a temp dir under os.TempDir (usually /tmp).
	OutputDir string

	Width  int
	Height int

	PlaylistName           string
	SegmentDurationSeconds int
	VideoPreset            int
	VideoCRF               int
	VideoBitrate           string
	DisableAudio           bool
	AudioBitrate           string
	ExtraArgs              []string
}

type HLSResult struct {
	OutputDir    string
	PlaylistPath string
	Stdout       []byte
	Stderr       []byte
}

func NewFFmpeg(path string, logger zerolog.Logger) *FFmpeg {
	f := &FFmpeg{
		Path:               path,
		Logger:             logger,
		MaxDurationSeconds: 3600,
	}
	f.Exec = f.defaultExec
	return f
}

func (f *FFmpeg) TranscodeHLS(ctx context.Context, req HLSRequest) (HLSResult, error) {
	if req.InputURL == "" {
		return HLSResult{}, fmt.Errorf("input url is required")
	}

	width := req.Width
	height := req.Height
	if width <= 0 {
		width = 128
	}
	if height <= 0 {
		height = 128
	}

	segmentDuration := req.SegmentDurationSeconds
	if segmentDuration <= 0 {
		segmentDuration = 4
	}

	preset := req.VideoPreset
	if preset == 0 {
		preset = 8
	}

	crf := req.VideoCRF
	if crf == 0 {
		crf = 35
	}

	playlistName := req.PlaylistName
	if playlistName == "" {
		playlistName = "index.m3u8"
	}

	outDir := req.OutputDir
	if outDir == "" {
		tmp, err := os.MkdirTemp("", "blobtube-hls-")
		if err != nil {
			return HLSResult{}, fmt.Errorf("create temp output dir: %w", err)
		}
		outDir = tmp
	} else {
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return HLSResult{}, fmt.Errorf("create output dir: %w", err)
		}
	}

	playlistPath := filepath.Join(outDir, playlistName)
	segmentPattern := filepath.Join(outDir, "segment_%05d.ts")

	args := []string{
		"-hide_banner",
		"-y",
		"-i",
		req.InputURL,
		"-t",
		strconv.Itoa(f.maxDurationSeconds()),
		"-vf",
		fmt.Sprintf("scale=%d:%d:flags=lanczos", width, height),
		"-c:v",
		"libsvtav1",
		"-preset",
		strconv.Itoa(preset),
		"-crf",
		strconv.Itoa(crf),
		"-pix_fmt",
		"yuv420p",
		"-sc_threshold",
		"0",
		"-force_key_frames",
		fmt.Sprintf("expr:gte(t,n_forced*%d)", segmentDuration),
	}

	if strings.TrimSpace(req.VideoBitrate) != "" {
		args = append(args,
			"-b:v",
			strings.TrimSpace(req.VideoBitrate),
		)
	}

	if req.DisableAudio {
		args = append(args, "-an")
	} else {
		bitrate := strings.TrimSpace(req.AudioBitrate)
		if bitrate == "" {
			bitrate = "48k"
		}
		args = append(args,
			"-c:a", "aac",
			"-b:a", bitrate,
		)
	}

	args = append(args,
		"-f",
		"hls",
		"-hls_time",
		strconv.Itoa(segmentDuration),
		"-hls_list_size",
		"0",
		"-hls_flags",
		"independent_segments",
		"-hls_segment_filename",
		segmentPattern,
	)

	if len(req.ExtraArgs) > 0 {
		args = append(args, req.ExtraArgs...)
	}

	args = append(args, playlistPath)

	stdout, stderr, err := f.Exec(ctx, f.Path, args...)
	if err != nil {
		trimmed := strings.TrimSpace(string(stderr))
		if trimmed == "" {
			return HLSResult{OutputDir: outDir, PlaylistPath: playlistPath, Stdout: stdout, Stderr: stderr}, err
		}
		return HLSResult{OutputDir: outDir, PlaylistPath: playlistPath, Stdout: stdout, Stderr: stderr}, fmt.Errorf("ffmpeg failed: %s", trimmed)
	}

	return HLSResult{OutputDir: outDir, PlaylistPath: playlistPath, Stdout: stdout, Stderr: stderr}, nil
}

func (f *FFmpeg) maxDurationSeconds() int {
	if f.MaxDurationSeconds <= 0 {
		return 3600
	}
	return f.MaxDurationSeconds
}

func (f *FFmpeg) defaultExec(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return stdout.Bytes(), stderr.Bytes(), nil
	}

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return stdout.Bytes(), stderr.Bytes(), err
	}
	return stdout.Bytes(), stderr.Bytes(), err
}
