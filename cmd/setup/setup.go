package setup

import (
	"context"
	_ "embed"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/tls"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
)

var (
	//go:embed steps.yaml
	defaultSteps []byte
	stepFiles    []string
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "setup ZITADEL instance",
		Long: `sets up data to start ZITADEL.
Requirements:
- cockroachdb`,
		Run: func(cmd *cobra.Command, args []string) {
			err := tls.ModeFromFlag(cmd)
			logging.OnError(err).Fatal("invalid tlsMode")

			config := MustNewConfig(viper.GetViper())
			steps := MustNewSteps(viper.New())

			masterKey, err := key.MasterKey(cmd)
			logging.OnError(err).Panic("No master key provided")

			Setup(config, steps, masterKey)
		},
	}

	Flags(cmd)

	return cmd
}

func Flags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringArrayVar(&stepFiles, "steps", nil, "paths to step files to overwrite default steps")
	key.AddMasterKeyFlag(cmd)
	tls.AddTLSModeFlag(cmd)
}

func Setup(config *Config, steps *Steps, masterKey string) {
	dbClient, err := database.Connect(config.Database)
	logging.OnError(err).Fatal("unable to connect to database")

	eventstoreClient, err := eventstore.Start(dbClient)
	logging.OnError(err).Fatal("unable to start eventstore")
	migration.RegisterMappers(eventstoreClient)

	steps.s1ProjectionTable = &ProjectionTable{dbClient: dbClient}
	steps.s2AssetsTable = &AssetTable{dbClient: dbClient}

	steps.S3DefaultInstance.instanceSetup = config.DefaultInstance
	steps.S3DefaultInstance.userEncryptionKey = config.EncryptionKeys.User
	steps.S3DefaultInstance.smtpEncryptionKey = config.EncryptionKeys.SMTP
	steps.S3DefaultInstance.masterKey = masterKey
	steps.S3DefaultInstance.db = dbClient
	steps.S3DefaultInstance.es = eventstoreClient
	steps.S3DefaultInstance.defaults = config.SystemDefaults
	steps.S3DefaultInstance.zitadelRoles = config.InternalAuthZ.RolePermissionMappings
	steps.S3DefaultInstance.externalDomain = config.ExternalDomain
	steps.S3DefaultInstance.externalSecure = config.ExternalSecure
	steps.S3DefaultInstance.externalPort = config.ExternalPort

	repeatableSteps := []migration.RepeatableMigration{
		&configChange{
			es:             eventstoreClient,
			ExternalDomain: config.ExternalDomain,
			ExternalPort:   config.ExternalPort,
			ExternalSecure: config.ExternalSecure,
		},
	}

	ctx := context.Background()
	err = migration.Migrate(ctx, eventstoreClient, steps.s1ProjectionTable)
	logging.OnError(err).Fatal("unable to migrate step 1")
	err = migration.Migrate(ctx, eventstoreClient, steps.s2AssetsTable)
	logging.OnError(err).Fatal("unable to migrate step 2")
	err = migration.Migrate(ctx, eventstoreClient, steps.S3DefaultInstance)
	logging.OnError(err).Fatal("unable to migrate step 3")

	for _, repeatableStep := range repeatableSteps {
		err = migration.Migrate(ctx, eventstoreClient, repeatableStep)
		logging.OnError(err).Fatalf("unable to migrate repeatable step: %s", repeatableStep.String())
	}
}
