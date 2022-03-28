package initialise

import (
	"database/sql"
	_ "embed"

	"github.com/caos/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/caos/zitadel/internal/database"
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
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.GetViper())

			InitAll(config)
		},
	}

	cmd.AddCommand(newZitadel(), newDatabase(), newUser(), newGrant())
	return cmd
}

func InitAll(config *Config) {
	err := initialise(config,
		VerifyUser(config.Database.Username, config.Database.Password),
		VerifyDatabase(config.Database.Database),
		VerifyGrant(config.Database.Database, config.Database.Username),
	)
	logging.OnError(err).Fatal("unable to initialize the database")

	err = verifyZitadel(config.Database)
	logging.OnError(err).Fatal("unable to initialize ZITADEL")
}

func initialise(config *Config, steps ...func(*sql.DB) error) error {
	logging.Info("initialization started")

	db, err := database.Connect(adminConfig(config))
	if err != nil {
		return err
	}
	err = Initialise(db, steps...)
	if err != nil {
		return err
	}
	return db.Close()
}

func Initialise(db *sql.DB, steps ...func(*sql.DB) error) error {
	for _, step := range steps {
		if err := step(db); err != nil {
			return err
		}
	}
	return nil
}
