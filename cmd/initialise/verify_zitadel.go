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

const (
	eventstoreSchema       = "eventstore"
	eventsTable            = "events"
	uniqueConstraintsTable = "unique_constraints"
	projectionsSchema      = "projections"
	systemSchema           = "system"
	encryptionKeysTable    = "encryption_keys"
)

func newZitadel() *cobra.Command {
	return &cobra.Command{
		Use:   "zitadel",
		Short: "initialize ZITADEL internals",
		Long: `initialize ZITADEL internals.

Prereqesits:
- cockroachdb with user and database
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := MustNewConfig(viper.GetViper())
			return verifyZitadel(config.Database)
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
	logging.WithFields("database", config.Database()).Info("verify zitadel")

	db, err := database.Connect(config, false)
	if err != nil {
		return err
	}

	if err := VerifyZitadel(db, config); err != nil {
		return nil
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
