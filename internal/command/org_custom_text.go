package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) SetOrgCustomText(ctx context.Context, resourceOwner string, text *domain.CustomText) (*domain.CustomText, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-3n8fs", "Errors.ResourceOwnerMissing")
	}
	if !text.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-o93Fs", "Errors.CustomText.Invalid")
	}
	setText := NewOrgCustomTextWriteModel(resourceOwner, text.Key, text.Language)
	err := c.eventstore.FilterToQueryReducer(ctx, setText)
	if err != nil {
		return nil, err
	}

	orgAgg := OrgAggregateFromWriteModel(&setText.CustomTextWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(
		ctx,
		org.NewCustomTextSetEvent(
			ctx,
			orgAgg,
			text.Template,
			text.Key,
			text.Text,
			text.Language))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(setText, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToCustomText(&setText.CustomTextWriteModel), nil
}

func (c *Commands) RemoveOrgCustomText(ctx context.Context, resourceOwner, template, key string, lang language.Tag) error {
	if resourceOwner == "" {
		return caos_errs.ThrowInvalidArgument(nil, "Org-2N7fd", "Errors.ResourceOwnerMissing")
	}
	if template == "" || key == "" || lang == language.Und {
		return caos_errs.ThrowInvalidArgument(nil, "Org-3n9fsd", "Errors.CustomText.Invalid")
	}
	customText := NewOrgCustomTextWriteModel(resourceOwner, template, lang)
	err := c.eventstore.FilterToQueryReducer(ctx, customText)
	if err != nil {
		return err
	}
	if customText.State == domain.CustomTextStateUnspecified || customText.State == domain.CustomTextStateRemoved {
		return caos_errs.ThrowNotFound(nil, "Org-3n8fs", "Errors.CustomText.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&customText.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, org.NewCustomTextRemovedEvent(ctx, orgAgg, template, key, lang))
	return err
}
