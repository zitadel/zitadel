package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/caos/zitadel/internal/api/saml"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/caos/logging"
	"github.com/getsentry/sentry-go"
	"github.com/rs/cors"

	admin_es "github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api"
	"github.com/caos/zitadel/internal/api/assets"
	internal_authz "github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/admin"
	"github.com/caos/zitadel/internal/api/grpc/auth"
	"github.com/caos/zitadel/internal/api/grpc/management"
	"github.com/caos/zitadel/internal/api/oidc"
	auth_es "github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/authz"
	authz_repo "github.com/caos/zitadel/internal/authz/repository"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/notification"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/caos/zitadel/internal/setup"
	"github.com/caos/zitadel/internal/static"
	static_config "github.com/caos/zitadel/internal/static/config"
	metrics "github.com/caos/zitadel/internal/telemetry/metrics/config"
	tracing "github.com/caos/zitadel/internal/telemetry/tracing/config"
	"github.com/caos/zitadel/internal/ui"
	"github.com/caos/zitadel/internal/ui/console"
	"github.com/caos/zitadel/internal/ui/login"
	"github.com/caos/zitadel/openapi"
)

// build argument
var version = "dev"

type Config struct {
	Log            logging.Config
	Tracing        tracing.TracingConfig
	Metrics        metrics.MetricsConfig
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

	API api.Config
	UI  ui.Config

	Notification notification.Config
}

type setupConfig struct {
	Log logging.Config

	Eventstore     types.SQL
	SystemDefaults sd.SystemDefaults
	SetUp          setup.IAMSetUp
	InternalAuthZ  internal_authz.Config
}

var (
	configPaths         = config.NewArrayFlags("authz.yaml", "startup.yaml", "system-defaults.yaml")
	setupPaths          = config.NewArrayFlags("authz.yaml", "system-defaults.yaml", "setup.yaml")
	adminEnabled        = flag.Bool("admin", true, "enable admin api")
	managementEnabled   = flag.Bool("management", true, "enable management api")
	authEnabled         = flag.Bool("auth", true, "enable auth api")
	oidcEnabled         = flag.Bool("oidc", true, "enable oidc api")
	samlEnabled         = flag.Bool("saml", true, "enable saml api")
	assetsEnabled       = flag.Bool("assets", true, "enable assets api")
	loginEnabled        = flag.Bool("login", true, "enable login ui")
	consoleEnabled      = flag.Bool("console", true, "enable console ui")
	notificationEnabled = flag.Bool("notification", true, "enable notification handler")
	localDevMode        = flag.Bool("localDevMode", false, "enable local development specific configs")
)

const (
	cmdStart = "start"
	cmdSetup = "setup"
)

func main() {
	enableSentry, _ := strconv.ParseBool(os.Getenv("SENTRY_USAGE"))

	if enableSentry {
		sentryVersion := version
		if !regexp.MustCompile("^v?[0-9]+.[0-9]+.[0-9]$").Match([]byte(version)) {
			sentryVersion = "dev"
		}
		err := sentry.Init(sentry.ClientOptions{
			Environment: os.Getenv("SENTRY_ENVIRONMENT"),
			Release:     fmt.Sprintf("zitadel-%s", sentryVersion),
		})
		if err != nil {
			logging.Log("MAIN-Gnzjw").WithError(err).Fatal("sentry init failed")
		}
		sentry.CaptureMessage("sentry started")
		logging.Log("MAIN-adgf3").Info("sentry started")
		defer func() {
			err := recover()

			if err != nil {
				sentry.CurrentHub().Recover(err)
				sentry.Flush(2 * time.Second)
				panic(err)
			}
		}()
	}
	flag.Var(configPaths, "config-files", "paths to the config files")
	flag.Var(setupPaths, "setup-files", "paths to the setup files")
	flag.Parse()
	arg := flag.Arg(0)
	switch arg {
	case cmdStart:
		startZitadel(configPaths.Values())
	case cmdSetup:
		startSetup(setupPaths.Values())
	default:
		logging.Log("MAIN-afEQ2").Fatal("please provide an valid argument [start, setup]")
	}
}

func startZitadel(configPaths []string) {
	conf := new(Config)
	err := config.Read(conf, configPaths...)
	logging.Log("ZITAD-EDz31").OnError(err).Fatal("cannot read config")

	logging.LogWithFields("MAIN-dsfg2",
		"HTTP_PROXY", os.Getenv("HTTP_PROXY") != "",
		"HTTPS_PROXY", os.Getenv("HTTPS_PROXY") != "",
		"NO_PROXY", os.Getenv("NO_PROXY")).Info("http proxy settings")

	ctx := context.Background()
	esQueries, err := eventstore.StartWithUser(conf.EventstoreBase, conf.Queries.Eventstore)
	if err != nil {
		logging.Log("MAIN-Ddv21").OnError(err).Fatal("cannot start eventstore for queries")
	}

	keyChan := make(chan interface{})
	certChan := make(chan interface{})
	queries, err := query.StartQueries(ctx, esQueries, conf.Projections, conf.SystemDefaults, keyChan, certChan, conf.InternalAuthZ.RolePermissionMappings)
	logging.Log("MAIN-WpeJY").OnError(err).Fatal("cannot start queries")

	authZRepo, err := authz.Start(conf.AuthZ, conf.SystemDefaults, queries)
	logging.Log("MAIN-s9KOw").OnError(err).Fatal("error starting authz repo")

	esCommands, err := eventstore.StartWithUser(conf.EventstoreBase, conf.Commands.Eventstore)
	logging.Log("ZITAD-iRCMm").OnError(err).Fatal("cannot start eventstore for commands")

	store, err := conf.AssetStorage.Config.NewStorage()
	logging.Log("ZITAD-Bfhe2").OnError(err).Fatal("Unable to start asset storage")

	commands, err := command.StartCommands(esCommands, conf.SystemDefaults, conf.InternalAuthZ, store, authZRepo)
	if err != nil {
		logging.Log("ZITAD-bmNiJ").OnError(err).Fatal("cannot start commands")
	}

	var authRepo *auth_es.EsRepository
	if *authEnabled || *oidcEnabled || *loginEnabled {
		authRepo, err = auth_es.Start(conf.Auth, conf.SystemDefaults, commands, queries)
		logging.Log("MAIN-9oRw6").OnError(err).Fatal("error starting auth repo")
	}

	repo := struct {
		authz_repo.Repository
		query.Queries
	}{
		authZRepo,
		*queries,
	}

	verifier := internal_authz.Start(&repo)
	startAPI(ctx, conf, verifier, authZRepo, authRepo, commands, queries, store, esQueries, conf.Projections.CRDB, keyChan, certChan)
	startUI(ctx, conf, authRepo, commands, queries, store)

	if *notificationEnabled {
		notification.Start(ctx, conf.Notification, conf.SystemDefaults, commands, queries, store != nil)
	}

	<-ctx.Done()
	logging.Log("MAIN-s8d2h").Info("stopping zitadel")
}

func startUI(ctx context.Context, conf *Config, authRepo *auth_es.EsRepository, command *command.Commands, query *query.Queries, staticStorage static.Storage) {
	uis := ui.Create(conf.UI)
	if *loginEnabled {
		login, prefix := login.Start(conf.UI.Login, command, query, authRepo, staticStorage, conf.SystemDefaults, *localDevMode)
		uis.RegisterHandler(prefix, login.Handler())
	}
	if *consoleEnabled {
		consoleHandler, prefix, err := console.Start(conf.UI.Console)
		logging.Log("API-AGD1f").OnError(err).Fatal("error starting console")
		uis.RegisterHandler(prefix, consoleHandler)
	}
	uis.Start(ctx)
}

func startAPI(ctx context.Context, conf *Config, verifier *internal_authz.TokenVerifier, authZRepo authz_repo.Repository, authRepo *auth_es.EsRepository, command *command.Commands, query *query.Queries, static static.Storage, es *eventstore.Eventstore, projections types.SQL, keyChan <-chan interface{}, certChan <-chan interface{}) {
	repo, err := admin_es.Start(ctx, conf.Admin, conf.SystemDefaults, command, static, *localDevMode)
	logging.Log("API-D42tq").OnError(err).Fatal("error starting auth repo")

	apis := api.Create(conf.API, conf.InternalAuthZ, query, authZRepo, authRepo, repo, conf.SystemDefaults)

	if *adminEnabled {
		apis.RegisterServer(ctx, admin.CreateServer(command, query, repo, conf.SystemDefaults.Domain, conf.API.Domain+"/assets/v1/"))
	}
	if *managementEnabled {
		apis.RegisterServer(ctx, management.CreateServer(command, query, conf.SystemDefaults, conf.API.Domain+"/assets/v1/"))
	}
	if *authEnabled {
		apis.RegisterServer(ctx, auth.CreateServer(command, query, authRepo, conf.SystemDefaults, conf.API.Domain+"/assets/v1/"))
	}
	if *oidcEnabled {
		op := oidc.NewProvider(ctx, conf.API.OIDC, command, query, authRepo, conf.SystemDefaults.KeyConfig, *localDevMode, es, projections, keyChan, conf.API.Domain+"/assets/v1/")
		apis.RegisterHandler("/oauth/v2", op.HttpHandler())
	}

	if *samlEnabled {
		idp, err := saml.NewProvider(ctx, conf.API.SAML, command, query, authRepo, conf.SystemDefaults.KeyConfig, es, projections, certChan, *localDevMode)
		if err != nil {
			logging.Log("API-pwaiks").OnError(err).Fatal("error starting saml")
		}
		apis.RegisterHandler("/saml", idp.HttpHandler())
	}

	if *assetsEnabled {
		assetsHandler := assets.NewHandler(command, verifier, conf.InternalAuthZ, id.SonyFlakeGenerator, static, query)
		apis.RegisterHandler("/assets/v1", assetsHandler)
	}

	openAPIHandler, err := openapi.Start()
	logging.Log("ZITAD-8pRk1").OnError(err).Fatal("Unable to start openapi handler")
	apis.RegisterHandler("/openapi/v2/swagger", cors.AllowAll().Handler(openAPIHandler))

	apis.Start(ctx)
}

func startSetup(configPaths []string) {
	conf := new(setupConfig)
	err := config.Read(conf, configPaths...)
	logging.Log("MAIN-FaF2r").OnError(err).Fatal("cannot read config")

	ctx := context.Background()

	es, err := eventstore.Start(conf.Eventstore)
	logging.Log("MAIN-Ddt3").OnError(err).Fatal("cannot start eventstore")

	commands, err := command.StartCommands(es, conf.SystemDefaults, conf.InternalAuthZ, nil, nil)
	logging.Log("MAIN-dsjrr").OnError(err).Fatal("cannot start command side")

	err = setup.Execute(ctx, conf.SetUp, conf.SystemDefaults.IamID, commands)
	logging.Log("MAIN-djs3R").OnError(err).Panic("failed to execute setup steps")
}
