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

	"github.com/caos/zitadel/internal/domain"

	"github.com/caos/zitadel/v2/internal/api"
	"github.com/caos/zitadel/v2/internal/api/grpc/admin"
	"github.com/caos/zitadel/v2/internal/api/grpc/auth"
	"github.com/caos/zitadel/v2/internal/api/grpc/management"
	"github.com/caos/zitadel/v2/internal/api/oidc"
	"github.com/caos/zitadel/v2/internal/api/ui/console"
	"github.com/caos/zitadel/v2/internal/api/ui/login"

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
	"github.com/caos/zitadel/internal/config"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/caos/zitadel/internal/static"
	static_config "github.com/caos/zitadel/internal/static/config"
	"github.com/caos/zitadel/openapi"
)

const (
	cmdStart = "start"
	cmdSetup = "setup"
)

// build argument
var version = "dev-v2"

type startConfig struct {
	Commands       command.ConfigV2
	Queries        query.ConfigV2
	Projections    projection.ConfigV2
	SystemDefaults systemdefaults.SystemDefaults
	AuthZ          authz.Config
	InternalAuthZ  internal_authz.Config
	EventstoreBase types.SQLBase2
	Auth           auth_es.Config
	AssetStorage   static_config.AssetStorageConfig
	Admin          admin_es.Config
	OIDC           api.Config
	Login          login.Config
	Console        console.Config
}

type querySQLConnection struct {
	types.SQL
}

const (
	eventstoreDB             = "eventstore"
	queryUser                = "queries"
	commandUser              = "eventstore"
	maxOpenConnection        = 3
	maxConnLifetime          = 30 * time.Minute
	maxConnIdleTime          = 30 * time.Minute
	commandMaxOpenConnection = 5
	commandMaxConnLifetime   = 30 * time.Minute
	commandMaxConnIdleTime   = 30 * time.Minute
	queryMaxOpenConnection   = 2
	queryMaxConnLifetime     = 30 * time.Minute
	queryMaxConnIdleTime     = 30 * time.Minute

	keyEventstoreHost   = "ZITADEL_EVENTSTORE_HOST"
	keyEventstorePort   = "ZITADEL_EVENTSTORE_PORT"
	keyCockroachOptions = "CR_OPTIONS"
	keyCommandPassword  = "CR_EVENTSTORE_PASSWORD"
	keyCommandCert      = "CR_EVENTSTORE_CERT"
	keyCommandKey       = "CR_EVENTSTORE_KEY"
	keyQueriesPassword  = "CR_QUERIES_PASSWORD"
	keyQueriesCert      = "CR_QUERIES_CERT"
	keyQueriesKey       = "CR_QUERIES_KEY"

	pathLogin   = "/ui/login"
	pathOAuthV2 = "/oauth/v2"
	pathConsole = "ui/console"

	defaultPort = "50002" //TODO: change to 80?
)

func defaultStartConfig() *startConfig {
	return &startConfig{
		Commands:       defaultCommandConfig(),
		Queries:        defaultQueryConfig(),
		Projections:    defaultProjectionConfig(),
		SystemDefaults: systemdefaults.SystemDefaults{}, //remove later; until then, read from startup.yaml
		AuthZ:          authz.Config{},                  //remove later; until then, read from startup.yaml
		InternalAuthZ:  internal_authz.Config{},
		EventstoreBase: eventstoreConfig(),
		Auth:           auth_es.Config{},
		AssetStorage:   static_config.AssetStorageConfig{},
		Admin:          admin_es.Config{},
		OIDC:           defaultOIDCConfig(),
		Login:          defaultLoginConfig(),
		Console:        defaultConsoleConfig(),
	}
}

func defaultConsoleConfig() console.Config {
	return console.Config{
		ConsoleOverwriteDir: os.Getenv("ZITADEL_CONSOLE_DIR"),
		ShortCache:          defaultShortCacheConfig(),
		LongCache:           defaultCacheConfig(),
		CSPDomain:           os.Getenv("ZITADEL_DEFAULT_DOMAIN"),
	}
}

func defaultShortCacheConfig() middleware.CacheConfig {
	return middleware.CacheConfig{
		MaxAge:       types.Duration{Duration: 5 * time.Minute},
		SharedMaxAge: types.Duration{Duration: 15 * time.Minute},
	}
}

func defaultLoginConfig() login.Config {
	return login.Config{
		BaseURL:             os.Getenv("ZITADEL_ACCOUNTS"),
		OidcAuthCallbackURL: pathOAuthV2 + "/authorize/callback?id=", //TODO: provide from op
		ZitadelURL:          pathConsole,
		LanguageCookieName:  "zitadel.login.lang", //TODO: constant? (console might need it as well)
		//DefaultLanguage:       language.Tag{},
		CSRF: login.CSRF{
			CookieName: "zitadel.login.csrf", //TODO: constant?
			Key: &crypto.KeyConfig{
				EncryptionKeyID: os.Getenv("ZITADEL_CSRF_KEY"),
			},
			Development: func() bool { //TODO: ???
				ok, _ := strconv.ParseBool(os.Getenv("ZITADEL_CSRF_DEV"))
				return ok
			}(),
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

func defaultCacheConfig() middleware.CacheConfig {
	return middleware.CacheConfig{
		MaxAge:       types.Duration{Duration: 12 * time.Hour},
		SharedMaxAge: types.Duration{Duration: 168 * time.Hour}, // 7 days
	}
}

func defaultUserAgentCookieConfig() *middleware.UserAgentCookieConfig {
	return &middleware.UserAgentCookieConfig{
		Name:   "zitadel.useragent", //TODO: constant?
		Domain: os.Getenv("ZITADEL_COOKIE_DOMAIN"),
		Key: &crypto.KeyConfig{
			EncryptionKeyID: os.Getenv("ZITADEL_COOKIE_KEY"),
		},
		MaxAge: types.Duration{Duration: 8760 * time.Hour}, // 365 days
	}
}

func defaultOIDCConfig() api.Config {
	return api.Config{
		OPHandlerConfig: oidc.OPHandlerConfig{
			OPConfig: &op.Config{
				Issuer:                   os.Getenv("ZITADEL_ISSUER"),
				CryptoKey:                [32]byte{},
				DefaultLogoutRedirectURI: pathLogin + "/logout/done",
				CodeMethodS256:           true,
				AuthMethodPost:           true,
				AuthMethodPrivateKeyJWT:  true,
				GrantTypeRefreshToken:    true,
				RequestObjectSupported:   true,
				SupportedUILocales:       nil, //TODO: change config type?
			},
			StorageConfig: oidc.StorageConfig{
				DefaultLoginURL:                   pathLogin + "/login?authRequestID=",
				SigningKeyAlgorithm:               "RS256",
				DefaultAccessTokenLifetime:        types.Duration{Duration: 12 * time.Hour},
				DefaultIdTokenLifetime:            types.Duration{Duration: 12 * time.Hour},
				DefaultRefreshTokenIdleExpiration: types.Duration{Duration: 720 * time.Hour},  // 30 days
				DefaultRefreshTokenExpiration:     types.Duration{Duration: 2160 * time.Hour}, // 90 days
			},
			UserAgentCookieConfig: defaultUserAgentCookieConfig(),
			Cache: &middleware.CacheConfig{
				MaxAge:       types.Duration{Duration: 12 * time.Hour},
				SharedMaxAge: types.Duration{Duration: 168 * time.Hour}, // 7 days
			},
			Endpoints: nil, // use default endpoints from OIDC library
		},
	}
}

const projectionDB = "zitadel"

const projectionSchema = "projections"

func defaultProjectionConfig() projection.ConfigV2 {
	return projection.ConfigV2{
		Config: projection.Config{
			RequeueEvery:     types.Duration{Duration: 10 * time.Second},
			RetryFailedAfter: types.Duration{Duration: 1 * time.Second},
			MaxFailureCount:  5,
			BulkLimit:        200,
			MaxIterators:     1,
			Customizations: map[string]projection.CustomConfig{
				"projects": {
					BulkLimit: func(i uint64) *uint64 { return &i }(2000),
				},
			},
		},
		CRDB: types.SQL2{
			Host:     os.Getenv(keyEventstoreHost),
			Port:     os.Getenv(keyEventstorePort),
			User:     queryUser,
			Password: os.Getenv(keyQueriesPassword),
			Database: projectionDB,
			Schema:   projectionSchema,
			SSL: &types.SSL{
				Mode:     os.Getenv("CR_SSL_MODE"),
				RootCert: os.Getenv("CR_ROOT_CERT"),
				Cert:     os.Getenv(keyQueriesCert),
				Key:      os.Getenv(keyQueriesKey),
			},
			Options:         os.Getenv(keyCockroachOptions),
			MaxOpenConns:    maxOpenConnection,
			MaxConnLifetime: types.Duration{Duration: maxConnLifetime},
			MaxConnIdleTime: types.Duration{Duration: maxConnIdleTime},
		},
	}
}

func eventstoreConfig() types.SQLBase2 {
	return types.SQLBase2{
		Host:            os.Getenv(keyEventstoreHost),
		Port:            os.Getenv(keyEventstorePort),
		Database:        eventstoreDB,
		Schema:          "",
		SSL:             defaultEventstoreSSL(),
		Options:         os.Getenv(keyCockroachOptions),
		MaxOpenConns:    maxOpenConnection,
		MaxConnLifetime: types.Duration{Duration: maxConnLifetime},
		MaxConnIdleTime: types.Duration{Duration: maxConnIdleTime},
	}
}

func defaultEventstoreSSL() types.SSLBase {
	return types.SSLBase{
		Mode:     os.Getenv("CR_SSL_MODE"),
		RootCert: os.Getenv("CR_ROOT_CERT"),
	}
}

func defaultCommandConfig() command.ConfigV2 {
	return command.ConfigV2{
		Eventstore: types.SQLUser2{
			User:            commandUser,
			Password:        os.Getenv(keyCommandPassword),
			MaxOpenConns:    commandMaxOpenConnection,
			MaxConnLifetime: types.Duration{Duration: commandMaxConnLifetime},
			MaxConnIdleTime: types.Duration{Duration: commandMaxConnIdleTime},
			SSL: types.SSLUser{
				Cert: os.Getenv(keyCommandCert),
				Key:  os.Getenv(keyCommandKey),
			},
		},
	}
}

func defaultQueryConfig() query.ConfigV2 {
	return query.ConfigV2{
		Eventstore: types.SQLUser2{
			User:            queryUser,
			Password:        os.Getenv(keyQueriesPassword),
			MaxOpenConns:    queryMaxOpenConnection,
			MaxConnLifetime: types.Duration{Duration: queryMaxConnLifetime},
			MaxConnIdleTime: types.Duration{Duration: queryMaxConnIdleTime},
			SSL: types.SSLUser{
				Cert: os.Getenv(keyQueriesCert),
				Key:  os.Getenv(keyQueriesKey),
			},
		},
	}
}

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
	sentryEnabled, _ := strconv.ParseBool(os.Getenv("SENTRY_USAGE")) //TODO: default false!!! do we want that
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

	esQueries, err := eventstore.StartWithUser2(conf.EventstoreBase, conf.Queries.Eventstore)
	logging.Log("MAIN-Ddv21").OnError(err).Fatal("cannot start eventstore for queries")

	projectionsDB, err := conf.Projections.CRDB.Start()
	logging.Log("MAIN-DAgw1").OnError(err).Fatal("cannot start client for projection")

	queries, err := query.StartQueries2(ctx, esQueries, projectionsDB, conf.Projections.Config, conf.SystemDefaults, keyChan, conf.InternalAuthZ.RolePermissionMappings)
	logging.Log("MAIN-WpeJY").OnError(err).Fatal("cannot start queries")

	authZRepo, err := authz.Start(conf.AuthZ, conf.SystemDefaults, queries)
	logging.Log("MAIN-s9KOw").OnError(err).Fatal("error starting authz repo")
	//
	esCommands, err := eventstore.StartWithUser2(conf.EventstoreBase, conf.Commands.Eventstore)
	logging.Log("ZITAD-iRCMm").OnError(err).Fatal("cannot start eventstore for commands")

	store, err := conf.AssetStorage.Config.NewStorage()
	logging.Log("ZITAD-Bfhe2").OnError(err).Fatal("Unable to start asset storage")

	commands, err := command.StartCommands(esCommands, conf.SystemDefaults, conf.InternalAuthZ, store, authZRepo)
	logging.Log("ZITAD-bmNiJ").OnError(err).Fatal("cannot start commands")

	router := mux.NewRouter()
	startAPIs(ctx, router, commands, queries, esQueries, projectionsDB, keyChan, conf, store, authZRepo)
	listen(ctx, router)
}

func configureStart() *startConfig {
	conf := defaultStartConfig()
	err := config.Read(conf, configPaths.Values()...)
	logging.Log("ZITAD-EDz31").OnError(err).Fatal("cannot read config")

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

	apis.RegisterServer(ctx, admin.CreateServer(commands, queries, adminRepo, conf.SystemDefaults.Domain, "/assets/v1/"))
	apis.RegisterServer(ctx, management.CreateServer(commands, queries, conf.SystemDefaults, "/assets/v1/"))
	apis.RegisterServer(ctx, auth.CreateServer(commands, queries, authRepo, conf.SystemDefaults, "/assets/v1/"))

	if store != nil {
		assetsHandler := assets.NewHandler(commands, verifier, conf.InternalAuthZ, id.SonyFlakeGenerator, store, queries)
		apis.RegisterHandler("/assets/v1", assetsHandler)
	}

	op := oidc.NewProvider(ctx, conf.OIDC.OPHandlerConfig, commands, queries, authRepo, conf.SystemDefaults.KeyConfig, *localDevMode, eventstore, projectionsDB, keyChan, "/assets/v1/")
	apis.RegisterHandler("/oauth/v2", op.HttpHandler())

	openAPIHandler, err := openapi.Start()
	logging.Log("MAIN-8pRk1").OnError(err).Fatal("Unable to start openapi handler")
	apis.RegisterHandler("/openapi/v2/swagger", cors.AllowAll().Handler(openAPIHandler))

	l, _ := login.CreateLogin(conf.Login, commands, queries, authRepo, store, conf.SystemDefaults, *localDevMode)
	apis.RegisterHandler(pathLogin, l.Handler())

	id, err := consoleClientID(ctx, queries)
	logging.Log("MAIN-Dgfqs").OnError(err).Fatal("unable to get client_id for console")
	c, consoleRoute, err := console.Start(conf.Console, conf.Auth.APIDomain, conf.OIDC.OPConfig.Issuer, id)
	apis.RegisterHandler(consoleRoute, c)

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

	//if errors.Is(err, http.ErrServerClosed) {
	//	fmt.Println("server closed")
	//} else if err != nil {
	//	panic(err)
	//}
}

func enableSentry() {
	sentryVersion := version
	if !regexp.MustCompile("^v?[0-9]+.[0-9]+.[0-9]$").Match([]byte(version)) {
		sentryVersion = version
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
