package hls

import (
	"fmt"
	"os"
	"sync"

	"github.com/grafov/m3u8"
)

type MediaPlaylist struct {
	mu sync.Mutex
	pl *m3u8.MediaPlaylist
}

func NewMediaPlaylist(window, capacity uint) (*MediaPlaylist, error) {
	pl, err := m3u8.NewMediaPlaylist(window, capacity)
	if err != nil {
		return nil, err
	}
	return &MediaPlaylist{pl: pl}, nil
}

func (m *MediaPlaylist) AppendSegment(uri string, duration float64) error {
	if m == nil || m.pl == nil {
		return fmt.Errorf("playlist is nil")
	}
	if uri == "" {
		return fmt.Errorf("segment uri is required")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.pl.Append(uri, duration, ""); err != nil {
		return err
	}
	return nil
}

func (m *MediaPlaylist) Close() {
	if m == nil || m.pl == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pl.Close()
}

func (m *MediaPlaylist) Bytes() []byte {
	if m == nil || m.pl == nil {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return []byte(m.pl.String())
}

func (m *MediaPlaylist) WriteFile(path string) error {
	b := m.Bytes()
	if b == nil {
		return fmt.Errorf("playlist is nil")
	}
	return os.WriteFile(path, b, 0o644)
}
