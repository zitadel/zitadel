package detection

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/signals"
)

// Runtime owns the lifecycle of detection infrastructure components:
// the signal emitter, DuckLake store, and compaction worker. It is
// extracted from Service so that cmd/start can access infrastructure
// without type-asserting the Evaluator back to *Service.
type Runtime struct {
	emitter          *signals.Emitter
	emitterCancel    context.CancelFunc
	duckLakeStore    *signals.DuckLakeStore
	compactionWorker *signals.CompactionWorker
}

// NewRuntime creates the detection infrastructure (DuckLake, emitter,
// compaction) when signal storage is enabled. Returns nil when disabled.
func NewRuntime(cfg Config, pgDSN string) (*Runtime, error) {
	if !cfg.SignalStore.Enabled || pgDSN == "" {
		return nil, nil
	}

	duckLakeStore, err := signals.NewDuckLakeStore(pgDSN, cfg.SnapshotConfig(), cfg.SignalStore.DuckLake)
	if err != nil {
		return nil, fmt.Errorf("ducklake signal store: %w", err)
	}
	duckLakeStore.LogInfo(context.Background())

	emitter := signals.NewEmitter(cfg.SignalStore, duckLakeStore)
	compactionWorker := signals.NewCompactionWorker(duckLakeStore, cfg.SignalStore.DuckLake.CompactionInterval)

	r := &Runtime{
		emitter:          emitter,
		duckLakeStore:    duckLakeStore,
		compactionWorker: compactionWorker,
	}

	// Start the emitter goroutine.
	ctx, cancel := context.WithCancel(context.Background())
	r.emitterCancel = cancel
	go emitter.Start(ctx)

	if instrumentation.IsStreamEnabled(instrumentation.StreamRisk) {
		logging.Info(ctx, "detection.signal_store.started",
			slog.Int("channel_size", cfg.SignalStore.ChannelSize),
			slog.String("mode", "ducklake"),
		)
	}

	return r, nil
}

// Close stops the emitter and closes the DuckLake store. Safe to call
// on a nil receiver.
func (r *Runtime) Close() {
	if r == nil {
		return
	}
	if r.emitterCancel != nil {
		r.emitterCancel()
		<-r.emitter.Done()
	}
	if r.duckLakeStore != nil {
		r.duckLakeStore.Close()
	}
}

// Emitter returns the signal emitter, or nil when the signal store is
// not enabled. Middleware uses this to emit fire-and-forget signals.
func (r *Runtime) Emitter() *signals.Emitter {
	if r == nil {
		return nil
	}
	return r.emitter
}

// DuckLakeStore returns the DuckLake signal store, or nil. Used by the
// Signals API for direct query access.
func (r *Runtime) DuckLakeStore() *signals.DuckLakeStore {
	if r == nil {
		return nil
	}
	return r.duckLakeStore
}

// CompactionWorker returns the compaction worker for registration with
// the job queue, or nil when DuckLake is not enabled.
func (r *Runtime) CompactionWorker() *signals.CompactionWorker {
	if r == nil {
		return nil
	}
	return r.compactionWorker
}
