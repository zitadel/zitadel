package config

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/cache"
	"github.com/caos/zitadel/internal/cache/bigcache"
	"github.com/caos/zitadel/internal/cache/fastcache"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CacheConfig struct {
	ID     string
	Type   string
	Config cache.Config
}

var caches = map[string]func() cache.Config{
	"bigcache":  func() cache.Config { return &bigcache.Config{} },
	"fastcache": func() cache.Config { return &fastcache.Config{} },
}

func (c *CacheConfig) UnmarshalJSON(data []byte) error {
	var rc struct {
		ID     string
		Type   string
		Config json.RawMessage
	}

	if err := json.Unmarshal(data, &rc); err != nil {
		return status.Errorf(codes.Internal, "%v parse config: %v", "CACHE-vmjS", err)
	}

	c.Type = rc.Type
	c.ID = rc.ID

	var err error
	c.Config, err = newCacheConfig(c.Type, rc.Config)
	if err != nil {
		return status.Errorf(codes.Internal, "%v parse config: %v", "CACHE-Ws9E", err)
	}

	return nil
}

func newCacheConfig(cacheType string, configData []byte) (cache.Config, error) {
	t, ok := caches[cacheType]
	if !ok {
		return nil, status.Errorf(codes.Internal, "%v No config: %v", "CACHE-HMEJ", cacheType)
	}

	cacheConfig := t()
	if len(configData) == 0 {
		return cacheConfig, nil
	}

	if err := json.Unmarshal(configData, cacheConfig); err != nil {
		return nil, status.Errorf(codes.Internal, "%v Could not read conifg: %v", "CACHE-1tSS", err)
	}

	return cacheConfig, nil
}
