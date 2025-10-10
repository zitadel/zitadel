package initialise

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

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
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.GetViper())

			err := initialise(cmd.Context(), config.Database, VerifyDatabase(config.Database.DatabaseName()))
			logging.OnError(err).Fatal("unable to initialize the database")
		},
	}
}

func VerifyDatabase(databaseName string) func(context.Context, *database.DB) error {
	return func(ctx context.Context, db *database.DB) error {
		logging.WithFields("database", databaseName).Info("verify database")

		// Check if database already exists first
		exists, err := databaseExists(ctx, db, databaseName)
		if err != nil {
			return fmt.Errorf("failed to check if database exists: %w", err)
		}

		if exists {
			logging.WithFields("database", databaseName).Info("database already exists, skipping creation")
			return nil
		}

		// Proceed with database creation
		return exec(ctx, db, fmt.Sprintf(databaseStmt, databaseName), []string{dbAlreadyExistsCode})
	}
}

func databaseExists(ctx context.Context, db *database.DB, databaseName string) (bool, error) {
	var exists int
	err := db.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&exists)
	}, `SELECT 1 FROM pg_database WHERE datname = $1`, databaseName)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	
	return true, nil
}
