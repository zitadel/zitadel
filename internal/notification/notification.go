package notification

import (
	"context"
	"github.com/caos/logging"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing"
	"github.com/caos/zitadel/internal/v2/command"
	"github.com/rakyll/statik/fs"

	_ "github.com/caos/zitadel/internal/notification/statik"
)

type Config struct {
	Repository eventsourcing.Config
}

func Start(ctx context.Context, config Config, systemDefaults sd.SystemDefaults, command *command.CommandSide) {
	statikFS, err := fs.NewWithNamespace("notification")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start listener")

	_, err = eventsourcing.Start(config.Repository, statikFS, systemDefaults, command)
	logging.Log("MAIN-9uBxp").OnError(err).Panic("unable to start app")
}
