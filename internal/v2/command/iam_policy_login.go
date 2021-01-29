package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) getDefaultLoginPolicy(ctx context.Context) (*domain.LoginPolicy, error) {
	policyWriteModel := NewIAMLoginPolicyWriteModel()
	err := r.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	policy := writeModelToLoginPolicy(&policyWriteModel.LoginPolicyWriteModel)
	policy.Default = true
	return policy, nil
}

func (r *CommandSide) AddDefaultLoginPolicy(ctx context.Context, policy *domain.LoginPolicy) (*domain.LoginPolicy, error) {
	addedPolicy := NewIAMLoginPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.WriteModel)
	err := r.addDefaultLoginPolicy(ctx, nil, addedPolicy, policy)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLoginPolicy(&addedPolicy.LoginPolicyWriteModel), nil
}

func (r *CommandSide) addDefaultLoginPolicy(ctx context.Context, iamAgg *iam_repo.Aggregate, addedPolicy *IAMLoginPolicyWriteModel, policy *domain.LoginPolicy) error {
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewLoginPolicyAddedEvent(ctx, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIDP, policy.ForceMFA, policy.PasswordlessType))

	return nil
}

func (r *CommandSide) ChangeDefaultLoginPolicy(ctx context.Context, policy *domain.LoginPolicy) (*domain.LoginPolicy, error) {
	existingPolicy := NewIAMLoginPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LoginPolicyWriteModel.WriteModel)
	err := r.changeDefaultLoginPolicy(ctx, iamAgg, existingPolicy, policy)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLoginPolicy(&existingPolicy.LoginPolicyWriteModel), nil
}

func (r *CommandSide) changeDefaultLoginPolicy(ctx context.Context, iamAgg *iam_repo.Aggregate, existingPolicy *IAMLoginPolicyWriteModel, policy *domain.LoginPolicy) error {
	err := r.defaultLoginPolicyWriteModelByID(ctx, existingPolicy)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "IAM-M0sif", "Errors.IAM.LoginPolicy.NotFound")
	}
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIDP, policy.ForceMFA, policy.PasswordlessType)
	if !hasChanged {
		return caos_errs.ThrowPreconditionFailed(nil, "IAM-5M9vdd", "Errors.IAM.LoginPolicy.NotChanged")
	}
	iamAgg.PushEvents(changedEvent)

	return nil
}

func (r *CommandSide) AddIDPProviderToDefaultLoginPolicy(ctx context.Context, idpProvider *domain.IDPProvider) (*domain.IDPProvider, error) {
	idpModel := NewIAMIdentityProviderWriteModel(idpProvider.IDPConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return nil, err
	}
	if idpModel.State == domain.IdentityProviderStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.IDP.AlreadyExists")
	}

	iamAgg := IAMAggregateFromWriteModel(&idpModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewIdentityProviderAddedEvent(ctx, idpProvider.IDPConfigID, domain.IdentityProviderType(idpProvider.Type)))

	if err = r.eventstore.PushAggregate(ctx, idpModel, iamAgg); err != nil {
		return nil, err
	}

	return writeModelToIDPProvider(&idpModel.IdentityProviderWriteModel), nil
}

func (r *CommandSide) RemoveIDPProviderFromDefaultLoginPolicy(ctx context.Context, idpProvider *domain.IDPProvider) error {
	idpModel := NewIAMIdentityProviderWriteModel(idpProvider.IDPConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return err
	}
	if idpModel.State == domain.IdentityProviderStateUnspecified || idpModel.State == domain.IdentityProviderStateRemoved {
		return caos_errs.ThrowNotFound(nil, "IAM-39fjs", "Errors.IAM.LoginPolicy.IDP.NotExisting")
	}
	iamAgg := IAMAggregateFromWriteModel(&idpModel.IdentityProviderWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewIdentityProviderRemovedEvent(ctx, idpProvider.IDPConfigID))

	return r.eventstore.PushAggregate(ctx, idpModel, iamAgg)
}

func (r *CommandSide) AddSecondFactorToDefaultLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType) (domain.SecondFactorType, error) {
	secondFactorModel := NewIAMSecondFactorWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
	err := r.addSecondFactorToDefaultLoginPolicy(ctx, nil, secondFactorModel, secondFactor)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, err
	}

	if err = r.eventstore.PushAggregate(ctx, secondFactorModel, iamAgg); err != nil {
		return domain.SecondFactorTypeUnspecified, err
	}

	return domain.SecondFactorType(secondFactorModel.MFAType), nil
}

func (r *CommandSide) addSecondFactorToDefaultLoginPolicy(ctx context.Context, iamAgg *iam_repo.Aggregate, secondFactorModel *IAMSecondFactorWriteModel, secondFactor domain.SecondFactorType) error {
	err := r.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return err
	}

	if secondFactorModel.State == domain.FactorStateActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewLoginPolicySecondFactorAddedEvent(ctx, domain.SecondFactorType(secondFactor)))

	return nil
}

func (r *CommandSide) RemoveSecondFactorFromDefaultLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType) error {
	secondFactorModel := NewIAMSecondFactorWriteModel()
	err := r.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return err
	}
	if secondFactorModel.State == domain.FactorStateUnspecified || secondFactorModel.State == domain.FactorStateRemoved {
		return caos_errs.ThrowNotFound(nil, "IAM-3M9od", "Errors.IAM.LoginPolicy.MFA.NotExisting")
	}
	iamAgg := IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewLoginPolicySecondFactorRemovedEvent(ctx, domain.SecondFactorType(secondFactor)))

	return r.eventstore.PushAggregate(ctx, secondFactorModel, iamAgg)
}

func (r *CommandSide) AddMultiFactorToDefaultLoginPolicy(ctx context.Context, multiFactor domain.MultiFactorType) (domain.MultiFactorType, error) {
	multiFactorModel := NewIAMMultiFactorWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&multiFactorModel.MultiFactoryWriteModel.WriteModel)
	err := r.addMultiFactorToDefaultLoginPolicy(ctx, iamAgg, multiFactorModel, multiFactor)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, err
	}

	if err = r.eventstore.PushAggregate(ctx, multiFactorModel, iamAgg); err != nil {
		return domain.MultiFactorTypeUnspecified, err
	}

	return domain.MultiFactorType(multiFactorModel.MultiFactoryWriteModel.MFAType), nil
}

func (r *CommandSide) addMultiFactorToDefaultLoginPolicy(ctx context.Context, iamAgg *iam_repo.Aggregate, multiFactorModel *IAMMultiFactorWriteModel, multiFactor domain.MultiFactorType) error {
	err := r.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return err
	}
	if multiFactorModel.State == domain.FactorStateActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-3M9od", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewLoginPolicyMultiFactorAddedEvent(ctx, domain.MultiFactorType(multiFactor)))

	return nil
}

func (r *CommandSide) RemoveMultiFactorFromDefaultLoginPolicy(ctx context.Context, multiFactor domain.MultiFactorType) error {
	multiFactorModel := NewIAMMultiFactorWriteModel()
	err := r.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return err
	}
	if multiFactorModel.State == domain.FactorStateUnspecified || multiFactorModel.State == domain.FactorStateRemoved {
		return caos_errs.ThrowNotFound(nil, "IAM-3M9df", "Errors.IAM.LoginPolicy.MFA.NotExisting")
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
