package notification

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
)

var _ types.ChannelChains = (*channels)(nil)

type channels struct {
	q *handlers.NotificationQueries
}

func (c *channels) Email(ctx context.Context) (*senders.Chain, *smtp.Config, error) {
	smtpCfg, err := c.q.GetSMTPConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	chain, err := senders.EmailChannels(
		ctx,
		smtpCfg,
		c.q.GetFileSystemProvider,
		c.q.GetLogProvider,
		metricSuccessfulDeliveriesEmail,
		metricFailedDeliveriesEmail,
	)
	return chain, smtpCfg, err
}

func (c *channels) SMS(ctx context.Context) (*senders.Chain, *twilio.Config, error) {
	twilioCfg, err := c.q.GetTwilioConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	chain, err := senders.SMSChannels(
		ctx,
		twilioCfg,
		c.q.GetFileSystemProvider,
		c.q.GetLogProvider,
		metricSuccessfulDeliveriesSMS,
		metricFailedDeliveriesSMS,
	)
	return chain, twilioCfg, err
}

func (c *channels) Webhook(ctx context.Context, cfg webhook.Config) (*senders.Chain, error) {
	return senders.WebhookChannels(
		ctx,
		cfg,
		c.q.GetFileSystemProvider,
		c.q.GetLogProvider,
		metricSuccessfulDeliveriesJSON,
		metricFailedDeliveriesJSON,
	)
}
