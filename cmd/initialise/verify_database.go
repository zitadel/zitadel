package initialise

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
)

func newDatabase() *cobra.Command {
	return &cobra.Command{
		Use:   "database",
		Short: "initialize only the database",
		Long: `Sets up the ZITADEL database.

Prerequisites:
- postgreSQL

The user provided by flags needs privileges to 
- create the database if it does not exist
- see other users and create a new one if the user does not exist
- grant all rights of the ZITADEL database to the user created if not yet set
`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				logging.OnError(cmd.Context(), err).Error("zitadel init verify database command failed")
			}()
			config, shutdown, err := NewConfig(cmd, viper.GetViper())
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, shutdown(cmd.Context()))
			}()

			return initialise(cmd.Context(), config.Database, VerifyDatabase(config.Database.DatabaseName()))
		},
	}
}

func VerifyDatabase(databaseName string) func(context.Context, *database.DB) error {
	return func(ctx context.Context, db *database.DB) error {
		var currentDatabase string
		err := db.QueryRowContext(ctx, func(r *sql.Row) error {
			return r.Scan(&currentDatabase)
		}, "SELECT current_database()")
		if err != nil {
			return fmt.Errorf("unable to get current database: %w", err)
		}
		if currentDatabase == databaseName {
			logging.Info(ctx, "database is same as config.database.postgres.admin.ExistingDatabase, skipping creation", "database", databaseName)
			return nil
		}

		// Check if the database already exists in the catalog before attempting CREATE DATABASE.
		// This handles the case where the database was provisioned externally and the admin
		// credentials are the same as the service user, which lacks the CREATEDB privilege.
		var exists bool
		err = db.QueryRowContext(ctx, func(r *sql.Row) error {
			return r.Scan(&exists)
		}, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", databaseName)
		if err != nil {
			return fmt.Errorf("unable to check if database exists: %w", err)
		}
		if exists {
			logging.Info(ctx, "database already exists, skipping creation", "database", databaseName)
			return nil
		}

		logging.Info(ctx, "verify database", "database", databaseName)

		return exec(ctx, db, fmt.Sprintf(databaseStmt, databaseName), []string{dbAlreadyExistsCode})
	}
}
