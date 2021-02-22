package eventsourcing

import (
	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
)

const (
	projectOwnerRole       = "PROJECT_OWNER"
	projectOwnerGlobalRole = "PROJECT_OWNER_GLOBAL"
)

type ProjectEventstore struct {
	es_int.Eventstore
	projectCache  *ProjectCache
	passwordAlg   crypto.HashAlgorithm
	pwGenerator   crypto.Generator
	idGenerator   id.Generator
	ClientKeySize int
}

type ProjectConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartProject(conf ProjectConfig, systemDefaults sd.SystemDefaults) (*ProjectEventstore, error) {
	projectCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}
	passwordAlg := crypto.NewBCrypt(systemDefaults.SecretGenerators.PasswordSaltCost)
	pwGenerator := crypto.NewHashGenerator(systemDefaults.SecretGenerators.ClientSecretGenerator, passwordAlg)
	return &ProjectEventstore{
		Eventstore:    conf.Eventstore,
		projectCache:  projectCache,
		passwordAlg:   passwordAlg,
		pwGenerator:   pwGenerator,
		idGenerator:   id.SonyFlakeGenerator,
		ClientKeySize: int(systemDefaults.SecretGenerators.ClientKeySize),
	}, nil
}

//
//func (es *ProjectEventstore) ProjectByID(ctx context.Context, id string) (*proj_model.Project, error) {
//	project := es.projectCache.getProject(id)
//
//	query, err := ProjectByIDQuery(project.AggregateID, project.Sequence)
//	if err != nil {
//		return nil, err
//	}
//	err = es_sdk.Filter(ctx, es.FilterEvents, project.AppendEvents, query)
//	if err != nil && !(caos_errs.IsNotFound(err) && project.Sequence != 0) {
//		return nil, err
//	}
//	if project.State == int32(proj_model.ProjectStateRemoved) {
//		return nil, caos_errs.ThrowNotFound(nil, "EVENT-dG8ie", "Errors.Project.NotFound")
//	}
//	es.projectCache.cacheProject(project)
//	return model.ProjectToModel(project), nil
//}
