package eventstore_test

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func claimedConstraint(action eventstore.UniqueConstraintAction, key, group string) *eventstore.UniqueConstraint {
	return (&eventstore.UniqueConstraint{Action: action, UniqueType: "test_claim", UniqueField: key}).WithClaimant("group", group)
}

func claimedCommand(typ eventstore.AggregateType, group string, revision *uint64, constraints ...*eventstore.UniqueConstraint) eventstore.Command {
	return generateCommand(typ, group, func(e *testEvent) {
		e.Agg.ExpectedSequence = revision
		e.uniqueConstraints = constraints
	})
}

func TestEventstore_ClaimedReservation_staleRelease(t *testing.T) {
	for name, pusher := range pushers {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			typ := eventstore.AggregateType(t.Name())
			a, b, c := strings.ToLower(t.Name())+"-a", strings.ToLower(t.Name())+"-b", strings.ToLower(t.Name())+"-c"
			push := func(group string, constraints ...*eventstore.UniqueConstraint) error {
				_, err := pusher.Push(ctx, nil, claimedCommand(typ, group, nil, constraints...))
				return err
			}
			require.NoError(t, push("g", claimedConstraint(eventstore.UniqueConstraintAdd, a, "g")))
			require.NoError(t, push("g", claimedConstraint(eventstore.UniqueConstraintRemove, a, "g"), claimedConstraint(eventstore.UniqueConstraintAdd, b, "g")))
			require.NoError(t, push("h", claimedConstraint(eventstore.UniqueConstraintAdd, a, "h")))
			// Deliberately omit a revision to test the ownership protection independently.
			require.ErrorIs(t, push("g", claimedConstraint(eventstore.UniqueConstraintRemove, a, "g"), claimedConstraint(eventstore.UniqueConstraintAdd, c, "g")), eventstore.ErrConcurrentModification)
			var claimant string
			require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT claimant_id FROM eventstore.unique_constraints WHERE unique_type = 'test_claim' AND unique_field = $1", a).Scan(&claimant))
			require.Equal(t, "h", claimant)
			// The rejected push cannot persist its event or acquire C.
			var count int
			require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.events2 WHERE aggregate_type = $1 AND aggregate_id = 'g'", typ).Scan(&count))
			require.Equal(t, 2, count)
			require.NoError(t, push("k", claimedConstraint(eventstore.UniqueConstraintAdd, c, "k")))
			// A legacy unowned release cannot remove H's owned reservation either.
			require.NoError(t, push("legacy", eventstore.NewRemoveUniqueConstraint("test_claim", a)))
			require.Error(t, push("k", claimedConstraint(eventstore.UniqueConstraintAdd, a, "k")))
		})
	}
}

func TestEventstore_ClaimedReservation_failedReplacementRollsBackRelease(t *testing.T) {
	for name, pusher := range pushers {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			typ := eventstore.AggregateType(t.Name())
			a, b := strings.ToLower(t.Name())+"-a", strings.ToLower(t.Name())+"-b"
			_, err := pusher.Push(ctx, nil, claimedCommand(typ, "g", nil, claimedConstraint(eventstore.UniqueConstraintAdd, a, "g")), claimedCommand(typ, "h", nil, claimedConstraint(eventstore.UniqueConstraintAdd, b, "h")))
			require.NoError(t, err)
			_, err = pusher.Push(ctx, nil, claimedCommand(typ, "g", gu.Ptr(uint64(1)), claimedConstraint(eventstore.UniqueConstraintRemove, a, "g"), claimedConstraint(eventstore.UniqueConstraintAdd, b, "g")))
			require.Error(t, err)
			var claimant string
			require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT claimant_id FROM eventstore.unique_constraints WHERE unique_type = 'test_claim' AND unique_field = $1", a).Scan(&claimant))
			require.Equal(t, "g", claimant)
		})
	}
}

func TestEventstore_ExpectedSequence_competingInitialBindings(t *testing.T) {
	for name, pusher := range pushers {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			typ := eventstore.AggregateType(t.Name())
			_, err := pusher.Push(ctx, nil, claimedCommand(typ, "g", gu.Ptr(uint64(0))))
			require.NoError(t, err)
			start := make(chan struct{})
			results := make(chan error, 2)
			var ready sync.WaitGroup
			ready.Add(2)
			for _, key := range []string{strings.ToLower(t.Name()) + "-b", strings.ToLower(t.Name()) + "-c"} {
				go func() {
					ready.Done()
					<-start
					_, err := pusher.Push(ctx, nil, claimedCommand(typ, "g", gu.Ptr(uint64(1)), claimedConstraint(eventstore.UniqueConstraintAdd, key, "g")))
					results <- err
				}()
			}
			ready.Wait()
			close(start)
			first, second := <-results, <-results
			require.NotEqual(t, first == nil, second == nil, "exactly one competing mutation must commit")
			var count int
			require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.events2 WHERE aggregate_type = $1", typ).Scan(&count))
			require.Equal(t, 2, count)
			require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.unique_constraints WHERE unique_type = 'test_claim' AND unique_field IN ($1, $2)", strings.ToLower(t.Name())+"-b", strings.ToLower(t.Name())+"-c").Scan(&count))
			require.Equal(t, 1, count)
			// A stale retry must still fail, even if it chooses a different fresh ID.
			_, err = pusher.Push(ctx, nil, claimedCommand(typ, "g", gu.Ptr(uint64(1)), claimedConstraint(eventstore.UniqueConstraintAdd, strings.ToLower(t.Name())+"-d", "g")))
			require.ErrorIs(t, err, eventstore.ErrConcurrentModification)
		})
	}
}

func TestEventstore_ExpectedSequence_bindingReturnsToSameOwner(t *testing.T) {
	for name, pusher := range pushers {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			typ := eventstore.AggregateType(t.Name())
			a, b, c := strings.ToLower(t.Name())+"-a", strings.ToLower(t.Name())+"-b", strings.ToLower(t.Name())+"-c"
			_, err := pusher.Push(ctx, nil, claimedCommand(typ, "g", gu.Ptr(uint64(0)), claimedConstraint(eventstore.UniqueConstraintAdd, a, "g")))
			require.NoError(t, err)
			_, err = pusher.Push(ctx, nil, claimedCommand(typ, "g", gu.Ptr(uint64(1)), claimedConstraint(eventstore.UniqueConstraintRemove, a, "g"), claimedConstraint(eventstore.UniqueConstraintAdd, b, "g")))
			require.NoError(t, err)
			_, err = pusher.Push(ctx, nil, claimedCommand(typ, "g", gu.Ptr(uint64(2)), claimedConstraint(eventstore.UniqueConstraintRemove, b, "g"), claimedConstraint(eventstore.UniqueConstraintAdd, a, "g")))
			require.NoError(t, err)
			// Ownership looks unchanged, but the original snapshot is obsolete.
			_, err = pusher.Push(ctx, nil, claimedCommand(typ, "g", gu.Ptr(uint64(1)), claimedConstraint(eventstore.UniqueConstraintRemove, a, "g"), claimedConstraint(eventstore.UniqueConstraintAdd, c, "g")))
			require.ErrorIs(t, err, eventstore.ErrConcurrentModification)
			_, err = pusher.Push(ctx, nil, claimedCommand(typ, "h", nil, claimedConstraint(eventstore.UniqueConstraintAdd, a, "h")))
			require.Error(t, err, "A must remain reserved by G")
			_, err = pusher.Push(ctx, nil, claimedCommand(typ, "h", nil, claimedConstraint(eventstore.UniqueConstraintAdd, c, "h")))
			require.NoError(t, err, "the stale push must not reserve C")
		})
	}
}

func TestEventstore_ClaimedReservation_schemaUpgrade(t *testing.T) {
	ctx := context.Background()
	statement, err := os.ReadFile("../../cmd/setup/76.sql")
	require.NoError(t, err)
	tx, err := testClient.BeginTx(ctx, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, tx.Rollback()) })
	// Exercise the migration against the old table shape without changing the
	// shared test database: all schema and fixture changes roll back at cleanup.
	_, err = tx.ExecContext(ctx, "ALTER TABLE eventstore.unique_constraints DROP COLUMN claimant_type, DROP COLUMN claimant_id")
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, "INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES ('migration', 'legacy', 'key')")
	require.NoError(t, err)
	for range 2 {
		_, err = tx.ExecContext(ctx, string(statement))
		require.NoError(t, err)
	}
	var count int
	err = tx.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.unique_constraints WHERE instance_id = 'migration' AND unique_type = 'legacy' AND unique_field = 'key' AND claimant_type IS NULL AND claimant_id IS NULL").Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}
