package initialise

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	searchGrant = "SELECT * FROM [SHOW GRANTS ON DATABASE %s] where grantee = $1 AND privilege_type = 'ALL'"
	//go:embed sql/grant_user.sql
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
			return initialise(config, verifyGrant(config.Database))
		},
	}
}

func verifyGrant(config database.Config) func(*sql.DB) error {
	return func(db *sql.DB) error {
		logging.WithFields("user", config.Username).Info("verify grant")
		return verify(db,
			exists(fmt.Sprintf(searchGrant, config.Database), config.Username),
			exec(fmt.Sprintf(grantStmt, config.Database, config.Username)),
		)
	}
}
