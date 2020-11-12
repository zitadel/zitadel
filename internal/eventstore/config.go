package eventstore

import (
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/eventstore/internal/repository/sql"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_v2 "github.com/caos/zitadel/internal/eventstore/v2"
	sql_v2 "github.com/caos/zitadel/internal/eventstore/v2/repository/sql"
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
		esV2:             es_v2.NewEventstore(sql_v2.NewCRDB(sqlClient)),
	}, nil
}
