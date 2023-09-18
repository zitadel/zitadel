package logstore

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
)

type EmitterConfig struct {
	Enabled  bool
	Debounce *DebouncerConfig
}

type emitter[T LogRecord[T]] struct {
	enabled   bool
	ctx       context.Context
	debouncer *debouncer[T]
	emitter   LogEmitter[T]
	clock     clock.Clock
}

type LogRecord[T any] interface {
	Normalize() T
}

type LogRecordFunc[T any] func() T

func (r LogRecordFunc[T]) Normalize() T {
	return r()
}

type LogEmitter[T LogRecord[T]] interface {
	Emit(ctx context.Context, bulk []T) error
}

type LogEmitterFunc[T LogRecord[T]] func(ctx context.Context, bulk []T) error

func (l LogEmitterFunc[T]) Emit(ctx context.Context, bulk []T) error {
	return l(ctx, bulk)
}

type LogCleanupper[T LogRecord[T]] interface {
	Cleanup(ctx context.Context, keep time.Duration) error
	LogEmitter[T]
}

// NewEmitter accepts Clock from github.com/benbjohnson/clock so we can control timers and tickers in the unit tests
func NewEmitter[T LogRecord[T]](ctx context.Context, clock clock.Clock, cfg *EmitterConfig, logger LogEmitter[T]) (*emitter[T], error) {
	svc := &emitter[T]{
		enabled: cfg != nil && cfg.Enabled,
		ctx:     ctx,
		emitter: logger,
		clock:   clock,
	}

	if !svc.enabled {
		return svc, nil
	}

	if cfg.Debounce != nil && (cfg.Debounce.MinFrequency > 0 || cfg.Debounce.MaxBulkSize > 0) {
		svc.debouncer = newDebouncer[T](ctx, *cfg.Debounce, clock, newStorageBulkSink(svc.emitter))
	}
	return svc, nil
}

func (s *emitter[T]) Emit(ctx context.Context, record T) (err error) {
	if !s.enabled {
		return nil
	}

	if s.debouncer != nil {
		s.debouncer.add(record)
		return nil
	}

	return s.emitter.Emit(ctx, []T{record})
}

func newStorageBulkSink[T LogRecord[T]](emitter LogEmitter[T]) bulkSinkFunc[T] {
	return func(ctx context.Context, bulk []T) error {
		return emitter.Emit(ctx, bulk)
	}
}
