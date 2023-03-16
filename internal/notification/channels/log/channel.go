package log

import (
	"context"
	"fmt"

	"github.com/k3a/html2text"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels"
)

func InitStdoutChannel(config LogConfig) channels.NotificationChannel {

	logging.Log("NOTIF-D0164").Debug("successfully initialized stdout email and sms channel")

	return channels.HandleMessageFunc(func(ctx context.Context, message channels.Message) error {

		content := message.GetContent()
		if config.Compact {
			content = html2text.HTML2Text(content)
		}

		logging.Log("NOTIF-c73ba").WithFields(map[string]interface{}{
			"instance": authz.GetInstance(ctx).InstanceID(),
			"type":     fmt.Sprintf("%T", message),
			"content":  content,
		}).Info("handling notification message")
		return nil
	})
}
