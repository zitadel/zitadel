package migrate

import (
	"context"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/projection"
	"github.com/zitadel/zitadel/internal/v2/readmodel"
	"github.com/zitadel/zitadel/internal/v2/system"
	"github.com/zitadel/zitadel/internal/v2/system/mirror"
)

func queryLastSuccessfulMigration(ctx context.Context, es *eventstore.EventStore, destination string) (*readmodel.LastSuccessfulMirror, error) {
	lastSuccess := readmodel.NewLastSuccessfulMirror(destination)
	if shouldIgnorePrevious {
		return lastSuccess, nil
	}
	_, err := es.Query(
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

func writeMigrationStart(ctx context.Context, es *eventstore.EventStore, id string, destination string) (float64, error) {
	var cmd eventstore.Command
	if len(instanceIDs) > 0 {
		cmd = mirror.NewStartedInstancesCommand(destination, instanceIDs)
	} else {
		cmd = mirror.NewStartedSystemCommand(destination)
	}

	var position projection.HighestPosition

	err := es.Push(
		ctx,
		eventstore.NewPushIntent(
			system.AggregateInstance,
			eventstore.AppendAggregate(
				system.AggregateOwner,
				system.AggregateType,
				id,
				// TODO: what sequence to use
				eventstore.IgnoreCurrentSequence(),
				eventstore.AppendCommands(cmd),
			),
			eventstore.PushReducer(&position),
		),
	)
	if err != nil {
		return 0, err
	}
	return position.Position, nil
}

func writeMigrationSucceeded(ctx context.Context, es *eventstore.EventStore, id string) error {
	return es.Push(
		ctx,
		eventstore.NewPushIntent(
			system.AggregateInstance,
			eventstore.AppendAggregate(
				system.AggregateOwner,
				system.AggregateType,
				id,
				// TODO: what sequence to use
				eventstore.CurrentSequenceAtLeast(1),
				eventstore.AppendCommands(mirror.NewSucceededCommand()),
			),
		),
	)
}

func writeMigrationFailed(ctx context.Context, es *eventstore.EventStore, id string, err error) error {
	return es.Push(
		ctx,
		eventstore.NewPushIntent(
			system.AggregateInstance,
			eventstore.AppendAggregate(
				system.AggregateOwner,
				system.AggregateType,
				id,
				// TODO: what sequence to use
				eventstore.CurrentSequenceAtLeast(1),
				eventstore.AppendCommands(mirror.NewFailedCommand(err)),
			),
		),
	)
}
