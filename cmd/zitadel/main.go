package main

import (
	"context"
	"flag"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"

	"github.com/caos/logging"

	authz "github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/config"
	tracing "github.com/caos/zitadel/internal/tracing/config"
	"github.com/caos/zitadel/pkg/admin"
	"github.com/caos/zitadel/pkg/auth"
	"github.com/caos/zitadel/pkg/console"
	"github.com/caos/zitadel/pkg/login"
	"github.com/caos/zitadel/pkg/management"
)

type Config struct {
	Mgmt    management.Config
	Auth    auth.Config
	Login   login.Config
	Admin   admin.Config
	Console console.Config

	Log            logging.Config
	Tracing        tracing.TracingConfig
	AuthZ          authz.Config
	SystemDefaults sd.SystemDefaults
}

func main() {
	var configPaths config.ArrayFlags
	flag.Var(&configPaths, "config-files", "path to the config files")
	managementEnabled := flag.Bool("management", true, "enable management api")
	authEnabled := flag.Bool("auth", true, "enable auth api")
	loginEnabled := flag.Bool("login", true, "enable login ui")
	adminEnabled := flag.Bool("admin", true, "enable admin api")
	consoleEnabled := flag.Bool("console", true, "enable console ui")
	flag.Parse()

	conf := new(Config)
	err := config.Read(conf, configPaths...)
	logging.Log("MAIN-FaF2r").OnError(err).Fatal("cannot read config")

	ctx := context.Background()
	if *managementEnabled {
		management.Start(ctx, conf.Mgmt, conf.AuthZ, conf.SystemDefaults)
	}
	if *authEnabled {
		auth.Start(ctx, conf.Auth, conf.AuthZ, conf.SystemDefaults)
	}
	if *loginEnabled {
		err = login.Start(ctx, conf.Login)
		logging.Log("MAIN-53RF2").OnError(err).Fatal("error starting login ui")
	}
	if *adminEnabled {
		admin.Start(ctx, conf.Admin, conf.AuthZ)
	}
	if *consoleEnabled {
		err = console.Start(ctx, conf.Console)
		logging.Log("MAIN-3Dfuc").OnError(err).Fatal("error starting console ui")
	}
	<-ctx.Done()
	logging.Log("MAIN-s8d2h").Info("stopping zitadel")
}
