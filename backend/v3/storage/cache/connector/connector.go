// Package connector provides glue between the [cache.Cache] interface and implementations from the connector sub-packages.
package connector

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/storage/cache"
	"github.com/zitadel/zitadel/backend/v3/storage/cache/connector/gomap"
	"github.com/zitadel/zitadel/backend/v3/storage/cache/connector/noop"
)

type CachesConfig struct {
	Connectors struct {
		Memory gomap.Config
	}
	Instance         *cache.Config
	Milestones       *cache.Config
	Organization     *cache.Config
	IdPFormCallbacks *cache.Config
}

type Connectors struct {
	Config CachesConfig
	Memory *gomap.Connector
}

func StartConnectors(conf *CachesConfig) (Connectors, error) {
	if conf == nil {
		return Connectors{}, nil
	}
	return Connectors{
		Config: *conf,
		Memory: gomap.NewConnector(conf.Connectors.Memory),
	}, nil
}

func StartCache[I ~int, K ~string, V cache.Entry[I, K]](background context.Context, indices []I, purpose cache.Purpose, conf *cache.Config, connectors Connectors) (cache.Cache[I, K, V], error) {
	if conf == nil || conf.Connector == cache.ConnectorUnspecified {
		return noop.NewCache[I, K, V](), nil
	}
	if conf.Connector == cache.ConnectorMemory && connectors.Memory != nil {
		c := gomap.NewCache[I, K, V](background, indices, *conf)
		connectors.Memory.Config.StartAutoPrune(background, c, purpose)
		return c, nil
	}

	return nil, fmt.Errorf("cache connector %q not enabled", conf.Connector)
}
