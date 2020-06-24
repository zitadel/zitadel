package auth

import (
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/pkg/auth/api"
)

type Config struct {
	API        api.Config
	Repository eventsourcing.Config
}

//func Start(ctx context.Context, config Config, authZRepo *authz_repo.EsRepository, authZ authz.Config, systemDefaults sd.SystemDefaults, authRepo *eventsourcing.EsRepository) {
//	api.Start(ctx, config.API, authZRepo, authZ, systemDefaults, authRepo)
//}
