package auth

import (
	"context"
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/auth"
	app "github.com/caos/zitadel/internal/auth"
	"github.com/caos/zitadel/pkg/auth/api"
)

type Config struct {
	App *app.Config
	API *api.Config
}

func Start(ctx context.Context, config *Config, authZ *auth.Config) error {
	err := api.Start(ctx, config.API)
	logging.Log("MAIN-BmOLI").OnError(err).Panic("unable to start api")
	return err
}
