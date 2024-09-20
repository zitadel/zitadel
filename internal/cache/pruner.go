package cache

import (
	"context"
	"math/rand"
	"time"

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

func (c AutoPruneConfig) StartAutoPrune(background context.Context, pruner Pruner) (close func()) {
	if c.Interval <= 0 {
		return func() {}
	}
	background, cancel := context.WithCancel(background)
	go c.pruneTimer(background, pruner)
	return cancel
}

func (c *AutoPruneConfig) pruneTimer(background context.Context, pruner Pruner) {
	// randomize the first interval
	timer := time.NewTimer(time.Duration(rand.Int63n(int64(c.Interval))))
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {
		case <-background.Done():
			return
		case <-timer.C:
			timer.Reset(c.Interval)
			err := c.doPrune(background, pruner)
			logging.OnError(err).WithField("name", "").Error("cache auto prune")
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
