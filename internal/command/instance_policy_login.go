package command

import (
	"context"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) ChangeDefaultLoginPolicy(ctx context.Context, policy *ChangeLoginPolicy) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareChangeDefaultLoginPolicy(instanceAgg, policy))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
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

	exists, err := ExistsInstanceIDP(ctx, c.eventstore.Filter, idpProvider.IDPConfigID)
	if err != nil || !exists {
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

func (c *Commands) AddSecondFactorToDefaultLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddSecondFactorToDefaultLoginPolicy(instanceAgg, secondFactor))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
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

func (c *Commands) AddMultiFactorToDefaultLoginPolicy(ctx context.Context, multiFactor domain.MultiFactorType) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddMultiFactorToDefaultLoginPolicy(instanceAgg, multiFactor))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
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

func prepareChangeDefaultLoginPolicy(a *instance.Aggregate, policy *ChangeLoginPolicy) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if ok := domain.ValidateDefaultRedirectURI(policy.DefaultRedirectURI); !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-SFdqd", "Errors.IAM.LoginPolicy.RedirectURIInvalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			wm := NewInstanceLoginPolicyWriteModel(ctx)
			if err := queryAndReduce(ctx, filter, wm); err != nil {
				return nil, err
			}
			if !wm.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-M0sif", "Errors.IAM.LoginPolicy.NotFound")
			}
			changedEvent, hasChanged := wm.NewChangedEvent(ctx, &a.Aggregate,
				policy.AllowUsernamePassword,
				policy.AllowRegister,
				policy.AllowExternalIDP,
				policy.ForceMFA,
				policy.HidePasswordReset,
				policy.IgnoreUnknownUsernames,
				policy.AllowDomainDiscovery,
				policy.DisableLoginWithEmail,
				policy.DisableLoginWithPhone,
				policy.PasswordlessType,
				policy.DefaultRedirectURI,
				policy.PasswordCheckLifetime,
				policy.ExternalLoginCheckLifetime,
				policy.MFAInitSkipLifetime,
				policy.SecondFactorCheckLifetime,
				policy.MultiFactorCheckLifetime)
			if !hasChanged {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-5M9vdd", "Errors.IAM.LoginPolicy.NotChanged")
			}
			return []eventstore.Command{changedEvent}, nil
		}, nil
	}
}

func prepareAddDefaultLoginPolicy(
	a *instance.Aggregate,
	allowUsernamePassword bool,
	allowRegister bool,
	allowExternalIDP bool,
	forceMFA bool,
	hidePasswordReset bool,
	ignoreUnknownUsernames bool,
	allowDomainDiscovery bool,
	disableLoginWithEmail bool,
	disableLoginWithPhone bool,
	passwordlessType domain.PasswordlessType,
	defaultRedirectURI string,
	passwordCheckLifetime time.Duration,
	externalLoginCheckLifetime time.Duration,
	mfaInitSkipLifetime time.Duration,
	secondFactorCheckLifetime time.Duration,
	multiFactorCheckLifetime time.Duration,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceLoginPolicyWriteModel(ctx)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if writeModel.State == domain.PolicyStateActive {
				return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-2B0ps", "Errors.Instance.LoginPolicy.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewLoginPolicyAddedEvent(ctx, &a.Aggregate,
					allowUsernamePassword,
					allowRegister,
					allowExternalIDP,
					forceMFA,
					hidePasswordReset,
					ignoreUnknownUsernames,
					allowDomainDiscovery,
					disableLoginWithEmail,
					disableLoginWithPhone,
					passwordlessType,
					defaultRedirectURI,
					passwordCheckLifetime,
					externalLoginCheckLifetime,
					mfaInitSkipLifetime,
					secondFactorCheckLifetime,
					multiFactorCheckLifetime,
				),
			}, nil
		}, nil
	}
}

func prepareAddSecondFactorToDefaultLoginPolicy(a *instance.Aggregate, factor domain.SecondFactorType) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if !factor.Valid() {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-5m9fs", "Errors.Instance.LoginPolicy.MFA.Unspecified")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceSecondFactorWriteModel(ctx, factor)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if writeModel.State == domain.FactorStateActive {
				return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-2B0ps", "Errors.Instance.MFA.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewLoginPolicySecondFactorAddedEvent(ctx, &a.Aggregate, factor),
			}, nil
		}, nil
	}
}

func prepareAddMultiFactorToDefaultLoginPolicy(a *instance.Aggregate, factor domain.MultiFactorType) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if !factor.Valid() {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-5m9fs", "Errors.Instance.LoginPolicy.MFA.Unspecified")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceMultiFactorWriteModel(ctx, factor)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if writeModel.State == domain.FactorStateActive {
				return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-3M9od", "Errors.Instance.MFA.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewLoginPolicyMultiFactorAddedEvent(ctx, &a.Aggregate, factor),
			}, nil
		}, nil
	}
}
