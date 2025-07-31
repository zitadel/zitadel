package log

import (
	"fmt"

	"github.com/k3a/html2text"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/notification/channels"
)

func InitStdoutChannel(config Config) channels.NotificationChannel {

	logging.WithFields("logID", "NOTIF-D0164").Debug("successfully initialized stdout email and sms channel")

	return channels.HandleMessageFunc(func(message channels.Message) error {

		content, err := message.GetContent()
		if err != nil {
			return err
		}
		if config.Compact {
			content = html2text.HTML2Text(content)
		}

		logging.WithFields("logID", "NOTIF-c73ba").WithFields(map[string]any{
			"type":    fmt.Sprintf("%T", message),
			"content": content,
		}).Info("handling notification message")
		return nil
	})
}
