package projection

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/migration"
)

func setupStepEvent(typ eventstore.EventType, name string) *migration.SetupStep {
	agg := eventstore.NewAggregate(context.Background(), migration.SystemAggregateID, migration.SystemAggregate, "v1")
	return &migration.SetupStep{
		BaseEvent: eventstore.BaseEvent{EventType: typ, Agg: agg},
		Name:      name,
	}
}

func TestFilterOffsetIsCursorFromEventstore(t *testing.T) {
	name := migration.EventstorePositionClockTimestampStep

	t.Run("77 stored as done", func(t *testing.T) {
		es := newMockEventStore().appendFilterResponse([]eventstore.Event{
			setupStepEvent(migration.DoneType, name),
		})
		got, err := FilterOffsetIsCursorFromEventstore(t.Context(), es)
		require.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("77 not stored", func(t *testing.T) {
		es := newMockEventStore().appendFilterResponse(nil)
		got, err := FilterOffsetIsCursorFromEventstore(t.Context(), es)
		require.NoError(t, err)
		assert.False(t, got)
	})

	t.Run("77 started without done errors", func(t *testing.T) {
		es := newMockEventStore().appendFilterResponse([]eventstore.Event{
			setupStepEvent(migration.StartedType, name),
		})
		got, err := FilterOffsetIsCursorFromEventstore(t.Context(), es)
		require.Error(t, err)
		assert.False(t, got)
	})
}

func TestApplyCustomConfigPreservesFilterOffsetIsCursor(t *testing.T) {
	cfg := handler.Config{FilterOffsetIsCursor: true}
	got := applyCustomConfig(cfg, CustomConfig{})
	assert.True(t, got.FilterOffsetIsCursor)
}
