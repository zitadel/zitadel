package sql

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
)

type SQL struct {
	client                   *database.DB
	allowOrderByCreationDate bool
}

func (db *SQL) Health(ctx context.Context) error {
	return db.client.Ping()
}
