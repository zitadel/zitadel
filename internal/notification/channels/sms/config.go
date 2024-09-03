package sms

import (
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
)

type Config struct {
	TwilioConfig  *twilio.Config
	WebhookConfig *webhook.Config
}
