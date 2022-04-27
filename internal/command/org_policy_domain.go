package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddOrgDomainPolicy(ctx context.Context, resourceOwner string, policy *domain.DomainPolicy) (*domain.DomainPolicy, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-4Jfsf", "Errors.ResourceOwnerMissing")
	}
	addedPolicy := NewOrgDomainPolicyWriteModel(resourceOwner)
	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.PolicyDomainWriteModel.WriteModel)
	event, err := c.addOrgDomainPolicy(ctx, orgAgg, addedPolicy, policy)
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
	return orgWriteModelToDomainPolicy(addedPolicy), nil
}

func (c *Commands) addOrgDomainPolicy(ctx context.Context, orgAgg *eventstore.Aggregate, addedPolicy *OrgDomainPolicyWriteModel, policy *domain.DomainPolicy) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "ORG-1M8ds", "Errors.Org.DomainPolicy.AlreadyExists")
	}
	return org.NewDomainPolicyAddedEvent(ctx, orgAgg, policy.UserLoginMustBeDomain, policy.ValidateOrgDomains), nil
}

func (c *Commands) ChangeOrgDomainPolicy(ctx context.Context, resourceOwner string, policy *domain.DomainPolicy) (*domain.DomainPolicy, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-5H8fs", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgDomainPolicyWriteModelByID(ctx, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-2N9sd", "Errors.Org.DomainPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.PolicyDomainWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, orgAgg, policy.UserLoginMustBeDomain, policy.ValidateOrgDomains)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-3M9ds", "Errors.Org.LabelPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return orgWriteModelToDomainPolicy(existingPolicy), nil
}

func (c *Commands) RemoveOrgDomainPolicy(ctx context.Context, orgID string) error {
	if orgID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "Org-3H8fs", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgDomainPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "ORG-Dvsh3", "Errors.Org.DomainPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.PolicyDomainWriteModel.WriteModel)
	_, err = c.eventstore.Push(ctx, org.NewDomainPolicyRemovedEvent(ctx, orgAgg))
	return err
}

func (c *Commands) getOrgDomainPolicy(ctx context.Context, orgID string) (*domain.DomainPolicy, error) {
	policy, err := c.orgDomainPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if policy.State == domain.PolicyStateActive {
		return orgWriteModelToDomainPolicy(policy), nil
	}
	return c.getDefaultDomainPolicy(ctx)
}

func (c *Commands) orgDomainPolicyWriteModelByID(ctx context.Context, orgID string) (policy *OrgDomainPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewOrgDomainPolicyWriteModel(orgID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
