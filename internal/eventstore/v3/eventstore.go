package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	// pushPlaceholderFmt defines how data are inserted into the events table
	pushPlaceholderFmt string
	// uniqueConstraintPlaceholderFmt defines the format of the unique constraint error returned from the database
	uniqueConstraintPlaceholderFmt string
)

type Eventstore struct {
	client *database.DB
}

func NewEventstore(client *database.DB) *Eventstore {
	switch client.Type() {
	case "cockroach":
		pushPlaceholderFmt = "($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, hlc_to_timestamp(cluster_logical_timestamp()), cluster_logical_timestamp(), $%d)"
		uniqueConstraintPlaceholderFmt = "('%s', '%s', '%s')"
	case "postgres":
		pushPlaceholderFmt = "($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp()), $%d)"
		uniqueConstraintPlaceholderFmt = "(%s, %s, %s)"
	}

	return &Eventstore{client: client}
}

func (es *Eventstore) Health(ctx context.Context) error {
	return es.client.PingContext(ctx)
}
