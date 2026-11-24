package setup

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/migration"
)

var _ migration.RepeatableMigration = (*ProjectionCursorFinalize)(nil)

func TestProjectionCursorFinalize_String(t *testing.T) {
	assert.Equal(t, "repeatable_projection_cursor_finalize", new(ProjectionCursorFinalize).String())
	assert.NotEqual(t, migration.EventstorePositionClockTimestampStep, new(ProjectionCursorFinalize).String())
}

func TestProjectionCursorFinalize_ExecuteSkipsBackfillWhenNotFinalizing(t *testing.T) {
	mig := &ProjectionCursorFinalize{}
	mig.Version = "v1.0.0"
	require.True(t, mig.Check(nil))
	assert.False(t, mig.WillFinalize())

	assert.NoError(t, mig.Execute(context.Background(), nil))
	assert.False(t, mig.Finalized)
}

func TestProjectionCursorFinalize_ExecuteFinalizes(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(migration.ProjectionCursorFinalizeSQL)).
		WillReturnResult(sqlmock.NewResult(0, 2))

	recorded := false
	mig := &ProjectionCursorFinalize{
		dbClient: &database.DB{DB: db},
		record77Done: func(context.Context) error {
			recorded = true
			return nil
		},
	}
	mig.Version = "v1.0.0"
	mig.ForceFinalize = true
	require.True(t, mig.Check(nil))
	require.True(t, mig.WillFinalize())

	assert.NoError(t, mig.Execute(context.Background(), nil))
	assert.True(t, recorded)
	assert.True(t, mig.Finalized)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectionCursorFinalize_ExecuteBackfillError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(migration.ProjectionCursorFinalizeSQL)).
		WillReturnError(errors.New("boom"))

	mig := &ProjectionCursorFinalize{
		dbClient: &database.DB{DB: db},
	}
	mig.Version = "v1.0.0"
	mig.ForceFinalize = true
	require.True(t, mig.Check(nil))

	assert.EqualError(t, mig.Execute(context.Background(), nil), "boom")
	assert.False(t, mig.Finalized)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectionCursorFinalizeJSONOmitsForceFinalize(t *testing.T) {
	mig := &ProjectionCursorFinalize{}
	mig.Version = "v2.0.0"
	mig.Finalized = true
	mig.ForceFinalize = true
	mig.LastVersion = "v1.0.0"
	mig.LastFinalized = true

	data, err := json.Marshal(mig)
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(data, &payload))
	assert.Equal(t, "v2.0.0", payload["version"])
	assert.Equal(t, true, payload["finalized"])
	assert.NotContains(t, payload, "ForceFinalize")
	assert.NotContains(t, payload, "forceFinalize")
}
