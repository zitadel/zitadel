package types

import (
	"context"
	"html"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/eventstore"
	zchannels "github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func generateEmail(
	ctx context.Context,
	channels ChannelChains,
	config *email.Config,
	rule smtp.Rule,
	user *query.NotifyUser,
	template string,
	data templates.TemplateData,
	args map[string]interface{},
	lastEmail bool,
	triggeringEventType eventstore.EventType,
) error {
	recipient := user.VerifiedEmail
	if lastEmail {
		recipient = user.LastEmail
	}
	// The channels connect to the provider and are closed when the message is handled.
	// They are therefore only created once the email is ready to be sent.
	emailChannels, err := channels.Email(ctx, config)
	logging.OnError(ctx, err).Error("could not create email channel")
	if emailChannels == nil || emailChannels.Len() == 0 {
		return zchannels.NewCancelError(
			zerrors.ThrowPreconditionFailed(nil, "MAIL-w8nfow", "Errors.Notification.Channels.NotPresent"),
		)
	}
	if config.SMTPConfig != nil {
		message := &messages.Email{
			Recipients:          []string{recipient},
			Subject:             data.Subject,
			Content:             html.UnescapeString(template),
			Headers:             rule.Headers,
			TriggeringEventType: triggeringEventType,
		}
		return emailChannels.HandleMessage(message)
	}
	if config.WebhookConfig != nil {
		caseArgs := make(map[string]interface{}, len(args))
		for k, v := range args {
			caseArgs[strings.ToLower(string(k[0]))+k[1:]] = v
		}
		contextInfo := map[string]interface{}{
			"recipientEmailAddress": recipient,
			"eventType":             triggeringEventType,
			"provider":              config.ProviderConfig,
		}

		message := &messages.JSON{
			Serializable: &serializableData{
				ContextInfo:  contextInfo,
				TemplateData: data,
				Args:         caseArgs,
			},
			TriggeringEventType: triggeringEventType,
		}
		webhookChannels, err := channels.Webhook(ctx, *config.WebhookConfig)
		if err != nil {
			return err
		}
		return webhookChannels.HandleMessage(message)
	}
	return zchannels.NewCancelError(
		zerrors.ThrowPreconditionFailed(nil, "MAIL-83nof", "Errors.Notification.Channels.NotPresent"),
	)
}

func mapNotifyUserToArgs(user *query.NotifyUser, args map[string]interface{}) map[string]interface{} {
	if args == nil {
		args = make(map[string]interface{})
	}
	args["UserID"] = user.ID
	args["OrgID"] = user.ResourceOwner
	args["UserName"] = user.Username
	args["FirstName"] = user.FirstName
	args["LastName"] = user.LastName
	args["NickName"] = user.NickName
	args["DisplayName"] = user.DisplayName
	args["LastEmail"] = user.LastEmail
	args["VerifiedEmail"] = user.VerifiedEmail
	args["LastPhone"] = user.LastPhone
	args["VerifiedPhone"] = user.VerifiedPhone
	args["PreferredLoginName"] = user.PreferredLoginName
	args["LoginName"] = user.PreferredLoginName // some endpoint promoted LoginName instead of PreferredLoginName
	args["LoginNames"] = user.LoginNames
	args["ChangeDate"] = user.ChangeDate
	args["CreationDate"] = user.CreationDate
	return args
}
