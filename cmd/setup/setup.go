package setup

import (
	"context"
	"embed"
	_ "embed"
	"errors"
	"net/http"
	"path"

	"github.com/jackc/pgx/v5/pgconn"
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
	"github.com/zitadel/zitadel/internal/cache/connector"
	"github.com/zitadel/zitadel/internal/command"
	cryptoDB "github.com/zitadel/zitadel/internal/crypto/database"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_es "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	new_es "github.com/zitadel/zitadel/internal/eventstore/v3"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/migration"
	notify_handler "github.com/zitadel/zitadel/internal/notification"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/queue"
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

	dbClient, err := database.Connect(config.Database, false)
	logging.OnError(err).Fatal("unable to connect to database")

	config.Eventstore.Querier = old_es.NewCRDB(dbClient)
	esV3 := new_es.NewEventstore(dbClient)
	config.Eventstore.Pusher = esV3
	config.Eventstore.Searcher = esV3
	eventstoreClient := eventstore.NewEventstore(config.Eventstore)

	logging.OnError(err).Fatal("unable to start eventstore")
	eventstoreV4 := es_v4.NewEventstoreFromOne(es_v4_pg.New(dbClient, &es_v4_pg.Config{
		MaxRetries: config.Eventstore.MaxRetries,
	}))

	steps.s1ProjectionTable = &ProjectionTable{dbClient: dbClient.DB}
	steps.s2AssetsTable = &AssetTable{dbClient: dbClient.DB}

	steps.FirstInstance.Skip = config.ForMirror || steps.FirstInstance.Skip
	steps.FirstInstance.instanceSetup = config.DefaultInstance
	steps.FirstInstance.userEncryptionKey = config.EncryptionKeys.User
	steps.FirstInstance.smtpEncryptionKey = config.EncryptionKeys.SMTP
	steps.FirstInstance.oidcEncryptionKey = config.EncryptionKeys.OIDC
	steps.FirstInstance.masterKey = masterKey
	steps.FirstInstance.db = dbClient
	steps.FirstInstance.es = eventstoreClient
	steps.FirstInstance.defaults = config.SystemDefaults
	steps.FirstInstance.zitadelRoles = config.InternalAuthZ.RolePermissionMappings
	steps.FirstInstance.externalDomain = config.ExternalDomain
	steps.FirstInstance.externalSecure = config.ExternalSecure
	steps.FirstInstance.externalPort = config.ExternalPort

	steps.s5LastFailed = &LastFailed{dbClient: dbClient.DB}
	steps.s6OwnerRemoveColumns = &OwnerRemoveColumns{dbClient: dbClient.DB}
	steps.s7LogstoreTables = &LogstoreTables{dbClient: dbClient.DB, username: config.Database.Username(), dbType: config.Database.Type()}
	steps.s8AuthTokens = &AuthTokenIndexes{dbClient: dbClient}
	steps.CorrectCreationDate.dbClient = dbClient
	steps.s12AddOTPColumns = &AddOTPColumns{dbClient: dbClient}
	steps.s13FixQuotaProjection = &FixQuotaConstraints{dbClient: dbClient}
	steps.s14NewEventsTable = &NewEventsTable{dbClient: dbClient}
	steps.s15CurrentStates = &CurrentProjectionState{dbClient: dbClient}
	steps.s16UniqueConstraintsLower = &UniqueConstraintToLower{dbClient: dbClient}
	steps.s17AddOffsetToUniqueConstraints = &AddOffsetToCurrentStates{dbClient: dbClient}
	steps.s18AddLowerFieldsToLoginNames = &AddLowerFieldsToLoginNames{dbClient: dbClient}
	steps.s19AddCurrentStatesIndex = &AddCurrentSequencesIndex{dbClient: dbClient}
	steps.s20AddByUserSessionIndex = &AddByUserIndexToSession{dbClient: dbClient}
	steps.s21AddBlockFieldToLimits = &AddBlockFieldToLimits{dbClient: dbClient}
	steps.s22ActiveInstancesIndex = &ActiveInstanceEvents{dbClient: dbClient}
	steps.s23CorrectGlobalUniqueConstraints = &CorrectGlobalUniqueConstraints{dbClient: dbClient}
	steps.s24AddActorToAuthTokens = &AddActorToAuthTokens{dbClient: dbClient}
	steps.s25User11AddLowerFieldsToVerifiedEmail = &User11AddLowerFieldsToVerifiedEmail{dbClient: dbClient}
	steps.s26AuthUsers3 = &AuthUsers3{dbClient: dbClient}
	steps.s27IDPTemplate6SAMLNameIDFormat = &IDPTemplate6SAMLNameIDFormat{dbClient: dbClient}
	steps.s28AddFieldTable = &AddFieldTable{dbClient: dbClient}
	steps.s29FillFieldsForProjectGrant = &FillFieldsForProjectGrant{eventstore: eventstoreClient}
	steps.s30FillFieldsForOrgDomainVerified = &FillFieldsForOrgDomainVerified{eventstore: eventstoreClient}
	steps.s31AddAggregateIndexToFields = &AddAggregateIndexToFields{dbClient: dbClient}
	steps.s32AddAuthSessionID = &AddAuthSessionID{dbClient: dbClient}
	steps.s33SMSConfigs3TwilioAddVerifyServiceSid = &SMSConfigs3TwilioAddVerifyServiceSid{dbClient: dbClient}
	steps.s34AddCacheSchema = &AddCacheSchema{dbClient: dbClient}
	steps.s35AddPositionToIndexEsWm = &AddPositionToIndexEsWm{dbClient: dbClient}
	steps.s36FillV2Milestones = &FillV3Milestones{dbClient: dbClient, eventstore: eventstoreClient}
	steps.s37Apps7OIDConfigsBackChannelLogoutURI = &Apps7OIDConfigsBackChannelLogoutURI{dbClient: dbClient}
	steps.s38BackChannelLogoutNotificationStart = &BackChannelLogoutNotificationStart{dbClient: dbClient, esClient: eventstoreClient}
	steps.s40InitPushFunc = &InitPushFunc{dbClient: dbClient}
	steps.s42Apps7OIDCConfigsLoginVersion = &Apps7OIDCConfigsLoginVersion{dbClient: dbClient}
	steps.s43CreateFieldsDomainIndex = &CreateFieldsDomainIndex{dbClient: dbClient}
	steps.s44ReplaceCurrentSequencesIndex = &ReplaceCurrentSequencesIndex{dbClient: dbClient}
	steps.s45CorrectProjectOwners = &CorrectProjectOwners{eventstore: eventstoreClient}
	steps.s46InitPermissionFunctions = &InitPermissionFunctions{eventstoreClient: dbClient}
	steps.s47FillMembershipFields = &FillMembershipFields{eventstore: eventstoreClient}
	steps.s48Apps7SAMLConfigsLoginVersion = &Apps7SAMLConfigsLoginVersion{dbClient: dbClient}
	steps.s49InitPermittedOrgsFunction = &InitPermittedOrgsFunction{eventstoreClient: dbClient}
	steps.s50IDPTemplate6UsePKCE = &IDPTemplate6UsePKCE{dbClient: dbClient}
	steps.s51IDPTemplate6RootCA = &IDPTemplate6RootCA{dbClient: dbClient}

	err = projection.Create(ctx, dbClient, eventstoreClient, config.Projections, nil, nil, nil)
	logging.OnError(err).Fatal("unable to start projections")

	for _, step := range []migration.Migration{
		steps.s14NewEventsTable,
		steps.s40InitPushFunc,
		steps.s1ProjectionTable,
		steps.s2AssetsTable,
		steps.s28AddFieldTable,
		steps.s31AddAggregateIndexToFields,
		steps.s46InitPermissionFunctions,
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
		steps.s34AddCacheSchema,
		steps.s35AddPositionToIndexEsWm,
		steps.s36FillV2Milestones,
		steps.s38BackChannelLogoutNotificationStart,
		steps.s44ReplaceCurrentSequencesIndex,
		steps.s45CorrectProjectOwners,
		steps.s47FillMembershipFields,
		steps.s49InitPermittedOrgsFunction,
		steps.s50IDPTemplate6UsePKCE,
		steps.s51IDPTemplate6RootCA,
	} {
		mustExecuteMigration(ctx, eventstoreClient, step, "migration failed")
	}

	commands, _, _, _ := startCommandsQueries(ctx, eventstoreClient, eventstoreV4, dbClient, masterKey, config)

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
		&DeleteStaleOrgFields{
			eventstore: eventstoreClient,
		},
		&FillFieldsForInstanceDomains{
			eventstore: eventstoreClient,
		},
		&SyncRolePermissions{
			commands:               commands,
			eventstore:             eventstoreClient,
			rolePermissionMappings: config.InternalAuthZ.RolePermissionMappings,
		},
		&RiverMigrateRepeatable{
			client: dbClient,
		},
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
		steps.s32AddAuthSessionID,
		steps.s33SMSConfigs3TwilioAddVerifyServiceSid,
		steps.s37Apps7OIDConfigsBackChannelLogoutURI,
		steps.s42Apps7OIDCConfigsLoginVersion,
		steps.s43CreateFieldsDomainIndex,
		steps.s48Apps7SAMLConfigsLoginVersion,
	} {
		mustExecuteMigration(ctx, eventstoreClient, step, "migration failed")
	}

	// projection initialization must be done last, since the steps above might add required columns to the projections
	if !config.ForMirror && config.InitProjections.Enabled {
		initProjections(
			ctx,
			eventstoreClient,
		)
	}
}

func mustExecuteMigration(ctx context.Context, eventstoreClient *eventstore.Eventstore, step migration.Migration, errorMsg string) {
	err := migration.Migrate(ctx, eventstoreClient, step)
	if err == nil {
		return
	}
	logFields := []any{
		"name", step.String(),
	}
	pgErr := new(pgconn.PgError)
	if errors.As(err, &pgErr) {
		logFields = append(logFields,
			"severity", pgErr.Severity,
			"code", pgErr.Code,
			"message", pgErr.Message,
			"detail", pgErr.Detail,
			"hint", pgErr.Hint,
		)
	}
	logging.WithFields(logFields...).WithError(err).Fatal(errorMsg)
}

// readStmt reads a single file from the embedded FS,
// under the folder/typ/filename path.
// Typ describes the database dialect and may be omitted if no
// dialect specific migration is specified.
func readStmt(fs embed.FS, folder, typ, filename string) (string, error) {
	stmt, err := fs.ReadFile(path.Join(folder, typ, filename))
	return string(stmt), err
}

type statement struct {
	file  string
	query string
}

// readStatements reads all files from the embedded FS,
// under the folder/type path.
// Typ describes the database dialect and may be omitted if no
// dialect specific migration is specified.
func readStatements(fs embed.FS, folder, typ string) ([]statement, error) {
	basePath := path.Join(folder, typ)
	dir, err := fs.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	statements := make([]statement, len(dir))
	for i, file := range dir {
		statements[i].file = file.Name()
		statements[i].query, err = readStmt(fs, folder, typ, file.Name())
		if err != nil {
			return nil, err
		}
	}
	return statements, nil
}

func startCommandsQueries(
	ctx context.Context,
	eventstoreClient *eventstore.Eventstore,
	eventstoreV4 *es_v4.EventStore,
	dbClient *database.DB,
	masterKey string,
	config *Config,
) (
	*command.Commands,
	*query.Queries,
	*admin_view.View,
	*auth_view.View,
) {
	keyStorage, err := cryptoDB.NewKeyStorage(dbClient, masterKey)
	logging.OnError(err).Fatal("unable to start key storage")

	keys, err := encryption.EnsureEncryptionKeys(ctx, config.EncryptionKeys, keyStorage)
	logging.OnError(err).Fatal("unable to ensure encryption keys")

	err = projection.Create(
		ctx,
		dbClient,
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

	staticStorage, err := config.AssetStorage.NewStorage(dbClient.DB)
	logging.OnError(err).Fatal("unable to start asset storage")

	adminView, err := admin_view.StartView(dbClient)
	logging.OnError(err).Fatal("unable to start admin view")
	admin_handler.Register(ctx,
		admin_handler.Config{
			Client:                dbClient,
			Eventstore:            eventstoreClient,
			BulkLimit:             config.InitProjections.BulkLimit,
			FailureCountUntilSkip: uint64(config.InitProjections.MaxFailureCount),
		},
		adminView,
		staticStorage,
	)

	sessionTokenVerifier := internal_authz.SessionTokenVerifier(keys.OIDC)

	cacheConnectors, err := connector.StartConnectors(config.Caches, dbClient)
	logging.OnError(err).Fatal("unable to start caches")

	queries, err := query.StartQueries(
		ctx,
		eventstoreClient,
		eventstoreV4.Querier,
		dbClient,
		dbClient,
		cacheConnectors,
		config.Projections,
		config.SystemDefaults,
		keys.IDPConfig,
		keys.OTP,
		keys.OIDC,
		keys.SAML,
		keys.Target,
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

	authView, err := auth_view.StartView(dbClient, keys.OIDC, queries, eventstoreClient)
	logging.OnError(err).Fatal("unable to start admin view")
	auth_handler.Register(ctx,
		auth_handler.Config{
			Client:                dbClient,
			Eventstore:            eventstoreClient,
			BulkLimit:             config.InitProjections.BulkLimit,
			FailureCountUntilSkip: uint64(config.InitProjections.MaxFailureCount),
		},
		authView,
		queries,
	)

	authZRepo, err := authz.Start(queries, eventstoreClient, dbClient, keys.OIDC, config.ExternalSecure)
	logging.OnError(err).Fatal("unable to start authz repo")
	permissionCheck := func(ctx context.Context, permission, orgID, resourceID string) (err error) {
		return internal_authz.CheckPermission(ctx, authZRepo, config.InternalAuthZ.RolePermissionMappings, permission, orgID, resourceID)
	}

	commands, err := command.StartCommands(ctx,
		eventstoreClient,
		cacheConnectors,
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
		keys.Target,
		&http.Client{},
		permissionCheck,
		sessionTokenVerifier,
		config.OIDC.DefaultAccessTokenLifetime,
		config.OIDC.DefaultRefreshTokenExpiration,
		config.OIDC.DefaultRefreshTokenIdleExpiration,
		config.DefaultInstance.SecretGenerators,
	)
	logging.OnError(err).Fatal("unable to start commands")

	if !config.Notifications.LegacyEnabled && dbClient.Type() == "cockroach" {
		logging.Fatal("notifications must be set to LegacyEnabled=true when using CockroachDB")
	}
	q, err := queue.NewQueue(&queue.Config{
		Client: dbClient,
	})
	logging.OnError(err).Fatal("unable to init queue")

	notify_handler.Register(
		ctx,
		config.Projections.Customizations["notifications"],
		config.Projections.Customizations["notificationsquotas"],
		config.Projections.Customizations["backchannel"],
		config.Projections.Customizations["telemetry"],
		config.Notifications,
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
		keys.OIDC,
		config.OIDC.DefaultBackChannelLogoutLifetime,
		dbClient,
		q,
	)

	return commands, queries, adminView, authView
}

func initProjections(
	ctx context.Context,
	eventstoreClient *eventstore.Eventstore,
) {
	logging.Info("init-projections is currently in beta")

	for _, p := range projection.Projections() {
		err := migration.Migrate(ctx, eventstoreClient, p)
		logging.WithFields("name", p.String()).OnError(err).Fatal("migration failed")
	}

	for _, p := range admin_handler.Projections() {
		err := migration.Migrate(ctx, eventstoreClient, p)
		logging.WithFields("name", p.String()).OnError(err).Fatal("migration failed")
	}

	for _, p := range auth_handler.Projections() {
		err := migration.Migrate(ctx, eventstoreClient, p)
		logging.WithFields("name", p.String()).OnError(err).Fatal("migration failed")
	}

	for _, p := range notify_handler.Projections() {
		err := migration.Migrate(ctx, eventstoreClient, p)
		logging.WithFields("name", p.String()).OnError(err).Fatal("migration failed")
	}
}
