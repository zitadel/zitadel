package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type IAMEventstore struct {
	es_int.Eventstore
	iamCache *IAMCache
}

type IAMConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartIAM(conf IAMConfig, systemDefaults sd.SystemDefaults) (*IAMEventstore, error) {
	iamCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}

	return &IAMEventstore{
		Eventstore: conf.Eventstore,
		iamCache:   iamCache,
	}, nil
}

func (es *IAMEventstore) IAMByID(ctx context.Context, id string) (_ *iam_model.IAM, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	iam := es.iamCache.getIAM(id)

	query, err := IAMByIDQuery(iam.AggregateID, iam.Sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, iam.AppendEvents, query)
	if err != nil && caos_errs.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	es.iamCache.cacheIAM(iam)
	return model.IAMToModel(iam), nil
}

func (es *IAMEventstore) IAMEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error) {
	query, err := IAMByIDQuery(id, sequence)
	if err != nil {
		return nil, err
	}
	return es.FilterEvents(ctx, query)
}

func (es *IAMEventstore) GetIDPConfig(ctx context.Context, aggregateID, idpConfigID string) (*iam_model.IDPConfig, error) {
	existing, err := es.IAMByID(ctx, aggregateID)
	if err != nil {
		return nil, err
	}
	if _, existingIDP := existing.GetIDP(idpConfigID); existingIDP != nil {
		return existingIDP, nil
	}
	return nil, caos_errs.ThrowNotFound(nil, "EVENT-Scj8s", "Errors.IAM.IdpNotExisting")
}

func (es *IAMEventstore) GetOrgIAMPolicy(ctx context.Context, iamID string) (*iam_model.OrgIAMPolicy, error) {
	existingIAM, err := es.IAMByID(ctx, iamID)
	if err != nil {
		return nil, err
	}
	if existingIAM.DefaultOrgIAMPolicy == nil {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-2Fj8s", "Errors.IAM.OrgIAMPolicy.NotExisting")
	}
	return existingIAM.DefaultOrgIAMPolicy, nil
}
