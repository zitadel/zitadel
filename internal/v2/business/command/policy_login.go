package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_login "github.com/caos/zitadel/internal/v2/repository/iam/policy/login"
	iam_factor "github.com/caos/zitadel/internal/v2/repository/iam/policy/login/factors"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/login/idpprovider"
	"github.com/caos/zitadel/internal/v2/repository/policy/login"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/factors"
	idpprovider2 "github.com/caos/zitadel/internal/v2/repository/policy/login/idpprovider"
)

func (r *CommandSide) AddDefaultLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-5Mv0s", "Errors.IAM.LoginPolicyInvalid")
	}

	addedPolicy := iam_login.NewWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.AlreadyExists")
	}

	iamAgg := AggregateFromWriteModel(&addedPolicy.WriteModel.WriteModel).
		PushLoginPolicyAddedEvent(ctx, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIdp, policy.ForceMFA, login.PasswordlessType(policy.PasswordlessType))

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLoginPolicy(addedPolicy), nil
}

func (r *CommandSide) ChangeDefaultLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-6M0od", "Errors.IAM.LoginPolicyInvalid")
	}

	existingPolicy, err := r.defaultLoginPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	iamAgg := AggregateFromWriteModel(&existingPolicy.WriteModel.WriteModel).
		PushLoginPolicyChangedFromExisting(ctx, existingPolicy, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIdp, policy.ForceMFA, login.PasswordlessType(policy.PasswordlessType))

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLoginPolicy(existingPolicy), nil
}

func (r *CommandSide) AddIDPProviderToDefaultLoginPolicy(ctx context.Context, idpProvider *iam_model.IDPProvider) (*iam_model.IDPProvider, error) {
	writeModel := idpprovider.NewWriteModel(idpProvider.AggregateID, idpProvider.IDPConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	aggregate := AggregateFromWriteModel(&writeModel.WriteModel.WriteModel).
		PushLoginPolicyIDPProviderAddedEvent(ctx, idpProvider.IDPConfigID, idpprovider2.Type(idpProvider.Type))

	if err = r.eventstore.PushAggregate(ctx, writeModel, aggregate); err != nil {
		return nil, err
	}

	return writeModelToIDPProvider(writeModel), nil
}

func (r *CommandSide) RemoveIDPProviderFromDefaultLoginPolicy(ctx context.Context, idpProvider *iam_model.IDPProvider) error {
	writeModel := idpprovider.NewWriteModel(idpProvider.AggregateID, idpProvider.IDPConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	aggregate := AggregateFromWriteModel(&writeModel.WriteModel.WriteModel).
		PushLoginPolicyIDPProviderAddedEvent(ctx, idpProvider.IDPConfigID, idpprovider2.Type(idpProvider.Type))

	return r.eventstore.PushAggregate(ctx, writeModel, aggregate)
}

func (r *CommandSide) AddSecondFactorToDefaultLoginPolicy(ctx context.Context, iamID string, secondFactor iam_model.SecondFactorType) (iam_model.SecondFactorType, error) {
	writeModel := iam_factor.NewSecondFactorWriteModel(iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return iam_model.SecondFactorTypeUnspecified, err
	}

	aggregate := AggregateFromWriteModel(&writeModel.SecondFactorWriteModel.WriteModel).
		PushLoginPolicySecondFactorAdded(ctx, factors.SecondFactorType(secondFactor))

	if err = r.eventstore.PushAggregate(ctx, writeModel, aggregate); err != nil {
		return iam_model.SecondFactorTypeUnspecified, err
	}

	return iam_model.SecondFactorType(writeModel.MFAType), nil
}

func (r *CommandSide) RemoveSecondFactorFromDefaultLoginPolicy(ctx context.Context, iamID string, secondFactor iam_model.SecondFactorType) error {
	writeModel := iam_factor.NewSecondFactorWriteModel(iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	aggregate := AggregateFromWriteModel(&writeModel.SecondFactorWriteModel.WriteModel).
		PushLoginPolicySecondFactorRemoved(ctx, factors.SecondFactorType(secondFactor))

	return r.eventstore.PushAggregate(ctx, writeModel, aggregate)
}

func (r *CommandSide) AddMultiFactorToDefaultLoginPolicy(ctx context.Context, iamID string, secondFactor iam_model.MultiFactorType) (iam_model.MultiFactorType, error) {
	writeModel := iam_factor.NewMultiFactorWriteModel(iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return iam_model.MultiFactorTypeUnspecified, err
	}

	aggregate := AggregateFromWriteModel(&writeModel.MultiFactoryWriteModel.WriteModel).
		PushLoginPolicyMultiFactorAdded(ctx, factors.MultiFactorType(secondFactor))

	if err = r.eventstore.PushAggregate(ctx, writeModel, aggregate); err != nil {
		return iam_model.MultiFactorTypeUnspecified, err
	}

	return iam_model.MultiFactorType(writeModel.MultiFactoryWriteModel.MFAType), nil
}

func (r *CommandSide) RemoveMultiFactorFromDefaultLoginPolicy(ctx context.Context, iamID string, secondFactor iam_model.MultiFactorType) error {
	writeModel := iam_factor.NewMultiFactorWriteModel(iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	aggregate := AggregateFromWriteModel(&writeModel.MultiFactoryWriteModel.WriteModel).
		PushLoginPolicyMultiFactorRemoved(ctx, factors.MultiFactorType(secondFactor))

	return r.eventstore.PushAggregate(ctx, writeModel, aggregate)
}

func (r *CommandSide) defaultLoginPolicyWriteModelByID(ctx context.Context, iamID string) (policy *iam_login.WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := iam_login.NewWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
