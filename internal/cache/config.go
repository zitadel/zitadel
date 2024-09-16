package cache

import (
	"context"
	"time"

	"github.com/zitadel/logging"
)

type Config struct {
	// Name used to refer to this cache.
	// May be used for logging or storage specific needs,
	// like a table name.
	Name string

	// Age since an object was added to the cache,
	// after which the object is considered invalid.
	// 0 disables max age checks.
	MaxAge time.Duration

	// Age since last use (Get) of an object,
	// after which the object is considered invalid.
	// 0 disables last use age checks.
	LastUseAge time.Duration

	// Interval at which the cache is automatically pruned.
	// 0 disables automatic pruning.
	AutoPruneInterval time.Duration

	// Timeout for an automatic prune.
	// It is recommended to keep the value shorter than AutoPruneInterval
	// 0 disables timeouts.
	AutoPruneTimeOut time.Duration
}

type Pruner interface {
	// Prune deletes all invalidated or expired objects.
	Prune(ctx context.Context) error
}

func (c Config) StartAutoPrune(background context.Context, pruner Pruner) (close func()) {
	if c.AutoPruneInterval <= 0 {
		return func() {}
	}
	background, cancel := context.WithCancel(background)
	go c.pruneTimer(background, pruner)
	return cancel
}

func (c *Config) pruneTimer(background context.Context, pruner Pruner) {
	ticker := time.NewTicker(c.AutoPruneInterval)
	for {
		select {
		case <-background.Done():
			return
		case <-ticker.C:
			err := c.doPrune(background, pruner)
			logging.OnError(err).WithField("name", "").Error("cache auto prune")
		}
	}
}

func (c *Config) doPrune(background context.Context, pruner Pruner) error {
	ctx, cancel := context.WithCancel(background)
	defer cancel()
	if c.AutoPruneTimeOut > 0 {
		ctx, cancel = context.WithTimeout(background, c.AutoPruneTimeOut)
		defer cancel()
	}
	return pruner.Prune(ctx)
}
