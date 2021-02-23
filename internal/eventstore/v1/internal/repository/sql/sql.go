package sql

import (
	"context"
	"database/sql"

	//sql import
	_ "github.com/lib/pq"
)

type SQL struct {
	client *sql.DB
}

func (db *SQL) Health(ctx context.Context) error {
	return db.client.Ping()
}
