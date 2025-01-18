package initialise

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

func newGrant() *cobra.Command {
	return &cobra.Command{
		Use:   "grant",
		Short: "set ALL grant to user",
		Long: `Sets ALL grant to the database user.

Prerequisites:
- cockroachDB or postgreSQL
`,
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.GetViper())

			err := initialise(cmd.Context(), config.Database, VerifyGrant(config.Database.DatabaseName(), config.Database.Username()))
			logging.OnError(err).Fatal("unable to set grant")
		},
	}
}

func VerifyGrant(databaseName, username string) func(context.Context, *database.DB) error {
	return func(ctx context.Context, db *database.DB) error {
		logging.WithFields("user", username, "database", databaseName).Info("verify grant")

		// the connection string is used to grant the user access to the public schema
		config := db.Pool.Config().ConnConfig
		config.Database = databaseName
		conn := strings.Replace(config.ConnString(), "dbname=postgres", fmt.Sprintf("dbname=%s", databaseName), 1)

		return exec(ctx, db, fmt.Sprintf(grantStmt, databaseName, username, conn), nil)
	}
}
