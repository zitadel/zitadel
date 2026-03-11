package signals

import (
	"context"
	"time"
)

type Store interface {
	Snapshot(ctx context.Context, signal Signal, cfg SnapshotConfig) (Snapshot, error)
	Save(ctx context.Context, signal Signal, findings []RecordedFinding, cfg SnapshotConfig) error
}

// effectiveSnapshotConfig merges cfg with fallback defaults.
func effectiveSnapshotConfig(cfg, fallback SnapshotConfig) SnapshotConfig {
	if cfg.HistoryWindow <= 0 {
		cfg.HistoryWindow = fallback.HistoryWindow
	}
	if cfg.ContextChangeWindow <= 0 {
		cfg.ContextChangeWindow = fallback.ContextChangeWindow
	}
	if cfg.MaxSignalsPerUser <= 0 {
		cfg.MaxSignalsPerUser = fallback.MaxSignalsPerUser
	}
	if cfg.MaxSignalsPerSession <= 0 {
		cfg.MaxSignalsPerSession = fallback.MaxSignalsPerSession
	}
	return cfg
}

// signalCutoff computes the cutoff time for filtering/pruning signals.
func signalCutoff(signalTime time.Time, historyWindow, contextChangeWindow time.Duration) time.Time {
	base := signalTime
	if base.IsZero() {
		base = time.Now().UTC()
	}
	return base.Add(-maxDuration(historyWindow, contextChangeWindow))
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
