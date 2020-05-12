package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	pol_model "github.com/caos/zitadel/internal/policy/model"
)

func (es *PolicyEventstore) GetPasswordAgePolicy(ctx context.Context, id string) (*pol_model.PasswordAgePolicy, error) {
	policy := es.policyCache.getAgePolicy(id)

	query := PasswordAgePolicyQuery(id, policy.Sequence)
	err := es_sdk.Filter(ctx, es.FilterEvents, policy.AppendEvents, query)
	if caos_errs.IsNotFound(err) {
		policy.Description = es.passwordAgePolicyDefault.Description
		policy.MaxAgeDays = es.passwordAgePolicyDefault.MaxAgeDays
		policy.ExpireWarnDays = es.passwordAgePolicyDefault.ExpireWarnDays
	} else if err != nil {
		return nil, err
	}
	es.policyCache.cacheAgePolicy(policy)
	return PasswordAgePolicyToModel(policy), nil
}

func (es *PolicyEventstore) CreatePasswordAgePolicy(ctx context.Context, policy *pol_model.PasswordAgePolicy) (*pol_model.PasswordAgePolicy, error) {
	ctxData := auth.GetCtxData(ctx)
	existingPolicy, err := es.GetPasswordAgePolicy(ctx, ctxData.OrgID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if existingPolicy != nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-yDJ5I", "Policy allready exists")
	}

	repoPolicy := PasswordAgePolicyFromModel(policy)

	createAggregate := PasswordAgePolicyCreateAggregate(es.AggregateCreator(), repoPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoPolicy.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.policyCache.cacheAgePolicy(repoPolicy)
	return PasswordAgePolicyToModel(repoPolicy), nil
}

func (es *PolicyEventstore) UpdatePasswordAgePolicy(ctx context.Context, policy *pol_model.PasswordAgePolicy) (*pol_model.PasswordAgePolicy, error) {
	ctxData := auth.GetCtxData(ctx)
	existingPolicy, err := es.GetPasswordAgePolicy(ctx, ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	repoExisting := PasswordAgePolicyFromModel(existingPolicy)
	repoNew := PasswordAgePolicyFromModel(policy)

	updateAggregate := PasswordAgePolicyUpdateAggregate(es.AggregateCreator(), repoExisting, repoNew)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.policyCache.cacheAgePolicy(repoExisting)
	return PasswordAgePolicyToModel(repoExisting), nil
}
