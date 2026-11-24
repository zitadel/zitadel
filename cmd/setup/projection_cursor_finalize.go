package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
)

var (
	//go:embed projection_cursor_finalize.sql
	projectionCursorFinalizeSQL string
)

// ProjectionCursorFinalize is a repeatable step that cuts projections over to
// cursor resume (filter_offset holds events2.in_tx_order) after pre-#12703
// projectors are gone. First run records Version without finalizing so mixed
// deployments stay on OFFSET. ForceFinalize or a later version bump rewrites
// filter_offset and records historical setup 77 as Done.
type ProjectionCursorFinalize struct {
	dbClient   *database.DB
	eventstore *eventstore.Eventstore

	Version       string `json:"version"`
	Finalized     bool   `json:"finalized"`
	ForceFinalize bool   `json:"-"` // YAML/env only; not persisted on lastRun

	lastVersion   string `json:"-"`
	lastFinalized bool   `json:"-"`
}

func (mig *ProjectionCursorFinalize) Check(lastRun map[string]interface{}) bool {
	if lastRun == nil {
		lastRun = map[string]interface{}{}
	}
	mig.lastVersion, _ = lastRun["version"].(string)
	mig.lastFinalized, _ = lastRun["finalized"].(bool)

	versionChanged := mig.lastVersion != mig.Version
	if mig.lastVersion == "" {
		return true // first run
	}
	if mig.lastFinalized {
		return false // never re-run after finalize, even on a new binary
	}
	return versionChanged || mig.ForceFinalize
}

func (mig *ProjectionCursorFinalize) willFinalize() bool {
	return mig.ForceFinalize || (mig.lastVersion != "" && mig.lastVersion != mig.Version)
}

func (mig *ProjectionCursorFinalize) finalizeAfterSuccess() {
	mig.Finalized = mig.lastFinalized || mig.willFinalize()
}

func (mig *ProjectionCursorFinalize) Execute(ctx context.Context, _ eventstore.Event) error {
	if mig.willFinalize() {
		if _, err := mig.dbClient.ExecContext(ctx, projectionCursorFinalizeSQL); err != nil {
			return err
		}
		if err := migration.Migrate(ctx, mig.eventstore, &eventstorePositionClockTimestampDone{}); err != nil {
			return err
		}
	}
	mig.finalizeAfterSuccess()
	return nil
}

func (mig *ProjectionCursorFinalize) String() string {
	return "repeatable_projection_cursor_finalize"
}

// eventstorePositionClockTimestampDone records historical setup 77 as Done so
// FilterOffsetIsCursorFromEventstore flips projectors to sort-key resume.
type eventstorePositionClockTimestampDone struct{}

func (*eventstorePositionClockTimestampDone) Execute(context.Context, eventstore.Event) error {
	return nil
}

func (*eventstorePositionClockTimestampDone) String() string {
	return migration.EventstorePositionClockTimestampStep
}
