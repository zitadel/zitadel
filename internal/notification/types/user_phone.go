package types

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type serializableData struct {
	templates.TemplateData

	Args map[string]interface{}
}

func generateSms(
	ctx context.Context,
	channels ChannelChains,
	user *query.NotifyUser,
	data templates.TemplateData,
	args map[string]interface{},
	lastPhone bool,
	triggeringEvent eventstore.Event,
) error {
	smsChannels, config, err := channels.SMS(ctx)
	logging.OnError(err).Error("could not create sms channel")
	if smsChannels == nil || smsChannels.Len() == 0 {
		return zerrors.ThrowPreconditionFailed(nil, "PHONE-w8nfow", "Errors.Notification.Channels.NotPresent")
	}

	if config.TwilioConfig != nil {
		number := ""
		if err == nil {
			number = config.TwilioConfig.SenderNumber
		}
		message := &messages.SMS{
			SenderPhoneNumber:    number,
			RecipientPhoneNumber: user.VerifiedPhone,
			Content:              data.Text,
			TriggeringEvent:      triggeringEvent,
		}
		if lastPhone {
			message.RecipientPhoneNumber = user.LastPhone
		}
		return smsChannels.HandleMessage(message)
	}
	if config.WebhookConfig != nil {
		message := &messages.JSON{
			Serializable: &serializableData{
				TemplateData: data,
				Args:         args,
			},
			TriggeringEvent: triggeringEvent,
		}
		webhookChannels, err := channels.Webhook(ctx, *config.WebhookConfig)
		if err != nil {
			return err
		}
		return webhookChannels.HandleMessage(message)
	}
	return zerrors.ThrowPreconditionFailed(nil, "PHONE-w8nfow", "Errors.Notification.Channels.NotPresent")
}
