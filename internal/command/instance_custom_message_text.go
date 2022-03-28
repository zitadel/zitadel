package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
	"golang.org/x/text/language"
)

func (c *Commands) SetDefaultMessageText(ctx context.Context, instanceID string, messageText *domain.CustomMessageText) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(instanceID)
	events, existingMessageText, err := c.setDefaultMessageText(ctx, &instanceAgg.Aggregate, messageText)
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

func (c *Commands) setDefaultMessageText(ctx context.Context, instanceAgg *eventstore.Aggregate, msg *domain.CustomMessageText) ([]eventstore.Command, *InstanceCustomMessageTextWriteModel, error) {
	if !msg.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-kd9fs", "Errors.CustomMessageText.Invalid")
	}

	existingMessageText, err := c.defaultCustomMessageTextWriteModelByID(ctx, msg.MessageTextType, msg.Language)
	if err != nil {
		return nil, nil, err
	}

	events := make([]eventstore.Command, 0)
	if existingMessageText.Greeting != msg.Greeting {
		if msg.Greeting != "" {
			events = append(events, instance.NewCustomTextSetEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageGreeting, msg.Greeting, msg.Language))
		} else {
			events = append(events, instance.NewCustomTextRemovedEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageGreeting, msg.Language))
		}
	}
	if existingMessageText.Subject != msg.Subject {
		if msg.Subject != "" {
			events = append(events, instance.NewCustomTextSetEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageSubject, msg.Subject, msg.Language))
		} else {
			events = append(events, instance.NewCustomTextRemovedEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageSubject, msg.Language))
		}
	}
	if existingMessageText.Title != msg.Title {
		if msg.Title != "" {
			events = append(events, instance.NewCustomTextSetEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageTitle, msg.Title, msg.Language))
		} else {
			events = append(events, instance.NewCustomTextRemovedEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageTitle, msg.Language))
		}
	}
	if existingMessageText.PreHeader != msg.PreHeader {
		if msg.PreHeader != "" {
			events = append(events, instance.NewCustomTextSetEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessagePreHeader, msg.PreHeader, msg.Language))
		} else {
			events = append(events, instance.NewCustomTextRemovedEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessagePreHeader, msg.Language))
		}
	}
	if existingMessageText.Text != msg.Text {
		if msg.Text != "" {
			events = append(events, instance.NewCustomTextSetEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageText, msg.Text, msg.Language))
		} else {
			events = append(events, instance.NewCustomTextRemovedEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageText, msg.Language))
		}
	}
	if existingMessageText.ButtonText != msg.ButtonText {
		if msg.ButtonText != "" {
			events = append(events, instance.NewCustomTextSetEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageButtonText, msg.ButtonText, msg.Language))
		} else {
			events = append(events, instance.NewCustomTextRemovedEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageButtonText, msg.Language))
		}
	}
	if existingMessageText.FooterText != msg.FooterText {
		if msg.FooterText != "" {
			events = append(events, instance.NewCustomTextSetEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageFooterText, msg.FooterText, msg.Language))
		} else {
			events = append(events, instance.NewCustomTextRemovedEvent(ctx, instanceAgg, msg.MessageTextType, domain.MessageFooterText, msg.Language))
		}
	}
	return events, existingMessageText, nil
}

func (c *Commands) RemoveInstanceMessageTexts(ctx context.Context, messageTextType string, lang language.Tag) (*domain.ObjectDetails, error) {
	if messageTextType == "" || lang == language.Und {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-fjw9b", "Errors.CustomMessageText.Invalid")
	}
	customText, err := c.defaultCustomMessageTextWriteModelByID(ctx, messageTextType, lang)
	if err != nil {
		return nil, err
	}
	if customText.State == domain.PolicyStateUnspecified || customText.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-fju90", "Errors.CustomMessageText.NotFound")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&customText.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewCustomTextTemplateRemovedEvent(ctx, instanceAgg, messageTextType, lang))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(customText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&customText.WriteModel), nil
}

func (c *Commands) defaultCustomMessageTextWriteModelByID(ctx context.Context, messageType string, lang language.Tag) (*InstanceCustomMessageTextWriteModel, error) {
	writeModel := NewInstanceCustomMessageTextWriteModel(messageType, lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
