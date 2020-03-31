package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/management/repository"
	"github.com/caos/zitadel/pkg/management/api"
)

type Config struct {
	App repository.Config
	API api.Config
}

func Start(ctx context.Context, config Config, authZ auth.Config) {
	api.Start(ctx, config.API)
}
