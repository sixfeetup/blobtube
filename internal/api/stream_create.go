package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/sixfeetup/blobtube/internal/config"
	"github.com/sixfeetup/blobtube/internal/stream"
	"github.com/sixfeetup/blobtube/internal/transcode"
)

type CreateStreamRequest struct {
	URL string `json:"url"`
}

type CreateStreamResponse struct {
	StreamID string `json:"stream_id"`
	Status   string `json:"status"`
}

type StreamOrchestrator struct {
	cfg      config.Config
	streams  *stream.Manager
	ytdlp    *transcode.YtDLP
	ffmpeg   *transcode.FFmpeg
	resource *stream.Resources
}

func serveCreateStream(orch *StreamOrchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateStreamRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Warn().Err(err).Msg("invalid json")
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			http.Error(w, `{"error":"url is required"}`, http.StatusBadRequest)
			return
		}

		// Create stream entry
		s, err := orch.streams.Create(time.Now())
		if err != nil {
			log.Error().Err(err).Msg("failed to create stream")
			http.Error(w, `{"error":"failed to create stream"}`, http.StatusInternalServerError)
			return
		}

		// Return stream ID immediately
		resp := CreateStreamResponse{
			StreamID: s.ID,
			Status:   string(stream.StateInitializing),
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(resp)

		// Start async processing
		go orch.processStream(s.ID, req.URL)
	}
}

func (orch *StreamOrchestrator) processStream(streamID string, youtubeURL string) {
	logger := log.With().Str("stream_id", streamID).Str("url", youtubeURL).Logger()
	logger.Info().Msg("stream processing started")

	// Extract video info using yt-dlp (no need to get stream URL)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	info, err := orch.ytdlp.Execute(ctx, youtubeURL)
	if err != nil {
		logger.Error().Err(err).Msg("yt-dlp extraction failed")
		orch.streams.SetState(streamID, stream.StateError, fmt.Sprintf("yt-dlp failed: %v", err))
		return
	}

	logger.Info().
		Str("video_id", info.VideoID).
		Str("title", info.Title).
		Int("duration", info.Duration).
		Str("format", info.FormatNote).
		Msg("yt-dlp extraction successful")

	// Create output directory for this stream
	streamDir := filepath.Join(orch.cfg.StreamsDir, streamID)
	if err := os.MkdirAll(streamDir, 0o755); err != nil {
		logger.Error().Err(err).Msg("failed to create stream directory")
		orch.streams.SetState(streamID, stream.StateError, fmt.Sprintf("failed to create directory: %v", err))
		return
	}

	// Start multi-quality transcoding using yt-dlp pipe
	// Instead of passing the stream URL directly, we use yt-dlp to pipe the video
	orch.streams.SetState(streamID, stream.StateActive, "")
	logger.Info().Str("youtube_url", youtubeURL).Msg("starting transcoding via yt-dlp pipe")

	transcodeCtx, transcodeCancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer transcodeCancel()

	result, err := transcode.TranscodeMultiQualityHLSFromYouTube(
		transcodeCtx,
		logger,
		orch.ffmpeg,
		orch.ytdlp,
		youtubeURL,
		streamDir,
		transcode.DefaultVariantConfigs(),
	)

	// Track the ffmpeg processes for cleanup
	// Note: We'll need to enhance the Resources tracking to handle this properly
	// For now, the janitor will clean up based on stream timeout

	if err != nil {
		logger.Error().Err(err).Msg("transcoding initialization failed")
		orch.streams.SetState(streamID, stream.StateError, fmt.Sprintf("transcoding failed: %v", err))
		return
	}

	// Check for errors in individual quality tiers
	hasErrors := false
	for tier, tierErr := range result.Errors {
		if tierErr != nil {
			logger.Error().Str("tier", string(tier)).Err(tierErr).Msg("quality tier failed")
			hasErrors = true
		}
	}

	if hasErrors {
		logger.Warn().Msg("some quality tiers failed, but stream may still be usable")
	}

	// Check if we have at least one successful quality
	if len(result.Results) == 0 || len(result.Results) == len(result.Errors) {
		logger.Error().Msg("all quality tiers failed")
		orch.streams.SetState(streamID, stream.StateError, "all quality tiers failed")
		return
	}

	logger.Info().Msg("transcoding completed successfully")
	orch.streams.SetState(streamID, stream.StateCompleted, "")
}
