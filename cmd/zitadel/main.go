package main

import (
	"context"
	"flag"

	"github.com/caos/logging"

	admin_es "github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api"
	internal_authz "github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/admin"
	"github.com/caos/zitadel/internal/api/grpc/auth"
	"github.com/caos/zitadel/internal/api/grpc/management"
	"github.com/caos/zitadel/internal/api/oidc"
	auth_es "github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/authz"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	mgmt_es "github.com/caos/zitadel/internal/management/repository/eventsourcing"
	"github.com/caos/zitadel/internal/notification"
	tracing "github.com/caos/zitadel/internal/tracing/config"
	"github.com/caos/zitadel/internal/ui"
	"github.com/caos/zitadel/internal/ui/console"
	"github.com/caos/zitadel/internal/ui/login"
)

type Config struct {
	Log            logging.Config
	Tracing        tracing.TracingConfig
	InternalAuthZ  internal_authz.Config
	SystemDefaults sd.SystemDefaults

	AuthZ authz.Config
	Auth  auth_es.Config
	Admin admin_es.Config
	Mgmt  mgmt_es.Config

	API api.Config
	UI  ui.Config

	Notification notification.Config
}

var (
	configPaths         = config.NewArrayFlags("authz.yaml", "startup.yaml", "system-defaults.yaml")
	adminEnabled        = flag.Bool("admin", true, "enable admin api")
	managementEnabled   = flag.Bool("management", true, "enable management api")
	authEnabled         = flag.Bool("auth", true, "enable auth api")
	oidcEnabled         = flag.Bool("oidc", true, "enable oidc api")
	loginEnabled        = flag.Bool("login", true, "enable login ui")
	consoleEnabled      = flag.Bool("console", true, "enable console ui")
	notificationEnabled = flag.Bool("notification", true, "enable notification handler")
)

func main() {
	flag.Var(configPaths, "config-files", "paths to the config files")
	flag.Parse()

	conf := new(Config)
	err := config.Read(conf, configPaths.Values()...)
	logging.Log("MAIN-FaF2r").OnError(err).Fatal("cannot read config")

	ctx := context.Background()
	authZRepo, err := authz.Start(ctx, conf.AuthZ, conf.InternalAuthZ, conf.SystemDefaults)
	logging.Log("MAIN-s9KOw").OnError(err).Fatal("error starting authz repo")
	var authRepo *auth_es.EsRepository
	if *authEnabled || *oidcEnabled || *loginEnabled {
		authRepo, err = auth_es.Start(conf.Auth, conf.InternalAuthZ, conf.SystemDefaults, authZRepo)
		logging.Log("MAIN-9oRw6").OnError(err).Fatal("error starting auth repo")
	}

	startAPI(ctx, conf, authZRepo, authRepo)
	startUI(ctx, conf, authRepo)

	if *notificationEnabled {
		notification.Start(ctx, conf.Notification, conf.SystemDefaults)
	}

	<-ctx.Done()
	logging.Log("MAIN-s8d2h").Info("stopping zitadel")
}

func startUI(ctx context.Context, conf *Config, authRepo *auth_es.EsRepository) {
	uis := ui.Create(conf.UI)
	if *loginEnabled {
		uis.RegisterHandler(ui.LoginHandler, login.Start(conf.UI.Login, authRepo, ui.LoginHandler).Handler())
	}
	if *consoleEnabled {
		consoleHandler, err := console.Start(conf.UI.Console)
		logging.Log("API-AGD1f").OnError(err).Fatal("error starting console")
		uis.RegisterHandler(ui.ConsoleHandler, consoleHandler)
	}
	uis.Start(ctx)
}

func startAPI(ctx context.Context, conf *Config, authZRepo *authz_repo.EsRepository, authRepo *auth_es.EsRepository) {
	apis := api.Create(conf.API, conf.InternalAuthZ, authZRepo, conf.SystemDefaults)
	roles := make([]string, len(conf.InternalAuthZ.RolePermissionMappings))
	for i, role := range conf.InternalAuthZ.RolePermissionMappings {
		roles[i] = role.Role
	}
	if *adminEnabled {
		adminRepo, err := admin_es.Start(ctx, conf.Admin, conf.SystemDefaults, roles)
		logging.Log("API-D42tq").OnError(err).Fatal("error starting auth repo")
		apis.RegisterServer(ctx, admin.CreateServer(adminRepo))
	}
	if *managementEnabled {
		managementRepo, err := mgmt_es.Start(conf.Mgmt, conf.SystemDefaults, roles)
		logging.Log("API-Gd2qq").OnError(err).Fatal("error starting management repo")
		apis.RegisterServer(ctx, management.CreateServer(managementRepo, conf.SystemDefaults))
	}
	if *authEnabled {
		apis.RegisterServer(ctx, auth.CreateServer(authRepo))
	}
	if *oidcEnabled {
		op := oidc.NewProvider(ctx, conf.API.OIDC, authRepo)
		apis.RegisterHandler("/oauth/v2", op.HttpHandler().Handler)
	}
	apis.Start(ctx)
}
