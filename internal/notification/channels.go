package notification

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/metrics"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/set"
	"github.com/zitadel/zitadel/internal/notification/channels/sms"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
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
	q         *handlers.NotificationQueries
	counters  counters
	smtpRules smtp.Rules
}

func newChannels(q *handlers.NotificationQueries, smtpRules smtp.Rules) *channels {
	c := &channels{
		q:         q,
		smtpRules: smtpRules,
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
	registerCounter(types.SuppressedEmailsCounter, "Emails not sent to the provider because the recipient domain is reserved")
	registerCounter(types.RejectedNotificationsCounter, "Notifications intentionally not sent, e.g. because a limit of the provider was exceeded")
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

func (c *channels) EmailConfig(ctx context.Context) (*email.Config, error) {
	return c.q.GetActiveEmailConfig(ctx)
}

func (c *channels) Email(ctx context.Context, config *email.Config) (*senders.Chain, error) {
	return senders.EmailChannels(
		ctx,
		config,
		c.q.GetFileSystemProvider,
		c.q.GetLogProvider,
		c.counters.success.email,
		c.counters.failed.email,
	)
}

func (c *channels) SMTPRule(ctx context.Context, config *smtp.Config, orgID string) smtp.Rule {
	return c.smtpRules.Match(config, smtp.RuleData{
		InstanceID: authz.GetInstance(ctx).InstanceID(),
		OrgID:      orgID,
	})
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
