package setup

import (
	"context"
	_ "embed"

	"github.com/caos/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/caos/zitadel/cmd/admin/key"
	http_util "github.com/caos/zitadel/internal/api/http"
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

			masterKey, err := key.MasterKey(cmd)
			logging.OnError(err).Panic("No master key provided")

			Setup(config, steps, masterKey)
		},
	}
}

func Setup(config *Config, steps *Steps, masterKey string) {
	dbClient, err := database.Connect(config.Database)
	logging.OnError(err).Fatal("unable to connect to database")

	eventstoreClient, err := eventstore.Start(dbClient)
	logging.OnError(err).Fatal("unable to start eventstore")
	migration.RegisterMappers(eventstoreClient)

	steps.S2DefaultInstance.es = eventstoreClient
	steps.S2DefaultInstance.db = dbClient
	steps.S2DefaultInstance.defaults = config.SystemDefaults
	steps.S2DefaultInstance.masterKey = masterKey
	steps.S2DefaultInstance.iamDomain = config.SystemDefaults.Domain
	steps.S2DefaultInstance.zitadelRoles = config.InternalAuthZ.RolePermissionMappings
	steps.S2DefaultInstance.userEncryptionKey = config.EncryptionKeys.User
	steps.S2DefaultInstance.InstanceSetup.Zitadel.IsDevMode = !config.ExternalSecure
	steps.S2DefaultInstance.InstanceSetup.Zitadel.BaseURL = http_util.BuildHTTP(config.ExternalDomain, config.ExternalPort, config.ExternalSecure)

	steps.S1ProjectionTable = &ProjectionTable{dbClient: dbClient}

	ctx := context.Background()
	migration.Migrate(ctx, eventstoreClient, steps.S1ProjectionTable)
	migration.Migrate(ctx, eventstoreClient, steps.S2DefaultInstance)
}
