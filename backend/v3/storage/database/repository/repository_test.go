package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/embedded"
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

var pool database.Pool

func runTests(m *testing.M) int {
	connector, stop, err := embedded.StartEmbedded()
	if err != nil {
		log.Fatalf("unable to start embedded postgres: %v", err)
	}
	defer stop()

	ctx := context.Background()

	pool, err = connector.Connect(ctx)
	if err != nil {
		log.Printf("unable to connect to embedded postgres: %v", err)
		return 1
	}

	err = pool.Migrate(ctx)
	if err != nil {
		log.Printf("unable to migrate database: %v", err)
		return 1
	}

	return m.Run()
}
