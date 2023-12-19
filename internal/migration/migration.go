package migration

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	StartedType        = eventstore.EventType("system.migration.started")
	doneType           = eventstore.EventType("system.migration.done")
	failedType         = eventstore.EventType("system.migration.failed")
	repeatableDoneType = eventstore.EventType("system.migration.repeatable.done")
	aggregateType      = eventstore.AggregateType("system")
	aggregateID        = "SYSTEM"
)

var (
	errMigrationAlreadyStarted = errors.New("already started")
)

type Migration interface {
	String() string
	Execute(context.Context) error
}

type errCheckerMigration interface {
	Migration
	ContinueOnErr(err error) bool
}

type RepeatableMigration interface {
	Migration
	SetLastExecution(lastRun map[string]interface{})
	Check() bool
}

func Migrate(ctx context.Context, es *eventstore.Eventstore, migration Migration) (err error) {
	logging.WithFields("name", migration.String()).Info("verify migration")

	continueOnErr := func(err error) bool {
		return false
	}
	errChecker, ok := migration.(errCheckerMigration)
	if ok {
		continueOnErr = errChecker.ContinueOnErr
	}

	// if should, err := checkExec(ctx, es, migration); !should || err != nil {
	should, err := checkExec(ctx, es, migration)
	if err != nil && !continueOnErr(err) {
		return err
	}
	if !should {
		return nil
	}

	if _, err = es.Push(ctx, setupStartedCmd(ctx, migration)); err != nil && !continueOnErr(err) {
		return err
	}

	logging.WithFields("name", migration.String()).Info("starting migration")
	err = migration.Execute(ctx)
	logging.WithFields("name", migration.String()).OnError(err).Error("migration failed")

	_, pushErr := es.Push(ctx, setupDoneCmd(ctx, migration, err))
	logging.WithFields("name", migration.String()).OnError(pushErr).Error("migration finish failed")
	if err != nil {
		return err
	}
	return pushErr
}

func LatestStep(ctx context.Context, es *eventstore.Eventstore) (*SetupStep, error) {
	events, err := es.Filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderDesc().
		Limit(1).
		AddQuery().
		AggregateTypes(aggregateType).
		AggregateIDs(aggregateID).
		EventTypes(StartedType, doneType, repeatableDoneType, failedType).
		Builder())
	if err != nil {
		return nil, err
	}
	step, ok := events[0].(*SetupStep)
	if !ok {
		return nil, zerrors.ThrowInternal(nil, "MIGRA-hppLM", "setup step is malformed")
	}
	return step, nil
}

var _ Migration = (*cancelMigration)(nil)

type cancelMigration struct {
	name string
}

// Execute implements Migration
func (*cancelMigration) Execute(context.Context) error {
	return nil
}

// String implements Migration
func (m *cancelMigration) String() string {
	return m.name
}

var errCancelStep = zerrors.ThrowError(nil, "MIGRA-zo86K", "migration canceled manually")

func CancelStep(ctx context.Context, es *eventstore.Eventstore, step *SetupStep) error {
	_, err := es.Push(ctx, setupDoneCmd(ctx, &cancelMigration{name: step.Name}, errCancelStep))
	return err
}

// checkExec ensures that only one setup step is done concurrently
// if a setup step is already started, it calls shouldExec after some time again
func checkExec(ctx context.Context, es *eventstore.Eventstore, migration Migration) (bool, error) {
	timer := time.NewTimer(0)
	for {
		select {
		case <-ctx.Done():
			return false, zerrors.ThrowInternal(nil, "MIGR-as3f7", "Errors.Internal")
		case <-timer.C:
			should, err := shouldExec(ctx, es, migration)
			if err != nil {
				if !errors.Is(err, errMigrationAlreadyStarted) {
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
		InstanceID("").
		AddQuery().
		AggregateTypes(aggregateType).
		AggregateIDs(aggregateID).
		EventTypes(StartedType, doneType, repeatableDoneType, failedType).
		Builder())
	if err != nil {
		return false, err
	}

	var isStarted bool
	for _, event := range events {
		e, ok := event.(*SetupStep)
		if !ok {
			return false, zerrors.ThrowInternal(nil, "MIGRA-IJY3D", "Errors.Internal")
		}

		if e.Name != migration.String() {
			continue
		}

		switch event.Type() {
		case StartedType, failedType:
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
