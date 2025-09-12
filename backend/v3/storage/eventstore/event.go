package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type LegacyEventstore interface {
	PushWithNewClient(ctx context.Context, client database.QueryExecutor, commands ...eventstore.Command) ([]eventstore.Event, error)
}

// Publish writes events to the eventstore using the provided pusher and database client.
func Publish(ctx context.Context, es LegacyEventstore, client database.QueryExecutor, commands ...eventstore.Command) error {
	_, err := es.PushWithNewClient(ctx, client, commands...)
	return err
}
