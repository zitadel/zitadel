package bigcache

import (
	"time"

	"github.com/zitadel/zitadel/v2/internal/cache"
)

type Config struct {
	MaxCacheSizeInMB int
	//CacheLifetime if set, entries older than the lifetime will be deleted on cleanup (every minute)
	CacheLifetime time.Duration
}

func (c *Config) NewCache() (cache.Cache, error) {
	return NewBigcache(c)
}
