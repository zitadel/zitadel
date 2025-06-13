package repository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/embedded"
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

var pool database.PoolTest

func runTests(m *testing.M) int {
	var stop func()
	var err error
	ctx := context.Background()
	pool, stop, err = newEmbeededDB(ctx)
	if err != nil {
		log.Print(err)
		return 1
	}
	defer stop()

	return m.Run()
}

func newEmbeededDB(ctx context.Context) (pool database.PoolTest, stop func(), err error) {
	var connector database.Connector
	connector, stop, err = embedded.StartEmbedded()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start embedded postgres: %w", err)
	}

	pool_, err := connector.Connect(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to embedded postgres: %w", err)
	}
	pool = pool_.(database.PoolTest)

	err = pool.MigrateTest(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to migrate database: %w", err)
	}
	return pool, stop, err
}
