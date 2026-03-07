package risk

import (
	"context"
	"sync"
	"time"
)

type Store interface {
	Snapshot(ctx context.Context, signal Signal) (Snapshot, error)
	Save(ctx context.Context, signal Signal, findings []Finding) error
}

type MemoryStore struct {
	mu             sync.RWMutex
	cfg            Config
	userSignals    map[string][]RecordedSignal
	sessionSignals map[string][]RecordedSignal
}

func NewMemoryStore(cfg Config) *MemoryStore {
	return &MemoryStore{
		cfg:            cfg,
		userSignals:    make(map[string][]RecordedSignal),
		sessionSignals: make(map[string][]RecordedSignal),
	}
}

func (s *MemoryStore) Snapshot(_ context.Context, signal Signal) (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cutoff := signal.Timestamp.Add(-maxDuration(s.cfg.HistoryWindow, s.cfg.ContextChangeWindow))
	if signal.Timestamp.IsZero() {
		cutoff = time.Now().UTC().Add(-maxDuration(s.cfg.HistoryWindow, s.cfg.ContextChangeWindow))
	}

	var snapshot Snapshot
	if signal.UserID != "" {
		snapshot.UserSignals = filterSignals(s.userSignals[signal.UserID], cutoff)
	}
	if signal.SessionID != "" {
		snapshot.SessionSignals = filterSignals(s.sessionSignals[signal.SessionID], cutoff)
	}
	return snapshot, nil
}

func (s *MemoryStore) Save(_ context.Context, signal Signal, findings []Finding) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := signal.Timestamp.Add(-maxDuration(s.cfg.HistoryWindow, s.cfg.ContextChangeWindow))
	if signal.Timestamp.IsZero() {
		cutoff = time.Now().UTC().Add(-maxDuration(s.cfg.HistoryWindow, s.cfg.ContextChangeWindow))
	}

	record := RecordedSignal{Signal: signal, Findings: append([]Finding(nil), findings...)}
	if signal.UserID != "" {
		records := append(s.userSignals[signal.UserID], record)
		s.userSignals[signal.UserID] = pruneSignals(records, cutoff, s.cfg.MaxSignalsPerUser)
	}
	if signal.SessionID != "" {
		records := append(s.sessionSignals[signal.SessionID], record)
		s.sessionSignals[signal.SessionID] = pruneSignals(records, cutoff, s.cfg.MaxSignalsPerSession)
	}
	return nil
}

func filterSignals(signals []RecordedSignal, cutoff time.Time) []RecordedSignal {
	filtered := make([]RecordedSignal, 0, len(signals))
	for _, signal := range signals {
		if signal.Timestamp.IsZero() || signal.Timestamp.Before(cutoff) {
			continue
		}
		filtered = append(filtered, signal)
	}
	return filtered
}

func pruneSignals(signals []RecordedSignal, cutoff time.Time, max int) []RecordedSignal {
	pruned := filterSignals(signals, cutoff)
	if len(pruned) <= max {
		return pruned
	}
	return pruned[len(pruned)-max:]
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
