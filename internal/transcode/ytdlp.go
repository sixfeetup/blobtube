package transcode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	ErrUnsupportedURL   = errors.New("unsupported url")
	ErrVideoUnavailable = errors.New("video unavailable")
	ErrRegionLocked     = errors.New("region locked")
)

type ExecFunc func(ctx context.Context, name string, args ...string) (stdout []byte, stderr []byte, err error)

type YtDLP struct {
	Path    string
	Exec    ExecFunc
	Logger  zerolog.Logger
	DevMode bool

	mu    sync.Mutex
	cache map[string]cacheEntry
}

type cacheEntry struct {
	info    StreamInfo
	expires time.Time
}

type StreamInfo struct {
	VideoID    string
	Title      string
	Duration   int
	Thumbnail  string
	StreamURL  string
	FormatID   string
	FormatNote string
}

func NewYtDLP(path string, logger zerolog.Logger, devMode bool) *YtDLP {
	y := &YtDLP{
		Path:    path,
		Logger:  logger,
		DevMode: devMode,
		cache:   map[string]cacheEntry{},
	}
	y.Exec = y.defaultExec
	return y
}

func (y *YtDLP) Execute(ctx context.Context, videoURL string) (StreamInfo, error) {
	if videoURL == "" {
		return StreamInfo{}, fmt.Errorf("video url is required")
	}

	if y.DevMode {
		if info, ok := y.getCached(videoURL); ok {
			y.Logger.Debug().Str("url", videoURL).Msg("yt-dlp cache hit")
			return info, nil
		}
	}

	args := []string{
		"-j",
		"--no-warnings",
		"--no-playlist",
		"--skip-download",
		// Prefer a muxed format (audio+video) when possible.
		"-f",
		"best[acodec!=none][vcodec!=none]/best",
		videoURL,
	}

	stdout, stderr, err := y.Exec(ctx, y.Path, args...)
	if err != nil {
		return StreamInfo{}, classifyYtDLPErr(stderr, err)
	}

	var payload ytDLPJSON
	if uerr := json.Unmarshal(stdout, &payload); uerr != nil {
		return StreamInfo{}, fmt.Errorf("parse yt-dlp json: %w", uerr)
	}

	info, err := extractStreamInfo(payload)
	if err != nil {
		return StreamInfo{}, err
	}

	if y.DevMode {
		y.setCached(videoURL, info, 5*time.Minute)
	}

	return info, nil
}

func (y *YtDLP) defaultExec(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	stdout, err := cmd.Output()
	if err == nil {
		return stdout, nil, nil
	}

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return stdout, ee.Stderr, err
	}
	return stdout, nil, err
}

func (y *YtDLP) getCached(videoURL string) (StreamInfo, bool) {
	y.mu.Lock()
	defer y.mu.Unlock()

	e, ok := y.cache[videoURL]
	if !ok {
		return StreamInfo{}, false
	}
	if time.Now().After(e.expires) {
		delete(y.cache, videoURL)
		return StreamInfo{}, false
	}
	return e.info, true
}

func (y *YtDLP) setCached(videoURL string, info StreamInfo, ttl time.Duration) {
	y.mu.Lock()
	defer y.mu.Unlock()

	y.cache[videoURL] = cacheEntry{info: info, expires: time.Now().Add(ttl)}
}

type ytDLPJSON struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Duration  int           `json:"duration"`
	Thumbnail string        `json:"thumbnail"`
	URL       string        `json:"url"`
	FormatID  string        `json:"format_id"`
	Format    string        `json:"format"`
	Formats   []ytDLPFormat `json:"formats"`
}

type ytDLPFormat struct {
	FormatID   string  `json:"format_id"`
	FormatNote string  `json:"format_note"`
	URL        string  `json:"url"`
	VCodec     string  `json:"vcodec"`
	ACodec     string  `json:"acodec"`
	TBR        float64 `json:"tbr"`
	Protocol   string  `json:"protocol"`
}

func extractStreamInfo(p ytDLPJSON) (StreamInfo, error) {
	// If yt-dlp already selected a concrete format, payload.URL is usually sufficient.
	// We still prefer a muxed format from formats when available.
	chosen := chooseBestMuxed(p.Formats)
	streamURL := p.URL
	formatID := p.FormatID
	formatNote := p.Format
	if chosen.URL != "" {
		streamURL = chosen.URL
		formatID = chosen.FormatID
		formatNote = chosen.FormatNote
	}
	if streamURL == "" {
		return StreamInfo{}, errors.New("yt-dlp did not provide a stream url")
	}

	return StreamInfo{
		VideoID:    p.ID,
		Title:      p.Title,
		Duration:   p.Duration,
		Thumbnail:  p.Thumbnail,
		StreamURL:  streamURL,
		FormatID:   formatID,
		FormatNote: formatNote,
	}, nil
}

func chooseBestMuxed(formats []ytDLPFormat) ytDLPFormat {
	muxed := make([]ytDLPFormat, 0, len(formats))
	for _, f := range formats {
		if f.URL == "" {
			continue
		}
		if f.ACodec == "none" || f.VCodec == "none" {
			continue
		}
		muxed = append(muxed, f)
	}
	if len(muxed) == 0 {
		return ytDLPFormat{}
	}

	sort.SliceStable(muxed, func(i, j int) bool {
		return muxed[i].TBR > muxed[j].TBR
	})
	return muxed[0]
}

func classifyYtDLPErr(stderr []byte, err error) error {
	msg := strings.ToLower(string(stderr))
	switch {
	case strings.Contains(msg, "unsupported url"):
		return fmt.Errorf("%w: %s", ErrUnsupportedURL, strings.TrimSpace(string(stderr)))
	case strings.Contains(msg, "video unavailable") || strings.Contains(msg, "private video") || strings.Contains(msg, "this video is private"):
		return fmt.Errorf("%w: %s", ErrVideoUnavailable, strings.TrimSpace(string(stderr)))
	case strings.Contains(msg, "not available in your country") || strings.Contains(msg, "geo"):
		return fmt.Errorf("%w: %s", ErrRegionLocked, strings.TrimSpace(string(stderr)))
	default:
		if len(stderr) > 0 {
			return fmt.Errorf("yt-dlp failed: %s", strings.TrimSpace(string(stderr)))
		}
		return err
	}
}
