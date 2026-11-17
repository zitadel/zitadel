package eventstore_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	new_es "github.com/zitadel/zitadel/internal/eventstore/v3"
)

func TestEventstore_LifecycleLateCommitPosition(t *testing.T) {
	ctx, cancel := context.WithTimeout(authz.WithInstanceID(t.Context(), t.Name()), 5*time.Second)
	defer cancel()
	conn, err := testClient.Conn(ctx)
	require.NoError(t, err)
	require.NoError(t, conn.Raw(func(driver any) error { return new_es.RegisterEventstoreTypes(ctx, driver.(*stdlib.Conn).Conn()) }))
	require.NoError(t, conn.Close())
	es := eventstore.NewEventstore(&eventstore.Config{Pusher: pushers["v3(inmemory)"], Querier: queriers["v2(inmemory)"]})
	typ := eventstore.AggregateType(t.Name())
	command := func(id string) eventstore.Command {
		return generateCommand(typ, id, func(e *testEvent) { e.Agg.InstanceID = t.Name() })
	}
	entered, release := make(chan struct{}), make(chan struct{})
	done := make(chan error, 1)
	var early eventstore.Event
	go func() {
		done <- es.WithInstanceLock(ctx, func(txctx context.Context) error {
			close(entered)
			select {
			case <-release:
			case <-ctx.Done():
				return ctx.Err()
			}
			events, err := es.Push(txctx, command("pending-delete"))
			if err == nil {
				early = events[0]
			}
			return err
		})
	}()
	<-entered
	// User deletion now starts its transaction before read/permission checks;
	// user profile changes and creation do not take the lifecycle lock.
	later, err := es.Push(ctx, command("concurrent-profile-change"))
	require.NoError(t, err)
	query := func() *eventstore.SearchQueryBuilder {
		return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).AwaitOpenTransactions().OrderAsc().AddQuery().AggregateTypes(typ).Builder()
	}
	visible, err := es.Filter(ctx, query())
	require.NoError(t, err)
	require.Len(t, visible, 1)
	cursor := visible[0].Position()
	require.Equal(t, later[0].Position(), cursor)
	close(release)
	require.NoError(t, <-done)
	require.True(t, early.Position().GreaterThan(cursor), "late commit must follow the cursor")
	visible, err = es.Filter(ctx, query().PositionAtLeast(cursor))
	require.NoError(t, err)
	require.Len(t, visible, 2)
	require.Equal(t, "pending-delete", visible[1].Aggregate().ID)
}

func TestEventstore_LifecycleReadsReuseConnection(t *testing.T) {
	previous := testClient.Stats().MaxOpenConnections
	testClient.SetMaxOpenConns(1)
	t.Cleanup(func() { testClient.SetMaxOpenConns(previous) })
	ctx, cancel := context.WithTimeout(authz.WithInstanceID(t.Context(), t.Name()), time.Second)
	defer cancel()
	es := eventstore.NewEventstore(&eventstore.Config{Pusher: pushers["v3(inmemory)"], Querier: queriers["v2(inmemory)"]})
	err := es.WithInstanceLock(ctx, func(txctx context.Context) error {
		if err := testClient.QueryRowContext(txctx, func(row *sql.Row) error { var n int; return row.Scan(&n) }, "SELECT 1"); err != nil {
			return err
		}
		return testClient.QueryContext(txctx, func(rows *sql.Rows) error {
			require.True(t, rows.Next())
			var n int
			return rows.Scan(&n)
		}, "SELECT 1")
	})
	require.NoError(t, err)
}
