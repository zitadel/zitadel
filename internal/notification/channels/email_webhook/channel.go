package email_webhook

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/messages"
)

func Connect(ctx context.Context, cfg *Config) (channels.NotificationChannel[*messages.Email], error) {
	whChannel, err := webhook.Connect(ctx, cfg.Webhook)
	if err != nil {
		return nil, err
	}
	return channels.HandleMessageFunc[*messages.Email](func(message *messages.Email) error {
		emailJSON, err := message.ToJSON(cfg.IncludeContent, cfg.IncludeSMTPMessage)
		if err != nil {
			return err
		}
		return whChannel.HandleMessage(&messages.JSON{
			Serializable:    emailJSON,
			TriggeringEvent: message.GetTriggeringEvent(),
		})
	}), nil
}
