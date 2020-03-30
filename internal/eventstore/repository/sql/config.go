package sql

import (
	"database/sql"

	"github.com/caos/zitadel/internal/eventstore/repository"
)

type Config struct {
	Client *sql.DB
}

func Start(conf Config) repository.Repository {
	return &SQL{
		client: conf.Client,
	}
}
