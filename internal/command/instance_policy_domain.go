package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDefaultDomainPolicy(ctx context.Context, policy *domain.DomainPolicy) (*domain.DomainPolicy, error) {
	addedPolicy := NewInstanceDomainPolicyWriteModel()
	instanceAgg := InstanceAggregateFromWriteModel(&addedPolicy.WriteModel)
	event, err := c.addDefaultDomainPolicy(ctx, instanceAgg, addedPolicy, policy)
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
	return writeModelToDomainPolicy(addedPolicy), nil
}

func (c *Commands) addDefaultDomainPolicy(ctx context.Context, instanceAgg *eventstore.Aggregate, addedPolicy *InstanceDomainPolicyWriteModel, policy *domain.DomainPolicy) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-Lk0dS", "Errors.IAM.DomainPolicy.AlreadyExists")
	}
	return iam_repo.NewInstnaceDomainPolicyAddedEvent(ctx, instanceAgg, policy.UserLoginMustBeDomain), nil
}

func (c *Commands) ChangeDefaultDomainPolicy(ctx context.Context, policy *domain.DomainPolicy) (*domain.DomainPolicy, error) {
	existingPolicy, err := c.defaultDomainPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if !existingPolicy.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-0Pl0d", "Errors.IAM.DomainPolicy.NotFound")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.PolicyDomainWriteModel.WriteModel)
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
	return writeModelToDomainPolicy(existingPolicy), nil
}

func (c *Commands) getDefaultDomainPolicy(ctx context.Context) (*domain.DomainPolicy, error) {
	policyWriteModel, err := c.defaultDomainPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if !policyWriteModel.State.Exists() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-3n8fs", "Errors.IAM.PasswordComplexityPolicy.NotFound")
	}
	policy := writeModelToDomainPolicy(policyWriteModel)
	policy.Default = true
	return policy, nil
}

func (c *Commands) defaultDomainPolicyWriteModelByID(ctx context.Context) (policy *InstanceDomainPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceDomainPolicyWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
