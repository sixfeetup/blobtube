package stream

import (
	"strings"
	"testing"
	"time"
)

func TestManager_Create_GeneratesUUIDv4(t *testing.T) {
	m := NewManager(5 * time.Minute)
	s, err := m.Create(time.Unix(0, 0))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if s.State != StateInitializing {
		t.Fatalf("expected initializing, got %q", s.State)
	}
	if parts := strings.Count(s.ID, "-"); parts != 4 {
		t.Fatalf("expected uuid with 4 hyphens, got %q", s.ID)
	}
}

func TestManager_ExpireInactive_MarksTimedOut(t *testing.T) {
	m := NewManager(5 * time.Minute)
	s, err := m.Create(time.Unix(0, 0))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	expired := m.ExpireInactive(time.Unix(0, 0).Add(10 * time.Minute))
	if len(expired) != 1 || expired[0] != s.ID {
		t.Fatalf("expected stream to expire")
	}
	got, ok := m.Get(s.ID)
	if !ok {
		t.Fatalf("expected stream to exist")
	}
	if got.State != StateTimedOut {
		t.Fatalf("expected timed_out, got %q", got.State)
	}
}
