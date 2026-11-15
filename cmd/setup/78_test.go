package setup

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/postgres"
)

func TestGroupLifecycleFields_Backfill(t *testing.T) {
	config, stop := postgres.StartEmbedded()
	defer stop()
	db, err := sql.Open("pgx", config.GetConnectionURL())
	require.NoError(t, err)
	defer db.Close()
	client := &database.DB{DB: db}
	ctx := context.Background()
	_, err = db.ExecContext(ctx, `CREATE SCHEMA eventstore;
 CREATE TABLE eventstore.events2 (instance_id TEXT, aggregate_type TEXT, aggregate_id TEXT, owner TEXT, event_type TEXT, sequence BIGINT, payload JSONB, PRIMARY KEY (instance_id, aggregate_type, aggregate_id, sequence));
 CREATE TABLE eventstore.unique_constraints (instance_id TEXT, unique_type TEXT, unique_field TEXT, claimant_type TEXT, claimant_id TEXT);`)
	require.NoError(t, err)
	fields, err := os.ReadFile("28.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, string(fields))
	require.NoError(t, err)
	aggregateIndex, err := os.ReadFile("31.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, string(aggregateIndex))
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO eventstore.events2 VALUES
 ('instance', 'group', 'live', 'org', 'group.added', 1, '{"name":"old"}'),
 ('instance', 'group', 'live', 'org', 'group.changed', 2, '{"name":"renamed"}'),
 ('instance', 'group', 'live', 'org', 'group.changed', 3, '{"description":"no rename"}'),
 ('instance', 'group', 'live', 'org', 'group.user.added', 4, '{"userId":"current"}'),
 ('instance', 'group', 'live', 'org', 'group.user.added', 5, '{"userId":"removed"}'),
 ('instance', 'group', 'live', 'org', 'group.user.removed', 6, '{"userId":"removed"}'),
 ('instance', 'group', 'live', 'org', 'group.user.added', 7, '{"userId":"deleted-user"}'),
 ('instance', 'group', 'live', 'org', 'group.users.added', 8, '{"userIds":["legacy-a","legacy-b"]}'),
 ('instance', 'group', 'live', 'org', 'group.users.removed', 9, '{"userIds":["legacy-b"]}'),
 ('instance', 'group', 'live', 'org', 'group.user.added', 10, '{"userId":"new-c"}'),
 ('instance', 'group', 'live', 'org', 'group.users.added', 11, '{"userIds":["legacy-d"]}'),
 ('instance', 'group', 'live', 'org', 'group.user.removed', 12, '{"userId":"legacy-d"}'),
 ('instance', 'user', 'deleted-user', 'org', 'user.removed', 1, '{}'),
 ('instance', 'group', 'gone', 'org', 'group.added', 1, '{"name":"gone"}'),
 ('instance', 'group', 'gone', 'org', 'group.user.added', 2, '{"userId":"current"}'),
 ('instance', 'group', 'gone', 'org', 'group.removed', 3, '{}'),
 ('instance', 'group', 'gone-org', 'deleted-org', 'group.added', 1, '{"name":"gone-org"}'),
 ('instance', 'org', 'deleted-org', 'deleted-org', 'org.removed', 1, '{}'),
 ('removed-instance', 'group', 'gone-instance', 'org', 'group.added', 1, '{"name":"gone-instance"}'),
 ('removed-instance', 'instance', 'removed-instance', '', 'instance.removed', 1, '{}'),
 ('other-instance', 'group', 'live', 'org', 'group.added', 1, '{"name":"other"}');
 INSERT INTO eventstore.unique_constraints VALUES
 ('instance', 'group_external_id', 'current-binding-key', 'group', 'live'),
 ('instance', 'group_external_id', 'removed-binding-key', 'group', 'gone'),
 ('instance', 'group_external_id', 'unclaimed-binding-key', NULL, NULL),
 ('other-instance', 'group_external_id', 'other-binding-key', 'group', 'live'),
 ('other-instance', 'group_external_id', 'other-removed-key', 'group', 'gone');`)
	require.NoError(t, err)
	migration := &GroupLifecycleFields{dbClient: client}
	for range 2 {
		require.NoError(t, migration.Execute(ctx, nil))
	}
	// Simulate the catalog state left by an interrupted concurrent index build.
	_, err = db.ExecContext(ctx, "UPDATE pg_index SET indisvalid = false WHERE indexrelid = 'eventstore.f_owner_object_type_idx'::regclass")
	require.NoError(t, err)
	require.NoError(t, migration.Execute(ctx, nil))
	var valid bool
	require.NoError(t, db.QueryRowContext(ctx, "SELECT indisvalid FROM pg_index WHERE indexrelid = 'eventstore.f_owner_object_type_idx'::regclass").Scan(&valid))
	require.True(t, valid)
	rows, err := db.QueryContext(ctx, `SELECT instance_id, aggregate_id, field_name, value #>> '{}' FROM eventstore.fields ORDER BY instance_id, aggregate_id, field_name, value`)
	require.NoError(t, err)
	defer rows.Close()
	var got [][4]string
	for rows.Next() {
		var row [4]string
		require.NoError(t, rows.Scan(&row[0], &row[1], &row[2], &row[3]))
		got = append(got, row)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, [][4]string{
		{"instance", "live", "external_id_key", "current-binding-key"},
		{"instance", "live", "member", "current"},
		{"instance", "live", "member", "deleted-user"},
		{"instance", "live", "member", "legacy-a"},
		{"instance", "live", "member", "new-c"},
		{"instance", "live", "name", "renamed"},
		{"other-instance", "live", "external_id_key", "other-binding-key"},
		{"other-instance", "live", "name", "other"},
	}, got)
	var keys []string
	require.NoError(t, client.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			var key string
			if err := rows.Scan(&key); err != nil {
				return err
			}
			keys = append(keys, key)
		}
		return rows.Err()
	}, "SELECT instance_id || '/' || unique_field FROM eventstore.unique_constraints ORDER BY 1"))
	require.Equal(t, []string{"instance/current-binding-key", "instance/unclaimed-binding-key", "other-instance/other-binding-key"}, keys,
		"claims of removed groups must be released")
	_, err = db.ExecContext(ctx, `INSERT INTO eventstore.fields
 (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, should_index)
 SELECT 'instance', 'org-' || n % 1000, 'group', 'g-' || n, 'group', 'u-' || n, 1, 'member', to_jsonb('u-' || n), true
 FROM generate_series(1, 20000) AS n`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "ANALYZE eventstore.fields")
	require.NoError(t, err)
	var plan string
	require.NoError(t, db.QueryRowContext(ctx, `EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)
 SELECT aggregate_id FROM eventstore.fields WHERE instance_id = 'instance' AND resource_owner = 'org'
 AND aggregate_type = 'group' AND object_type = 'group' AND field_name = 'member' AND text_value = 'current'`).Scan(&plan))
	require.Contains(t, plan, `"Index Name": "f_text_value_idx"`)
	require.Contains(t, plan, `"Actual Rows": 1`)
	conn, err := db.Conn(ctx)
	require.NoError(t, err)
	defer conn.Close()
	_, err = conn.ExecContext(ctx, `SET plan_cache_mode = force_generic_plan;
 PREPARE org_groups AS SELECT aggregate_id FROM eventstore.fields WHERE instance_id = $1 AND aggregate_type = $2 AND object_type = $3 AND resource_owner = $4`)
	require.NoError(t, err)
	require.NoError(t, conn.QueryRowContext(ctx, `EXPLAIN (ANALYZE, FORMAT JSON) EXECUTE org_groups('instance', 'group', 'group', 'org-7')`).Scan(&plan))
	require.Contains(t, plan, `"Index Name": "f_owner_object_type_idx"`)
	require.Contains(t, plan, `"Actual Rows": 20`)
	_, err = db.ExecContext(ctx, `INSERT INTO eventstore.events2
 SELECT 'bulk-' || n % 10, 'user', 'u-' || n, 'org', 'user.added', 1, '{}' FROM generate_series(1, 20000) AS n;
 ANALYZE eventstore.events2`)
	require.NoError(t, err)
	require.NoError(t, db.QueryRowContext(ctx, "EXPLAIN (FORMAT JSON) "+groupInstancesQuery).Scan(&plan))
	require.NotContains(t, plan, `"Seq Scan"`)

	_, err = db.ExecContext(ctx, `INSERT INTO eventstore.fields
 (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, should_index)
 SELECT 'instance', 'org', 'group', 'big', 'group', 'u-' || n, 1, 'member', to_jsonb('u-' || n), true
 FROM generate_series(1, 30000) AS n;
 INSERT INTO eventstore.fields
 (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, should_index)
 SELECT 'instance', 'org', 'project', 'p-' || n % 50000, 'project_grant', 'o-' || n, 1, 'f', to_jsonb('v-' || n), true
 FROM generate_series(1, 400000) AS n;
 ANALYZE eventstore.fields`)
	require.NoError(t, err)
	// One batch of 500 member changes, the size of a large SCIM patch.
	columns := make([][]string, 9)
	indexed := make([]bool, 500)
	for i := range 500 {
		user := fmt.Sprintf("u-%d", i*7)
		for column, value := range []string{"group", "big", "instance", "org", "group", user, "1", "member", `"` + user + `"`} {
			columns[column] = append(columns[column], value)
		}
		indexed[i] = true
	}
	key := make([]any, 8)
	for i := range key {
		key[i] = columns[i]
	}
	for name, args := range map[string][]any{
		"field_batch_set.sql":    append(key, columns[8], indexed),
		"field_batch_remove.sql": key,
	} {
		stmt, err := os.ReadFile("../../internal/eventstore/v3/" + name)
		require.NoError(t, err)
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NoError(t, tx.QueryRowContext(ctx, "EXPLAIN (ANALYZE, FORMAT JSON) "+string(stmt), args...).Scan(&plan))
		require.NoError(t, tx.Rollback())
		require.Contains(t, plan, `"Index Name": "f_aggregate_object_idx"`, name)
		require.NotRegexp(t, `"Rows Removed by Filter": [1-9]`, plan, name)
	}

	// Group pushes take the same key in internal/eventstore/v3/transaction.go.
	lock, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer lock.Rollback()
	_, err = lock.ExecContext(ctx, "SELECT pg_advisory_xact_lock(hashtextextended('zitadel.group-lifecycle:' || $1, 0))", "instance")
	require.NoError(t, err)
	waitCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	require.Error(t, migration.backfillInstance(waitCtx, "instance"), "the backfill must wait for concurrent group pushes")
}

func TestGroupLifecycleFields_Check(t *testing.T) {
	migration := &GroupLifecycleFields{Rerun: 2}
	require.True(t, migration.Check(nil), "a run without a recorded counter must repeat")
	require.True(t, migration.Check(map[string]any{"rerun": float64(1)}))
	require.False(t, migration.Check(map[string]any{"rerun": float64(2)}))

	lastRun := map[string]any{}
	data, err := json.Marshal(migration)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &lastRun))
	require.False(t, migration.Check(lastRun), "the recorded done payload must stop the next run")
}
