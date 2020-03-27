package main

import (
	"context"
	"flag"

	"github.com/caos/logging"

	authz "github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/config"
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

	//Log //TODO: add
	//Tracing tracing.TracingConfig //TODO: add
	AuthZ authz.Config
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
		err = management.Start(ctx, conf.Mgmt, conf.AuthZ)
		logging.Log("MAIN-39Nv5").OnError(err).Fatal("error starting management api")
	}
	if *authEnabled {
		err = auth.Start(ctx, conf.Auth, conf.AuthZ)
		logging.Log("MAIN-x0nD2").OnError(err).Fatal("error starting auth api")
	}
	if *loginEnabled {
		err = login.Start(ctx, conf.Login)
		logging.Log("MAIN-53RF2").OnError(err).Fatal("error starting login ui")
	}
	if *adminEnabled {
		err = admin.Start(ctx, conf.Admin, conf.AuthZ)
		logging.Log("MAIN-0na71").OnError(err).Fatal("error starting admin api")
	}
	if *consoleEnabled {
		err = console.Start(ctx, conf.Console)
		logging.Log("MAIN-3Dfuc").OnError(err).Fatal("error starting console ui")
	}
	<-ctx.Done()
	logging.Log("MAIN-s8d2h").Info("stopping zitadel")
}
