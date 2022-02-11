package sql

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Start(client *sql.DB) *SQL {
	return &SQL{
		client: client,
	}
}
