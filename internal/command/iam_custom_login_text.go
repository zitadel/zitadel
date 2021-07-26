package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

func (c *Commands) SetCustomIAMLoginText(ctx context.Context, loginText *domain.CustomLoginText) (*domain.ObjectDetails, error) {
	iamAgg := iam.NewAggregate()
	events, existingMailText, err := c.setCustomIAMLoginText(ctx, &iamAgg.Aggregate, loginText)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMailText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingMailText.WriteModel), nil
}

func (c *Commands) RemoveCustomIAMLoginTexts(ctx context.Context, lang language.Tag) (*domain.ObjectDetails, error) {
	if lang == language.Und {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-Gfbg3", "Errors.CustomText.Invalid")
	}
	customText, err := c.defaultLoginTextWriteModelByID(ctx, lang)
	if err != nil {
		return nil, err
	}
	if customText.State == domain.PolicyStateUnspecified || customText.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-fru44", "Errors.CustomText.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&customText.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam.NewCustomTextTemplateRemovedEvent(ctx, iamAgg, domain.LoginCustomText, lang))
	err = AppendAndReduce(customText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&customText.WriteModel), nil
}

func (c *Commands) setCustomIAMLoginText(ctx context.Context, iamAgg *eventstore.Aggregate, text *domain.CustomLoginText) ([]eventstore.EventPusher, *IAMCustomLoginTextReadModel, error) {
	if !text.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "IAM-kd9fs", "Errors.CustomText.Invalid")
	}

	existingLoginText, err := c.defaultLoginTextWriteModelByID(ctx, text.Language)
	if err != nil {
		return nil, nil, err
	}
	events := c.createAllLoginTextEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, true)
	return events, existingLoginText, nil
}

func (c *Commands) defaultLoginTextWriteModelByID(ctx context.Context, lang language.Tag) (*IAMCustomLoginTextReadModel, error) {
	writeModel := NewIAMCustomLoginTextReadModel(lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
