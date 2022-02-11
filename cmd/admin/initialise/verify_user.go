package initialise

import (
	"database/sql"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newUser() *cobra.Command {
	return &cobra.Command{
		Use:   "user",
		Short: "initialize only the database user",
		Long: `Sets up the ZITADEL database user.

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
			return initialise(config, verifyUser)
		},
	}
}

func verifyUser(db *sql.DB, config database.Config) error {
	logging.Info("verify user")
	exists, err := existsUser(db, config)
	if exists || err != nil {
		return err
	}
	return createUser(db, config)
}

func existsUser(db *sql.DB, config database.Config) (exists bool, err error) {
	row := db.QueryRow("SELECT EXISTS(SELECT username FROM [show roles] WHERE username = $1)", config.Username)
	err = row.Scan(&exists)
	return exists, err
}

func createUser(db *sql.DB, config database.Config) error {
	_, err := db.Exec("CREATE USER $1 WITH PASSWORD $2", config.Username, &sql.NullString{String: config.Password, Valid: config.Password != ""})
	return err
}
