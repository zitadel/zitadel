package cache

import (
	"context"
	"math/rand"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/zitadel/logging"
)

// Pruner is an optional [Cache] interface.
type Pruner interface {
	// Prune deletes all invalidated or expired objects.
	Prune(ctx context.Context) error
}

type PrunerCache[I, K comparable, V Entry[I, K]] interface {
	Cache[I, K, V]
	Pruner
}

type AutoPruneConfig struct {
	// Interval at which the cache is automatically pruned.
	// 0 or lower disables automatic pruning.
	Interval time.Duration

	// Timeout for an automatic prune.
	// It is recommended to keep the value shorter than AutoPruneInterval
	// 0 or lower disables automatic pruning.
	Timeout time.Duration
}

func (c AutoPruneConfig) StartAutoPrune(background context.Context, pruner Pruner, purpose Purpose) (close func()) {
	return c.startAutoPrune(background, pruner, purpose, clockwork.NewRealClock())
}

func (c *AutoPruneConfig) startAutoPrune(background context.Context, pruner Pruner, purpose Purpose, clock clockwork.Clock) (close func()) {
	if c.Interval <= 0 {
		return func() {}
	}
	background, cancel := context.WithCancel(background)
	// randomize the first interval
	timer := clock.NewTimer(time.Duration(rand.Int63n(int64(c.Interval))))
	go c.pruneTimer(background, pruner, purpose, timer)
	return cancel
}

func (c *AutoPruneConfig) pruneTimer(background context.Context, pruner Pruner, purpose Purpose, timer clockwork.Timer) {
	defer func() {
		if !timer.Stop() {
			<-timer.Chan()
		}
	}()

	for {
		select {
		case <-background.Done():
			return
		case <-timer.Chan():
			err := c.doPrune(background, pruner)
			logging.OnError(err).WithField("purpose", purpose).Error("cache auto prune")
			timer.Reset(c.Interval)
		}
	}
}

func (c *AutoPruneConfig) doPrune(background context.Context, pruner Pruner) error {
	ctx, cancel := context.WithCancel(background)
	defer cancel()
	if c.Timeout > 0 {
		ctx, cancel = context.WithTimeout(background, c.Timeout)
		defer cancel()
	}
	return pruner.Prune(ctx)
}
