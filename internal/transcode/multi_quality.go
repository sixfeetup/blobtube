package transcode

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
)

type QualityTier string

const (
	Quality64  QualityTier = "64x64"
	Quality128 QualityTier = "128x128"
	Quality256 QualityTier = "256x256"
)

type VariantConfig struct {
	Tier         QualityTier
	Width        int
	Height       int
	VideoBitrate string
}

type MultiQualityResult struct {
	OutputDir string
	Results   map[QualityTier]HLSResult
	Errors    map[QualityTier]error
}

func DefaultVariantConfigs() []VariantConfig {
	return []VariantConfig{
		{Tier: Quality64, Width: 64, Height: 64, VideoBitrate: "50k"},
		{Tier: Quality128, Width: 128, Height: 128, VideoBitrate: "100k"},
		{Tier: Quality256, Width: 256, Height: 256, VideoBitrate: "200k"},
	}
}

func TranscodeMultiQualityHLS(ctx context.Context, logger zerolog.Logger, ff *FFmpeg, inputURL string, outputDir string, variants []VariantConfig) (MultiQualityResult, error) {
	if ff == nil {
		return MultiQualityResult{}, fmt.Errorf("ffmpeg is required")
	}
	if inputURL == "" {
		return MultiQualityResult{}, fmt.Errorf("input url is required")
	}
	if outputDir == "" {
		return MultiQualityResult{}, fmt.Errorf("output dir is required")
	}
	if len(variants) == 0 {
		variants = DefaultVariantConfigs()
	}

	res := MultiQualityResult{
		OutputDir: outputDir,
		Results:   map[QualityTier]HLSResult{},
		Errors:    map[QualityTier]error{},
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, v := range variants {
		v := v
		wg.Add(1)
		go func() {
			defer wg.Done()

			out := filepath.Join(outputDir, string(v.Tier))
			logger.Debug().Str("tier", string(v.Tier)).Str("dir", out).Msg("ffmpeg transcode starting")

			hlsRes, err := ff.TranscodeHLS(ctx, HLSRequest{
				InputURL:               inputURL,
				OutputDir:              out,
				Width:                  v.Width,
				Height:                 v.Height,
				VideoBitrate:           v.VideoBitrate,
				PlaylistName:           "index.m3u8",
				DisableAudio:           false,
				AudioBitrate:           "32k",
				VideoPreset:            8,
				VideoCRF:               35,
				SegmentDurationSeconds: 4,
			})

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				res.Errors[v.Tier] = err
				res.Results[v.Tier] = hlsRes
				logger.Error().Str("tier", string(v.Tier)).Err(err).Msg("ffmpeg transcode failed")
				return
			}
			res.Results[v.Tier] = hlsRes
			logger.Debug().Str("tier", string(v.Tier)).Msg("ffmpeg transcode completed")
		}()
	}

	wg.Wait()

	return res, nil
}
