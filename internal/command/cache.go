package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/gomap"
	"github.com/zitadel/zitadel/internal/cache/noop"
	"github.com/zitadel/zitadel/internal/cache/pg"
	"github.com/zitadel/zitadel/internal/database"
)

type Caches struct {
	connectors *cacheConnectors
	milestones cache.Cache[milestoneIndex, string, *MilestonesReached]
}

func startCaches(background context.Context, conf *cache.CachesConfig, client *database.DB) (_ *Caches, err error) {
	caches := &Caches{
		milestones: noop.NewCache[milestoneIndex, string, *MilestonesReached](),
	}
	if conf == nil {
		return caches, nil
	}
	caches.connectors, err = startCacheConnectors(background, conf, client)
	if err != nil {
		return nil, err
	}
	caches.milestones, err = startCache[milestoneIndex, string, *MilestonesReached](background, []milestoneIndex{milestoneIndexInstanceID}, "milestones", conf.Instance, caches.connectors)
	if err != nil {
		return nil, err
	}
	return caches, nil
}

type cacheConnectors struct {
	memory   *cache.AutoPruneConfig
	postgres *pgxPoolCacheConnector
}

type pgxPoolCacheConnector struct {
	*cache.AutoPruneConfig
	client *database.DB
}

func startCacheConnectors(_ context.Context, conf *cache.CachesConfig, client *database.DB) (_ *cacheConnectors, err error) {
	connectors := new(cacheConnectors)
	if conf.Connectors.Memory.Enabled {
		connectors.memory = &conf.Connectors.Memory.AutoPrune
	}
	if conf.Connectors.Postgres.Enabled {
		connectors.postgres = &pgxPoolCacheConnector{
			AutoPruneConfig: &conf.Connectors.Postgres.AutoPrune,
			client:          client,
		}
	}
	return connectors, nil
}

func startCache[I ~int, K ~string, V cache.Entry[I, K]](background context.Context, indices []I, name string, conf *cache.CacheConfig, connectors *cacheConnectors) (cache.Cache[I, K, V], error) {
	if conf == nil || conf.Connector == "" {
		return noop.NewCache[I, K, V](), nil
	}
	if strings.EqualFold(conf.Connector, "memory") && connectors.memory != nil {
		c := gomap.NewCache[I, K, V](background, indices, *conf)
		connectors.memory.StartAutoPrune(background, c, name)
		return c, nil
	}
	if strings.EqualFold(conf.Connector, "postgres") && connectors.postgres != nil {
		client := connectors.postgres.client
		c, err := pg.NewCache[I, K, V](background, name, *conf, indices, client.Pool, client.Type())
		if err != nil {
			return nil, fmt.Errorf("query start cache: %w", err)
		}
		connectors.postgres.StartAutoPrune(background, c, name)
		return c, nil
	}

	return nil, fmt.Errorf("cache connector %q not enabled", conf.Connector)
}
