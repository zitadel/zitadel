package iam

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
	iam_login "github.com/caos/zitadel/internal/v2/repository/iam/policy/login"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/login/idpprovider"
	iam_multi_factor "github.com/caos/zitadel/internal/v2/repository/iam/policy/login/multi_factors"
	iam_second_factor "github.com/caos/zitadel/internal/v2/repository/iam/policy/login/second_factors"
	"github.com/caos/zitadel/internal/v2/repository/policy/login"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/multi_factors"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/second_factors"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

func (r *Repository) AddLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-5Mv0s", "Errors.IAM.LoginPolicyInvalid")
	}

	addedPolicy := iam_login.NewLoginPolicyWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.AlreadyExists")
	}

	iamAgg := iam_repo.AggregateFromWriteModel(&addedPolicy.WriteModel).
		PushLoginPolicyAddedEvent(ctx, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIdp, policy.ForceMFA, login.PasswordlessType(policy.PasswordlessType))

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLoginPolicy(addedPolicy), nil
}

func (r *Repository) ChangeLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-6M0od", "Errors.IAM.LoginPolicyInvalid")
	}

	existingPolicy, err := r.loginPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	iamAgg := iam_repo.AggregateFromWriteModel(&existingPolicy.WriteModel).
		PushLoginPolicyChangedFromExisting(ctx, existingPolicy, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIdp, policy.ForceMFA, login.PasswordlessType(policy.PasswordlessType))

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLoginPolicy(existingPolicy), nil
}

func (r *Repository) AddIDPProviderToLoginPolicy(ctx context.Context, idpProvider *iam_model.IDPProvider) (*iam_model.IDPProvider, error) {
	writeModel := idpprovider.NewLoginPolicyIDPProviderWriteModel(idpProvider.AggregateID, idpProvider.IdpConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	aggregate := iam_repo.AggregateFromWriteModel(&writeModel.WriteModel).
		PushLoginPolicyIDPProviderAddedEvent(ctx, idpProvider.IdpConfigID, provider.Type(idpProvider.Type))

	if err = r.eventstore.PushAggregate(ctx, writeModel, aggregate); err != nil {
		return nil, err
	}

	return writeModelToIDPProvider(writeModel), nil
}

func (r *Repository) RemoveIDPProviderFromLoginPolicy(ctx context.Context, idpProvider *iam_model.IDPProvider) error {
	writeModel := idpprovider.NewLoginPolicyIDPProviderWriteModel(idpProvider.AggregateID, idpProvider.IdpConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	aggregate := iam_repo.AggregateFromWriteModel(&writeModel.WriteModel).
		PushLoginPolicyIDPProviderAddedEvent(ctx, idpProvider.IdpConfigID, provider.Type(idpProvider.Type))

	return r.eventstore.PushAggregate(ctx, writeModel, aggregate)
}

func (r *Repository) AddSecondFactorToLoginPolicy(ctx context.Context, iamID string, secondFactor iam_model.SecondFactorType) (iam_model.SecondFactorType, error) {
	writeModel := iam_second_factor.NewSecondFactorWriteModel(iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return iam_model.SecondFactorTypeUnspecified, err
	}

	aggregate := iam_repo.AggregateFromWriteModel(&writeModel.WriteModel).
		PushLoginPolicySecondFactorAdded(ctx, second_factors.SecondFactorType(secondFactor))

	if err = r.eventstore.PushAggregate(ctx, writeModel, aggregate); err != nil {
		return iam_model.SecondFactorTypeUnspecified, err
	}

	return iam_model.SecondFactorType(writeModel.SecondFactor.MFAType), nil
}

func (r *Repository) RemoveSecondFactorFromLoginPolicy(ctx context.Context, iamID string, secondFactor iam_model.SecondFactorType) error {
	writeModel := iam_second_factor.NewSecondFactorWriteModel(iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	aggregate := iam_repo.AggregateFromWriteModel(&writeModel.WriteModel).
		PushLoginPolicySecondFactorRemoved(ctx, second_factors.SecondFactorType(secondFactor))

	return r.eventstore.PushAggregate(ctx, writeModel, aggregate)
}

func (r *Repository) AddMultiFactorToLoginPolicy(ctx context.Context, iamID string, secondFactor iam_model.MultiFactorType) (iam_model.MultiFactorType, error) {
	writeModel := iam_multi_factor.NewMultiFactorWriteModel(iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return iam_model.MultiFactorTypeUnspecified, err
	}

	aggregate := iam_repo.AggregateFromWriteModel(&writeModel.WriteModel).
		PushLoginPolicyMultiFactorAdded(ctx, multi_factors.MultiFactorType(secondFactor))

	if err = r.eventstore.PushAggregate(ctx, writeModel, aggregate); err != nil {
		return iam_model.MultiFactorTypeUnspecified, err
	}

	return iam_model.MultiFactorType(writeModel.MultiFactor.MFAType), nil
}

func (r *Repository) RemoveMultiFactorFromLoginPolicy(ctx context.Context, iamID string, secondFactor iam_model.MultiFactorType) error {
	writeModel := iam_multi_factor.NewMultiFactorWriteModel(iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	aggregate := iam_repo.AggregateFromWriteModel(&writeModel.WriteModel).
		PushLoginPolicyMultiFactorRemoved(ctx, multi_factors.MultiFactorType(secondFactor))

	return r.eventstore.PushAggregate(ctx, writeModel, aggregate)
}

func (r *Repository) loginPolicyWriteModelByID(ctx context.Context, iamID string) (policy *iam_login.LoginPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := iam_login.NewLoginPolicyWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
