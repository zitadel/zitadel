package debouncer

import (
	"sync"
	"time"
)

type asyncDebouncer struct {
	mux               sync.Mutex
	cfg               *Config
	shipper           Shipper
	cache             []any
	cacheLen          uint
	shipSynchronously bool
	ticker            *time.Ticker
}

func newAsyncDebouncer(cfg *Config, ship Shipper) *asyncDebouncer {
	a := &asyncDebouncer{
		cfg:     cfg,
		shipper: ship,
	}

	if cfg.MinFrequency > 0 {
		a.ticker = time.NewTicker(cfg.MinFrequency)
		go a.shipWhenTimerFires()
	}
	return a
}

func (a *asyncDebouncer) Add(item any) {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.cache = append(a.cache, item)
	a.cacheLen++
	if a.cacheLen >= a.cfg.MaxBulkSize {
		// Add should not block and release the lock
		go a.ship()
	}
}

func (a *asyncDebouncer) ship() {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.shipper.Ship(a.cache)
	a.cache = nil
	a.cacheLen = 0
	a.ticker.Reset(a.cfg.MinFrequency)
}

func (a *asyncDebouncer) shipWhenTimerFires() {
	for range a.ticker.C {
		a.ship()
	}
}
