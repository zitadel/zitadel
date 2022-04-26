package types

import (
	"html"

	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/senders"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

func generateEmail(user *view_model.NotifyUser, subject, content string, config systemdefaults.Notifications, lastEmail bool) error {
	content = html.UnescapeString(content)
	message := &messages.Email{
		SenderEmail: config.Providers.Email.From,
		Recipients:  []string{user.VerifiedEmail},
		Subject:     subject,
		Content:     content,
	}
	if lastEmail {
		message.Recipients = []string{user.LastEmail}
	}

	channels, err := senders.EmailChannels(config)
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
