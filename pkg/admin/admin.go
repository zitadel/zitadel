package admin

import (
	"context"
	"github.com/caos/logging"

	app "github.com/caos/zitadel/internal/admin"
	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/pkg/admin/api"
)

type Config struct {
	App app.Config
	API api.Config
}

func Start(ctx context.Context, config Config, authZ auth.Config) error {
	err := api.Start(ctx, config.API)
	logging.Log("MAIN-lfo5h").OnError(err).Panic("unable to start api")
	return err
}
