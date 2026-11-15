package query

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"strings"
	"testing"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database/postgres"
)

func TestGroupExternalIDSearchQuery(t *testing.T) {
	query, err := NewGroupExternalIDSearchQuery("source:entra", "CaseSensitive' OR true --")
	require.NoError(t, err)
	sql, args, err := query.toQuery(sq.Select("id").From("projections.groups1")).ToSql()
	require.NoError(t, err)
	require.Contains(t, sql, "projections.group_external_ids2.instance_id = projections.groups1.instance_id")
	require.Contains(t, sql, "projections.group_external_ids2.resource_owner = projections.groups1.resource_owner")
	require.Contains(t, sql, "projections.group_external_ids2.scope = ?")
	require.Contains(t, sql, "projections.group_external_ids2.external_id_sha256 = ?")
	require.Contains(t, sql, "projections.group_external_ids2.external_id = ?")
	require.NotContains(t, sql, "CaseSensitive")
	digest := sha256.Sum256([]byte("CaseSensitive' OR true --"))
	require.Equal(t, []interface{}{"source:entra", digest[:], "CaseSensitive' OR true --"}, args)
}

func TestGroupExternalIDSearchQuery_database(t *testing.T) {
	config, stop := postgres.StartEmbedded()
	defer stop()
	db, err := sql.Open("pgx", config.GetConnectionURL())
	require.NoError(t, err)
	defer db.Close()
	ctx := context.Background()
	_, err = db.ExecContext(ctx, `CREATE SCHEMA projections;
		CREATE TABLE projections.groups1 (id TEXT, instance_id TEXT, resource_owner TEXT);
		CREATE TABLE projections.group_external_ids2 (
			group_id TEXT, instance_id TEXT, resource_owner TEXT, scope TEXT,
			external_id TEXT, external_id_sha256 BYTEA
		)`)
	require.NoError(t, err)
	value := `CaseSensitive,"{}\雪' OR true --` + strings.Repeat("long-id", 1000)
	digest := sha256.Sum256([]byte(value))
	for _, fixture := range []struct {
		id, instance, org, scope, externalID string
	}{
		{"match", "instance", "org", "entra", value},
		{"other-instance", "other", "org", "entra", value},
		{"other-org", "instance", "other", "entra", value},
		{"other-source", "instance", "org", "okta", value},
		{"wrong-case", "instance", "org", "entra", strings.ToLower(value)},
		{"collision", "instance", "org", "entra", "different-value"},
	} {
		// Keep candidate groups in the requested scope to test the binding correlation.
		_, err = db.ExecContext(ctx, `INSERT INTO projections.groups1 VALUES ($1, $2, $3)`, fixture.id, "instance", "org")
		require.NoError(t, err)
		// Deliberately give every row the same digest to exercise the exact-value guard.
		_, err = db.ExecContext(ctx, `INSERT INTO projections.group_external_ids2 VALUES ($1, $2, $3, $4, $5, $6)`, fixture.id, fixture.instance, fixture.org, fixture.scope, fixture.externalID, digest[:])
		require.NoError(t, err)
	}
	query, err := NewGroupExternalIDSearchQuery("entra", value)
	require.NoError(t, err)
	statement, args, err := query.toQuery(sq.Select("id").From("projections.groups1").Where(sq.Eq{"instance_id": "instance", "resource_owner": "org"})).PlaceholderFormat(sq.Dollar).ToSql()
	require.NoError(t, err)
	rows, err := db.QueryContext(ctx, statement, args...)
	require.NoError(t, err)
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		require.NoError(t, rows.Scan(&id))
		ids = append(ids, id)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []string{"match"}, ids)
}
