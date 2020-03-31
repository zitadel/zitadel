package auth

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	app "github.com/caos/zitadel/internal/auth"
	"github.com/caos/zitadel/pkg/auth/api"
)

type Config struct {
	App app.Config
	API api.Config
}

func Start(ctx context.Context, config Config, authZ auth.Config) {
	api.Start(ctx, config.API)
}
