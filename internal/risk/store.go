package risk

import (
	"context"
	"sort"
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

	cutoff := signalCutoff(signal.Timestamp, s.cfg.HistoryWindow, s.cfg.ContextChangeWindow)

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

	cutoff := signalCutoff(signal.Timestamp, s.cfg.HistoryWindow, s.cfg.ContextChangeWindow)

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

// PruneSessions removes session entries whose most recent signal is older than
// the configured history window. Call periodically to prevent unbounded growth
// of the sessionSignals map from finished sessions.
func (s *MemoryStore) PruneSessions(now time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := now.Add(-maxDuration(s.cfg.HistoryWindow, s.cfg.ContextChangeWindow))
	pruned := 0
	for id, signals := range s.sessionSignals {
		if len(signals) == 0 || signals[len(signals)-1].Timestamp.Before(cutoff) {
			delete(s.sessionSignals, id)
			pruned++
		}
	}
	return pruned
}

// signalCutoff computes the cutoff time for filtering/pruning signals.
func signalCutoff(signalTime time.Time, historyWindow, contextChangeWindow time.Duration) time.Time {
	base := signalTime
	if base.IsZero() {
		base = time.Now().UTC()
	}
	return base.Add(-maxDuration(historyWindow, contextChangeWindow))
}

// filterSignals returns only signals at or after the cutoff time.
// Signals are stored in chronological order, so we use binary search to find
// the cutoff index and return a sub-slice (zero allocation for large histories).
func filterSignals(signals []RecordedSignal, cutoff time.Time) []RecordedSignal {
	if len(signals) == 0 {
		return nil
	}
	// Binary search: find the first signal at or after cutoff.
	idx := sort.Search(len(signals), func(i int) bool {
		return !signals[i].Timestamp.Before(cutoff)
	})
	if idx >= len(signals) {
		return nil
	}
	return signals[idx:]
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
