package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

// Publish writes events to the eventstore using the provided pusher and database client.
func Publish(ctx context.Context, es *eventstore.Eventstore, client database.QueryExecutor, commands ...eventstore.Command) error {
	_, err := es.PushWithNewClient(ctx, client, commands...)
	return err
}
