package command

import (
	"context"
	"reflect"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) AddLoginPolicy(ctx context.Context, resourceOwner string, policy *domain.LoginPolicy) (*domain.LoginPolicy, error) {
	addedPolicy := NewOrgLoginPolicyWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-Dgfb2", "Errors.Org.LoginPolicy.AlreadyExists")
	}

	err = c.checkLoginPolicyAllowed(ctx, resourceOwner, nil)
	if err != nil {
		return nil, err
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(
		ctx,
		org.NewLoginPolicyAddedEvent(
			ctx,
			orgAgg,
			policy.AllowUsernamePassword,
			policy.AllowRegister,
			policy.AllowExternalIDP,
			policy.ForceMFA,
			policy.PasswordlessType))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLoginPolicy(&addedPolicy.LoginPolicyWriteModel), nil
}

func (c *Commands) ChangeLoginPolicy(ctx context.Context, resourceOwner string, policy *domain.LoginPolicy) (*domain.LoginPolicy, error) {
	existingPolicy := NewOrgLoginPolicyWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-M0sif", "Errors.Org.LoginPolicy.NotFound")
	}

	err = c.checkLoginPolicyAllowed(ctx, resourceOwner, policy)
	if err != nil {
		return nil, err
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LoginPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, orgAgg, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIDP, policy.ForceMFA, policy.PasswordlessType)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-5M9vdd", "Errors.Org.LoginPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLoginPolicy(&existingPolicy.LoginPolicyWriteModel), nil
}

func (c *Commands) checkLoginPolicyAllowed(ctx context.Context, resourceOwner string, policy *domain.LoginPolicy) error {
	defaultPolicy, err := c.getDefaultLoginPolicy(ctx)
	if err != nil {
		return err
	}
	requiredFeatures := make([]string, 0)
	if defaultPolicy.ForceMFA != policy.ForceMFA || !reflect.DeepEqual(defaultPolicy.MultiFactors, policy.MultiFactors) || !reflect.DeepEqual(defaultPolicy.SecondFactors, policy.SecondFactors) {
		requiredFeatures = append(requiredFeatures, domain.FeatureLoginPolicyFactors)
	}
	if defaultPolicy.AllowExternalIDP != policy.AllowExternalIDP || !reflect.DeepEqual(defaultPolicy.IDPProviders, policy.IDPProviders) {
		requiredFeatures = append(requiredFeatures, domain.FeatureLoginPolicyIDP)
	}
	if defaultPolicy.AllowRegister != policy.AllowRegister {
		requiredFeatures = append(requiredFeatures, domain.FeatureLoginPolicyRegistration)
	}
	if defaultPolicy.PasswordlessType != policy.PasswordlessType {
		requiredFeatures = append(requiredFeatures, domain.FeatureLoginPolicyPasswordless)
	}
	if defaultPolicy.AllowUsernamePassword != policy.AllowUsernamePassword {
		requiredFeatures = append(requiredFeatures, domain.FeatureLoginPolicyUsernameLogin)
	}
	return authz.CheckOrgFeatures(ctx, c.tokenVerifier, resourceOwner, requiredFeatures...)
}

func (c *Commands) RemoveLoginPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	existingPolicy := NewOrgLoginPolicyWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-GHB37", "Errors.Org.LoginPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LoginPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLoginPolicyRemovedEvent(ctx, orgAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LoginPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddIDPProviderToLoginPolicy(ctx context.Context, resourceOwner string, idpProvider *domain.IDPProvider) (*domain.IDPProvider, error) {
	idpModel := NewOrgIdentityProviderWriteModel(resourceOwner, idpProvider.IDPConfigID)
	err := c.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return nil, err
	}
	if idpModel.State == domain.IdentityProviderStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-2B0ps", "Errors.Org.LoginPolicy.IDP.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&idpModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewIdentityProviderAddedEvent(ctx, orgAgg, idpProvider.IDPConfigID, idpProvider.Type))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(idpModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPProvider(&idpModel.IdentityProviderWriteModel), nil
}

func (c *Commands) RemoveIDPProviderFromLoginPolicy(ctx context.Context, resourceOwner string, idpProvider *domain.IDPProvider, cascadeExternalIDPs ...*domain.ExternalIDP) (*domain.ObjectDetails, error) {
	idpModel := NewOrgIdentityProviderWriteModel(resourceOwner, idpProvider.IDPConfigID)
	err := c.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return nil, err
	}
	if idpModel.State == domain.IdentityProviderStateUnspecified || idpModel.State == domain.IdentityProviderStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-39fjs", "Errors.Org.LoginPolicy.IDP.NotExisting")
	}

	orgAgg := OrgAggregateFromWriteModel(&idpModel.IdentityProviderWriteModel.WriteModel)
	events := c.removeIDPProviderFromLoginPolicy(ctx, orgAgg, idpProvider.IDPConfigID, false, cascadeExternalIDPs...)

	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(idpModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&idpModel.WriteModel), nil
}

func (c *Commands) removeIDPProviderFromLoginPolicy(ctx context.Context, orgAgg *eventstore.Aggregate, idpConfigID string, cascade bool, cascadeExternalIDPs ...*domain.ExternalIDP) []eventstore.EventPusher {
	var events []eventstore.EventPusher
	if cascade {
		events = append(events, org.NewIdentityProviderCascadeRemovedEvent(ctx, orgAgg, idpConfigID))
	} else {
		events = append(events, org.NewIdentityProviderRemovedEvent(ctx, orgAgg, idpConfigID))
	}

	for _, idp := range cascadeExternalIDPs {
		event, _, err := c.removeHumanExternalIDP(ctx, idp, true)
		if err != nil {
			logging.LogWithFields("COMMAND-n8RRf", "userid", idp.AggregateID, "idpconfigid", idp.IDPConfigID).WithError(err).Warn("could not cascade remove external idp")
			continue
		}
		events = append(events, event)
	}
	return events
}

func (c *Commands) AddSecondFactorToLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType, orgID string) (domain.SecondFactorType, error) {
	secondFactorModel := NewOrgSecondFactorWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, err
	}

	if secondFactorModel.State == domain.FactorStateActive {
		return domain.SecondFactorTypeUnspecified, caos_errs.ThrowAlreadyExists(nil, "Org-2B0ps", "Errors.Org.LoginPolicy.MFA.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)

	if _, err = c.eventstore.PushEvents(ctx, org.NewLoginPolicySecondFactorAddedEvent(ctx, orgAgg, secondFactor)); err != nil {
		return domain.SecondFactorTypeUnspecified, err
	}

	return secondFactorModel.MFAType, nil
}

func (c *Commands) RemoveSecondFactorFromLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType, orgID string) error {
	secondFactorModel := NewOrgSecondFactorWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return err
	}
	if secondFactorModel.State == domain.FactorStateUnspecified || secondFactorModel.State == domain.FactorStateRemoved {
		return caos_errs.ThrowNotFound(nil, "Org-3M9od", "Errors.Org.LoginPolicy.MFA.NotExisting")
	}
	orgAgg := OrgAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)

	_, err = c.eventstore.PushEvents(ctx, org.NewLoginPolicySecondFactorRemovedEvent(ctx, orgAgg, secondFactor))
	return err
}

func (c *Commands) AddMultiFactorToLoginPolicy(ctx context.Context, multiFactor domain.MultiFactorType, orgID string) (domain.MultiFactorType, error) {
	multiFactorModel := NewOrgMultiFactorWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, err
	}
	if multiFactorModel.State == domain.FactorStateActive {
		return domain.MultiFactorTypeUnspecified, caos_errs.ThrowAlreadyExists(nil, "Org-3M9od", "Errors.Org.LoginPolicy.MFA.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&multiFactorModel.WriteModel)

	if _, err = c.eventstore.PushEvents(ctx, org.NewLoginPolicyMultiFactorAddedEvent(ctx, orgAgg, multiFactor)); err != nil {
		return domain.MultiFactorTypeUnspecified, err
	}

	return multiFactorModel.MFAType, nil
}

func (c *Commands) RemoveMultiFactorFromLoginPolicy(ctx context.Context, multiFactor domain.MultiFactorType, orgID string) error {
	multiFactorModel := NewOrgMultiFactorWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return err
	}
	if multiFactorModel.State == domain.FactorStateUnspecified || multiFactorModel.State == domain.FactorStateRemoved {
		return caos_errs.ThrowNotFound(nil, "Org-3M9df", "Errors.Org.LoginPolicy.MFA.NotExisting")
	}
	orgAgg := OrgAggregateFromWriteModel(&multiFactorModel.MultiFactoryWriteModel.WriteModel)

	_, err = c.eventstore.PushEvents(ctx, org.NewLoginPolicyMultiFactorRemovedEvent(ctx, orgAgg, multiFactor))
	return err
}
