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
	old_es "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	new_es "github.com/zitadel/zitadel/internal/eventstore/v3"
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

	cmd.AddCommand(NewCleanup())

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

	zitadelDBClient, err := database.Connect(config.Database, false, false)
	logging.OnError(err).Fatal("unable to connect to database")
	esPusherDBClient, err := database.Connect(config.Database, false, true)
	logging.OnError(err).Fatal("unable to connect to database")

	config.Eventstore.Querier = old_es.NewCRDB(zitadelDBClient)
	config.Eventstore.Pusher = new_es.NewEventstore(esPusherDBClient)
	eventstoreClient := eventstore.NewEventstore(config.Eventstore)
	logging.OnError(err).Fatal("unable to start eventstore")
	migration.RegisterMappers(eventstoreClient)

	steps.s1ProjectionTable = &ProjectionTable{dbClient: zitadelDBClient.DB}
	steps.s2AssetsTable = &AssetTable{dbClient: zitadelDBClient.DB}

	steps.FirstInstance.instanceSetup = config.DefaultInstance
	steps.FirstInstance.userEncryptionKey = config.EncryptionKeys.User
	steps.FirstInstance.smtpEncryptionKey = config.EncryptionKeys.SMTP
	steps.FirstInstance.oidcEncryptionKey = config.EncryptionKeys.OIDC
	steps.FirstInstance.masterKey = masterKey
	steps.FirstInstance.db = zitadelDBClient
	steps.FirstInstance.es = eventstoreClient
	steps.FirstInstance.defaults = config.SystemDefaults
	steps.FirstInstance.zitadelRoles = config.InternalAuthZ.RolePermissionMappings
	steps.FirstInstance.externalDomain = config.ExternalDomain
	steps.FirstInstance.externalSecure = config.ExternalSecure
	steps.FirstInstance.externalPort = config.ExternalPort

	steps.s5LastFailed = &LastFailed{dbClient: zitadelDBClient.DB}
	steps.s6OwnerRemoveColumns = &OwnerRemoveColumns{dbClient: zitadelDBClient.DB}
	steps.s7LogstoreTables = &LogstoreTables{dbClient: zitadelDBClient.DB, username: config.Database.Username(), dbType: config.Database.Type()}
	steps.s8AuthTokens = &AuthTokenIndexes{dbClient: zitadelDBClient}
	steps.CorrectCreationDate.dbClient = esPusherDBClient
	steps.s12AddOTPColumns = &AddOTPColumns{dbClient: zitadelDBClient}
	steps.s13FixQuotaProjection = &FixQuotaConstraints{dbClient: zitadelDBClient}
	steps.s14NewEventsTable = &NewEventsTable{dbClient: esPusherDBClient}
	steps.s15CurrentStates = &CurrentProjectionState{dbClient: zitadelDBClient}

	err = projection.Create(ctx, zitadelDBClient, eventstoreClient, config.Projections, nil, nil, nil)
	logging.OnError(err).Fatal("unable to start projections")

	repeatableSteps := []migration.RepeatableMigration{
		&externalConfigChange{
			es:             eventstoreClient,
			ExternalDomain: config.ExternalDomain,
			ExternalPort:   config.ExternalPort,
			ExternalSecure: config.ExternalSecure,
			defaults:       config.SystemDefaults,
		},
		&projectionTables{
			es:      eventstoreClient,
			Version: build.Version(),
		},
	}

	err = migration.Migrate(ctx, eventstoreClient, steps.s14NewEventsTable)
	logging.WithFields("name", steps.s14NewEventsTable.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.s1ProjectionTable)
	logging.WithFields("name", steps.s1ProjectionTable.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.s2AssetsTable)
	logging.WithFields("name", steps.s2AssetsTable.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.FirstInstance)
	logging.WithFields("name", steps.FirstInstance.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.s5LastFailed)
	logging.WithFields("name", steps.s5LastFailed.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.s6OwnerRemoveColumns)
	logging.WithFields("name", steps.s6OwnerRemoveColumns.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.s7LogstoreTables)
	logging.WithFields("name", steps.s7LogstoreTables.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.s8AuthTokens)
	logging.WithFields("name", steps.s8AuthTokens.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.s12AddOTPColumns)
	logging.WithFields("name", steps.s12AddOTPColumns.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.s13FixQuotaProjection)
	logging.WithFields("name", steps.s13FixQuotaProjection.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.s15CurrentStates)
	logging.WithFields("name", steps.s15CurrentStates.String()).OnError(err).Fatal("migration failed")

	for _, repeatableStep := range repeatableSteps {
		err = migration.Migrate(ctx, eventstoreClient, repeatableStep)
		logging.OnError(err).Fatalf("unable to migrate repeatable step: %s", repeatableStep.String())
	}
}

func readStmt(fs embed.FS, folder, typ, filename string) (string, error) {
	stmt, err := fs.ReadFile(folder + "/" + typ + "/" + filename)
	return string(stmt), err
}
