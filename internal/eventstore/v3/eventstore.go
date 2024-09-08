package eventstore

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	// pushPlaceholderFmt defines how data are inserted into the events table
	pushPlaceholderFmt string
	// uniqueConstraintPlaceholderFmt defines the format of the unique constraint error returned from the database
	uniqueConstraintPlaceholderFmt string
)

type Eventstore struct {
	client *database.DB
	// used to send a pgnotify event on push to the postgres channel named after the event type
	// the channels can be used to send a trigger to the projection
	subscriptions *subscriptions
}

func NewEventstore(client *database.DB) *Eventstore {
	switch client.Type() {
	case "cockroach":
		pushPlaceholderFmt = "($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, hlc_to_timestamp(cluster_logical_timestamp()), cluster_logical_timestamp(), $%d)"
		uniqueConstraintPlaceholderFmt = "('%s', '%s', '%s')"
	case "postgres":
		pushPlaceholderFmt = "($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, statement_timestamp(), EXTRACT(EPOCH FROM transaction_timestamp()), $%d)"
		uniqueConstraintPlaceholderFmt = "(%s, %s, %s)"
	}
	return &Eventstore{client: client, subscriptions: newSubscriptions(client.Pool)}
}

func (es *Eventstore) Subscribe(queue chan<- decimal.Decimal, eventTypes ...eventstore.EventType) {
	es.subscriptions.Add(queue, eventTypes...)
}

func (es *Eventstore) Health(ctx context.Context) error {
	return es.client.PingContext(ctx)
}
