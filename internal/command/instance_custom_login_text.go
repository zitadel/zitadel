package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

func (c *Commands) SetCustomInstanceLoginText(ctx context.Context, instanceID string, loginText *domain.CustomLoginText) (*domain.ObjectDetails, error) {
	iamAgg := instance.NewAggregate(instanceID)
	events, existingMailText, err := c.setCustomInstanceLoginText(ctx, &iamAgg.Aggregate, loginText)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMailText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingMailText.WriteModel), nil
}

func (c *Commands) RemoveCustomInstanceLoginTexts(ctx context.Context, lang language.Tag) (*domain.ObjectDetails, error) {
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
	iamAgg := InstanceAggregateFromWriteModel(&customText.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewCustomTextTemplateRemovedEvent(ctx, iamAgg, domain.LoginCustomText, lang))
	err = AppendAndReduce(customText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&customText.WriteModel), nil
}

func (c *Commands) setCustomInstanceLoginText(ctx context.Context, instanceAgg *eventstore.Aggregate, text *domain.CustomLoginText) ([]eventstore.Command, *InstanceCustomLoginTextReadModel, error) {
	if !text.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "Instance-kd9fs", "Errors.CustomText.Invalid")
	}

	existingLoginText, err := c.defaultLoginTextWriteModelByID(ctx, text.Language)
	if err != nil {
		return nil, nil, err
	}
	events := c.createAllLoginTextEvents(ctx, instanceAgg, &existingLoginText.CustomLoginTextReadModel, text, true)
	return events, existingLoginText, nil
}

func (c *Commands) defaultLoginTextWriteModelByID(ctx context.Context, lang language.Tag) (*InstanceCustomLoginTextReadModel, error) {
	writeModel := NewInstanceCustomLoginTextReadModel(lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
