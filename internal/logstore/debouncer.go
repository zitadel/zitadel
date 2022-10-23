package logstore

import (
	"context"
	"sync"
	"time"

	"github.com/zitadel/logging"
)

type BulkSink interface {
	SendBulk(ctx context.Context, bulk []any) error
}

var _ BulkSink = BulkSinkFunc(nil)

type BulkSinkFunc func(ctx context.Context, items []any) error

func (s BulkSinkFunc) SendBulk(ctx context.Context, items []any) error {
	return s(ctx, items)
}

type Debouncer struct {
	ctx               context.Context
	mux               sync.Mutex
	cfg               *DebouncerConfig
	storage           BulkSink
	cache             []any
	cacheLen          uint
	shipSynchronously bool
	ticker            *time.Ticker
}

type DebouncerConfig struct {
	MinFrequency time.Duration
	MaxBulkSize  uint
}

func NewDebouncer(ctx context.Context, cfg *DebouncerConfig, ship BulkSink) *Debouncer {
	a := &Debouncer{
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

func (d *Debouncer) Add(item any) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.cache = append(d.cache, item)
	d.cacheLen++
	if d.cacheLen >= d.cfg.MaxBulkSize {
		// Add should not block and release the lock
		go d.ship()
	}
}

func (d *Debouncer) ship() {
	if d.cacheLen == 0 {
		return
	}
	d.mux.Lock()
	defer d.mux.Unlock()
	if err := d.storage.SendBulk(d.ctx, d.cache); err != nil {
		logging.WithError(err).Warnf("storing bulk of size %d failed", len(d.cache))
	}
	d.cache = nil
	d.cacheLen = 0
	if d.cfg.MinFrequency > 0 {
		d.ticker.Reset(d.cfg.MinFrequency)
	}
}

func (d *Debouncer) shipWhenTimerFires() {
	for range d.ticker.C {
		d.ship()
	}
}
