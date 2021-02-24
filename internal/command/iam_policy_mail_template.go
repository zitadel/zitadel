package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDefaultMailTemplate(ctx context.Context, policy *domain.MailTemplate) (*domain.MailTemplate, error) {
	addedPolicy := NewIAMMailTemplateWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.MailTemplateWriteModel.WriteModel)
	event, err := c.addDefaultMailTemplate(ctx, iamAgg, addedPolicy, policy)
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
	return writeModelToMailTemplatePolicy(&addedPolicy.MailTemplateWriteModel), nil
}

func (c *Commands) addDefaultMailTemplate(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMMailTemplateWriteModel, policy *domain.MailTemplate) (eventstore.EventPusher, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-fm9sd", "Errors.IAM.MailTemplate.Invalid")
	}
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-5n8fs", "Errors.IAM.MailTemplate.AlreadyExists")
	}

	return iam_repo.NewMailTemplateAddedEvent(ctx, iamAgg, policy.Template), nil
}

func (c *Commands) ChangeDefaultMailTemplate(ctx context.Context, policy *domain.MailTemplate) (*domain.MailTemplate, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4m9ds", "Errors.IAM.MailTemplate.Invalid")
	}
	existingPolicy, err := c.defaultMailTemplateWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-2N8fs", "Errors.IAM.MailTemplate.NotFound")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.MailTemplateWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, iamAgg, policy.Template)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-3nfsG", "Errors.IAM.MailTemplate.NotChanged")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToMailTemplatePolicy(&existingPolicy.MailTemplateWriteModel), nil
}

func (c *Commands) defaultMailTemplateWriteModelByID(ctx context.Context) (policy *IAMMailTemplateWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMMailTemplateWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
