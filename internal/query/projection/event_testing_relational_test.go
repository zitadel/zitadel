package projection

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/embedded"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
)

var (
	pool  database.PoolTest
	rawDB *sql.DB
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
		return 1
	}
	defer stop()

	rawDB = pool.RawDB()
	defer rawDB.Close()
	defer pool.Close(ctx)

	return m.Run()
}

func newEmbeddedDB(ctx context.Context) (pool database.PoolTest, stop func(), err error) {
	connector, stop, err := embedded.StartEmbedded()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start embedded postgres: %w", err)
	}

	dummyPool, err := connector.Connect(ctx)
	if err != nil {
		return nil, stop, fmt.Errorf("unable to connect to embedded postgres: %w", err)
	}

	pool, ok := dummyPool.(database.PoolTest)
	if !ok {
		return nil, stop, fmt.Errorf("expecting database.PoolTest, got %T", dummyPool)
	}

	err = pool.MigrateTest(ctx)
	if err != nil {
		return nil, stop, fmt.Errorf("unable to migrate database: %w", err)
	}

	return pool, stop, err
}

func getTransactions(t *testing.T) (rawTx *sql.Tx, v3SQLTx *v3_sql.Transaction) {
	rawTx, err := rawDB.Begin()
	require.NoError(t, err)
	v3SQLTx = v3_sql.SQLTx(rawTx)
	return
}

// callReduce on the input projection tests that:
//
//  1. the aggregate type of the input event exists
//  2. the event type of the input event exists for that aggregate type
//  3. the reducer function for the input event is not nil
//  4. the reducer function for the input event executes without errors
//
// Returns false if any of the above checks fails
func callReduce(t *testing.T, ctx context.Context, tx *sql.Tx, projection handler.Projection, event eventstore.Event) bool {
	reducers := projection.Reducers()
	aggregateReducerIdx := slices.IndexFunc(reducers, func(r handler.AggregateReducer) bool {
		return r.Aggregate == event.Aggregate().Type
	})
	if !assert.Greater(t, aggregateReducerIdx, -1) {
		return false
	}

	eventReducerIdx := slices.IndexFunc(reducers[aggregateReducerIdx].EventReducers, func(er handler.EventReducer) bool {
		return er.Event == event.Type()
	})
	if !assert.Greater(t, eventReducerIdx, -1) {
		return false
	}

	reduceFn := reducers[aggregateReducerIdx].EventReducers[eventReducerIdx].Reduce
	if !assert.NotNil(t, reduceFn) {
		return false
	}

	stmt, err := reduceFn(event)
	if !assert.NoError(t, err) {
		return false
	}

	err = stmt.Execute(ctx, tx, "")
	return assert.NoError(t, err)
}
