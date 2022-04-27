package setup

import (
	"context"
	_ "embed"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/admin/key"
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
}

func Setup(config *Config, steps *Steps, masterKey string) {
	dbClient, err := database.Connect(config.Database)
	logging.OnError(err).Fatal("unable to connect to database")

	eventstoreClient, err := eventstore.Start(dbClient)
	logging.OnError(err).Fatal("unable to start eventstore")
	migration.RegisterMappers(eventstoreClient)

	steps.s1ProjectionTable = &ProjectionTable{dbClient: dbClient}
	steps.s2AssetsTable = &AssetTable{dbClient: dbClient}

	instanceSetup := config.DefaultInstance
	instanceSetup.InstanceName = steps.S3DefaultInstance.InstanceSetup.InstanceName
	instanceSetup.CustomDomain = steps.S3DefaultInstance.InstanceSetup.CustomDomain
	instanceSetup.Org = steps.S3DefaultInstance.InstanceSetup.Org
	steps.S3DefaultInstance.InstanceSetup = instanceSetup

	steps.S3DefaultInstance.InstanceSetup.Org.Human.Email.Address = strings.TrimSpace(steps.S3DefaultInstance.InstanceSetup.Org.Human.Email.Address)
	if steps.S3DefaultInstance.InstanceSetup.Org.Human.Email.Address == "" {
		steps.S3DefaultInstance.InstanceSetup.Org.Human.Email.Address = "admin@" + instanceSetup.CustomDomain
	}

	steps.S3DefaultInstance.es = eventstoreClient
	steps.S3DefaultInstance.db = dbClient
	steps.S3DefaultInstance.defaults = config.SystemDefaults
	steps.S3DefaultInstance.masterKey = masterKey
	steps.S3DefaultInstance.domain = instanceSetup.CustomDomain
	steps.S3DefaultInstance.zitadelRoles = config.InternalAuthZ.RolePermissionMappings
	steps.S3DefaultInstance.userEncryptionKey = config.EncryptionKeys.User
	steps.S3DefaultInstance.externalSecure = config.ExternalSecure
	steps.S3DefaultInstance.externalPort = config.ExternalPort

	ctx := context.Background()
	err = migration.Migrate(ctx, eventstoreClient, steps.s1ProjectionTable)
	logging.OnError(err).Fatal("unable to migrate step 1")
	err = migration.Migrate(ctx, eventstoreClient, steps.s2AssetsTable)
	logging.OnError(err).Fatal("unable to migrate step 3")
	err = migration.Migrate(ctx, eventstoreClient, steps.S3DefaultInstance)
	logging.OnError(err).Fatal("unable to migrate step 4")
}

func initSteps(v *viper.Viper, files ...string) func() {
	return func() {
		for _, file := range files {
			v.SetConfigFile(file)
			err := v.MergeInConfig()
			logging.WithFields("file", file).OnError(err).Warn("unable to read setup file")
		}
	}
}
