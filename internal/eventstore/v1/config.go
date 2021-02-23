package v1

import (
	"github.com/caos/zitadel/internal/cache/config"
	eventstore2 "github.com/caos/zitadel/internal/eventstore"
	sql_v2 "github.com/caos/zitadel/internal/eventstore/repository/sql"
	"github.com/caos/zitadel/internal/eventstore/v1/internal/repository/sql"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Config struct {
	Repository  sql.Config
	ServiceName string
	Cache       *config.CacheConfig
}

func Start(conf Config) (Eventstore, error) {
	repo, sqlClient, err := sql.Start(conf.Repository)
	if err != nil {
		return nil, err
	}

	return &eventstore{
		repo:             repo,
		aggregateCreator: models.NewAggregateCreator(conf.ServiceName),
		esV2:             eventstore2.NewEventstore(sql_v2.NewCRDB(sqlClient)),
	}, nil
}
