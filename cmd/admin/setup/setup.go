package setup

import (
	"context"
	_ "embed"

	"github.com/caos/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/caos/zitadel/internal/api/authz"
	http_util "github.com/caos/zitadel/internal/api/http"
	command "github.com/caos/zitadel/internal/command/v2"
	"github.com/caos/zitadel/internal/database"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/migration"
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
			config := MustNewConfig(viper.GetViper())
			steps := MustNewSteps(viper.New())

			Setup(config, steps)
		},
	}
}

func Setup(config *Config, steps *Steps) {
	dbClient, err := database.Connect(config.Database)
	logging.OnError(err).Fatal("unable to connect to database")

	eventstoreClient, err := eventstore.Start(dbClient)
	logging.OnError(err).Fatal("unable to start eventstore")
	migration.RegisterMappers(eventstoreClient)

	cmd := command.New(eventstoreClient, "localhost", config.SystemDefaults)

	steps.S2DefaultInstance.cmd = cmd
	steps.S1ProjectionTable = &ProjectionTable{dbClient: dbClient}
	steps.S2DefaultInstance.InstanceSetup.Zitadel.IsDevMode = !config.ExternalSecure
	steps.S2DefaultInstance.InstanceSetup.Zitadel.BaseURL = http_util.BuildHTTP(config.ExternalDomain, config.ExternalPort, config.ExternalSecure)

	ctx := authz.WithInstance(context.Background(), "system")
	migration.Migrate(ctx, eventstoreClient, steps.S2DefaultInstance)
}
