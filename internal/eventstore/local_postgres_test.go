package eventstore_test

import (
	"context"
	"encoding/json"
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
	"github.com/zitadel/zitadel/internal/eventstore"
	es_sql "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	new_es "github.com/zitadel/zitadel/internal/eventstore/v3"
)

var (
	testClient *database.DB
	queriers   map[string]eventstore.Querier = make(map[string]eventstore.Querier)
	pushers    map[string]eventstore.Pusher  = make(map[string]eventstore.Pusher)
	clients    map[string]*database.DB       = make(map[string]*database.DB)
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		config, cleanup := postgres.StartEmbedded()
		defer cleanup()

		testClient = &database.DB{
			Database: new(testDB),
		}

		connConfig, err := pgxpool.ParseConfig(config.GetConnectionURL())
		logging.OnError(err).Fatal("unable to parse db url")

		connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			pgxdecimal.Register(conn.TypeMap())
			return new_es.RegisterEventstoreTypes(ctx, conn)
		}
		pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
		logging.OnError(err).Fatal("unable to create db pool")

		testClient.DB = stdlib.OpenDBFromPool(pool)
		err = testClient.Ping()
		logging.OnError(err).Fatal("unable to ping db")

		v2 := &es_sql.Postgres{DB: testClient}
		queriers["v2(inmemory)"] = v2
		clients["v2(inmemory)"] = testClient

		pushers["v3(inmemory)"] = new_es.NewEventstore(testClient)
		clients["v3(inmemory)"] = testClient

		if localDB, err := connectLocalhost(); err == nil {
			err = initDB(context.Background(), localDB)
			logging.OnError(err).Fatal("migrations failed")

			pushers["v3(singlenode)"] = new_es.NewEventstore(localDB)
			clients["v3(singlenode)"] = localDB
		}

		defer func() {
			logging.OnError(testClient.Close()).Error("unable to close db")
		}()

		err = initDB(context.Background(), &database.DB{DB: testClient.DB, Database: &postgres.Config{Database: "zitadel"}})
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

	err = initialise.VerifyZitadel(ctx, db, *config)
	if err != nil {
		return err
	}

	// create old events
	_, err = db.Exec(oldEventsTable)
	return err
}

func connectLocalhost() (*database.DB, error) {
	config, err := pgxpool.ParseConfig("postgresql://postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return nil, err
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return new_es.RegisterEventstoreTypes(ctx, conn)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	client := stdlib.OpenDBFromPool(pool)
	if err = client.Ping(); err != nil {
		return nil, err
	}

	return &database.DB{
		DB:       client,
		Database: new(testDB),
	}, nil
}

type testDB struct{}

func (_ *testDB) Timetravel(time.Duration) string { return " AS OF SYSTEM TIME '-1 ms' " }

func (*testDB) DatabaseName() string { return "db" }

func (*testDB) Username() string { return "user" }

func (*testDB) Type() dialect.DatabaseType { return dialect.DatabaseTypePostgres }

func generateCommand(aggregateType eventstore.AggregateType, aggregateID string, opts ...func(*testEvent)) eventstore.Command {
	e := &testEvent{
		BaseEvent: eventstore.BaseEvent{
			Agg: &eventstore.Aggregate{
				ID:            aggregateID,
				Type:          aggregateType,
				ResourceOwner: "ro",
				Version:       "v1",
			},
			Service:   "svc",
			EventType: "test.created",
		},
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

type testEvent struct {
	eventstore.BaseEvent
	uniqueConstraints []*eventstore.UniqueConstraint
}

func (e *testEvent) Payload() any {
	return e.BaseEvent.Data
}

func (e *testEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return e.uniqueConstraints
}

func canceledCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func fillUniqueData(unique_type, field, instanceID string) error {
	_, err := testClient.Exec("INSERT INTO eventstore.unique_constraints (unique_type, unique_field, instance_id) VALUES ($1, $2, $3)", unique_type, field, instanceID)
	return err
}

func generateAddUniqueConstraint(table, uniqueField string) func(e *testEvent) {
	return func(e *testEvent) {
		e.uniqueConstraints = append(e.uniqueConstraints,
			&eventstore.UniqueConstraint{
				UniqueType:  table,
				UniqueField: uniqueField,
				Action:      eventstore.UniqueConstraintAdd,
			},
		)
	}
}

func generateRemoveUniqueConstraint(table, uniqueField string) func(e *testEvent) {
	return func(e *testEvent) {
		e.uniqueConstraints = append(e.uniqueConstraints,
			&eventstore.UniqueConstraint{
				UniqueType:  table,
				UniqueField: uniqueField,
				Action:      eventstore.UniqueConstraintRemove,
			},
		)
	}
}

func withTestData(data any) func(e *testEvent) {
	return func(e *testEvent) {
		d, err := json.Marshal(data)
		if err != nil {
			panic("marshal data failed")
		}
		e.BaseEvent.Data = d
	}
}

func cleanupEventstore(client *database.DB) func() {
	return func() {
		_, err := client.Exec("TRUNCATE eventstore.events")
		if err != nil {
			logging.Warnf("unable to truncate events: %v", err)
		}
		_, err = client.Exec("TRUNCATE eventstore.events2")
		if err != nil {
			logging.Warnf("unable to truncate events: %v", err)
		}
		_, err = client.Exec("TRUNCATE eventstore.unique_constraints")
		if err != nil {
			logging.Warnf("unable to truncate unique constraints: %v", err)
		}
	}
}

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
