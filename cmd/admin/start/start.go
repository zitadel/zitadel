package start

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caos/logging"
	"github.com/caos/oidc/v2/pkg/op"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/caos/zitadel/cmd/admin/key"
	admin_es "github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api"
	"github.com/caos/zitadel/internal/api/assets"
	internal_authz "github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/admin"
	"github.com/caos/zitadel/internal/api/grpc/auth"
	"github.com/caos/zitadel/internal/api/grpc/management"
	"github.com/caos/zitadel/internal/api/grpc/system"
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/api/oidc"
	"github.com/caos/zitadel/internal/api/ui/console"
	"github.com/caos/zitadel/internal/api/ui/login"
	auth_es "github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/authz"
	authz_repo "github.com/caos/zitadel/internal/authz/repository"
	"github.com/caos/zitadel/internal/command"
	cryptoDB "github.com/caos/zitadel/internal/crypto/database"
	"github.com/caos/zitadel/internal/database"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/notification"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/static"
	"github.com/caos/zitadel/internal/webauthn"
	"github.com/caos/zitadel/openapi"
)

func New() *cobra.Command {
	start := &cobra.Command{
		Use:   "start",
		Short: "starts ZITADEL instance",
		Long: `starts ZITADEL.
Requirements:
- cockroachdb`,
		RunE: func(cmd *cobra.Command, args []string) error {
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

	dbClient, err := database.Connect(config.Database)
	if err != nil {
		return fmt.Errorf("cannot start client for projection: %w", err)
	}

	keyStorage, err := cryptoDB.NewKeyStorage(dbClient, masterKey)
	if err != nil {
		return fmt.Errorf("cannot start key storage: %w", err)
	}
	keys, err := ensureEncryptionKeys(config.EncryptionKeys, keyStorage)
	if err != nil {
		return err
	}

	eventstoreClient, err := eventstore.Start(dbClient)
	if err != nil {
		return fmt.Errorf("cannot start eventstore for queries: %w", err)
	}

	queries, err := query.StartQueries(ctx, eventstoreClient, dbClient, config.Projections, keys.OIDC, config.InternalAuthZ.RolePermissionMappings)
	if err != nil {
		return fmt.Errorf("cannot start queries: %w", err)
	}

	authZRepo, err := authz.Start(config.AuthZ, config.SystemDefaults, queries, dbClient, keys.OIDC)
	if err != nil {
		return fmt.Errorf("error starting authz repo: %w", err)
	}

	storage, err := config.AssetStorage.NewStorage(dbClient)
	if err != nil {
		return fmt.Errorf("cannot start asset storage client: %w", err)
	}
	webAuthNConfig := &webauthn.Config{
		DisplayName:    config.WebAuthNName,
		ExternalSecure: config.ExternalSecure,
	}
	commands, err := command.StartCommands(eventstoreClient, config.SystemDefaults, config.InternalAuthZ.RolePermissionMappings, storage, authZRepo, webAuthNConfig, keys.IDPConfig, keys.OTP, keys.SMTP, keys.SMS, keys.User, keys.DomainVerification, keys.OIDC)
	if err != nil {
		return fmt.Errorf("cannot start commands: %w", err)
	}

	notification.Start(config.Notification, config.ExternalPort, config.ExternalSecure, commands, queries, dbClient, assets.HandlerPrefix, keys.User, keys.SMTP, keys.SMS)

	router := mux.NewRouter()
	err = startAPIs(ctx, router, commands, queries, eventstoreClient, dbClient, config, storage, authZRepo, keys)
	if err != nil {
		return err
	}
	return listen(ctx, router, config.Port)
}

func startAPIs(ctx context.Context, router *mux.Router, commands *command.Commands, queries *query.Queries, eventstore *eventstore.Eventstore, dbClient *sql.DB, config *Config, store static.Storage, authZRepo authz_repo.Repository, keys *encryptionKeys) error {
	repo := struct {
		authz_repo.Repository
		*query.Queries
	}{
		authZRepo,
		queries,
	}
	verifier := internal_authz.Start(repo)

	authenticatedAPIs := api.New(config.Port, router, &repo, config.InternalAuthZ, config.ExternalSecure, config.HTTP2HostHeader)
	authRepo, err := auth_es.Start(config.Auth, config.SystemDefaults, commands, queries, dbClient, assets.HandlerPrefix, keys.OIDC, keys.User)
	if err != nil {
		return fmt.Errorf("error starting auth repo: %w", err)
	}
	adminRepo, err := admin_es.Start(config.Admin, store, dbClient, login.HandlerPrefix)
	if err != nil {
		return fmt.Errorf("error starting admin repo: %w", err)
	}
	if err := authenticatedAPIs.RegisterServer(ctx, system.CreateServer(commands, queries, adminRepo, config.DefaultInstance, config.ExternalPort, config.ExternalDomain, config.ExternalSecure)); err != nil {
		return err
	}
	if err := authenticatedAPIs.RegisterServer(ctx, admin.CreateServer(commands, queries, adminRepo, assets.HandlerPrefix, keys.User)); err != nil {
		return err
	}
	if err := authenticatedAPIs.RegisterServer(ctx, management.CreateServer(commands, queries, config.SystemDefaults, assets.HandlerPrefix, keys.User, config.ExternalSecure, oidc.HandlerPrefix)); err != nil {
		return err
	}
	if err := authenticatedAPIs.RegisterServer(ctx, auth.CreateServer(commands, queries, authRepo, config.SystemDefaults, assets.HandlerPrefix, keys.User, config.ExternalSecure)); err != nil {
		return err
	}

	instanceInterceptor := middleware.InstanceInterceptor(queries, config.HTTP1HostHeader)
	authenticatedAPIs.RegisterHandler(assets.HandlerPrefix, assets.NewHandler(commands, verifier, config.InternalAuthZ, id.SonyFlakeGenerator, store, queries, instanceInterceptor.Handler))

	userAgentInterceptor, err := middleware.NewUserAgentHandler(config.UserAgentCookie, keys.UserAgentCookieKey, id.SonyFlakeGenerator, config.ExternalSecure)
	if err != nil {
		return err
	}

	oidcProvider, err := oidc.NewProvider(ctx, config.OIDC, login.DefaultLoggedOutPath, config.ExternalSecure, commands, queries, authRepo, keys.OIDC, keys.OIDCKey, eventstore, dbClient, userAgentInterceptor, instanceInterceptor.Handler)
	if err != nil {
		return fmt.Errorf("unable to start oidc provider: %w", err)
	}
	authenticatedAPIs.RegisterHandler(oidc.HandlerPrefix, oidcProvider.HttpHandler())

	openAPIHandler, err := openapi.Start()
	if err != nil {
		return fmt.Errorf("unable to start openapi handler: %w", err)
	}
	authenticatedAPIs.RegisterHandler(openapi.HandlerPrefix, openAPIHandler)

	baseURL := http_util.BuildHTTP(config.ExternalDomain, config.ExternalPort, config.ExternalSecure)
	c, err := console.Start(config.Console, config.ExternalSecure, oidcProvider.IssuerFromRequest, instanceInterceptor.Handler)
	if err != nil {
		return fmt.Errorf("unable to start console: %w", err)
	}
	authenticatedAPIs.RegisterHandler(console.HandlerPrefix, c)

	l, err := login.CreateLogin(config.Login, commands, queries, authRepo, store, config.SystemDefaults, console.HandlerPrefix+"/", config.ExternalDomain, baseURL, op.AuthCallbackURL(oidcProvider), config.ExternalSecure, userAgentInterceptor, op.NewIssuerInterceptor(oidcProvider.IssuerFromRequest).Handler, instanceInterceptor.Handler, keys.User, keys.IDPConfig, keys.CSRFCookieKey)
	if err != nil {
		return fmt.Errorf("unable to start login: %w", err)
	}
	authenticatedAPIs.RegisterHandler(login.HandlerPrefix, l.Handler())

	return nil
}

func listen(ctx context.Context, router *mux.Router, port uint16) error {
	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(router, http2Server)}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("tcp listener on %d failed: %w", port, err)
	}

	errCh := make(chan error)

	go func() {
		logging.Infof("server is listening on %s", lis.Addr().String())
		errCh <- http1Server.Serve(lis)
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
