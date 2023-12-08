package bigcache

import (
	"bytes"
	"encoding/gob"
	"reflect"

	a_cache "github.com/allegro/bigcache"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type Bigcache struct {
	cache *a_cache.BigCache
}

func NewBigcache(c *Config) (*Bigcache, error) {
	cacheConfig := a_cache.DefaultConfig(c.CacheLifetime)
	cacheConfig.HardMaxCacheSize = c.MaxCacheSizeInMB
	cache, err := a_cache.NewBigCache(cacheConfig)
	if err != nil {
		return nil, err
	}

	return &Bigcache{
		cache: cache,
	}, nil
}

func (c *Bigcache) Set(key string, object interface{}) error {
	if key == "" || reflect.ValueOf(object).IsNil() {
		return zerrors.ThrowInvalidArgument(nil, "BIGCA-du73s", "key or value should not be empty")
	}
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(object); err != nil {
		return zerrors.ThrowInvalidArgument(err, "BIGCA-RUyxI", "unable to encode object")
	}
	return c.cache.Set(key, b.Bytes())
}

func (c *Bigcache) Get(key string, ptrToObject interface{}) error {
	if key == "" || reflect.ValueOf(ptrToObject).IsNil() {
		return zerrors.ThrowInvalidArgument(nil, "BIGCA-dksoe", "key or value should not be empty")
	}
	value, err := c.cache.Get(key)
	if err == a_cache.ErrEntryNotFound {
		return zerrors.ThrowNotFound(err, "BIGCA-we32s", "not in cache")
	}
	if err != nil {
		logging.Log("BIGCA-ftofbc").WithError(err).Info("read from cache failed")
		return zerrors.ThrowInvalidArgument(err, "BIGCA-3idls", "error in reading from cache")
	}

	b := bytes.NewBuffer(value)
	dec := gob.NewDecoder(b)

	return dec.Decode(ptrToObject)
}

func (c *Bigcache) Delete(key string) error {
	if key == "" {
		return zerrors.ThrowInvalidArgument(nil, "BIGCA-clsi2", "key should not be empty")
	}
	return c.cache.Delete(key)
}
