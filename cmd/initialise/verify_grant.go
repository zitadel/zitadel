package initialise

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
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
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.New())

			err := initialise(config, VerifyGrant(config.Database.Database, config.Database.Username))
			logging.OnError(err).Fatal("unable to set grant")
		},
	}
}

func VerifyGrant(database, username string) func(*sql.DB) error {
	return func(db *sql.DB) error {
		logging.WithFields("user", username, "database", database).Info("verify grant")
		return verify(db,
			exists(fmt.Sprintf(searchGrant, database), username),
			exec(fmt.Sprintf(grantStmt, database, username)),
		)
	}
}
