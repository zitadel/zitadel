package bigcache

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"

	a_cache "github.com/allegro/bigcache"
)

type Bigcache struct {
	cache *a_cache.BigCache
}

func NewBigcache(c *Config) (*Bigcache, error) {
	cacheConfig := a_cache.DefaultConfig(c.CacheLifetime)
	cacheConfig.HardMaxCacheSize = c.MaxCacheSizeInMB
	if c.CacheLifetime > 0 {
		cacheConfig.CleanWindow = 1 * time.Minute
	}
	cache, err := a_cache.NewBigCache(cacheConfig)
	if err != nil {
		return nil, err
	}
	return &Bigcache{
		cache: cache,
	}, nil
}

func (c *Bigcache) Set(key string, object interface{}) error {
	marshalled, err := json.Marshal(object)
	if err != nil {
		logging.Log("BIGCA-j6Vkhm").Debug("unable to marshall object into json")
		return errors.ThrowInvalidArgument(err, "BIGCA-ie83s", "unable to marshall object into json")
	}
	return c.cache.Set(key, marshalled)
}

func (c *Bigcache) Get(key string, ptrToObject interface{}) error {
	value, err := c.cache.Get(key)
	if err == a_cache.ErrEntryNotFound {
		return errors.ThrowNotFound(err, "BIGCA-we32s", "not in cache")
	}
	if err != nil {
		logging.Log("BIGCA-ftofbc").WithError(err).Info("read from cache failed")
		return errors.ThrowInvalidArgument(err, "BIGCA-3idls", "error in reading from cache")
	}
	return json.Unmarshal(value, ptrToObject)
}

func (c *Bigcache) Delete(key string) error {
	return c.cache.Delete(key)
}
