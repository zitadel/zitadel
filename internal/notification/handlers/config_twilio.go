package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// GetTwilioConfig reads the iam Twilio provider config
func (n *NotificationQueries) GetTwilioConfig(ctx context.Context) (*twilio.Config, error) {
	active, err := query.NewSMSProviderStateQuery(domain.SMSConfigStateActive)
	if err != nil {
		return nil, err
	}
	config, err := n.SMSProviderConfig(ctx, active)
	if err != nil {
		return nil, err
	}
	if config.TwilioConfig == nil {
		return nil, zerrors.ThrowNotFound(nil, "HANDLER-8nfow", "Errors.SMS.Twilio.NotFound")
	}
	token, err := crypto.DecryptString(config.TwilioConfig.Token, n.SMSTokenCrypto)
	if err != nil {
		return nil, err
	}
	return &twilio.Config{
		SID:              config.TwilioConfig.SID,
		Token:            token,
		SenderNumber:     config.TwilioConfig.SenderNumber,
		VerifyServiceSID: config.TwilioConfig.VerifyServiceSID,
	}, nil
}
