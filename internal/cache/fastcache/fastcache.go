package fastcache

import (
	"bytes"
	"encoding/gob"
	"reflect"

	"github.com/VictoriaMetrics/fastcache"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type Fastcache struct {
	cache *fastcache.Cache
}

func NewFastcache(config *Config) (*Fastcache, error) {
	return &Fastcache{
		cache: fastcache.New(config.MaxCacheSizeInByte),
	}, nil
}

func (fc *Fastcache) Set(key string, object interface{}) error {
	if key == "" || reflect.ValueOf(object).IsNil() {
		return zerrors.ThrowInvalidArgument(nil, "FASTC-87dj3", "key or value should not be empty")
	}
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(object); err != nil {
		return zerrors.ThrowInvalidArgument(err, "FASTC-RUyxI", "unable to encode object")
	}
	fc.cache.Set([]byte(key), b.Bytes())
	return nil
}

func (fc *Fastcache) Get(key string, ptrToObject interface{}) error {
	if key == "" || reflect.ValueOf(ptrToObject).IsNil() {
		return zerrors.ThrowInvalidArgument(nil, "FASTC-di8es", "key or value should not be empty")
	}
	data := fc.cache.Get(nil, []byte(key))
	if len(data) == 0 {
		return zerrors.ThrowNotFound(nil, "FASTC-xYzSm", "key not found")
	}

	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)

	return dec.Decode(ptrToObject)
}

func (fc *Fastcache) Delete(key string) error {
	if key == "" {
		return zerrors.ThrowInvalidArgument(nil, "FASTC-lod92", "key should not be empty")
	}
	fc.cache.Del([]byte(key))
	return nil
}
