package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) GetDefaultLoginPolicy(ctx context.Context) (*domain.LoginPolicy, error) {
	policyWriteModel := NewIAMLoginPolicyWriteModel(r.iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	policy := writeModelToLoginPolicy(policyWriteModel)
	policy.Default = true
	return policy, nil
}

func (r *CommandSide) AddDefaultLoginPolicy(ctx context.Context, policy *domain.LoginPolicy) (*domain.LoginPolicy, error) {
	policy.AggregateID = r.iamID
	addedPolicy := NewIAMLoginPolicyWriteModel(policy.AggregateID)
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.WriteModel)
	err := r.addDefaultLoginPolicy(ctx, nil, addedPolicy, policy)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLoginPolicy(addedPolicy), nil
}

func (r *CommandSide) addDefaultLoginPolicy(ctx context.Context, iamAgg *iam_repo.Aggregate, addedPolicy *IAMLoginPolicyWriteModel, policy *domain.LoginPolicy) error {
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return err
	}
	if addedPolicy.IsActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewLoginPolicyAddedEvent(ctx, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIdp, policy.ForceMFA, policy.PasswordlessType))

	return nil
}

func (r *CommandSide) ChangeDefaultLoginPolicy(ctx context.Context, policy *domain.LoginPolicy) (*domain.LoginPolicy, error) {
	policy.AggregateID = r.iamID
	existingPolicy := NewIAMLoginPolicyWriteModel(r.iamID)
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LoginPolicyWriteModel.WriteModel)
	err := r.changeDefaultLoginPolicy(ctx, iamAgg, existingPolicy, policy)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLoginPolicy(existingPolicy), nil
}

func (r *CommandSide) changeDefaultLoginPolicy(ctx context.Context, iamAgg *iam_repo.Aggregate, existingPolicy *IAMLoginPolicyWriteModel, policy *domain.LoginPolicy) error {
	policy.AggregateID = r.iamID
	err := r.defaultLoginPolicyWriteModelByID(ctx, existingPolicy)
	if err != nil {
		return err
	}
	if !existingPolicy.IsActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-M0sif", "Errors.IAM.LoginPolicy.NotFound")
	}
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIdp, policy.ForceMFA, domain.PasswordlessType(policy.PasswordlessType))
	if !hasChanged {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-5M9vdd", "Errors.IAM.LoginPolicy.NotChanged")
	}
	iamAgg.PushEvents(changedEvent)

	return nil
}

func (r *CommandSide) AddIDPProviderToDefaultLoginPolicy(ctx context.Context, idpProvider *domain.IDPProvider) (*domain.IDPProvider, error) {
	idpProvider.AggregateID = r.iamID
	idpModel := NewIAMIdentityProviderWriteModel(idpProvider.AggregateID, idpProvider.IDPConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return nil, err
	}
	if idpModel.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.IDP.AlreadyExists")
	}

	iamAgg := IAMAggregateFromWriteModel(&idpModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewIdentityProviderAddedEvent(ctx, idpProvider.IDPConfigID, domain.IdentityProviderType(idpProvider.Type)))

	if err = r.eventstore.PushAggregate(ctx, idpModel, iamAgg); err != nil {
		return nil, err
	}

	return writeModelToIDPProvider(idpModel), nil
}

func (r *CommandSide) RemoveIDPProviderFromDefaultLoginPolicy(ctx context.Context, idpProvider *iam_model.IDPProvider) error {
	idpProvider.AggregateID = r.iamID
	idpModel := NewIAMIdentityProviderWriteModel(idpProvider.AggregateID, idpProvider.IDPConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return err
	}
	if !idpModel.IsActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-39fjs", "Errors.IAM.LoginPolicy.IDP.NotExisting")
	}
	iamAgg := IAMAggregateFromWriteModel(&idpModel.IdentityProviderWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewIdentityProviderRemovedEvent(ctx, idpProvider.IDPConfigID))

	return r.eventstore.PushAggregate(ctx, idpModel, iamAgg)
}

func (r *CommandSide) AddSecondFactorToDefaultLoginPolicy(ctx context.Context, secondFactor iam_model.SecondFactorType) (iam_model.SecondFactorType, error) {
	secondFactorModel := NewIAMSecondFactorWriteModel(r.iamID)
	iamAgg := IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
	err := r.addSecondFactorToDefaultLoginPolicy(ctx, nil, secondFactorModel, secondFactor)
	if err != nil {
		return iam_model.SecondFactorTypeUnspecified, err
	}

	if err = r.eventstore.PushAggregate(ctx, secondFactorModel, iamAgg); err != nil {
		return iam_model.SecondFactorTypeUnspecified, err
	}

	return iam_model.SecondFactorType(secondFactorModel.MFAType), nil
}

func (r *CommandSide) addSecondFactorToDefaultLoginPolicy(ctx context.Context, iamAgg *iam_repo.Aggregate, secondFactorModel *IAMSecondFactorWriteModel, secondFactor iam_model.SecondFactorType) error {
	err := r.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return err
	}

	if secondFactorModel.IsActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewLoginPolicySecondFactorAddedEvent(ctx, domain.SecondFactorType(secondFactor)))

	return nil
}

func (r *CommandSide) RemoveSecondFactorFromDefaultLoginPolicy(ctx context.Context, secondFactor iam_model.SecondFactorType) error {
	secondFactorModel := NewIAMSecondFactorWriteModel(r.iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return err
	}
	if !secondFactorModel.IsActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-3M9od", "Errors.IAM.LoginPolicy.MFA.NotExisting")
	}
	iamAgg := IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewLoginPolicySecondFactorRemovedEvent(ctx, domain.SecondFactorType(secondFactor)))

	return r.eventstore.PushAggregate(ctx, secondFactorModel, iamAgg)
}

func (r *CommandSide) AddMultiFactorToDefaultLoginPolicy(ctx context.Context, multiFactor iam_model.MultiFactorType) (iam_model.MultiFactorType, error) {
	multiFactorModel := NewIAMMultiFactorWriteModel(r.iamID)
	iamAgg := IAMAggregateFromWriteModel(&multiFactorModel.MultiFactoryWriteModel.WriteModel)
	err := r.addMultiFactorToDefaultLoginPolicy(ctx, iamAgg, multiFactorModel, multiFactor)
	if err != nil {
		return iam_model.MultiFactorTypeUnspecified, err
	}

	if err = r.eventstore.PushAggregate(ctx, multiFactorModel, iamAgg); err != nil {
		return iam_model.MultiFactorTypeUnspecified, err
	}

	return iam_model.MultiFactorType(multiFactorModel.MultiFactoryWriteModel.MFAType), nil
}

func (r *CommandSide) addMultiFactorToDefaultLoginPolicy(ctx context.Context, iamAgg *iam_repo.Aggregate, multiFactorModel *IAMMultiFactorWriteModel, multiFactor iam_model.MultiFactorType) error {
	err := r.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return err
	}
	if multiFactorModel.IsActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-3M9od", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewLoginPolicyMultiFactorAddedEvent(ctx, domain.MultiFactorType(multiFactor)))

	return nil
}

func (r *CommandSide) RemoveMultiFactorFromDefaultLoginPolicy(ctx context.Context, multiFactor iam_model.MultiFactorType) error {
	multiFactorModel := NewIAMMultiFactorWriteModel(r.iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return err
	}
	if multiFactorModel.IsActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-3M9df", "Errors.IAM.LoginPolicy.MFA.NotExisting")
	}
	iamAgg := IAMAggregateFromWriteModel(&multiFactorModel.MultiFactoryWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewLoginPolicyMultiFactorRemovedEvent(ctx, domain.MultiFactorType(multiFactor)))

	return r.eventstore.PushAggregate(ctx, multiFactorModel, iamAgg)
}

func (r *CommandSide) defaultLoginPolicyWriteModelByID(ctx context.Context, writeModel *IAMLoginPolicyWriteModel) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	return nil
}
