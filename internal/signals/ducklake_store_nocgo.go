//go:build !cgo

// When CGO is disabled the DuckDB driver cannot be linked.
// This stub provides the DuckLakeStore type so the rest of the signals
// package compiles, but NewDuckLakeStore always returns an error.

package signals

import (
	"context"
	"fmt"
	"time"
)

// DuckLakeStore is a stub for non-CGO builds.
type DuckLakeStore struct{}

// NewDuckLakeStore returns an error because DuckDB requires CGO.
func NewDuckLakeStore(_ string, _ DuckLakeConfig) (*DuckLakeStore, error) {
	return nil, fmt.Errorf("ducklake: DuckDB requires CGO_ENABLED=1; rebuild with CGO support")
}

func (s *DuckLakeStore) WriteBatch(_ context.Context, _ []RecordedSignal) error {
	return fmt.Errorf("ducklake: not available (CGO disabled)")
}

func (s *DuckLakeStore) SearchSignals(_ context.Context, _ SignalFilters, _, _ int) ([]RecordedSignal, int64, error) {
	return nil, 0, fmt.Errorf("ducklake: not available (CGO disabled)")
}

func (s *DuckLakeStore) AggregateSignals(_ context.Context, _ SignalFilters, _ AggregateRequest) ([]AggregationBucket, error) {
	return nil, fmt.Errorf("ducklake: not available (CGO disabled)")
}

func (s *DuckLakeStore) PruneStream(_ context.Context, _ string, _ SignalStream, _ time.Duration) (int64, error) {
	return 0, fmt.Errorf("ducklake: not available (CGO disabled)")
}

func (s *DuckLakeStore) Close() error { return nil }

func (s *DuckLakeStore) LogInfo(_ context.Context) {}

// CompactionWorker is a no-op stub for non-CGO builds.
type CompactionWorker struct{ done chan struct{} }

func NewCompactionWorker(_ *DuckLakeStore, _ time.Duration, _ *SignalMetrics) *CompactionWorker {
	ch := make(chan struct{})
	close(ch)
	return &CompactionWorker{done: ch}
}

func (w *CompactionWorker) Start(_ context.Context) {}
func (w *CompactionWorker) Done() <-chan struct{}    { return w.done }

// RetentionWorker is a no-op stub for non-CGO builds.
type RetentionWorker struct{ done chan struct{} }

func NewRetentionWorker(_ *DuckLakeStore, _ StreamsConfig, _ RetentionConfig, _ *SignalMetrics) *RetentionWorker {
	ch := make(chan struct{})
	close(ch)
	return &RetentionWorker{done: ch}
}

func (w *RetentionWorker) Start(_ context.Context) {}
func (w *RetentionWorker) Done() <-chan struct{}    { return w.done }
