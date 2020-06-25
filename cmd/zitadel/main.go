package main

import (
	"context"
	"flag"

	"github.com/caos/zitadel/internal/api"
	management2 "github.com/caos/zitadel/internal/api/grpc/management"
	auth_es "github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/authz"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/login"
	mgmt_es "github.com/caos/zitadel/internal/management/repository/eventsourcing"

	"github.com/caos/logging"

	internal_authz "github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config"
	"github.com/caos/zitadel/internal/notification"
	tracing "github.com/caos/zitadel/internal/tracing/config"
	"github.com/caos/zitadel/pkg/admin"
	"github.com/caos/zitadel/pkg/auth"
	"github.com/caos/zitadel/pkg/console"
)

type Config struct {
	API api.Config

	Mgmt         management2.Config
	Auth         auth.Config
	Login        login.Config
	AuthZ        authz.Config
	Admin        admin.Config
	Console      console.Config
	Notification notification.Config

	Log            logging.Config
	Tracing        tracing.TracingConfig
	InternalAuthZ  internal_authz.Config
	SystemDefaults sd.SystemDefaults
}

func main() {
	configPaths := config.NewArrayFlags("authz.yaml", "startup.yaml", "system-defaults.yaml")
	flag.Var(configPaths, "config-files", "paths to the config files")
	managementEnabled := flag.Bool("management", true, "enable management api")
	authEnabled := flag.Bool("auth", true, "enable auth api")
	oidcEnabled := flag.Bool("oidc", true, "enable oidc api")
	//loginEnabled := flag.Bool("login", true, "enable login ui")
	adminEnabled := flag.Bool("admin", true, "enable admin api")
	//consoleEnabled := flag.Bool("console", true, "enable console ui")
	notificationEnabled := flag.Bool("notification", true, "enable notification handler")
	flag.Parse()

	conf := new(Config)
	err := config.Read(conf, configPaths.Values()...)
	logging.Log("MAIN-FaF2r").OnError(err).Fatal("cannot read config")

	ctx := context.Background()
	authZRepo, err := authz.Start(ctx, conf.AuthZ, conf.InternalAuthZ, conf.SystemDefaults)
	logging.Log("MAIN-s9KOw").OnError(err).Fatal("error starting authz repo")

	var authRepo *auth_es.EsRepository
	if *authEnabled || *oidcEnabled {
		authRepo, err = auth_es.Start(conf.Auth, conf.InternalAuthZ, conf.SystemDefaults, authZRepo)
		logging.Log("API-9oRw6").OnError(err).Fatal("error starting auth repo")
	}
	api.Start(ctx, conf.API, conf.InternalAuthZ, authZRepo, conf.SystemDefaults, authRepo, *adminEnabled, *managementEnabled, *authEnabled, *oidcEnabled)

	//var authRepo *eventsourcing.EsRepository
	//if *authEnabled || *loginEnabled {
	//	authRepo, err = eventsourcing.Start(conf.Auth.Repository, conf.InternalAuthZ, conf.SystemDefaults, authZRepo)
	//	logging.Log("MAIN-9oRw6").OnError(err).Fatal("error starting auth repo")
	//}
	//if *authEnabled {
	//	auth.Start(ctx, conf.Auth, authZRepo, conf.InternalAuthZ, conf.SystemDefaults, authRepo)
	//}
	//if *loginEnabled {
	//	login.Start(ctx, conf.Login, conf.SystemDefaults, authRepo)
	//}
	if *notificationEnabled {
		notification.Start(ctx, conf.Notification, conf.SystemDefaults)
	}
	//if *consoleEnabled {
	//	err = console.Start(ctx, conf.Console)
	//	logging.Log("MAIN-3Dfuc").OnError(err).Fatal("error starting console ui")
	//}
	<-ctx.Done()
	logging.Log("MAIN-s8d2h").Info("stopping zitadel")
}
