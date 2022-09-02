package sql

import (
	"context"
	"database/sql"
)

type SQL struct {
	client *sql.DB
}

func (db *SQL) Health(ctx context.Context) error {
	return db.client.Ping()
}
