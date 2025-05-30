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
	DoneType           = eventstore.EventType("system.migration.done")
	failedType         = eventstore.EventType("system.migration.failed")
	repeatableDoneType = eventstore.EventType("system.migration.repeatable.done")
	SystemAggregate    = eventstore.AggregateType("system")
	SystemAggregateID  = "SYSTEM"
)

var (
	errMigrationAlreadyStarted = errors.New("already started")
)

type Migration interface {
	String() string
	Execute(ctx context.Context, startedEvent eventstore.Event) error
}

type errCheckerMigration interface {
	Migration
	ContinueOnErr(err error) bool
}

type RepeatableMigration interface {
	Migration

	// Check if the migration should be executed again.
	// True will repeat the migration, false will not.
	Check(lastRun map[string]any) bool
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

	should, err := checkExec(ctx, es, migration)
	if err != nil && !continueOnErr(err) {
		return err
	}
	if !should {
		return nil
	}

	startedEvent, err := es.Push(ctx, setupStartedCmd(ctx, migration))
	if err != nil && !continueOnErr(err) {
		return err
	}

	logging.WithFields("name", migration.String()).Info("starting migration")
	err = migration.Execute(ctx, startedEvent[0])
	logging.WithFields("name", migration.String()).OnError(err).Error("migration failed")

	_, pushErr := es.Push(ctx, setupDoneCmd(ctx, migration, err))
	logging.WithFields("name", migration.String()).OnError(pushErr).Error("migration finish failed")
	if err != nil {
		return err
	}
	return pushErr
}

func LastStuckStep(ctx context.Context, es *eventstore.Eventstore) (*SetupStep, error) {
	var states StepStates
	err := es.FilterToQueryReducer(ctx, &states)
	if err != nil {
		return nil, err
	}
	step := states.lastByState(StepStarted)
	if step == nil {
		return nil, nil
	}

	return step.SetupStep, nil
}

var _ Migration = (*cancelMigration)(nil)

type cancelMigration struct {
	name string
}

// Execute implements Migration
func (*cancelMigration) Execute(context.Context, eventstore.Event) error {
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
	var states StepStates
	err = es.FilterToQueryReducer(ctx, &states)
	if err != nil {
		return false, err
	}
	step := states.byName(migration.String())
	if step == nil {
		return true, nil
	}
	if step.state == StepFailed {
		return true, nil
	}
	if step.state == StepStarted {
		return false, errMigrationAlreadyStarted
	}

	repeatable, ok := migration.(RepeatableMigration)
	if !ok {
		return step.state != StepDone, nil
	}
	lastRun, _ := step.LastRun.(map[string]interface{})
	return repeatable.Check(lastRun), nil
}
