package setup

import (
	"context"
	"database/sql"
	"os"
	"testing"

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
 CREATE TABLE eventstore.events2 (instance_id TEXT, aggregate_type TEXT, aggregate_id TEXT, owner TEXT, event_type TEXT, sequence BIGINT, payload JSONB);
 CREATE TABLE eventstore.unique_constraints (instance_id TEXT, unique_type TEXT, unique_field TEXT, claimant_type TEXT, claimant_id TEXT);`)
	require.NoError(t, err)
	fields, err := os.ReadFile("28.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, string(fields))
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO eventstore.events2 VALUES
 ('instance', 'group', 'live', 'org', 'group.added', 1, '{"name":"old"}'),
 ('instance', 'group', 'live', 'org', 'group.changed', 2, '{"name":"renamed"}'),
 ('instance', 'group', 'live', 'org', 'group.changed', 3, '{"description":"no rename"}'),
 ('instance', 'group', 'live', 'org', 'group.user.added', 4, '{"userId":"current"}'),
 ('instance', 'group', 'live', 'org', 'group.user.added', 5, '{"userId":"removed"}'),
 ('instance', 'group', 'live', 'org', 'group.user.removed', 6, '{"userId":"removed"}'),
 ('instance', 'group', 'live', 'org', 'group.user.added', 7, '{"userId":"deleted-user"}'),
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
 ('instance', 'group_external_id', 'removed-binding-key', 'group', 'gone');`)
	require.NoError(t, err)
	migration := &GroupLifecycleFields{dbClient: client}
	for range 2 {
		require.NoError(t, migration.Execute(ctx, nil))
	}
	// Simulate the catalog state left by an interrupted concurrent index build.
	_, err = db.ExecContext(ctx, "UPDATE pg_index SET indisvalid = false WHERE indexrelid = 'eventstore.f_group_lifecycle_idx'::regclass")
	require.NoError(t, err)
	require.NoError(t, migration.Execute(ctx, nil))
	var valid bool
	require.NoError(t, db.QueryRowContext(ctx, "SELECT indisvalid FROM pg_index WHERE indexrelid = 'eventstore.f_group_lifecycle_idx'::regclass").Scan(&valid))
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
		{"instance", "live", "name", "renamed"},
		{"other-instance", "live", "name", "other"},
	}, got)
	_, err = db.ExecContext(ctx, `INSERT INTO eventstore.fields
 (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, should_index)
 SELECT 'instance', 'org', 'group', 'g-' || n, 'group', 'u-' || n, 1, 'member', to_jsonb('u-' || n), true
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
	t.Log(plan)

}
