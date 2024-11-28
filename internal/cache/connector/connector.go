// Package connector provides glue between the [cache.Cache] interface and implementations from the connector sub-packages.
package connector

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/connector/gomap"
	"github.com/zitadel/zitadel/internal/cache/connector/noop"
	"github.com/zitadel/zitadel/internal/cache/connector/pg"
	"github.com/zitadel/zitadel/internal/cache/connector/redis"
	"github.com/zitadel/zitadel/internal/database"
)

type CachesConfig struct {
	Connectors struct {
		Memory   gomap.Config
		Postgres pg.Config
		Redis    redis.Config
	}
	Instance     *cache.Config
	Milestones   *cache.Config
	Organization *cache.Config
}

type Connectors struct {
	Config   CachesConfig
	Memory   *gomap.Connector
	Postgres *pg.Connector
	Redis    *redis.Connector
}

func StartConnectors(conf *CachesConfig, client *database.DB) (Connectors, error) {
	if conf == nil {
		return Connectors{}, nil
	}
	return Connectors{
		Config:   *conf,
		Memory:   gomap.NewConnector(conf.Connectors.Memory),
		Postgres: pg.NewConnector(conf.Connectors.Postgres, client),
		Redis:    redis.NewConnector(conf.Connectors.Redis),
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
	if conf.Connector == cache.ConnectorPostgres && connectors.Postgres != nil {
		c, err := pg.NewCache[I, K, V](background, purpose, *conf, indices, connectors.Postgres)
		if err != nil {
			return nil, fmt.Errorf("start cache: %w", err)
		}
		connectors.Postgres.Config.AutoPrune.StartAutoPrune(background, c, purpose)
		return c, nil
	}
	if conf.Connector == cache.ConnectorRedis && connectors.Redis != nil {
		db := connectors.Redis.Config.DBOffset + int(purpose)
		c := redis.NewCache[I, K, V](*conf, connectors.Redis, db, indices)
		return c, nil
	}

	return nil, fmt.Errorf("cache connector %q not enabled", conf.Connector)
}
