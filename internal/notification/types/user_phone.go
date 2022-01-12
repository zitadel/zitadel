package types

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/messages"
	"github.com/caos/zitadel/internal/notification/senders"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

func generateSms(user *view_model.NotifyUser, content string, config systemdefaults.Notifications, lastPhone bool) error {
	message := &messages.SMS{
		SenderPhoneNumber:    config.Providers.Twilio.From,
		RecipientPhoneNumber: user.VerifiedPhone,
		Content:              content,
	}
	if lastPhone {
		message.RecipientPhoneNumber = user.LastPhone
	}

	channels, err := senders.SMSChannels(config)
	if err != nil {
		return err
	}
	return channels.HandleMessage(message)
}
