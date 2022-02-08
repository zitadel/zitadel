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
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	"github.com/caos/zitadel/internal/cache/bigcache"
	cache_config "github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
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

const (
	cmdStart = "start"
	cmdSetup = "setup"
)

// build argument
var version = "dev-v2"

type startConfig struct {
	Domain         string
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
	OIDC           oidc.Config
	Login          login.Config
	Console        console.Config
	Notification   notification.Config
	WebAuthN       webauthn.Config
}

type projectionConfig struct {
	projection.Config
	CRDB      types.SQL
	KeyConfig *crypto.KeyConfig
}

const (
	eventstoreDB             = "eventstore"
	projectionDB             = "zitadel"
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

	envDomain              = "ZITADEL_DOMAIN"
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

	defaultPort = "8080"
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
	//case cmdSetup:
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
	var storage static.Storage
	if *assetsEnabled {
		storage, err = conf.AssetStorage.Config.NewStorage()
		logging.Log("MAIN-Bfhe2").OnError(err).Fatal("Unable to start asset storage")
	}
	esQueries, err := eventstore.StartWithUser(conf.EventstoreBase, conf.Queries)
	logging.Log("MAIN-Ddv21").OnError(err).Fatal("cannot start eventstore for queries")
	queries, err := query.StartQueries(ctx, esQueries, projectionsDB, conf.Projections.Config, conf.SystemDefaults, conf.Projections.KeyConfig, keyChan, conf.InternalAuthZ.RolePermissionMappings)
	logging.Log("MAIN-WpeJY").OnError(err).Fatal("cannot start queries")

	authZRepo, err := authz.Start(conf.AuthZ, conf.SystemDefaults, queries, conf.OIDC.KeyConfig)
	logging.Log("MAIN-s9KOw").OnError(err).Fatal("error starting authz repo")

	esCommands, err := eventstore.StartWithUser(conf.EventstoreBase, conf.Commands)
	logging.Log("MAIN-iRCMm").OnError(err).Fatal("cannot start eventstore for commands")
	commands, err := command.StartCommands(esCommands, conf.SystemDefaults, conf.InternalAuthZ, storage, authZRepo, conf.OIDC.KeyConfig, conf.WebAuthN)
	logging.Log("MAIN-bmNiJ").OnError(err).Fatal("cannot start commands")

	if *notificationEnabled {
		notification.Start(ctx, conf.Notification, conf.SystemDefaults, commands, queries, storage != nil)
	}

	router := mux.NewRouter()
	startAPIs(ctx, router, commands, queries, esQueries, projectionsDB, keyChan, conf, storage, authZRepo)
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
	verifier := internal_authz.Start(repo)

	apis := api.New(*port, router, &repo, conf.InternalAuthZ, conf.SystemDefaults)

	authRepo, err := auth_es.Start(conf.Auth, conf.SystemDefaults, commands, queries, conf.OIDC.KeyConfig)
	logging.Log("MAIN-9oRw6").OnError(err).Fatal("error starting auth repo")
	if *adminEnabled {
		adminRepo, err := admin_es.Start(ctx, conf.Admin, conf.SystemDefaults, commands, store, *localDevMode)
		logging.Log("MAIN-D42tq").OnError(err).Fatal("error starting auth repo")
		apis.RegisterServer(ctx, admin.CreateServer(commands, queries, adminRepo, conf.SystemDefaults.Domain, assets.HandlerPrefix))
	}
	if *managementEnabled {
		apis.RegisterServer(ctx, management.CreateServer(commands, queries, conf.SystemDefaults, assets.HandlerPrefix))
	}
	if *authEnabled {
		apis.RegisterServer(ctx, auth.CreateServer(commands, queries, authRepo, conf.SystemDefaults, assets.HandlerPrefix))
	}

	if *assetsEnabled {
		assetsHandler := assets.NewHandler(commands, verifier, conf.InternalAuthZ, id.SonyFlakeGenerator, store, queries)
		apis.RegisterHandler(assets.HandlerPrefix, assetsHandler)
	}

	if *oidcEnabled {
		oidcProvider := oidc.NewProvider(ctx, conf.OIDC, commands, queries, authRepo, conf.SystemDefaults.KeyConfig, *localDevMode, eventstore, projectionsDB, keyChan, assets.HandlerPrefix, conf.Domain)
		apis.RegisterHandler(oidc.HandlerPrefix, oidcProvider.HttpHandler())
	}

	openAPIHandler, err := openapi.Start()
	logging.Log("MAIN-8pRk1").OnError(err).Fatal("Unable to start openapi handler")
	apis.RegisterHandler(openapi.HandlerPrefix, cors.AllowAll().Handler(openAPIHandler))

	if *consoleEnabled {
		consoleID, err := consoleClientID(ctx, queries)
		logging.Log("MAIN-Dgfqs").OnError(err).Fatal("unable to get client_id for console")
		c, err := console.Start(conf.Console, local(conf.Domain), url(local(conf.Domain)), conf.OIDC.Issuer, consoleID)
		apis.RegisterHandler(console.HandlerPrefix, c)
	}

	if *loginEnabled {
		l := login.CreateLogin(conf.Login, commands, queries, authRepo, store, conf.SystemDefaults, *localDevMode, conf.Domain, console.HandlerPrefix)
		apis.RegisterHandler(login.HandlerPrefix, l.Handler())
	}
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

	logging.Log("MAIN-dsGra").Info("server closed")
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
		Domain:         os.Getenv(envDomain),
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
		WebAuthN:       defaultWebAuthNConfig(),
	}
}

func defaultWebAuthNConfig() webauthn.Config {
	return webauthn.Config{
		ID:          os.Getenv(envDomain),
		Origin:      url(local(os.Getenv(envDomain))),
		DisplayName: "ZITADEL",
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
		KeyConfig: &crypto.KeyConfig{
			EncryptionKeyID: os.Getenv(envOIDCKey),
		},
	}
}

func defaultOIDCConfig() oidc.Config {
	return oidc.Config{
		Issuer:                            url(local(os.Getenv(envDomain)) + oidc.HandlerPrefix), //TODO: Domain/oauth/v2/ ??
		DefaultLogoutRedirectURI:          login.HandlerPrefix + "/logout/done",                  //TODO: still config?
		CodeMethodS256:                    true,
		AuthMethodPost:                    true,
		AuthMethodPrivateKeyJWT:           true,
		GrantTypeRefreshToken:             true,
		RequestObjectSupported:            true,
		DefaultLoginURL:                   fmt.Sprintf("%s%s?%s=", login.HandlerPrefix, login.EndpointLogin, login.QueryAuthRequestID), //TODO: still config?
		SigningKeyAlgorithm:               "RS256",
		DefaultAccessTokenLifetime:        types.Duration{Duration: 12 * time.Hour},
		DefaultIdTokenLifetime:            types.Duration{Duration: 12 * time.Hour},
		DefaultRefreshTokenIdleExpiration: types.Duration{Duration: 720 * time.Hour},  // 30 days
		DefaultRefreshTokenExpiration:     types.Duration{Duration: 2160 * time.Hour}, // 90 days
		UserAgentCookieConfig:             defaultUserAgentCookieConfig(),
		Cache: &middleware.CacheConfig{
			MaxAge:       types.Duration{Duration: 12 * time.Hour},
			SharedMaxAge: types.Duration{Duration: 168 * time.Hour}, // 7 days
		},
		KeyConfig: &crypto.KeyConfig{
			EncryptionKeyID: os.Getenv(envOIDCKey),
		},
		CustomEndpoints: nil, // use default endpoints from OIDC library
	}
}

func local(host string) string {
	if *localDevMode {
		return host + ":" + *port
	}
	return host
}

func url(url string) string {
	if *localDevMode {
		return "http://" + url
	}
	return "https://" + url
}

func defaultLoginConfig() login.Config {
	return login.Config{
		OidcAuthCallbackURL: oidc.HandlerPrefix + "/authorize/callback?id=", //TODO: provide from op
		LanguageCookieName:  "zitadel.login.lang",                           //TODO: constant? (console might need it as well)
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
		MaxAge:       types.Duration{Duration: 5 * time.Minute},
		SharedMaxAge: types.Duration{Duration: 15 * time.Minute},
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
