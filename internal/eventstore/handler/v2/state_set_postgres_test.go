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
	_, err = db.Exec(`
INSERT INTO projections.current_states (
    projection_name, instance_id, aggregate_id, aggregate_type, "sequence",
    event_date, "position", last_updated, filter_offset
) VALUES (
    'users14', 'inst', 'agg', 'user', 1, now(), 1.0, now(), 3
)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO eventstore.events2 VALUES ('inst', 'agg', 'user', 1, 9)`)
	require.NoError(t, err)

	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	stepSQL, err := os.ReadFile(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..", "cmd", "setup", "82.sql"))
	require.NoError(t, err)
	_, err = db.Exec(string(stepSQL))
	require.NoError(t, err)
	assertInTxOrder(t, db, 9)

	_, err = db.Exec(`
INSERT INTO projections.current_states (
    projection_name, instance_id, aggregate_id, aggregate_type, "sequence",
    event_date, "position", last_updated, filter_offset, in_tx_order
) VALUES (
    'users15', 'inst2', 'agg', 'user', 1, now(), 1.0, now(), 3, 7
)`)
	require.NoError(t, err)
	var inserted sql.NullInt64
	err = db.QueryRow(`SELECT in_tx_order FROM projections.current_states WHERE projection_name = $1 AND instance_id = $2`, "users15", "inst2").Scan(&inserted)
	require.NoError(t, err)
	require.True(t, inserted.Valid)
	assert.Equal(t, int64(7), inserted.Int64)

	oldStyleUpdate := `
UPDATE projections.current_states SET
    aggregate_id = $3,
    aggregate_type = $4,
    "sequence" = $5,
    filter_offset = $6
WHERE projection_name = $1 AND instance_id = $2`
	_, err = db.Exec(oldStyleUpdate, "users14", "inst", "agg", "user", 2, 4)
	require.NoError(t, err)
	assertInTxOrderNull(t, db)

	tx, err := db.BeginTx(context.Background(), nil)
	require.NoError(t, err)
	defer func() { _ = tx.Rollback() }()

	res, err := tx.Exec(updateStateStmt,
		"users14",
		"inst",
		"agg",
		"user",
		uint64(3),
		sql.NullTime{},
		1.0,
		uint32(5),
		uint32(11),
	)
	require.NoError(t, err)
	affected, err := res.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), affected)
	require.NoError(t, tx.Commit())
	assertInTxOrder(t, db, 11)
}

func assertInTxOrder(t *testing.T, db *sql.DB, want int64) {
	t.Helper()
	var inTxOrder sql.NullInt64
	err := db.QueryRow(`SELECT in_tx_order FROM projections.current_states WHERE projection_name = $1 AND instance_id = $2`, "users14", "inst").Scan(&inTxOrder)
	require.NoError(t, err)
	require.True(t, inTxOrder.Valid)
	assert.Equal(t, want, inTxOrder.Int64)
}

func assertInTxOrderNull(t *testing.T, db *sql.DB) {
	t.Helper()
	var inTxOrder sql.NullInt64
	err := db.QueryRow(`SELECT in_tx_order FROM projections.current_states WHERE projection_name = $1 AND instance_id = $2`, "users14", "inst").Scan(&inTxOrder)
	require.NoError(t, err)
	assert.False(t, inTxOrder.Valid, "old-style update without in_tx_order in SET must leave in_tx_order null")
}
