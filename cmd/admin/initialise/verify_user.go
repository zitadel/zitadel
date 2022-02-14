package initialise

import (
	"database/sql"
	_ "embed"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	searchUser = "SELECT username FROM [show roles] WHERE username = $1"
	//go:embed sql/user.sql
	createUserStmt string
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
			return initialise(config, verifyUser(config.Database))
		},
	}
}

func verifyUser(config database.Config) func(*sql.DB) error {
	return func(db *sql.DB) error {
		logging.WithFields("username", config.Username).Info("verify user")
		return verify(db,
			exists(searchUser, config.Username),
			exec(createUserStmt, config.Username, &sql.NullString{String: config.Password, Valid: config.Password != ""}),
		)
	}
}
