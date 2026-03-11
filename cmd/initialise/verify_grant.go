package initialise

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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

// GrantPublicSchemaCreate grants CREATE ON SCHEMA public to the given user
// in the target database. This is needed because DuckLake's postgres extension
// creates catalog tables (ducklake_metadata, etc.) in the public schema, and
// PostgreSQL 15+ no longer grants CREATE ON SCHEMA public to all users.
//
// The admin connection (db) targets the admin database, but the grant must
// execute on the ZITADEL database where the public schema lives. A temporary
// pool is created by copying the admin credentials and switching the database.
func GrantPublicSchemaCreate(databaseName, username string) func(context.Context, *database.DB) error {
	return func(ctx context.Context, db *database.DB) error {
		if db.Pool == nil {
			return errors.New("grant public schema: admin pool not available")
		}

		// Copy the admin pool config and retarget to the ZITADEL database.
		cfg := db.Pool.Config().Copy()
		cfg.ConnConfig.Database = databaseName

		pool, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			return fmt.Errorf("grant public schema: connect to %s: %w", databaseName, err)
		}
		defer pool.Close()

		stmt := fmt.Sprintf(`GRANT CREATE ON SCHEMA public TO "%s"`, username)
		logging.Info(ctx, "verify public schema grant for DuckLake", "user", username, "database", databaseName)

		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("grant CREATE ON SCHEMA public to %s: %w", username, err)
		}
		return nil
	}
}
