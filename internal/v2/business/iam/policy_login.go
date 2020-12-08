package iam

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
	iam_login "github.com/caos/zitadel/internal/v2/repository/iam/policy/login"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/login/idpprovider"
	"github.com/caos/zitadel/internal/v2/repository/policy/login"

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
