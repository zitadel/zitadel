package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddDefaultMailTemplate(ctx context.Context, policy *domain.MailTemplate) (*domain.MailTemplate, error) {
	addedPolicy := NewInstanceMailTemplateWriteModel(ctx)
	instanceAgg := InstanceAggregateFromWriteModel(&addedPolicy.MailTemplateWriteModel.WriteModel)
	event, err := c.addDefaultMailTemplate(ctx, instanceAgg, addedPolicy, policy)
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
	return writeModelToMailTemplatePolicy(&addedPolicy.MailTemplateWriteModel), nil
}

func (c *Commands) addDefaultMailTemplate(ctx context.Context, instanceAgg *eventstore.Aggregate, addedPolicy *InstanceMailTemplateWriteModel, policy *domain.MailTemplate) (eventstore.Command, error) {
	if !policy.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-fm9sd", "Errors.IAM.MailTemplate.Invalid")
	}
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, zerrors.ThrowAlreadyExists(nil, "INSTANCE-5n8fs", "Errors.IAM.MailTemplate.AlreadyExists")
	}

	return instance.NewMailTemplateAddedEvent(ctx, instanceAgg, policy.Template), nil
}

func (c *Commands) ChangeDefaultMailTemplate(ctx context.Context, policy *domain.MailTemplate) (*domain.MailTemplate, error) {
	existingPolicy, changedEvent, err := c.changeDefaultMailTemplate(ctx, policy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToMailTemplatePolicy(&existingPolicy.MailTemplateWriteModel), nil
}

func (c *Commands) changeDefaultMailTemplate(ctx context.Context, policy *domain.MailTemplate) (*InstanceMailTemplateWriteModel, eventstore.Command, error) {
	if !policy.IsValid() {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-4m9ds", "Errors.IAM.MailTemplate.Invalid")
	}
	existingPolicy, err := c.defaultMailTemplateWriteModelByID(ctx)
	if err != nil {
		return nil, nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, nil, zerrors.ThrowNotFound(nil, "INSTANCE-2N8fs", "Errors.IAM.MailTemplate.NotFound")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.MailTemplateWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, instanceAgg, policy.Template)
	if !hasChanged {
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "INSTANCE-3nfsG", "Errors.IAM.MailTemplate.NotChanged")
	}

	return existingPolicy, changedEvent, nil
}

func (c *Commands) defaultMailTemplateWriteModelByID(ctx context.Context) (policy *InstanceMailTemplateWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceMailTemplateWriteModel(ctx)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func prepareAddDefaultEmailTemplate(
	a *instance.Aggregate,
	template []byte,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if template == nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-fm9sd", "Errors.Instance.MailTemplate.Invalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceMailTemplateWriteModel(ctx)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if writeModel.State == domain.PolicyStateActive {
				return nil, zerrors.ThrowAlreadyExists(nil, "INSTANCE-5n8fs", "Errors.Instance.MailTemplate.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewMailTemplateAddedEvent(ctx, &a.Aggregate,
					template,
				),
			}, nil
		}, nil
	}
}
