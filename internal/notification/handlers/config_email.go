package handlers

import (
	"context"
	"net/http"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// GetSMTPConfig reads the iam SMTP provider config
func (n *NotificationQueries) GetActiveEmailConfig(ctx context.Context) (*email.Config, error) {
	config, err := n.SMTPConfigActive(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	provider := &email.Provider{
		ID:          config.ID,
		Description: config.Description,
	}
	if config.SMTPConfig != nil {
		if config.SMTPConfig.Password == nil {
			return nil, zerrors.ThrowNotFound(err, "QUERY-Wrs3gw", "Errors.SMTPConfig.NotFound")
		}
		password, err := crypto.DecryptString(config.SMTPConfig.Password, n.SMTPPasswordCrypto)
		if err != nil {
			return nil, err
		}
		return &email.Config{
			ProviderConfig: provider,
			SMTPConfig: &smtp.Config{
				From:           config.SMTPConfig.SenderAddress,
				FromName:       config.SMTPConfig.SenderName,
				ReplyToAddress: config.SMTPConfig.ReplyToAddress,
				Tls:            config.SMTPConfig.TLS,
				SMTP: smtp.SMTP{
					Host:     config.SMTPConfig.Host,
					User:     config.SMTPConfig.User,
					Password: password,
				},
			},
		}, nil
	}
	if config.HTTPConfig != nil {
		return &email.Config{
			ProviderConfig: provider,
			WebhookConfig: &webhook.Config{
				CallURL:    config.HTTPConfig.Endpoint,
				Method:     http.MethodPost,
				Headers:    nil,
				SigningKey: config.HTTPConfig.SigningKey,
			},
		}, nil
	}
	return nil, zerrors.ThrowNotFound(err, "QUERY-KPQleOckOV", "Errors.SMTPConfig.NotFound")
}
