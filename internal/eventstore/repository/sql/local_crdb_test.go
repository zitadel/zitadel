package sql

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/cockroach"
)

var (
	testCRDBClient *sql.DB
)

func TestMain(m *testing.M) {
	opts := make([]testserver.TestServerOpt, 0, 1)
	if version := os.Getenv("ZITADEL_CRDB_VERSION"); version != "" {
		opts = append(opts, testserver.CustomVersionOpt(version))
	}
	ts, err := testserver.NewTestServer(opts...)
	if err != nil {
		logging.WithFields("error", err).Fatal("unable to start db")
	}

	testCRDBClient, err = sql.Open("postgres", ts.PGURL().String())
	if err != nil {
		logging.WithFields("error", err).Fatal("unable to connect to db")
	}
	if err = testCRDBClient.Ping(); err != nil {
		logging.WithFields("error", err).Fatal("unable to ping db")
	}

	defer func() {
		testCRDBClient.Close()
		ts.Stop()
	}()

	if err = initDB(&database.DB{DB: testCRDBClient, Database: &cockroach.Config{Database: "zitadel"}}); err != nil {
		logging.WithFields("error", err).Fatal("migrations failed")
	}

	os.Exit(m.Run())
}

func initDB(db *database.DB) error {
	config := new(database.Config)
	config.SetConnector(&cockroach.Config{User: cockroach.User{Username: "zitadel"}, Database: "zitadel"})

	if err := initialise.ReadStmts("cockroach"); err != nil {
		return err
	}

	err := initialise.Init(db,
		initialise.VerifyUser(config.Username(), ""),
		initialise.VerifyDatabase(config.DatabaseName()),
		initialise.VerifyGrant(config.DatabaseName(), config.Username()),
		initialise.VerifySettings(config.DatabaseName(), config.Username()))
	if err != nil {
		return err
	}

	err = initialise.VerifyZitadel(context.Background(), db, *config)
	if err != nil {
		return err
	}

	// create old events
	_, err = db.Exec(oldEventsTable)
	return err
}

type testDB struct{}

func (_ *testDB) Timetravel(time.Duration) string { return " AS OF SYSTEM TIME '-1 ms' " }

func (*testDB) DatabaseName() string { return "db" }

func (*testDB) Username() string { return "user" }

func (*testDB) Type() string { return "cockroach" }

const oldEventsTable = `CREATE TABLE IF NOT EXISTS eventstore.events (
	id UUID DEFAULT gen_random_uuid()
	, event_type TEXT NOT NULL
	, aggregate_type TEXT NOT NULL
	, aggregate_id TEXT NOT NULL
	, aggregate_version TEXT NOT NULL
	, event_sequence BIGINT NOT NULL
	, previous_aggregate_sequence BIGINT
	, previous_aggregate_type_sequence INT8
	, creation_date TIMESTAMPTZ NOT NULL DEFAULT now()
	, created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
	, event_data JSONB
	, editor_user TEXT NOT NULL 
	, editor_service TEXT
	, resource_owner TEXT NOT NULL
	, instance_id TEXT NOT NULL
	, "position" DECIMAL NOT NULL
	, in_tx_order INTEGER NOT NULL

	, PRIMARY KEY (instance_id, aggregate_type, aggregate_id, event_sequence DESC)
);`
