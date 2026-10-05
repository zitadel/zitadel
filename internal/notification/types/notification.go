package types

import (
	"context"
	"html"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/i18n"
	zchannels "github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/set"
	"github.com/zitadel/zitadel/internal/notification/channels/sms"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Notify func(
	url string,
	args map[string]interface{},
	messageType string,
	allowUnverifiedNotificationChannel bool,
) error

type ChannelChains interface {
	// EmailConfig returns the active email provider of the instance.
	EmailConfig(context.Context) (*email.Config, error)
	// Email returns the channels of the email provider.
	// The channels connect to the provider, so they are only created right before an email is sent.
	Email(context.Context, *email.Config) (*senders.Chain, error)
	SMS(context.Context) (*senders.Chain, *sms.Config, error)
	Webhook(context.Context, webhook.Config) (*senders.Chain, error)
	SecurityTokenEvent(context.Context, set.Config) (*senders.Chain, error)
	// SMTPRule returns the rule defined by the operator for the SMTP provider.
	// If no rule matches the provider, the zero value is returned.
	SMTPRule(ctx context.Context, config *smtp.Config, orgID string) smtp.Rule
}

// SuppressedEmailsCounter counts the emails, which were not sent to the provider,
// because the recipient domain is reserved and the operator rule of the provider suppresses them.
const SuppressedEmailsCounter = "suppressed_deliveries_email"

// SendEmail returns a [Notify] function, which sends the notification as email.
// deliveryInfo is optional and states the provider the email was sent through,
// or that the email was accepted but not sent to the provider,
// because the recipient domain is reserved and the operator rule of the provider suppresses them.
func SendEmail(
	ctx context.Context,
	channels ChannelChains,
	mailhtml string,
	translator *i18n.Translator,
	user *query.NotifyUser,
	colors *query.LabelPolicy,
	triggeringEventType eventstore.EventType,
	deliveryInfo *senders.DeliveryInfo,
) Notify {
	return func(
		urlTmpl string,
		args map[string]interface{},
		messageType string,
		allowUnverifiedNotificationChannel bool,
	) error {
		// The provider is resolved before the content is rendered,
		// because the rule of the provider defines how it is rendered.
		// The channels are only created in [generateEmail], as they connect to the provider.
		config, err := channels.EmailConfig(ctx)
		if err != nil {
			logging.OnError(ctx, err).Error("could not get email provider")
			return zchannels.NewCancelError(
				zerrors.ThrowPreconditionFailed(err, "MAIL-w8nfow", "Errors.Notification.Channels.NotPresent"),
			)
		}
		rule := smtpRule(ctx, channels, config, user)
		args = mapNotifyUserToArgs(user, args)
		sanitizeArgsForHTML(args)
		url, err := urlFromTemplate(urlTmpl, args)
		if err != nil {
			return err
		}
		data := GetTemplateData(ctx, translator, args, url, messageType, user.PreferredLanguage.String(), colors)
		if rule.RestrictCustomHTML {
			// The texts contain the custom message texts of the instance / org with the already escaped arguments.
			// Only the texts are restricted, the arguments remain escaped and are never rendered.
			data.RestrictHTML()
		}
		template, err := templates.GetParsedTemplate(mailhtml, data)
		if err != nil {
			return err
		}
		return generateEmail(
			ctx,
			channels,
			config,
			rule,
			user,
			template,
			data,
			args,
			allowUnverifiedNotificationChannel,
			triggeringEventType,
			deliveryInfo,
		)
	}
}

// smtpRule returns the rule of the operator for the provider.
// Rules are only applied to SMTP providers.
func smtpRule(ctx context.Context, channels ChannelChains, config *email.Config, user *query.NotifyUser) smtp.Rule {
	if config == nil || config.SMTPConfig == nil {
		return smtp.Rule{}
	}
	return channels.SMTPRule(ctx, config.SMTPConfig, user.ResourceOwner)
}

func sanitizeArgsForHTML(args map[string]any) {
	for key, arg := range args {
		switch a := arg.(type) {
		case string:
			args[key] = html.EscapeString(a)
		case []string:
			for i, s := range a {
				a[i] = html.EscapeString(s)
			}
		case database.TextArray[string]:
			for i, s := range a {
				a[i] = html.EscapeString(s)
			}
		}
	}
}

func urlFromTemplate(urlTmpl string, args map[string]interface{}) (string, error) {
	var buf strings.Builder
	if err := domain.RenderURLTemplate(&buf, urlTmpl, args); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func SendSMS(
	ctx context.Context,
	channels ChannelChains,
	translator *i18n.Translator,
	user *query.NotifyUser,
	colors *query.LabelPolicy,
	triggeringEventType eventstore.EventType,
	instanceID string,
	jobID string,
	generatorInfo *senders.CodeGeneratorInfo,
) Notify {
	return func(
		urlTmpl string,
		args map[string]interface{},
		messageType string,
		allowUnverifiedNotificationChannel bool,
	) error {
		args = mapNotifyUserToArgs(user, args)
		url, err := urlFromTemplate(urlTmpl, args)
		if err != nil {
			return err
		}
		data := GetTemplateData(ctx, translator, args, url, messageType, user.PreferredLanguage.String(), colors)
		return generateSms(
			ctx,
			channels,
			user,
			data,
			args,
			allowUnverifiedNotificationChannel,
			triggeringEventType,
			instanceID,
			jobID,
			generatorInfo,
		)
	}
}

func SendJSON(
	ctx context.Context,
	webhookConfig webhook.Config,
	channels ChannelChains,
	serializable interface{},
	triggeringEventType eventstore.EventType,
) Notify {
	return func(_ string, _ map[string]interface{}, _ string, _ bool) error {
		return handleWebhook(
			ctx,
			webhookConfig,
			channels,
			serializable,
			triggeringEventType,
		)
	}
}

func SendSecurityTokenEvent(
	ctx context.Context,
	setConfig set.Config,
	channels ChannelChains,
	token any,
	triggeringEventType eventstore.EventType,
) Notify {
	return func(_ string, _ map[string]interface{}, _ string, _ bool) error {
		return handleSecurityTokenEvent(
			ctx,
			setConfig,
			channels,
			token,
			triggeringEventType,
		)
	}
}
