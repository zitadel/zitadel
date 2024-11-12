package pg

import (
	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/database"
)

type Config struct {
	Enabled   bool
	AutoPrune cache.AutoPruneConfig
}

type Connector struct {
	PGXPool
	Dialect string
	Config  Config
}

func NewConnector(config Config, client *database.DB) *Connector {
	if !config.Enabled {
		return nil
	}
	return &Connector{
		PGXPool: client.Pool,
		Dialect: client.Type(),
		Config:  config,
	}
}
