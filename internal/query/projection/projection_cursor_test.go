package projection

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/migration"
)

func TestFilterOffsetIsCursorFromEventstore(t *testing.T) {
	t.Run("77 stored as done", func(t *testing.T) {
		es := newMockEventStore().appendFilterResponse([]eventstore.Event{
			&eventstore.BaseEvent{EventType: migration.DoneType},
		})
		got, err := filterOffsetIsCursorFromEventstore(t.Context(), es)
		require.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("77 not stored", func(t *testing.T) {
		es := newMockEventStore().appendFilterResponse(nil)
		got, err := filterOffsetIsCursorFromEventstore(t.Context(), es)
		require.NoError(t, err)
		assert.False(t, got)
	})
}

func TestApplyCustomConfigPreservesFilterOffsetIsCursor(t *testing.T) {
	cfg := handler.Config{FilterOffsetIsCursor: true}
	got := applyCustomConfig(cfg, CustomConfig{})
	assert.True(t, got.FilterOffsetIsCursor)
}
