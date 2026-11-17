package mirror

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	db "github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/postgres"
)

func TestCopyUniqueConstraints(t *testing.T) {
	config, stop := postgres.StartEmbedded()
	defer stop()
	ctx := context.Background()
	open := func(name string) *db.DB {
		client, err := sql.Open("pgx", strings.TrimSuffix(config.GetConnectionURL(), "postgres")+name)
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })
		return &db.DB{DB: client}
	}
	dest := open("postgres")
	for _, name := range []string{"current", "legacy"} {
		_, err := dest.ExecContext(ctx, "CREATE DATABASE "+name)
		require.NoError(t, err)
	}
	const table = "CREATE SCHEMA eventstore; CREATE TABLE eventstore.unique_constraints (instance_id TEXT, unique_type TEXT, unique_field TEXT, owners TEXT[] NOT NULL DEFAULT '{}', PRIMARY KEY (instance_id, unique_type, unique_field));"
	_, err := dest.ExecContext(ctx, table)
	require.NoError(t, err)
	current := open("current")
	_, err = current.ExecContext(ctx, table+"INSERT INTO eventstore.unique_constraints VALUES ('i1', 'group_external_id', 'key', '{org:o1,group:g1}'), ('i1', 'username', 'u', '{}')")
	require.NoError(t, err)
	// Releases before setup step 78 have no owners column.
	legacy := open("legacy")
	_, err = legacy.ExecContext(ctx, `CREATE SCHEMA eventstore; CREATE TABLE eventstore.unique_constraints (instance_id TEXT, unique_type TEXT, unique_field TEXT);
 INSERT INTO eventstore.unique_constraints VALUES ('i2', 'username', 'u')`)
	require.NoError(t, err)

	isSystem = true
	defer func() { isSystem = false }()
	copyUniqueConstraints(ctx, current, dest)
	copyUniqueConstraints(ctx, legacy, dest)

	rows, err := dest.DB.QueryContext(ctx, "SELECT instance_id, unique_type, unique_field, array_to_string(owners, ',') FROM eventstore.unique_constraints ORDER BY 1, 2")
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
		{"i1", "group_external_id", "key", "org:o1,group:g1"},
		{"i1", "username", "u", ""},
		{"i2", "username", "u", ""},
	}, got)
}
