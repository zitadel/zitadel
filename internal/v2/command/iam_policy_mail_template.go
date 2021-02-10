package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddDefaultMailTemplate(ctx context.Context, policy *domain.MailTemplate) (*domain.MailTemplate, error) {
	addedPolicy := NewIAMMailTemplateWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.MailTemplateWriteModel.WriteModel)
	err := r.addDefaultMailTemplate(ctx, nil, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToMailTemplatePolicy(&addedPolicy.MailTemplateWriteModel), nil
}

func (r *CommandSide) addDefaultMailTemplate(ctx context.Context, iamAgg *iam_repo.Aggregate, addedPolicy *IAMMailTemplateWriteModel, policy *domain.MailTemplate) error {
	if !policy.IsValid() {
		return caos_errs.ThrowPreconditionFailed(nil, "IAM-fm9sd", "Errors.IAM.MailTemplate.Invalid")
	}
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-5n8fs", "Errors.IAM.MailTemplate.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewMailTemplateAddedEvent(ctx, policy.Template))

	return nil
}

func (r *CommandSide) ChangeDefaultMailTemplate(ctx context.Context, policy *domain.MailTemplate) (*domain.MailTemplate, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-4m9ds", "Errors.IAM.MailTemplate.Invalid")
	}
	existingPolicy, err := r.defaultMailTemplateWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-2N8fs", "Errors.IAM.MailTemplate.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.Template)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-3nfsG", "Errors.IAM.MailTemplate.NotChanged")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.MailTemplateWriteModel.WriteModel)
	iamAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToMailTemplatePolicy(&existingPolicy.MailTemplateWriteModel), nil
}

func (r *CommandSide) defaultMailTemplateWriteModelByID(ctx context.Context) (policy *IAMMailTemplateWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMMailTemplateWriteModel()
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
