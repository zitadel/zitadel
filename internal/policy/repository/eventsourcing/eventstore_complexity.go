package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
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

	query := PasswordComplexityPolicyQuery(id, policy.Sequence)
	err := es_sdk.Filter(ctx, es.FilterEvents, policy.AppendEvents, query)
	if err != nil {
		return nil, err
	}
	es.policyCache.cachePolicy(policy)
	return PasswordComplexityPolicyToModel(policy), nil
}

func (es *PolicyEventstore) CreatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Description is required")
	}
	ctxData := auth.GetCtxData(ctx)
	existingPolicy, err := es.GetPasswordComplexityPolicy(ctx, ctxData.OrgID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if existingPolicy != nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-yDJ5I", "Policy allready exists")
	}

	repoPolicy := PasswordComplexityPolicyFromModel(policy)

	createAggregate := PasswordComplexityPolicyCreateAggregate(es.AggregateCreator(), repoPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoPolicy.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.policyCache.cachePolicy(repoPolicy)
	return PasswordComplexityPolicyToModel(repoPolicy), nil
}

func (es *PolicyEventstore) UpdatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Description is required")
	}
	ctxData := auth.GetCtxData(ctx)
	existingPolicy, err := es.GetPasswordComplexityPolicy(ctx, ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	repoExisting := PasswordComplexityPolicyFromModel(existingPolicy)
	repoNew := PasswordComplexityPolicyFromModel(policy)

	updateAggregate := PasswordComplexityPolicyUpdateAggregate(es.AggregateCreator(), repoExisting, repoNew)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.policyCache.cachePolicy(repoExisting)
	return PasswordComplexityPolicyToModel(repoExisting), nil
}
