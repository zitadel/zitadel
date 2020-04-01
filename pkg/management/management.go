package management

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing"
	"github.com/caos/zitadel/pkg/management/api"
)

type Config struct {
	Repository eventsourcing.Config
	API        api.Config
}

func Start(ctx context.Context, config Config, authZ auth.Config) {
	repo, err := eventsourcing.Start(config.Repository)
	logging.Log("MAIN-9uBxp").OnError(err).Panic("unable to start app")

	api.Start(ctx, config.API, repo)
	eventsourcing.Start(config.Repository)
}
