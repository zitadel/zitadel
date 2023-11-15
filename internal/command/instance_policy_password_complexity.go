package command

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/command/preparation"
	"github.com/zitadel/zitadel/v2/internal/domain"
	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/repository/instance"
	"github.com/zitadel/zitadel/v2/internal/telemetry/tracing"
)

func (c *Commands) AddDefaultPasswordComplexityPolicy(ctx context.Context, minLength uint64, hasLowercase, hasUppercase, hasNumber, hasSymbol bool) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddDefaultPasswordComplexityPolicy(instanceAgg, minLength, hasLowercase, hasUppercase, hasNumber, hasSymbol))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) ChangeDefaultPasswordComplexityPolicy(ctx context.Context, policy *domain.PasswordComplexityPolicy) (*domain.PasswordComplexityPolicy, error) {
	if err := policy.IsValid(); err != nil {
		return nil, err
	}

	existingPolicy, err := c.defaultPasswordComplexityPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-0oPew", "Errors.IAM.PasswordComplexityPolicy.NotFound")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.PasswordComplexityPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, instanceAgg, policy.MinLength, policy.HasLowercase, policy.HasUppercase, policy.HasNumber, policy.HasSymbol)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-9jlsf", "Errors.IAM.PasswordComplexityPolicy.NotChanged")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPasswordComplexityPolicy(&existingPolicy.PasswordComplexityPolicyWriteModel), nil
}

func prepareAddDefaultPasswordComplexityPolicy(
	a *instance.Aggregate,
	minLength uint64,
	hasLowercase,
	hasUppercase,
	hasNumber,
	hasSymbol bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if minLength == 0 || minLength > 72 {
			return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-Lsp0e", "Errors.Instance.PasswordComplexityPolicy.MinLengthNotAllowed")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstancePasswordComplexityPolicyWriteModel(ctx)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if writeModel.State == domain.PolicyStateActive {
				return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-Lk0dS", "Errors.Instance.PasswordComplexityPolicy.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewPasswordComplexityPolicyAddedEvent(ctx, &a.Aggregate,
					minLength,
					hasLowercase,
					hasUppercase,
					hasNumber,
					hasSymbol,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) getDefaultPasswordComplexityPolicy(ctx context.Context) (*domain.PasswordComplexityPolicy, error) {
	policyWriteModel := NewInstancePasswordComplexityPolicyWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	if !policyWriteModel.State.Exists() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-M0gsf", "Errors.IAM.PasswordComplexityPolicy.NotFound")
	}
	policy := writeModelToPasswordComplexityPolicy(&policyWriteModel.PasswordComplexityPolicyWriteModel)
	policy.Default = true
	return policy, nil
}

func (c *Commands) defaultPasswordComplexityPolicyWriteModelByID(ctx context.Context) (policy *InstancePasswordComplexityPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstancePasswordComplexityPolicyWriteModel(ctx)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
