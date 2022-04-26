package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	iam_repo "github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) SetIAMCustomText(ctx context.Context, customText *domain.CustomText) (*domain.CustomText, error) {
	setText := NewIAMCustomTextWriteModel(customText.Key, customText.Language)
	iamAgg := IAMAggregateFromWriteModel(&setText.CustomTextWriteModel.WriteModel)
	event, err := c.setDefaultCustomText(ctx, iamAgg, setText, customText)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(setText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToCustomText(&setText.CustomTextWriteModel), nil
}

func (c *Commands) setDefaultCustomText(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMCustomTextWriteModel, text *domain.CustomText) (eventstore.Command, error) {
	if !text.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-3MN0s", "Errors.CustomText.Invalid")
	}
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	return iam_repo.NewCustomTextSetEvent(
		ctx,
		iamAgg,
		text.Template,
		text.Key,
		text.Text,
		text.Language), nil
}

func (c *Commands) defaultCustomTextWriteModelByID(ctx context.Context, key string, language language.Tag) (policy *IAMCustomTextWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMCustomTextWriteModel(key, language)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
