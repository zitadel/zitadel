package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

func (c *Commands) SetDefaultMessageText(ctx context.Context, messageText *domain.CustomMessageText) (*domain.ObjectDetails, error) {
	iamAgg := iam.NewAggregate()
	events, existingMessageText, err := c.setDefaultMessageText(ctx, &iamAgg.Aggregate, messageText)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMessageText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingMessageText.WriteModel), nil
}

func (c *Commands) setDefaultMessageText(ctx context.Context, iamAgg *eventstore.Aggregate, msg *domain.CustomMessageText) ([]eventstore.EventPusher, *IAMCustomMessageTextReadModel, error) {
	if !msg.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "IAM-kd9fs", "Errors.CustomMessageText.Invalid")
	}

	existingMessageText, err := c.defaultCustomMessageTextWriteModelByID(ctx, msg.MessageTextType, msg.Language)
	if err != nil {
		return nil, nil, err
	}
	events := make([]eventstore.EventPusher, 0)
	if existingMessageText.Greeting != msg.Greeting {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageGreeting, msg.Greeting, msg.Language))
	}
	if existingMessageText.Subject != msg.Subject {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageSubject, msg.Subject, msg.Language))
	}
	if existingMessageText.Title != msg.Title {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageTitle, msg.Title, msg.Language))
	}
	if existingMessageText.PreHeader != msg.PreHeader {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessagePreHeader, msg.PreHeader, msg.Language))
	}
	if existingMessageText.Text != msg.Text {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageText, msg.Text, msg.Language))
	}
	if existingMessageText.ButtonText != msg.ButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageButtonText, msg.ButtonText, msg.Language))
	}
	if existingMessageText.FooterText != msg.FooterText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, msg.MessageTextType, domain.MessageFooterText, msg.FooterText, msg.Language))
	}
	return events, existingMessageText, nil
}

func (c *Commands) defaultCustomMessageTextWriteModelByID(ctx context.Context, messageType string, lang language.Tag) (*IAMCustomMessageTextReadModel, error) {
	writeModel := NewIAMCustomMessageTextWriteModel(messageType, lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
