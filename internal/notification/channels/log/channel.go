package log

import (
	"fmt"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/notification/channels"
)

func InitStdoutChannel() channels.NotificationChannel {
	return channels.HandleMessageFunc(func(message channels.Message) error {
		logging.Log("NOTIF-c73ba").WithFields(map[string]interface{}{
			"type":    fmt.Sprintf("%T", message),
			"content": message.GetContent(),
		}).Info("handling notification message")
		return nil
	})
}
