package types

import (
	"context"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/notification/messages"
	"github.com/zitadel/zitadel/v2/internal/notification/templates"
	"github.com/zitadel/zitadel/v2/internal/query"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

type serializableData struct {
	ContextInfo  map[string]interface{} `json:"contextInfo,omitempty"`
	TemplateData templates.TemplateData `json:"templateData,omitempty"`
	Args         map[string]interface{} `json:"args,omitempty"`
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
	recipient := user.VerifiedPhone
	if lastPhone {
		recipient = user.LastPhone
	}
	if config.TwilioConfig != nil {
		number := ""
		if err == nil {
			number = config.TwilioConfig.SenderNumber
		}
		message := &messages.SMS{
			SenderPhoneNumber:    number,
			RecipientPhoneNumber: recipient,
			Content:              data.Text,
			TriggeringEvent:      triggeringEvent,
		}
		return smsChannels.HandleMessage(message)
	}
	if config.WebhookConfig != nil {
		caseArgs := make(map[string]interface{}, len(args))
		for k, v := range args {
			caseArgs[strings.ToLower(string(k[0]))+k[1:]] = v
		}
		contextInfo := map[string]interface{}{
			"recipientPhoneNumber": recipient,
			"eventType":            triggeringEvent.Type(),
			"provider":             config.ProviderConfig,
		}

		message := &messages.JSON{
			Serializable: &serializableData{
				TemplateData: data,
				Args:         caseArgs,
				ContextInfo:  contextInfo,
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
