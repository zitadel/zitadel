package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/caos/logging"
	admin_es "github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	api_v1 "github.com/caos/zitadel/internal/api"
	"github.com/caos/zitadel/internal/api/assets"
	internal_authz "github.com/caos/zitadel/internal/api/authz"
	admin_grpc "github.com/caos/zitadel/internal/api/grpc/admin"
	"github.com/caos/zitadel/internal/api/grpc/auth"
	"github.com/caos/zitadel/internal/api/grpc/management"
	"github.com/caos/zitadel/internal/api/oidc"
	auth_es "github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/authz"
	authz_es "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	mgmt_es "github.com/caos/zitadel/internal/management/repository/eventsourcing"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/caos/zitadel/internal/static"
	static_config "github.com/caos/zitadel/internal/static/config"
	"github.com/caos/zitadel/internal/ui"
	"github.com/caos/zitadel/openapi"
	"github.com/caos/zitadel/v2/api"
	"github.com/caos/zitadel/v2/api/ui/console"
	"github.com/caos/zitadel/v2/api/ui/login"
)

var (
	mgmtSvc  *management.Server
	adminSvc *admin_grpc.Server
	authSvc  *auth.Server

	mgmtRepo  *mgmt_es.EsRepository
	adminRepo *admin_es.EsRepository
	authRepo  *auth_es.EsRepository
	authZRepo *authz_es.EsRepository

	queries    *query.Queries
	commands   *command.Commands
	assetStore static.Storage

	esQueries *eventstore.Eventstore

	verifier *internal_authz.TokenVerifier
)

func main() {
	flag.Var(configPaths, "config-files", "paths to the config files")
	flag.Var(setupPaths, "setup-files", "paths to the setup files")
	flag.Parse()
	conf := configure()

	ctx := context.TODO()

	startCQRS(ctx, conf)
	startRepos(ctx, conf)
	startSvcs(ctx, conf)

	listen(ctx, conf)
}

func startSvcs(ctx context.Context, conf *Config) {
	mgmtSvc = management.CreateServer(commands, queries, mgmtRepo, conf.SystemDefaults, conf.Mgmt.APIDomain+"/assets/v1/")
	authSvc = auth.CreateServer(commands, queries, authRepo, conf.SystemDefaults)
	adminSvc = admin_grpc.CreateServer(commands, queries, adminRepo, conf.SystemDefaults.Domain, "/assets/v1")

	repo := struct {
		authz_es.EsRepository
		*query.Queries
	}{
		*authZRepo,
		queries,
	}

	verifier = internal_authz.Start(&repo)
}

func startCQRS(ctx context.Context, conf *Config) {
	var err error

	esQueries, err = eventstore.StartWithUser(conf.EventstoreBase, conf.Queries.Eventstore)
	logging.Log("MAIN-Ddv21").OnError(err).Fatal("cannot start eventstore for queries")

	queries, err = query.StartQueries(ctx, esQueries, conf.Projections, conf.SystemDefaults)
	logging.Log("MAIN-WpeJY").OnError(err).Fatal("cannot start queries")

	esCommands, err := eventstore.StartWithUser(conf.EventstoreBase, conf.Commands.Eventstore)
	logging.Log("ZITAD-iRCMm").OnError(err).Fatal("cannot start eventstore for commands")

	assetStore, err = conf.AssetStorage.Config.NewStorage()
	logging.Log("ZITAD-Bfhe2").OnError(err).Fatal("Unable to start asset storage")

	commands, err = command.StartCommands(esCommands, conf.SystemDefaults, conf.InternalAuthZ, assetStore, authZRepo)
	logging.Log("ZITAD-bmNiJ").OnError(err).Fatal("cannot start commands")
}

func startRepos(ctx context.Context, conf *Config) {
	roles := make([]string, len(conf.InternalAuthZ.RolePermissionMappings))
	for i, role := range conf.InternalAuthZ.RolePermissionMappings {
		roles[i] = role.Role
	}

	var err error

	authZRepo, err = authz.Start(ctx, conf.AuthZ, conf.InternalAuthZ, conf.SystemDefaults, queries)
	logging.Log("MAIN-s9KOw").OnError(err).Fatal("error starting authz repo")

	mgmtRepo, err = mgmt_es.Start(conf.Mgmt, conf.SystemDefaults, roles, queries, assetStore)
	logging.Log("API-Gd2qq").OnError(err).Fatal("error starting management repo")

	authRepo, err = auth_es.Start(conf.Auth, conf.InternalAuthZ, conf.SystemDefaults, commands, queries, authZRepo, esQueries)
	logging.Log("MAIN-9oRw6").OnError(err).Fatal("error starting auth repo")
}

func configure() *Config {
	conf := new(Config)
	err := config.Read(conf, configPaths.Values()...)
	logging.Log("ZITAD-EDz31").OnError(err).Fatal("cannot read config")
	return conf
}

func listen(ctx context.Context, config *Config) {
	baseRouter := mux.NewRouter().StrictSlash(true)

	api.New(ctx, baseRouter, mgmtSvc, adminSvc, authSvc, verifier, config.InternalAuthZ)

	l, loginPrefix := login.CreateLogin(baseRouter, login.Config(config.UI.Login), commands, queries, authRepo, assetStore, config.SystemDefaults, true)
	baseRouter.PathPrefix(loginPrefix).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(loginPrefix, l.Handler()).ServeHTTP(w, r)
	})

	oidcHandler := oidc.NewProvider(ctx, config.API.OIDC, commands, queries, authRepo, config.SystemDefaults.KeyConfig.EncryptionConfig, true)
	baseRouter.PathPrefix("/.well-known/openid-configuration").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oidcHandler.HttpHandler().ServeHTTP(w, r)
	})
	baseRouter.PathPrefix("/oauth/v2").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/oauth/v2", oidcHandler.HttpHandler()).ServeHTTP(w, r)
	})

	assetsHandler := assets.NewHandler(commands, verifier, config.InternalAuthZ, id.SonyFlakeGenerator, assetStore, mgmtRepo, queries)
	baseRouter.PathPrefix("/assets/v1").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets/v1", assetsHandler).ServeHTTP(w, r)
	})

	openAPIHandler, err := openapi.Start()
	logging.Log("ZITAD-qXpND").OnError(err).Fatal("unable to start openapi")
	baseRouter.PathPrefix("/openapi/v2/swagger").
		Handler(http.StripPrefix("/openapi/v2/swagger", cors.AllowAll().Handler(openAPIHandler)))

	uiRouter := baseRouter.PathPrefix("/ui").Subrouter()
	consoleDir := "./console/"
	if config.UI.Console.ConsoleOverwriteDir != "" {
		consoleDir = config.UI.Console.ConsoleOverwriteDir
		// consoleDir = "/Users/adlerhurst/Downloads/zitadel-console"
	}
	console.New(uiRouter, console.Config{
		ConsoleOverwriteDir: consoleDir,
		Environment: console.Environment{
			AuthServiceUrl:         config.Auth.APIDomain,
			MgmtServiceUrl:         config.Mgmt.APIDomain,
			AdminServiceUrl:        config.Admin.APIDomain,
			SubscriptionServiceUrl: config.Mgmt.APIDomain,
			AssetServiceUrl:        config.Mgmt.APIDomain,
			Issuer:                 config.API.OIDC.OPConfig.Issuer,
			Clientid:               "141602932889026980@zitadel",
		},
	})

	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(baseRouter, http2Server)}
	lis, err := net.Listen("tcp", ":50002")
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Println("listening on " + lis.Addr().String())
		err = http1Server.Serve(lis)
		fmt.Println("stopped serving")
	}()

	<-ctx.Done()

	if shutdownErr := http1Server.Shutdown(ctx); shutdownErr != nil && !errors.Is(shutdownErr, context.Canceled) {
		panic(shutdownErr)
	}

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		panic(err)
	}
}

var (
	configPaths         = config.NewArrayFlags("authz.yaml", "startup.yaml", "system-defaults.yaml")
	setupPaths          = config.NewArrayFlags("authz.yaml", "system-defaults.yaml", "setup.yaml")
	adminEnabled        = flag.Bool("admin", true, "enable admin api")
	managementEnabled   = flag.Bool("management", true, "enable management api")
	authEnabled         = flag.Bool("auth", true, "enable auth api")
	oidcEnabled         = flag.Bool("oidc", true, "enable oidc api")
	assetsEnabled       = flag.Bool("assets", true, "enable assets api")
	loginEnabled        = flag.Bool("login", true, "enable login ui")
	consoleEnabled      = flag.Bool("console", true, "enable console ui")
	notificationEnabled = flag.Bool("notification", true, "enable notification handler")
	localDevMode        = flag.Bool("localDevMode", false, "enable local development specific configs")
)

type Config struct {
	AssetStorage   static_config.AssetStorageConfig
	InternalAuthZ  internal_authz.Config
	SystemDefaults sd.SystemDefaults

	EventstoreBase types.SQLBase
	Commands       command.Config
	Queries        query.Config
	Projections    projection.Config

	AuthZ authz.Config
	Auth  auth_es.Config
	Admin admin_es.Config
	Mgmt  mgmt_es.Config

	UI  ui.Config
	API api_v1.Config
}
