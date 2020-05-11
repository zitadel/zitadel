package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

type IamEventstore struct {
	es_int.Eventstore
	iamCache *IamCache
}

type IamConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartIam(conf IamConfig, systemDefaults sd.SystemDefaults) (*IamEventstore, error) {
	iamCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}

	return &IamEventstore{
		Eventstore: conf.Eventstore,
		iamCache:   iamCache,
	}, nil
}

func (es *IamEventstore) IamByID(ctx context.Context, id string) (*iam_model.Iam, error) {
	iam := es.iamCache.getIam(id)

	query, err := IamByIDQuery(iam.AggregateID, iam.Sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, iam.AppendEvents, query)
	if err != nil && caos_errs.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	es.iamCache.cacheIam(iam)
	return model.IamToModel(iam), nil
}

func (es *IamEventstore) StartSetup(ctx context.Context, iamID string) (*iam_model.Iam, error) {
	iam, err := es.IamByID(ctx, iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}

	if iam != nil && iam.SetUpStarted {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9so34", "Setup already started")
	}

	repoIam := &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: iamID}}
	createAggregate := IamSetupStartedAggregate(es.AggregateCreator(), repoIam)
	err = es_sdk.Push(ctx, es.PushAggregates, repoIam.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.iamCache.cacheIam(repoIam)
	return model.IamToModel(repoIam), nil
}
