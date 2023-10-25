package notification

import (
	"context"
	"github.com/zitadel/zitadel/internal/notification/messages"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

var _ types.ChannelChains = (*allChannels)(nil)

type counters struct {
	success deliveryMetrics
	failed  deliveryMetrics
}

type deliveryMetrics struct {
	smtp         string
	emailWebhook string
	sms          string
	json         string
}

type allChannels struct {
	q        *handlers.NotificationQueries
	counters counters
}

func newChannels(q *handlers.NotificationQueries) *allChannels {
	c := &allChannels{
		q: q,
		counters: counters{
			success: deliveryMetrics{
				smtp:         "successful_deliveries_email",
				emailWebhook: "successful_deliveries_email_webhook",
				sms:          "successful_deliveries_sms",
				json:         "successful_deliveries_json",
			},
			failed: deliveryMetrics{
				smtp:         "failed_deliveries_email",
				emailWebhook: "failed_deliveries_email_webhook",
				sms:          "failed_deliveries_sms",
				json:         "failed_deliveries_json",
			},
		},
	}
	registerCounter(c.counters.success.smtp, "Successfully delivered emails over SMTP")
	registerCounter(c.counters.failed.smtp, "Failed email deliveries over SMTP")
	registerCounter(c.counters.success.emailWebhook, "Successfully delivered emails over webhook")
	registerCounter(c.counters.failed.emailWebhook, "Failed email deliveries over webhook")
	registerCounter(c.counters.success.sms, "Successfully delivered SMS")
	registerCounter(c.counters.failed.sms, "Failed SMS deliveries")
	registerCounter(c.counters.success.json, "Successfully delivered JSON messages over webhook")
	registerCounter(c.counters.failed.json, "Failed JSON message deliveries over webhook")
	return c
}

func registerCounter(counter, desc string) {
	err := metrics.RegisterCounter(counter, desc)
	logging.WithFields("metric", counter).OnError(err).Panic("unable to register counter")
}

func (c *allChannels) Email(ctx context.Context) (*senders.Chain[*messages.Email], *smtp.Config, error) {
	smtpCfg, err := c.q.GetSMTPConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	chain, err := senders.EmailChannels(
		ctx,
		c.q,
		smtpCfg,
		webhookCfg,
		c.counters.success.smtp,
		c.counters.failed.smtp,
		c.counters.success.emailWebhook,
		c.counters.failed.emailWebhook,
	)
	return chain, smtpCfg, err
}

func (c *allChannels) SMS(ctx context.Context) (*senders.Chain[*messages.SMS], *twilio.Config, error) {
	twilioCfg, err := c.q.GetTwilioConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	chain, err := senders.SMSChannels(ctx, c.q, twilioCfg, c.counters.success.sms, c.counters.failed.sms)
	return chain, twilioCfg, err
}

func (c *allChannels) Webhook(ctx context.Context, cfg webhook.Config) (*senders.Chain[*messages.JSON], error) {
	return senders.WebhookChannels(ctx, c.q, cfg, c.counters.success.json, c.counters.failed.json)
}
