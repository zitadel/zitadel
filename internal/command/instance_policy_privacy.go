package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) getDefaultPrivacyPolicy(ctx context.Context) (*domain.PrivacyPolicy, error) {
	policyWriteModel := NewInstancePrivacyPolicyWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	if !policyWriteModel.State.Exists() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-559os", "Errors.IAM.PasswordAgePolicy.NotFound")
	}
	policy := writeModelToPrivacyPolicy(&policyWriteModel.PrivacyPolicyWriteModel)
	policy.Default = true
	return policy, nil
}

func (c *Commands) AddDefaultPrivacyPolicy(ctx context.Context, policy *domain.PrivacyPolicy) (*domain.PrivacyPolicy, error) {
	addedPolicy := NewInstancePrivacyPolicyWriteModel(ctx)
	instanceAgg := InstanceAggregateFromWriteModel(&addedPolicy.WriteModel)
	events, err := c.addDefaultPrivacyPolicy(ctx, instanceAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, events)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPrivacyPolicy(&addedPolicy.PrivacyPolicyWriteModel), nil
}

func (c *Commands) addDefaultPrivacyPolicy(ctx context.Context, instanceAgg *eventstore.Aggregate, addedPolicy *InstancePrivacyPolicyWriteModel, policy *domain.PrivacyPolicy) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-M00rJ", "Errors.IAM.PrivacyPolicy.AlreadyExists")
	}

	return instance.NewPrivacyPolicyAddedEvent(ctx, instanceAgg, policy.TOSLink, policy.PrivacyLink, policy.HelpLink), nil
}

func (c *Commands) ChangeDefaultPrivacyPolicy(ctx context.Context, policy *domain.PrivacyPolicy) (*domain.PrivacyPolicy, error) {
	existingPolicy, err := c.defaultPrivacyPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-0oPew", "Errors.IAM.PasswordAgePolicy.NotFound")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.PrivacyPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, instanceAgg, policy.TOSLink, policy.PrivacyLink, policy.HelpLink)
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
	return writeModelToPrivacyPolicy(&existingPolicy.PrivacyPolicyWriteModel), nil
}

func (c *Commands) defaultPrivacyPolicyWriteModelByID(ctx context.Context) (policy *InstancePrivacyPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstancePrivacyPolicyWriteModel(ctx)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
