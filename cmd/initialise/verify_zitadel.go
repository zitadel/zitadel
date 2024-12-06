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
	"github.com/zitadel/zitadel/internal/database/dialect"
	es_v3 "github.com/zitadel/zitadel/internal/eventstore/v3"
)

func newZitadel() *cobra.Command {
	return &cobra.Command{
		Use:   "zitadel",
		Short: "initialize ZITADEL internals",
		Long: `initialize ZITADEL internals.

Prerequisites:
- cockroachDB or postgreSQL with user and database
`,
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.GetViper())
			err := verifyZitadel(cmd.Context(), config.Database)
			logging.OnError(err).Fatal("unable to init zitadel")
		},
	}
}

func VerifyZitadel(ctx context.Context, db *database.DB, config database.Config) error {
	err := ReadStmts(config.Type())
	if err != nil {
		return err
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	logging.WithFields().Info("verify system")
	if err := exec(ctx, conn, fmt.Sprintf(createSystemStmt, config.Username()), nil); err != nil {
		return err
	}

	logging.WithFields().Info("verify encryption keys")
	if err := createEncryptionKeys(ctx, conn); err != nil {
		return err
	}

	logging.WithFields().Info("verify projections")
	if err := exec(ctx, conn, fmt.Sprintf(createProjectionsStmt, config.Username()), nil); err != nil {
		return err
	}

	logging.WithFields().Info("verify eventstore")
	if err := exec(ctx, conn, fmt.Sprintf(createEventstoreStmt, config.Username()), nil); err != nil {
		return err
	}

	logging.WithFields().Info("verify events tables")
	if err := createEvents(ctx, conn); err != nil {
		return err
	}

	logging.WithFields().Info("verify system sequence")
	if err := exec(ctx, conn, createSystemSequenceStmt, nil); err != nil {
		return err
	}

	logging.WithFields().Info("verify unique constraints")
	if err := exec(ctx, conn, createUniqueConstraints, nil); err != nil {
		return err
	}

	return nil
}

func verifyZitadel(ctx context.Context, config database.Config) error {
	logging.WithFields("database", config.DatabaseName()).Info("verify zitadel")

	db, err := database.Connect(config, false, dialect.DBPurposeQuery)
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
		logging.OnError(rollbackErr).Error("rollback failed")
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
			logging.OnError(rollbackErr).Error("rollback failed")
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
