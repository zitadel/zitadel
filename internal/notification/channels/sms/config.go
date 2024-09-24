package sms

import (
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
)

type Config struct {
	ProviderConfig *Provider
	TwilioConfig   *twilio.Config
	WebhookConfig  *webhook.Config
}

type Provider struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

//
//func (config *Config) VerifyCode(phoneNumber, code string) error {
//	if config.TwilioConfig != nil {
//		return config.TwilioConfig.VerifyCode(phoneNumber, code)
//	}
//	//if config.WebhookConfig != nil {
//	//
//	//}
//	return zerrors.ThrowInvalidArgument(nil, "SMS-HKk31", "Errors.SMS.InvalidConfig")
//}
