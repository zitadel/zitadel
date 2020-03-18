package main

import (
	"context"
	"flag"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/config"
	"github.com/caos/zitadel/pkg/admin"
	"github.com/caos/zitadel/pkg/auth"
	"github.com/caos/zitadel/pkg/eventstore"
	"github.com/caos/zitadel/pkg/management"
)

type Config struct {
	Eventstore eventstore.Config
	Management management.Config
	Auth       auth.Config
	Admin      admin.Config
}

func main() {
	configPath := flag.String("config-file", "/zitadel/config/startup.yaml", "path to the config file")
	eventstoreEnabled := flag.Bool("eventstore", true, "enable eventstore")
	managementEnabled := flag.Bool("management", true, "enable management api")
	authEnabled := flag.Bool("auth", true, "enable auth api")
	adminEnabled := flag.Bool("admin", true, "enable admin api")

	flag.Parse()

	conf := new(Config)
	err := config.Read(conf, *configPath)
	logging.Log("MAIN-FaF2r").OnError(err).Fatal("cannot read config")

	ctx := context.Background()
	if *eventstoreEnabled {
		err = eventstore.Start(ctx, conf.Eventstore)
		logging.Log("MAIN-sj2Sd").OnError(err).Fatal("error starting eventstore")
	}
	if *managementEnabled {
		err = management.Start(ctx, conf.Management)
		logging.Log("MAIN-39Nv5").OnError(err).Fatal("error starting management api")
	}
	if *authEnabled {
		err = auth.Start(ctx, conf.Auth)
		logging.Log("MAIN-x0nD2").OnError(err).Fatal("error starting auth api")
	}
	if *adminEnabled {
		err = admin.Start(ctx, conf.Admin)
		logging.Log("MAIN-0na71").OnError(err).Fatal("error starting admin api")
	}
	<-ctx.Done()
	logging.Log("MAIN-s8d2h").Info("stopping zitadel")
}
