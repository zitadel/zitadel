package smtp_relay

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/templates"
)

type SMTPRelayMessage struct {
	Recipients     []string               `json:"recipients,omitempty"`
	BCC            []string               `json:"bcc,omitempty"`
	CC             []string               `json:"cc,omitempty"`
	SenderEmail    string                 `json:"senderEmail,omitempty"`
	SenderName     string                 `json:"senderName,omitempty"`
	ReplyToAddress string                 `json:"replyToAddress,omitempty"`
	Subject        string                 `json:"subject,omitempty"`
	SMTPContent    string                 `json:"content,omitempty"`
	TemplateData   templates.TemplateData `json:"templateData,omitempty"`
}

func newSMTPRelayMessage(email *messages.Email, includeSMTPContent bool) (msg SMTPRelayMessage, err error) {
	msg = SMTPRelayMessage{
		Recipients:     email.Recipients,
		BCC:            email.BCC,
		CC:             email.CC,
		SenderEmail:    email.SenderEmail,
		SenderName:     email.SenderName,
		ReplyToAddress: email.ReplyToAddress,
		Subject:        email.Subject,
		TemplateData:   email.TemplateData,
	}
	if includeSMTPContent {
		msg.SMTPContent, err = email.GetContent()
	}
	return msg, err
}

func InitChannel(ctx context.Context, cfg *Config) (channels.NotificationChannel[*messages.Email], error) {
	whChannel, err := webhook.InitChannel(ctx, cfg.Connection)
	if err != nil {
		return nil, err
	}
	return channels.HandleMessageFunc[*messages.Email](func(message *messages.Email) error {
		relayMsg, err := newSMTPRelayMessage(message, cfg.IncludeSMTPMessage)
		if err != nil {
			return err
		}
		return whChannel.HandleMessage(&messages.JSON{
			Serializable:    relayMsg,
			TriggeringEvent: message.GetTriggeringEvent(),
		})
	}), nil
}
