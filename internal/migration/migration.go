package migration

import (
	"context"
	errs "errors"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	startedType        = eventstore.EventType("system.migration.started")
	doneType           = eventstore.EventType("system.migration.done")
	failedType         = eventstore.EventType("system.migration.failed")
	repeatableDoneType = eventstore.EventType("system.migration.repeatable.done")
	aggregateType      = eventstore.AggregateType("system")
	aggregateID        = "SYSTEM"
)

var (
	errMigrationAlreadyStarted = errs.New("already started")
)

type Migration interface {
	String() string
	Execute(context.Context) error
}

type RepeatableMigration interface {
	Migration
	SetLastExecution(lastRun map[string]interface{})
	Check() bool
}

func Migrate(ctx context.Context, es *eventstore.Eventstore, migration Migration) (err error) {
	logging.Infof("verify migration %s", migration.String())

	if should, err := checkExec(ctx, es, migration); !should || err != nil {
		return err
	}

	if _, err = es.Push(ctx, setupStartedCmd(migration)); err != nil {
		return err
	}

	logging.Infof("starting migration %s", migration.String())
	err = migration.Execute(ctx)
	logging.OnError(err).Error("migration failed")

	_, pushErr := es.Push(ctx, setupDoneCmd(migration, err))
	logging.OnError(pushErr).Error("migration failed")
	if err != nil {
		return err
	}
	return pushErr
}

// checkExec ensures that only one setup step is done concurrently
// if a setup step is already started, it calls shouldExec after some time again
func checkExec(ctx context.Context, es *eventstore.Eventstore, migration Migration) (bool, error) {
	timer := time.NewTimer(0)
	for {
		select {
		case <-ctx.Done():
			return false, errors.ThrowInternal(nil, "MIGR-as3f7", "Errors.Internal")
		case <-timer.C:
			should, err := shouldExec(ctx, es, migration)
			if err != nil {
				if !errs.Is(err, errMigrationAlreadyStarted) {
					return false, err
				}
				logging.WithFields("migration step", migration.String()).
					Warn("migration already started, will check again in 5 seconds")
				timer.Reset(5 * time.Second)
				break
			}
			return should, nil
		}
	}
}

func shouldExec(ctx context.Context, es *eventstore.Eventstore, migration Migration) (should bool, err error) {
	events, err := es.Filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderAsc().
		AddQuery().
		AggregateTypes(aggregateType).
		AggregateIDs(aggregateID).
		EventTypes(startedType, doneType, repeatableDoneType, failedType).
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
		case doneType,
			repeatableDoneType:
			repeatable, ok := migration.(RepeatableMigration)
			if !ok {
				return false, nil
			}
			isStarted = false
			repeatable.SetLastExecution(e.LastRun.(map[string]interface{}))
		}
	}

	if isStarted {
		return false, errMigrationAlreadyStarted
	}
	repeatable, ok := migration.(RepeatableMigration)
	if !ok {
		return true, nil
	}
	return repeatable.Check(), nil
}
