package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/migration"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestIsUniqueConstraintOwnerDeleteReady(t *testing.T) {
	t.Run("nil checker is not ready", func(t *testing.T) {
		c := &Commands{}
		ready, err := c.isUniqueConstraintOwnerDeleteReady(t.Context())
		require.NoError(t, err)
		assert.False(t, ready)
	})

	t.Run("ready is cached", func(t *testing.T) {
		calls := 0
		c := &Commands{
			ownerDeleteReady: func(context.Context) (bool, error) {
				calls++
				return true, nil
			},
		}
		ready, err := c.isUniqueConstraintOwnerDeleteReady(t.Context())
		require.NoError(t, err)
		assert.True(t, ready)
		ready, err = c.isUniqueConstraintOwnerDeleteReady(t.Context())
		require.NoError(t, err)
		assert.True(t, ready)
		assert.Equal(t, 1, calls)
	})

	t.Run("not ready is not cached", func(t *testing.T) {
		calls := 0
		c := &Commands{
			ownerDeleteReady: func(context.Context) (bool, error) {
				calls++
				return false, nil
			},
		}
		_, err := c.isUniqueConstraintOwnerDeleteReady(t.Context())
		require.NoError(t, err)
		_, err = c.isUniqueConstraintOwnerDeleteReady(t.Context())
		require.NoError(t, err)
		assert.Equal(t, 2, calls)
	})
}

func TestUniqueConstraintOwnersBackfillMatchesVersion(t *testing.T) {
	t.Run("no backfill event", func(t *testing.T) {
		c := &Commands{eventstore: expectEventstore(expectFilter())(t)}
		ready, err := c.uniqueConstraintOwnersBackfillMatchesVersion(t.Context())
		require.NoError(t, err)
		assert.False(t, ready)
	})

	t.Run("matching version", func(t *testing.T) {
		c := &Commands{eventstore: expectEventstore(
			expectFilter(uniqueConstraintOwnersBackfillDoneEvent(build.Version())),
		)(t)}
		ready, err := c.uniqueConstraintOwnersBackfillMatchesVersion(t.Context())
		require.NoError(t, err)
		assert.True(t, ready)
	})

	t.Run("different version", func(t *testing.T) {
		c := &Commands{eventstore: expectEventstore(
			expectFilter(uniqueConstraintOwnersBackfillDoneEvent("other-version")),
		)(t)}
		ready, err := c.uniqueConstraintOwnersBackfillMatchesVersion(t.Context())
		require.NoError(t, err)
		assert.False(t, ready)
	})

	t.Run("filter error", func(t *testing.T) {
		c := &Commands{eventstore: expectEventstore(
			expectFilterError(zerrors.ThrowInternal(nil, "id", "err")),
		)(t)}
		_, err := c.uniqueConstraintOwnersBackfillMatchesVersion(t.Context())
		assert.Error(t, err)
	})
}

func uniqueConstraintOwnersBackfillDoneEvent(version string) *repository.Event {
	ctx := authz.WithInstanceID(context.Background(), "")
	cmd := &migration.SetupStep{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			eventstore.NewAggregate(ctx, migration.SystemAggregateID, migration.SystemAggregate, "v1"),
			eventstore.EventType("system.migration.repeatable.done"),
		),
		Name:    eventstore.UniqueConstraintOwnersBackfillStep,
		LastRun: map[string]any{"version": version},
	}
	data, err := eventstore.EventData(cmd)
	if err != nil {
		panic(err)
	}
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
