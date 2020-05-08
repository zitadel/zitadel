package auth

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/pkg/auth/api"
)

type Config struct {
	API        api.Config
	Repository eventsourcing.Config
}

func Start(ctx context.Context, config Config, authZ auth.Config, systemDefaults sd.SystemDefaults) {
	repo, err := eventsourcing.Start(config.Repository, systemDefaults)
	logging.Log("MAIN-9uBxp").OnError(err).Panic("unable to start app")

	api.Start(ctx, config.API, authZ, repo)
}
