package types

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/providers"
	"github.com/caos/zitadel/internal/notification/providers/chat"
	"github.com/caos/zitadel/internal/notification/providers/twilio"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

func generateSms(user *view_model.NotifyUser, content string, config systemdefaults.Notifications, lastPhone bool) error {
	provider := twilio.InitTwilioProvider(config.Providers.Twilio)
	message := &twilio.TwilioMessage{
		SenderPhoneNumber:    config.Providers.Twilio.From,
		RecipientPhoneNumber: user.VerifiedPhone,
		Content:              content,
	}
	if lastPhone {
		message.RecipientPhoneNumber = user.LastPhone
	}
	if provider.CanHandleMessage(message) {
		if config.DebugMode {
			return sendDebugPhone(message, config)
		}
		return provider.HandleMessage(message)
	}
	return caos_errs.ThrowInternalf(nil, "NOTIF-s8ipw", "Could not send init message: userid: %v", user.ID)
}

func sendDebugPhone(message providers.Message, config systemdefaults.Notifications) error {
	provider, err := chat.InitChatProvider(config.Providers.Chat)
	if err != nil {
		return err
	}
	return provider.HandleMessage(message)
}
