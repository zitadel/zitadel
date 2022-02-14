package initialise

import (
	"database/sql"
	_ "embed"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//sql import
	_ "github.com/lib/pq"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize ZITADEL instance",
		Long: `Sets up the minimum requirements to start ZITADEL.

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
			if err := initialise(config, verifyUser, verifyDB, verifyGrant); err != nil {
				return err
			}

			return verifyZitadel(config.Database)
		},
	}

	cmd.AddCommand(newZitadel(), newDatabase(), newUser(), newGrant())
	return cmd
}

func initialise(config Config, steps ...func(*sql.DB, database.Config) error) error {
	logging.Info("initialization started")

	db, err := database.Connect(adminConfig(config))
	if err != nil {
		return err
	}

	for _, step := range steps {
		if err = step(db, config.Database); err != nil {
			return err
		}
	}

	return db.Close()
}
