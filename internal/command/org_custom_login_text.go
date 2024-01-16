package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// SetOrgLoginText only validates if the language is supported, not if it is allowed.
// This enables setting texts before allowing a language
func (c *Commands) SetOrgLoginText(ctx context.Context, resourceOwner string, loginText *domain.CustomLoginText) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-m29rF", "Errors.ResourceOwnerMissing")
	}
	iamAgg := org.NewAggregate(resourceOwner)
	events, existingLoginText, err := c.setOrgLoginText(ctx, &iamAgg.Aggregate, loginText)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingLoginText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingLoginText.WriteModel), nil
}

func (c *Commands) setOrgLoginText(ctx context.Context, orgAgg *eventstore.Aggregate, loginText *domain.CustomLoginText) ([]eventstore.Command, *OrgCustomLoginTextReadModel, error) {
	if err := loginText.IsValid(i18n.SupportedLanguages()); err != nil {
		return nil, nil, err
	}
	existingLoginText, err := c.orgCustomLoginTextWriteModelByID(ctx, orgAgg.ID, loginText.Language)
	if err != nil {
		return nil, nil, err
	}
	events := c.createAllLoginTextEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)
	return events, existingLoginText, nil
}

func (c *Commands) RemoveOrgLoginTexts(ctx context.Context, resourceOwner string, lang language.Tag) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-1B8dw", "Errors.ResourceOwnerMissing")
	}
	if lang == language.Und {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-5ZZmo", "Errors.CustomText.Invalid")
	}
	customText, err := c.orgCustomLoginTextWriteModelByID(ctx, resourceOwner, lang)
	if err != nil {
		return nil, err
	}
	if customText.State == domain.PolicyStateUnspecified || customText.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "Org-9ru44", "Errors.CustomText.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&customText.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewCustomTextTemplateRemovedEvent(ctx, orgAgg, domain.LoginCustomText, lang))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(customText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&customText.WriteModel), nil
}

func (c *Commands) removeOrgLoginTextsIfExists(ctx context.Context, orgID string) ([]eventstore.Command, error) {
	msgTemplates := NewOrgCustomLoginTextsReadModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, msgTemplates)
	if err != nil {
		return nil, err
	}

	orgAgg := OrgAggregateFromWriteModel(&msgTemplates.WriteModel)
	events := make([]eventstore.Command, 0, len(msgTemplates.CustomLoginTexts))
	for _, tmpl := range msgTemplates.CustomLoginTexts {
		events = append(events, org.NewCustomTextTemplateRemovedEvent(ctx, orgAgg, tmpl.Template, tmpl.Language))
	}
	return events, nil
}

func (c *Commands) orgCustomLoginTextWriteModelByID(ctx context.Context, orgID string, lang language.Tag) (*OrgCustomLoginTextReadModel, error) {
	writeModel := NewOrgCustomLoginTextReadModel(orgID, lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
