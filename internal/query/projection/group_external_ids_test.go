package projection

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func Test_GroupExternalIDsReduces(t *testing.T) {
	projection := &groupExternalIDsProjection{}
	digest := sha256.Sum256([]byte("object-id"))
	tests := []struct {
		name      string
		event     func(*testing.T) eventstore.Event
		reduce    func(eventstore.Event) (*handler.Statement, error)
		statement string
		args      []interface{}
	}{
		{name: "persist source binding", event: getEvent(testEvent(group.GroupExternalIDSetEventType, group.AggregateType, []byte(`{"scope":"entra","externalId":"object-id"}`)), eventstore.GenericEventMapper[group.GroupExternalIDSetEvent]), reduce: projection.reduceExternalIDSet, statement: "INSERT INTO projections.group_external_ids2 (instance_id, group_id, resource_owner, scope, external_id, external_id_sha256) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (instance_id, group_id, scope) DO UPDATE SET (resource_owner, external_id, external_id_sha256) = (EXCLUDED.resource_owner, EXCLUDED.external_id, EXCLUDED.external_id_sha256)", args: []interface{}{"instance-id", "agg-id", "ro-id", "entra", "object-id", digest[:]}},
		{name: "remove only one source", event: getEvent(testEvent(group.GroupExternalIDSetEventType, group.AggregateType, []byte(`{"scope":"entra","externalId":""}`)), eventstore.GenericEventMapper[group.GroupExternalIDSetEvent]), reduce: projection.reduceExternalIDSet, statement: "DELETE FROM projections.group_external_ids2 WHERE (instance_id = $1) AND (group_id = $2) AND (scope = $3)", args: []interface{}{"instance-id", "agg-id", "entra"}},
		{name: "delete group bindings", event: getEvent(testEvent(group.GroupRemovedEventType, group.AggregateType, nil), eventstore.GenericEventMapper[group.GroupRemovedEvent]), reduce: projection.reduceGroupRemoved, statement: "DELETE FROM projections.group_external_ids2 WHERE (group_id = $1) AND (resource_owner = $2) AND (instance_id = $3)", args: []interface{}{"agg-id", "ro-id", "instance-id"}},
		{name: "delete organization bindings", event: getEvent(testEvent(org.OrgRemovedEventType, org.AggregateType, nil), org.OrgRemovedEventMapper), reduce: projection.reduceOwnerRemoved, statement: "DELETE FROM projections.group_external_ids2 WHERE (resource_owner = $1) AND (instance_id = $2)", args: []interface{}{"agg-id", "instance-id"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.event(t)
			statement, err := tt.reduce(event)
			assertReduce(t, statement, err, GroupExternalIDsProjectionTable, wantReduce{aggregateType: event.Aggregate().Type, sequence: 15, executer: &testExecuter{executions: []execution{{expectedStmt: tt.statement, expectedArgs: tt.args}}}})
		})
	}
}

func Test_GroupExternalIDsDigest_database(t *testing.T) {
	ctx := context.Background()
	tx, err := rawDB.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer func() { require.NoError(t, tx.Rollback()) }()
	_, err = tx.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS projections")
	require.NoError(t, err)
	projection := &groupExternalIDsProjection{}
	for range 2 {
		for _, execute := range projection.Init().Executes {
			next, err := execute(ctx, tx, projection.Name())
			require.NoError(t, err)
			if !next {
				break
			}
		}
	}
	var definition string
	require.NoError(t, tx.QueryRowContext(ctx, "SELECT indexdef FROM pg_indexes WHERE schemaname = 'projections' AND tablename = 'group_external_ids2' AND indexname = 'group_external_ids2_lookup_idx'").Scan(&definition))
	require.Contains(t, definition, "USING btree (instance_id, resource_owner, scope, external_id_sha256)")
	require.NotContains(t, definition, "INCLUDE")

	// An incompressible ID would exceed the B-tree entry limit if indexed directly.
	var longID strings.Builder
	for i := range 512 {
		fmt.Fprintf(&longID, "%x", sha256.Sum256([]byte(fmt.Sprint(i))))
	}
	aggregate := eventstore.NewAggregate(ctx, "group", group.AggregateType, group.AggregateVersion, eventstore.WithInstanceID("instance"), eventstore.WithResourceOwner("org"))
	for _, value := range []string{longID.String(), `MiXeD,"{}\雪`, "replacement"} {
		event := group.NewGroupExternalIDSetEvent(ctx, aggregate, "entra", "", value)
		statement, err := projection.reduceExternalIDSet(event)
		require.NoError(t, err)
		require.NoError(t, statement.Execute(ctx, tx, projection.Name()))
		var storedValue string
		var storedDigest []byte
		require.NoError(t, tx.QueryRowContext(ctx, `SELECT external_id, external_id_sha256 FROM projections.group_external_ids2 WHERE group_id = 'group'`).Scan(&storedValue, &storedDigest))
		require.Equal(t, value, storedValue)
		digest := sha256.Sum256([]byte(value))
		require.Equal(t, digest[:], storedDigest)
		require.Len(t, storedDigest, 32)
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO projections.group_external_ids2
		(instance_id, group_id, resource_owner, scope, external_id, external_id_sha256)
		SELECT 'instance', 'filler-' || n, 'org', 'entra', n::text, sha256(n::text::bytea)
		FROM generate_series(1, 20000) AS n`)
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, "ANALYZE projections.group_external_ids2")
	require.NoError(t, err)
	digest := sha256.Sum256([]byte("replacement"))
	var plan string
	require.NoError(t, tx.QueryRowContext(ctx, `EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)
		SELECT group_id FROM projections.group_external_ids2
		WHERE instance_id = 'instance' AND resource_owner = 'org' AND scope = 'entra'
		AND external_id_sha256 = $1 AND external_id = $2`, digest[:], "replacement").Scan(&plan))
	require.Contains(t, plan, `"Index Name": "group_external_ids2_lookup_idx"`)
	require.Contains(t, plan, "external_id_sha256 =")
	t.Log(plan)
}
