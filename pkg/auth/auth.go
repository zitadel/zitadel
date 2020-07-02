package auth

import (
	"context"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/pkg/auth/api"
)

type Config struct {
	API        api.Config
	Repository eventsourcing.Config
}

func Start(ctx context.Context, config Config, authZRepo *authz_repo.EsRepository, authZ auth.Config, systemDefaults sd.SystemDefaults, authRepo *eventsourcing.EsRepository) {
	api.Start(ctx, config.API, authZRepo, authZ, systemDefaults, authRepo)
}
