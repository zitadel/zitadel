package command

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) getDefaultLoginPolicy(ctx context.Context) (*domain.LoginPolicy, error) {
	policyWriteModel := NewInstanceLoginPolicyWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	policy := writeModelToLoginPolicy(&policyWriteModel.LoginPolicyWriteModel)
	policy.Default = true
	return policy, nil
}

func (c *Commands) AddDefaultLoginPolicy(ctx context.Context, policy *domain.LoginPolicy) (*domain.LoginPolicy, error) {
	addedPolicy := NewInstanceLoginPolicyWriteModel(ctx)
	instanceAgg := InstanceAggregateFromWriteModel(&addedPolicy.WriteModel)
	event, err := c.addDefaultLoginPolicy(ctx, instanceAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLoginPolicy(&addedPolicy.LoginPolicyWriteModel), nil
}

func (c *Commands) addDefaultLoginPolicy(ctx context.Context, instanceAgg *eventstore.Aggregate, addedPolicy *InstanceLoginPolicyWriteModel, policy *domain.LoginPolicy) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-2B0ps", "Errors.IAM.LoginPolicy.AlreadyExists")
	}

	return instance.NewLoginPolicyAddedEvent(ctx,
		instanceAgg,
		policy.AllowUsernamePassword,
		policy.AllowRegister,
		policy.AllowExternalIDP,
		policy.ForceMFA,
		policy.HidePasswordReset,
		policy.PasswordlessType,
		policy.PasswordCheckLifetime,
		policy.ExternalLoginCheckLifetime,
		policy.MFAInitSkipLifetime,
		policy.SecondFactorCheckLifetime,
		policy.MultiFactorCheckLifetime), nil
}

func (c *Commands) ChangeDefaultLoginPolicy(ctx context.Context, policy *domain.LoginPolicy) (*domain.LoginPolicy, error) {
	existingPolicy := NewInstanceLoginPolicyWriteModel(ctx)
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LoginPolicyWriteModel.WriteModel)
	event, err := c.changeDefaultLoginPolicy(ctx, instanceAgg, existingPolicy, policy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLoginPolicy(&existingPolicy.LoginPolicyWriteModel), nil
}

func (c *Commands) changeDefaultLoginPolicy(ctx context.Context, instanceAgg *eventstore.Aggregate, existingPolicy *InstanceLoginPolicyWriteModel, policy *domain.LoginPolicy) (eventstore.Command, error) {
	err := c.defaultLoginPolicyWriteModelByID(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-M0sif", "Errors.IAM.LoginPolicy.NotFound")
	}
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx,
		instanceAgg,
		policy.AllowUsernamePassword,
		policy.AllowRegister,
		policy.AllowExternalIDP,
		policy.ForceMFA,
		policy.HidePasswordReset,
		policy.PasswordlessType,
		policy.PasswordCheckLifetime,
		policy.ExternalLoginCheckLifetime,
		policy.MFAInitSkipLifetime,
		policy.SecondFactorCheckLifetime,
		policy.MultiFactorCheckLifetime)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-5M9vdd", "Errors.IAM.LoginPolicy.NotChanged")
	}
	return changedEvent, nil
}

func (c *Commands) AddIDPProviderToDefaultLoginPolicy(ctx context.Context, idpProvider *domain.IDPProvider) (*domain.IDPProvider, error) {
	if !idpProvider.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-9nf88", "Errors.IAM.LoginPolicy.IDP.Invalid")
	}
	existingPolicy := NewInstanceLoginPolicyWriteModel(ctx)
	err := c.defaultLoginPolicyWriteModelByID(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-GVDfe", "Errors.IAM.LoginPolicy.NotFound")
	}

	_, err = c.getInstanceIDPConfigByID(ctx, idpProvider.IDPConfigID)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "INSTANCE-m8fsd", "Errors.IDPConfig.NotExisting")
	}
	idpModel := NewInstanceIdentityProviderWriteModel(ctx, idpProvider.IDPConfigID)
	err = c.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return nil, err
	}
	if idpModel.State == domain.IdentityProviderStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-2B0ps", "Errors.IAM.LoginPolicy.IDP.AlreadyExists")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&idpModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewIdentityProviderAddedEvent(ctx, instanceAgg, idpProvider.IDPConfigID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(idpModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPProvider(&idpModel.IdentityProviderWriteModel), nil
}

func (c *Commands) RemoveIDPProviderFromDefaultLoginPolicy(ctx context.Context, idpProvider *domain.IDPProvider, cascadeExternalIDPs ...*domain.UserIDPLink) (*domain.ObjectDetails, error) {
	if !idpProvider.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-66m9s", "Errors.IAM.LoginPolicy.IDP.Invalid")
	}
	existingPolicy := NewInstanceLoginPolicyWriteModel(ctx)
	err := c.defaultLoginPolicyWriteModelByID(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-Dfg4t", "Errors.IAM.LoginPolicy.NotFound")
	}

	idpModel := NewInstanceIdentityProviderWriteModel(ctx, idpProvider.IDPConfigID)
	err = c.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return nil, err
	}
	if idpModel.State == domain.IdentityProviderStateUnspecified || idpModel.State == domain.IdentityProviderStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-39fjs", "Errors.IAM.LoginPolicy.IDP.NotExisting")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&idpModel.IdentityProviderWriteModel.WriteModel)
	events := c.removeIDPProviderFromDefaultLoginPolicy(ctx, instanceAgg, idpProvider, false, cascadeExternalIDPs...)
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(idpModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&idpModel.IdentityProviderWriteModel.WriteModel), nil
}

func (c *Commands) removeIDPProviderFromDefaultLoginPolicy(ctx context.Context, instanceAgg *eventstore.Aggregate, idpProvider *domain.IDPProvider, cascade bool, cascadeExternalIDPs ...*domain.UserIDPLink) []eventstore.Command {
	var events []eventstore.Command
	if cascade {
		events = append(events, instance.NewIdentityProviderCascadeRemovedEvent(ctx, instanceAgg, idpProvider.IDPConfigID))
	} else {
		events = append(events, instance.NewIdentityProviderRemovedEvent(ctx, instanceAgg, idpProvider.IDPConfigID))
	}

	for _, idp := range cascadeExternalIDPs {
		userEvent, _, err := c.removeUserIDPLink(ctx, idp, true)
		if err != nil {
			logging.WithFields("COMMAND-4nfsf", "userid", idp.AggregateID, "idp-id", idp.IDPConfigID).WithError(err).Warn("could not cascade remove externalidp in remove provider from policy")
			continue
		}
		events = append(events, userEvent)
	}
	return events
}

func (c *Commands) AddSecondFactorToDefaultLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType) (domain.SecondFactorType, *domain.ObjectDetails, error) {
	if !secondFactor.Valid() {
		return domain.SecondFactorTypeUnspecified, nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-5m9fs", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	secondFactorModel := NewInstanceSecondFactorWriteModel(ctx, secondFactor)
	instanceAgg := InstanceAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
	event, err := c.addSecondFactorToDefaultLoginPolicy(ctx, instanceAgg, secondFactorModel, secondFactor)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, nil, err
	}
	err = AppendAndReduce(secondFactorModel, pushedEvents...)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, nil, err
	}
	return secondFactorModel.MFAType, writeModelToObjectDetails(&secondFactorModel.WriteModel), nil
}

func (c *Commands) addSecondFactorToDefaultLoginPolicy(ctx context.Context, instanceAgg *eventstore.Aggregate, secondFactorModel *InstanceSecondFactorWriteModel, secondFactor domain.SecondFactorType) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return nil, err
	}

	if secondFactorModel.State == domain.FactorStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-2B0ps", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}
	return instance.NewLoginPolicySecondFactorAddedEvent(ctx, instanceAgg, secondFactor), nil
}

func (c *Commands) RemoveSecondFactorFromDefaultLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType) (*domain.ObjectDetails, error) {
	if !secondFactor.Valid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-55n8s", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	secondFactorModel := NewInstanceSecondFactorWriteModel(ctx, secondFactor)
	err := c.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return nil, err
	}
	if secondFactorModel.State == domain.FactorStateUnspecified || secondFactorModel.State == domain.FactorStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-3M9od", "Errors.IAM.LoginPolicy.MFA.NotExisting")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLoginPolicySecondFactorRemovedEvent(ctx, instanceAgg, secondFactor))
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
		return domain.MultiFactorTypeUnspecified, nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-5m9fs", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	multiFactorModel := NewInstanceMultiFactorWriteModel(ctx, multiFactor)
	instanceAgg := InstanceAggregateFromWriteModel(&multiFactorModel.MultiFactorWriteModel.WriteModel)
	event, err := c.addMultiFactorToDefaultLoginPolicy(ctx, instanceAgg, multiFactorModel, multiFactor)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, nil, err
	}
	err = AppendAndReduce(multiFactorModel, pushedEvents...)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, nil, err
	}
	return multiFactorModel.MultiFactorWriteModel.MFAType, writeModelToObjectDetails(&multiFactorModel.WriteModel), nil
}

func (c *Commands) addMultiFactorToDefaultLoginPolicy(ctx context.Context, instanceAgg *eventstore.Aggregate, multiFactorModel *InstanceMultiFactorWriteModel, multiFactor domain.MultiFactorType) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return nil, err
	}
	if multiFactorModel.State == domain.FactorStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-3M9od", "Errors.IAM.LoginPolicy.MFA.AlreadyExists")
	}

	return instance.NewLoginPolicyMultiFactorAddedEvent(ctx, instanceAgg, multiFactor), nil
}

func (c *Commands) RemoveMultiFactorFromDefaultLoginPolicy(ctx context.Context, multiFactor domain.MultiFactorType) (*domain.ObjectDetails, error) {
	if !multiFactor.Valid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-33m9F", "Errors.IAM.LoginPolicy.MFA.Unspecified")
	}
	multiFactorModel := NewInstanceMultiFactorWriteModel(ctx, multiFactor)
	err := c.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return nil, err
	}
	if multiFactorModel.State == domain.FactorStateUnspecified || multiFactorModel.State == domain.FactorStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-3M9df", "Errors.IAM.LoginPolicy.MFA.NotExisting")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&multiFactorModel.MultiFactorWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLoginPolicyMultiFactorRemovedEvent(ctx, instanceAgg, multiFactor))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(multiFactorModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&multiFactorModel.WriteModel), nil
}

func (c *Commands) defaultLoginPolicyWriteModelByID(ctx context.Context, writeModel *InstanceLoginPolicyWriteModel) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	return nil
}
