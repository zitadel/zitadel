package projection

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/embedded"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_es "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	es_v3 "github.com/zitadel/zitadel/internal/eventstore/v3"
	// Register migrations
	_ "github.com/zitadel/zitadel/internal/query/projection/migrations"
)

var (
	pool database.PoolTest
	es   *eventstore.Eventstore
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	var stop func()
	var err error
	ctx := context.Background()
	pool, stop, err = newEmbeddedDB(ctx)
	if err != nil {
		log.Printf("error with embedded postgres database: %v", err)
		os.Exit(1)
	}
	defer stop()

	pusher := es_v3.NewEventstoreFromPool(pool)
	es = eventstore.NewEventstore(&eventstore.Config{
		Pusher:   pusher,
		Querier:  old_es.NewPostgres(pool.InternalDB()),
		Searcher: pusher,
	})

	CreateRelational(ctx, pool.InternalDB(), es, Config{})

	if err := Start(ctx); err != nil {
		log.Printf("error starting projections: %v", err)
		return 1
	}

	return m.Run()
}

func newEmbeddedDB(ctx context.Context) (pool database.PoolTest, stop func(), err error) {
	connector, stop, err := embedded.StartEmbedded()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start embedded postgres: %w", err)
	}

	dummyPool, err := connector.Connect(ctx, postgres.WithAfterConnectFunc(func(ctx context.Context, c *pgx.Conn) error {
		pgxdecimal.Register(c.TypeMap())
		return es_v3.RegisterEventstoreTypes(ctx, c)
	}))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to embedded postgres: %w", err)
	}

	pool = dummyPool.(database.PoolTest)
	err = pool.MigrateTest(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to migrate database: %w", err)
	}

	return pool, stop, err
}

// Copied from [integration.WaitForAndTickWithMaxDuration] to avoid an import loop
func waitForAndTickWithMaxDuration(ctx context.Context, max time.Duration) (time.Duration, time.Duration) {
	// interval which is used to retry the test
	tick := time.Second
	// tolerance which is used to stop the test for the timeout
	tolerance := tick * 5
	// default of the WaitFor is always a defined duration, shortened if the context would time out before
	waitFor := max

	if ctxDeadline, ok := ctx.Deadline(); ok {
		// if the context has a deadline, set the WaitFor to the shorter duration
		if until := time.Until(ctxDeadline); until < waitFor {
			// ignore durations which are smaller than the tolerance
			if until < tolerance {
				waitFor = 0
			} else {
				// always let the test stop with tolerance before the context is in timeout
				waitFor = until - tolerance
			}
		}
	}
	return waitFor, tick
}
