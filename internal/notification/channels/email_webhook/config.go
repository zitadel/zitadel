package email_webhook

import (
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
)

type Config struct {
	Webhook                            webhook.Config
	IncludeContent, IncludeSMTPMessage bool
}
