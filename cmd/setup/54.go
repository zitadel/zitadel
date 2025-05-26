package setup

import (
	"context"
	_ "embed"

	v3_db "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type TransactionalTables struct {
	dbClient *database.DB
}

func (mig *TransactionalTables) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS zitadel")
	if err != nil {
		return err
	}

	config := &postgres.Config{Pool: mig.dbClient.Pool}
	pool, err := config.Connect(ctx)
	if err != nil {
		return err
	}

	return pool.(v3_db.Migrator).Migrate(ctx)
}

func (mig *TransactionalTables) String() string {
	return "54_repeatable_transactional_tables"
}

func (mig *TransactionalTables) Check(lastRun map[string]interface{}) bool {
	return true
}
