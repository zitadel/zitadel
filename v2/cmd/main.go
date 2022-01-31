package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/caos/logging"
	"github.com/caos/oidc/pkg/op"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	admin_es "github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api/assets"
	internal_authz "github.com/caos/zitadel/internal/api/authz"
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/http/middleware"
	auth_es "github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/authz"
	authz_repo "github.com/caos/zitadel/internal/authz/repository"
	"github.com/caos/zitadel/internal/cache/bigcache"
	cache_config "github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	types_v1 "github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	es_sql "github.com/caos/zitadel/internal/eventstore/repository/sql"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/notification"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/caos/zitadel/internal/static"
	static_config "github.com/caos/zitadel/internal/static/config"
	"github.com/caos/zitadel/openapi"
	"github.com/caos/zitadel/v2/internal/api"
	"github.com/caos/zitadel/v2/internal/api/grpc/admin"
	"github.com/caos/zitadel/v2/internal/api/grpc/auth"
	"github.com/caos/zitadel/v2/internal/api/grpc/management"
	middlewareV2 "github.com/caos/zitadel/v2/internal/api/http/middleware"
	"github.com/caos/zitadel/v2/internal/api/oidc"
	"github.com/caos/zitadel/v2/internal/api/ui/console"
	"github.com/caos/zitadel/v2/internal/api/ui/login"
	"github.com/caos/zitadel/v2/internal/config"
	"github.com/caos/zitadel/v2/internal/config/types"
)

const (
	cmdStart = "start"
	cmdSetup = "setup"
)

// build argument
var version = "dev-v2"

type startConfig struct {
	BaseDomain     string
	Commands       types.SQLUser
	Queries        types.SQLUser
	Projections    projectionConfig
	SystemDefaults systemdefaults.SystemDefaults
	AuthZ          authz.Config
	InternalAuthZ  internal_authz.Config
	EventstoreBase types.SQLBase
	Auth           auth_es.Config
	AssetStorage   static_config.AssetStorageConfig
	Admin          admin_es.Config
	OIDC           api.Config
	Login          login.Config
	Console        console.Config
	Notification   notification.Config
}

type projectionConfig struct {
	projection.Config
	CRDB types.SQL
}

const (
	eventstoreDB             = "eventstore"
	queryUser                = "queries"
	commandUser              = "eventstore"
	projectionSchema         = "projections"
	maxOpenConnection        = 3
	maxConnLifetime          = 30 * time.Minute
	maxConnIdleTime          = 30 * time.Minute
	commandMaxOpenConnection = 5
	commandMaxConnLifetime   = 30 * time.Minute
	commandMaxConnIdleTime   = 30 * time.Minute
	queryMaxOpenConnection   = 2
	queryMaxConnLifetime     = 30 * time.Minute
	queryMaxConnIdleTime     = 30 * time.Minute

	envBaseDomain          = "ZITADEL_BASE_DOMAIN"
	envEventstoreHost      = "ZITADEL_EVENTSTORE_HOST"
	envEventstorePort      = "ZITADEL_EVENTSTORE_PORT"
	envEventstoreSSLMode   = "CR_SSL_MODE"
	envEventstoreRootCert  = "CR_ROOT_CERT"
	envCockroachOptions    = "CR_OPTIONS"
	envCommandPassword     = "CR_EVENTSTORE_PASSWORD"
	envCommandCert         = "CR_EVENTSTORE_CERT"
	envCommandKey          = "CR_EVENTSTORE_KEY"
	envQueriesPassword     = "CR_QUERIES_PASSWORD"
	envQueriesCert         = "CR_QUERIES_CERT"
	envQueriesKey          = "CR_QUERIES_KEY"
	envConsoleOverwriteDir = "ZITADEL_CONSOLE_DIR"
	envCSRFKey             = "ZITADEL_CSRF_KEY"
	envCookieKey           = "ZITADEL_COOKIE_KEY"
	envOIDCKey             = "ZITADEL_OIDC_KEYS_ID"
	envSentryUsage         = "SENTRY_USAGE"
	envSentryEnvironment   = "SENTRY_ENVIRONMENT"
	pathLogin              = "/ui/login"
	pathOAuthV2            = "/oauth/v2"
	pathAssetAPI           = "/assets/v1"
	projectionDB           = "zitadel"

	defaultPort = "50002" //TODO: change to 80?
)

var (
	configPaths = config.NewArrayFlags("authz.yaml", "system-defaults.yaml", "startup.yaml")
	setupPaths  = config.NewArrayFlags("authz.yaml", "system-defaults.yaml", "setup.yaml")
	port        = flag.String("port", defaultPort, "port to run ZITADEL on")
	//TODO: do we still need these flags:
	adminEnabled        = flag.Bool("admin", true, "enable admin api")
	managementEnabled   = flag.Bool("management", true, "enable management api")
	authEnabled         = flag.Bool("auth", true, "enable auth api")
	oidcEnabled         = flag.Bool("oidc", true, "enable oidc api")
	assetsEnabled       = flag.Bool("assets", true, "enable assets api")
	loginEnabled        = flag.Bool("login", true, "enable login ui")
	consoleEnabled      = flag.Bool("console", true, "enable console ui")
	notificationEnabled = flag.Bool("notification", true, "enable notification handler")
	//TODO: especially this one
	localDevMode = flag.Bool("localDevMode", false, "enable local development specific configs")
)

func main() {
	sentryEnabled, _ := strconv.ParseBool(os.Getenv(envSentryUsage)) //TODO: default false!!! do we want that
	if sentryEnabled {
		enableSentry()
	}
	flag.Var(configPaths, "config-files", "paths to the config files")
	flag.Var(setupPaths, "setup-files", "paths to the setup files")
	flag.Parse()

	arg := flag.Arg(0)
	switch arg {
	case cmdStart:
		startZitadel()
	case cmdSetup:
		//startSetup()
	default:
		logging.Log("MAIN-afEQ2").Fatal("please provide an valid argument [start, setup]")
	}
}

func startZitadel() {
	conf := configureStart()

	ctx := context.Background()

	keyChan := make(chan interface{})
	projectionsDB, err := conf.Projections.CRDB.Start()
	logging.Log("MAIN-DAgw1").OnError(err).Fatal("cannot start client for projection")
	store, err := conf.AssetStorage.Config.NewStorage()
	logging.Log("MAIN-Bfhe2").OnError(err).Fatal("Unable to start asset storage")

	sqlClient, err := conf.Queries.Start(conf.EventstoreBase)
	logging.Log("MAIN-Ddv21").OnError(err).Fatal("cannot start eventstore for queries")
	esQueries := eventstore.NewEventstore(es_sql.NewCRDB(sqlClient))
	queries, err := query.StartQueries2(ctx, esQueries, projectionsDB, conf.Projections.Config, conf.SystemDefaults, keyChan, conf.InternalAuthZ.RolePermissionMappings)
	logging.Log("MAIN-WpeJY").OnError(err).Fatal("cannot start queries")

	authZRepo, err := authz.Start(conf.AuthZ, conf.SystemDefaults, queries)
	logging.Log("MAIN-s9KOw").OnError(err).Fatal("error starting authz repo")

	sqlClient, err = conf.Commands.Start(conf.EventstoreBase)
	logging.Log("MAIN-iRCMm").OnError(err).Fatal("cannot start eventstore for commands")
	esCommands := eventstore.NewEventstore(es_sql.NewCRDB(sqlClient))
	commands, err := command.StartCommands(esCommands, conf.SystemDefaults, conf.InternalAuthZ, store, authZRepo)
	logging.Log("MAIN-bmNiJ").OnError(err).Fatal("cannot start commands")

	notification.Start(ctx, conf.Notification, conf.SystemDefaults, commands, queries, store != nil)

	router := mux.NewRouter()
	startAPIs(ctx, router, commands, queries, esQueries, projectionsDB, keyChan, conf, store, authZRepo)
	listen(ctx, router)
}

func configureStart() *startConfig {
	conf := defaultStartConfig()
	err := config.Read(conf, configPaths.Values()...)
	logging.Log("MAIN-EDz31").OnError(err).Fatal("cannot read config")

	return conf
}

func startAPIs(ctx context.Context, router *mux.Router, commands *command.Commands, queries *query.Queries, eventstore *eventstore.Eventstore, projectionsDB *sql.DB, keyChan chan interface{}, conf *startConfig, store static.Storage, authZRepo authz_repo.Repository) {
	repo := struct {
		authz_repo.Repository
		*query.Queries
	}{
		authZRepo,
		queries,
	}
	verifier := internal_authz.Start(&repo)

	apis := api.New(ctx, *port, router, verifier, conf.InternalAuthZ, conf.SystemDefaults)

	adminRepo, err := admin_es.Start(ctx, conf.Admin, conf.SystemDefaults, commands, store, *localDevMode)
	logging.Log("MAIN-D42tq").OnError(err).Fatal("error starting auth repo")
	authRepo, err := auth_es.Start(conf.Auth, conf.SystemDefaults, commands, queries)
	logging.Log("MAIN-9oRw6").OnError(err).Fatal("error starting auth repo")

	apis.RegisterServer(ctx, admin.CreateServer(commands, queries, adminRepo, conf.SystemDefaults.Domain, pathAssetAPI))
	apis.RegisterServer(ctx, management.CreateServer(commands, queries, conf.SystemDefaults, pathAssetAPI))
	apis.RegisterServer(ctx, auth.CreateServer(commands, queries, authRepo, conf.SystemDefaults, pathAssetAPI))

	if store != nil {
		assetsHandler := assets.NewHandler(commands, verifier, conf.InternalAuthZ, id.SonyFlakeGenerator, store, queries)
		apis.RegisterHandler(pathAssetAPI, assetsHandler)
	}

	oidcProvider := oidc.NewProvider(ctx, conf.OIDC.OPHandlerConfig, commands, queries, authRepo, conf.SystemDefaults.KeyConfig, *localDevMode, eventstore, projectionsDB, keyChan, pathAssetAPI, conf.BaseDomain)
	apis.RegisterHandler("/oauth/v2", oidcProvider.HttpHandler())

	openAPIHandler, err := openapi.Start()
	logging.Log("MAIN-8pRk1").OnError(err).Fatal("Unable to start openapi handler")
	apis.RegisterHandler("/openapi/v2/swagger", cors.AllowAll().Handler(openAPIHandler))

	consoleID, err := consoleClientID(ctx, queries)
	logging.Log("MAIN-Dgfqs").OnError(err).Fatal("unable to get client_id for console")
	c, err := console.Start(conf.Console, conf.BaseDomain, *port, conf.OIDC.OPConfig.Issuer, consoleID)
	apis.RegisterHandler(console.HandlerPrefix, c)

	l := login.CreateLogin(conf.Login, commands, queries, authRepo, store, conf.SystemDefaults, *localDevMode, conf.BaseDomain, console.HandlerPrefix)
	apis.RegisterHandler(login.HandlerPrefix, l.Handler())

	apis.Router()
}

func listen(ctx context.Context, router *mux.Router) {
	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(router, http2Server)}
	lis := http_util.CreateListener(*port)

	go func() {
		logging.LogWithFields("MAIN-DG2FG", "port", lis.Addr().String()).Info("server is listening")
		fmt.Println("listening on " + lis.Addr().String())
		err := http1Server.Serve(lis)
		logging.Log("MAIN-Dvf21").OnError(err).Panic("grpc server serve failed")
	}()

	<-ctx.Done()

	if shutdownErr := http1Server.Shutdown(ctx); shutdownErr != nil && !errors.Is(shutdownErr, context.Canceled) {
		logging.Log("MAIN-SWE2A").WithError(shutdownErr).Panic("graceful server shutdown failed")
	}

	logging.Log("MAIN-dsGra").Info("server close")
}

func enableSentry() {
	sentryVersion := version
	if !regexp.MustCompile("^v?[0-9]+.[0-9]+.[0-9]$").Match([]byte(version)) {
		sentryVersion = version
	}
	err := sentry.Init(sentry.ClientOptions{
		Environment: os.Getenv(envSentryEnvironment),
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

func defaultStartConfig() *startConfig {
	return &startConfig{
		BaseDomain:     os.Getenv(envBaseDomain),
		EventstoreBase: eventstoreConfig(),
		Commands:       defaultCommandConfig(),
		Queries:        defaultQueryConfig(),
		Projections:    defaultProjectionConfig(),
		OIDC:           defaultOIDCConfig(),
		Login:          defaultLoginConfig(),
		Console:        defaultConsoleConfig(),
		SystemDefaults: systemdefaults.SystemDefaults{},    //remove later; until then, read from system-defaults.yaml
		AuthZ:          authz.Config{},                     //remove later; until then, read from startup.yaml
		InternalAuthZ:  internal_authz.Config{},            //remove later?
		Auth:           auth_es.Config{},                   //remove later; until then, read from startup.yaml
		AssetStorage:   static_config.AssetStorageConfig{}, //TODO: default config?
		Admin:          admin_es.Config{},                  //remove later; until then, read from startup.yaml
	}
}

func eventstoreConfig() types.SQLBase {
	return types.SQLBase{
		Host:            os.Getenv(envEventstoreHost),
		Port:            os.Getenv(envEventstorePort),
		Database:        eventstoreDB,
		SSL:             defaultEventstoreSSL(),
		Options:         os.Getenv(envCockroachOptions),
		MaxOpenConns:    maxOpenConnection,
		MaxConnLifetime: types.Duration{Duration: maxConnLifetime},
		MaxConnIdleTime: types.Duration{Duration: maxConnIdleTime},
	}
}

func defaultEventstoreSSL() types.SSLBase {
	return types.SSLBase{
		Mode:     os.Getenv(envEventstoreSSLMode),
		RootCert: os.Getenv(envEventstoreRootCert),
	}
}

func defaultCommandConfig() types.SQLUser {
	return types.SQLUser{
		User:            commandUser,
		Password:        os.Getenv(envCommandPassword),
		MaxOpenConns:    commandMaxOpenConnection,
		MaxConnLifetime: types.Duration{Duration: commandMaxConnLifetime},
		MaxConnIdleTime: types.Duration{Duration: commandMaxConnIdleTime},
		SSL: types.SSLUser{
			Cert: os.Getenv(envCommandCert),
			Key:  os.Getenv(envCommandKey),
		},
	}
}

func defaultQueryConfig() types.SQLUser {
	return types.SQLUser{
		User:            queryUser,
		Password:        os.Getenv(envQueriesPassword),
		MaxOpenConns:    queryMaxOpenConnection,
		MaxConnLifetime: types.Duration{Duration: queryMaxConnLifetime},
		MaxConnIdleTime: types.Duration{Duration: queryMaxConnIdleTime},
		SSL: types.SSLUser{
			Cert: os.Getenv(envQueriesCert),
			Key:  os.Getenv(envQueriesKey),
		},
	}
}

func defaultProjectionConfig() projectionConfig {
	return projectionConfig{
		Config: projection.Config{
			RequeueEvery:     types_v1.Duration{Duration: 10 * time.Second},
			RetryFailedAfter: types_v1.Duration{Duration: 1 * time.Second},
			MaxFailureCount:  5,
			BulkLimit:        200,
			MaxIterators:     1,
			Customizations: map[string]projection.CustomConfig{
				"projects": {
					BulkLimit: func(i uint64) *uint64 { return &i }(2000),
				},
			},
		},
		CRDB: types.SQL{
			Host:     os.Getenv(envEventstoreHost),
			Port:     os.Getenv(envEventstorePort),
			User:     queryUser,
			Password: os.Getenv(envQueriesPassword),
			Database: projectionDB,
			Schema:   projectionSchema,
			SSL: &types.SSL{
				Mode:     os.Getenv(envEventstoreSSLMode),
				RootCert: os.Getenv(envEventstoreRootCert),
				Cert:     os.Getenv(envQueriesCert),
				Key:      os.Getenv(envQueriesKey),
			},
			Options:         os.Getenv(envCockroachOptions),
			MaxOpenConns:    maxOpenConnection,
			MaxConnLifetime: types.Duration{Duration: maxConnLifetime},
			MaxConnIdleTime: types.Duration{Duration: maxConnIdleTime},
		},
	}
}

func defaultOIDCConfig() api.Config {
	return api.Config{
		OPHandlerConfig: oidc.OPHandlerConfig{
			OPConfig: &op.Config{
				Issuer:                   os.Getenv(envBaseDomain) + pathOAuthV2, //TODO: BaseDomain/oauth/v2/ ??
				CryptoKey:                [32]byte{},                             //TODO: change config type?
				DefaultLogoutRedirectURI: pathLogin + "/logout/done",             //TODO: still config?
				CodeMethodS256:           true,
				AuthMethodPost:           true,
				AuthMethodPrivateKeyJWT:  true,
				GrantTypeRefreshToken:    true,
				RequestObjectSupported:   true,
				SupportedUILocales:       nil, //TODO: change config type?
			},
			StorageConfig: oidc.StorageConfig{
				DefaultLoginURL:                   fmt.Sprintf("%s%s?%s=", pathLogin, login.EndpointLogin, login.QueryAuthRequestID), //TODO: still config?
				SigningKeyAlgorithm:               "RS256",
				DefaultAccessTokenLifetime:        types.Duration{Duration: 12 * time.Hour},
				DefaultIdTokenLifetime:            types.Duration{Duration: 12 * time.Hour},
				DefaultRefreshTokenIdleExpiration: types.Duration{Duration: 720 * time.Hour},  // 30 days
				DefaultRefreshTokenExpiration:     types.Duration{Duration: 2160 * time.Hour}, // 90 days
			},
			UserAgentCookieConfig: defaultUserAgentCookieConfig(),
			Cache: &middleware.CacheConfig{
				MaxAge:       types_v1.Duration{Duration: 12 * time.Hour},
				SharedMaxAge: types_v1.Duration{Duration: 168 * time.Hour}, // 7 days
			},
			KeyConfig: &crypto.KeyConfig{
				EncryptionKeyID: os.Getenv(envOIDCKey),
			},
			CustomEndpoints: nil, // use default endpoints from OIDC library
		},
	}
}

func defaultLoginConfig() login.Config {
	return login.Config{
		OidcAuthCallbackURL: pathOAuthV2 + "/authorize/callback?id=", //TODO: provide from op
		LanguageCookieName:  "zitadel.login.lang",                    //TODO: constant? (console might need it as well)
		CSRF: login.CSRF{
			CookieName: "zitadel.login.csrf", //TODO: constant?
			Key: &crypto.KeyConfig{
				EncryptionKeyID: os.Getenv(envCSRFKey),
			},
			Development: *localDevMode, //TODO: ?
		},
		UserAgentCookieConfig: defaultUserAgentCookieConfig(),
		Cache:                 defaultCacheConfig(),
		StaticCache: cache_config.CacheConfig{
			Type: "bigcache",
			Config: &bigcache.Config{
				MaxCacheSizeInMB: 52428800, //50MB
			},
		},
	}
}

func defaultConsoleConfig() console.Config {
	return console.Config{
		ConsoleOverwriteDir: os.Getenv(envConsoleOverwriteDir),
		ShortCache:          defaultShortCacheConfig(),
		LongCache:           defaultCacheConfig(),
	}
}

func defaultShortCacheConfig() middleware.CacheConfig {
	return middleware.CacheConfig{
		MaxAge:       types_v1.Duration{Duration: 5 * time.Minute},
		SharedMaxAge: types_v1.Duration{Duration: 15 * time.Minute},
	}
}

func defaultCacheConfig() middleware.CacheConfig {
	return middleware.CacheConfig{
		MaxAge:       types_v1.Duration{Duration: 12 * time.Hour},
		SharedMaxAge: types_v1.Duration{Duration: 168 * time.Hour}, // 7 days
	}
}

func defaultUserAgentCookieConfig() *middlewareV2.UserAgentCookieConfig {
	return &middlewareV2.UserAgentCookieConfig{
		Name: "zitadel.useragent", //TODO: constant?
		Key: &crypto.KeyConfig{
			EncryptionKeyID: os.Getenv(envCookieKey),
		},
		MaxAge: types.Duration{Duration: 8760 * time.Hour}, // 365 days
	}
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
