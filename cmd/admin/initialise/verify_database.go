package initialise

import (
	"database/sql"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newDatabase() *cobra.Command {
	return &cobra.Command{
		Use:   "database",
		Short: "initialize only the database",
		Long: `Sets up the ZITADEL database.

Prereqesits:
- cockroachdb

The user provided by flags needs priviledge to 
- create the database if it does not exist
- see other users and create a new one if the user does not exist
- grant all rights of the ZITADEL database to the user created if not yet set
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := Config{}
			if err := viper.Unmarshal(&config); err != nil {
				return err
			}
			return initialise(config, verifyDB)
		},
	}
}

func verifyDB(db *sql.DB, config database.Config) error {
	logging.Info("verify database")
	exists, err := existsDatabase(db, config)
	if exists || err != nil {
		return err
	}
	return createDatabase(db, config)
}

func existsDatabase(db *sql.DB, config database.Config) (exists bool, err error) {
	row := db.QueryRow("SELECT EXISTS(SELECT database_name FROM [show databases] WHERE database_name = $1)", config.Database)
	err = row.Scan(&exists)
	return exists, err
}

func createDatabase(db *sql.DB, config database.Config) error {
	_, err := db.Exec("CREATE DATABASE " + config.Database)
	return err
}
