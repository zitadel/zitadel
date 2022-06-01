package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDefaultPrivacyPolicy(ctx context.Context, tosLink, privacyLink, helpLink string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddDefaultPrivacyPolicy(instanceAgg, tosLink, privacyLink, helpLink))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) ChangeDefaultPrivacyPolicy(ctx context.Context, policy *domain.PrivacyPolicy) (*domain.PrivacyPolicy, error) {
	existingPolicy, err := c.defaultPrivacyPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-0oPew", "Errors.IAM.PrivacyPolicy.NotFound")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.PrivacyPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, instanceAgg, policy.TOSLink, policy.PrivacyLink, policy.HelpLink)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-9jJfs", "Errors.IAM.PrivacyPolicy.NotChanged")
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

func (c *Commands) getDefaultPrivacyPolicy(ctx context.Context) (*domain.PrivacyPolicy, error) {
	policyWriteModel := NewInstancePrivacyPolicyWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	if !policyWriteModel.State.Exists() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-559os", "Errors.IAM.PrivacyPolicy.NotFound")
	}
	policy := writeModelToPrivacyPolicy(&policyWriteModel.PrivacyPolicyWriteModel)
	policy.Default = true
	return policy, nil
}

func prepareAddDefaultPrivacyPolicy(
	a *instance.Aggregate,
	tosLink,
	privacyLink,
	helpLink string,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstancePrivacyPolicyWriteModel(ctx)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if writeModel.State == domain.PolicyStateActive {
				return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-M00rJ", "Errors.Instance.PrivacyPolicy.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewPrivacyPolicyAddedEvent(ctx, &a.Aggregate, tosLink, privacyLink, helpLink),
			}, nil
		}, nil
	}
}
