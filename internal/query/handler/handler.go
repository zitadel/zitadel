package handler

import (
	"context"
	"database/sql"

	"github.com/caos/zitadel/internal/eventstore"
)

func Start(ctx context.Context, es *eventstore.Eventstore, client *sql.DB) {
	NewOrgHandler(ctx, es, client)
}
