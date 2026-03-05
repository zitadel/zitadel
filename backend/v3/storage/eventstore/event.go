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
// All commands are wrapped with the V3Command marker to indicate they were written through the v3 storage adapter.
func Publish(ctx context.Context, es LegacyEventstore, client database.QueryExecutor, commands ...eventstore.Command) error {
	v3Commands := make([]eventstore.Command, len(commands))
	for i, cmd := range commands {
		v3Commands[i] = &v3CommandWrapper{Command: cmd}
	}
	_, err := es.PushWithNewClient(ctx, client, v3Commands...)
	return err
}

// v3CommandWrapper wraps a Command to implement the V3Command marker interface.
type v3CommandWrapper struct {
	eventstore.Command
}

// IsV3Command implements eventstore.V3Command marker interface.
func (*v3CommandWrapper) IsV3Command() {}
