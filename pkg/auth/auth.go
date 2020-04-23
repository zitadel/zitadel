package auth

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	app "github.com/caos/zitadel/internal/auth"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/pkg/auth/api"
)

type Config struct {
	App app.Config
	API api.Config
}

func Start(ctx context.Context, config Config, authZ auth.Config, systemDefaults sd.SystemDefaults) {
	api.Start(ctx, config.API)
}
