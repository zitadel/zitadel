package bigcache

import (
	"github.com/caos/zitadel/internal/cache"
	"time"
)

type Config struct {
	MaxCacheSizeInMB int
	//CacheLifetime if set, entries older than the lifetime will be deleted on cleanup (every minute)
	CacheLifetime time.Duration
}

func (c *Config) NewCache() (cache.Cache, error) {
	return NewBigcache(c)
}
