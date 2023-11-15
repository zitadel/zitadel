package fastcache

import "github.com/zitadel/zitadel/v2/internal/cache"

type Config struct {
	MaxCacheSizeInByte int
}

func (c *Config) NewCache() (cache.Cache, error) {
	return NewFastcache(c)
}
