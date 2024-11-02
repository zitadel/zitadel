package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// SetCustomInstanceLoginText only validates if the language is supported, not if it is allowed.
// This enables setting texts before allowing a language
func (c *Commands) SetCustomInstanceLoginText(ctx context.Context, loginText *domain.CustomLoginText) (*domain.ObjectDetails, error) {
	iamAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	events, existingMailText, err := c.setCustomInstanceLoginText(ctx, &iamAgg.Aggregate, loginText)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModelToObjectDetails(&existingMailText.WriteModel), nil
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
		return nil, zerrors.ThrowInvalidArgument(nil, "IAM-Gfbg3", "Errors.CustomText.Invalid")
	}
	customText, err := c.defaultLoginTextWriteModelByID(ctx, lang)
	if err != nil {
		return nil, err
	}
	if customText.State == domain.PolicyStateUnspecified || customText.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "IAM-fru44", "Errors.CustomText.NotFound")
	}
	iamAgg := InstanceAggregateFromWriteModel(&customText.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewCustomTextTemplateRemovedEvent(ctx, iamAgg, domain.LoginCustomText, lang))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(customText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&customText.WriteModel), nil
}

func (c *Commands) setCustomInstanceLoginText(ctx context.Context, instanceAgg *eventstore.Aggregate, text *domain.CustomLoginText) ([]eventstore.Command, *InstanceCustomLoginTextReadModel, error) {
	if err := text.IsValid(i18n.SupportedLanguages()); err != nil {
		return nil, nil, err
	}
	existingLoginText, err := c.defaultLoginTextWriteModelByID(ctx, text.Language)
	if err != nil {
		return nil, nil, err
	}
	events := c.createAllLoginTextEvents(ctx, instanceAgg, &existingLoginText.CustomLoginTextReadModel, text, true)
	return events, existingLoginText, nil
}

func (c *Commands) defaultLoginTextWriteModelByID(ctx context.Context, lang language.Tag) (*InstanceCustomLoginTextReadModel, error) {
	writeModel := NewInstanceCustomLoginTextReadModel(ctx, lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
