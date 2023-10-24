package smtp_relay

import (
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
)

type Config struct {
	Connection         webhook.Config
	SMTP               smtp.Config
	IncludeSMTPMessage bool
}
