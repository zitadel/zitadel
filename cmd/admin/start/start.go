package start

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caos/logging"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	admin_es "github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api"
	"github.com/caos/zitadel/internal/api/assets"
	internal_authz "github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/admin"
	"github.com/caos/zitadel/internal/api/grpc/auth"
	"github.com/caos/zitadel/internal/api/grpc/management"
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/api/oidc"
	"github.com/caos/zitadel/internal/api/ui/console"
	"github.com/caos/zitadel/internal/api/ui/login"
	auth_es "github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/authz"
	authz_repo "github.com/caos/zitadel/internal/authz/repository"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/database"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/notification"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/caos/zitadel/internal/static"
	static_config "github.com/caos/zitadel/internal/static/config"
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
			config := new(startConfig)
			err := viper.Unmarshal(config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToSliceHookFunc(":"),
			)))
			if err != nil {
				return err
			}
			err = config.Log.SetLogger()
			if err != nil {
				return err
			}
			return startZitadel(config)
		},
	}
	bindUint16Flag(start, "port", "port to run ZITADEL on")
	bindStringFlag(start, "externalDomain", "domain ZITADEL will be exposed on")
	bindStringFlag(start, "externalPort", "port ZITADEL will be exposed on")
	bindBoolFlag(start, "externalSecure", "if ZITADEL will be served on HTTPS")

	return start
}

func bindStringFlag(cmd *cobra.Command, name, description string) {
	cmd.PersistentFlags().String(name, viper.GetString(name), description)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func bindUint16Flag(cmd *cobra.Command, name, description string) {
	cmd.PersistentFlags().Uint16(name, uint16(viper.GetUint(name)), description)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func bindBoolFlag(cmd *cobra.Command, name, description string) {
	cmd.PersistentFlags().Bool(name, viper.GetBool(name), description)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

type startConfig struct {
	Log             *logging.Config
	Port            uint16
	ExternalPort    uint16
	ExternalDomain  string
	ExternalSecure  bool
	Database        database.Config
	Projections     projectionConfig
	AuthZ           authz.Config
	Auth            auth_es.Config
	Admin           admin_es.Config
	UserAgentCookie *middleware.UserAgentCookieConfig
	OIDC            oidc.Config
	Login           login.Config
	Console         console.Config
	Notification    notification.Config
	AssetStorage    static_config.AssetStorageConfig
	InternalAuthZ   internal_authz.Config
	SystemDefaults  systemdefaults.SystemDefaults
}

type projectionConfig struct {
	projection.Config
	KeyConfig *crypto.KeyConfig
}

func startZitadel(config *startConfig) error {
	ctx := context.Background()
	keyChan := make(chan interface{})

	dbClient, err := database.Connect(config.Database)
	if err != nil {
		return fmt.Errorf("cannot start client for projection: %w", err)
	}
	storage, err := config.AssetStorage.NewStorage(dbClient)
	if err != nil {
		return fmt.Errorf("cannot start asset storage client: %w", err)
	}

	eventstoreClient, err := eventstore.Start(dbClient)
	if err != nil {
		return fmt.Errorf("cannot start eventstore for queries: %w", err)
	}
	smtpPasswordCrypto, err := crypto.NewAESCrypto(config.SystemDefaults.SMTPPasswordVerificationKey)
	if err != nil {
		return fmt.Errorf("cannot create smtp crypto: %w", err)
	}

	smsCrypto, err := crypto.NewAESCrypto(config.SystemDefaults.SMSVerificationKey)
	if err != nil {
		return fmt.Errorf("cannot create smtp crypto: %w", err)
	}

	queries, err := query.StartQueries(ctx, eventstoreClient, dbClient, config.Projections.Config, config.SystemDefaults, config.Projections.KeyConfig, keyChan, config.InternalAuthZ.RolePermissionMappings)
	if err != nil {
		return fmt.Errorf("cannot start queries: %w", err)
	}

	authZRepo, err := authz.Start(config.AuthZ, config.SystemDefaults, queries, dbClient, config.OIDC.KeyConfig)
	if err != nil {
		return fmt.Errorf("error starting authz repo: %w", err)
	}
	webAuthNConfig := webauthn.Config{
		ID:          config.ExternalDomain,
		Origin:      http_util.BuildHTTP(config.ExternalDomain, config.ExternalPort, config.ExternalSecure),
		DisplayName: "ZITADEL",
	}
	commands, err := command.StartCommands(eventstoreClient, config.SystemDefaults, config.InternalAuthZ, storage, authZRepo, config.OIDC.KeyConfig, webAuthNConfig, smtpPasswordCrypto, smsCrypto)
	if err != nil {
		return fmt.Errorf("cannot start commands: %w", err)
	}

	notification.Start(config.Notification, config.SystemDefaults, commands, queries, dbClient, assets.HandlerPrefix, smtpPasswordCrypto)

	router := mux.NewRouter()
	err = startAPIs(ctx, router, commands, queries, eventstoreClient, dbClient, keyChan, config, storage, authZRepo)
	if err != nil {
		return err
	}
	return listen(ctx, router, config.Port)
}

func startAPIs(ctx context.Context, router *mux.Router, commands *command.Commands, queries *query.Queries, eventstore *eventstore.Eventstore, dbClient *sql.DB, keyChan chan interface{}, config *startConfig, store static.Storage, authZRepo authz_repo.Repository) error {
	repo := struct {
		authz_repo.Repository
		*query.Queries
	}{
		authZRepo,
		queries,
	}
	verifier := internal_authz.Start(repo)

	apis := api.New(config.Port, router, &repo, config.InternalAuthZ, config.ExternalSecure)
	userEncryptionAlgorithm, err := crypto.NewAESCrypto(config.SystemDefaults.UserVerificationKey)
	if err != nil {
		return nil
	}
	authRepo, err := auth_es.Start(config.Auth, config.SystemDefaults, commands, queries, dbClient, config.OIDC.KeyConfig, assets.HandlerPrefix, userEncryptionAlgorithm)
	if err != nil {
		return fmt.Errorf("error starting auth repo: %w", err)
	}
	adminRepo, err := admin_es.Start(config.Admin, store, dbClient, login.HandlerPrefix)
	if err != nil {
		return fmt.Errorf("error starting admin repo: %w", err)
	}
	if err := apis.RegisterServer(ctx, admin.CreateServer(commands, queries, adminRepo, config.SystemDefaults.Domain, assets.HandlerPrefix, userEncryptionAlgorithm)); err != nil {
		return err
	}
	if err := apis.RegisterServer(ctx, management.CreateServer(commands, queries, config.SystemDefaults, assets.HandlerPrefix, userEncryptionAlgorithm)); err != nil {
		return err
	}
	if err := apis.RegisterServer(ctx, auth.CreateServer(commands, queries, authRepo, config.SystemDefaults, assets.HandlerPrefix, userEncryptionAlgorithm)); err != nil {
		return err
	}

	apis.RegisterHandler(assets.HandlerPrefix, assets.NewHandler(commands, verifier, config.InternalAuthZ, id.SonyFlakeGenerator, store, queries))

	userAgentInterceptor, err := middleware.NewUserAgentHandler(config.UserAgentCookie, config.ExternalDomain, id.SonyFlakeGenerator, config.ExternalSecure)
	if err != nil {
		return err
	}

	issuer := oidc.Issuer(config.ExternalDomain, config.ExternalPort, config.ExternalSecure)
	oidcProvider, err := oidc.NewProvider(ctx, config.OIDC, issuer, login.DefaultLoggedOutPath, commands, queries, authRepo, config.SystemDefaults.KeyConfig, eventstore, dbClient, keyChan, userAgentInterceptor)
	if err != nil {
		return fmt.Errorf("unable to start oidc provider: %w", err)
	}
	apis.RegisterHandler(oidc.HandlerPrefix, oidcProvider.HttpHandler())

	openAPIHandler, err := openapi.Start()
	if err != nil {
		return fmt.Errorf("unable to start openapi handler: %w", err)
	}
	apis.RegisterHandler(openapi.HandlerPrefix, openAPIHandler)

	consoleID, err := consoleClientID(ctx, queries)
	if err != nil {
		return fmt.Errorf("unable to get client_id for console: %w", err)
	}
	c, err := console.Start(config.Console, config.ExternalDomain, http_util.BuildHTTP(config.ExternalDomain, config.ExternalPort, config.ExternalSecure), issuer, consoleID)
	if err != nil {
		return fmt.Errorf("unable to start console: %w", err)
	}
	apis.RegisterHandler(console.HandlerPrefix, c)

	l, err := login.CreateLogin(config.Login, commands, queries, authRepo, store, config.SystemDefaults, console.HandlerPrefix, config.ExternalDomain, oidc.AuthCallback, config.ExternalSecure, userAgentInterceptor, userEncryptionAlgorithm)
	if err != nil {
		return fmt.Errorf("unable to start login: %w", err)
	}
	apis.RegisterHandler(login.HandlerPrefix, l.Handler())

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

//TODO:!!??!!
func consoleClientID(ctx context.Context, queries *query.Queries) (string, error) {
	iam, err := queries.IAMByID(ctx, domain.IAMID)
	if err != nil {
		return "", err
	}
	projectID, err := query.NewAppProjectIDSearchQuery(iam.IAMProjectID)
	if err != nil {
		return "", err
	}
	name, err := query.NewAppNameSearchQuery(query.TextContainsIgnoreCase, "console") //TODO:!!??!!
	if err != nil {
		return "", err
	}
	apps, err := queries.SearchApps(ctx, &query.AppSearchQueries{
		Queries: []query.SearchQuery{projectID, name},
	})
	if err != nil {
		return "", err
	}
	if len(apps.Apps) != 1 || apps.Apps[0].OIDCConfig == nil {
		return "", errors.New("invalid app")
	}
	return apps.Apps[0].OIDCConfig.ClientID, nil
}
