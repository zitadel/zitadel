package setup

import (
	_ "embed"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/database"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "setup ZITADEL instance",
		Long: `sets up data to start ZITADEL.
Requirements:
- cockroachdb`,
		Run: func(cmd *cobra.Command, args []string) {
			config := new(Config)
			err := viper.Unmarshal(config)
			logging.OnError(err).Fatal("unable to read config")

			setup(config)
		},
	}
}

func setup(config *Config) {
	dbClient, err := database.Connect(config.Database)
	logging.OnError(err).Fatal("unable to connect to database")

	eventstoreClient, err := eventstore.Start(dbClient)
	logging.OnError(err).Fatal("unable to start eventstore")

	commands, err := command.StartCommands(eventstoreClient, config.SystemDefaults, config.InternalAuthZ, nil, nil, nil, nil, nil, nil)
	logging.OnError(err).Fatal("unable to start commands")

}
