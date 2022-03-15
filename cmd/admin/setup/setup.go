package setup

import (
	"bytes"
	"context"
	_ "embed"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/authz"
	command "github.com/caos/zitadel/internal/command/v2"
	"github.com/caos/zitadel/internal/config/hook"
	"github.com/caos/zitadel/internal/database"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/migration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//go:embed steps.yaml
	defaultSteps []byte
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

			v := viper.New()
			v.SetConfigType("yaml")
			err = v.ReadConfig(bytes.NewBuffer(defaultSteps))
			logging.OnError(err).Fatal("unable to read setup steps")

			steps := new(Steps)
			err = v.Unmarshal(steps,
				viper.DecodeHook(hook.Base64ToBytesHookFunc()),
				viper.DecodeHook(hook.TagToLanguageHookFunc()),
			)
			logging.OnError(err).Fatal("unable to read steps")

			setup(config, steps)
		},
	}
}

func setup(config *Config, steps *Steps) {
	dbClient, err := database.Connect(config.Database)
	logging.OnError(err).Fatal("unable to connect to database")

	eventstoreClient, err := eventstore.Start(dbClient)
	logging.OnError(err).Fatal("unable to start eventstore")

	cmd := command.New(eventstoreClient, "localhost")

	steps.S1DefaultInstance.cmd = cmd

	ctx := authz.WithTenant(context.Background(), "system")
	migration.Migrate(ctx, eventstoreClient, steps.S1DefaultInstance)
}
