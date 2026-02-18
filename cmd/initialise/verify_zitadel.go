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
	es_v3 "github.com/zitadel/zitadel/internal/eventstore/v3"
)

func newZitadel() *cobra.Command {
	return &cobra.Command{
		Use:     "schema",
		Aliases: []string{"zitadel"},
		Short:   "bootstrap the ZITADEL database schema",
		Long: `Bootstrap the ZITADEL database schema.

Creates all required schemas (eventstore, projections, system) and base tables
using the service user credentials. No admin/superuser privileges are required.

Use this command when you have provisioned the database user and database
yourself (e.g. on a managed PostgreSQL service) and want to skip the
admin-credential requirement of the full 'zitadel init' command.

Prerequisites:
- PostgreSQL user exists and has ownership of the target database
`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				logging.OnError(cmd.Context(), err).Error("zitadel verify zitadel command failed")
			}()
			config, shutdown, err := NewConfig(cmd, viper.GetViper())
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, shutdown(cmd.Context()))
			}()
			return verifyZitadel(cmd.Context(), config.Database)
		},
	}
}

func VerifyZitadel(ctx context.Context, db *database.DB, config database.Config) error {
	err := ReadStmts()
	if err != nil {
		return err
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	logging.Info(ctx, "verify system")
	if err := exec(ctx, conn, fmt.Sprintf(createSystemStmt, config.Username()), nil); err != nil {
		return err
	}

	logging.Info(ctx, "verify encryption keys")
	if err := createEncryptionKeys(ctx, conn); err != nil {
		return err
	}

	logging.Info(ctx, "verify projections")
	if err := exec(ctx, conn, fmt.Sprintf(createProjectionsStmt, config.Username()), nil); err != nil {
		return err
	}

	logging.Info(ctx, "verify eventstore")
	if err := exec(ctx, conn, fmt.Sprintf(createEventstoreStmt, config.Username()), nil); err != nil {
		return err
	}

	logging.Info(ctx, "verify events tables")
	if err := createEvents(ctx, conn); err != nil {
		return err
	}

	logging.Info(ctx, "verify unique constraints")
	if err := exec(ctx, conn, createUniqueConstraints, nil); err != nil {
		return err
	}

	return nil
}

func verifyZitadel(ctx context.Context, config database.Config) error {
	logging.Info(ctx, "verify zitadel", "database", config.DatabaseName())

	db, err := database.Connect(config, false)
	if err != nil {
		return err
	}

	if err := VerifyZitadel(ctx, db, config); err != nil {
		return err
	}

	return db.Close()
}

func createEncryptionKeys(ctx context.Context, db database.Beginner) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err = tx.Exec(createEncryptionKeysStmt); err != nil {
		rollbackErr := tx.Rollback()
		logging.WithError(ctx, rollbackErr).Error("rollback failed")
		return err
	}

	return tx.Commit()
}

func createEvents(ctx context.Context, conn *sql.Conn) (err error) {
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			logging.WithError(ctx, rollbackErr).Error("rollback failed")
			return
		}
		err = tx.Commit()
	}()

	// if events already exists events2 is created during a setup job
	var count int
	row := tx.QueryRow("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'eventstore' AND table_name like 'events%'")
	if err = row.Scan(&count); err != nil {
		return err
	}
	if row.Err() != nil || count >= 1 {
		return row.Err()
	}
	_, err = tx.Exec(createEventsStmt)
	if err != nil {
		return err
	}
	return es_v3.CheckExecutionPlan(ctx, conn)
}
