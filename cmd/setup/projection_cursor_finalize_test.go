package setup

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/migration"
)

var _ migration.RepeatableMigration = (*ProjectionCursorFinalize)(nil)

func TestProjectionCursorFinalize_String(t *testing.T) {
	assert.Equal(t, "repeatable_projection_cursor_finalize", new(ProjectionCursorFinalize).String())
	assert.NotEqual(t, migration.EventstorePositionClockTimestampStep, new(ProjectionCursorFinalize).String())
	assert.Equal(t, migration.EventstorePositionClockTimestampStep, new(eventstorePositionClockTimestampDone).String())
}

func TestProjectionCursorFinalizeSQL(t *testing.T) {
	assert.Contains(t, projectionCursorFinalizeSQL, "UPDATE projections.current_states")
	assert.Contains(t, projectionCursorFinalizeSQL, "SET filter_offset = in_tx_order")
	assert.Contains(t, projectionCursorFinalizeSQL, "WHERE in_tx_order IS NOT NULL")
	assert.NotContains(t, projectionCursorFinalizeSQL, "events2")
}

func TestProjectionCursorFinalize_CheckAndFinalize(t *testing.T) {
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
			mig := &ProjectionCursorFinalize{Version: "v2.0.0", ForceFinalize: tt.forceFinalize}
			assert.Equal(t, tt.wantRun, mig.Check(tt.lastRun))
			assert.False(t, mig.Finalized, "Check must not set finalized before Execute")
			if tt.wantRun {
				assert.Equal(t, tt.wantFinalize, mig.willFinalize())
			}
			mig.finalizeAfterSuccess()
			assert.Equal(t, tt.wantFinalized, mig.Finalized)
		})
	}
}

func TestProjectionCursorFinalize_SequentialVersionChanges(t *testing.T) {
	mig := &ProjectionCursorFinalize{Version: "v1.0.0"}
	require.True(t, mig.Check(nil))
	assert.False(t, mig.willFinalize())
	mig.finalizeAfterSuccess()
	assert.False(t, mig.Finalized)

	lastRun := map[string]interface{}{"version": mig.Version, "finalized": mig.Finalized}
	mig = &ProjectionCursorFinalize{Version: "v2.0.0"}
	require.True(t, mig.Check(lastRun))
	assert.True(t, mig.willFinalize())
	mig.finalizeAfterSuccess()
	assert.True(t, mig.Finalized)

	lastRun = map[string]interface{}{"version": mig.Version, "finalized": mig.Finalized}
	mig = &ProjectionCursorFinalize{Version: "v3.0.0"}
	assert.False(t, mig.Check(lastRun), "later version must skip after previous run finalized")
}

func TestProjectionCursorFinalize_ForceFinalizeRestart(t *testing.T) {
	mig := &ProjectionCursorFinalize{Version: "v1.0.0"}
	require.True(t, mig.Check(nil))
	mig.finalizeAfterSuccess()
	require.False(t, mig.Finalized)

	lastRun := map[string]interface{}{"version": "v1.0.0", "finalized": false}
	mig = &ProjectionCursorFinalize{Version: "v1.0.0", ForceFinalize: true}
	require.True(t, mig.Check(lastRun), "restart with ForceFinalize should re-run")
	assert.False(t, mig.Finalized, "Check must not set finalized before Execute")
	assert.True(t, mig.willFinalize())
	mig.finalizeAfterSuccess()
	assert.True(t, mig.Finalized)

	lastRun = map[string]interface{}{"version": "v1.0.0", "finalized": true}
	mig = &ProjectionCursorFinalize{Version: "v1.0.0", ForceFinalize: true}
	assert.False(t, mig.Check(lastRun), "later restart with ForceFinalize should skip")
}

func TestProjectionCursorFinalizeJSONOmitsForceFinalize(t *testing.T) {
	mig := &ProjectionCursorFinalize{
		Version:       "v2.0.0",
		Finalized:     true,
		ForceFinalize: true,
		lastVersion:   "v1.0.0",
		lastFinalized: true,
	}
	data, err := json.Marshal(mig)
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(data, &payload))
	assert.Equal(t, "v2.0.0", payload["version"])
	assert.Equal(t, true, payload["finalized"])
	assert.NotContains(t, payload, "ForceFinalize")
	assert.NotContains(t, payload, "forceFinalize")
	assert.NotContains(t, payload, "lastVersion")
	assert.NotContains(t, payload, "lastFinalized")
	assert.False(t, strings.Contains(string(data), "v1.0.0"))
}
