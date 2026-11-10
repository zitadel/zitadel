package sql

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/database/postgres"
	new_es "github.com/zitadel/zitadel/internal/eventstore/v3"
)

var (
	testClient *sql.DB
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		config, cleanup := postgres.StartEmbedded()
		defer cleanup()

		connConfig, err := pgxpool.ParseConfig(config.GetConnectionURL())
		logging.OnError(err).Fatal("unable to parse db url")

		connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			pgxdecimal.Register(conn.TypeMap())
			return new_es.RegisterEventstoreTypes(ctx, conn)
		}

		pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
		logging.OnError(err).Fatal("unable to create db pool")

		testClient = stdlib.OpenDBFromPool(pool)

		err = testClient.Ping()
		logging.OnError(err).Fatal("unable to ping db")

		defer func() {
			logging.OnError(testClient.Close()).Error("unable to close db")
		}()

		err = initDB(context.Background(), &database.DB{DB: testClient, Database: &postgres.Config{Database: "zitadel"}})
		logging.OnError(err).Fatal("migrations failed")

		return m.Run()
	}())
}

func initDB(ctx context.Context, db *database.DB) error {
	config := new(database.Config)
	config.SetConnector(&postgres.Config{User: postgres.User{Username: "zitadel"}, Database: "zitadel"})

	if err := initialise.ReadStmts(); err != nil {
		return err
	}

	err := initialise.Init(ctx, db,
		initialise.VerifyUser(config.Username(), ""),
		initialise.VerifyDatabase(config.DatabaseName()),
		initialise.VerifyGrant(config.DatabaseName(), config.Username()))
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

func (*testDB) Type() dialect.DatabaseType { return dialect.DatabaseTypePostgres }

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

	, PRIMARY KEY (instance_id, aggregate_type, aggregate_id, event_sequence)
);`
