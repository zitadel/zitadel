package handler

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database/postgres"
)

func TestStateSetTrigger_optInKeepsInTxOrder(t *testing.T) {
	config, cleanup := postgres.StartEmbedded()
	t.Cleanup(cleanup)

	db, err := sql.Open("pgx", config.GetConnectionURL())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.Exec(`CREATE SCHEMA IF NOT EXISTS projections`)
	require.NoError(t, err)
	_, err = db.Exec(`CREATE SCHEMA IF NOT EXISTS eventstore`)
	require.NoError(t, err)
	_, err = db.Exec(`
CREATE TABLE eventstore.events2 (
    instance_id TEXT,
    aggregate_id TEXT,
    aggregate_type TEXT,
    "sequence" INT8,
    in_tx_order INTEGER
)`)
	require.NoError(t, err)
	_, err = db.Exec(`
CREATE TABLE projections.current_states (
    projection_name TEXT NOT NULL,
    instance_id TEXT NOT NULL,
    last_updated TIMESTAMPTZ,
    aggregate_id TEXT,
    aggregate_type TEXT,
    "sequence" INT8,
    event_date TIMESTAMPTZ,
    "position" DECIMAL,
    filter_offset INTEGER,
    PRIMARY KEY (projection_name, instance_id)
)`)
	require.NoError(t, err)

	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	stepSQL, err := os.ReadFile(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..", "cmd", "setup", "82.sql"))
	require.NoError(t, err)
	_, err = db.Exec(string(stepSQL))
	require.NoError(t, err)

	oldStyleInsert := `
INSERT INTO projections.current_states (
    projection_name, instance_id, aggregate_id, aggregate_type, "sequence",
    event_date, "position", last_updated, filter_offset, in_tx_order
) VALUES (
    $1, $2, $3, $4, $5, now(), $6, now(), $7, $8
)`
	_, err = db.Exec(oldStyleInsert, "users14", "inst", "agg", "user", 1, 1.0, 3, 9)
	require.NoError(t, err)

	var inTxOrder sql.NullInt64
	err = db.QueryRow(`SELECT in_tx_order FROM projections.current_states WHERE projection_name = $1 AND instance_id = $2`, "users14", "inst").Scan(&inTxOrder)
	require.NoError(t, err)
	assert.False(t, inTxOrder.Valid, "old-style write without set_config must leave in_tx_order null")

	tx, err := db.BeginTx(context.Background(), nil)
	require.NoError(t, err)
	defer func() { _ = tx.Rollback() }()

	res, err := tx.Exec(updateStateStmt,
		"users14",
		"inst",
		"agg",
		"user",
		uint64(2),
		sql.NullTime{},
		1.0,
		uint32(4),
		uint32(11),
	)
	require.NoError(t, err)
	affected, err := res.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), affected)
	require.NoError(t, tx.Commit())

	err = db.QueryRow(`SELECT in_tx_order FROM projections.current_states WHERE projection_name = $1 AND instance_id = $2`, "users14", "inst").Scan(&inTxOrder)
	require.NoError(t, err)
	require.True(t, inTxOrder.Valid)
	assert.Equal(t, int64(11), inTxOrder.Int64)
}
