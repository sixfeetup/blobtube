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

// TranscodeHLSFromYtDlpPipe uses yt-dlp to download and pipe video to FFmpeg.
// This avoids 403 Forbidden errors from YouTube's stream URLs.
func (f *FFmpeg) TranscodeHLSFromYtDlpPipe(ctx context.Context, youtubeURL string, ytdlpPath string, req HLSRequest) (HLSResult, error) {
	if youtubeURL == "" {
		return HLSResult{}, fmt.Errorf("youtube url is required")
	}
	if ytdlpPath == "" {
		ytdlpPath = "yt-dlp"
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
	segmentPattern := filepath.Join(outDir, "segment_%05d.m4s")

	// yt-dlp command to download and output to stdout
	ytdlpCmd := exec.CommandContext(ctx, ytdlpPath,
		"--quiet",
		"--no-warnings",
		"--format", "best[acodec!=none][vcodec!=none]/best",
		"--output", "-",
		youtubeURL,
	)

	// FFmpeg command to read from stdin
	ffmpegArgs := []string{
		"-hide_banner",
		"-y",
		"-i", "pipe:0", // Read from stdin
		"-t", strconv.Itoa(f.maxDurationSeconds()),
		"-vf", fmt.Sprintf("scale=%d:%d:flags=lanczos", width, height),
		"-c:v", "libsvtav1",
		"-preset", strconv.Itoa(preset),
		"-crf", strconv.Itoa(crf),
		"-pix_fmt", "yuv420p",
		"-sc_threshold", "0",
		"-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%d)", segmentDuration),
	}

	if strings.TrimSpace(req.VideoBitrate) != "" {
		ffmpegArgs = append(ffmpegArgs, "-b:v", strings.TrimSpace(req.VideoBitrate))
	}

	if req.DisableAudio {
		ffmpegArgs = append(ffmpegArgs, "-an")
	} else {
		bitrate := strings.TrimSpace(req.AudioBitrate)
		if bitrate == "" {
			bitrate = "48k"
		}
		ffmpegArgs = append(ffmpegArgs, "-c:a", "aac", "-b:a", bitrate)
	}

	ffmpegArgs = append(ffmpegArgs,
		"-f", "hls",
		"-hls_time", strconv.Itoa(segmentDuration),
		"-hls_list_size", "0",
		"-hls_segment_type", "fmp4",
		"-hls_fmp4_init_filename", "init.mp4",
		"-hls_flags", "independent_segments",
		"-hls_segment_filename", segmentPattern,
	)

	if len(req.ExtraArgs) > 0 {
		ffmpegArgs = append(ffmpegArgs, req.ExtraArgs...)
	}

	ffmpegArgs = append(ffmpegArgs, playlistPath)

	ffmpegCmd := exec.CommandContext(ctx, f.Path, ffmpegArgs...)

	// Pipe yt-dlp stdout to FFmpeg stdin
	pipe, err := ytdlpCmd.StdoutPipe()
	if err != nil {
		return HLSResult{OutputDir: outDir, PlaylistPath: playlistPath}, fmt.Errorf("create pipe: %w", err)
	}
	ffmpegCmd.Stdin = pipe

	// Capture stderr from both commands
	var ytdlpStderr, ffmpegStderr bytes.Buffer
	ytdlpCmd.Stderr = &ytdlpStderr
	ffmpegCmd.Stderr = &ffmpegStderr

	// Start yt-dlp
	if err := ytdlpCmd.Start(); err != nil {
		return HLSResult{OutputDir: outDir, PlaylistPath: playlistPath}, fmt.Errorf("start yt-dlp: %w", err)
	}

	// Start FFmpeg
	if err := ffmpegCmd.Start(); err != nil {
		ytdlpCmd.Process.Kill() // Kill yt-dlp if FFmpeg fails to start
		return HLSResult{OutputDir: outDir, PlaylistPath: playlistPath}, fmt.Errorf("start ffmpeg: %w", err)
	}

	// Wait for both processes
	ytdlpErr := ytdlpCmd.Wait()
	ffmpegErr := ffmpegCmd.Wait()

	// Check for errors
	if ytdlpErr != nil {
		trimmed := strings.TrimSpace(ytdlpStderr.String())
		if trimmed != "" {
			return HLSResult{OutputDir: outDir, PlaylistPath: playlistPath, Stderr: ffmpegStderr.Bytes()},
				fmt.Errorf("yt-dlp failed: %s", trimmed)
		}
		return HLSResult{OutputDir: outDir, PlaylistPath: playlistPath, Stderr: ffmpegStderr.Bytes()},
			fmt.Errorf("yt-dlp failed: %w", ytdlpErr)
	}

	if ffmpegErr != nil {
		trimmed := strings.TrimSpace(ffmpegStderr.String())
		if trimmed != "" {
			return HLSResult{OutputDir: outDir, PlaylistPath: playlistPath, Stderr: ffmpegStderr.Bytes()},
				fmt.Errorf("ffmpeg failed: %s", trimmed)
		}
		return HLSResult{OutputDir: outDir, PlaylistPath: playlistPath, Stderr: ffmpegStderr.Bytes()}, ffmpegErr
	}

	return HLSResult{OutputDir: outDir, PlaylistPath: playlistPath, Stderr: ffmpegStderr.Bytes()}, nil
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
	segmentPattern := filepath.Join(outDir, "segment_%05d.m4s")

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
		"-hls_segment_type",
		"fmp4",
		"-hls_fmp4_init_filename",
		"init.mp4",
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
