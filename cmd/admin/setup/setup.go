package setup

import (
	"context"
	_ "embed"
	"strings"

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

	steps.s1ProjectionTable = &ProjectionTable{dbClient: dbClient}
	steps.s2AssetsTable = &AssetTable{dbClient: dbClient}

	steps.S3DefaultInstance.InstanceSetup.Org.Human.Email.Address = strings.TrimSpace(steps.S3DefaultInstance.InstanceSetup.Org.Human.Email.Address)
	if steps.S3DefaultInstance.InstanceSetup.Org.Human.Email.Address == "" {
		steps.S3DefaultInstance.InstanceSetup.Org.Human.Email.Address = "admin@" + config.ExternalDomain
	}

	steps.S3DefaultInstance.es = eventstoreClient
	steps.S3DefaultInstance.db = dbClient
	steps.S3DefaultInstance.defaults = config.SystemDefaults
	steps.S3DefaultInstance.masterKey = masterKey
	steps.S3DefaultInstance.domain = config.ExternalDomain
	steps.S3DefaultInstance.zitadelRoles = config.InternalAuthZ.RolePermissionMappings
	steps.S3DefaultInstance.userEncryptionKey = config.EncryptionKeys.User
	steps.S3DefaultInstance.InstanceSetup.Zitadel.IsDevMode = !config.ExternalSecure
	steps.S3DefaultInstance.InstanceSetup.Zitadel.BaseURL = http_util.BuildHTTP(config.ExternalDomain, config.ExternalPort, config.ExternalSecure)
	steps.S3DefaultInstance.InstanceSetup.Zitadel.IsDevMode = !config.ExternalSecure
	steps.S3DefaultInstance.InstanceSetup.Zitadel.BaseURL = http_util.BuildHTTP(config.ExternalDomain, config.ExternalPort, config.ExternalSecure)

	ctx := context.Background()
	migration.Migrate(ctx, eventstoreClient, steps.s1ProjectionTable)
	migration.Migrate(ctx, eventstoreClient, steps.s2AssetsTable)
	migration.Migrate(ctx, eventstoreClient, steps.S3DefaultInstance)
}
