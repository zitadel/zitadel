package start

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zitadel/saml/pkg/provider"

	clockpkg "github.com/benbjohnson/clock"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/op"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

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
	"github.com/zitadel/zitadel/internal/api/grpc/system"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/api/saml"
	"github.com/zitadel/zitadel/internal/api/ui/console"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	auth_es "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/authz"
	authz_repo "github.com/zitadel/zitadel/internal/authz/repository"
	"github.com/zitadel/zitadel/internal/command"
	cryptoDB "github.com/zitadel/zitadel/internal/crypto/database"
	"github.com/zitadel/zitadel/internal/database"
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

func New() *cobra.Command {
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

			return startZitadel(config, masterKey)
		},
	}

	startFlags(start)

	return start
}

func startZitadel(config *Config, masterKey string) error {
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

	queries, err := query.StartQueries(ctx, eventstoreClient, dbClient, config.Projections, config.SystemDefaults, keys.IDPConfig, keys.OTP, keys.OIDC, keys.SAML, config.InternalAuthZ.RolePermissionMappings)
	if err != nil {
		return fmt.Errorf("cannot start queries: %w", err)
	}

	authZRepo, err := authz.Start(queries, dbClient, keys.OIDC, config.ExternalSecure)
	if err != nil {
		return fmt.Errorf("error starting authz repo: %w", err)
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

	usageReporter := logstore.UsageReporterFunc(commands.ReportUsage)
	actionsLogstoreSvc := logstore.New(queries, usageReporter, actionsExecutionDBEmitter, actionsExecutionStdoutEmitter)
	if actionsLogstoreSvc.Enabled() {
		logging.Warn("execution logs are currently in beta")
	}
	actions.SetLogstoreService(actionsLogstoreSvc)

	notification.Start(ctx, config.Projections.Customizations["notifications"], config.ExternalPort, config.ExternalSecure, commands, queries, eventstoreClient, assets.AssetAPIFromDomain(config.ExternalSecure, config.ExternalPort), config.SystemDefaults.Notifications.FileSystemPath, keys.User, keys.SMTP, keys.SMS)

	router := mux.NewRouter()
	tlsConfig, err := config.TLS.Config()
	if err != nil {
		return err
	}
	err = startAPIs(ctx, clock, router, commands, queries, eventstoreClient, dbClient, config, storage, authZRepo, keys, queries, usageReporter)
	if err != nil {
		return err
	}
	return listen(ctx, router, config.Port, tlsConfig)
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
	accessInterceptor := middleware.NewAccessInterceptor(accessSvc, config.Quotas.Access)
	apis := api.New(config.Port, router, queries, verifier, config.InternalAuthZ, tlsConfig, config.HTTP2HostHeader, config.HTTP1HostHeader, accessSvc)
	authRepo, err := auth_es.Start(ctx, config.Auth, config.SystemDefaults, commands, queries, dbClient, eventstore, keys.OIDC, keys.User)
	if err != nil {
		return fmt.Errorf("error starting auth repo: %w", err)
	}
	adminRepo, err := admin_es.Start(ctx, config.Admin, store, dbClient, eventstore)
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
	instanceInterceptor := middleware.InstanceInterceptor(queries, config.HTTP1HostHeader, login.IgnoreInstanceEndpoints...)
	assetsCache := middleware.AssetsCacheInterceptor(config.AssetStorage.Cache.MaxAge, config.AssetStorage.Cache.SharedMaxAge)
	apis.RegisterHandler(assets.HandlerPrefix, assets.NewHandler(commands, verifier, config.InternalAuthZ, id.SonyFlakeGenerator(), store, queries, middleware.CallDurationHandler, instanceInterceptor.Handler, assetsCache.Handler, accessInterceptor.Handle))

	userAgentInterceptor, err := middleware.NewUserAgentHandler(config.UserAgentCookie, keys.UserAgentCookieKey, id.SonyFlakeGenerator(), config.ExternalSecure, login.EndpointResources)
	if err != nil {
		return err
	}

	// TODO: Record openapi access logs?
	openAPIHandler, err := openapi.Start()
	if err != nil {
		return fmt.Errorf("unable to start openapi handler: %w", err)
	}
	apis.RegisterHandler(openapi.HandlerPrefix, openAPIHandler)

	oidcProvider, err := oidc.NewProvider(ctx, config.OIDC, login.DefaultLoggedOutPath, config.ExternalSecure, commands, queries, authRepo, keys.OIDC, keys.OIDCKey, eventstore, dbClient, userAgentInterceptor, instanceInterceptor.Handler, accessInterceptor.Handle)
	if err != nil {
		return fmt.Errorf("unable to start oidc provider: %w", err)
	}

	samlProvider, err := saml.NewProvider(ctx, config.SAML, config.ExternalSecure, commands, queries, authRepo, keys.OIDC, keys.SAML, eventstore, dbClient, instanceInterceptor.Handler, userAgentInterceptor, accessInterceptor.Handle)
	if err != nil {
		return fmt.Errorf("unable to start saml provider: %w", err)
	}
	apis.RegisterHandler(saml.HandlerPrefix, samlProvider.HttpHandler())

	c, err := console.Start(config.Console, config.ExternalSecure, oidcProvider.IssuerFromRequest, middleware.CallDurationHandler, instanceInterceptor.Handler, accessInterceptor.Handle, config.CustomerPortal)
	if err != nil {
		return fmt.Errorf("unable to start console: %w", err)
	}
	apis.RegisterHandler(console.HandlerPrefix, c)

	l, err := login.CreateLogin(config.Login, commands, queries, authRepo, store, console.HandlerPrefix+"/", op.AuthCallbackURL(oidcProvider), provider.AuthCallbackURL(samlProvider), config.ExternalSecure, userAgentInterceptor, op.NewIssuerInterceptor(oidcProvider.IssuerFromRequest).Handler, provider.NewIssuerInterceptor(samlProvider.IssuerFromRequest).Handler, instanceInterceptor.Handler, assetsCache.Handler, accessInterceptor.Handle, keys.User, keys.IDPConfig, keys.CSRFCookieKey)
	if err != nil {
		return fmt.Errorf("unable to start login: %w", err)
	}
	apis.RegisterHandler(login.HandlerPrefix, l.Handler())

	//handle oidc at last, to be able to handle the root
	//we might want to change that in the future
	//esp. if we want to have multiple well-known endpoints
	//it might make sense to handle the discovery endpoint and oauth and oidc prefixes individually
	//but this will require a change in the oidc lib
	apis.RegisterHandler("", oidcProvider.HttpHandler())
	return nil
}

func listen(ctx context.Context, router *mux.Router, port uint16, tlsConfig *tls.Config) error {
	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(router, http2Server), TLSConfig: tlsConfig}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

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
