package query

import (
	"context"

	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/gomap"
)

type CachesConfig struct {
	Instance cache.CacheConfig
}

func (conf CachesConfig) Start(background context.Context) (_ Caches, err error) {
	return Caches{
		instance: gomap.NewCache[cache.InstanceIndex, string, *authzInstance](
			background,
			cache.InstanceIndexValues(),
			conf.Instance,
		),
	}, nil
}

type Caches struct {
	instance cache.Cache[cache.InstanceIndex, string, *authzInstance]
}
