package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	pol_model "github.com/caos/zitadel/internal/policy/model"
)

func (es *PolicyEventstore) GetPasswordComplexityPolicy(ctx context.Context, id string) (*pol_model.PasswordComplexityPolicy, error) {
	policy := es.policyCache.getComplexityPolicy(id)

	query := PasswordComplexityPolicyQuery(id, policy.Sequence)
	err := es_sdk.Filter(ctx, es.FilterEvents, policy.AppendEvents, query)
	if caos_errs.IsNotFound(err) && es.passwordComplexityPolicyDefault.Description != "" {
		policy.Description = es.passwordComplexityPolicyDefault.Description
		policy.MinLength = es.passwordComplexityPolicyDefault.MinLength
		policy.HasLowercase = es.passwordComplexityPolicyDefault.HasLowercase
		policy.HasUppercase = es.passwordComplexityPolicyDefault.HasUppercase
		policy.HasNumber = es.passwordComplexityPolicyDefault.HasNumber
		policy.HasSymbol = es.passwordComplexityPolicyDefault.HasSymbol
	} else if err != nil {
		return nil, err
	}
	es.policyCache.cacheComplexityPolicy(policy)
	return PasswordComplexityPolicyToModel(policy), nil
}

func (es *PolicyEventstore) CreatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	ctxData := auth.GetCtxData(ctx)
	existingPolicy, err := es.GetPasswordComplexityPolicy(ctx, ctxData.OrgID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if existingPolicy != nil && existingPolicy.Sequence > 0 {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-yDJ5I", "Errors.Policy.AlreadyExists")
	}

	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	policy.AggregateID = id

	repoPolicy := PasswordComplexityPolicyFromModel(policy)

	createAggregate := PasswordComplexityPolicyCreateAggregate(es.AggregateCreator(), repoPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoPolicy.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.policyCache.cacheComplexityPolicy(repoPolicy)
	return PasswordComplexityPolicyToModel(repoPolicy), nil
}

func (es *PolicyEventstore) UpdatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	ctxData := auth.GetCtxData(ctx)
	existingPolicy, err := es.GetPasswordComplexityPolicy(ctx, ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	if existingPolicy.Sequence <= 0 {
		return es.CreatePasswordComplexityPolicy(ctx, policy)
	}
	repoExisting := PasswordComplexityPolicyFromModel(existingPolicy)
	repoNew := PasswordComplexityPolicyFromModel(policy)

	updateAggregate := PasswordComplexityPolicyUpdateAggregate(es.AggregateCreator(), repoExisting, repoNew)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.policyCache.cacheComplexityPolicy(repoExisting)
	return PasswordComplexityPolicyToModel(repoExisting), nil
}
