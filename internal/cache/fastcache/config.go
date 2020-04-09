package fastcache

import "github.com/caos/zitadel/internal/cache"

type Config struct {
	MaxCacheSizeInByte int
}

func (c *Config) NewCache() (cache.Cache, error) {
	return NewFastcache(c)
}
