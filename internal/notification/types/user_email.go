package types

import (
	"context"
	"html"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func generateEmail(
	ctx context.Context,
	channels ChannelChains,
	user *query.NotifyUser,
	template string,
	data templates.TemplateData,
	args map[string]interface{},
	lastEmail bool,
	triggeringEvent eventstore.Event,
) error {
	emailChannels, config, err := channels.Email(ctx)
	logging.OnError(err).Error("could not create email channel")
	if emailChannels == nil || emailChannels.Len() == 0 {
		return zerrors.ThrowPreconditionFailed(nil, "PHONE-w8nfow", "Errors.Notification.Channels.NotPresent")
	}
	recipient := user.VerifiedEmail
	if lastEmail {
		recipient = user.LastEmail
	}
	if config.SMTPConfig != nil {
		message := &messages.Email{
			Recipients:      []string{recipient},
			Subject:         data.Subject,
			Content:         html.UnescapeString(template),
			TriggeringEvent: triggeringEvent,
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
			"eventType":             triggeringEvent.Type(),
			"provider":              config.ProviderConfig,
		}

		message := &messages.JSON{
			Serializable: &serializableData{
				ContextInfo:  contextInfo,
				TemplateData: data,
				Args:         caseArgs,
			},
			TriggeringEvent: triggeringEvent,
		}
		webhookChannels, err := channels.Webhook(ctx, *config.WebhookConfig)
		if err != nil {
			return err
		}
		return webhookChannels.HandleMessage(message)
	}
	return zerrors.ThrowPreconditionFailed(nil, "MAIL-83nof", "Errors.Notification.Channels.NotPresent")
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
