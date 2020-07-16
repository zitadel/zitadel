package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/service_account/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/id"
	svcacc_model "github.com/caos/zitadel/internal/service_account/model"
)

type ServiceAccountEventstore struct {
	es_int.Eventstore
	// userCache                *ServiceAccountCache
	idGenerator id.Generator
}

type ServiceAccountConfig struct {
	es_int.Eventstore
}

func StartServiceAccount(conf ServiceAccountConfig) (*ServiceAccountEventstore, error) {
	// userCache, err := StartCache(conf.Cache)
	// if err != nil {
	// 	return nil, err
	// }

	return &ServiceAccountEventstore{
		Eventstore: conf.Eventstore,
	}, nil
}

func (es *ServiceAccountEventstore) ServiceAccountByID(ctx context.Context, id string) (*svcacc_model.ServiceAccount, error) {
	serviceAccount := new(model.ServiceAccount)

	query, err := ServiceAccountByIDQuery(id, 0)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, serviceAccount.AppendEvents, query)
	if err != nil && errors.IsNotFound(err) && serviceAccount.Sequence == 0 {
		return nil, err
	}
	return model.ServiceAccountToModel(serviceAccount), nil
}
