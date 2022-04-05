package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/telemetry/tracing"
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-fm9sd", "Errors.IAM.MailTemplate.Invalid")
	}
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-5n8fs", "Errors.IAM.MailTemplate.AlreadyExists")
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
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-4m9ds", "Errors.IAM.MailTemplate.Invalid")
	}
	existingPolicy, err := c.defaultMailTemplateWriteModelByID(ctx)
	if err != nil {
		return nil, nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, nil, caos_errs.ThrowNotFound(nil, "INSTANCE-2N8fs", "Errors.IAM.MailTemplate.NotFound")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.MailTemplateWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, instanceAgg, policy.Template)
	if !hasChanged {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-3nfsG", "Errors.IAM.MailTemplate.NotChanged")
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
