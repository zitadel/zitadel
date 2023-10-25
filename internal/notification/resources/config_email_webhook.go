package resources

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels/email_webhook"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"

	"github.com/zitadel/zitadel/internal/api/authz"
)

// GetEmailWebhookConfig reads the iam email webhook provider config
func (n *NotificationQueries) GetEmailWebhookConfig(ctx context.Context) (*email_webhook.Config, error) {
	_, err := n.EmailWebhookConfigByAggregateID(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	/*	_, err = crypto.DecryptString(config.Password, n.EmailWebhookHeaderCrypto)
		if err != nil {
			return nil, err
		}*/
	return &email_webhook.Config{
		Webhook:            webhook.Config{},
		IncludeContent:     false,
		IncludeSMTPMessage: false,
	}, nil
}
