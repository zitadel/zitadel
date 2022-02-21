package types

import (
	"context"
	"html"

	"github.com/caos/zitadel/internal/notification/channels/smtp"
	"github.com/caos/zitadel/internal/notification/messages"
	"github.com/caos/zitadel/internal/notification/senders"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

func generateEmail(ctx context.Context, user *view_model.NotifyUser, subject, content string, config systemdefaults.Notifications, smtpConfig func(ctx context.Context) (*smtp.EmailConfig, error), lastEmail bool) error {
	content = html.UnescapeString(content)
	message := &messages.Email{
		Recipients: []string{user.VerifiedEmail},
		Subject:    subject,
		Content:    content,
	}
	if lastEmail {
		message.Recipients = []string{user.LastEmail}
	}

	channels, err := senders.EmailChannels(ctx, config, smtpConfig)
	if err != nil {
		return err
	}

	return channels.HandleMessage(message)
}

func mapNotifyUserToArgs(user *view_model.NotifyUser) map[string]interface{} {
	return map[string]interface{}{
		"UserName":           user.UserName,
		"FirstName":          user.FirstName,
		"LastName":           user.LastName,
		"NickName":           user.NickName,
		"DisplayName":        user.DisplayName,
		"LastEmail":          user.LastEmail,
		"VerifiedEmail":      user.VerifiedEmail,
		"LastPhone":          user.LastPhone,
		"VerifiedPhone":      user.VerifiedPhone,
		"PreferredLoginName": user.PreferredLoginName,
		"LoginNames":         user.LoginNames,
		"ChangeDate":         user.ChangeDate,
		"CreationDate":       user.CreationDate,
	}
}
