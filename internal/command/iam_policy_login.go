package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
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
	event, err := r.addDefaultLoginPolicy(ctx, iamAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}
	_, err = r.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}

	return writeModelToLoginPolicy(&addedPolicy.LoginPolicyWriteModel), nil
}

func (r *CommandSide) addDefaultLoginPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMLoginPolicyWriteModel, policy *domain.LoginPolicy) (eventstore.EventPusher, error) {
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.AlreadyExists")
	}

	return iam_repo.NewLoginPolicyAddedEvent(ctx, iamAgg, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIDP, policy.ForceMFA, policy.PasswordlessType), nil
}

func (r *CommandSide) ChangeDefaultLoginPolicy(ctx context.Context, policy *domain.LoginPolicy) (*domain.LoginPolicy, error) {
	existingPolicy := NewIAMLoginPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LoginPolicyWriteModel.WriteModel)
	event, err := r.changeDefaultLoginPolicy(ctx, iamAgg, existingPolicy, policy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := r.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLoginPolicy(&existingPolicy.LoginPolicyWriteModel), nil
}

func (r *CommandSide) changeDefaultLoginPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, existingPolicy *IAMLoginPolicyWriteModel, policy *domain.LoginPolicy) (eventstore.EventPusher, error) {
	err := r.defaultLoginPolicyWriteModelByID(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-M0sif", "Errors.IAM.LoginPolicy.NotFound")
	}
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, iamAgg, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIDP, policy.ForceMFA, policy.PasswordlessType)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-5M9vdd", "Errors.IAM.LoginPolicy.NotChanged")
	}
	return changedEvent, nil
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
	pushedEvents, err := r.eventstore.PushEvents(ctx, iam_repo.NewIdentityProviderAddedEvent(ctx, iamAgg, idpProvider.IDPConfigID, idpProvider.Type))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(idpModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPProvider(&idpModel.IdentityProviderWriteModel), nil
}

func (r *CommandSide) RemoveIDPProviderFromDefaultLoginPolicy(ctx context.Context, idpProvider *domain.IDPProvider, cascadeExternalIDPs ...*domain.ExternalIDP) error {
	idpModel := NewIAMIdentityProviderWriteModel(idpProvider.IDPConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return err
	}
	if idpModel.State == domain.IdentityProviderStateUnspecified || idpModel.State == domain.IdentityProviderStateRemoved {
		return caos_errs.ThrowNotFound(nil, "IAM-39fjs", "Errors.IAM.LoginPolicy.IDP.NotExisting")
	}

	iamAgg := IAMAggregateFromWriteModel(&idpModel.IdentityProviderWriteModel.WriteModel)
	events := []eventstore.EventPusher{
		iam_repo.NewIdentityProviderRemovedEvent(ctx, iamAgg, idpProvider.IDPConfigID),
	}

	userEvents := r.removeIDPProviderFromDefaultLoginPolicy(ctx, iamAgg, idpProvider, false, cascadeExternalIDPs...)
	events = append(events, userEvents...)
	_, err = r.eventstore.PushEvents(ctx, events...)
	return err
}

func (r *CommandSide) removeIDPProviderFromDefaultLoginPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, idpProvider *domain.IDPProvider, cascade bool, cascadeExternalIDPs ...*domain.ExternalIDP) []eventstore.EventPusher {
	var events []eventstore.EventPusher
	if cascade {
		events = append(events, iam_repo.NewIdentityProviderCascadeRemovedEvent(ctx, iamAgg, idpProvider.IDPConfigID))
	} else {
		events = append(events, iam_repo.NewIdentityProviderRemovedEvent(ctx, iamAgg, idpProvider.IDPConfigID))
	}

	for _, idp := range cascadeExternalIDPs {
		userEvent, err := r.removeHumanExternalIDP(ctx, idp, true)
		if err != nil {
			logging.LogWithFields("COMMAND-4nfsf", "userid", idp.AggregateID, "idp-id", idp.IDPConfigID).WithError(err).Warn("could not cascade remove externalidp in remove provider from policy")
			continue
		}
		events = append(events, userEvent)
	}
	return events
}

func (r *CommandSide) AddSecondFactorToDefaultLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType) (domain.SecondFactorType, error) {
	secondFactorModel := NewIAMSecondFactorWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
	event, err := r.addSecondFactorToDefaultLoginPolicy(ctx, iamAgg, secondFactorModel, secondFactor)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, err
	}

	if _, err = r.eventstore.PushEvents(ctx, event); err != nil {
		return domain.SecondFactorTypeUnspecified, err
	}

	return secondFactorModel.MFAType, nil
}

func (r *CommandSide) addSecondFactorToDefaultLoginPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, secondFactorModel *IAMSecondFactorWriteModel, secondFactor domain.SecondFactorType) (eventstore.EventPusher, error) {
	err := r.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return nil, err
	}

	if secondFactorModel.State == domain.FactorStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}
	return iam_repo.NewLoginPolicySecondFactorAddedEvent(ctx, iamAgg, secondFactor), nil
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
	_, err = r.eventstore.PushEvents(ctx, iam_repo.NewLoginPolicySecondFactorRemovedEvent(ctx, iamAgg, secondFactor))
	return err
}

func (r *CommandSide) AddMultiFactorToDefaultLoginPolicy(ctx context.Context, multiFactor domain.MultiFactorType) (domain.MultiFactorType, error) {
	multiFactorModel := NewIAMMultiFactorWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&multiFactorModel.MultiFactoryWriteModel.WriteModel)
	event, err := r.addMultiFactorToDefaultLoginPolicy(ctx, iamAgg, multiFactorModel, multiFactor)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, err
	}

	if _, err = r.eventstore.PushEvents(ctx, event); err != nil {
		return domain.MultiFactorTypeUnspecified, err
	}

	return multiFactorModel.MultiFactoryWriteModel.MFAType, nil
}

func (r *CommandSide) addMultiFactorToDefaultLoginPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, multiFactorModel *IAMMultiFactorWriteModel, multiFactor domain.MultiFactorType) (eventstore.EventPusher, error) {
	err := r.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return nil, err
	}
	if multiFactorModel.State == domain.FactorStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-3M9od", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}

	return iam_repo.NewLoginPolicyMultiFactorAddedEvent(ctx, iamAgg, multiFactor), nil
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
	_, err = r.eventstore.PushEvents(ctx, iam_repo.NewLoginPolicyMultiFactorRemovedEvent(ctx, iamAgg, multiFactor))
	return err
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
