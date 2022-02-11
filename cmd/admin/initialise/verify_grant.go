package initialise

import (
	"database/sql"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			return initialise(config, verifyGrant)
		},
	}
}

func verifyGrant(db *sql.DB, config database.Config) error {
	logging.Info("verify grant")
	exists, err := hasGrant(db, config)
	if exists || err != nil {
		return err
	}
	return grant(db, config)
}

func hasGrant(db *sql.DB, config database.Config) (has bool, err error) {
	row := db.QueryRow("SELECT EXISTS(SELECT * FROM [SHOW GRANTS ON DATABASE "+config.Database+"] where grantee = $1 AND privilege_type = 'ALL')", config.Username)
	err = row.Scan(&has)
	return has, err
}

func grant(db *sql.DB, config database.Config) error {
	_, err := db.Exec("GRANT ALL ON DATABASE " + config.Database + " TO " + config.Username)
	return err
}
