package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
)

// GetSMTPConfig reads the iam SMTP provider config
func (n *NotificationQueries) GetSMTPConfig(ctx context.Context) (*smtp.Config, error) {
	config, err := n.SMTPConfigByAggregateID(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	password, err := crypto.DecryptString(config.Password, n.SMTPPasswordCrypto)
	if err != nil {
		return nil, err
	}
	return &smtp.Config{
		From:           config.SenderAddress,
		FromName:       config.SenderName,
		ReplyToAddress: config.ReplyToAddress,
		Tls:            config.TLS,
		SMTP: smtp.SMTP{
			Host:     config.Host,
			User:     config.User,
			Password: password,
		},
	}, nil
}
