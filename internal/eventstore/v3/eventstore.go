package eventstore

import (
	"context"

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
	subscribedEventTypes map[eventstore.EventType][]chan<- float64
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

	return &Eventstore{client: client, subscribedEventTypes: make(map[eventstore.EventType][]chan<- float64)}
}

func (es *Eventstore) Subscribe(queue chan<- float64, eventTypes ...eventstore.EventType) {
	for _, eventType := range eventTypes {
		es.subscribedEventTypes[eventType] = append(es.subscribedEventTypes[eventType], queue)
	}
}

func (es *Eventstore) Health(ctx context.Context) error {
	return es.client.PingContext(ctx)
}
