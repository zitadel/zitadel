package setup

import (
	"context"
	"embed"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/cmd/encryption"
	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/tls"
	admin_handler "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/handler"
	admin_view "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	auth_handler "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/handler"
	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/authz"
	authz_es "github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/eventstore"
	"github.com/zitadel/zitadel/internal/command"
	cryptoDB "github.com/zitadel/zitadel/internal/crypto/database"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_es "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	new_es "github.com/zitadel/zitadel/internal/eventstore/v3"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/migration"
	notify_handler "github.com/zitadel/zitadel/internal/notification"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
	es_v4 "github.com/zitadel/zitadel/internal/v2/eventstore"
	es_v4_pg "github.com/zitadel/zitadel/internal/v2/eventstore/postgres"
	"github.com/zitadel/zitadel/internal/webauthn"
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

			err = BindInitProjections(cmd)
			logging.OnError(err).Fatal("unable to bind \"init-projections\" flag")

			err = bindForMirror(cmd)
			logging.OnError(err).Fatal("unable to bind \"for-mirror\" flag")

			config := MustNewConfig(viper.GetViper())
			steps := MustNewSteps(viper.New())

			masterKey, err := key.MasterKey(cmd)
			logging.OnError(err).Panic("No master key provided")

			Setup(cmd.Context(), config, steps, masterKey)
		},
	}

	cmd.AddCommand(NewCleanup())

	Flags(cmd)

	return cmd
}

func Flags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringArrayVar(&stepFiles, "steps", nil, "paths to step files to overwrite default steps")
	cmd.Flags().Bool("init-projections", viper.GetBool("InitProjections"), "beta feature: initializes projections after they are created, allows smooth start as projections are up to date")
	cmd.Flags().Bool("for-mirror", viper.GetBool("ForMirror"), "use this flag if you want to mirror your existing data")
	key.AddMasterKeyFlag(cmd)
	tls.AddTLSModeFlag(cmd)
}

func BindInitProjections(cmd *cobra.Command) error {
	return viper.BindPFlag("InitProjections.Enabled", cmd.Flags().Lookup("init-projections"))
}

func bindForMirror(cmd *cobra.Command) error {
	return viper.BindPFlag("ForMirror", cmd.Flags().Lookup("for-mirror"))
}

func Setup(ctx context.Context, config *Config, steps *Steps, masterKey string) {
	logging.Info("setup started")

	i18n.MustLoadSupportedLanguagesFromDir()

	queryDBClient, err := database.Connect(config.Database, false, dialect.DBPurposeQuery)
	logging.OnError(err).Fatal("unable to connect to database")
	esPusherDBClient, err := database.Connect(config.Database, false, dialect.DBPurposeEventPusher)
	logging.OnError(err).Fatal("unable to connect to database")
	projectionDBClient, err := database.Connect(config.Database, false, dialect.DBPurposeProjectionSpooler)
	logging.OnError(err).Fatal("unable to connect to database")

	config.Eventstore.Querier = old_es.NewCRDB(queryDBClient)
	esV3 := new_es.NewEventstore(esPusherDBClient)
	config.Eventstore.Pusher = esV3
	config.Eventstore.Searcher = esV3
	eventstoreClient := eventstore.NewEventstore(config.Eventstore)

	logging.OnError(err).Fatal("unable to start eventstore")
	eventstoreV4 := es_v4.NewEventstoreFromOne(es_v4_pg.New(queryDBClient, &es_v4_pg.Config{
		MaxRetries: config.Eventstore.MaxRetries,
	}))

	steps.s1ProjectionTable = &ProjectionTable{dbClient: queryDBClient.DB}
	steps.s2AssetsTable = &AssetTable{dbClient: queryDBClient.DB}

	steps.FirstInstance.Skip = config.ForMirror || steps.FirstInstance.Skip
	steps.FirstInstance.instanceSetup = config.DefaultInstance
	steps.FirstInstance.userEncryptionKey = config.EncryptionKeys.User
	steps.FirstInstance.smtpEncryptionKey = config.EncryptionKeys.SMTP
	steps.FirstInstance.oidcEncryptionKey = config.EncryptionKeys.OIDC
	steps.FirstInstance.masterKey = masterKey
	steps.FirstInstance.db = queryDBClient
	steps.FirstInstance.es = eventstoreClient
	steps.FirstInstance.defaults = config.SystemDefaults
	steps.FirstInstance.zitadelRoles = config.InternalAuthZ.RolePermissionMappings
	steps.FirstInstance.externalDomain = config.ExternalDomain
	steps.FirstInstance.externalSecure = config.ExternalSecure
	steps.FirstInstance.externalPort = config.ExternalPort

	steps.s5LastFailed = &LastFailed{dbClient: queryDBClient.DB}
	steps.s6OwnerRemoveColumns = &OwnerRemoveColumns{dbClient: queryDBClient.DB}
	steps.s7LogstoreTables = &LogstoreTables{dbClient: queryDBClient.DB, username: config.Database.Username(), dbType: config.Database.Type()}
	steps.s8AuthTokens = &AuthTokenIndexes{dbClient: queryDBClient}
	steps.CorrectCreationDate.dbClient = esPusherDBClient
	steps.s12AddOTPColumns = &AddOTPColumns{dbClient: queryDBClient}
	steps.s13FixQuotaProjection = &FixQuotaConstraints{dbClient: queryDBClient}
	steps.s14NewEventsTable = &NewEventsTable{dbClient: esPusherDBClient}
	steps.s15CurrentStates = &CurrentProjectionState{dbClient: queryDBClient}
	steps.s16UniqueConstraintsLower = &UniqueConstraintToLower{dbClient: queryDBClient}
	steps.s17AddOffsetToUniqueConstraints = &AddOffsetToCurrentStates{dbClient: queryDBClient}
	steps.s18AddLowerFieldsToLoginNames = &AddLowerFieldsToLoginNames{dbClient: queryDBClient}
	steps.s19AddCurrentStatesIndex = &AddCurrentSequencesIndex{dbClient: queryDBClient}
	steps.s20AddByUserSessionIndex = &AddByUserIndexToSession{dbClient: queryDBClient}
	steps.s21AddBlockFieldToLimits = &AddBlockFieldToLimits{dbClient: queryDBClient}
	steps.s22ActiveInstancesIndex = &ActiveInstanceEvents{dbClient: queryDBClient}
	steps.s23CorrectGlobalUniqueConstraints = &CorrectGlobalUniqueConstraints{dbClient: esPusherDBClient}
	steps.s24AddActorToAuthTokens = &AddActorToAuthTokens{dbClient: queryDBClient}
	steps.s25User11AddLowerFieldsToVerifiedEmail = &User11AddLowerFieldsToVerifiedEmail{dbClient: esPusherDBClient}
	steps.s26AuthUsers3 = &AuthUsers3{dbClient: esPusherDBClient}
	steps.s27IDPTemplate6SAMLNameIDFormat = &IDPTemplate6SAMLNameIDFormat{dbClient: esPusherDBClient}
	steps.s28AddFieldTable = &AddFieldTable{dbClient: esPusherDBClient}
	steps.s29FillFieldsForProjectGrant = &FillFieldsForProjectGrant{eventstore: eventstoreClient}
	steps.s30FillFieldsForOrgDomainVerified = &FillFieldsForOrgDomainVerified{eventstore: eventstoreClient}
	steps.s31AddAggregateIndexToFields = &AddAggregateIndexToFields{dbClient: esPusherDBClient}
	steps.s32SMSConfigs2TwilioAddVerifyServiceSid = &SMSConfigs2TwilioAddVerifyServiceSid{dbClient: esPusherDBClient}

	err = projection.Create(ctx, projectionDBClient, eventstoreClient, config.Projections, nil, nil, nil)
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

	for _, step := range []migration.Migration{
		steps.s14NewEventsTable,
		steps.s1ProjectionTable,
		steps.s2AssetsTable,
		steps.s28AddFieldTable,
		steps.s31AddAggregateIndexToFields,
		steps.FirstInstance,
		steps.s5LastFailed,
		steps.s6OwnerRemoveColumns,
		steps.s7LogstoreTables,
		steps.s8AuthTokens,
		steps.s12AddOTPColumns,
		steps.s13FixQuotaProjection,
		steps.s15CurrentStates,
		steps.s16UniqueConstraintsLower,
		steps.s17AddOffsetToUniqueConstraints,
		steps.s19AddCurrentStatesIndex,
		steps.s20AddByUserSessionIndex,
		steps.s22ActiveInstancesIndex,
		steps.s23CorrectGlobalUniqueConstraints,
		steps.s24AddActorToAuthTokens,
		steps.s26AuthUsers3,
		steps.s29FillFieldsForProjectGrant,
		steps.s30FillFieldsForOrgDomainVerified,
	} {
		mustExecuteMigration(ctx, eventstoreClient, step, "migration failed")
	}

	for _, repeatableStep := range repeatableSteps {
		mustExecuteMigration(ctx, eventstoreClient, repeatableStep, "unable to migrate repeatable step")
	}

	// These steps are executed after the repeatable steps because they add fields projections
	for _, step := range []migration.Migration{
		steps.s18AddLowerFieldsToLoginNames,
		steps.s21AddBlockFieldToLimits,
		steps.s25User11AddLowerFieldsToVerifiedEmail,
		steps.s27IDPTemplate6SAMLNameIDFormat,
		steps.s32SMSConfigs2TwilioAddVerifyServiceSid,
	} {
		mustExecuteMigration(ctx, eventstoreClient, step, "migration failed")
	}

	// projection initialization must be done last, since the steps above might add required columns to the projections
	if !config.ForMirror && config.InitProjections.Enabled {
		initProjections(
			ctx,
			eventstoreClient,
			eventstoreV4,
			queryDBClient,
			projectionDBClient,
			masterKey,
			config,
		)
	}
}

func mustExecuteMigration(ctx context.Context, eventstoreClient *eventstore.Eventstore, step migration.Migration, errorMsg string) {
	err := migration.Migrate(ctx, eventstoreClient, step)
	logging.WithFields("name", step.String()).OnError(err).Fatal(errorMsg)
}

func readStmt(fs embed.FS, folder, typ, filename string) (string, error) {
	stmt, err := fs.ReadFile(folder + "/" + typ + "/" + filename)
	return string(stmt), err
}

func initProjections(
	ctx context.Context,
	eventstoreClient *eventstore.Eventstore,
	eventstoreV4 *es_v4.EventStore,
	queryDBClient,
	projectionDBClient *database.DB,
	masterKey string,
	config *Config,
) {
	logging.Info("init-projections is currently in beta")

	keyStorage, err := cryptoDB.NewKeyStorage(queryDBClient, masterKey)
	logging.OnError(err).Fatal("unable to start key storage")

	keys, err := encryption.EnsureEncryptionKeys(ctx, config.EncryptionKeys, keyStorage)
	logging.OnError(err).Fatal("unable to ensure encryption keys")

	err = projection.Create(
		ctx,
		queryDBClient,
		eventstoreClient,
		projection.Config{
			RetryFailedAfter: config.InitProjections.RetryFailedAfter,
			MaxFailureCount:  config.InitProjections.MaxFailureCount,
			BulkLimit:        config.InitProjections.BulkLimit,
		},
		keys.OIDC,
		keys.SAML,
		config.SystemAPIUsers,
	)
	logging.OnError(err).Fatal("unable to start projections")
	for _, p := range projection.Projections() {
		err := migration.Migrate(ctx, eventstoreClient, p)
		logging.WithFields("name", p.String()).OnError(err).Fatal("migration failed")
	}

	staticStorage, err := config.AssetStorage.NewStorage(queryDBClient.DB)
	logging.OnError(err).Fatal("unable to start asset storage")

	adminView, err := admin_view.StartView(queryDBClient)
	logging.OnError(err).Fatal("unable to start admin view")
	admin_handler.Register(ctx,
		admin_handler.Config{
			Client:                queryDBClient,
			Eventstore:            eventstoreClient,
			BulkLimit:             config.InitProjections.BulkLimit,
			FailureCountUntilSkip: uint64(config.InitProjections.MaxFailureCount),
		},
		adminView,
		staticStorage,
	)
	for _, p := range admin_handler.Projections() {
		err := migration.Migrate(ctx, eventstoreClient, p)
		logging.WithFields("name", p.String()).OnError(err).Fatal("migration failed")
	}

	sessionTokenVerifier := internal_authz.SessionTokenVerifier(keys.OIDC)
	queries, err := query.StartQueries(
		ctx,
		eventstoreClient,
		eventstoreV4.Querier,
		queryDBClient,
		projectionDBClient,
		config.Projections,
		config.SystemDefaults,
		keys.IDPConfig,
		keys.OTP,
		keys.OIDC,
		keys.SAML,
		config.InternalAuthZ.RolePermissionMappings,
		sessionTokenVerifier,
		func(q *query.Queries) domain.PermissionCheck {
			return func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				return internal_authz.CheckPermission(ctx, &authz_es.UserMembershipRepo{Queries: q}, config.InternalAuthZ.RolePermissionMappings, permission, orgID, resourceID)
			}
		},
		0,   // not needed for projections
		nil, // not needed for projections
		false,
	)
	logging.OnError(err).Fatal("unable to start queries")

	authView, err := auth_view.StartView(queryDBClient, keys.OIDC, queries, eventstoreClient)
	logging.OnError(err).Fatal("unable to start admin view")
	auth_handler.Register(ctx,
		auth_handler.Config{
			Client:                queryDBClient,
			Eventstore:            eventstoreClient,
			BulkLimit:             config.InitProjections.BulkLimit,
			FailureCountUntilSkip: uint64(config.InitProjections.MaxFailureCount),
		},
		authView,
		queries,
	)
	for _, p := range auth_handler.Projections() {
		err := migration.Migrate(ctx, eventstoreClient, p)
		logging.WithFields("name", p.String()).OnError(err).Fatal("migration failed")
	}

	authZRepo, err := authz.Start(queries, eventstoreClient, queryDBClient, keys.OIDC, config.ExternalSecure)
	logging.OnError(err).Fatal("unable to start authz repo")
	permissionCheck := func(ctx context.Context, permission, orgID, resourceID string) (err error) {
		return internal_authz.CheckPermission(ctx, authZRepo, config.InternalAuthZ.RolePermissionMappings, permission, orgID, resourceID)
	}
	commands, err := command.StartCommands(
		eventstoreClient,
		config.SystemDefaults,
		config.InternalAuthZ.RolePermissionMappings,
		staticStorage,
		&webauthn.Config{
			DisplayName:    config.WebAuthNName,
			ExternalSecure: config.ExternalSecure,
		},
		config.ExternalDomain,
		config.ExternalSecure,
		config.ExternalPort,
		keys.IDPConfig,
		keys.OTP,
		keys.SMTP,
		keys.SMS,
		keys.User,
		keys.DomainVerification,
		keys.OIDC,
		keys.SAML,
		&http.Client{},
		permissionCheck,
		sessionTokenVerifier,
		config.OIDC.DefaultAccessTokenLifetime,
		config.OIDC.DefaultRefreshTokenExpiration,
		config.OIDC.DefaultRefreshTokenIdleExpiration,
		config.DefaultInstance.SecretGenerators,
	)
	logging.OnError(err).Fatal("unable to start commands")
	notify_handler.Register(
		ctx,
		config.Projections.Customizations["notifications"],
		config.Projections.Customizations["notificationsquotas"],
		config.Projections.Customizations["telemetry"],
		*config.Telemetry,
		config.ExternalDomain,
		config.ExternalPort,
		config.ExternalSecure,
		commands,
		queries,
		eventstoreClient,
		config.Login.DefaultOTPEmailURLV2,
		config.SystemDefaults.Notifications.FileSystemPath,
		keys.User,
		keys.SMTP,
		keys.SMS,
	)
	for _, p := range notify_handler.Projections() {
		err := migration.Migrate(ctx, eventstoreClient, p)
		logging.WithFields("name", p.String()).OnError(err).Fatal("migration failed")
	}
}
