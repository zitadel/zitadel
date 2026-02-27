package projection

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/embedded"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
)

var pool database.PoolTest

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
	return m.Run()
}

func newEmbeddedDB(ctx context.Context) (pool database.PoolTest, stop func(), err error) {
	connector, stop, err := embedded.StartEmbedded()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start embedded postgres: %w", err)
	}

	dummyPool, err := connector.Connect(ctx)
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

func getTransactions(t *testing.T, pool database.PoolTest) (rawTx *sql.Tx, v3SQLTx *v3_sql.Transaction) {
	rawTx, err := pool.RawDB().Begin()
	require.NoError(t, err)
	v3SQLTx = v3_sql.SQLTx(rawTx)
	return
}
