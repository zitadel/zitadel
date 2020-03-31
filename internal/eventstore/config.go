package eventstore

import (
	"github.com/caos/zitadel/internal/eventstore/internal/repository/sql"
)

type Config struct {
	Repository sql.Config
}

func Start(conf Config) App {
	return &app{
		repo: sql.Start(conf.Repository),
	}
}
