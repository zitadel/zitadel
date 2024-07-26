package start

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	clockpkg "github.com/benbjohnson/clock"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/op"
	"github.com/zitadel/saml/pkg/provider"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/cmd/encryption"
	"github.com/zitadel/zitadel/cmd/key"
	cmd_tls "github.com/zitadel/zitadel/cmd/tls"
	"github.com/zitadel/zitadel/internal/actions"
	admin_es "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/api"
	"github.com/zitadel/zitadel/internal/api/assets"
	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	action_v3_alpha "github.com/zitadel/zitadel/internal/api/grpc/action/v3alpha"
	"github.com/zitadel/zitadel/internal/api/grpc/admin"
	"github.com/zitadel/zitadel/internal/api/grpc/auth"
	feature_v2 "github.com/zitadel/zitadel/internal/api/grpc/feature/v2"
	feature_v2beta "github.com/zitadel/zitadel/internal/api/grpc/feature/v2beta"
	"github.com/zitadel/zitadel/internal/api/grpc/management"
	oidc_v2 "github.com/zitadel/zitadel/internal/api/grpc/oidc/v2"
	oidc_v2beta "github.com/zitadel/zitadel/internal/api/grpc/oidc/v2beta"
	org_v2 "github.com/zitadel/zitadel/internal/api/grpc/org/v2"
	org_v2beta "github.com/zitadel/zitadel/internal/api/grpc/org/v2beta"
	session_v2 "github.com/zitadel/zitadel/internal/api/grpc/session/v2"
	session_v2beta "github.com/zitadel/zitadel/internal/api/grpc/session/v2beta"
	settings_v2 "github.com/zitadel/zitadel/internal/api/grpc/settings/v2"
	settings_v2beta "github.com/zitadel/zitadel/internal/api/grpc/settings/v2beta"
	"github.com/zitadel/zitadel/internal/api/grpc/system"
	user_schema_v3_alpha "github.com/zitadel/zitadel/internal/api/grpc/user/schema/v3alpha"
	user_v2 "github.com/zitadel/zitadel/internal/api/grpc/user/v2"
	user_v2beta "github.com/zitadel/zitadel/internal/api/grpc/user/v2beta"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/api/idp"
	"github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/api/robots_txt"
	"github.com/zitadel/zitadel/internal/api/saml"
	"github.com/zitadel/zitadel/internal/api/ui/console"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	auth_es "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/authz"
	authz_repo "github.com/zitadel/zitadel/internal/authz/repository"
	authz_es "github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/eventstore"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	cryptoDB "github.com/zitadel/zitadel/internal/crypto/database"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_es "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	new_es "github.com/zitadel/zitadel/internal/eventstore/v3"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/emitters/access"
	"github.com/zitadel/zitadel/internal/logstore/emitters/execution"
	"github.com/zitadel/zitadel/internal/logstore/emitters/stdout"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/net"
	"github.com/zitadel/zitadel/internal/notification"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/static"
	es_v4 "github.com/zitadel/zitadel/internal/v2/eventstore"
	es_v4_pg "github.com/zitadel/zitadel/internal/v2/eventstore/postgres"
	"github.com/zitadel/zitadel/internal/webauthn"
	"github.com/zitadel/zitadel/openapi"
)

func New(server chan<- *Server) *cobra.Command {
	start := &cobra.Command{
		Use:   "start",
		Short: "starts ZITADEL instance",
		Long: `starts ZITADEL.
Requirements:
- cockroachdb`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd_tls.ModeFromFlag(cmd)
			if err != nil {
				return err
			}
			config := MustNewConfig(viper.GetViper())
			masterKey, err := key.MasterKey(cmd)
			if err != nil {
				return err
			}
			return startZitadel(cmd.Context(), config, masterKey, server)
		},
	}

	startFlags(start)

	return start
}

type Server struct {
	Config     *Config
	DB         *database.DB
	KeyStorage crypto.KeyStorage
	Keys       *encryption.EncryptionKeys
	Eventstore *eventstore.Eventstore
	Queries    *query.Queries
	AuthzRepo  authz_repo.Repository
	Storage    static.Storage
	Commands   *command.Commands
	Router     *mux.Router
	TLSConfig  *tls.Config
	Shutdown   chan<- os.Signal
}

func startZitadel(ctx context.Context, config *Config, masterKey string, server chan<- *Server) error {
	showBasicInformation(config)

	i18n.MustLoadSupportedLanguagesFromDir()

	queryDBClient, err := database.Connect(config.Database, false, dialect.DBPurposeQuery)
	if err != nil {
		return fmt.Errorf("cannot start DB client for queries: %w", err)
	}
	esPusherDBClient, err := database.Connect(config.Database, false, dialect.DBPurposeEventPusher)
	if err != nil {
		return fmt.Errorf("cannot start client for event store pusher: %w", err)
	}
	projectionDBClient, err := database.Connect(config.Database, false, dialect.DBPurposeProjectionSpooler)
	if err != nil {
		return fmt.Errorf("cannot start client for projection spooler: %w", err)
	}

	keyStorage, err := cryptoDB.NewKeyStorage(queryDBClient, masterKey)
	if err != nil {
		return fmt.Errorf("cannot start key storage: %w", err)
	}
	keys, err := encryption.EnsureEncryptionKeys(ctx, config.EncryptionKeys, keyStorage)
	if err != nil {
		return err
	}

	config.Eventstore.Pusher = new_es.NewEventstore(esPusherDBClient)
	config.Eventstore.Searcher = new_es.NewEventstore(queryDBClient)
	config.Eventstore.Querier = old_es.NewCRDB(queryDBClient)
	eventstoreClient := eventstore.NewEventstore(config.Eventstore)
	eventstoreV4 := es_v4.NewEventstoreFromOne(es_v4_pg.New(queryDBClient, &es_v4_pg.Config{
		MaxRetries: config.Eventstore.MaxRetries,
	}))

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
		config.AuditLogRetention,
		config.SystemAPIUsers,
		true,
	)
	if err != nil {
		return fmt.Errorf("cannot start queries: %w", err)
	}

	authZRepo, err := authz.Start(queries, eventstoreClient, queryDBClient, keys.OIDC, config.ExternalSecure)
	if err != nil {
		return fmt.Errorf("error starting authz repo: %w", err)
	}
	permissionCheck := func(ctx context.Context, permission, orgID, resourceID string) (err error) {
		return internal_authz.CheckPermission(ctx, authZRepo, config.InternalAuthZ.RolePermissionMappings, permission, orgID, resourceID)
	}

	storage, err := config.AssetStorage.NewStorage(queryDBClient.DB)
	if err != nil {
		return fmt.Errorf("cannot start asset storage client: %w", err)
	}
	webAuthNConfig := &webauthn.Config{
		DisplayName:    config.WebAuthNName,
		ExternalSecure: config.ExternalSecure,
	}
	commands, err := command.StartCommands(
		eventstoreClient,
		config.SystemDefaults,
		config.InternalAuthZ.RolePermissionMappings,
		storage,
		webAuthNConfig,
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
	if err != nil {
		return fmt.Errorf("cannot start commands: %w", err)
	}
	defer commands.Close(ctx) // wait for background jobs

	clock := clockpkg.New()
	actionsExecutionStdoutEmitter, err := logstore.NewEmitter[*record.ExecutionLog](ctx, clock, &logstore.EmitterConfig{Enabled: config.LogStore.Execution.Stdout.Enabled}, stdout.NewStdoutEmitter[*record.ExecutionLog]())
	if err != nil {
		return err
	}
	actionsExecutionDBEmitter, err := logstore.NewEmitter[*record.ExecutionLog](ctx, clock, config.Quotas.Execution, execution.NewDatabaseLogStorage(queryDBClient, commands, queries))
	if err != nil {
		return err
	}

	actionsLogstoreSvc := logstore.New(queries, actionsExecutionDBEmitter, actionsExecutionStdoutEmitter)
	actions.SetLogstoreService(actionsLogstoreSvc)

	notification.Register(
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
	notification.Start(ctx)

	router := mux.NewRouter()
	tlsConfig, err := config.TLS.Config()
	if err != nil {
		return err
	}
	api, err := startAPIs(
		ctx,
		clock,
		router,
		commands,
		queries,
		eventstoreClient,
		queryDBClient,
		config,
		storage,
		authZRepo,
		keys,
		permissionCheck,
	)
	if err != nil {
		return err
	}
	commands.GrpcMethodExisting = checkExisting(api.ListGrpcMethods())
	commands.GrpcServiceExisting = checkExisting(api.ListGrpcServices())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	if server != nil {
		server <- &Server{
			Config:     config,
			DB:         queryDBClient,
			KeyStorage: keyStorage,
			Keys:       keys,
			Eventstore: eventstoreClient,
			Queries:    queries,
			AuthzRepo:  authZRepo,
			Storage:    storage,
			Commands:   commands,
			Router:     router,
			TLSConfig:  tlsConfig,
			Shutdown:   shutdown,
		}
		close(server)
	}

	return listen(ctx, router, config.Port, tlsConfig, shutdown)
}

func startAPIs(
	ctx context.Context,
	clock clockpkg.Clock,
	router *mux.Router,
	commands *command.Commands,
	queries *query.Queries,
	eventstore *eventstore.Eventstore,
	dbClient *database.DB,
	config *Config,
	store static.Storage,
	authZRepo authz_repo.Repository,
	keys *encryption.EncryptionKeys,
	permissionCheck domain.PermissionCheck,
) (*api.API, error) {
	repo := struct {
		authz_repo.Repository
		*query.Queries
	}{
		authZRepo,
		queries,
	}
	oidcPrefixes := []string{"/.well-known/openid-configuration", "/oidc/v1", "/oauth/v2"}
	// always set the origin in the context if available in the http headers, no matter for what protocol
	router.Use(middleware.WithOrigin(config.ExternalSecure))
	systemTokenVerifier, err := internal_authz.StartSystemTokenVerifierFromConfig(http_util.BuildHTTP(config.ExternalDomain, config.ExternalPort, config.ExternalSecure), config.SystemAPIUsers)
	if err != nil {
		return nil, err
	}
	accessTokenVerifer := internal_authz.StartAccessTokenVerifierFromRepo(repo)
	verifier := internal_authz.StartAPITokenVerifier(repo, accessTokenVerifer, systemTokenVerifier)
	tlsConfig, err := config.TLS.Config()
	if err != nil {
		return nil, err
	}

	accessStdoutEmitter, err := logstore.NewEmitter[*record.AccessLog](ctx, clock, &logstore.EmitterConfig{Enabled: config.LogStore.Access.Stdout.Enabled}, stdout.NewStdoutEmitter[*record.AccessLog]())
	if err != nil {
		return nil, err
	}
	accessDBEmitter, err := logstore.NewEmitter[*record.AccessLog](ctx, clock, &config.Quotas.Access.EmitterConfig, access.NewDatabaseLogStorage(dbClient, commands, queries))
	if err != nil {
		return nil, err
	}

	accessSvc := logstore.New[*record.AccessLog](queries, accessDBEmitter, accessStdoutEmitter)
	exhaustedCookieHandler := http_util.NewCookieHandler(
		http_util.WithUnsecure(),
		http_util.WithNonHttpOnly(),
		http_util.WithMaxAge(int(math.Floor(config.Quotas.Access.ExhaustedCookieMaxAge.Seconds()))),
	)
	limitingAccessInterceptor := middleware.NewAccessInterceptor(accessSvc, exhaustedCookieHandler, &config.Quotas.Access.AccessConfig)
	apis, err := api.New(ctx, config.Port, router, queries, verifier, config.InternalAuthZ, tlsConfig, config.HTTP2HostHeader, config.HTTP1HostHeader, config.ExternalDomain, limitingAccessInterceptor)
	if err != nil {
		return nil, fmt.Errorf("error creating api %w", err)
	}

	config.Auth.Spooler.Client = dbClient
	config.Auth.Spooler.Eventstore = eventstore
	authRepo, err := auth_es.Start(ctx, config.Auth, config.SystemDefaults, commands, queries, dbClient, eventstore, keys.OIDC, keys.User)
	if err != nil {
		return nil, fmt.Errorf("error starting auth repo: %w", err)
	}

	config.Admin.Spooler.Client = dbClient
	config.Admin.Spooler.Eventstore = eventstore
	err = admin_es.Start(ctx, config.Admin, store, dbClient)
	if err != nil {
		return nil, fmt.Errorf("error starting admin repo: %w", err)
	}

	if err := apis.RegisterServer(ctx, system.CreateServer(commands, queries, config.Database.DatabaseName(), config.DefaultInstance, config.ExternalDomain), tlsConfig); err != nil {
		return nil, err
	}
	if err := apis.RegisterServer(ctx, admin.CreateServer(config.Database.DatabaseName(), commands, queries, config.SystemDefaults, config.ExternalSecure, keys.User, config.AuditLogRetention), tlsConfig); err != nil {
		return nil, err
	}
	if err := apis.RegisterServer(ctx, management.CreateServer(commands, queries, config.SystemDefaults, keys.User, config.ExternalSecure), tlsConfig); err != nil {
		return nil, err
	}
	if err := apis.RegisterServer(ctx, auth.CreateServer(commands, queries, authRepo, config.SystemDefaults, keys.User, config.ExternalSecure), tlsConfig); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, user_v2beta.CreateServer(commands, queries, keys.User, keys.IDPConfig, idp.CallbackURL(config.ExternalSecure), idp.SAMLRootURL(config.ExternalSecure), assets.AssetAPI(config.ExternalSecure), permissionCheck)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, user_v2.CreateServer(commands, queries, keys.User, keys.IDPConfig, idp.CallbackURL(config.ExternalSecure), idp.SAMLRootURL(config.ExternalSecure), assets.AssetAPI(config.ExternalSecure), permissionCheck)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, session_v2beta.CreateServer(commands, queries)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, settings_v2beta.CreateServer(commands, queries, config.ExternalSecure)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, org_v2beta.CreateServer(commands, queries, permissionCheck)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, feature_v2beta.CreateServer(commands, queries)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, session_v2.CreateServer(commands, queries)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, settings_v2.CreateServer(commands, queries, config.ExternalSecure)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, org_v2.CreateServer(commands, queries, permissionCheck)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, feature_v2.CreateServer(commands, queries)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, action_v3_alpha.CreateServer(commands, queries, domain.AllFunctions, apis.ListGrpcMethods, apis.ListGrpcServices)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, user_schema_v3_alpha.CreateServer(commands, queries)); err != nil {
		return nil, err
	}
	instanceInterceptor := middleware.InstanceInterceptor(queries, config.HTTP1HostHeader, config.ExternalDomain, login.IgnoreInstanceEndpoints...)
	assetsCache := middleware.AssetsCacheInterceptor(config.AssetStorage.Cache.MaxAge, config.AssetStorage.Cache.SharedMaxAge)
	apis.RegisterHandlerOnPrefix(assets.HandlerPrefix, assets.NewHandler(commands, verifier, config.InternalAuthZ, id.SonyFlakeGenerator(), store, queries, middleware.CallDurationHandler, instanceInterceptor.Handler, assetsCache.Handler, limitingAccessInterceptor.Handle))

	apis.RegisterHandlerOnPrefix(idp.HandlerPrefix, idp.NewHandler(commands, queries, keys.IDPConfig, config.ExternalSecure, instanceInterceptor.Handler))

	userAgentInterceptor, err := middleware.NewUserAgentHandler(config.UserAgentCookie, keys.UserAgentCookieKey, id.SonyFlakeGenerator(), config.ExternalSecure, login.EndpointResources, login.EndpointExternalLoginCallbackFormPost, login.EndpointSAMLACS)
	if err != nil {
		return nil, err
	}

	// robots.txt handler
	robotsTxtHandler, err := robots_txt.Start()
	if err != nil {
		return nil, fmt.Errorf("unable to start robots txt handler: %w", err)
	}
	apis.RegisterHandlerOnPrefix(robots_txt.HandlerPrefix, robotsTxtHandler)

	// TODO: Record openapi access logs?
	openAPIHandler, err := openapi.Start()
	if err != nil {
		return nil, fmt.Errorf("unable to start openapi handler: %w", err)
	}
	apis.RegisterHandlerOnPrefix(openapi.HandlerPrefix, openAPIHandler)

	oidcServer, err := oidc.NewServer(ctx, config.OIDC, login.DefaultLoggedOutPath, config.ExternalSecure, commands, queries, authRepo, keys.OIDC, keys.OIDCKey, eventstore, dbClient, userAgentInterceptor, instanceInterceptor.Handler, limitingAccessInterceptor, config.Log.Slog(), config.SystemDefaults.SecretHasher)
	if err != nil {
		return nil, fmt.Errorf("unable to start oidc provider: %w", err)
	}
	apis.RegisterHandlerPrefixes(oidcServer, oidcPrefixes...)

	samlProvider, err := saml.NewProvider(config.SAML, config.ExternalSecure, commands, queries, authRepo, keys.OIDC, keys.SAML, eventstore, dbClient, instanceInterceptor.Handler, userAgentInterceptor, limitingAccessInterceptor)
	if err != nil {
		return nil, fmt.Errorf("unable to start saml provider: %w", err)
	}
	apis.RegisterHandlerOnPrefix(saml.HandlerPrefix, samlProvider.HttpHandler())

	c, err := console.Start(config.Console, config.ExternalSecure, oidcServer.IssuerFromRequest, middleware.CallDurationHandler, instanceInterceptor.Handler, limitingAccessInterceptor, config.CustomerPortal)
	if err != nil {
		return nil, fmt.Errorf("unable to start console: %w", err)
	}
	apis.RegisterHandlerOnPrefix(console.HandlerPrefix, c)
	consolePath := console.HandlerPrefix + "/"
	l, err := login.CreateLogin(
		config.Login,
		commands,
		queries,
		authRepo,
		store,
		consolePath,
		oidcServer.AuthCallbackURL(),
		provider.AuthCallbackURL(samlProvider),
		config.ExternalSecure,
		userAgentInterceptor,
		op.NewIssuerInterceptor(oidcServer.IssuerFromRequest).Handler,
		provider.NewIssuerInterceptor(samlProvider.IssuerFromRequest).Handler,
		instanceInterceptor.Handler,
		assetsCache.Handler,
		limitingAccessInterceptor.WithRedirect(consolePath).Handle,
		keys.User,
		keys.IDPConfig,
		keys.CSRFCookieKey,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to start login: %w", err)
	}
	apis.RegisterHandlerOnPrefix(login.HandlerPrefix, l.Handler())
	apis.HandleFunc(login.EndpointDeviceAuth, login.RedirectDeviceAuthToPrefix)

	// After OIDC provider so that the callback endpoint can be used
	if err := apis.RegisterService(ctx, oidc_v2beta.CreateServer(commands, queries, oidcServer, config.ExternalSecure)); err != nil {
		return nil, err
	}
	if err := apis.RegisterService(ctx, oidc_v2.CreateServer(commands, queries, oidcServer, config.ExternalSecure)); err != nil {
		return nil, err
	}
	// handle grpc at last to be able to handle the root, because grpc and gateway require a lot of different prefixes
	apis.RouteGRPC()
	return apis, nil
}

func listen(ctx context.Context, router *mux.Router, port uint16, tlsConfig *tls.Config, shutdown <-chan os.Signal) error {
	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(router, http2Server), TLSConfig: tlsConfig}

	lc := net.ListenConfig()
	lis, err := lc.Listen(ctx, "tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("tcp listener on %d failed: %w", port, err)
	}

	errCh := make(chan error)

	go func() {
		logging.Infof("server is listening on %s", lis.Addr().String())
		if tlsConfig != nil {
			//we don't need to pass the files here, because we already initialized the TLS config on the server
			errCh <- http1Server.ServeTLS(lis, "", "")
		} else {
			errCh <- http1Server.Serve(lis)
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("error starting server: %w", err)
	case <-shutdown:
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return shutdownServer(ctx, http1Server)
	case <-ctx.Done():
		return shutdownServer(ctx, http1Server)
	}
}

func shutdownServer(ctx context.Context, server *http.Server) error {
	err := server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("could not shutdown gracefully: %w", err)
	}
	logging.New().Info("server shutdown gracefully")
	return nil
}

func showBasicInformation(startConfig *Config) {
	fmt.Println(color.MagentaString(figure.NewFigure("ZITADEL", "", true).String()))
	http := "http"
	if startConfig.TLS.Enabled || startConfig.ExternalSecure {
		http = "https"
	}

	consoleURL := fmt.Sprintf("%s://%s:%v/ui/console\n", http, startConfig.ExternalDomain, startConfig.ExternalPort)
	healthCheckURL := fmt.Sprintf("%s://%s:%v/debug/healthz\n", http, startConfig.ExternalDomain, startConfig.ExternalPort)
	machineIdMethod := id.MachineIdentificationMethod()

	insecure := !startConfig.TLS.Enabled && !startConfig.ExternalSecure

	fmt.Printf(" ===============================================================\n\n")
	fmt.Printf(" Version          	: %s\n", build.Version())
	fmt.Printf(" TLS enabled      	: %v\n", startConfig.TLS.Enabled)
	fmt.Printf(" External Secure 	: %v\n", startConfig.ExternalSecure)
	fmt.Printf(" Machine Id Method	: %v\n", machineIdMethod)
	fmt.Printf(" Console URL      	: %s", color.BlueString(consoleURL))
	fmt.Printf(" Health Check URL 	: %s", color.BlueString(healthCheckURL))
	if insecure {
		fmt.Printf("\n %s: you're using plain http without TLS. Be aware this is \n", color.RedString("Warning"))
		fmt.Printf(" not a secure setup and should only be used for test systems.         \n")
		fmt.Printf(" Visit: %s    \n", color.CyanString("https://zitadel.com/docs/self-hosting/manage/tls_modes"))
	}
	fmt.Printf("\n ===============================================================\n\n")
}

func checkExisting(values []string) func(string) bool {
	return func(value string) bool {
		return slices.Contains(values, value)
	}
}
