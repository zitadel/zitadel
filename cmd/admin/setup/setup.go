package setup

import (
	"context"
	_ "embed"

	"github.com/caos/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/command/v2"
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

	steps.s1ProjectionTable = &ProjectionTable{dbClient: dbClient}
	steps.s2AssetsTable = &AssetTable{dbClient: dbClient}
	steps.S3DefaultInstance.cmd = cmd
	steps.S3DefaultInstance.InstanceSetup.Zitadel.IsDevMode = !config.ExternalSecure
	steps.S3DefaultInstance.InstanceSetup.Zitadel.BaseURL = http_util.BuildHTTP(config.ExternalDomain, config.ExternalPort, config.ExternalSecure)

	ctx := context.Background()
	migration.Migrate(ctx, eventstoreClient, steps.s1ProjectionTable)
	migration.Migrate(ctx, eventstoreClient, steps.s2AssetsTable)
	migration.Migrate(ctx, eventstoreClient, steps.S3DefaultInstance)
}
