package eventstore

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/postgres"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_lockRetryDelay(t *testing.T) {
	for attempt, want := range []time.Duration{5, 10, 20, 40, 80, 100, 100, 100} {
		want *= time.Millisecond
		seen := map[time.Duration]bool{}
		for range 50 {
			delay := lockRetryDelay(attempt)
			assert.GreaterOrEqual(t, delay, want/2)
			assert.LessOrEqual(t, delay, want)
			seen[delay] = true
		}
		assert.Greater(t, len(seen), 1, "attempt %d must be jittered", attempt)
	}
	assert.LessOrEqual(t, lockRetryDelay(1000), maxLockRetryDelay)
}

func TestWithInstanceLock(t *testing.T) {
	config, stop := postgres.StartEmbedded()
	defer stop()
	db, err := sql.Open("pgx", config.GetConnectionURL())
	require.NoError(t, err)
	defer db.Close()
	es := NewEventstore(&database.DB{DB: db})
	ctx := authz.WithInstanceID(t.Context(), "instance")

	t.Run("panic releases the lock", func(t *testing.T) {
		func() {
			defer func() { require.Equal(t, "boom", recover()) }()
			_ = es.WithInstanceLock(ctx, func(context.Context) error { panic("boom") })
		}()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		require.NoError(t, es.WithInstanceLock(ctx, func(context.Context) error { return nil }))
	})

	t.Run("wait is capped and uses the setup lock key", func(t *testing.T) {
		previous := maxInstanceLockWait
		maxInstanceLockWait = 200 * time.Millisecond
		t.Cleanup(func() { maxInstanceLockWait = previous })
		holder, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = holder.Rollback() }()
		// The blocking form cmd/setup/81.go uses must take the same lock.
		_, err = holder.ExecContext(ctx, "SELECT pg_advisory_xact_lock(hashtextextended('zitadel.group-lifecycle:' || $1, 0))", "instance")
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		start := time.Now()
		err = es.WithInstanceLock(ctx, func(context.Context) error { return nil })
		require.True(t, zerrors.IsDeadlineExceeded(err), "got %v", err)
		require.False(t, errors.Is(err, eventstore.ErrConcurrentModification), "callers must not retry a capped wait")
		require.Less(t, time.Since(start), 2*time.Second)
		require.NoError(t, es.WithInstanceLock(authz.WithInstanceID(ctx, "other"), func(context.Context) error { return nil }))
	})
}
