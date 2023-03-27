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

type emitter struct {
	enabled   bool
	ctx       context.Context
	debouncer *debouncer
	emitter   LogEmitter
	clock     clock.Clock
}

type LogRecord interface {
	Normalize() LogRecord
}

type LogRecordFunc func() LogRecord

func (r LogRecordFunc) Normalize() LogRecord {
	return r()
}

type LogEmitter interface {
	Emit(ctx context.Context, bulk []LogRecord) error
}

type LogEmitterFunc func(ctx context.Context, bulk []LogRecord) error

func (l LogEmitterFunc) Emit(ctx context.Context, bulk []LogRecord) error {
	return l(ctx, bulk)
}

type LogCleanupper interface {
	LogEmitter
	Cleanup(ctx context.Context, keep time.Duration) error
}

// NewEmitter accepts Clock from github.com/benbjohnson/clock so we can control timers and tickers in the unit tests
func NewEmitter(ctx context.Context, clock clock.Clock, cfg *EmitterConfig, logger LogEmitter) (*emitter, error) {
	svc := &emitter{
		enabled: cfg != nil && cfg.Enabled,
		ctx:     ctx,
		emitter: logger,
		clock:   clock,
	}

	if !svc.enabled {
		return svc, nil
	}

	if cfg.Debounce != nil && (cfg.Debounce.MinFrequency > 0 || cfg.Debounce.MaxBulkSize > 0) {
		svc.debouncer = newDebouncer(ctx, *cfg.Debounce, clock, newStorageBulkSink(svc.emitter))
	}

	cleanupper, ok := logger.(LogCleanupper)
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

func (s *emitter) startCleanupping(cleanupper LogCleanupper, cleanupInterval, keep time.Duration) {
	for range s.clock.Tick(cleanupInterval) {
		if err := cleanupper.Cleanup(s.ctx, keep); err != nil {
			logging.WithError(err).Error("cleaning up logs failed")
		}
	}
}

func (s *emitter) Emit(ctx context.Context, record LogRecord) (err error) {
	if !s.enabled {
		return nil
	}

	if s.debouncer != nil {
		s.debouncer.add(record)
		return nil
	}

	return s.emitter.Emit(ctx, []LogRecord{record})
}

func newStorageBulkSink(emitter LogEmitter) bulkSinkFunc {
	return func(ctx context.Context, bulk []LogRecord) error {
		return emitter.Emit(ctx, bulk)
	}
}
