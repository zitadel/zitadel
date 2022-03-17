package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDefaultOrgIAMPolicy(ctx context.Context, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	addedPolicy := NewInstanceOrgIAMPolicyWriteModel()
	instanceAgg := InstanceAggregateFromWriteModel(&addedPolicy.WriteModel)
	event, err := c.addDefaultOrgIAMPolicy(ctx, instanceAgg, addedPolicy, policy)
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
	return writeModelToOrgIAMPolicy(addedPolicy), nil
}

func (c *Commands) addDefaultOrgIAMPolicy(ctx context.Context, instanceAgg *eventstore.Aggregate, addedPolicy *InstanceOrgIAMPolicyWriteModel, policy *domain.OrgIAMPolicy) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-Lk0dS", "Errors.IAM.OrgIAMPolicy.AlreadyExists")
	}
	return iam_repo.NewOrgIAMPolicyAddedEvent(ctx, instanceAgg, policy.UserLoginMustBeDomain), nil
}

func (c *Commands) ChangeDefaultOrgIAMPolicy(ctx context.Context, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	existingPolicy, err := c.defaultOrgIAMPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if !existingPolicy.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-0Pl0d", "Errors.IAM.OrgIAMPolicy.NotFound")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.PolicyOrgIAMWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, instanceAgg, policy.UserLoginMustBeDomain)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-3n8fs", "Errors.IAM.PasswordComplexityPolicy.NotFound")
	}
	policy := writeModelToOrgIAMPolicy(policyWriteModel)
	policy.Default = true
	return policy, nil
}

func (c *Commands) defaultOrgIAMPolicyWriteModelByID(ctx context.Context) (policy *InstanceOrgIAMPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceOrgIAMPolicyWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
