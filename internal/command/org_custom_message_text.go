package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// SetOrgMessageText only validates if the language is supported, not if it is allowed.
// This enables setting texts before allowing a language
func (c *Commands) SetOrgMessageText(ctx context.Context, resourceOwner string, messageText *domain.CustomMessageText) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-2biiR", "Errors.ResourceOwnerMissing")
	}
	orgAgg := org.NewAggregate(resourceOwner)
	events, existingMessageText, err := c.setOrgMessageText(ctx, &orgAgg.Aggregate, messageText)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModelToObjectDetails(&existingMessageText.WriteModel), nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMessageText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingMessageText.WriteModel), nil
}

func (c *Commands) setOrgMessageText(ctx context.Context, orgAgg *eventstore.Aggregate, message *domain.CustomMessageText) ([]eventstore.Command, *OrgCustomMessageTextReadModel, error) {
	if err := message.IsValid(i18n.SupportedLanguages()); err != nil {
		return nil, nil, err
	}
	existingMessageText, err := c.orgCustomMessageTextWriteModelByID(ctx, orgAgg.ID, message.MessageTextType, message.Language)
	if err != nil {
		return nil, nil, err
	}
	events := make([]eventstore.Command, 0)
	if existingMessageText.Greeting != message.Greeting {
		if message.Greeting != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, message.MessageTextType, domain.MessageGreeting, message.Greeting, message.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, message.MessageTextType, domain.MessageGreeting, message.Language))
		}
	}
	if existingMessageText.Subject != message.Subject {
		if message.Subject != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, message.MessageTextType, domain.MessageSubject, message.Subject, message.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, message.MessageTextType, domain.MessageSubject, message.Language))
		}
	}
	if existingMessageText.Title != message.Title {
		if message.Title != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, message.MessageTextType, domain.MessageTitle, message.Title, message.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, message.MessageTextType, domain.MessageTitle, message.Language))
		}
	}
	if existingMessageText.PreHeader != message.PreHeader {
		if message.PreHeader != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, message.MessageTextType, domain.MessagePreHeader, message.PreHeader, message.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, message.MessageTextType, domain.MessagePreHeader, message.Language))
		}
	}
	if existingMessageText.Text != message.Text {
		if message.Text != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, message.MessageTextType, domain.MessageText, message.Text, message.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, message.MessageTextType, domain.MessageText, message.Language))
		}
	}
	if existingMessageText.ButtonText != message.ButtonText {
		if message.ButtonText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, message.MessageTextType, domain.MessageButtonText, message.ButtonText, message.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, message.MessageTextType, domain.MessageButtonText, message.Language))
		}
	}
	if existingMessageText.FooterText != message.FooterText {
		if message.FooterText != "" {
			events = append(events, org.NewCustomTextSetEvent(ctx, orgAgg, message.MessageTextType, domain.MessageFooterText, message.FooterText, message.Language))
		} else {
			events = append(events, org.NewCustomTextRemovedEvent(ctx, orgAgg, message.MessageTextType, domain.MessageFooterText, message.Language))
		}
	}
	return events, existingMessageText, nil
}

func (c *Commands) RemoveOrgMessageTexts(ctx context.Context, resourceOwner, messageTextType string, lang language.Tag) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-3mfsf", "Errors.ResourceOwnerMissing")
	}
	if messageTextType == "" || lang == language.Und {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-j59f", "Errors.CustomMessageText.Invalid")
	}
	customText, err := c.orgCustomMessageTextWriteModelByID(ctx, resourceOwner, messageTextType, lang)
	if err != nil {
		return nil, err
	}
	if customText.State == domain.PolicyStateUnspecified || customText.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "Org-3b8Jf", "Errors.CustomMessageText.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&customText.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewCustomTextTemplateRemovedEvent(ctx, orgAgg, messageTextType, lang))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(customText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&customText.WriteModel), nil
}

func (c *Commands) removeOrgMessageTextsIfExists(ctx context.Context, orgID string) ([]eventstore.Command, error) {
	msgTemplates := NewOrgCustomMessageTextsWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, msgTemplates)
	if err != nil {
		return nil, err
	}

	orgAgg := OrgAggregateFromWriteModel(&msgTemplates.WriteModel)
	events := make([]eventstore.Command, 0, len(msgTemplates.CustomMessageTemplate))
	for _, tmpl := range msgTemplates.CustomMessageTemplate {
		events = append(events, org.NewCustomTextTemplateRemovedEvent(ctx, orgAgg, tmpl.Template, tmpl.Language))
	}
	return events, nil
}

func (c *Commands) orgCustomMessageTextWriteModelByID(ctx context.Context, orgID, messageType string, lang language.Tag) (*OrgCustomMessageTextReadModel, error) {
	writeModel := NewOrgCustomMessageTextWriteModel(orgID, messageType, lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
