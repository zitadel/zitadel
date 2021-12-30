package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/caos/logging"
	admin_es "github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	api_v1 "github.com/caos/zitadel/internal/api"
	internal_authz "github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/oidc"
	auth_es "github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/authz"
	"github.com/caos/zitadel/internal/ui"

	// authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore"
	mgmt_es "github.com/caos/zitadel/internal/management/repository/eventsourcing"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/query/projection"
	static_config "github.com/caos/zitadel/internal/static/config"
	"github.com/caos/zitadel/v2/api"
	"github.com/caos/zitadel/v2/api/ui/console"
	"github.com/caos/zitadel/v2/api/ui/login"
)

func main() {
	flag.Var(configPaths, "config-files", "paths to the config files")
	flag.Var(setupPaths, "setup-files", "paths to the setup files")
	flag.Parse()
	conf := configure()

	ctx := context.TODO()

	esQueries, err := eventstore.StartWithUser(conf.EventstoreBase, conf.Queries.Eventstore)
	if err != nil {
		logging.Log("MAIN-Ddv21").OnError(err).Fatal("cannot start eventstore for queries")
	}

	queries, err := query.StartQueries(ctx, esQueries, conf.Projections, conf.SystemDefaults)
	logging.Log("MAIN-WpeJY").OnError(err).Fatal("cannot start queries")

	authZRepo, err := authz.Start(ctx, conf.AuthZ, conf.InternalAuthZ, conf.SystemDefaults, queries)
	logging.Log("MAIN-s9KOw").OnError(err).Fatal("error starting authz repo")

	esCommands, err := eventstore.StartWithUser(conf.EventstoreBase, conf.Commands.Eventstore)
	logging.Log("ZITAD-iRCMm").OnError(err).Fatal("cannot start eventstore for commands")

	store, err := conf.AssetStorage.Config.NewStorage()
	logging.Log("ZITAD-Bfhe2").OnError(err).Fatal("Unable to start asset storage")

	commands, err := command.StartCommands(esCommands, conf.SystemDefaults, conf.InternalAuthZ, store, authZRepo)
	if err != nil {
		logging.Log("ZITAD-bmNiJ").OnError(err).Fatal("cannot start commands")
	}

	authRepo, err := auth_es.Start(conf.Auth, conf.InternalAuthZ, conf.SystemDefaults, commands, queries, authZRepo, esQueries)
	logging.Log("MAIN-9oRw6").OnError(err).Fatal("error starting auth repo")

	// repo := struct {
	// 	authz_repo.EsRepository
	// 	query.Queries
	// }{
	// 	*authZRepo,
	// 	*queries,
	// }

	// verifier := internal_authz.Start(&repo)

	baseRouter := mux.NewRouter().StrictSlash(true)

	l, loginPrefix := login.CreateLogin(baseRouter, login.Config(conf.UI.Login), commands, queries, authRepo, store, conf.SystemDefaults, true)
	baseRouter.PathPrefix(loginPrefix).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/login", l.Handler()).ServeHTTP(w, r)
	})

	oidcHandler := oidc.NewProvider(ctx, conf.API.OIDC, commands, queries, authRepo, conf.SystemDefaults.KeyConfig.EncryptionConfig, true)
	baseRouter.PathPrefix("/oauth/v2").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/oauth/v2", oidcHandler.HttpHandler()).ServeHTTP(w, r)
	})

	listen(ctx, baseRouter, conf)
}

func configure() *Config {
	conf := new(Config)
	err := config.Read(conf, configPaths.Values()...)
	logging.Log("ZITAD-EDz31").OnError(err).Fatal("cannot read config")
	return conf
}

func listen(ctx context.Context, baseRouter *mux.Router, config *Config) {
	api.New(ctx, baseRouter)
	uiRouter := baseRouter.PathPrefix("/ui").Subrouter()
	consoleDir := "./console/"
	if config.UI.Console.ConsoleOverwriteDir != "" {
		consoleDir = config.UI.Console.ConsoleOverwriteDir
	}
	console.New(uiRouter, console.Config{
		ConsoleOverwriteDir: consoleDir,
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
