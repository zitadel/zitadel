package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"golang.org/x/text/language"
)

func (c *Commands) SetDefaultMessageText(ctx context.Context, messageText *domain.CustomMessageText) (*domain.ObjectDetails, error) {
	iamAgg := iam.NewAggregate()
	events, existingMessageText, err := c.setDefaultMessageText(ctx, &iamAgg.Aggregate, messageText)
	if err != nil {
		return nil, err
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

func (c *Commands) setDefaultMessageText(ctx context.Context, iamAgg *eventstore.Aggregate, msg *domain.CustomMessageText) ([]eventstore.Command, *IAMCustomMessageTextReadModel, error) {
	if !msg.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "IAM-kd9fs", "Errors.CustomMessageText.Invalid")
	}

	existingMessageText, err := c.defaultCustomMessageTextWriteModelByID(ctx, msg.MessageTextType, msg.Language)
	if err != nil {
		return nil, nil, err
	}

	events := make([]eventstore.Command, 0)
	if existingMessageText.Greeting != msg.Greeting {
		if msg.Greeting != "" {
			events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageGreeting, msg.Greeting, msg.Language))
		} else {
			events = append(events, iam.NewCustomTextRemovedEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageGreeting, msg.Language))
		}
	}
	if existingMessageText.Subject != msg.Subject {
		if msg.Subject != "" {
			events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageSubject, msg.Subject, msg.Language))
		} else {
			events = append(events, iam.NewCustomTextRemovedEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageSubject, msg.Language))
		}
	}
	if existingMessageText.Title != msg.Title {
		if msg.Title != "" {
			events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageTitle, msg.Title, msg.Language))
		} else {
			events = append(events, iam.NewCustomTextRemovedEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageTitle, msg.Language))
		}
	}
	if existingMessageText.PreHeader != msg.PreHeader {
		if msg.PreHeader != "" {
			events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessagePreHeader, msg.PreHeader, msg.Language))
		} else {
			events = append(events, iam.NewCustomTextRemovedEvent(ctx, iamAgg, msg.MessageTextType, domain.MessagePreHeader, msg.Language))
		}
	}
	if existingMessageText.Text != msg.Text {
		if msg.Text != "" {
			events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageText, msg.Text, msg.Language))
		} else {
			events = append(events, iam.NewCustomTextRemovedEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageText, msg.Language))
		}
	}
	if existingMessageText.ButtonText != msg.ButtonText {
		if msg.ButtonText != "" {
			events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageButtonText, msg.ButtonText, msg.Language))
		} else {
			events = append(events, iam.NewCustomTextRemovedEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageButtonText, msg.Language))
		}
	}
	if existingMessageText.FooterText != msg.FooterText {
		if msg.FooterText != "" {
			events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageFooterText, msg.FooterText, msg.Language))
		} else {
			events = append(events, iam.NewCustomTextRemovedEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageFooterText, msg.Language))
		}
	}
	return events, existingMessageText, nil
}

func (c *Commands) RemoveIAMMessageTexts(ctx context.Context, messageTextType string, lang language.Tag) (*domain.ObjectDetails, error) {
	if messageTextType == "" || lang == language.Und {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-fjw9b", "Errors.CustomMessageText.Invalid")
	}
	customText, err := c.defaultCustomMessageTextWriteModelByID(ctx, messageTextType, lang)
	if err != nil {
		return nil, err
	}
	if customText.State == domain.PolicyStateUnspecified || customText.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-fju90", "Errors.CustomMessageText.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&customText.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewCustomTextTemplateRemovedEvent(ctx, iamAgg, messageTextType, lang))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(customText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&customText.WriteModel), nil
}

func (c *Commands) defaultCustomMessageTextWriteModelByID(ctx context.Context, messageType string, lang language.Tag) (*IAMCustomMessageTextReadModel, error) {
	writeModel := NewIAMCustomMessageTextWriteModel(messageType, lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
