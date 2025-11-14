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

func newGrant() *cobra.Command {
	return &cobra.Command{
		Use:   "grant",
		Short: "set ALL grant to user",
		Long: `Sets ALL grant to the database user.

Prerequisites:
- postgreSQL
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
		var currentUser string
		err := db.QueryRowContext(ctx, func(r *sql.Row) error {
			return r.Scan(&currentUser)
		}, "SELECT current_user")
		if err != nil {
			return fmt.Errorf("unable to get current user: %w", err)
		}
		if currentUser == username {
			logging.WithFields("username", username).Info("config.database.postgres.user.username is same as config.database.postgres.admin.username, skipping grant")
			return nil
		}
		logging.WithFields("user", username, "database", databaseName).Info("verify grant")

		return exec(ctx, db, fmt.Sprintf(grantStmt, databaseName, username), nil)
	}
}
