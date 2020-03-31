package eventstore

import (
	"github.com/caos/zitadel/internal/eventstore/internal/repository/sql"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type Config struct {
	Repository  sql.Config
	ServiceName string
}

func Start(conf Config) App {
	return &app{
		repo:             sql.Start(conf.Repository),
		aggregateCreator: models.NewAggregateCreator(conf.ServiceName),
	}
}
