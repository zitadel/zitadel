package setup

import (
	"context"
	"embed"
	_ "embed"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/tls"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
	"github.com/zitadel/zitadel/internal/query/projection"
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
	ctx := context.Background()
	logging.Info("setup started")

	dbClient, err := database.Connect(config.Database, false)
	logging.OnError(err).Fatal("unable to connect to database")

	eventstoreClient, err := eventstore.Start(&eventstore.Config{Client: dbClient})
	logging.OnError(err).Fatal("unable to start eventstore")
	migration.RegisterMappers(eventstoreClient)

	steps.s1ProjectionTable = &ProjectionTable{dbClient: dbClient.DB}
	steps.s2AssetsTable = &AssetTable{dbClient: dbClient.DB}

	steps.FirstInstance.instanceSetup = config.DefaultInstance
	steps.FirstInstance.userEncryptionKey = config.EncryptionKeys.User
	steps.FirstInstance.smtpEncryptionKey = config.EncryptionKeys.SMTP
	steps.FirstInstance.masterKey = masterKey
	steps.FirstInstance.db = dbClient.DB
	steps.FirstInstance.es = eventstoreClient
	steps.FirstInstance.defaults = config.SystemDefaults
	steps.FirstInstance.zitadelRoles = config.InternalAuthZ.RolePermissionMappings
	steps.FirstInstance.externalDomain = config.ExternalDomain
	steps.FirstInstance.externalSecure = config.ExternalSecure
	steps.FirstInstance.externalPort = config.ExternalPort

	steps.s4EventstoreIndexes = New04(dbClient)
	steps.s5LastFailed = &LastFailed{dbClient: dbClient.DB}
	steps.s6OwnerRemoveColumns = &OwnerRemoveColumns{dbClient: dbClient.DB}
	steps.s7LogstoreTables = &LogstoreTables{dbClient: dbClient.DB, username: config.Database.Username(), dbType: config.Database.Type()}
	steps.s8AuthTokens = &AuthTokenIndexes{dbClient: dbClient}
	steps.s9EventstoreIndexes2 = New09(dbClient)

	err = projection.Create(ctx, dbClient, eventstoreClient, config.Projections, nil, nil)
	logging.OnError(err).Fatal("unable to start projections")

	repeatableSteps := []migration.RepeatableMigration{
		&externalConfigChange{
			es:             eventstoreClient,
			ExternalDomain: config.ExternalDomain,
			ExternalPort:   config.ExternalPort,
			ExternalSecure: config.ExternalSecure,
		},
		&projectionTables{
			es:      eventstoreClient,
			Version: build.Version(),
		},
	}

	err = migration.Migrate(ctx, eventstoreClient, steps.s1ProjectionTable)
	logging.OnError(err).Fatal("unable to migrate step 1")
	err = migration.Migrate(ctx, eventstoreClient, steps.s2AssetsTable)
	logging.OnError(err).Fatal("unable to migrate step 2")
	err = migration.Migrate(ctx, eventstoreClient, steps.FirstInstance)
	logging.OnError(err).Fatal("unable to migrate step 3")
	err = migration.Migrate(ctx, eventstoreClient, steps.s4EventstoreIndexes)
	logging.OnError(err).Fatal("unable to migrate step 4")
	err = migration.Migrate(ctx, eventstoreClient, steps.s5LastFailed)
	logging.OnError(err).Fatal("unable to migrate step 5")
	err = migration.Migrate(ctx, eventstoreClient, steps.s6OwnerRemoveColumns)
	logging.OnError(err).Fatal("unable to migrate step 6")
	err = migration.Migrate(ctx, eventstoreClient, steps.s7LogstoreTables)
	logging.OnError(err).Fatal("unable to migrate step 7")
	err = migration.Migrate(ctx, eventstoreClient, steps.s8AuthTokens)
	logging.OnError(err).Fatal("unable to migrate step 8")
	err = migration.Migrate(ctx, eventstoreClient, steps.s9EventstoreIndexes2)
	logging.OnError(err).Fatal("unable to migrate step 9")

	for _, repeatableStep := range repeatableSteps {
		err = migration.Migrate(ctx, eventstoreClient, repeatableStep)
		logging.OnError(err).Fatalf("unable to migrate repeatable step: %s", repeatableStep.String())
	}
}

func readStmt(fs embed.FS, folder, typ, filename string) (string, error) {
	stmt, err := fs.ReadFile(folder + "/" + typ + "/" + filename)
	return string(stmt), err
}
