package logstore

import (
	"context"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/zitadel/logging"
)

type bulkSink interface {
	sendBulk(ctx context.Context, bulk []LogRecord) error
}

var _ bulkSink = bulkSinkFunc(nil)

type bulkSinkFunc func(ctx context.Context, items []LogRecord) error

func (s bulkSinkFunc) sendBulk(ctx context.Context, items []LogRecord) error {
	return s(ctx, items)
}

type debouncer struct {
	// Storing context.Context in a struct is generally bad practice
	// https://go.dev/blog/context-and-structs
	// However, debouncer starts a go routine that triggers side effects itself.
	// So, there is no incoming context.Context available when these events trigger.
	// The only context we can use for the side effects is the app context.
	// Because this can be cancelled by os signals, it's the better solution than creating new background contexts.
	binarySignaledCtx context.Context
	clock             clock.Clock
	ticker            *clock.Ticker
	mux               sync.Mutex
	cfg               DebouncerConfig
	storage           bulkSink
	cache             []LogRecord
	cacheLen          uint
}

type DebouncerConfig struct {
	MinFrequency time.Duration
	MaxBulkSize  uint
}

func newDebouncer(binarySignaledCtx context.Context, cfg DebouncerConfig, clock clock.Clock, ship bulkSink) *debouncer {
	a := &debouncer{
		binarySignaledCtx: binarySignaledCtx,
		clock:             clock,
		cfg:               cfg,
		storage:           ship,
	}

	if cfg.MinFrequency > 0 {
		a.ticker = clock.Ticker(cfg.MinFrequency)
		go a.shipOnTicks()
	}
	return a
}

func (d *debouncer) add(item LogRecord) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.cache = append(d.cache, item)
	d.cacheLen++
	if d.cfg.MaxBulkSize > 0 && d.cacheLen >= d.cfg.MaxBulkSize {
		// Add should not block and release the lock
		go d.ship()
	}
}

func (d *debouncer) ship() {
	if d.cacheLen == 0 {
		return
	}
	d.mux.Lock()
	defer d.mux.Unlock()
	if err := d.storage.sendBulk(d.binarySignaledCtx, d.cache); err != nil {
		logging.WithError(err).WithField("size", len(d.cache)).Error("storing bulk failed")
	}
	d.cache = nil
	d.cacheLen = 0
	if d.cfg.MinFrequency > 0 {
		d.ticker.Reset(d.cfg.MinFrequency)
	}
}

func (d *debouncer) shipOnTicks() {
	for range d.ticker.C {
		d.ship()
	}
}
