package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

func (c *Commands) SetDefaultMessageText(ctx context.Context, mailText *domain.CustomMessageText) (*domain.ObjectDetails, error) {
	iamAgg := iam.NewAggregate()
	events, existingMailText, err := c.setDefaultMessageText(ctx, &iamAgg.Aggregate, mailText)
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

func (c *Commands) setDefaultMessageText(ctx context.Context, iamAgg *eventstore.Aggregate, mailText *domain.CustomMessageText) ([]eventstore.EventPusher, *IAMCustomMessageTextReadModel, error) {
	//TODO: Check variablen
	if !mailText.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "IAM-kd9fs", "Errors.CustomText.Invalid")
	}

	existingMailText, err := c.defaultCustomMessageTextWriteModelByID(ctx)
	if err != nil {
		return nil, nil, err
	}
	events := make([]eventstore.EventPusher, 0)
	if existingMailText.Greeting != mailText.Greeting {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, mailText.MessageTextType+domain.MailGreeting, mailText.Greeting, mailText.Language))
	}
	if existingMailText.Subject != mailText.Subject {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, mailText.MessageTextType+domain.MailSubject, mailText.Subject, mailText.Language))
	}
	if existingMailText.Title != mailText.Title {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, mailText.MessageTextType+domain.MailTitle, mailText.Title, mailText.Language))
	}
	if existingMailText.PreHeader != mailText.PreHeader {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, mailText.MessageTextType+domain.MailPreHeader, mailText.PreHeader, mailText.Language))
	}
	if existingMailText.Text != mailText.Text {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, mailText.MessageTextType+domain.MailText, mailText.Text, mailText.Language))
	}
	if existingMailText.ButtonText != mailText.ButtonText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, mailText.MessageTextType+domain.MailButtonText, mailText.ButtonText, mailText.Language))
	}
	if existingMailText.FooterText != mailText.FooterText {
		events = append(events, iam.NewCustomTextSetEvent(ctx, iamAgg, mailText.MessageTextType+domain.MailFooterText, mailText.FooterText, mailText.Language))
	}
	return events, existingMailText, nil
}

func (c *Commands) defaultCustomMessageTextWriteModelByID(ctx context.Context) (*IAMCustomMessageTextReadModel, error) {
	writeModel := NewIAMCustomMessageTextWriteModel()
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
