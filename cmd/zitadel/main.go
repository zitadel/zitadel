package main

import (
	"context"
	"flag"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/login"

	"github.com/caos/logging"

	authz "github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/config"
	"github.com/caos/zitadel/internal/notification"
	tracing "github.com/caos/zitadel/internal/tracing/config"
	"github.com/caos/zitadel/pkg/admin"
	"github.com/caos/zitadel/pkg/auth"
	"github.com/caos/zitadel/pkg/console"
	"github.com/caos/zitadel/pkg/management"
)

type Config struct {
	Mgmt         management.Config
	Auth         auth.Config
	Login        login.Config
	Admin        admin.Config
	Console      console.Config
	Notification notification.Config

	Log            logging.Config
	Tracing        tracing.TracingConfig
	AuthZ          authz.Config
	SystemDefaults sd.SystemDefaults
}

func main() {
	configPaths := config.NewArrayFlags("authz.yaml", "startup.yaml", "system-defaults.yaml")
	flag.Var(configPaths, "config-files", "paths to the config files")
	managementEnabled := flag.Bool("management", true, "enable management api")
	authEnabled := flag.Bool("auth", true, "enable auth api")
	loginEnabled := flag.Bool("login", true, "enable login ui")
	adminEnabled := flag.Bool("admin", true, "enable admin api")
	consoleEnabled := flag.Bool("console", true, "enable console ui")
	notificationEnabled := flag.Bool("notification", true, "enable notification handler")
	flag.Parse()

	conf := new(Config)
	err := config.Read(conf, configPaths.Values()...)
	logging.Log("MAIN-FaF2r").OnError(err).Fatal("cannot read config")

	ctx := context.Background()
	if *managementEnabled {
		management.Start(ctx, conf.Mgmt, conf.AuthZ, conf.SystemDefaults)
	}
	var authRepo *eventsourcing.EsRepository
	if *authEnabled || *loginEnabled {
		authRepo, err = eventsourcing.Start(conf.Auth.Repository, conf.SystemDefaults)
		logging.Log("MAIN-9oRw6").OnError(err).Fatal("error starting auth repo")
	}
	if *authEnabled {
		auth.Start(ctx, conf.Auth, conf.AuthZ, conf.SystemDefaults, authRepo)
	}
	if *loginEnabled {
		login.Start(ctx, conf.Login, conf.SystemDefaults, authRepo)
	}
	if *adminEnabled {
		admin.Start(ctx, conf.Admin, conf.AuthZ, conf.SystemDefaults)
	}
	if *notificationEnabled {
		notification.Start(ctx, conf.Notification, conf.SystemDefaults)
	}
	if *consoleEnabled {
		err = console.Start(ctx, conf.Console)
		logging.Log("MAIN-3Dfuc").OnError(err).Fatal("error starting console ui")
	}
	<-ctx.Done()
	logging.Log("MAIN-s8d2h").Info("stopping zitadel")
}
