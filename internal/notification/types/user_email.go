package types

import (
	"context"
	"html"

	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/query"
)

func generateEmail(ctx context.Context, user *query.NotifyUser, subject, content string, smtpConfig func(ctx context.Context) (*smtp.Config, error), getFileSystemProvider func(ctx context.Context) (*fs.Config, error), getLogProvider func(ctx context.Context) (*log.Config, error), lastEmail bool) error {
	content = html.UnescapeString(content)
	message := &messages.Email{
		Recipients: []string{user.VerifiedEmail},
		Subject:    subject,
		Content:    content,
	}
	if lastEmail {
		message.Recipients = []string{user.LastEmail}
	}

	channelChain, err := senders.EmailChannels(ctx, smtpConfig, getFileSystemProvider, getLogProvider)
	if err != nil {
		return err
	}

	if channelChain.Len() == 0 {
		return caos_errors.ThrowPreconditionFailed(nil, "MAIL-83nof", "Errors.Notification.Channels.NotPresent")
	}
	return channelChain.HandleMessage(message)
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
