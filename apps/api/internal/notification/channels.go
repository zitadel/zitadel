package notification

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/set"
	"github.com/zitadel/zitadel/internal/notification/channels/sms"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

var _ types.ChannelChains = (*channels)(nil)

type counters struct {
	success deliveryMetrics
	failed  deliveryMetrics
}

type deliveryMetrics struct {
	email string
	sms   string
	json  string
}

type channels struct {
	q        *handlers.NotificationQueries
	counters counters
}

func newChannels(q *handlers.NotificationQueries) *channels {
	c := &channels{
		q: q,
		counters: counters{
			success: deliveryMetrics{
				email: "successful_deliveries_email",
				sms:   "successful_deliveries_sms",
				json:  "successful_deliveries_json",
			},
			failed: deliveryMetrics{
				email: "failed_deliveries_email",
				sms:   "failed_deliveries_sms",
				json:  "failed_deliveries_json",
			},
		},
	}
	registerCounter(c.counters.success.email, "Successfully delivered emails")
	registerCounter(c.counters.failed.email, "Failed email deliveries")
	registerCounter(c.counters.success.sms, "Successfully delivered SMS")
	registerCounter(c.counters.failed.sms, "Failed SMS deliveries")
	registerCounter(c.counters.success.json, "Successfully delivered JSON messages")
	registerCounter(c.counters.failed.json, "Failed JSON message deliveries")
	return c
}

func registerCounter(counter, desc string) {
	err := metrics.RegisterCounter(counter, desc)
	logging.WithFields("metric", counter).OnError(err).Panic("unable to register counter")
}

func (c *channels) Email(ctx context.Context) (*senders.Chain, *email.Config, error) {
	emailCfg, err := c.q.GetActiveEmailConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	chain, err := senders.EmailChannels(
		ctx,
		emailCfg,
		c.q.GetFileSystemProvider,
		c.q.GetLogProvider,
		c.counters.success.email,
		c.counters.failed.email,
	)
	return chain, emailCfg, err
}

func (c *channels) SMS(ctx context.Context) (*senders.Chain, *sms.Config, error) {
	smsCfg, err := c.q.GetActiveSMSConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	chain, err := senders.SMSChannels(
		ctx,
		smsCfg,
		c.q.GetFileSystemProvider,
		c.q.GetLogProvider,
		c.counters.success.sms,
		c.counters.failed.sms,
	)
	return chain, smsCfg, err
}

func (c *channels) Webhook(ctx context.Context, cfg webhook.Config) (*senders.Chain, error) {
	return senders.WebhookChannels(
		ctx,
		cfg,
		c.q.GetFileSystemProvider,
		c.q.GetLogProvider,
		c.counters.success.json,
		c.counters.failed.json,
	)
}

func (c *channels) SecurityTokenEvent(ctx context.Context, cfg set.Config) (*senders.Chain, error) {
	return senders.SecurityEventTokenChannels(
		ctx,
		cfg,
		c.q.GetFileSystemProvider,
		c.q.GetLogProvider,
		c.counters.success.json,
		c.counters.failed.json,
	)
}
