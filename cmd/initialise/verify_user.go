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

func newUser() *cobra.Command {
	return &cobra.Command{
		Use:   "user",
		Short: "initialize only the database user",
		Long: `Sets up the ZITADEL database user.

Prerequisites:
- postgreSQL

The user provided by flags needs privileges to 
- create the database if it does not exist
- see other users and create a new one if the user does not exist
- grant all rights of the ZITADEL database to the user created if not yet set
`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				logging.OnError(cmd.Context(), err).Error("zitadel verify user command failed")
			}()
			config, shutdown, err := NewConfig(cmd.Context(), viper.GetViper())
			if err != nil {
				return err
			}
			// Set logger again to include changes from config
			cmd.SetContext(logging.NewCtx(cmd.Context(), logging.StreamRuntime))
			defer func() {
				err = errors.Join(err, shutdown(cmd.Context()))
			}()

			return initialise(cmd.Context(), config.Database, VerifyUser(config.Database.Username(), config.Database.Password()))
		},
	}
}

func VerifyUser(username, password string) func(context.Context, *database.DB) error {
	return func(ctx context.Context, db *database.DB) error {
		var currentUser string
		err := db.QueryRowContext(ctx, func(r *sql.Row) error {
			return r.Scan(&currentUser)
		}, "SELECT current_user")
		if err != nil {
			return fmt.Errorf("unable to get current user: %w", err)
		}
		if currentUser == username {
			logging.Info(ctx, "config.database.postgres.user.username is same as config.database.postgres.admin.username, skipping create user", "username", username)
			return nil
		}
		logging.Info(ctx, "verify user", "username", username)
		if password != "" {
			createUserStmt += " WITH PASSWORD '" + password + "'"
		}

		return exec(ctx, db, fmt.Sprintf(createUserStmt, username), []string{roleAlreadyExistsCode})
	}
}
