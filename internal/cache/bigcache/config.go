package bigcache

import (
	"github.com/caos/zitadel/internal/cache"
	"time"
)

type Config struct {
	MaxCacheSizeInMB int
	// CacheLifetimeSeconds if set, cache makes cleanup every minute
	CacheLifetime time.Duration
}

func (c *Config) NewCache() (cache.Cache, error) {
	return NewBigcache(c)
}
