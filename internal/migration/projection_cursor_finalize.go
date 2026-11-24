package migration

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/eventstore"
)

// ProjectionCursorFinalizeSQL copies dual-written in_tx_order into filter_offset
// when cutting over to cursor resume (setup 77 Done).
const ProjectionCursorFinalizeSQL = `
UPDATE projections.current_states
SET filter_offset = in_tx_order
WHERE in_tx_order IS NOT NULL;
`

// ProjectionCursorFinalizeState mirrors the 79-style version/Finalized/ForceFinalize
// lastRun payload for the projection cursor cutover repeatable step.
type ProjectionCursorFinalizeState struct {
	Version       string `json:"version"`
	Finalized     bool   `json:"finalized"`
	ForceFinalize bool   `json:"-"`

	LastVersion   string `json:"-"`
	LastFinalized bool   `json:"-"`
}

// Check reports whether the repeatable step should run again.
func (s *ProjectionCursorFinalizeState) Check(lastRun map[string]interface{}) bool {
	if lastRun == nil {
		lastRun = map[string]interface{}{}
	}
	s.LastVersion, _ = lastRun["version"].(string)
	s.LastFinalized, _ = lastRun["finalized"].(bool)

	versionChanged := s.LastVersion != s.Version
	if s.LastVersion == "" {
		return true
	}
	if s.LastFinalized {
		return false
	}
	return versionChanged || s.ForceFinalize
}

// WillFinalize is true when this run should rewrite filter_offset and mark 77 Done.
func (s *ProjectionCursorFinalizeState) WillFinalize() bool {
	return s.ForceFinalize || (s.LastVersion != "" && s.LastVersion != s.Version)
}

// FinalizeAfterSuccess updates Finalized for the repeatable lastRun payload.
func (s *ProjectionCursorFinalizeState) FinalizeAfterSuccess() {
	s.Finalized = s.LastFinalized || s.WillFinalize()
}

type execContext interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// BackfillFilterOffsetFromInTxOrder rewrites current_states.filter_offset from in_tx_order.
func BackfillFilterOffsetFromInTxOrder(ctx context.Context, db execContext) error {
	_, err := db.ExecContext(ctx, ProjectionCursorFinalizeSQL)
	return err
}

// RecordEventstorePositionClockTimestampDone records historical setup 77 as Done.
func RecordEventstorePositionClockTimestampDone(ctx context.Context, es *eventstore.Eventstore) error {
	return Migrate(ctx, es, &eventstorePositionClockTimestampDone{})
}

type eventstorePositionClockTimestampDone struct{}

func (*eventstorePositionClockTimestampDone) Execute(context.Context, eventstore.Event) error {
	return nil
}

func (*eventstorePositionClockTimestampDone) String() string {
	return EventstorePositionClockTimestampStep
}
