package hls

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/grafov/m3u8"
)

type Variant struct {
	URI        string
	Bandwidth  uint32
	Resolution string
}

func BuildMasterPlaylist(variants []Variant) ([]byte, error) {
	if len(variants) == 0 {
		return nil, fmt.Errorf("variants are required")
	}

	mp := m3u8.NewMasterPlaylist()
	for _, v := range variants {
		if strings.TrimSpace(v.URI) == "" {
			return nil, fmt.Errorf("variant uri is required")
		}
		if v.Bandwidth == 0 {
			return nil, fmt.Errorf("variant bandwidth is required")
		}
		params := m3u8.VariantParams{
			Bandwidth:  v.Bandwidth,
			Resolution: v.Resolution,
		}
		mp.Append(v.URI, nil, params)
	}

	return []byte(mp.String()), nil
}

func WriteMasterPlaylist(path string, variants []Variant) error {
	b, err := BuildMasterPlaylist(variants)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func ParseBitrate(bitrate string) (uint32, error) {
	s := strings.TrimSpace(strings.ToLower(bitrate))
	if s == "" {
		return 0, fmt.Errorf("bitrate is required")
	}

	// Common shorthand: 50k, 100k, 2m.
	mult := uint64(1)
	if strings.HasSuffix(s, "k") {
		mult = 1000
		s = strings.TrimSuffix(s, "k")
	} else if strings.HasSuffix(s, "m") {
		mult = 1000 * 1000
		s = strings.TrimSuffix(s, "m")
	}

	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parse bitrate %q: %w", bitrate, err)
	}

	bps := n * mult
	if bps == 0 {
		return 0, fmt.Errorf("bitrate must be > 0")
	}
	if bps > uint64(^uint32(0)) {
		return 0, fmt.Errorf("bitrate too large")
	}

	return uint32(bps), nil
}
