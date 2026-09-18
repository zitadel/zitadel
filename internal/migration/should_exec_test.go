package migration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/repository/mock"
)

type checkRecordingMigration struct {
	name       string
	checkCalls int
	lastRun    map[string]any
}

func (m *checkRecordingMigration) String() string { return m.name }
func (m *checkRecordingMigration) Execute(context.Context, eventstore.Event) error {
	return nil
}
func (m *checkRecordingMigration) Check(lastRun map[string]any) bool {
	m.checkCalls++
	m.lastRun = lastRun
	return false
}

func TestShouldExec_failedRepeatableCallsCheckAndRetries(t *testing.T) {
	mig := &checkRecordingMigration{name: "repeatable-step"}
	done := migrationSetupEvent(t, mig.name, repeatableDoneType, map[string]any{"version": "v1.0.0", "finalized": false})
	failed := migrationSetupEvent(t, mig.name, failedType, map[string]any{"version": "v2.0.0", "finalized": true})

	repo := mock.NewRepo(t)
	repo.ExpectFilterEvents(done, failed)
	es := eventstore.NewEventstore(&eventstore.Config{
		Querier: repo.MockQuerier,
		Pusher:  repo.MockPusher,
	})

	should, err := shouldExec(t.Context(), es, mig)
	require.NoError(t, err)
	assert.True(t, should)
	assert.Equal(t, 1, mig.checkCalls)
	require.NotNil(t, mig.lastRun)
	assert.Equal(t, "v1.0.0", mig.lastRun["version"])
	assert.Equal(t, false, mig.lastRun["finalized"])
}

func migrationSetupEvent(t *testing.T, name string, typ eventstore.EventType, lastRun map[string]any) *repository.Event {
	t.Helper()
	cmd := &SetupStep{
		BaseEvent: *eventstore.NewBaseEventForPush(
			t.Context(),
			eventstore.NewAggregate(t.Context(), SystemAggregateID, SystemAggregate, "v1"),
			typ,
		),
		Name:    name,
		LastRun: lastRun,
	}
	data, err := eventstore.EventData(cmd)
	require.NoError(t, err)
	return &repository.Event{
		InstanceID:    cmd.Aggregate().InstanceID,
		Typ:           cmd.Type(),
		Data:          data,
		EditorUser:    cmd.Creator(),
		Version:       cmd.Aggregate().Version,
		AggregateID:   cmd.Aggregate().ID,
		AggregateType: cmd.Aggregate().Type,
	}
}
