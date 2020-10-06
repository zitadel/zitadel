package types

import (
	"html"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/providers"
	"github.com/caos/zitadel/internal/notification/providers/chat"
	"github.com/caos/zitadel/internal/notification/providers/email"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

func generateEmail(user *view_model.NotifyUser, subject, content string, config systemdefaults.Notifications, lastEmail bool) error {
	provider, err := email.InitEmailProvider(config.Providers.Email)
	if err != nil {
		return err
	}
	content = html.UnescapeString(content)
	message := &email.EmailMessage{
		SenderEmail: config.Providers.Email.From,
		Recipients:  []string{user.VerifiedEmail},
		Subject:     subject,
		Content:     content,
	}
	if lastEmail {
		message.Recipients = []string{user.LastEmail}
	}
	if provider.CanHandleMessage(message) {
		if config.DebugMode {
			return sendDebugEmail(message, config)
		}
		return provider.HandleMessage(message)
	}
	return caos_errs.ThrowInternalf(nil, "NOTIF-s8ipw", "Could not send init message: userid: %v", user.ID)
}

func sendDebugEmail(message providers.Message, config systemdefaults.Notifications) error {
	provider, err := chat.InitChatProvider(config.Providers.Chat)
	if err != nil {
		return err
	}
	return provider.HandleMessage(message)
}
