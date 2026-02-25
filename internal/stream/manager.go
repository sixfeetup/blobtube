package stream

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type State string

const (
	StateInitializing State = "initializing"
	StateActive       State = "active"
	StateCompleted    State = "completed"
	StateError        State = "error"
	StateTimedOut     State = "timed_out"
)

type Stream struct {
	ID         string    `json:"id"`
	Qualities  []string  `json:"qualities"`
	State      State     `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	LastAccess time.Time `json:"last_access"`
	Error      string    `json:"error,omitempty"`
}

var defaultQualities = []string{"64x64", "128x128", "256x256"}

type Manager struct {
	mu      sync.Mutex
	streams map[string]*Stream
	timeout time.Duration
}

func NewManager(inactivityTimeout time.Duration) *Manager {
	if inactivityTimeout <= 0 {
		inactivityTimeout = 5 * time.Minute
	}
	return &Manager{streams: map[string]*Stream{}, timeout: inactivityTimeout}
}

func (m *Manager) InactivityTimeout() time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.timeout
}

func (m *Manager) Create(now time.Time) (Stream, error) {
	if now.IsZero() {
		now = time.Now()
	}
	id, err := newUUIDv4()
	if err != nil {
		return Stream{}, err
	}

	s := &Stream{ID: id, State: StateInitializing, CreatedAt: now, LastAccess: now}
	s.Qualities = append([]string(nil), defaultQualities...)
	m.mu.Lock()
	m.streams[id] = s
	m.mu.Unlock()

	return *s, nil
}

func (m *Manager) Register(id string, now time.Time) (Stream, bool) {
	if id == "" {
		return Stream{}, false
	}
	if now.IsZero() {
		now = time.Now()
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.streams[id]; ok {
		s.LastAccess = now
		return *s, true
	}

	s := &Stream{ID: id, State: StateActive, CreatedAt: now, LastAccess: now}
	s.Qualities = append([]string(nil), defaultQualities...)
	m.streams[id] = s
	return *s, true
}

func (m *Manager) Touch(id string, now time.Time) bool {
	if id == "" {
		return false
	}
	if now.IsZero() {
		now = time.Now()
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.streams[id]
	if !ok {
		return false
	}
	s.LastAccess = now
	return true
}

func (m *Manager) Get(id string) (Stream, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.streams[id]
	if !ok {
		return Stream{}, false
	}
	return *s, true
}

func (m *Manager) IDs() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	ids := make([]string, 0, len(m.streams))
	for id := range m.streams {
		ids = append(ids, id)
	}
	return ids
}

func (m *Manager) SetState(id string, state State, errMsg string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.streams[id]
	if !ok {
		return false
	}
	s.State = state
	s.Error = errMsg
	return true
}

func (m *Manager) ExpireInactive(now time.Time) []string {
	if now.IsZero() {
		now = time.Now()
	}

	expired := []string{}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range m.streams {
		if s.State == StateCompleted || s.State == StateError || s.State == StateTimedOut {
			continue
		}
		if now.Sub(s.LastAccess) <= m.timeout {
			continue
		}
		s.State = StateTimedOut
		s.Error = "inactive timeout"
		expired = append(expired, s.ID)
	}
	return expired
}

func (m *Manager) StartJanitor(ctx context.Context, interval time.Duration, onTimeout func(streamID string)) {
	if interval <= 0 {
		interval = 30 * time.Second
	}

	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			expired := m.ExpireInactive(now)
			if onTimeout == nil {
				continue
			}
			for _, id := range expired {
				onTimeout(id)
			}
		}
	}
}

func newUUIDv4() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("uuid rand: %w", err)
	}
	// Set version (4) and variant bits.
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	hexStr := hex.EncodeToString(b[:])
	return hexStr[0:8] + "-" + hexStr[8:12] + "-" + hexStr[12:16] + "-" + hexStr[16:20] + "-" + hexStr[20:32], nil
}
