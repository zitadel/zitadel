package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddDefaultPrivacyPolicy(ctx context.Context, tosLink, privacyLink, helpLink string, supportEmail domain.EmailAddress, docsLink, customLink, customLinkText string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())

	if supportEmail != "" {
		if err := supportEmail.Validate(); err != nil {
			return nil, err
		}
		supportEmail = supportEmail.Normalize()
	}

	writeModel := NewInstancePrivacyPolicyWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	if writeModel.State.Exists() {
		return nil, zerrors.ThrowAlreadyExists(nil, "INSTANCE-M00rJ", "Errors.Instance.PrivacyPolicy.AlreadyExists")
	}

	event := instance.NewPrivacyPolicyAddedEvent(ctx, &instanceAgg.Aggregate, tosLink, privacyLink, helpLink, supportEmail, docsLink, customLink, customLinkText)

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(writeModel, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) ChangeDefaultPrivacyPolicy(ctx context.Context, policy *domain.PrivacyPolicy) (*domain.PrivacyPolicy, error) {
	if policy.SupportEmail != "" {
		if err := policy.SupportEmail.Validate(); err != nil {
			return nil, err
		}
		policy.SupportEmail = policy.SupportEmail.Normalize()
	}

	existingPolicy, err := c.defaultPrivacyPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-0oPew", "Errors.IAM.PrivacyPolicy.NotFound")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.PrivacyPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, instanceAgg, policy.TOSLink, policy.PrivacyLink, policy.HelpLink, policy.SupportEmail, policy.DocsLink, policy.CustomLink, policy.CustomLinkText)
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "INSTANCE-9jJfs", "Errors.IAM.PrivacyPolicy.NotChanged")
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
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-559os", "Errors.IAM.PrivacyPolicy.NotFound")
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
	supportEmail domain.EmailAddress,
	docsLink, customLink, customLinkText string,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if supportEmail != "" {
			if err := supportEmail.Validate(); err != nil {
				return nil, err
			}
			supportEmail = supportEmail.Normalize()
		}
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
				return nil, zerrors.ThrowAlreadyExists(nil, "INSTANCE-M00rJ", "Errors.Instance.PrivacyPolicy.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewPrivacyPolicyAddedEvent(ctx, &a.Aggregate, tosLink, privacyLink, helpLink, supportEmail, docsLink, customLink, customLinkText),
			}, nil
		}, nil
	}
}
