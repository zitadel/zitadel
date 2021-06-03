package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) SetOrgFeatures(ctx context.Context, resourceOwner string, features *domain.Features) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Features-G5tg", "Errors.ResourceOwnerMissing")
	}
	existingFeatures := NewOrgFeaturesWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, existingFeatures)
	if err != nil {
		return nil, err
	}
	setEvent, hasChanged := existingFeatures.NewSetEvent(
		ctx,
		OrgAggregateFromWriteModel(&existingFeatures.FeaturesWriteModel.WriteModel),
		features.TierName,
		features.TierDescription,
		features.State,
		features.StateDescription,
		features.AuditLogRetention,
		features.LoginPolicyFactors,
		features.LoginPolicyIDP,
		features.LoginPolicyPasswordless,
		features.LoginPolicyRegistration,
		features.LoginPolicyUsernameLogin,
		features.LoginPolicyPasswordReset,
		features.PasswordComplexityPolicy,
		features.LabelPolicy,
		features.CustomDomain,
	)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Features-GE4h2", "Errors.Features.NotChanged")
	}

	events, err := c.ensureOrgSettingsToFeatures(ctx, resourceOwner, features)
	if err != nil {
		return nil, err
	}
	events = append(events, setEvent)

	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingFeatures, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingFeatures.WriteModel), nil
}

func (c *Commands) RemoveOrgFeatures(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Features-G5tg", "Errors.ResourceOwnerMissing")
	}
	existingFeatures := NewOrgFeaturesWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, existingFeatures)
	if err != nil {
		return nil, err
	}
	if existingFeatures.State == domain.FeaturesStateUnspecified || existingFeatures.State == domain.FeaturesStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Features-Bg32G", "Errors.Features.NotFound")
	}
	removedEvent := org.NewFeaturesRemovedEvent(ctx, OrgAggregateFromWriteModel(&existingFeatures.FeaturesWriteModel.WriteModel))

	features, err := c.getDefaultFeatures(ctx)
	if err != nil {
		return nil, err
	}
	events, err := c.ensureOrgSettingsToFeatures(ctx, orgID, features)
	if err != nil {
		return nil, err
	}

	events = append(events, removedEvent)
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingFeatures, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingFeatures.WriteModel), nil
}

func (c *Commands) ensureOrgSettingsToFeatures(ctx context.Context, orgID string, features *domain.Features) ([]eventstore.EventPusher, error) {
	events, err := c.setAllowedLoginPolicy(ctx, orgID, features)
	if err != nil {
		return nil, err
	}
	if !features.PasswordComplexityPolicy {
		removePasswordComplexityEvent, err := c.removePasswordComplexityPolicyIfExists(ctx, orgID)
		if err != nil {
			return nil, err
		}
		if removePasswordComplexityEvent != nil {
			events = append(events, removePasswordComplexityEvent)
		}
	}
	if !features.LabelPolicy {
		removeLabelPolicyEvent, err := c.removeLabelPolicyIfExists(ctx, orgID)
		if err != nil {
			return nil, err
		}
		if removeLabelPolicyEvent != nil {
			events = append(events, removeLabelPolicyEvent)
		}
	}
	if !features.CustomDomain {
		removeCustomDomainsEvents, err := c.removeCustomDomains(ctx, orgID)
		if err != nil {
			return nil, err
		}
		if removeCustomDomainsEvents != nil {
			events = append(events, removeCustomDomainsEvents...)
		}
	}
	return events, nil
}

func (c *Commands) setAllowedLoginPolicy(ctx context.Context, orgID string, features *domain.Features) ([]eventstore.EventPusher, error) {
	events := make([]eventstore.EventPusher, 0)
	existingPolicy, err := c.orgLoginPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, nil
	}
	defaultPolicy, err := c.getDefaultLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	policy := *existingPolicy
	if !features.LoginPolicyFactors {
		if defaultPolicy.ForceMFA != existingPolicy.ForceMFA {
			policy.ForceMFA = defaultPolicy.ForceMFA
		}
		authFactorsEvents, err := c.setDefaultAuthFactorsInCustomLoginPolicy(ctx, orgID)
		if err != nil {
			return nil, err
		}
		events = append(events, authFactorsEvents...)
	}
	if !features.LoginPolicyIDP {
		if defaultPolicy.AllowExternalIDP != existingPolicy.AllowExternalIDP {
			policy.AllowExternalIDP = defaultPolicy.AllowExternalIDP
		}
		//TODO: handle idps
	}
	if !features.LoginPolicyRegistration && defaultPolicy.AllowRegister != existingPolicy.AllowRegister {
		policy.AllowRegister = defaultPolicy.AllowRegister
	}
	if !features.LoginPolicyPasswordless && defaultPolicy.PasswordlessType != existingPolicy.PasswordlessType {
		policy.PasswordlessType = defaultPolicy.PasswordlessType
	}
	if !features.LoginPolicyUsernameLogin && defaultPolicy.AllowUsernamePassword != existingPolicy.AllowUserNamePassword {
		policy.AllowUserNamePassword = defaultPolicy.AllowUsernamePassword
	}
	if !features.LoginPolicyPasswordReset && defaultPolicy.HidePasswordReset != existingPolicy.HidePasswordReset {
		policy.HidePasswordReset = defaultPolicy.HidePasswordReset
	}
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, OrgAggregateFromWriteModel(&existingPolicy.WriteModel), policy.AllowUserNamePassword, policy.AllowRegister, policy.AllowExternalIDP, policy.ForceMFA, policy.HidePasswordReset, policy.PasswordlessType)
	if hasChanged {
		events = append(events, changedEvent)
	}
	return events, nil
}

func (c *Commands) setDefaultAuthFactorsInCustomLoginPolicy(ctx context.Context, orgID string) ([]eventstore.EventPusher, error) {
	orgAuthFactors, err := c.orgLoginPolicyAuthFactorsWriteModel(ctx, orgID)
	if err != nil {
		return nil, err
	}
	events := make([]eventstore.EventPusher, 0)
	for _, factor := range domain.SecondFactorTypes() {
		state := orgAuthFactors.SecondFactors[factor]
		if state == nil || state.IAM == state.Org {
			continue
		}
		secondFactorWriteModel := orgAuthFactors.ToSecondFactorWriteModel(factor)
		if state.IAM == domain.FactorStateActive {
			event, err := c.addSecondFactorToLoginPolicy(ctx, secondFactorWriteModel, factor)
			if err != nil {
				return nil, err
			}
			if event != nil {
				events = append(events, event)
			}
			continue
		}
		event, err := c.removeSecondFactorFromLoginPolicy(ctx, secondFactorWriteModel, factor)
		if err != nil {
			return nil, err
		}
		if event != nil {
			events = append(events, event)
		}
	}

	for _, factor := range domain.MultiFactorTypes() {
		state := orgAuthFactors.MultiFactors[factor]
		if state == nil || state.IAM == state.Org {
			continue
		}
		multiFactorWriteModel := orgAuthFactors.ToMultiFactorWriteModel(factor)
		if state.IAM == domain.FactorStateActive {
			event, err := c.addMultiFactorToLoginPolicy(ctx, multiFactorWriteModel, factor)
			if err != nil {
				return nil, err
			}
			if event != nil {
				events = append(events, event)
			}
			continue
		}
		event, err := c.removeMultiFactorFromLoginPolicy(ctx, multiFactorWriteModel, factor)
		if err != nil {
			return nil, err
		}
		if event != nil {
			events = append(events, event)
		}
	}
	return events, nil
}
