package setup

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
)

// ProjectionCursorFinalize is a repeatable step that cuts projections over to
// cursor resume (filter_offset holds events2.in_tx_order) after pre-#12703
// projectors are gone. First run records Version without finalizing so mixed
// deployments stay on OFFSET. ForceFinalize or a later version bump rewrites
// filter_offset and records historical setup 77 as Done.
type ProjectionCursorFinalize struct {
	dbClient   *database.DB
	eventstore *eventstore.Eventstore

	migration.ProjectionCursorFinalizeState

	// record77Done overrides RecordEventstorePositionClockTimestampDone in tests.
	record77Done func(context.Context) error
}

func (mig *ProjectionCursorFinalize) Check(lastRun map[string]interface{}) bool {
	return mig.ProjectionCursorFinalizeState.Check(lastRun)
}

func (mig *ProjectionCursorFinalize) Execute(ctx context.Context, _ eventstore.Event) error {
	if mig.WillFinalize() {
		if err := migration.BackfillFilterOffsetFromInTxOrder(ctx, mig.dbClient); err != nil {
			return err
		}
		if err := mig.doRecord77Done(ctx); err != nil {
			return err
		}
	}
	mig.FinalizeAfterSuccess()
	return nil
}

func (mig *ProjectionCursorFinalize) doRecord77Done(ctx context.Context) error {
	if mig.record77Done != nil {
		return mig.record77Done(ctx)
	}
	return migration.RecordEventstorePositionClockTimestampDone(ctx, mig.eventstore)
}

func (mig *ProjectionCursorFinalize) String() string {
	return "repeatable_projection_cursor_finalize"
}
