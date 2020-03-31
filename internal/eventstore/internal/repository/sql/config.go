package sql

import (
	"database/sql"
)

type Config struct {
	Client *sql.DB
}

func Start(conf Config) *SQL {
	return &SQL{
		client: conf.Client,
	}
}
