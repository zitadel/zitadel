package management

import (
	//"context"
	//
	//"github.com/caos/logging"
	//
	//"github.com/caos/zitadel/internal/api/authz"
	//authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	//sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing"
	//"github.com/caos/zitadel/pkg/management/api"
)

type Config struct {
	Repository eventsourcing.Config
	//API        api.Config
}

//
//func Start(ctx context.Context, config Config, authZRepo *authz_repo.EsRepository, authZ authz.Config, systemDefaults sd.SystemDefaults) {
//	roles := make([]string, len(authZ.RolePermissionMappings))
//	for i, role := range authZ.RolePermissionMappings {
//		roles[i] = role.Role
//	}
//	repo, err := eventsourcing.Start(config.Repository, systemDefaults, roles)
//	logging.Log("MAIN-9uBxp").OnError(err).Panic("unable to start app")
//
//	api.Start(ctx, config.API, authZRepo, authZ, systemDefaults, repo)
//}
