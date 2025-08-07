package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) getOrgPrivacyPolicy(ctx context.Context, orgID string) (*domain.PrivacyPolicy, error) {
	policy, err := c.orgPrivacyPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if policy.State == domain.PolicyStateActive {
		return orgWriteModelToPrivacyPolicy(policy), nil
	}
	return c.getDefaultPrivacyPolicy(ctx)
}

func (c *Commands) orgPrivacyPolicyWriteModelByID(ctx context.Context, orgID string) (*OrgPrivacyPolicyWriteModel, error) {
	policy := NewOrgPrivacyPolicyWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

func (c *Commands) AddPrivacyPolicy(ctx context.Context, resourceOwner string, policy *domain.PrivacyPolicy) (_ *domain.PrivacyPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if policy.SupportEmail != "" {
		if err := policy.SupportEmail.Validate(); err != nil {
			return nil, err
		}
		policy.SupportEmail = policy.SupportEmail.Normalize()
	}

	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-MMk9fs", "Errors.ResourceOwnerMissing")
	}
	addedPolicy := NewOrgPrivacyPolicyWriteModel(resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, zerrors.ThrowAlreadyExists(nil, "Org-0oLpd", "Errors.Org.PrivacyPolicy.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.WriteModel)
	pushedEvents, err := c.eventstore.Push(
		ctx,
		org.NewPrivacyPolicyAddedEvent(
			ctx,
			orgAgg,
			policy.TOSLink,
			policy.PrivacyLink,
			policy.HelpLink,
			policy.SupportEmail,
			policy.DocsLink,
			policy.CustomLink,
			policy.CustomLinkText))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPrivacyPolicy(&addedPolicy.PrivacyPolicyWriteModel), nil
}

func (c *Commands) ChangePrivacyPolicy(ctx context.Context, resourceOwner string, policy *domain.PrivacyPolicy) (*domain.PrivacyPolicy, error) {

	if policy.SupportEmail != "" {
		if err := policy.SupportEmail.Validate(); err != nil {
			return nil, err
		}
		policy.SupportEmail = policy.SupportEmail.Normalize()
	}

	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-22N89f", "Errors.ResourceOwnerMissing")
	}

	existingPolicy, err := c.orgPrivacyPolicyWriteModelByID(ctx, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-Ng8sf", "Errors.Org.PrivacyPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.PrivacyPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, orgAgg, policy.TOSLink, policy.PrivacyLink, policy.HelpLink, policy.SupportEmail, policy.DocsLink, policy.CustomLink, policy.CustomLinkText)
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "Org-4N9fs", "Errors.Org.PrivacyPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPrivacyPolicy(&existingPolicy.PrivacyPolicyWriteModel), nil
}

func (c *Commands) RemovePrivacyPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-Nf9sf", "Errors.ResourceOwnerMissing")
	}
	existingPolicy := NewOrgPrivacyPolicyWriteModel(orgID)
	event, err := c.removePrivacyPolicy(ctx, existingPolicy)
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
	return writeModelToObjectDetails(&existingPolicy.PrivacyPolicyWriteModel.WriteModel), nil
}

func (c *Commands) removePrivacyPolicy(ctx context.Context, existingPolicy *OrgPrivacyPolicyWriteModel) (*org.PrivacyPolicyRemovedEvent, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-Ze9gs", "Errors.Org.PrivacyPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	return org.NewPrivacyPolicyRemovedEvent(ctx, orgAgg), nil
}

func (c *Commands) removePrivacyPolicyIfExists(ctx context.Context, orgID string) (*org.PrivacyPolicyRemovedEvent, error) {
	existingPolicy, err := c.orgPrivacyPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State != domain.PolicyStateActive {
		return nil, nil
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	return org.NewPrivacyPolicyRemovedEvent(ctx, orgAgg), nil
}
