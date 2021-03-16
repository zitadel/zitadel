package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDefaultOrgIAMPolicy(ctx context.Context, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	addedPolicy := NewIAMOrgIAMPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.WriteModel)
	event, err := c.addDefaultOrgIAMPolicy(ctx, iamAgg, addedPolicy, policy)
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
	return writeModelToOrgIAMPolicy(addedPolicy), nil
}

func (c *Commands) addDefaultOrgIAMPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMOrgIAMPolicyWriteModel, policy *domain.OrgIAMPolicy) (eventstore.EventPusher, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-Lk0dS", "Errors.IAM.OrgIAMPolicy.AlreadyExists")
	}
	return iam_repo.NewOrgIAMPolicyAddedEvent(ctx, iamAgg, policy.UserLoginMustBeDomain), nil
}

func (c *Commands) ChangeDefaultOrgIAMPolicy(ctx context.Context, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	existingPolicy, err := c.defaultOrgIAMPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if !existingPolicy.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-0Pl0d", "Errors.IAM.OrgIAMPolicy.NotFound")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.PolicyOrgIAMWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, iamAgg, policy.UserLoginMustBeDomain)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToOrgIAMPolicy(existingPolicy), nil
}

func (c *Commands) getDefaultOrgIAMPolicy(ctx context.Context) (*domain.OrgIAMPolicy, error) {
	policyWriteModel, err := c.defaultOrgIAMPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if !policyWriteModel.State.Exists() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-3n8fs", "Errors.IAM.PasswordComplexityPolicy.NotFound")
	}
	policy := writeModelToOrgIAMPolicy(policyWriteModel)
	policy.Default = true
	return policy, nil
}

func (c *Commands) defaultOrgIAMPolicyWriteModelByID(ctx context.Context) (policy *IAMOrgIAMPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMOrgIAMPolicyWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
