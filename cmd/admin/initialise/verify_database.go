package initialise

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	searchDatabase = "SELECT database_name FROM [show databases] WHERE database_name = $1"

	//go:embed sql/02_database.sql
	databaseStmt string
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
			return initialise(config, VerifyDatabase(config.Database.Database))
		},
	}
}

func VerifyDatabase(database string) func(*sql.DB) error {
	return func(db *sql.DB) error {
		return verify(db,
			exists(searchDatabase, database),
			exec(fmt.Sprintf(databaseStmt, database)),
		)
	}
}
