package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

func AddEmailTemplate(
	a *instance.Aggregate,
	tempalte []byte,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			//TODO: check if already exists
			return []eventstore.Command{
				instance.NewMailTemplateAddedEvent(ctx, &a.Aggregate,
					tempalte,
				),
			}, nil
		}, nil
	}
}

func SetInstanceCustomTexts(
	a *instance.Aggregate,
	msg *domain.CustomMessageText,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			existing, err := existingInstanceCustomMessageText(ctx, filter, msg.MessageTextType, msg.Language)
			if err != nil {
				return nil, err
			}

			cmds := make([]eventstore.Command, 0, 7)
			if existing.Greeting != msg.Greeting {
				if msg.Greeting != "" {
					cmds = append(cmds, instance.NewCustomTextSetEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageGreeting, msg.Greeting, msg.Language))
				} else {
					cmds = append(cmds, instance.NewCustomTextRemovedEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageGreeting, msg.Language))
				}
			}
			if existing.Subject != msg.Subject {
				if msg.Subject != "" {
					cmds = append(cmds, instance.NewCustomTextSetEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageSubject, msg.Subject, msg.Language))
				} else {
					cmds = append(cmds, instance.NewCustomTextRemovedEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageSubject, msg.Language))
				}
			}
			if existing.Title != msg.Title {
				if msg.Title != "" {
					cmds = append(cmds, instance.NewCustomTextSetEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageTitle, msg.Title, msg.Language))
				} else {
					cmds = append(cmds, instance.NewCustomTextRemovedEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageTitle, msg.Language))
				}
			}
			if existing.PreHeader != msg.PreHeader {
				if msg.PreHeader != "" {
					cmds = append(cmds, instance.NewCustomTextSetEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessagePreHeader, msg.PreHeader, msg.Language))
				} else {
					cmds = append(cmds, instance.NewCustomTextRemovedEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessagePreHeader, msg.Language))
				}
			}
			if existing.Text != msg.Text {
				if msg.Text != "" {
					cmds = append(cmds, instance.NewCustomTextSetEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageText, msg.Text, msg.Language))
				} else {
					cmds = append(cmds, instance.NewCustomTextRemovedEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageText, msg.Language))
				}
			}
			if existing.ButtonText != msg.ButtonText {
				if msg.ButtonText != "" {
					cmds = append(cmds, instance.NewCustomTextSetEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageButtonText, msg.ButtonText, msg.Language))
				} else {
					cmds = append(cmds, instance.NewCustomTextRemovedEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageButtonText, msg.Language))
				}
			}
			if existing.FooterText != msg.FooterText {
				if msg.FooterText != "" {
					cmds = append(cmds, instance.NewCustomTextSetEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageFooterText, msg.FooterText, msg.Language))
				} else {
					cmds = append(cmds, instance.NewCustomTextRemovedEvent(ctx, &a.Aggregate, msg.MessageTextType, domain.MessageFooterText, msg.Language))
				}
			}
			// TODO: what if no text changed? len(events) == 0
			return cmds, nil
		}, nil
	}
}

func existingInstanceCustomMessageText(ctx context.Context, filter preparation.FilterToQueryReducer, textType string, lang language.Tag) (*command.InstanceCustomMessageTextWriteModel, error) {
	writeModel := command.NewInstanceCustomMessageTextWriteModel(textType, lang)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	writeModel.AppendEvents(events...)
	writeModel.Reduce()
	return writeModel, nil
}
