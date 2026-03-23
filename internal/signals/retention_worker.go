//go:build cgo

package signals

import (
	"context"
	"log/slog"
	"time"
)

// RetentionWorker periodically deletes signals older than the configured
// per-stream retention period. After pruning, it is recommended to run
// compaction to reclaim Parquet file space.
type RetentionWorker struct {
	store   *DuckLakeStore
	streams StreamsConfig
	cfg     RetentionConfig
	metrics *SignalMetrics
	done    chan struct{}
}

// NewRetentionWorker creates a retention worker.
func NewRetentionWorker(store *DuckLakeStore, streams StreamsConfig, cfg RetentionConfig, m *SignalMetrics) *RetentionWorker {
	interval := cfg.PruneInterval
	if interval <= 0 {
		interval = 6 * time.Hour
	}
	return &RetentionWorker{
		store:   store,
		streams: streams,
		cfg:     RetentionConfig{PruneInterval: interval},
		metrics: m,
		done:    make(chan struct{}),
	}
}

// Start runs the retention loop. It blocks until ctx is cancelled.
func (w *RetentionWorker) Start(ctx context.Context) {
	defer close(w.done)

	ticker := time.NewTicker(w.cfg.PruneInterval)
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
func (w *RetentionWorker) Done() <-chan struct{} {
	return w.done
}

func (w *RetentionWorker) safeRun(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "identity_signals.retention_panic",
				slog.Any("panic", r),
			)
		}
	}()
	w.run(ctx)
}

func (w *RetentionWorker) run(ctx context.Context) {
	for _, stream := range w.streams.EnabledStreams() {
		retention := w.streams.RetentionForStream(stream)
		if retention <= 0 {
			continue // keep forever
		}

		pruned, err := w.store.PruneStream(ctx, "", stream, retention)
		if err != nil {
			slog.ErrorContext(ctx, "identity_signals.retention_prune_failed",
				slog.String("stream", string(stream)),
				slog.String("error", err.Error()),
			)
			continue
		}

		if pruned > 0 {
			w.metrics.RecordPruned(ctx, string(stream), pruned)
			slog.InfoContext(ctx, "identity_signals.retention_pruned",
				slog.String("stream", string(stream)),
				slog.Int64("rows_deleted", pruned),
				slog.Duration("retention", retention),
			)
		}
	}
}
