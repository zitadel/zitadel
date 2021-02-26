package v1

import (
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/eventstore/v1/internal/repository/sql"
)

type Config struct {
	Repository  sql.Config
	ServiceName string
	Cache       *config.CacheConfig
}

func Start(conf Config) (Eventstore, error) {
	repo, _, err := sql.Start(conf.Repository)
	if err != nil {
		return nil, err
	}

	return &eventstore{
		repo: repo,
	}, nil
}
