package notification

import (
	"context"
	"github.com/caos/logging"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing"
)

type Config struct {
	Repository eventsourcing.Config
}

func Start(ctx context.Context, config Config, systemDefaults sd.SystemDefaults) {
	_, err := eventsourcing.Start(config.Repository, systemDefaults)
	logging.Log("MAIN-9uBxp").OnError(err).Panic("unable to start app")
}
