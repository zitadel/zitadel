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

func newGrant() *cobra.Command {
	return &cobra.Command{
		Use:   "grant",
		Short: "set ALL grant to user",
		Long: `Sets ALL grant to the database user.

Prerequisites:
- postgreSQL
`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				logging.OnError(cmd.Context(), err).Error("zitadel verify grant command failed")
			}()
			config, shutdown, err := NewConfig(cmd, viper.GetViper())
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, shutdown(cmd.Context()))
			}()

			return initialise(cmd.Context(), config.Database, VerifyGrant(config.Database.DatabaseName(), config.Database.Username()))
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
			logging.Info(ctx, "config.database.postgres.user.username is same as config.database.postgres.admin.username, skipping grant", "username", username)
			return nil
		}
		logging.Info(ctx, "verify grant", "user", username, "database", databaseName)

		return exec(ctx, db, fmt.Sprintf(grantStmt, databaseName, username), nil)
	}
}
