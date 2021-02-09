package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddDefaultMailText(ctx context.Context, policy *domain.MailText) (*domain.MailText, error) {
	addedPolicy := NewIAMMailTextWriteModel(policy.MailTextType, policy.Language)
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.MailTextWriteModel.WriteModel)
	err := r.addDefaultMailText(ctx, nil, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToMailTextPolicy(&addedPolicy.MailTextWriteModel), nil
}

func (r *CommandSide) addDefaultMailText(ctx context.Context, iamAgg *iam_repo.Aggregate, addedPolicy *IAMMailTextWriteModel, mailText *domain.MailText) error {
	if !mailText.IsValid() {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-3n8fs", "Errors.IAM.MailText.Invalid")
	}
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-9o0pM", "Errors.IAM.MailText.AlreadyExists")
	}

	iamAgg.PushEvents(
		iam_repo.NewMailTextAddedEvent(
			ctx,
			mailText.MailTextType,
			mailText.Language,
			mailText.Title,
			mailText.PreHeader,
			mailText.Subject,
			mailText.Greeting,
			mailText.Text,
			mailText.ButtonText),
	)

	return nil
}

func (r *CommandSide) ChangeDefaultMailText(ctx context.Context, mailText *domain.MailText) (*domain.MailText, error) {
	if !mailText.IsValid() {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-kd9fs", "Errors.IAM.MailText.Invalid")
	}
	existingPolicy, err := r.defaultMailTextWriteModelByID(ctx, mailText.MailTextType, mailText.Language)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-2N8fs", "Errors.IAM.MailText.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(
		ctx,
		mailText.MailTextType,
		mailText.Language,
		mailText.Title,
		mailText.PreHeader,
		mailText.Subject,
		mailText.Greeting,
		mailText.Text,
		mailText.ButtonText)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-m9L0s", "Errors.IAM.MailText.NotChanged")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.MailTextWriteModel.WriteModel)
	iamAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToMailTextPolicy(&existingPolicy.MailTextWriteModel), nil
}

func (r *CommandSide) defaultMailTextWriteModelByID(ctx context.Context, mailTextType, language string) (policy *IAMMailTextWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMMailTextWriteModel(mailTextType, language)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
