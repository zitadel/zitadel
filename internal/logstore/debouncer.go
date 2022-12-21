package logstore

import (
	"context"
	"sync"
	"time"

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
	ctx               context.Context
	mux               sync.Mutex
	cfg               *DebouncerConfig
	storage           bulkSink
	cache             []LogRecord
	cacheLen          uint
	shipSynchronously bool
	ticker            *time.Ticker
}

type DebouncerConfig struct {
	MinFrequency time.Duration
	MaxBulkSize  uint
}

func newDebouncer(ctx context.Context, cfg *DebouncerConfig, ship bulkSink) *debouncer {
	a := &debouncer{
		ctx:     ctx,
		cfg:     cfg,
		storage: ship,
	}

	if cfg.MinFrequency > 0 {
		a.ticker = time.NewTicker(cfg.MinFrequency)
		go a.shipWhenTimerFires()
	}
	return a
}

func (d *debouncer) add(item LogRecord) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.cache = append(d.cache, item)
	d.cacheLen++
	if d.cacheLen >= d.cfg.MaxBulkSize {
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
	if err := d.storage.sendBulk(d.ctx, d.cache); err != nil {
		logging.WithError(err).Warnf("storing bulk of size %d failed", len(d.cache))
	}
	d.cache = nil
	d.cacheLen = 0
	if d.cfg.MinFrequency > 0 {
		d.ticker.Reset(d.cfg.MinFrequency)
	}
}

func (d *debouncer) shipWhenTimerFires() {
	for range d.ticker.C {
		d.ship()
	}
}
