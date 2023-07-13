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
	"syscall"
	"time"

	clockpkg "github.com/benbjohnson/clock"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/op"
	"github.com/zitadel/saml/pkg/provider"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/cmd/key"
	cmd_tls "github.com/zitadel/zitadel/cmd/tls"
	"github.com/zitadel/zitadel/internal/actions"
	admin_es "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/api"
	"github.com/zitadel/zitadel/internal/api/assets"
	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/admin"
	"github.com/zitadel/zitadel/internal/api/grpc/auth"
	"github.com/zitadel/zitadel/internal/api/grpc/management"
	oidc_v2 "github.com/zitadel/zitadel/internal/api/grpc/oidc/v2"
	"github.com/zitadel/zitadel/internal/api/grpc/session/v2"
	"github.com/zitadel/zitadel/internal/api/grpc/settings/v2"
	"github.com/zitadel/zitadel/internal/api/grpc/system"
	"github.com/zitadel/zitadel/internal/api/grpc/user/v2"
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
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/emitters/access"
	"github.com/zitadel/zitadel/internal/logstore/emitters/execution"
	"github.com/zitadel/zitadel/internal/logstore/emitters/stdout"
	"github.com/zitadel/zitadel/internal/notification"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/static"
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

			return startZitadel(config, masterKey, server)
		},
	}

	startFlags(start)

	return start
}

type Server struct {
	Config     *Config
	DB         *database.DB
	KeyStorage crypto.KeyStorage
	Keys       *encryptionKeys
	Eventstore *eventstore.Eventstore
	Queries    *query.Queries
	AuthzRepo  authz_repo.Repository
	Storage    static.Storage
	Commands   *command.Commands
	LogStore   *logstore.Service
	Router     *mux.Router
	TLSConfig  *tls.Config
	Shutdown   chan<- os.Signal
}

func startZitadel(config *Config, masterKey string, server chan<- *Server) error {
	showBasicInformation(config)

	ctx := context.Background()

	dbClient, err := database.Connect(config.Database, false)
	if err != nil {
		return fmt.Errorf("cannot start client for projection: %w", err)
	}

	keyStorage, err := cryptoDB.NewKeyStorage(dbClient.DB, masterKey)
	if err != nil {
		return fmt.Errorf("cannot start key storage: %w", err)
	}
	keys, err := ensureEncryptionKeys(config.EncryptionKeys, keyStorage)
	if err != nil {
		return err
	}

	config.Eventstore.Client = dbClient
	eventstoreClient, err := eventstore.Start(config.Eventstore)
	if err != nil {
		return fmt.Errorf("cannot start eventstore for queries: %w", err)
	}

	sessionTokenVerifier := internal_authz.SessionTokenVerifier(keys.OIDC)

	queries, err := query.StartQueries(
		ctx,
		eventstoreClient,
		dbClient,
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
	)
	if err != nil {
		return fmt.Errorf("cannot start queries: %w", err)
	}

	authZRepo, err := authz.Start(queries, dbClient, keys.OIDC, config.ExternalSecure, config.Eventstore.AllowOrderByCreationDate)
	if err != nil {
		return fmt.Errorf("error starting authz repo: %w", err)
	}
	permissionCheck := func(ctx context.Context, permission, orgID, resourceID string) (err error) {
		return internal_authz.CheckPermission(ctx, authZRepo, config.InternalAuthZ.RolePermissionMappings, permission, orgID, resourceID)
	}

	storage, err := config.AssetStorage.NewStorage(dbClient.DB)
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
	)
	if err != nil {
		return fmt.Errorf("cannot start commands: %w", err)
	}

	clock := clockpkg.New()
	actionsExecutionStdoutEmitter, err := logstore.NewEmitter(ctx, clock, config.LogStore.Execution.Stdout, stdout.NewStdoutEmitter())
	if err != nil {
		return err
	}
	actionsExecutionDBEmitter, err := logstore.NewEmitter(ctx, clock, config.LogStore.Execution.Database, execution.NewDatabaseLogStorage(dbClient))
	if err != nil {
		return err
	}

	usageReporter := logstore.UsageReporterFunc(commands.ReportQuotaUsage)
	actionsLogstoreSvc := logstore.New(queries, usageReporter, actionsExecutionDBEmitter, actionsExecutionStdoutEmitter)
	if actionsLogstoreSvc.Enabled() {
		logging.Warn("execution logs are currently in beta")
	}
	actions.SetLogstoreService(actionsLogstoreSvc)

	notification.Start(ctx, config.Projections.Customizations["notifications"], config.Projections.Customizations["notificationsquotas"], config.Projections.Customizations["telemetry"], *config.Telemetry, config.ExternalDomain, config.ExternalPort, config.ExternalSecure, commands, queries, eventstoreClient, assets.AssetAPIFromDomain(config.ExternalSecure, config.ExternalPort), config.SystemDefaults.Notifications.FileSystemPath, keys.User, keys.SMTP, keys.SMS)

	router := mux.NewRouter()
	tlsConfig, err := config.TLS.Config()
	if err != nil {
		return err
	}
	err = startAPIs(
		ctx,
		clock,
		router,
		commands,
		queries,
		eventstoreClient,
		dbClient,
		config,
		storage,
		authZRepo,
		keys,
		queries,
		usageReporter,
		permissionCheck,
	)
	if err != nil {
		return err
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	if server != nil {
		server <- &Server{
			Config:     config,
			DB:         dbClient,
			KeyStorage: keyStorage,
			Keys:       keys,
			Eventstore: eventstoreClient,
			Queries:    queries,
			AuthzRepo:  authZRepo,
			Storage:    storage,
			Commands:   commands,
			LogStore:   actionsLogstoreSvc,
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
	keys *encryptionKeys,
	quotaQuerier logstore.QuotaQuerier,
	usageReporter logstore.UsageReporter,
	permissionCheck domain.PermissionCheck,
) error {
	repo := struct {
		authz_repo.Repository
		*query.Queries
	}{
		authZRepo,
		queries,
	}
	verifier := internal_authz.Start(repo, http_util.BuildHTTP(config.ExternalDomain, config.ExternalPort, config.ExternalSecure), config.SystemAPIUsers)
	tlsConfig, err := config.TLS.Config()
	if err != nil {
		return err
	}

	accessStdoutEmitter, err := logstore.NewEmitter(ctx, clock, config.LogStore.Access.Stdout, stdout.NewStdoutEmitter())
	if err != nil {
		return err
	}
	accessDBEmitter, err := logstore.NewEmitter(ctx, clock, config.LogStore.Access.Database, access.NewDatabaseLogStorage(dbClient))
	if err != nil {
		return err
	}

	accessSvc := logstore.New(quotaQuerier, usageReporter, accessDBEmitter, accessStdoutEmitter)
	if accessSvc.Enabled() {
		logging.Warn("access logs are currently in beta")
	}
	exhaustedCookieHandler := http_util.NewCookieHandler(
		http_util.WithUnsecure(),
		http_util.WithNonHttpOnly(),
		http_util.WithMaxAge(int(math.Floor(config.Quotas.Access.ExhaustedCookieMaxAge.Seconds()))),
	)
	limitingAccessInterceptor := middleware.NewAccessInterceptor(accessSvc, exhaustedCookieHandler, config.Quotas.Access)
	apis, err := api.New(ctx, config.Port, router, queries, verifier, config.InternalAuthZ, tlsConfig, config.HTTP2HostHeader, config.HTTP1HostHeader, limitingAccessInterceptor)
	if err != nil {
		return fmt.Errorf("error creating api %w", err)
	}
	authRepo, err := auth_es.Start(ctx, config.Auth, config.SystemDefaults, commands, queries, dbClient, eventstore, keys.OIDC, keys.User, config.Eventstore.AllowOrderByCreationDate)
	if err != nil {
		return fmt.Errorf("error starting auth repo: %w", err)
	}
	adminRepo, err := admin_es.Start(ctx, config.Admin, store, dbClient, eventstore, config.Eventstore.AllowOrderByCreationDate)
	if err != nil {
		return fmt.Errorf("error starting admin repo: %w", err)
	}
	if err := apis.RegisterServer(ctx, system.CreateServer(commands, queries, adminRepo, config.Database.DatabaseName(), config.DefaultInstance, config.ExternalDomain)); err != nil {
		return err
	}
	if err := apis.RegisterServer(ctx, admin.CreateServer(config.Database.DatabaseName(), commands, queries, config.SystemDefaults, adminRepo, config.ExternalSecure, keys.User, config.AuditLogRetention)); err != nil {
		return err
	}
	if err := apis.RegisterServer(ctx, management.CreateServer(commands, queries, config.SystemDefaults, keys.User, config.ExternalSecure, config.AuditLogRetention)); err != nil {
		return err
	}
	if err := apis.RegisterServer(ctx, auth.CreateServer(commands, queries, authRepo, config.SystemDefaults, keys.User, config.ExternalSecure, config.AuditLogRetention)); err != nil {
		return err
	}
	if err := apis.RegisterService(ctx, user.CreateServer(commands, queries, keys.User, keys.IDPConfig, idp.CallbackURL(config.ExternalSecure))); err != nil {
		return err
	}
	if err := apis.RegisterService(ctx, session.CreateServer(commands, queries, permissionCheck)); err != nil {
		return err
	}

	if err := apis.RegisterService(ctx, settings.CreateServer(commands, queries, config.ExternalSecure)); err != nil {
		return err
	}
	instanceInterceptor := middleware.InstanceInterceptor(queries, config.HTTP1HostHeader, login.IgnoreInstanceEndpoints...)
	assetsCache := middleware.AssetsCacheInterceptor(config.AssetStorage.Cache.MaxAge, config.AssetStorage.Cache.SharedMaxAge)
	apis.RegisterHandlerOnPrefix(assets.HandlerPrefix, assets.NewHandler(commands, verifier, config.InternalAuthZ, id.SonyFlakeGenerator(), store, queries, middleware.CallDurationHandler, instanceInterceptor.Handler, assetsCache.Handler, limitingAccessInterceptor.Handle))

	apis.RegisterHandlerOnPrefix(idp.HandlerPrefix, idp.NewHandler(commands, queries, keys.IDPConfig, config.ExternalSecure, instanceInterceptor.Handler))

	userAgentInterceptor, err := middleware.NewUserAgentHandler(config.UserAgentCookie, keys.UserAgentCookieKey, id.SonyFlakeGenerator(), config.ExternalSecure, login.EndpointResources)
	if err != nil {
		return err
	}

	// robots.txt handler
	robotsTxtHandler, err := robots_txt.Start()
	if err != nil {
		return fmt.Errorf("unable to start robots txt handler: %w", err)
	}
	apis.RegisterHandlerOnPrefix(robots_txt.HandlerPrefix, robotsTxtHandler)

	// TODO: Record openapi access logs?
	openAPIHandler, err := openapi.Start()
	if err != nil {
		return fmt.Errorf("unable to start openapi handler: %w", err)
	}
	apis.RegisterHandlerOnPrefix(openapi.HandlerPrefix, openAPIHandler)

	oidcProvider, err := oidc.NewProvider(config.OIDC, login.DefaultLoggedOutPath, config.ExternalSecure, commands, queries, authRepo, keys.OIDC, keys.OIDCKey, eventstore, dbClient, userAgentInterceptor, instanceInterceptor.Handler, limitingAccessInterceptor.Handle)
	if err != nil {
		return fmt.Errorf("unable to start oidc provider: %w", err)
	}
	apis.RegisterHandlerPrefixes(oidcProvider.HttpHandler(), "/.well-known/openid-configuration", "/oidc/v1", "/oauth/v2")

	samlProvider, err := saml.NewProvider(config.SAML, config.ExternalSecure, commands, queries, authRepo, keys.OIDC, keys.SAML, eventstore, dbClient, instanceInterceptor.Handler, userAgentInterceptor, limitingAccessInterceptor.Handle)
	if err != nil {
		return fmt.Errorf("unable to start saml provider: %w", err)
	}
	apis.RegisterHandlerOnPrefix(saml.HandlerPrefix, samlProvider.HttpHandler())

	c, err := console.Start(config.Console, config.ExternalSecure, oidcProvider.IssuerFromRequest, middleware.CallDurationHandler, instanceInterceptor.Handler, limitingAccessInterceptor, config.CustomerPortal)
	if err != nil {
		return fmt.Errorf("unable to start console: %w", err)
	}
	apis.RegisterHandlerOnPrefix(console.HandlerPrefix, c)

	l, err := login.CreateLogin(config.Login, commands, queries, authRepo, store, console.HandlerPrefix+"/", op.AuthCallbackURL(oidcProvider), provider.AuthCallbackURL(samlProvider), config.ExternalSecure, userAgentInterceptor, op.NewIssuerInterceptor(oidcProvider.IssuerFromRequest).Handler, provider.NewIssuerInterceptor(samlProvider.IssuerFromRequest).Handler, instanceInterceptor.Handler, assetsCache.Handler, limitingAccessInterceptor.Handle, keys.User, keys.IDPConfig, keys.CSRFCookieKey)
	if err != nil {
		return fmt.Errorf("unable to start login: %w", err)
	}
	apis.RegisterHandlerOnPrefix(login.HandlerPrefix, l.Handler())
	apis.HandleFunc(login.EndpointDeviceAuth, login.RedirectDeviceAuthToPrefix)

	// After OIDC provider so that the callback endpoint can be used
	if err := apis.RegisterService(ctx, oidc_v2.CreateServer(commands, queries, oidcProvider, config.ExternalSecure)); err != nil {
		return err
	}

	// handle grpc at last to be able to handle the root, because grpc and gateway require a lot of different prefixes
	apis.RouteGRPC()
	return nil
}

func listen(ctx context.Context, router *mux.Router, port uint16, tlsConfig *tls.Config, shutdown <-chan os.Signal) error {
	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(router, http2Server), TLSConfig: tlsConfig}

	lc := listenConfig()
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
	fmt.Println(color.MagentaString(figure.NewFigure("Zitadel", "", true).String()))
	http := "http"
	if startConfig.TLS.Enabled || startConfig.ExternalSecure {
		http = "https"
	}

	consoleURL := fmt.Sprintf("%s://%s:%v/ui/console\n", http, startConfig.ExternalDomain, startConfig.ExternalPort)
	healthCheckURL := fmt.Sprintf("%s://%s:%v/debug/healthz\n", http, startConfig.ExternalDomain, startConfig.ExternalPort)

	insecure := !startConfig.TLS.Enabled && !startConfig.ExternalSecure

	fmt.Printf(" ===============================================================\n\n")
	fmt.Printf(" Version          : %s\n", build.Version())
	fmt.Printf(" TLS enabled      : %v\n", startConfig.TLS.Enabled)
	fmt.Printf(" External Secure  : %v\n", startConfig.ExternalSecure)
	fmt.Printf(" Console URL      : %s", color.BlueString(consoleURL))
	fmt.Printf(" Health Check URL : %s", color.BlueString(healthCheckURL))
	if insecure {
		fmt.Printf("\n %s: you're using plain http without TLS. Be aware this is \n", color.RedString("Warning"))
		fmt.Printf(" not a secure setup and should only be used for test systems.         \n")
		fmt.Printf(" Visit: %s    \n", color.CyanString("https://zitadel.com/docs/self-hosting/manage/tls_modes"))
	}
	fmt.Printf("\n ===============================================================\n\n")
}
