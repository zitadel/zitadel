package initialise

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"

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
			config, shutdown, err := NewConfig(cmd, viper.GetViper())
			if err != nil {
				return err
			}
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

		// Check if the role already exists in the catalog before attempting CREATE USER.
		// This avoids PostgreSQL logging the full CREATE USER statement (including the
		// plaintext password) when the role already exists on re-deploy / re-init.
		var exists bool
		err = db.QueryRowContext(ctx, func(r *sql.Row) error {
			return r.Scan(&exists)
		}, "SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = $1)", username)
		if err != nil {
			return fmt.Errorf("unable to check if user exists: %w", err)
		}
		if exists {
			logging.Info(ctx, "user already exists, skipping creation", "username", username)
			return nil
		}

		logging.Info(ctx, "verify user", "username", username)
		// Format the username first so the password is never part of a fmt format string
		// (passwords may contain '%' which would be interpreted as format verbs).
		stmt := fmt.Sprintf(createUserStmt, username)
		if password != "" {
			stmt += " WITH PASSWORD " + quotePostgresLiteral(password)
		}

		return exec(ctx, db, stmt, []string{roleAlreadyExistsCode})
	}
}

// quotePostgresLiteral returns s as a single-quoted PostgreSQL string literal,
// escaping embedded single quotes by doubling them.
func quotePostgresLiteral(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
