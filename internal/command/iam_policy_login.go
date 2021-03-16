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

func (c *Commands) getDefaultLoginPolicy(ctx context.Context) (*domain.LoginPolicy, error) {
	policyWriteModel := NewIAMLoginPolicyWriteModel()
	err := c.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	policy := writeModelToLoginPolicy(&policyWriteModel.LoginPolicyWriteModel)
	policy.Default = true
	return policy, nil
}

func (c *Commands) AddDefaultLoginPolicy(ctx context.Context, policy *domain.LoginPolicy) (*domain.LoginPolicy, error) {
	addedPolicy := NewIAMLoginPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.WriteModel)
	event, err := c.addDefaultLoginPolicy(ctx, iamAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLoginPolicy(&addedPolicy.LoginPolicyWriteModel), nil
}

func (c *Commands) addDefaultLoginPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMLoginPolicyWriteModel, policy *domain.LoginPolicy) (eventstore.EventPusher, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.AlreadyExists")
	}

	return iam_repo.NewLoginPolicyAddedEvent(ctx, iamAgg, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIDP, policy.ForceMFA, policy.PasswordlessType), nil
}

func (c *Commands) ChangeDefaultLoginPolicy(ctx context.Context, policy *domain.LoginPolicy) (*domain.LoginPolicy, error) {
	existingPolicy := NewIAMLoginPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LoginPolicyWriteModel.WriteModel)
	event, err := c.changeDefaultLoginPolicy(ctx, iamAgg, existingPolicy, policy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLoginPolicy(&existingPolicy.LoginPolicyWriteModel), nil
}

func (c *Commands) changeDefaultLoginPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, existingPolicy *IAMLoginPolicyWriteModel, policy *domain.LoginPolicy) (eventstore.EventPusher, error) {
	err := c.defaultLoginPolicyWriteModelByID(ctx, existingPolicy)
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

func (c *Commands) AddIDPProviderToDefaultLoginPolicy(ctx context.Context, idpProvider *domain.IDPProvider) (*domain.IDPProvider, error) {
	if !idpProvider.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-9nf88", "Errors.IAM.LoginPolicy.IDP.")
	}
	_, err := c.getIAMIDPConfigByID(ctx, idpProvider.IDPConfigID)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "IAM-m8fsd", "Errors.IDPConfig.NotExisting")
	}
	idpModel := NewIAMIdentityProviderWriteModel(idpProvider.IDPConfigID)
	err = c.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return nil, err
	}
	if idpModel.State == domain.IdentityProviderStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.IDP.AlreadyExists")
	}

	iamAgg := IAMAggregateFromWriteModel(&idpModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewIdentityProviderAddedEvent(ctx, iamAgg, idpProvider.IDPConfigID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(idpModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPProvider(&idpModel.IdentityProviderWriteModel), nil
}

func (c *Commands) RemoveIDPProviderFromDefaultLoginPolicy(ctx context.Context, idpProvider *domain.IDPProvider, cascadeExternalIDPs ...*domain.ExternalIDP) (*domain.ObjectDetails, error) {
	if !idpProvider.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-66m9s", "Errors.IAM.LoginPolicy.IDP.Invalid")
	}
	idpModel := NewIAMIdentityProviderWriteModel(idpProvider.IDPConfigID)
	err := c.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return nil, err
	}
	if idpModel.State == domain.IdentityProviderStateUnspecified || idpModel.State == domain.IdentityProviderStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-39fjs", "Errors.IAM.LoginPolicy.IDP.NotExisting")
	}

	iamAgg := IAMAggregateFromWriteModel(&idpModel.IdentityProviderWriteModel.WriteModel)
	events := c.removeIDPProviderFromDefaultLoginPolicy(ctx, iamAgg, idpProvider, false, cascadeExternalIDPs...)
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(idpModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&idpModel.IdentityProviderWriteModel.WriteModel), nil
}

func (c *Commands) removeIDPProviderFromDefaultLoginPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, idpProvider *domain.IDPProvider, cascade bool, cascadeExternalIDPs ...*domain.ExternalIDP) []eventstore.EventPusher {
	var events []eventstore.EventPusher
	if cascade {
		events = append(events, iam_repo.NewIdentityProviderCascadeRemovedEvent(ctx, iamAgg, idpProvider.IDPConfigID))
	} else {
		events = append(events, iam_repo.NewIdentityProviderRemovedEvent(ctx, iamAgg, idpProvider.IDPConfigID))
	}

	for _, idp := range cascadeExternalIDPs {
		userEvent, _, err := c.removeHumanExternalIDP(ctx, idp, true)
		if err != nil {
			logging.LogWithFields("COMMAND-4nfsf", "userid", idp.AggregateID, "idp-id", idp.IDPConfigID).WithError(err).Warn("could not cascade remove externalidp in remove provider from policy")
			continue
		}
		events = append(events, userEvent)
	}
	return events
}

func (c *Commands) AddSecondFactorToDefaultLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType) (domain.SecondFactorType, *domain.ObjectDetails, error) {
	if !secondFactor.Valid() {
		return domain.SecondFactorTypeUnspecified, nil, caos_errs.ThrowInvalidArgument(nil, "IAM-5m9fs", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	secondFactorModel := NewIAMSecondFactorWriteModel(secondFactor)
	iamAgg := IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
	event, err := c.addSecondFactorToDefaultLoginPolicy(ctx, iamAgg, secondFactorModel, secondFactor)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, nil, err
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, nil, err
	}
	err = AppendAndReduce(secondFactorModel, pushedEvents...)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, nil, err
	}
	return secondFactorModel.MFAType, writeModelToObjectDetails(&secondFactorModel.WriteModel), nil
}

func (c *Commands) addSecondFactorToDefaultLoginPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, secondFactorModel *IAMSecondFactorWriteModel, secondFactor domain.SecondFactorType) (eventstore.EventPusher, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return nil, err
	}

	if secondFactorModel.State == domain.FactorStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}
	return iam_repo.NewLoginPolicySecondFactorAddedEvent(ctx, iamAgg, secondFactor), nil
}

func (c *Commands) RemoveSecondFactorFromDefaultLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType) (*domain.ObjectDetails, error) {
	if !secondFactor.Valid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-55n8s", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	secondFactorModel := NewIAMSecondFactorWriteModel(secondFactor)
	err := c.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return nil, err
	}
	if secondFactorModel.State == domain.FactorStateUnspecified || secondFactorModel.State == domain.FactorStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-3M9od", "Errors.IAM.LoginPolicy.MFA.NotExisting")
	}
	iamAgg := IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLoginPolicySecondFactorRemovedEvent(ctx, iamAgg, secondFactor))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(secondFactorModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&secondFactorModel.WriteModel), nil
}

func (c *Commands) AddMultiFactorToDefaultLoginPolicy(ctx context.Context, multiFactor domain.MultiFactorType) (domain.MultiFactorType, *domain.ObjectDetails, error) {
	if !multiFactor.Valid() {
		return domain.MultiFactorTypeUnspecified, nil, caos_errs.ThrowInvalidArgument(nil, "IAM-5m9fs", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	multiFactorModel := NewIAMMultiFactorWriteModel(multiFactor)
	iamAgg := IAMAggregateFromWriteModel(&multiFactorModel.MultiFactoryWriteModel.WriteModel)
	event, err := c.addMultiFactorToDefaultLoginPolicy(ctx, iamAgg, multiFactorModel, multiFactor)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, nil, err
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, nil, err
	}
	err = AppendAndReduce(multiFactorModel, pushedEvents...)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, nil, err
	}
	return multiFactorModel.MultiFactoryWriteModel.MFAType, writeModelToObjectDetails(&multiFactorModel.WriteModel), nil
}

func (c *Commands) addMultiFactorToDefaultLoginPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, multiFactorModel *IAMMultiFactorWriteModel, multiFactor domain.MultiFactorType) (eventstore.EventPusher, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return nil, err
	}
	if multiFactorModel.State == domain.FactorStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-3M9od", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}

	return iam_repo.NewLoginPolicyMultiFactorAddedEvent(ctx, iamAgg, multiFactor), nil
}

func (c *Commands) RemoveMultiFactorFromDefaultLoginPolicy(ctx context.Context, multiFactor domain.MultiFactorType) (*domain.ObjectDetails, error) {
	if !multiFactor.Valid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-33m9F", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	multiFactorModel := NewIAMMultiFactorWriteModel(multiFactor)
	err := c.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return nil, err
	}
	if multiFactorModel.State == domain.FactorStateUnspecified || multiFactorModel.State == domain.FactorStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-3M9df", "Errors.IAM.LoginPolicy.MFA.NotExisting")
	}
	iamAgg := IAMAggregateFromWriteModel(&multiFactorModel.MultiFactoryWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLoginPolicyMultiFactorRemovedEvent(ctx, iamAgg, multiFactor))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(multiFactorModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&multiFactorModel.WriteModel), nil
}

func (c *Commands) defaultLoginPolicyWriteModelByID(ctx context.Context, writeModel *IAMLoginPolicyWriteModel) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	return nil
}
