package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type TransactionalTables struct {
	dbClient *database.DB

	ShouldRecreateSchema bool
}

func (mig *TransactionalTables) Execute(ctx context.Context, _ eventstore.Event) error {
	// TODO(adlerhurst): revert changes made in https://github.com/zitadel/zitadel/pull/11833 before v5 release.
	if mig.ShouldRecreateSchema {
		logging.Info(ctx, "dropping schema of relational tables")
		if err := mig.dropSchema(ctx); err != nil {
			return err
		}
	}

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

func (mig *TransactionalTables) dropSchema(ctx context.Context) (err error) {
	tx, err := mig.dbClient.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			logging.OnError(ctx, rollbackErr).Debug("rollback failed")
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.ExecContext(ctx, "DROP SCHEMA IF EXISTS zitadel CASCADE")
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM projections.current_states WHERE projection_name LIKE $1 || '%' OR projection_name = $2", "zitadel.", "relational_tables")
	return err
}
