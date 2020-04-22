package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/cache/config"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	pol_model "github.com/caos/zitadel/internal/policy/model"
)

type PolicyEventstore struct {
	es_int.Eventstore
	policyCache *PolicyCache
}

type PolicyConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartPolicy(conf PolicyConfig) (*PolicyEventstore, error) {
	policyCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}
	return &PolicyEventstore{
		Eventstore:  conf.Eventstore,
		policyCache: policyCache,
	}, nil
}

func (es *PolicyEventstore) GetPasswordComplexityPolicy(ctx context.Context, id string) (*pol_model.PasswordComplexityPolicy, error) {
	policy := es.policyCache.getPolicy(id)

	query := PolicyQuery(policy.Sequence)
	err := es_sdk.Filter(ctx, es.FilterEvents, policy.AppendEvents, query)
	if err != nil {
		return nil, err
	}
	es.policyCache.cachePolicy(policy)
	return PolicyToModel(policy), nil
}

func (es *PolicyEventstore) CreatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	if !policy.IsValid() { // Brauchts das????? war bei Project ob Name vorganden
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Description is required")
	}
	repoPolicy := PolicyFromModel(policy)

	createAggregate := PolicyCreateAggregate(es.AggregateCreator(), repoPolicy)
	err := es_sdk.Push(ctx, es.PushAggregates, repoPolicy.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.policyCache.cachePolicy(repoPolicy)
	return PolicyToModel(repoPolicy), nil
}

func (es *PolicyEventstore) UpdatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	if !policy.IsValid() { // Brauchts das?????  war bei Project ob Name vorganden
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Description is required")
	}
	existingPolicy, err := es.GetPasswordComplexityPolicy(ctx, policy.ID)
	if err != nil {
		return nil, err
	}
	repoExisting := PolicyFromModel(existingPolicy)
	repoNew := PolicyFromModel(policy)

	updateAggregate := PolicyUpdateAggregate(es.AggregateCreator(), repoExisting, repoNew)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.policyCache.cachePolicy(repoExisting)
	return PolicyToModel(repoExisting), nil
}
