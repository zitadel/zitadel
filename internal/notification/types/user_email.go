package types

import (
	"context"
	"html"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/query"
)

func generateEmail(
	ctx context.Context,
	channels ChannelChains,
	user *query.NotifyUser,
	subject,
	content string,
	lastEmail bool,
	triggeringEvent eventstore.Event,
) error {
	content = html.UnescapeString(content)
	message := &messages.Email{
		Recipients:      []string{user.VerifiedEmail},
		Subject:         subject,
		Content:         content,
		TriggeringEvent: triggeringEvent,
	}
	if lastEmail {
		message.Recipients = []string{user.LastEmail}
	}
	emailChannels, _, err := channels.Email(ctx)
	if err != nil {
		return err
	}
	if emailChannels == nil || emailChannels.Len() == 0 {
		return errors.ThrowPreconditionFailed(nil, "MAIL-83nof", "Errors.Notification.Channels.NotPresent")
	}
	return emailChannels.HandleMessage(message)
}

func mapNotifyUserToArgs(user *query.NotifyUser, args map[string]interface{}) map[string]interface{} {
	if args == nil {
		args = make(map[string]interface{})
	}
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
	args["LoginNames"] = user.LoginNames
	args["ChangeDate"] = user.ChangeDate
	args["CreationDate"] = user.CreationDate
	return args
}
