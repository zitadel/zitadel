package migration

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
)

const (
	startedType   = eventstore.EventType("system.migration.started")
	doneType      = eventstore.EventType("system.migration.done")
	failedType    = eventstore.EventType("system.migration.failed")
	aggregateType = eventstore.AggregateType("system")
	aggregateID   = "SYSTEM"
)

type Migration interface {
	String() string
	Execute(context.Context) error
}

func Migrate(ctx context.Context, es *eventstore.Eventstore, migration Migration) (err error) {
	if should, err := shouldExec(ctx, es, migration); !should || err != nil {
		return err
	}

	if _, err = es.Push(ctx, setupStartedCmd(migration)); err != nil {
		return err
	}

	err = migration.Execute(ctx)
	logging.OnError(err).Error("migration failed")

	_, err = es.Push(ctx, setupDoneCmd(migration, err))
	return err
}

func shouldExec(ctx context.Context, es *eventstore.Eventstore, migration Migration) (should bool, err error) {
	events, err := es.Filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderDesc().
		AddQuery().
		AggregateTypes(aggregateType).
		AggregateIDs(aggregateID).
		EventTypes(startedType, doneType, failedType).
		Builder())
	if err != nil {
		return false, err
	}

	var isStarted bool
	for _, event := range events {
		e, ok := event.(*SetupStep)
		if !ok {
			return false, errors.ThrowInternal(nil, "MIGRA-IJY3D", "Errors.Internal")
		}

		if e.Name != migration.String() {
			continue
		}

		switch event.Type() {
		case startedType, failedType:
			isStarted = !isStarted
		case doneType:
			return false, nil
		}
	}

	return !isStarted, nil
}
