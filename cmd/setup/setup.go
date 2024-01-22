package setup

import (
	"context"
	"embed"
	_ "embed"
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
	"github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	iam_repo "github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/repository/restrictions"
	"github.com/zitadel/zitadel/internal/repository/session"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
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
	cmd.PersistentFlags().Bool("init-projections", viper.GetBool("init-projections"), "initializes projections after they are created")
	err := viper.BindPFlag("InitProjections", cmd.PersistentFlags().Lookup("init-projections"))
	logging.OnError(err).Fatal("unable to bind \"init-projections\" flag")
	key.AddMasterKeyFlag(cmd)
	tls.AddTLSModeFlag(cmd)
}

func Setup(config *Config, steps *Steps, masterKey string) {
	ctx := context.Background()
	logging.Info("setup started")

	i18n.MustLoadSupportedLanguagesFromDir()

	queryDBClient, err := database.Connect(config.Database, false, dialect.DBPurposeQuery)
	logging.OnError(err).Fatal("unable to connect to database")
	esPusherDBClient, err := database.Connect(config.Database, false, dialect.DBPurposeEventPusher)
	logging.OnError(err).Fatal("unable to connect to database")
	projectionDBClient, err := database.Connect(config.Database, false, dialect.DBPurposeProjectionSpooler)
	logging.OnError(err).Fatal("unable to connect to database")

	config.Eventstore.Querier = old_es.NewCRDB(queryDBClient)
	config.Eventstore.Pusher = new_es.NewEventstore(esPusherDBClient)
	eventstoreClient := eventstore.NewEventstore(config.Eventstore)
	logging.OnError(err).Fatal("unable to start eventstore")
	migration.RegisterMappers(eventstoreClient)

	steps.s1ProjectionTable = &ProjectionTable{dbClient: queryDBClient.DB}
	steps.s2AssetsTable = &AssetTable{dbClient: queryDBClient.DB}

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
	err = migration.Migrate(ctx, eventstoreClient, steps.s16UniqueConstraintsLower)
	logging.WithFields("name", steps.s16UniqueConstraintsLower.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.s17AddOffsetToUniqueConstraints)
	logging.WithFields("name", steps.s17AddOffsetToUniqueConstraints.String()).OnError(err).Fatal("migration failed")
	err = migration.Migrate(ctx, eventstoreClient, steps.s19AddCurrentStatesIndex)
	logging.WithFields("name", steps.s19AddCurrentStatesIndex.String()).OnError(err).Fatal("migration failed")

	for _, repeatableStep := range repeatableSteps {
		err = migration.Migrate(ctx, eventstoreClient, repeatableStep)
		logging.OnError(err).Fatalf("unable to migrate repeatable step: %s", repeatableStep.String())
	}

	if config.InitProjections {
		initProjections(
			ctx,
			eventstoreClient,
			queryDBClient,
			projectionDBClient,
			masterKey,
			config,
		)
	}

	// This step is executed after the repeatable steps because it adds fields to the login_names3 projection
	err = migration.Migrate(ctx, eventstoreClient, steps.s18AddLowerFieldsToLoginNames)
	logging.WithFields("name", steps.s18AddLowerFieldsToLoginNames.String()).OnError(err).Fatal("migration failed")
}

func readStmt(fs embed.FS, folder, typ, filename string) (string, error) {
	stmt, err := fs.ReadFile(folder + "/" + typ + "/" + filename)
	return string(stmt), err
}

func initProjections(
	ctx context.Context,
	eventstoreClient *eventstore.Eventstore,
	queryDBClient,
	projectionDBClient *database.DB,
	masterKey string,
	config *Config,
) {
	iam_repo.RegisterEventMappers(eventstoreClient)
	usr_repo.RegisterEventMappers(eventstoreClient)
	org.RegisterEventMappers(eventstoreClient)
	project.RegisterEventMappers(eventstoreClient)
	action.RegisterEventMappers(eventstoreClient)
	keypair.RegisterEventMappers(eventstoreClient)
	usergrant.RegisterEventMappers(eventstoreClient)
	session.RegisterEventMappers(eventstoreClient)
	idpintent.RegisterEventMappers(eventstoreClient)
	authrequest.RegisterEventMappers(eventstoreClient)
	oidcsession.RegisterEventMappers(eventstoreClient)
	quota.RegisterEventMappers(eventstoreClient)
	limits.RegisterEventMappers(eventstoreClient)
	restrictions.RegisterEventMappers(eventstoreClient)
	deviceauth.RegisterEventMappers(eventstoreClient)

	keyStorage, err := cryptoDB.NewKeyStorage(queryDBClient, masterKey)
	logging.OnError(err).Fatal("unable to start key storage")

	keys, err := encryption.EnsureEncryptionKeys(ctx, config.EncryptionKeys, keyStorage)
	logging.OnError(err).Fatal("unable to ensure encryption keys")

	err = projection.Create(
		ctx,
		queryDBClient,
		eventstoreClient,
		config.Projections,
		keys.OIDC,
		keys.SAML,
		nil, // system users are only used for milestone checks
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
	config.Admin.Spooler.Eventstore = eventstoreClient
	config.Admin.Spooler.Client = queryDBClient
	admin_handler.Register(ctx,
		config.Admin.Spooler,
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
	)
	logging.OnError(err).Fatal("unable to start queries")

	config.Auth.Spooler.Eventstore = eventstoreClient
	config.Auth.Spooler.Client = queryDBClient
	authView, err := auth_view.StartView(queryDBClient, keys.OIDC, queries, eventstoreClient)
	logging.OnError(err).Fatal("unable to start admin view")
	auth_handler.Register(ctx,
		config.Auth.Spooler,
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
