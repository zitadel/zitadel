package admin

import (
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/pkg/admin/api"
)

type Config struct {
	Repository eventsourcing.Config
	API        api.Config
}

//func Start(ctx context.Context, config Config, authZRepo *authz_repo.EsRepository, authZ auth.Config, systemDefaults sd.SystemDefaults) {
//	roles := make([]string, len(authZ.RolePermissionMappings))
//	for i, role := range authZ.RolePermissionMappings {
//		roles[i] = role.Role
//	}
//
//	repo, err := eventsourcing.Start(ctx, config.Repository, systemDefaults, roles)
//	logging.Log("MAIN-9uBxp").OnError(err).Panic("unable to start app")
//
//	api.Start(ctx, config.API, authZRepo, authZ, systemDefaults, repo)
//}
