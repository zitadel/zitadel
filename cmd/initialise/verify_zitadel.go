package initialise

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/database"
)

func newZitadel() *cobra.Command {
	return &cobra.Command{
		Use:   "zitadel",
		Short: "initialize ZITADEL internals",
		Long: `initialize ZITADEL internals.

Prereqesits:
- cockroachDB or postgreSQL with user and database
`,
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.GetViper())
			err := verifyZitadel(config.Database)
			logging.OnError(err).Fatal("unable to init zitadel")
		},
	}
}

func VerifyZitadel(db *sql.DB, config database.Config) error {
	err := ReadStmts(config.Type())
	if err != nil {
		return err
	}

	if err := exec(db, fmt.Sprintf(createSystemStmt, config.Username()), nil); err != nil {
		return err
	}

	if err := createEncryptionKeys(db); err != nil {
		return err
	}

	if err := exec(db, fmt.Sprintf(createProjectionsStmt, config.Username()), nil); err != nil {
		return err
	}

	if err := exec(db, fmt.Sprintf(createEventstoreStmt, config.Username()), nil); err != nil {
		return err
	}

	if err := createEvents(db); err != nil {
		return err
	}

	if err := exec(db, createSystemSequenceStmt, nil); err != nil {
		return err
	}

	if err := exec(db, createUniqueConstraints, nil); err != nil {
		return err
	}
	return nil
}

func verifyZitadel(config database.Config) error {
	logging.WithFields("database", config.DatabaseName()).Info("verify zitadel")

	db, err := database.Connect(config, false)
	if err != nil {
		return err
	}

	if err := VerifyZitadel(db.DB, config); err != nil {
		return err
	}

	return db.Close()
}

func createEncryptionKeys(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec(createEncryptionKeysStmt); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func createEvents(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(createEventsStmt); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
