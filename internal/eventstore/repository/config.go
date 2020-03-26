package repository

import (
	"database/sql"

	"github.com/caos/eventstore-lib/pkg/repository"
)

type Config struct {
	Client *sql.DB
}

func Start(conf Config) repository.Repository {
	return &SQL{
		client: conf.Client,
	}
}
