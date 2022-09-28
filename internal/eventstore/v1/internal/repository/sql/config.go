package sql

import (
	"database/sql"
)

func Start(client *sql.DB) *SQL {
	return &SQL{
		client: client,
	}
}
