package migration

import (
	"context"
	"encoding/json"

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

	if len(events) == 0 {
		return true, nil
	}

	if events[len(events)-1].Type() == startedType {
		return false, nil
	}

	for _, e := range events {
		step := new(SetupStep)

		err = json.Unmarshal(e.DataAsBytes(), step)
		if err != nil {
			return false, err
		}

		if step.Name != migration.String() {
			continue
		}

		switch e.Type() {
		case startedType, doneType:
			//TODO: if started should we wait until done/failed?
			return false, nil
		case failedType:
			//TODO: how to allow retries?
			logging.WithFields("migration", migration.String()).Error("failed before")
			return false, errors.ThrowInternal(nil, "MIGRA-mjI2E", "migration failed before")
		}
	}
	return true, nil
}
