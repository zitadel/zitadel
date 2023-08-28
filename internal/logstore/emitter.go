package logstore

import (
	"context"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/zitadel/logging"
)

type EmitterConfig struct {
	Enabled         bool
	Keep            time.Duration
	CleanupInterval time.Duration
	Debounce        *DebouncerConfig
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

	cleanupper, ok := logger.(LogCleanupper[T])
	if !ok {
		if cfg.Keep != 0 {
			return nil, fmt.Errorf("cleaning up for this storage type is not supported, so keep duration must be 0, but is %d", cfg.Keep)
		}
		if cfg.CleanupInterval != 0 {
			return nil, fmt.Errorf("cleaning up for this storage type is not supported, so cleanup interval duration must be 0, but is %d", cfg.Keep)
		}

		return svc, nil
	}

	if cfg.Keep != 0 && cfg.CleanupInterval != 0 {
		go svc.startCleanupping(cleanupper, cfg.CleanupInterval, cfg.Keep)
	}
	return svc, nil
}

func (s *emitter[T]) startCleanupping(cleanupper LogCleanupper[T], cleanupInterval, keep time.Duration) {
	for range s.clock.Tick(cleanupInterval) {
		if err := cleanupper.Cleanup(s.ctx, keep); err != nil {
			logging.WithError(err).Error("cleaning up logs failed")
		}
	}
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
