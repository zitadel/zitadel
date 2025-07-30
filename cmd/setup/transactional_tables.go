package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type TransactionalTables struct {
	dbClient *database.DB
}

func (mig *TransactionalTables) Execute(ctx context.Context, _ eventstore.Event) error {
	config := &postgres.Config{Pool: mig.dbClient.Pool}
	pool, err := config.Connect(ctx)
	if err != nil {
		return err
	}

	return pool.Migrate(ctx)
}

func (mig *TransactionalTables) String() string {
	return "repeatable_transactional_tables"
}

func (mig *TransactionalTables) Check(lastRun map[string]interface{}) bool {
	return true
}
