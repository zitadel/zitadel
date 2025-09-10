package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type Event struct {
	AggregateType string `json:"aggregateType"`
	AggregateID   string `json:"aggregateId"`
	Type          string `json:"type"`
	Payload       any    `json:"payload,omitempty"`
}

func Publish(ctx context.Context, events []*Event, db database.Executor) error {
	for _, event := range events {
		_, err := db.Exec(ctx, `INSERT INTO events (aggregate_type, aggregate_id) VALUES ($1, $2)`, event.AggregateType, event.AggregateID)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteLegacyEvents(ctx context.Context, es eventstore.Pusher, client database.QueryExecutor, commands ...eventstore.Command) error {
	_, err := es.Push(ctx, client, commands...)
	return err
}
