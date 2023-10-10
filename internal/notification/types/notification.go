package types

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
)

type Notify func(
	url string,
	args map[string]interface{},
	messageType string,
	allowUnverifiedNotificationChannel bool,
) error

type ChannelChains interface {
	Email(context.Context) (*senders.Chain, *smtp.Config, error)
	SMS(context.Context) (*senders.Chain, *twilio.Config, error)
	Webhook(context.Context, webhook.Config) (*senders.Chain, error)
}

func SendEmail(
	ctx context.Context,
	channels ChannelChains,
	mailhtml string,
	translator *i18n.Translator,
	user *query.NotifyUser,
	colors *query.LabelPolicy,
	triggeringEvent eventstore.Event,
) Notify {
	return func(
		url string,
		args map[string]interface{},
		messageType string,
		allowUnverifiedNotificationChannel bool,
	) error {
		args = mapNotifyUserToArgs(user, args)
		data := GetTemplateData(ctx, translator, args, url, messageType, user.PreferredLanguage.String(), colors)
		template, err := templates.GetParsedTemplate(mailhtml, data)
		if err != nil {
			return err
		}
		return generateEmail(
			ctx,
			channels,
			user,
			data.Subject,
			template,
			allowUnverifiedNotificationChannel,
			triggeringEvent,
		)
	}
}

func SendSMSTwilio(
	ctx context.Context,
	channels ChannelChains,
	translator *i18n.Translator,
	user *query.NotifyUser,
	colors *query.LabelPolicy,
	triggeringEvent eventstore.Event,
) Notify {
	return func(
		url string,
		args map[string]interface{},
		messageType string,
		allowUnverifiedNotificationChannel bool,
	) error {
		args = mapNotifyUserToArgs(user, args)
		data := GetTemplateData(ctx, translator, args, url, messageType, user.PreferredLanguage.String(), colors)
		return generateSms(
			ctx,
			channels,
			user,
			data.Text,
			allowUnverifiedNotificationChannel,
			triggeringEvent,
		)
	}
}

func SendJSON(
	ctx context.Context,
	webhookConfig webhook.Config,
	channels ChannelChains,
	serializable interface{},
	triggeringEvent eventstore.Event,
) Notify {
	return func(_ string, _ map[string]interface{}, _ string, _ bool) error {
		return handleWebhook(
			ctx,
			webhookConfig,
			channels,
			serializable,
			triggeringEvent,
		)
	}
}
