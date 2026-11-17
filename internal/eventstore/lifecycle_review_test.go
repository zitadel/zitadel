package eventstore_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"

	new_db "github.com/zitadel/zitadel/backend/v3/storage/database"
	new_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
)

func TestEventstore_LifecycleSequenceCollisionIsConcurrentModification(t *testing.T) {
	ctx, cancel := context.WithTimeout(authz.WithInstanceID(t.Context(), t.Name()), 10*time.Second)
	defer cancel()
	pusher := pushers["v3(inmemory)"]
	es := eventstore.NewEventstore(&eventstore.Config{Pusher: pusher, Querier: queriers["v2(inmemory)"], MaxRetries: 3})
	typ := eventstore.AggregateType(t.Name())
	command := func() eventstore.Command {
		return generateCommand(typ, "same", func(e *testEvent) { e.Agg.InstanceID = t.Name() })
	}
	// An unlocked writer holds the next sequence of the same aggregate until the lifecycle push collides with it.
	other, err := testClient.BeginTx(ctx, nil)
	require.NoError(t, err)
	_, err = pusher.Push(ctx, new_sql.SQLTx(other), command())
	require.NoError(t, err)
	go func() {
		time.Sleep(300 * time.Millisecond)
		_ = other.Commit()
	}()
	err = es.WithInstanceLock(ctx, func(txctx context.Context) error {
		_, err := es.Push(txctx, command())
		return err
	})
	require.ErrorIs(t, err, eventstore.ErrConcurrentModification)
	var pgErr *pgconn.PgError
	require.ErrorAs(t, err, &pgErr)
	require.Equal(t, "23505", pgErr.SQLState())
}

func TestEventstore_FilterKeepsExplicitTransaction(t *testing.T) {
	ctx, cancel := context.WithTimeout(authz.WithInstanceID(t.Context(), t.Name()), 10*time.Second)
	defer cancel()
	es := eventstore.NewEventstore(&eventstore.Config{Pusher: pushers["v3(inmemory)"], Querier: queriers["v2(inmemory)"]})
	typ := eventstore.AggregateType(t.Name())
	query := func() *eventstore.SearchQueryBuilder {
		return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).InstanceID(t.Name()).AddQuery().AggregateTypes(typ).Builder()
	}
	rollback := errors.New("rollback")
	err := es.WithInstanceLock(ctx, func(txctx context.Context) error {
		_, err := es.Push(txctx, generateCommand(typ, "x", func(e *testEvent) { e.Agg.InstanceID = t.Name() }))
		require.NoError(t, err)
		events, err := es.Filter(txctx, query())
		require.NoError(t, err)
		require.Len(t, events, 1, "command reads see the lifecycle transaction")
		// Projection handlers read through their own transaction.
		handlerTx, err := testClient.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = handlerTx.Rollback() }()
		events, err = es.Filter(txctx, query().SetTx(handlerTx))
		require.NoError(t, err)
		require.Empty(t, events, "uncommitted lifecycle events must not leak into an explicit transaction")
		reducer := &countReducer{}
		require.NoError(t, es.FilterToReducer(txctx, query().SetTx(handlerTx), reducer))
		require.Zero(t, reducer.count)
		return rollback
	})
	require.ErrorIs(t, err, rollback)
}

type countReducer struct{ count int }

func (r *countReducer) AppendEvents(events ...eventstore.Event) { r.count += len(events) }
func (*countReducer) Reduce() error                             { return nil }

func TestEventstore_ClaimedReservation_releaseMissingRow(t *testing.T) {
	for name, pusher := range pushers {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			typ := eventstore.AggregateType(t.Name())
			a, b := strings.ToLower(t.Name())+"-a", strings.ToLower(t.Name())+"-b"
			push := func(group string, constraints ...*eventstore.UniqueConstraint) error {
				_, err := pusher.Push(ctx, nil, claimedCommand(typ, group, nil, constraints...))
				return err
			}
			require.NoError(t, push("g", claimedConstraint(eventstore.UniqueConstraintAdd, a, "g")))
			_, err := clients[name].ExecContext(ctx, "DELETE FROM eventstore.unique_constraints WHERE unique_type = 'test_claim' AND unique_field = $1", a)
			require.NoError(t, err)
			require.NoError(t, push("g", claimedConstraint(eventstore.UniqueConstraintRemove, a, "g"), claimedConstraint(eventstore.UniqueConstraintAdd, b, "g")), "a missing reservation counts as released")
			require.NoError(t, push("h", claimedConstraint(eventstore.UniqueConstraintAdd, a, "h")))
			require.ErrorIs(t, push("g", claimedConstraint(eventstore.UniqueConstraintRemove, a, "g")), eventstore.ErrConcurrentModification, "another claimant still blocks the release")
		})
	}
}

func TestEventstore_ClaimedReservation_addThenReleaseSamePush(t *testing.T) {
	for name, pusher := range pushers {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			typ := eventstore.AggregateType(t.Name())
			a, b := strings.ToLower(t.Name())+"-a", strings.ToLower(t.Name())+"-b"
			count := func(key string) (n int) {
				require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.unique_constraints WHERE unique_type = 'test_claim' AND unique_field = $1", key).Scan(&n))
				return n
			}
			_, err := pusher.Push(ctx, nil,
				claimedCommand(typ, "g", nil, claimedConstraint(eventstore.UniqueConstraintAdd, a, "g")),
				claimedCommand(typ, "g", nil, claimedConstraint(eventstore.UniqueConstraintRemove, strings.ToUpper(a), "g")))
			require.NoError(t, err)
			require.Zero(t, count(a))
			_, err = pusher.Push(ctx, nil,
				claimedCommand(typ, "g", nil, claimedConstraint(eventstore.UniqueConstraintAdd, b, "g")),
				claimedCommand(typ, "g", nil, claimedConstraint(eventstore.UniqueConstraintRemove, b, "g")),
				claimedCommand(typ, "g", nil, claimedConstraint(eventstore.UniqueConstraintAdd, b, "g")))
			require.NoError(t, err)
			require.Equal(t, 1, count(b))
			_, err = pusher.Push(ctx, nil,
				claimedCommand(typ, "g", nil, claimedConstraint(eventstore.UniqueConstraintAdd, a, "g")),
				claimedCommand(typ, "h", nil, claimedConstraint(eventstore.UniqueConstraintRemove, a, "h")))
			require.ErrorIs(t, err, eventstore.ErrConcurrentModification)
			require.Zero(t, count(a))

			stored := func(key, claimant string) {
				_, err := pusher.Push(ctx, nil, claimedCommand(typ, claimant, nil, claimedConstraint(eventstore.UniqueConstraintAdd, key, claimant)))
				require.NoError(t, err)
			}
			c, d := strings.ToLower(t.Name())+"-c", strings.ToLower(t.Name())+"-d"
			stored(c, "g")
			_, err = pusher.Push(ctx, nil,
				claimedCommand(typ, "g", nil, claimedConstraint(eventstore.UniqueConstraintAdd, c, "g")),
				claimedCommand(typ, "g", nil, claimedConstraint(eventstore.UniqueConstraintRemove, c, "g")))
			require.NoError(t, err)
			require.Zero(t, count(c), "the stored row of the same claimant is released")
			stored(d, "h")
			_, err = pusher.Push(ctx, nil,
				claimedCommand(typ, "g", nil, claimedConstraint(eventstore.UniqueConstraintAdd, d, "g")),
				claimedCommand(typ, "g", nil, claimedConstraint(eventstore.UniqueConstraintRemove, d, "g")))
			require.ErrorIs(t, err, eventstore.ErrConcurrentModification, "a stored row of another claimant blocks the push")
			require.Equal(t, 1, count(d))
		})
	}
}

func TestEventstore_ClaimedReservation_releaseOrderInOnePush(t *testing.T) {
	for name, pusher := range pushers {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			typ := eventstore.AggregateType(t.Name())
			key := strings.ToLower(t.Name())
			push := func(group string, constraints ...*eventstore.UniqueConstraint) error {
				_, err := pusher.Push(ctx, nil, claimedCommand(typ, group, nil, constraints...))
				return err
			}
			require.NoError(t, push("g", claimedConstraint(eventstore.UniqueConstraintAdd, key, "g")))
			require.ErrorIs(t, push("h", claimedConstraint(eventstore.UniqueConstraintRemove, key, "h"), claimedConstraint(eventstore.UniqueConstraintRemove, key, "g")),
				eventstore.ErrConcurrentModification, "a release before the holder's release sees the held key")
			require.NoError(t, push("g", claimedConstraint(eventstore.UniqueConstraintRemove, key, "g"), claimedConstraint(eventstore.UniqueConstraintRemove, key, "h")),
				"a release after the holder's release sees a free key")
		})
	}
}

type countingTx struct {
	new_db.Transaction
	statements int
}

func (tx *countingTx) Exec(ctx context.Context, stmt string, args ...any) (int64, error) {
	tx.statements++
	return tx.Transaction.Exec(ctx, stmt, args...)
}

func (tx *countingTx) Query(ctx context.Context, stmt string, args ...any) (new_db.Rows, error) {
	tx.statements++
	return tx.Transaction.Query(ctx, stmt, args...)
}

func (tx *countingTx) QueryRow(ctx context.Context, stmt string, args ...any) new_db.Row {
	tx.statements++
	return tx.Transaction.QueryRow(ctx, stmt, args...)
}

// An organization removal pushes one group.removed per group inside the
// instance lock, so its statement count must not grow with the group count.
func TestEventstore_GroupCascadeStatementsDoNotGrow(t *testing.T) {
	instance := t.Name()
	ctx := authz.WithInstanceID(context.Background(), instance)
	pusher := pushers["v3(inmemory)"]
	var fields *string
	require.NoError(t, testClient.DB.QueryRowContext(ctx, "SELECT to_regclass('eventstore.fields')::TEXT").Scan(&fields))
	if fields == nil {
		ddl, err := os.ReadFile("../../cmd/setup/28.sql")
		require.NoError(t, err)
		_, err = testClient.ExecContext(ctx, string(ddl))
		require.NoError(t, err)
	}
	statements := func(groups int) int {
		org := fmt.Sprintf("org-%d", groups)
		removals := make([]eventstore.Command, groups)
		for i := range groups {
			agg := &group.NewAggregate(fmt.Sprintf("%s-%d", org, i), org).Aggregate
			agg.InstanceID = instance
			set := group.NewGroupExternalIDSetEvent(ctx, agg, "scope", "", "ext-"+agg.ID)
			_, err := pusher.Push(ctx, nil, group.NewGroupAddedEvent(ctx, agg, "name-"+agg.ID, ""), set,
				group.NewGroupUserAddedEvent(ctx, agg, "u1"), group.NewGroupUserAddedEvent(ctx, agg, "u2"))
			require.NoError(t, err)
			var keys []string
			for _, constraint := range set.UniqueConstraints() {
				keys = append(keys, constraint.UniqueField)
			}
			removals[i] = group.NewGroupRemovedEvent(ctx, agg, "name-"+agg.ID).WithMembers([]string{"u1", "u2"}).WithExternalIDKeys(keys)
		}
		sqlTx, err := testClient.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = sqlTx.Rollback() }()
		tx := &countingTx{Transaction: new_sql.SQLTx(sqlTx)}
		_, err = pusher.Push(ctx, tx, removals...)
		require.NoError(t, err)
		var left int
		require.NoError(t, sqlTx.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.fields WHERE instance_id = $1 AND resource_owner = $2", instance, org).Scan(&left))
		require.Zero(t, left)
		require.NoError(t, sqlTx.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = 'group_external_id' AND EXISTS (SELECT 1 FROM unnest(owners) AS tag WHERE tag LIKE $2)", instance, "claimant:group:"+org+"-%").Scan(&left))
		require.Zero(t, left)
		return tx.statements
	}
	require.Equal(t, statements(2), statements(50))
}
