package migration

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectionCursorFinalizeSQL(t *testing.T) {
	assert.Contains(t, ProjectionCursorFinalizeSQL, "UPDATE projections.current_states")
	assert.Contains(t, ProjectionCursorFinalizeSQL, "SET filter_offset = in_tx_order")
	assert.Contains(t, ProjectionCursorFinalizeSQL, "WHERE in_tx_order IS NOT NULL")
	assert.NotContains(t, ProjectionCursorFinalizeSQL, "events2")
}

func TestProjectionCursorFinalizeState_CheckAndFinalize(t *testing.T) {
	tests := []struct {
		name          string
		forceFinalize bool
		lastRun       map[string]interface{}
		wantRun       bool
		wantFinalize  bool
		wantFinalized bool
	}{
		{
			name:          "missing lastRun runs without finalizing",
			lastRun:       nil,
			wantRun:       true,
			wantFinalize:  false,
			wantFinalized: false,
		},
		{
			name:          "empty lastRun runs without finalizing",
			lastRun:       map[string]interface{}{},
			wantRun:       true,
			wantFinalize:  false,
			wantFinalized: false,
		},
		{
			name:          "missing lastRun with ForceFinalize finalizes",
			forceFinalize: true,
			lastRun:       nil,
			wantRun:       true,
			wantFinalize:  true,
			wantFinalized: true,
		},
		{
			name:          "same version not finalized skips",
			lastRun:       map[string]interface{}{"version": "v2.0.0", "finalized": false},
			wantRun:       false,
			wantFinalize:  false,
			wantFinalized: false,
		},
		{
			name:          "same version missing finalized skips",
			lastRun:       map[string]interface{}{"version": "v2.0.0"},
			wantRun:       false,
			wantFinalize:  false,
			wantFinalized: false,
		},
		{
			name:          "same version not finalized with ForceFinalize runs and finalizes",
			forceFinalize: true,
			lastRun:       map[string]interface{}{"version": "v2.0.0", "finalized": false},
			wantRun:       true,
			wantFinalize:  true,
			wantFinalized: true,
		},
		{
			name:          "same version already finalized skips",
			lastRun:       map[string]interface{}{"version": "v2.0.0", "finalized": true},
			wantRun:       false,
			wantFinalize:  false,
			wantFinalized: true,
		},
		{
			name:          "same version already finalized with ForceFinalize skips",
			forceFinalize: true,
			lastRun:       map[string]interface{}{"version": "v2.0.0", "finalized": true},
			wantRun:       false,
			wantFinalize:  false,
			wantFinalized: true,
		},
		{
			name:          "version changed runs and finalizes",
			lastRun:       map[string]interface{}{"version": "v1.0.0", "finalized": false},
			wantRun:       true,
			wantFinalize:  true,
			wantFinalized: true,
		},
		{
			name:          "version changed after finalize skips",
			lastRun:       map[string]interface{}{"version": "v1.0.0", "finalized": true},
			wantRun:       false,
			wantFinalize:  false,
			wantFinalized: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ProjectionCursorFinalizeState{Version: "v2.0.0", ForceFinalize: tt.forceFinalize}
			assert.Equal(t, tt.wantRun, s.Check(tt.lastRun))
			assert.False(t, s.Finalized, "Check must not set finalized before Execute")
			if tt.wantRun {
				assert.Equal(t, tt.wantFinalize, s.WillFinalize())
			}
			s.FinalizeAfterSuccess()
			assert.Equal(t, tt.wantFinalized, s.Finalized)
		})
	}
}

func TestProjectionCursorFinalizeState_SequentialVersionChanges(t *testing.T) {
	s := &ProjectionCursorFinalizeState{Version: "v1.0.0"}
	require.True(t, s.Check(nil))
	assert.False(t, s.WillFinalize())
	s.FinalizeAfterSuccess()
	assert.False(t, s.Finalized)

	lastRun := map[string]interface{}{"version": s.Version, "finalized": s.Finalized}
	s = &ProjectionCursorFinalizeState{Version: "v2.0.0"}
	require.True(t, s.Check(lastRun))
	assert.True(t, s.WillFinalize())
	s.FinalizeAfterSuccess()
	assert.True(t, s.Finalized)

	lastRun = map[string]interface{}{"version": s.Version, "finalized": s.Finalized}
	s = &ProjectionCursorFinalizeState{Version: "v3.0.0"}
	assert.False(t, s.Check(lastRun), "later version must skip after previous run finalized")
}

func TestProjectionCursorFinalizeState_ForceFinalizeRestart(t *testing.T) {
	s := &ProjectionCursorFinalizeState{Version: "v1.0.0"}
	require.True(t, s.Check(nil))
	s.FinalizeAfterSuccess()
	require.False(t, s.Finalized)

	lastRun := map[string]interface{}{"version": "v1.0.0", "finalized": false}
	s = &ProjectionCursorFinalizeState{Version: "v1.0.0", ForceFinalize: true}
	require.True(t, s.Check(lastRun), "restart with ForceFinalize should re-run")
	assert.False(t, s.Finalized, "Check must not set finalized before Execute")
	assert.True(t, s.WillFinalize())
	s.FinalizeAfterSuccess()
	assert.True(t, s.Finalized)

	lastRun = map[string]interface{}{"version": "v1.0.0", "finalized": true}
	s = &ProjectionCursorFinalizeState{Version: "v1.0.0", ForceFinalize: true}
	assert.False(t, s.Check(lastRun), "later restart with ForceFinalize should skip")
}

func TestProjectionCursorFinalizeStateJSONOmitsForceFinalize(t *testing.T) {
	s := &ProjectionCursorFinalizeState{
		Version:       "v2.0.0",
		Finalized:     true,
		ForceFinalize: true,
		LastVersion:   "v1.0.0",
		LastFinalized: true,
	}
	data, err := json.Marshal(s)
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(data, &payload))
	assert.Equal(t, "v2.0.0", payload["version"])
	assert.Equal(t, true, payload["finalized"])
	assert.NotContains(t, payload, "ForceFinalize")
	assert.NotContains(t, payload, "forceFinalize")
	assert.NotContains(t, payload, "LastVersion")
	assert.NotContains(t, payload, "lastVersion")
	assert.NotContains(t, payload, "LastFinalized")
	assert.NotContains(t, payload, "lastFinalized")
	assert.False(t, strings.Contains(string(data), "v1.0.0"))
}

func TestBackfillFilterOffsetFromInTxOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(ProjectionCursorFinalizeSQL)).
		WillReturnResult(sqlmock.NewResult(0, 3))

	assert.NoError(t, BackfillFilterOffsetFromInTxOrder(context.Background(), db))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBackfillFilterOffsetFromInTxOrder_error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(ProjectionCursorFinalizeSQL)).
		WillReturnError(errors.New("boom"))

	assert.EqualError(t, BackfillFilterOffsetFromInTxOrder(context.Background(), db), "boom")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEventstorePositionClockTimestampDone(t *testing.T) {
	step := &eventstorePositionClockTimestampDone{}
	assert.Equal(t, EventstorePositionClockTimestampStep, step.String())
	assert.NoError(t, step.Execute(context.Background(), nil))
}
