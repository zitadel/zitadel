package initialise

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/caos/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	searchGrant = "SELECT * FROM [SHOW GRANTS ON DATABASE %s] where grantee = $1 AND privilege_type = 'ALL'"
	//go:embed sql/03_grant_user.sql
	grantStmt string
)

func newGrant() *cobra.Command {
	return &cobra.Command{
		Use:   "grant",
		Short: "set ALL grant to user",
		Long: `Sets ALL grant to the database user.

Prereqesits:
- cockroachdb
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := Config{}
			if err := viper.Unmarshal(&config); err != nil {
				return err
			}
			return initialise(config, VerifyGrant(config.Database.Database, config.Database.User.Username))
		},
	}
}

func VerifyGrant(database, username string) func(*sql.DB) error {
	return func(db *sql.DB) error {
		logging.WithFields("user", username).Info("verify grant")
		return verify(db,
			exists(fmt.Sprintf(searchGrant, database), username),
			exec(fmt.Sprintf(grantStmt, database, username)),
		)
	}
}
