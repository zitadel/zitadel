package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
)

type Eventstore struct {
	client *database.DB
}

func NewEventstore(client *database.DB) *Eventstore {
	return &Eventstore{client: client}
}

func (es *Eventstore) Health(ctx context.Context) error {
	return es.client.PingContext(ctx)
}
