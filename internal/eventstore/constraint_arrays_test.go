package eventstore_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestEventstore_UniqueConstraintArrays_largeBatch(t *testing.T) {
	for name, pusher := range pushers {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			typ := eventstore.AggregateType(t.Name())
			// Exceed both the old five-parameter insert and three-parameter delete limits.
			const size = 22000
			adds := make([]*eventstore.UniqueConstraint, size)
			removes := make([]*eventstore.UniqueConstraint, size)
			for i := range adds {
				key := fmt.Sprintf("key-%d", i)
				adds[i] = eventstore.NewAddEventUniqueConstraint(string(typ), key, "duplicate")
				removes[i] = eventstore.NewRemoveUniqueConstraint(string(typ), key)
			}
			push := func(constraints ...*eventstore.UniqueConstraint) error {
				_, err := pusher.Push(ctx, nil, claimedCommand(typ, "group", nil, constraints...))
				return err
			}
			count := func(wantConstraints, wantEvents int) {
				var got int
				require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.unique_constraints WHERE unique_type = $1", typ).Scan(&got))
				require.Equal(t, wantConstraints, got)
				require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.events2 WHERE aggregate_type = $1", typ).Scan(&got))
				require.Equal(t, wantEvents, got)
			}
			require.NoError(t, push(adds...))
			count(size, 1)
			// A duplicate after the bulk release must roll back both rows and event.
			conflict := eventstore.NewAddEventUniqueConstraint(string(typ), "conflict", "duplicate")
			require.Error(t, push(append(removes, conflict, conflict)...))
			count(size, 1)
			require.NoError(t, push(removes...))
			count(0, 2)
		})
	}
}

func TestEventstore_UniqueConstraintArrays_legacyMatching(t *testing.T) {
	for name, pusher := range pushers {
		t.Run(name, func(t *testing.T) {
			for _, tc := range []struct {
				name       string
				exact      bool
				owned      bool
				wantFields []string
			}{
				{name: "exact first", exact: true, wantFields: []string{"mixed"}},
				{name: "owned exact blocks fallback", exact: true, owned: true, wantFields: []string{"MiXeD", "mixed"}},
				{name: "lowercase fallback", wantFields: []string{}},
				{name: "owned lowercase", owned: true, wantFields: []string{"mixed"}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					ctx := context.Background()
					typ := eventstore.AggregateType(t.Name())
					fields := []string{"mixed"}
					if tc.exact {
						fields = append(fields, "MiXeD")
					}
					for _, field := range fields {
						owner := ""
						if tc.owned && (field == "MiXeD" || !tc.exact) {
							owner = "group"
						}
						_, err := clients[name].ExecContext(ctx, `INSERT INTO eventstore.unique_constraints
							(instance_id, unique_type, unique_field, owners)
							VALUES ('', $1, $2, CASE WHEN $3::TEXT = '' THEN '{}'::TEXT[] ELSE ARRAY['claimant:group:' || $3::TEXT] END)`, typ, field, owner)
						require.NoError(t, err)
					}
					_, err := pusher.Push(ctx, nil, claimedCommand(typ, "group", nil, eventstore.NewRemoveUniqueConstraint(string(typ), "MiXeD")))
					require.NoError(t, err)
					rows, err := clients[name].DB.QueryContext(ctx, "SELECT unique_field FROM eventstore.unique_constraints WHERE unique_type = $1 ORDER BY unique_field COLLATE \"C\"", typ)
					require.NoError(t, err)
					defer rows.Close()
					got := []string{}
					for rows.Next() {
						var field string
						require.NoError(t, rows.Scan(&field))
						got = append(got, field)
					}
					require.NoError(t, rows.Err())
					require.Equal(t, tc.wantFields, got)
				})
			}
		})
	}
}

func TestEventstore_UniqueConstraintArrays_scopesAndEscaping(t *testing.T) {
	for name, pusher := range pushers {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			typ := eventstore.AggregateType(t.Name())
			instance := t.Name() + `,"{}\雪`
			otherInstance := instance + "-other"
			key := `MiXeD,"{}\雪`
			push := func(instance string, constraints ...*eventstore.UniqueConstraint) error {
				command := claimedCommand(typ, "group", nil, constraints...)
				command.Aggregate().InstanceID = instance
				_, err := pusher.Push(ctx, nil, command)
				return err
			}
			add := func(owned bool) *eventstore.UniqueConstraint {
				constraint := eventstore.NewAddEventUniqueConstraint(string(typ), key, "duplicate")
				if owned {
					constraint.WithClaimant("group", key)
				}
				return constraint
			}
			require.NoError(t, push(instance, add(false), eventstore.NewAddGlobalUniqueConstraint(string(typ), key, "duplicate")))
			require.NoError(t, push(otherInstance, add(true)))
			var field, owner string
			require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT unique_field, owners[1] FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2", otherInstance, typ).Scan(&field, &owner))
			require.Equal(t, strings.ToLower(key), field)
			require.Equal(t, "claimant:group:"+key, owner)
			require.NoError(t, push(instance, eventstore.NewRemoveUniqueConstraint(string(typ), key)))
			// The global and other instance's owned reservation must survive the unowned release.
			var count int
			require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.unique_constraints WHERE unique_type = $1", typ).Scan(&count))
			require.Equal(t, 2, count)
			// Instance removal clears both ownership kinds, then additions in the same push take effect.
			require.NoError(t, push(otherInstance, eventstore.NewAddEventUniqueConstraint(string(typ), "unowned", "duplicate")))
			require.NoError(t, push(otherInstance, eventstore.NewRemoveInstanceUniqueConstraints(), add(false)))
			require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.unique_constraints WHERE unique_type = $1", typ).Scan(&count))
			require.Equal(t, 2, count)
			require.NoError(t, push(otherInstance, eventstore.NewRemoveUniqueConstraint(string(typ), key), eventstore.NewRemoveGlobalUniqueConstraint(string(typ), key)))
			require.NoError(t, clients[name].DB.QueryRowContext(ctx, "SELECT count(*) FROM eventstore.unique_constraints WHERE unique_type = $1", typ).Scan(&count))
			require.Zero(t, count)
		})
	}
}
