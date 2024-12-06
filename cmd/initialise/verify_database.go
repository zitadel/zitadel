package initialise

import (
	"context"
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
- cockroachDB or postgreSQL

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

		return exec(ctx, db, fmt.Sprintf(databaseStmt, databaseName), []string{dbAlreadyExistsCode})
	}
}
