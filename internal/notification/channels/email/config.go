package email

import (
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
)

type Config struct {
	SMTPConfig    *smtp.Config
	WebhookConfig *webhook.Config
}
