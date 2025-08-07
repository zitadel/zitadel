package mirror

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/readmodel"
	"github.com/zitadel/zitadel/internal/v2/system"
	mirror_event "github.com/zitadel/zitadel/internal/v2/system/mirror"
)

func queryLastSuccessfulMigration(ctx context.Context, destinationES *eventstore.EventStore, source string) (*readmodel.LastSuccessfulMirror, error) {
	lastSuccess := readmodel.NewLastSuccessfulMirror(source)
	if shouldIgnorePrevious {
		return lastSuccess, nil
	}
	_, err := destinationES.Query(
		ctx,
		eventstore.NewQuery(
			system.AggregateInstance,
			lastSuccess,
			eventstore.SetFilters(lastSuccess.Filter()),
		),
	)
	if err != nil {
		return nil, err
	}

	return lastSuccess, nil
}

func writeMigrationSucceeded(ctx context.Context, destinationES *eventstore.EventStore, id, source string, position decimal.Decimal) error {
	return destinationES.Push(
		ctx,
		eventstore.NewPushIntent(
			system.AggregateInstance,
			eventstore.AppendAggregate(
				system.AggregateOwner,
				system.AggregateType,
				id,
				eventstore.CurrentSequenceMatches(0),
				eventstore.AppendCommands(mirror_event.NewSucceededCommand(source, position)),
			),
		),
	)
}

func writeMigrationFailed(ctx context.Context, destinationES *eventstore.EventStore, id, source string, err error) error {
	return destinationES.Push(
		ctx,
		eventstore.NewPushIntent(
			system.AggregateInstance,
			eventstore.AppendAggregate(
				system.AggregateOwner,
				system.AggregateType,
				id,
				eventstore.CurrentSequenceMatches(0),
				eventstore.AppendCommands(mirror_event.NewFailedCommand(source, err)),
			),
		),
	)
}
