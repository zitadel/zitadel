package query

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/gomap"
)

type CachesConfig struct {
	Instance cache.CacheConfig
}

func (conf CachesConfig) Start(background context.Context) (_ Caches, err error) {
	caches := Caches{
		instance: gomap.NewCache[instanceIndex, string, *authzInstance](
			background,
			instanceIndexValues(),
			conf.Instance,
		),
	}
	caches.registerInstanceInvalidation()
	return caches, nil
}

type Caches struct {
	instance cache.Cache[instanceIndex, string, *authzInstance]
}

type invalidator[I comparable] interface {
	Invalidate(ctx context.Context, index I, key ...string) error
}

func cacheInvalidationFunc[I comparable](cache invalidator[I], index I) func(context.Context, ...string) {
	return func(ctx context.Context, aggregateIDs ...string) {
		err := cache.Invalidate(ctx, index, aggregateIDs...)
		logging.OnError(err).Warn("cache invalidation failed")
	}
}
