package bigcache

import (
	"encoding/json"
	"github.com/caos/logging"
	"time"

	a_cache "github.com/allegro/bigcache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Bigcache struct {
	cache *a_cache.BigCache
}

func NewBigcache(c *Config) (*Bigcache, error) {
	cacheConfig := a_cache.DefaultConfig(c.CacheLifetimeSeconds * time.Second) // Only clean if HardMaxCacheSize is reached
	cacheConfig.HardMaxCacheSize = c.MaxCacheSizeInMB
	if c.CacheLifetimeSeconds > 0 {
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
		return status.Error(codes.InvalidArgument, "unable to marshall object into json")
	}
	return c.cache.Set(key, marshalled)
}

func (c *Bigcache) Get(key string, ptrToObject interface{}) error {
	value, err := c.cache.Get(key)
	if err == a_cache.ErrEntryNotFound {
		return status.Error(codes.NotFound, "not in cache")
	}
	if err != nil {
		logging.Log("BIGCA-ftofbc").WithError(err).Info("read from cache failed")
		return status.Error(codes.Internal, "error in reading from cache")
	}
	return json.Unmarshal(value, ptrToObject)
}

func (c *Bigcache) Delete(key string) error {
	return c.cache.Delete(key)
}
