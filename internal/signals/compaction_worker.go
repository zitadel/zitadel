//go:build cgo

package signals

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

// CompactionWorker merges small Parquet files written by DuckLake into
// larger time-aligned files. This reduces file count on S3/filesystem
// and improves query performance.
type CompactionWorker struct {
	store    *DuckLakeStore
	interval time.Duration
	metrics  *SignalMetrics
	done     chan struct{}
}

// NewCompactionWorker creates a compaction worker.
func NewCompactionWorker(store *DuckLakeStore, interval time.Duration, m *SignalMetrics) *CompactionWorker {
	if interval <= 0 {
		interval = 1 * time.Hour
	}
	return &CompactionWorker{
		store:    store,
		interval: interval,
		metrics:  m,
		done:     make(chan struct{}),
	}
}

// Start runs the compaction loop. It blocks until ctx is cancelled.
func (w *CompactionWorker) Start(ctx context.Context) {
	defer close(w.done)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.safeRun(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// Done returns a channel closed when the worker has stopped.
func (w *CompactionWorker) Done() <-chan struct{} {
	return w.done
}

func (w *CompactionWorker) safeRun(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "identity_signals.compaction_panic",
				slog.Any("panic", r),
			)
		}
	}()
	w.run(ctx)
}

func (w *CompactionWorker) run(ctx context.Context) {
	if w.store == nil {
		return
	}

	start := time.Now()
	compacted, err := w.store.Compact(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "identity_signals.compaction_failed",
			slog.String("error", err.Error()),
		)
		return
	}

	if compacted > 0 {
		w.metrics.RecordCompactionDuration(ctx, time.Since(start).Seconds())
		slog.InfoContext(ctx, "identity_signals.compaction_complete",
			slog.Int("files_compacted", compacted),
		)
	}
}

func runCompaction(_ context.Context, _ *sql.DB) (int, error) {
	// Deprecated: use DuckLakeStore.Compact() instead.
	return 0, fmt.Errorf("use DuckLakeStore.Compact()")
}
