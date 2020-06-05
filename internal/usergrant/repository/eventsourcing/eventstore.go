package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/cache/config"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/id"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
)

type UserGrantEventStore struct {
	es_int.Eventstore
	userGrantCache *UserGrantCache
	idGenerator    id.Generator
}

type UserGrantConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartUserGrant(conf UserGrantConfig) (*UserGrantEventStore, error) {
	userGrantCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}
	return &UserGrantEventStore{
		Eventstore:     conf.Eventstore,
		userGrantCache: userGrantCache,
		idGenerator:    id.SonyFlakeGenerator,
	}, nil
}

func (es *UserGrantEventStore) UserGrantByID(ctx context.Context, id string) (*grant_model.UserGrant, error) {
	grant := es.userGrantCache.getUserGrant(id)

	query, err := UserGrantByIDQuery(grant.AggregateID, grant.Sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, grant.AppendEvents, query)
	if err != nil && caos_errs.IsNotFound(err) && grant.Sequence == 0 {
		return nil, err
	}
	es.userGrantCache.cacheUserGrant(grant)
	if grant.State == int32(grant_model.USERGRANTSTATE_REMOVED) {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-2ks8d", "UserGrant not found")
	}
	return model.UserGrantToModel(grant), nil
}

func (es *UserGrantEventStore) AddUserGrant(ctx context.Context, grant *grant_model.UserGrant) (*grant_model.UserGrant, error) {
	if grant == nil || !grant.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-sdiw3", "User grant invalid")
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	grant.AggregateID = id

	repoGrant := model.UserGrantFromModel(grant)

	addAggregates, err := UserGrantAddedAggregate(ctx, es.Eventstore.AggregateCreator(), repoGrant)
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoGrant.AppendEvents, addAggregates...)
	if err != nil {
		return nil, err
	}
	return model.UserGrantToModel(repoGrant), nil
}

func (es *UserGrantEventStore) ChangeUserGrant(ctx context.Context, grant *grant_model.UserGrant) (*grant_model.UserGrant, error) {
	if grant == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-lo0s9", "invalid grant")
	}
	existing, err := es.UserGrantByID(ctx, grant.AggregateID)
	if err != nil {
		return nil, err
	}
	repoExisting := model.UserGrantFromModel(existing)
	repoGrant := model.UserGrantFromModel(grant)

	projectAggregate := UserGrantChangedAggregate(es.Eventstore.AggregateCreator(), repoExisting, repoGrant)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.userGrantCache.cacheUserGrant(repoExisting)
	return model.UserGrantToModel(repoExisting), nil
}

func (es *UserGrantEventStore) RemoveUserGrant(ctx context.Context, grantID string) error {
	existing, err := es.UserGrantByID(ctx, grantID)
	if err != nil {
		return err
	}
	repoExisting := model.UserGrantFromModel(existing)
	repoGrant := &model.UserGrant{ObjectRoot: models.ObjectRoot{AggregateID: grantID}}
	projectAggregates, err := UserGrantRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoExisting, repoGrant)
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoExisting.AppendEvents, projectAggregates...)
	if err != nil {
		return err
	}
	es.userGrantCache.cacheUserGrant(repoExisting)
	return nil
}

func (es *UserGrantEventStore) DeactivateUserGrant(ctx context.Context, grantID string) (*grant_model.UserGrant, error) {
	if grantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-8si34", "grantID missing")
	}
	existing, err := es.UserGrantByID(ctx, grantID)
	if err != nil {
		return nil, err
	}
	if !existing.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-lo9sw", "deactivate only possible for active grant")
	}
	repoExisting := model.UserGrantFromModel(existing)
	repoGrant := &model.UserGrant{ObjectRoot: models.ObjectRoot{AggregateID: grantID}}

	projectAggregate := UserGrantDeactivatedAggregate(es.Eventstore.AggregateCreator(), repoExisting, repoGrant)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.userGrantCache.cacheUserGrant(repoGrant)
	return model.UserGrantToModel(repoExisting), nil
}

func (es *UserGrantEventStore) ReactivateUserGrant(ctx context.Context, grantID string) (*grant_model.UserGrant, error) {
	if grantID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-sksiw", "grantID missing")
	}
	existing, err := es.UserGrantByID(ctx, grantID)
	if err != nil {
		return nil, err
	}
	if !existing.IsInactive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-lo9sw", "reactivate only possible for inactive grant")
	}
	repoExisting := model.UserGrantFromModel(existing)
	repoGrant := &model.UserGrant{ObjectRoot: models.ObjectRoot{AggregateID: grantID}}

	projectAggregate := UserGrantReactivatedAggregate(es.Eventstore.AggregateCreator(), repoExisting, repoGrant)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.userGrantCache.cacheUserGrant(repoExisting)
	return model.UserGrantToModel(repoExisting), nil
}
