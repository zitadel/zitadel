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

func newSettings() *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "Ensures proper settings on the database",
		Long: `Ensures proper settings on the database.

Prerequisites:
- cockroachDB or postgreSQL

Cockroach
- Sets enable_durable_locking_for_serializable to on for the zitadel user and database
`,
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.GetViper())

			err := initialise(cmd.Context(), config.Database, VerifySettings(config.Database.DatabaseName(), config.Database.Username()))
			logging.OnError(err).Fatal("unable to set settings")
		},
	}
}

func VerifySettings(databaseName, username string) func(context.Context, *database.DB) error {
	return func(ctx context.Context, db *database.DB) error {
		if db.Type() == "postgres" {
			return nil
		}
		logging.WithFields("user", username, "database", databaseName).Info("verify settings")

		return exec(ctx, db, fmt.Sprintf(settingsStmt, databaseName, username), nil)
	}
}
