package initialise

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
)

func newGrant() *cobra.Command {
	return &cobra.Command{
		Use:   "grant",
		Short: "set ALL grant to user",
		Long: `Sets ALL grant to the database user.

Prereqesits:
- cockroachDB or postgreSQL
`,
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.GetViper())

			err := initialise(config.Database, VerifyGrant(config.Database.DatabaseName(), config.Database.Username()))
			logging.OnError(err).Fatal("unable to set grant")
		},
	}
}

func VerifyGrant(databaseName, username string) func(*sql.DB) error {
	return func(db *sql.DB) error {
		logging.WithFields("user", username, "database", databaseName).Info("verify grant")

		return exec(db, fmt.Sprintf(grantStmt, databaseName, username), nil)
	}
}
